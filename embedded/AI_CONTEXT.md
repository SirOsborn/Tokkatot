# 🤖 AI Context: ESP32 Embedded Firmware

**Component**: `embedded/` - IoT edge devices for farm equipment control  
**Tech Stack**: ESP32, ESP-IDF (C/C++), FreeRTOS, MQTT, GPIO/PWM  
**Purpose**: Real-time device control (pumps, lights, fans), sensor data collection (DHT22), MQTT communication with local hub  

---

## 📖 Read First

**Before reading this file**, understand the project context:
- **Project overview**: Read [`../AI_INSTRUCTIONS.md`](../AI_INSTRUCTIONS.md) for business model, Farm→Coop→Device hierarchy, why IoT devices matter to farmers
- **Full embedded spec**: See [`../docs/implementation/EMBEDDED.md`](../docs/implementation/EMBEDDED.md) for complete hardware pinout, MQTT topics, OTA updates

**This file contains**: ESP-IDF patterns, FreeRTOS tasks, MQTT communication, DHT22 sensor reading, GPIO/PWM control

---

## 📚 Full Documentation

| Document | Purpose |
|----------|---------|
| [`docs/implementation/EMBEDDED.md`](../docs/implementation/EMBEDDED.md) | Complete embedded specs (pinout, MQTT protocol, OTA updates) |
| [`docs/implementation/API.md`](../docs/implementation/API.md) | Backend device endpoints that sync with ESP32 state |
| [`docs/AUTOMATION_USE_CASES.md`](../docs/AUTOMATION_USE_CASES.md) | Real farmer automation (pulse feeding, climate control) |
| [`AI_INSTRUCTIONS.md`](../AI_INSTRUCTIONS.md) | Why farmers need IoT devices (labor reduction, precision control) |

---

## 📁 Quick File Reference

```
embedded/
├── CMakeLists.txt                     # Build configuration
├── sdkconfig, sdkconfig.old           # ESP32 configuration
├── main/
│   ├── CMakeLists.txt
│   ├── idf_component.yml              # Component dependencies
│   ├── include/
│   │   └── device_config.h            # GPIO pin definitions, device types
│   └── src/
│       ├── main.c                     # Entry point (app_main, FreeRTOS tasks)
│       ├── mqtt_handler.c             # MQTT client, command handling
│       ├── device_control.c           # GPIO relay control, PWM fan/light speed
│       ├── sensor_reader.c            # DHT22 polling (30s interval)
│       └── ota_updater.c              # Over-The-Air firmware updates
├── components/
│   └── dht/                           # DHT22 driver (temperature, humidity)
├── managed_components/
│   └── espressif__servo/              # Servo motor control (optional)
├── build/                             # Build artifacts (NOT COMMITTED - see .gitignore)
└── AI_CONTEXT.md                      # This file
```

---

## 🎯 IoT Device Purpose (Farmer Context)

**Why ESP32 Devices Exist**:
- Farmers cannot manually control equipment 24/7 (sleep, farm labor, family time)
- Precision control needed: Turn on feeder at 6AM for exactly 30 seconds (prevent waste)
- Remote monitoring: Check coop temperature from home (no need to walk 500 meters)
- Emergency response: Automatic fan activation if temperature > 35°C (prevent heat stroke in birds)

**Device Types**:
1. **Water Pumps**: GPIO relay (on/off), scheduled watering times
2. **Lights**: GPIO relay or PWM dimming (day/night simulation for egg production)
3. **Fans**: PWM speed control (0-100% based on temperature)
4. **Feeders**: Stepper motor or conveyor belt (pulse sequences, see [`docs/AUTOMATION_USE_CASES.md`](../docs/AUTOMATION_USE_CASES.md))
5. **Heaters**: GPIO relay with temperature feedback (winter warming)
6. **Sensors**: DHT22 (temp/humidity), ultrasonic (water level), analog (light intensity)

**Communication Flow**:
```
ESP32 Device
    ↓ MQTT (local WiFi)
Raspberry Pi Local Hub (MQTT broker)
    ↓ WebSocket (internet)
Go Middleware (cloud server)
    ↓ HTTP API / WebSocket
Frontend (farmer's phone)
```

