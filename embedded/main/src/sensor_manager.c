#include "sensor_manager.h"
#include "esp_adc/adc_oneshot.h"
#include "esp_adc/adc_cali.h"
#include "esp_adc/adc_cali_scheme.h"
#include "esp_timer.h"
#include "esp_log.h"
#include "esp_rom_sys.h"
#include "dht.h"
#include "adc.h"

#include "freertos/FreeRTOS.h"
#include "freertos/task.h"

static const char *TAG = "sensor_manager";
static adc_oneshot_config_t adc_config;
static sensor_history_t history = {0};

void sensor_manager_init(void)
{
    ESP_LOGI(TAG, "Initializing DHT22 sensor on GPIO %d", DHT22_PIN);
    
    // Initialize DHT22 sensor GPIO with pin from header
    setDHTgpio(DHT22_PIN);
    
    // Set initial state and wait for sensor to stabilize
    gpio_reset_pin(DHT22_PIN);
    gpio_set_direction(DHT22_PIN, GPIO_MODE_OUTPUT);
    gpio_set_level(DHT22_PIN, 1);
    vTaskDelay(pdMS_TO_TICKS(2000));  // 2 second delay for power stabilization

    // Initialize water sensor ADC
    adc_config = init_adc(WATER_SENSOR_ADC_UNIT, WATER_SENSOR_ADC_CHANNEL);
}

void read_dht22(float *temperature, float *humidity)
{
    int ret = readDHT();
    
    if (ret == DHT_OK) {
        *temperature = getTemperature();
        *humidity = getHumidity();
    } else {
        errorHandler(ret);
        *temperature = 0.0f;
        *humidity = 0.0f;
    }
}

int read_water_level(void)
{
    return adc_read_raw(adc_config, WATER_SENSOR_ADC_CHANNEL);
}

void get_current_sensor_data(sensor_data_t *data)
{
    float temp, humid;
    read_dht22(&temp, &humid);
    
    data->timestamp = esp_timer_get_time() / 1000; // Convert to milliseconds
    data->temperature = temp;
    data->humidity = humid;
    data->water_level = read_water_level();
        
    update_sensor_history(data);
}

void update_sensor_history(sensor_data_t *data)
{
    history.data[history.index] = *data;
    history.index = (history.index + 1) % QUEUE_SIZE;
    if (history.count < QUEUE_SIZE) {
        history.count++;
    }
}

void get_sensor_history(sensor_history_t *hist)
{
    memcpy(hist, &history, sizeof(sensor_history_t));
}