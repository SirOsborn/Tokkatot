#ifndef DEVICE_CONTROL_H
#define DEVICE_CONTROL_H

#include <stdbool.h>
#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "driver/gpio.h"
#include "driver/ledc.h"

// Pin definitions
#define CONVEYER_PIN     GPIO_NUM_25
#define SERVO_PIN        GPIO_NUM_23
#define FAN_PIN         GPIO_NUM_26
#define LIGHTBULB_PIN   GPIO_NUM_14
#define WATERPUMP_PIN   GPIO_NUM_27

// Device states
typedef struct {
    bool auto_mode;
    bool fan;
    bool bulb;
    bool feeder;
    bool pump;
    bool conveyer;
} device_state_t;

// Function declarations
void device_control_init(void);
void open_feeder(void);
void close_feeder(void);
void toggle_device(gpio_num_t pin, bool *state);
void update_device_state(device_state_t *state);

#endif // DEVICE_CONTROL_H