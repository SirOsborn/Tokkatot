#ifndef DEVICE_CONTROL_H
#define DEVICE_CONTROL_H

#include <stdbool.h>
#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "driver/gpio.h"

// Pin definitions
#define CONVEYER_PIN     GPIO_NUM_25
#define FAN_PIN          GPIO_NUM_26
#define HEATER_PIN       GPIO_NUM_14
#define FEEDER_MOTOR_PIN GPIO_NUM_27

// Device states
typedef struct {
    bool auto_mode;
    bool fan;
    bool heater;
    bool feeder_motor;
    bool conveyer;
} device_state_t;

// Function declarations
void device_control_init(void);
void toggle_device(gpio_num_t pin, bool *state);
void set_device(gpio_num_t pin, bool *state, bool on);
void set_device_timed(gpio_num_t pin, bool *state, bool on, uint32_t duration_sec);
void update_device_state(device_state_t *state);

#endif // DEVICE_CONTROL_H