**Impact**: 80% labor reduction (manual watering 5x/day → scheduled automation).

---

## 🔌 Hardware Pinout (Standard Tokkatot Device)

**GPIO Mapping** (see `main/include/device_config.h`):

```c
// Relay outputs (devices)
#define RELAY_PIN_1       GPIO_NUM_12  // Water pump, light, etc
#define RELAY_PIN_2       GPIO_NUM_13  // Backup device

// PWM outputs (variable speed)
#define PWM_FAN_PIN       GPIO_NUM_14  // Fan speed (0-100%)
#define PWM_LIGHT_PIN     GPIO_NUM_15  // Light dimming (0-100%)

// Sensor inputs
#define DHT22_PIN         GPIO_NUM_26  // Temperature & humidity
#define BUTTON_PIN        GPIO_NUM_27  // Manual control button
#define ANALOG_SENSOR_PIN GPIO_NUM_25  // Analog ADC (water level, light)

// Status LEDs
#define LED_RED_PIN       GPIO_NUM_2   // Error indicator
#define LED_GREEN_PIN     GPIO_NUM_4   // Online indicator
```

**Wiring Example** (Water Pump):
```
ESP32 GPIO 12 ──── [5V Relay Module IN] ──── [Water Pump 220V]
ESP32 GND     ──── [Relay GND]
ESP32 5V      ──── [Relay VCC]
```

**DHT22 Sensor** (Temperature/Humidity):
```
DHT22 VCC  ──── ESP32 3.3V
DHT22 DATA ──── ESP32 GPIO 26 + 4.7kΩ pull-up resistor to 3.3V
DHT22 GND  ──── ESP32 GND
```

---

## 🛠️ Core ESP-IDF Patterns

### Initialization (`main.c`)

```c
void app_main(void) {
    // 1. Initialize non-volatile storage (NVS - WiFi credentials, device config)
    esp_err_t ret = nvs_flash_init();
    if (ret == ESP_ERR_NVS_NO_FREE_PAGES || ret == ESP_ERR_NVS_NEW_VERSION_FOUND) {
        ESP_ERROR_CHECK(nvs_flash_erase());
        ret = nvs_flash_init();
    }
    ESP_ERROR_CHECK(ret);
    
    // 2. Initialize GPIO pins (relays, LEDs, buttons)
    init_gpio();
    
    // 3. Connect to WiFi (credentials from NVS or provisioning)
    wifi_init_sta();
    
    // 4. Initialize MQTT client (local hub: Raspberry Pi)
    mqtt_app_start();
    
    // 5. Create FreeRTOS tasks (sensor polling, MQTT keepalive)
    xTaskCreate(sensor_polling_task, "sensor", 4096, NULL, 5, NULL);
    xTaskCreate(mqtt_keepalive_task, "mqtt", 4096, NULL, 4, NULL);
    
    // 6. Set status LED (green = online)
    gpio_set_level(LED_GREEN_PIN, 1);
    
    ESP_LOGI(TAG, "Tokkatot device initialized successfully!");
}
```

### GPIO Relay Control (`device_control.c`)

```c
void init_gpio(void) {
    // Configure relay pins as output
    gpio_config_t io_conf = {
        .pin_bit_mask = (1ULL << RELAY_PIN_1) | (1ULL << RELAY_PIN_2),
        .mode = GPIO_MODE_OUTPUT,
        .pull_up_en = GPIO_PULLUP_DISABLE,
        .pull_down_en = GPIO_PULLDOWN_DISABLE,
        .intr_type = GPIO_INTR_DISABLE
    };
    gpio_config(&io_conf);
    
    // Initialize to OFF
    gpio_set_level(RELAY_PIN_1, 0);
    gpio_set_level(RELAY_PIN_2, 0);
}

void turn_relay_on(gpio_num_t pin) {
    gpio_set_level(pin, 1);
    ESP_LOGI(TAG, "Relay ON: GPIO %d", pin);
}

void turn_relay_off(gpio_num_t pin) {
    gpio_set_level(pin, 0);
    ESP_LOGI(TAG, "Relay OFF: GPIO %d", pin);
}
```

### PWM Fan/Light Control (`device_control.c`)

