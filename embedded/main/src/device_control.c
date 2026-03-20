#include "device_control.h"
#include "esp_log.h"
static const char *tag = "device_control";

static device_state_t device_state = {
    .auto_mode = false,
    .fan = false,
    .heater = false,
    .feeder_motor = false,
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
                       (1ULL << HEATER_PIN) |
                       (1ULL << FEEDER_MOTOR_PIN),
        .pull_down_en = 0,
        .pull_up_en = 0
    };
    gpio_config(&io_conf);

    // Initialize all devices to OFF state
    gpio_set_level(CONVEYER_PIN, 1);
    gpio_set_level(FAN_PIN, 1);
    gpio_set_level(HEATER_PIN, 1);
    gpio_set_level(FEEDER_MOTOR_PIN, 1);
}

void toggle_device(gpio_num_t pin, bool *state)
{
    *state = !*state;
    gpio_set_level(pin, *state ? 0 : 1);
    ESP_LOGI(tag, "Toggled pin %d to %s", pin, *state ? "ON" : "OFF");
}

void set_device(gpio_num_t pin, bool *state, bool on)
{
    *state = on;
    gpio_set_level(pin, *state ? 0 : 1);
    ESP_LOGI(tag, "Set pin %d to %s", pin, *state ? "ON" : "OFF");
}

void update_device_state(device_state_t *state)
{
    memcpy(&device_state, state, sizeof(device_state_t));
    
    gpio_set_level(FAN_PIN, state->fan ? 0 : 1);
    gpio_set_level(HEATER_PIN, state->heater ? 0 : 1);
    gpio_set_level(FEEDER_MOTOR_PIN, state->feeder_motor ? 0 : 1);
    gpio_set_level(CONVEYER_PIN, state->conveyer ? 0 : 1);
}
