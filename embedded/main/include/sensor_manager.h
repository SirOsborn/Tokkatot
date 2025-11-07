#ifndef SENSOR_MANAGER_H
#define SENSOR_MANAGER_H

#include "driver/gpio.h"

#include <stdint.h>
//#include "esp_err.h"

#define DHT22_PIN       GPIO_NUM_32
// GPIO 35
#define WATER_SENSOR_ADC_UNIT ADC_UNIT_1
#define WATER_SENSOR_ADC_CHANNEL ADC_CHANNEL_7

// Water level thresholds
#define WATER_LEVEL_LOW    300
#define WATER_LEVEL_FULL   2000

#define QUEUE_SIZE 10

typedef struct {
    uint64_t timestamp;
    float temperature;
    float humidity;
    int water_level;
} sensor_data_t;

typedef struct {
    sensor_data_t data[QUEUE_SIZE];
    int index;
    int count;
} sensor_history_t;

void sensor_manager_init(void);
void read_dht22(float *temperature, float *humidity);
void get_current_sensor_data(sensor_data_t *data);
void update_sensor_history(sensor_data_t *data);
void get_sensor_history(sensor_history_t *history);
void read_dht22(float *temperature, float *humidity);
int read_water_level(void);

#endif // SENSOR_MANAGER_H