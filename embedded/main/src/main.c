#include <stdio.h>
#include "esp_system.h"
#include "nvs_flash.h"
#include "esp_event.h"
#include "esp_log.h"
#include "esp_timer.h"

#include "wifi_manager.h"
#include "sensor_manager.h"
#include "device_control.h"
#include "server_handlers.c"

void app_main(void)
{
    // Initialize NVS
    nvs_flash_init();
    sensor_manager_init();
    
    ESP_ERROR_CHECK(wifi_init_sta());   // Start WiFi
    ESP_ERROR_CHECK(server_init()); // Start HTTPS server

    ESP_LOGI(TAG, "System initialization complete");

    static bool was_water_low = false;  // Track previous water state to prevent rapid switching
    static uint32_t last_sensors_read_ms = 0;
    sensor_data_t current_data;

    // Main loop
    while (1) {
        uint32_t current_time = esp_timer_get_time() / 1000;  // Current time in milliseconds

        if (current_time - last_sensors_read_ms >= 2000) {  // Read sensors every 2 seconds
            get_current_sensor_data(&current_data);
            last_sensors_read_ms = current_time;

            // ESP_LOGI(TAG, "Temp: %.2f°C, Humidity: %.2f%%, Water Level: %d", 
            //          current_data.temperature, 
            //          current_data.humidity, 
            //          current_data.water_level);
        }

        // Auto mode logic
        if (device_state.auto_mode) {
            // Temperature control with hysteresis
            float temp = current_data.temperature;
            
            // Temperature control
            if (temp <= 28.0f) {
                // Cold condition: Turn on bulb, turn off fan
                device_state.bulb = true;
                device_state.fan = false;
                // ESP_LOGI(TAG, "Auto: Cold (%.2f°C) - Bulb ON, Fan OFF", temp);
            } else if (temp >= 32.0f) {
                // Hot condition: Turn on fan, turn off bulb
                device_state.bulb = false;
                device_state.fan = true;
                // ESP_LOGI(TAG, "Auto: Hot (%.2f°C) - Bulb OFF, Fan ON", temp);
            } else {
                // Comfortable range: Both off
                device_state.bulb = false;
                device_state.fan = false;
                // ESP_LOGI(TAG, "Auto: Normal (%.2f°C) - Both OFF", temp);
            }

            // Water level control with hysteresis
            // if (current_data.water_level <= WATER_LEVEL_LOW) {
            //     if (!was_water_low) {
            //         // Water level just became low
            //         device_states.pump = true;
            //         was_water_low = true;
            //         ESP_LOGI(TAG, "Auto: Water Low (%.2fV) - Pump ON", current_data.water_level);
            //     }
            // } else if (current_data.water_level >= WATER_LEVEL_FULL) {
            //     if (was_water_low) {
            //         // Water level reached full
            //         device_states.pump = false;
            //         was_water_low = false;
            //         ESP_LOGI(TAG, "Auto: Water Full (%.2fV) - Pump OFF", current_data.water_level);
            //     }
            // }

            // Update device states
            update_device_state(&device_state);
        }

        vTaskDelay(pdMS_TO_TICKS(100));  // Small delay to prevent watchdog timer resets
    }
}