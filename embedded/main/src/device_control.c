#include "device_control.h"
#include "esp_log.h"
#include "iot_servo.h"

static const char *tag = "device_control";

// Servo calibration angles
static uint16_t CALIBRATION_0 = 30;    // Real 0 degree angle
static uint16_t CALIBRATION_180 = 195; // Real 180 degree angle

static device_state_t device_state = {
    .auto_mode = false,
    .fan = false,
    .bulb = false,
    .feeder = false,
    .pump = false,
    .conveyer = false
};

void device_control_init(void)
{
    // Configure GPIO pins
    gpio_config_t io_conf = {
        .intr_type = GPIO_INTR_DISABLE,
        .mode = GPIO_MODE_OUTPUT,
        .pin_bit_mask = (1ULL << CONVEYER_PIN) |
                       (1ULL << FAN_PIN) |
                       (1ULL << LIGHTBULB_PIN) |
                       (1ULL << WATERPUMP_PIN),
        .pull_down_en = 0,
        .pull_up_en = 0
    };
    gpio_config(&io_conf);

    servo_config_t servo_cfg = {
        .max_angle = 180,
        .min_width_us = 500,
        .max_width_us = 2500,
        .freq = 50,
        .timer_number = LEDC_TIMER_0,
        .channels = {
            .servo_pin = {
                SERVO_PIN
            },
            .ch = {
                LEDC_CHANNEL_0,
            },
        },
        .channel_number = 1,
    };

    // Initialize the servo
    iot_servo_init(LEDC_LOW_SPEED_MODE, &servo_cfg);

    // Initialize all devices to OFF state
    gpio_set_level(CONVEYER_PIN, 1);
    gpio_set_level(FAN_PIN, 1);
    gpio_set_level(LIGHTBULB_PIN, 1);
    gpio_set_level(WATERPUMP_PIN, 1);
    iot_servo_write_angle(LEDC_LOW_SPEED_MODE, 0, CALIBRATION_0 + 60);
}

void close_feeder(void)
{
    for (uint16_t i = CALIBRATION_0; i <= (CALIBRATION_0 + 60); i++) {
        iot_servo_write_angle(LEDC_LOW_SPEED_MODE, 0, i);
        vTaskDelay(20 / portTICK_PERIOD_MS);
    }
}

void open_feeder(void)
{
    for (uint16_t i = (CALIBRATION_0 + 60); i >= CALIBRATION_0; i--) {
        iot_servo_write_angle(LEDC_LOW_SPEED_MODE, 0, i);
        vTaskDelay(20 / portTICK_PERIOD_MS);
    }
}

void toggle_device(gpio_num_t pin, bool *state)
{
    *state = !*state;
    gpio_set_level(pin, *state ? 0 : 1);
    ESP_LOGI(tag, "Toggled pin %d to %s", pin, *state ? "ON" : "OFF");
}

void update_device_state(device_state_t *state)
{
    memcpy(&device_state, state, sizeof(device_state_t));
    
    gpio_set_level(FAN_PIN, state->fan ? 0 : 1);
    gpio_set_level(LIGHTBULB_PIN, state->bulb ? 0 : 1);
    gpio_set_level(WATERPUMP_PIN, state->pump ? 0 : 1);
    gpio_set_level(CONVEYER_PIN, state->conveyer ? 0 : 1);
}