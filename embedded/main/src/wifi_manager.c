#include "wifi_manager.h"
#include "esp_wifi.h"
#include "esp_log.h"
#include "esp_netif.h"
#include "mdns.h"
#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include <string.h>

static const char *TAG = "wifi_manager";

static void wifi_scan_and_log(void)
{
    wifi_scan_config_t scan_config = {
        .ssid = 0,
        .bssid = 0,
        .channel = 0,         // all channels
        .show_hidden = true,  // include hidden SSIDs
        .scan_type = WIFI_SCAN_TYPE_ACTIVE,
    };

    esp_err_t err = esp_wifi_scan_start(&scan_config, true /* block */);
    if (err != ESP_OK) {
        ESP_LOGW(TAG, "WiFi scan failed: %s", esp_err_to_name(err));
        return;
    }

    uint16_t ap_count = 0;
    esp_wifi_scan_get_ap_num(&ap_count);
    if (ap_count == 0) {
        ESP_LOGW(TAG, "WiFi scan found 0 APs");
        return;
    }

    uint16_t max_records = ap_count > 10 ? 10 : ap_count;
    wifi_ap_record_t *ap_records = calloc(max_records, sizeof(wifi_ap_record_t));
    if (!ap_records) {
        ESP_LOGW(TAG, "WiFi scan: OOM allocating records");
        return;
    }

    if (esp_wifi_scan_get_ap_records(&max_records, ap_records) != ESP_OK) {
        free(ap_records);
        ESP_LOGW(TAG, "WiFi scan get records failed");
        return;
    }

    ESP_LOGI(TAG, "WiFi scan found %d APs (showing %d):", ap_count, max_records);
    for (int i = 0; i < max_records; i++) {
        ESP_LOGI(TAG, "  #%d SSID='%s' RSSI=%d CH=%d AUTH=%d", i + 1,
                 (char *)ap_records[i].ssid, ap_records[i].rssi, ap_records[i].primary, ap_records[i].authmode);
    }

    free(ap_records);
}

static void wifi_scan_task(void *arg)
{
    (void)arg;
    wifi_scan_and_log();
    // After scan completes, attempt to connect using the configured SSID/PASS.
    esp_wifi_connect();
    vTaskDelete(NULL);
}

static void initialize_mdns(void)
{
    // Initialize mDNS
    ESP_ERROR_CHECK(mdns_init());
    // Set mDNS hostname (will be tokkatot-sensor.local)
    ESP_ERROR_CHECK(mdns_hostname_set("tokkatot-sensor"));
    // Set default instance name
    ESP_ERROR_CHECK(mdns_instance_name_set("Tokkatot Sensor Node"));

    // Add HTTP service
    mdns_service_add("tokkatot-sensor", "_http", "_tcp", 80, NULL, 0);
    ESP_LOGI(TAG, "mDNS initialized: tokkatot-sensor.local");
}

static void wifi_event_handler(void* arg, esp_event_base_t event_base,
                             int32_t event_id, void* event_data)
{
    if (event_base == WIFI_EVENT && event_id == WIFI_EVENT_STA_START) {
        ESP_LOGI(TAG, "WiFi STA start, connecting...");
        // Don't scan inside the system event task; it can overflow stack and block WiFi events.
        // Also don't call connect before scanning, otherwise scanning is not allowed.
        xTaskCreate(wifi_scan_task, "wifi_scan", 4096, NULL, tskIDLE_PRIORITY + 1, NULL);
    } else if (event_base == WIFI_EVENT && event_id == WIFI_EVENT_STA_DISCONNECTED) {
        wifi_event_sta_disconnected_t *disc = (wifi_event_sta_disconnected_t *)event_data;
        ESP_LOGW(TAG, "WiFi disconnected (reason=%d). Reconnecting...", disc ? disc->reason : -1);
        esp_wifi_connect();
    } else if (event_base == IP_EVENT && event_id == IP_EVENT_STA_GOT_IP) {
        ip_event_got_ip_t* event = (ip_event_got_ip_t*) event_data;
        ESP_LOGI(TAG, "Got IP:" IPSTR, IP2STR(&event->ip_info.ip));
    }
}

esp_err_t wifi_init_sta(void)
{
    ESP_ERROR_CHECK(esp_netif_init());
    ESP_ERROR_CHECK(esp_event_loop_create_default());
    esp_netif_t *sta_netif = esp_netif_create_default_wifi_sta();

    wifi_init_config_t cfg = WIFI_INIT_CONFIG_DEFAULT();
    ESP_ERROR_CHECK(esp_wifi_init(&cfg));

    ESP_ERROR_CHECK(esp_event_handler_instance_register(WIFI_EVENT,
                                                      ESP_EVENT_ANY_ID,
                                                      &wifi_event_handler,
                                                      NULL,
                                                      NULL));
    ESP_ERROR_CHECK(esp_event_handler_instance_register(IP_EVENT,
                                                      IP_EVENT_STA_GOT_IP,
                                                      &wifi_event_handler,
                                                      NULL,
                                                      NULL));

    wifi_config_t wifi_config = {
        .sta = {
            .ssid = WIFI_SSID,
            .password = WIFI_PASS,
            // Use a permissive minimum auth mode so WPA2-only hotspots are not filtered out.
            // (Filtering too high can produce reason=211: NO_AP_FOUND_IN_AUTHMODE_THRESHOLD.)
            .threshold.authmode = WIFI_AUTH_WPA2_PSK,
        },
    };

    ESP_LOGI(TAG, "Configuring WiFi SSID: '%s'", WIFI_SSID);

    ESP_ERROR_CHECK(esp_wifi_set_mode(WIFI_MODE_STA));
    ESP_ERROR_CHECK(esp_wifi_set_config(WIFI_IF_STA, &wifi_config));
    ESP_ERROR_CHECK(esp_wifi_start());

    // Initialize discovery
    initialize_mdns();

    ESP_LOGI(TAG, "wifi_init_sta finished.");
    return ESP_OK;
}