```c
void init_pwm(void) {
    // Configure LEDC timer (for PWM)
    ledc_timer_config_t ledc_timer = {
        .speed_mode = LEDC_LOW_SPEED_MODE,
        .timer_num = LEDC_TIMER_0,
        .duty_resolution = LEDC_TIMER_8_BIT,  // 0-255
        .freq_hz = 5000,                      // 5 kHz
        .clk_cfg = LEDC_AUTO_CLK
    };
    ESP_ERROR_CHECK(ledc_timer_config(&ledc_timer));
    
    // Configure LEDC channel (fan)
    ledc_channel_config_t ledc_channel = {
        .speed_mode = LEDC_LOW_SPEED_MODE,
        .channel = LEDC_CHANNEL_0,
        .timer_sel = LEDC_TIMER_0,
        .gpio_num = PWM_FAN_PIN,
        .duty = 0,               // Start at 0% (off)
        .hpoint = 0
    };
    ESP_ERROR_CHECK(ledc_channel_config(&ledc_channel));
}

void set_fan_speed(uint8_t percent) {
    // percent: 0-100
    uint32_t duty = (percent * 255) / 100;
    ESP_ERROR_CHECK(ledc_set_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_0, duty));
    ESP_ERROR_CHECK(ledc_update_duty(LEDC_LOW_SPEED_MODE, LEDC_CHANNEL_0));
    ESP_LOGI(TAG, "Fan speed: %d%%", percent);
}
```

### DHT22 Sensor Reading (`sensor_reader.c`)

```c
void sensor_polling_task(void *pvParameters) {
    while (1) {
        float temperature, humidity;
        
        // Read DHT22 sensor
        esp_err_t ret = dht_read_float_data(DHT_TYPE_DHT22, DHT22_PIN, &humidity, &temperature);
        
        if (ret == ESP_OK) {
            ESP_LOGI(TAG, "Temperature: %.1f°C, Humidity: %.1f%%", temperature, humidity);
            
            // Publish to MQTT
            char payload[128];
            snprintf(payload, sizeof(payload), 
                     "{\"temperature\": %.1f, \"humidity\": %.1f}", 
                     temperature, humidity);
            
            char topic[128];
            snprintf(topic, sizeof(topic), 
                     "farm/%s/devices/%s/sensor", 
                     FARM_ID, DEVICE_ID);
            
            mqtt_publish(topic, payload);
        } else {
            ESP_LOGE(TAG, "Failed to read DHT22 sensor");
        }
        
        // Wait 30 seconds before next reading
        vTaskDelay(30000 / portTICK_PERIOD_MS);
    }
}
```

---

## 📡 MQTT Communication

**MQTT Broker**: Raspberry Pi local hub (192.168.1.100:1883)

### MQTT Topics (Subscribe - Receive Commands)

**Device Control**:
```
Topic: farm/{farm_id}/devices/{device_id}/command
Payload: {
  "command": "on",           // "on", "off", "pwm"
  "pwm_value": 75,           // 0-100% (for PWM devices)
  "duration_seconds": 3600   // Auto-turn-off after 1 hour (optional)
}

Example: farm/farm123/devices/pump001/command
         {"command": "on", "duration_seconds": 300}
```

**Configuration Updates**:
```
Topic: farm/{farm_id}/devices/{device_id}/config
Payload: {
  "relay_pin": 12,
  "device_type": "pump",
  "sensor_interval_seconds": 30
}
```

**OTA Firmware Updates**:
```
Topic: farm/{farm_id}/devices/{device_id}/ota
Payload: {
  "firmware_url": "https://server.com/firmware_v1.2.bin",
  "version": "1.2.0"
}
```

### MQTT Topics (Publish - Send Status)

**Device Status** (Heartbeat - Every 30 seconds):
```
Topic: farm/{farm_id}/devices/{device_id}/status
Payload: {
  "is_online": true,
  "command_state": "on",
  "uptime_seconds": 12345,
  "rssi": -65,              // WiFi signal strength
  "free_heap": 100000       // Free memory (bytes)
}
```

**Sensor Data** (Every 30 seconds):
```
Topic: farm/{farm_id}/devices/{device_id}/sensor
Payload: {
  "type": "temperature_humidity",
  "temperature": 28.5,
  "humidity": 65.2,
  "timestamp": "2025-02-01T12:00:00Z"
}
```

