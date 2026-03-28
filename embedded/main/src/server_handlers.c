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
    cJSON_AddNumberToObject(root, "heater", device_state.heater);
    cJSON_AddNumberToObject(root, "feeder_motor", device_state.feeder_motor);
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
    cJSON_AddNumberToObject(root, "water_level", current_data.water_level);

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
        cJSON_AddNumberToObject(entry, "water_level", history.data[idx].water_level);
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
    device_state.heater = false;
    device_state.fan = false;
    device_state.feeder_motor = false;
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

static esp_err_t toggle_heater_handler(httpd_req_t *req)
{
    toggle_device(HEATER_PIN, &device_state.heater);
    update_device_state(&device_state);
    return send_text_response(req, device_state.heater ? "true" : "false");
}

static esp_err_t toggle_feeder_motor_handler(httpd_req_t *req)
{
    toggle_device(FEEDER_MOTOR_PIN, &device_state.feeder_motor);
    update_device_state(&device_state);
    return send_text_response(req, device_state.feeder_motor ? "true" : "false");
}

static esp_err_t parse_state_body(httpd_req_t *req, bool *state_out, uint32_t *duration_out)
{
    int len = req->content_len;
    if (len <= 0 || len > 256) {
        return ESP_FAIL;
    }

    char buf[257];
    int ret = httpd_req_recv(req, buf, len);
    if (ret <= 0) {
        return ESP_FAIL;
    }
    buf[ret] = '\0';

    cJSON *root = cJSON_Parse(buf);
    if (!root) {
        return ESP_FAIL;
    }

    cJSON *state = cJSON_GetObjectItem(root, "state");
    if (!cJSON_IsBool(state)) {
        cJSON_Delete(root);
        return ESP_FAIL;
    }

    *state_out = cJSON_IsTrue(state);
    
    // Optional duration parsing
    cJSON *duration = cJSON_GetObjectItem(root, "duration");
    if (cJSON_IsNumber(duration)) {
        *duration_out = (uint32_t)duration->valuedouble;
    } else {
        *duration_out = 0;
    }

    cJSON_Delete(root);
    return ESP_OK;
}

static esp_err_t set_fan_handler(httpd_req_t *req)
{
    bool on = false;
    uint32_t duration = 0;
    if (parse_state_body(req, &on, &duration) != ESP_OK) {
        httpd_resp_send_err(req, HTTPD_400_BAD_REQUEST, "Invalid body");
        return ESP_FAIL;
    }
    set_device_timed(FAN_PIN, &device_state.fan, on, duration);
    update_device_state(&device_state);
    return send_text_response(req, device_state.fan ? "true" : "false");
}

static esp_err_t set_heater_handler(httpd_req_t *req)
{
    bool on = false;
    uint32_t duration = 0;
    if (parse_state_body(req, &on, &duration) != ESP_OK) {
        httpd_resp_send_err(req, HTTPD_400_BAD_REQUEST, "Invalid body");
        return ESP_FAIL;
    }
    set_device_timed(HEATER_PIN, &device_state.heater, on, duration);
    update_device_state(&device_state);
    return send_text_response(req, device_state.heater ? "true" : "false");
}

static esp_err_t set_feeder_motor_handler(httpd_req_t *req)
{
    bool on = false;
    uint32_t duration = 0;
    if (parse_state_body(req, &on, &duration) != ESP_OK) {
        httpd_resp_send_err(req, HTTPD_400_BAD_REQUEST, "Invalid body");
        return ESP_FAIL;
    }
    set_device_timed(FEEDER_MOTOR_PIN, &device_state.feeder_motor, on, duration);
    update_device_state(&device_state);
    return send_text_response(req, device_state.feeder_motor ? "true" : "false");
}

static esp_err_t set_conveyer_handler(httpd_req_t *req)
{
    bool on = false;
    uint32_t duration = 0;
    if (parse_state_body(req, &on, &duration) != ESP_OK) {
        httpd_resp_send_err(req, HTTPD_400_BAD_REQUEST, "Invalid body");
        return ESP_FAIL;
    }
    set_device_timed(CONVEYER_PIN, &device_state.conveyer, on, duration);
    update_device_state(&device_state);
    return send_text_response(req, device_state.conveyer ? "true" : "false");
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

    httpd_uri_t uri_toggle_heater = {
        .uri = "/toggle-heater",
        .method = HTTP_GET,
        .handler = toggle_heater_handler
    };
    httpd_register_uri_handler(server, &uri_toggle_heater);

    httpd_uri_t uri_toggle_feeder = {
        .uri = "/toggle-feeder",
        .method = HTTP_GET,
        .handler = toggle_feeder_motor_handler
    };
    httpd_register_uri_handler(server, &uri_toggle_feeder);

    /* Register explicit actuator set endpoints (POST with {"state":true/false}) */
    httpd_uri_t uri_set_fan = {
        .uri = "/actuators/fan",
        .method = HTTP_POST,
        .handler = set_fan_handler
    };
    httpd_register_uri_handler(server, &uri_set_fan);

    httpd_uri_t uri_set_heater = {
        .uri = "/actuators/heater",
        .method = HTTP_POST,
        .handler = set_heater_handler
    };
    httpd_register_uri_handler(server, &uri_set_heater);

    httpd_uri_t uri_set_feeder = {
        .uri = "/actuators/feeder_motor",
        .method = HTTP_POST,
        .handler = set_feeder_motor_handler
    };
    httpd_register_uri_handler(server, &uri_set_feeder);

    httpd_uri_t uri_set_conveyer = {
        .uri = "/actuators/conveyor_belt",
        .method = HTTP_POST,
        .handler = set_conveyer_handler
    };
    httpd_register_uri_handler(server, &uri_set_conveyer);

    ESP_LOGI(TAG, "Server started and URIs registered");
    return ESP_OK;
}
