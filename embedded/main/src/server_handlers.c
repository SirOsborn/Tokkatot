#include "server_handlers.h"
#include "esp_log.h"
#include "esp_https_server.h"
#include "sensor_manager.h"
#include "device_control.h"
#include "cJSON.h"
#include <string.h>

static const char *TAG = "server_handlers";
static httpd_handle_t server = NULL;

/* SSL certificate - to be replaced with the generated certificate */
extern const unsigned char servercert_start[] asm("_binary_cert_pem_start");
extern const unsigned char servercert_end[] asm("_binary_cert_pem_end");
extern const unsigned char serverkey_start[] asm("_binary_key_pem_start");
extern const unsigned char serverkey_end[] asm("_binary_key_pem_end");

/* Local copy of device states managed by device_control */
static device_state_t device_state;

/* Helper: send JSON response (already present) */
esp_err_t send_json_response(httpd_req_t *req, cJSON *root)
{
    char *json_str = cJSON_Print(root);
    httpd_resp_set_type(req, "application/json");
    httpd_resp_sendstr(req, json_str);
    free(json_str);
    cJSON_Delete(root);
    return ESP_OK;
}

/* Helper: send plain text response (used for toggles) */
static esp_err_t send_text_response(httpd_req_t *req, const char *text)
{
    httpd_resp_set_type(req, "text/plain");
    httpd_resp_sendstr(req, text);
    return ESP_OK;
}

/* ====== DATA HANDLERS ====== */
esp_err_t get_initial_state_handler(httpd_req_t *req)
{
    cJSON *root = cJSON_CreateObject();
    cJSON_AddNumberToObject(root, "auto_mode", device_state.auto_mode);
    cJSON_AddNumberToObject(root, "fan", device_state.fan);
    cJSON_AddNumberToObject(root, "bulb", device_state.bulb);
    cJSON_AddNumberToObject(root, "feeder", device_state.feeder);
    cJSON_AddNumberToObject(root, "pump", device_state.pump);
    cJSON_AddNumberToObject(root, "conveyer", device_state.conveyer);

    return send_json_response(req, root);
}

esp_err_t get_current_data_handler(httpd_req_t *req)
{
    sensor_data_t current_data;
    get_current_sensor_data(&current_data);

    cJSON *root = cJSON_CreateObject();
    cJSON_AddNumberToObject(root, "timestamp", current_data.timestamp);
    cJSON_AddNumberToObject(root, "temperature", current_data.temperature);
    cJSON_AddNumberToObject(root, "humidity", current_data.humidity);

    return send_json_response(req, root);
}

esp_err_t get_historical_data_handler(httpd_req_t *req)
{
    sensor_history_t history;
    get_sensor_history(&history);

    cJSON *root = cJSON_CreateArray();

    for (int i = 0; i < history.count; i++) {
        int idx = (history.index - history.count + i + QUEUE_SIZE) % QUEUE_SIZE;
        cJSON *entry = cJSON_CreateObject();
        cJSON_AddNumberToObject(entry, "timestamp", history.data[idx].timestamp);
        cJSON_AddNumberToObject(entry, "temperature", history.data[idx].temperature);
        cJSON_AddNumberToObject(entry, "humidity", history.data[idx].humidity);
        cJSON_AddItemToArray(root, entry);
    }

    return send_json_response(req, root);
}

/* ====== TOGGLE HANDLERS ======
   Each toggle handler flips the relevant state and returns plain text "true" or "false"
   so the upstream API can consume the raw response body directly. */

static esp_err_t toggle_auto_handler(httpd_req_t *req)
{
    device_state.auto_mode = !device_state.auto_mode;
    device_state.bulb = false;
    device_state.fan = false;
    device_state.pump = false;
    device_state.conveyer = false;
    update_device_state(&device_state);

    return send_text_response(req, device_state.auto_mode ? "true" : "false");
}

static esp_err_t toggle_belt_handler(httpd_req_t *req)
{
    // Toggle conveyer (belt)
    toggle_device(CONVEYER_PIN, &device_state.conveyer);
    update_device_state(&device_state);
    return send_text_response(req, device_state.conveyer ? "true" : "false");
}

