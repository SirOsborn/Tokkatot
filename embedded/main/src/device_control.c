#include "device_control.h"
#include "esp_log.h"
#include "esp_timer.h"
#include "string.h"

static const char *tag = "device_control";

typedef struct {
    esp_timer_handle_t timer;
    gpio_num_t pin;
    bool *state_ptr;
} timed_actuator_t;

// Max number of physical actuators (Fan, Heater, Feeder, Conveyer)
#define MAX_TIMED_ACTUATORS 4
static timed_actuator_t timed_actuators[MAX_TIMED_ACTUATORS];
static int actuator_count = 0;

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
    // Cancel any active timer for this pin if it exists
    for (int i = 0; i < actuator_count; i++) {
        if (timed_actuators[i].pin == pin) {
            esp_timer_stop(timed_actuators[i].timer);
            ESP_LOGD(tag, "Cancelled active timer for pin %d due to manual override", pin);
        }
    }

    *state = on;
    gpio_set_level(pin, *state ? 0 : 1);
    ESP_LOGI(tag, "Set pin %d to %s", pin, *state ? "ON" : "OFF");
}

static void actuator_timer_callback(void* arg)
{
    timed_actuator_t* act = (timed_actuator_t*)arg;
    ESP_LOGI(tag, "Timed event: Turning OFF pin %d", act->pin);
    *(act->state_ptr) = false;
    gpio_set_level(act->pin, 1); // Turn OFF (active low)
}

void set_device_timed(gpio_num_t pin, bool *state, bool on, uint32_t duration_sec)
{
    // 1. Initial set
    set_device(pin, state, on);

    // 2. If ON and duration provided, start timer
    if (on && duration_sec > 0) {
        timed_actuator_t* act = NULL;
        
        // Find existing or new slot
        for (int i = 0; i < actuator_count; i++) {
            if (timed_actuators[i].pin == pin) {
                act = &timed_actuators[i];
                break;
            }
        }

        if (act == NULL && actuator_count < MAX_TIMED_ACTUATORS) {
            act = &timed_actuators[actuator_count++];
            act->pin = pin;
            act->state_ptr = state;
            
            const esp_timer_create_args_t timer_args = {
                .callback = &actuator_timer_callback,
                .arg = (void*)act,
                .name = "actuator_timer"
            };
            esp_timer_create(&timer_args, &act->timer);
        }

        if (act) {
            ESP_LOGI(tag, "Scheduling OFF for pin %d in %lu seconds", pin, duration_sec);
            esp_timer_start_once(act->timer, (uint64_t)duration_sec * 1000000ULL);
        }
    }
}

void update_device_state(device_state_t *state)
{
    memcpy(&device_state, state, sizeof(device_state_t));
    
    gpio_set_level(FAN_PIN, state->fan ? 0 : 1);
    gpio_set_level(HEATER_PIN, state->heater ? 0 : 1);
    gpio_set_level(FEEDER_MOTOR_PIN, state->feeder_motor ? 0 : 1);
    gpio_set_level(CONVEYER_PIN, state->conveyer ? 0 : 1);
}