**Error Reporting**:
```
Topic: farm/{farm_id}/devices/{device_id}/error
Payload: {
  "error": "DHT22 sensor read failed",
  "severity": "warning",
  "timestamp": "2025-02-01T12:00:00Z"
}
```

### MQTT Handler (`mqtt_handler.c`)

```c
static void mqtt_event_handler(void *handler_args, esp_event_base_t base, 
                                int32_t event_id, void *event_data) {
    esp_mqtt_event_handle_t event = event_data;
    
    switch ((esp_mqtt_event_id_t)event_id) {
        case MQTT_EVENT_CONNECTED:
            ESP_LOGI(TAG, "MQTT connected!");
            
            // Subscribe to command topic
            char topic[128];
            snprintf(topic, sizeof(topic), "farm/%s/devices/%s/command", FARM_ID, DEVICE_ID);
            esp_mqtt_client_subscribe(event->client, topic, 1);
            
            gpio_set_level(LED_GREEN_PIN, 1);  // Green LED = online
            publish_status();
            break;
            
        case MQTT_EVENT_DATA:
            ESP_LOGI(TAG, "MQTT message received: %.*s", event->data_len, event->data);
            handle_mqtt_command(event->data, event->data_len);
            break;
            
        case MQTT_EVENT_DISCONNECTED:
            ESP_LOGW(TAG, "MQTT disconnected!");
            gpio_set_level(LED_RED_PIN, 1);   // Red LED = error
            break;
    }
}

void handle_mqtt_command(const char *payload, int len) {
    // Parse JSON payload
    cJSON *root = cJSON_ParseWithLength(payload, len);
    if (!root) {
        ESP_LOGE(TAG, "Failed to parse JSON");
        return;
    }
    
    const char *command = cJSON_GetStringValue(cJSON_GetObjectItem(root, "command"));
    
    if (strcmp(command, "on") == 0) {
        turn_relay_on(RELAY_PIN_1);
    } else if (strcmp(command, "off") == 0) {
        turn_relay_off(RELAY_PIN_1);
    } else if (strcmp(command, "pwm") == 0) {
        int pwm_value = cJSON_GetNumberValue(cJSON_GetObjectItem(root, "pwm_value"));
        set_fan_speed(pwm_value);
    }
    
    cJSON_Delete(root);
    
    // Publish status update
    publish_status();
}
```

---

## 🔄 Multi-Step Automation (Action Sequences)

**Context**: Farmers need pulse sequences for feeders (ON 30s, pause 10s, repeat) to prevent chicken congestion.

**Implementation Strategy**:
1. **Backend creates schedule** with `action_sequence` JSON (see [`docs/AUTOMATION_USE_CASES.md`](../docs/AUTOMATION_USE_CASES.md))
2. **Go middleware sends MQTT command** with full sequence
3. **ESP32 executes sequence** locally (no internet dependency)

**MQTT Command Example** (Pulse Feeding):
```json
Topic: farm/farm123/devices/feeder001/sequence
Payload: {
  "action_sequence": [
    {"action": "ON", "duration": 30},
    {"action": "OFF", "duration": 10},
    {"action": "ON", "duration": 30},
    {"action": "OFF", "duration": 10},
    {"action": "OFF", "duration": 0}
  ]
}
```

**ESP32 Execution** (`device_control.c`):
```c
void execute_action_sequence(cJSON *sequence_array) {
    int count = cJSON_GetArraySize(sequence_array);
    
    for (int i = 0; i < count; i++) {
        cJSON *step = cJSON_GetArrayItem(sequence_array, i);
        const char *action = cJSON_GetStringValue(cJSON_GetObjectItem(step, "action"));
        int duration = cJSON_GetNumberValue(cJSON_GetObjectItem(step, "duration"));
        
        ESP_LOGI(TAG, "Step %d: %s for %d seconds", i+1, action, duration);
        
        if (strcmp(action, "ON") == 0) {
            turn_relay_on(RELAY_PIN_1);
        } else if (strcmp(action, "OFF") == 0) {
            turn_relay_off(RELAY_PIN_1);
        }
        
        // Wait for duration (blocking)
        if (duration > 0) {
            vTaskDelay((duration * 1000) / portTICK_PERIOD_MS);
        } else {
            break;  // duration=0 means "until next schedule"
        }
    }
    
    ESP_LOGI(TAG, "Action sequence complete");
}
```