static esp_err_t toggle_fan_handler(httpd_req_t *req)
{
    toggle_device(FAN_PIN, &device_state.fan);
    update_device_state(&device_state);
    return send_text_response(req, device_state.fan ? "true" : "false");
}

static esp_err_t toggle_bulb_handler(httpd_req_t *req)
{
    toggle_device(LIGHTBULB_PIN, &device_state.bulb);
    update_device_state(&device_state);
    return send_text_response(req, device_state.bulb ? "true" : "false");
}

static esp_err_t toggle_pump_handler(httpd_req_t *req)
{
    toggle_device(WATERPUMP_PIN, &device_state.pump);
    update_device_state(&device_state);
    return send_text_response(req, device_state.pump ? "true" : "false");
}

static esp_err_t toggle_feeder_handler(httpd_req_t *req)
{
    device_state.feeder = !device_state.feeder;
    device_state.feeder ? open_feeder() : close_feeder();
    update_device_state(&device_state);
    return send_text_response(req, device_state.feeder ? "true" : "false");
}

/* Server initialization: start HTTPS server, register data and toggle endpoints */
esp_err_t server_init(void)
{
    httpd_ssl_config_t config = HTTPD_SSL_CONFIG_DEFAULT();
    config.httpd.max_uri_handlers = 12; // Adjust based on number of handlers

    config.servercert = servercert_start;
    config.servercert_len = servercert_end - servercert_start;
    config.prvtkey_pem = serverkey_start;
    config.prvtkey_len = serverkey_end - serverkey_start;

    /* Initialize device control and read current state */
    device_control_init();
    memset(&device_state, 0, sizeof(device_state));
    update_device_state(&device_state);

    esp_err_t ret = httpd_ssl_start(&server, &config);
    if (ret != ESP_OK) {
        ESP_LOGE(TAG, "Failed to start server!");
        return ret;
    }

    /* Register data endpoints */
    httpd_uri_t initial_state = {
        .uri = "/get-initial-state",
        .method = HTTP_GET,
        .handler = get_initial_state_handler
    };
    httpd_register_uri_handler(server, &initial_state);

    httpd_uri_t current_data = {
        .uri = "/get-current-data",
        .method = HTTP_GET,
        .handler = get_current_data_handler
    };
    httpd_register_uri_handler(server, &current_data);

    httpd_uri_t historical_data = {
        .uri = "/get-historical-data",
        .method = HTTP_GET,
        .handler = get_historical_data_handler
    };
    httpd_register_uri_handler(server, &historical_data);

    /* Register toggle endpoints (no verify step; return plain "true"/"false") */
    httpd_uri_t uri_toggle_auto = {
        .uri = "/toggle-auto",
        .method = HTTP_GET,
        .handler = toggle_auto_handler
    };
    httpd_register_uri_handler(server, &uri_toggle_auto);

    httpd_uri_t uri_toggle_belt = {
        .uri = "/toggle-belt",
        .method = HTTP_GET,
        .handler = toggle_belt_handler
    };
    httpd_register_uri_handler(server, &uri_toggle_belt);

    httpd_uri_t uri_toggle_fan = {
        .uri = "/toggle-fan",
        .method = HTTP_GET,
        .handler = toggle_fan_handler
    };
    httpd_register_uri_handler(server, &uri_toggle_fan);

    httpd_uri_t uri_toggle_bulb = {
        .uri = "/toggle-bulb",
        .method = HTTP_GET,
        .handler = toggle_bulb_handler
    };
    httpd_register_uri_handler(server, &uri_toggle_bulb);

    httpd_uri_t uri_toggle_pump = {
        .uri = "/toggle-pump",
        .method = HTTP_GET,
        .handler = toggle_pump_handler
    };
    httpd_register_uri_handler(server, &uri_toggle_pump);

    httpd_uri_t uri_toggle_feeder = {
        .uri = "/toggle-feeder",
        .method = HTTP_GET,
        .handler = toggle_feeder_handler
    };
    httpd_register_uri_handler(server, &uri_toggle_feeder);

    ESP_LOGI(TAG, "Server started and URIs registered");
    return ESP_OK;
}