**Benefit**: Feeder pulse prevents crowding at feed bowl, reduces waste by 30% (see farmer testimonials in [`docs/AUTOMATION_USE_CASES.md`](../docs/AUTOMATION_USE_CASES.md)).

---

## 🔒 Security Best Practices

- ✅ **WiFi Credentials**: Store in NVS (encrypted), never hardcode in firmware
- ✅ **MQTT Auth**: Username/password authentication to broker (not public)
- ✅ **Device Identification**: MAC address or TLS certificate for device identity
- ✅ **OTA Signing**: Verify firmware signature before flashing (prevent malicious updates)
- ✅ **Input Validation**: Validate MQTT JSON commands (prevent buffer overflows)
- ✅ **Watchdog Timer**: Auto-reboot if firmware hangs (prevent permanent failures)

---

## 🆘 Common Issues & Solutions

| Issue | Solution |
|-------|----------|
| **WiFi won't connect** | Check SSID/password in NVS, verify signal strength, reboot device |
| **MQTT disconnects** | Verify Raspberry Pi broker running, check firewall (port 1883), reduce publish frequency |
| **DHT22 sensor fails** | Add 4.7kΩ pull-up resistor, verify GPIO pin, slow reading interval (not < 2s) |
| **Relay stuck/won't respond** | Check relay module power (3.3V or 5V), test GPIO with simple blink, verify wiring |
| **OTA update fails** | Check firmware binary signature, ensure flash storage available, verify internet access |
| **Device reboots randomly** | Feed watchdog timer, reduce memory usage, check power supply stability |

---

## 🧪 Development Tasks

### Build & Flash Firmware

```bash
cd embedded

# Configure project (WiFi, MQTT broker IP, device ID)
idf.py menuconfig
# Component config → Tokkatot Device Config → Set WiFi SSID, password, MQTT broker IP

# Build firmware
idf.py build

# Flash to ESP32
idf.py -p COM3 flash  # Windows (check Device Manager for port)
idf.py -p /dev/ttyUSB0 flash  # Linux

# Monitor serial output
idf.py -p COM3 monitor
# Press Ctrl+] to exit monitor
```

### Test GPIO Relay (Simple Blink)

```c
// Quick test: Blink relay to verify wiring
void test_relay_blink(void) {
    while (1) {
        gpio_set_level(RELAY_PIN_1, 1);
        vTaskDelay(1000 / portTICK_PERIOD_MS);  // 1 second ON
        gpio_set_level(RELAY_PIN_1, 0);
        vTaskDelay(1000 / portTICK_PERIOD_MS);  // 1 second OFF
    }
}
```

### Test MQTT Publish (Local)

```bash
# On Raspberry Pi (MQTT broker):
mosquitto_sub -h localhost -t "farm/+/devices/+/sensor" -v

# Expected output when ESP32 publishes:
# farm/farm123/devices/pump001/sensor {"temperature": 28.5, "humidity": 65.2}
```

---

## 📘 Documentation Map

**AI Context Files** (component-specific guides):
- **This file**: [`embedded/AI_CONTEXT.md`](./AI_CONTEXT.md) - ESP-IDF patterns, MQTT, GPIO/PWM
- [`middleware/AI_CONTEXT.md`](../middleware/AI_CONTEXT.md) - Go API, how it sends MQTT commands
- [`frontend/AI_CONTEXT.md`](../frontend/AI_CONTEXT.md) - Vue.js device control UI
- [`ai-service/AI_CONTEXT.md`](../ai-service/AI_CONTEXT.md) - PyTorch disease detection (no direct connection to ESP32)
- [`docs/AI_CONTEXT.md`](../docs/AI_CONTEXT.md) - Documentation maintenance guide

**Master Guide**: [`AI_INSTRUCTIONS.md`](../AI_INSTRUCTIONS.md) - Read first for project overview

---

**Happy coding! 🚀 Test on real hardware, not just in IDE!**
