# 🤖 AI Context: ESP32 Embedded Firmware

**Directory**: `embedded/`  
**Your Role**: Microcontroller firmware, hardware control, sensor reading, MQTT communication  
**Tech Stack**: ESP32, ESP-IDF (C/C++), MQTT, GPIO, PWM  

---

## 🎯 What You're Building

**IoT Edge Devices** (ESP32 Microcontrollers)

Each ESP32 controls one or more farm devices:
- **Water Pumps**: GPIO relay control (on/off)
- **Lights**: GPIO on/off or PWM dimming
- **Fans**: GPIO on/off or PWM speed control
- **Heaters**: GPIO on/off with temperature feedback
- **Feeders}: PWM stepper motor control
- **Conveyors**: GPIO on/off
- **Sensors**: DHT22 (temperature, humidity), analog sensors
- **Emergency**: Button inputs, buzzer outputs

**Key Capabilities**:
- MQTT communication with Raspberry Pi local hub
- Sensor data reading every 30 seconds
- Command execution with status feedback
- Over-the-Air (OTA) firmware updates
- Offline operation (queue commands locally)
- Status LED indicators
- Energy efficiency (sleep modes)

---

## 📁 File Structure

```
embedded/
├── CMakeLists.txt              # Build configuration
├── sdkconfig                   # ESP32 configuration
├── sdkconfig.old               # Previous config (backup)
├── dependencies.lock           # Dependency versions
├── main/
│   ├── CMakeLists.txt
│   ├── idf_component.yml
│   ├── include/                # Header files
│   │   └── device_config.h
│   └── src/
│       ├── main.c              # Entry point
│       ├── mqtt_handler.c      # MQTT communication
│       ├── device_control.c    # GPIO, PWM control
│       ├── sensor_reader.c     # Sensor polling
│       └── ota_updater.c       # Firmware updates
├── components/
│   ├── dht/                    # DHT22 driver
│   │   ├── CMakeLists.txt
│   │   └── dht.c
│   └── (other custom components)
├── managed_components/         # IDF component registry
├── build/                      # Build artifacts (DO NOT COMMIT)
└── AI_CONTEXT.md              # This file
```

---

## 🚀 Getting Started

### Local Setup

```bash
cd embedded

# Install ESP-IDF (if not already done)
# https://docs.espressif.com/projects/esp-idf/en/latest/esp32/get-started/

# Configure project
idf.py menuconfig
# Navigate to: Component config → Configure WiFi, MQTT settings

# Build
idf.py build

# Flash to device
idf.py -p /dev/ttyUSB0 flash  # Linux
idf.py -p COM3 flash          # Windows
idf.py -p /dev/ttyUSB0 monitor # Watch output
```

### Building for Different Devices

```bash
# Set target
idf.py set-target esp32

# Build just bootloader
idf.py build bootloader

# Build just partition table
idf.py build partition-table
```

---

## 🔌 Hardware Pinout

### Standard Tokkatot Device Setup

```
ESP32 GPIO Mapping:
- GPIO 12: Relay 1 (pump, light, etc)
- GPIO 13: Relay 2 (backup device)
- GPIO 14: PWM output (fan speed, dimming)
- GPIO 15: PWM output 2
- GPIO 25: ADC input (sensor analog)
- GPIO 26: DHT22 data (temperature/humidity)
- GPIO 27: Button input (manual control)
- GPIO 2:  Status LED (red = error)
- GPIO 4:  Status LED (green = online)

UART:
- GPIO 1: TX (debug)
- GPIO 3: RX (debug)
```

### Wiring Example (Water Pump with Relay)

```
ESP32 GPIO 12 ──── [Relay Module] ──── [Water Pump]
ESP32 GND ──────── [Relay GND]
ESP32 3.3V ──────── [Relay Power]

Status feedback (optional):
Pump status pin ──── GPIO 34 (input ADC)
```

---

## 💾 Software Architecture

### Initialization Flow

```c
void setup() {
  1. Initialize WiFi (from NVS or provisioning)
  2. Connect to local MQTT broker (Raspberry Pi)
  3. Initialize GPIO pins
  4. Initialize sensor readers
  5. Subscribe to MQTT topics: farm/{farm_id}/devices/{device_id}/+
  6. Send heartbeat: farm/{farm_id}/devices/{device_id}/status
  7. Start sensor polling task (every 30 seconds)
}
```

### MQTT Topics

**Subscribe** (Listen for commands):
```
farm/{farm_id}/devices/{device_id}/command
  Payload: {"command": "on|off|pwm", "pwm_value": 255, "duration_seconds": 3600}
  
farm/{farm_id}/devices/{device_id}/config
  Payload: {"relay_pin": 12, "type": "pump"}
```

**Publish** (Report status):
```
farm/{farm_id}/devices/{device_id}/status
  Payload: {"is_online": true, "command_state": "on", "uptime_seconds": 12345}
  Published: Every 30 seconds (heartbeat)
  
farm/{farm_id}/devices/{device_id}/sensor
  Payload: {"type": "temperature", "value": 28.5, "humidity": 65.2}
  Published: Every 30 seconds (from DHT22)
  
farm/{farm_id}/devices/{device_id}/error
  Payload: {"error": "MQTT disconnected", "timestamp": "2026-02-19T10:30:00Z"}
  Published: On errors (sensor failures, etc)
```

---

## 🔧 Key Functions

### `main.c`

```c
void app_main(void) {
  // 1. Initialize NVS (non-volatile storage)
  nvs_flash_init();
  
  // 2. Initialize board (GPIO, UART)
  init_gpio();
  init_uart();
  
  // 3. Connect WiFi
  enum wifi_error_t err = connect_wifi();
  if (err != WIFI_OK) handle_error("WiFi failed");
  
  // 4. Initialize MQTT
  mqtt_client_config_t mqtt_cfg = {
    .broker_address = "192.168.1.100",  // Local hub IP
    .port = 1883,
    .username = "esp32_device",
    .password = "mqtt_password"
  };
  mqtt_app_start(&mqtt_cfg);
  
  // 5. Create tasks
  xTaskCreate(sensor_polling_task, "sensor", 4096, NULL, 5, NULL);
  xTaskCreate(mqtt_keepalive_task, "mqtt", 4096, NULL, 4, NULL);
  
  // 6. Ready
  set_led(GREEN, true);
  log_info("Device initialized successfully!");
}
```

### `mqtt_handler.c`

```c
void mqtt_event_handler(esp_mqtt_event_handle_t event) {
  switch(event->event_id) {
    case MQTT_EVENT_CONNECTED:
      log_info("MQTT connected!");
      subscribe_to_commands();
      publish_status();
      break;
      
    case MQTT_EVENT_DATA:
      handle_mqtt_command(event->topic, event->data);
      break;
      
    case MQTT_EVENT_DISCONNECTED:
      log_error("MQTT disconnected!");
      set_led(RED, true);
      break;
  }
}

void handle_mqtt_command(char *topic, char *payload) {
  // Parse JSON: {"command": "on|off", "pwm_value": 255}
  cJSON *root = cJSON_Parse(payload);
  char *command = cJSON_GetStringValue(root, "command");
  
  if (strcmp(command, "on") == 0) {
    gpio_set_level(RELAY_PIN, 1);
  } else if (strcmp(command, "off") == 0) {
    gpio_set_level(RELAY_PIN, 0);
  } else if (strcmp(command, "pwm") == 0) {
    int value = cJSON_GetNumberValue(root, "pwm_value");
    ledc_set_duty(LEDC_MODE, PWM_CHANNEL, value);
  }
  
  // Publish status back
  publish_status();
}
```

### `device_control.c`

```c
void init_gpio(void) {
  // Configure GPIO
  gpio_config_t io_conf = {
    .pin_bit_mask = (1ULL << RELAY_PIN) | (1ULL << BUTTON_PIN),
    .mode = GPIO_MODE_OUTPUT,  // or GPIO_MODE_INPUT for buttons
    .pull_up_en = 0,
    .intr_type = GPIO_INTR_DISABLE
  };
  gpio_config(&io_conf);
  
  // Configure PWM
  ledc_timer_config_t timer_conf = {
    .speed_mode = LEDC_MODE,
    .timer_num = LEDC_TIMER,
    .freq_hz = 5000,
    .duty_resolution = LEDC_TIMER_8_BIT
  };
  ledc_timer_config(&timer_conf);
  
  ledc_channel_config_t channel_conf = {
    .speed_mode = LEDC_MODE,
    .channel = PWM_CHANNEL,
    .gpio_num = PWM_PIN,
    .timer_sel = LEDC_TIMER,
    .duty = 0
  };
  ledc_channel_config(&channel_conf);
}

void turn_relay_on(int pin) {
  gpio_set_level(pin, 1);
  log_info("Relay ON: GPIO %d", pin);
}

void turn_relay_off(int pin) {
  gpio_set_level(pin, 0);
  log_info("Relay OFF: GPIO %d", pin);
}

void set_pwm_duty(int value) {
  // value: 0-255
  ledc_set_duty(LEDC_MODE, PWM_CHANNEL, (value * 255) / 100);
  ledc_update_duty(LEDC_MODE, PWM_CHANNEL);
}
```

### `sensor_reader.c`

```c
void sensor_polling_task(void *pvParameters) {
  while (1) {
    // Read DHT22 temperature/humidity
    float temp, humidity;
    int err = dht_read_float_data(DHT_TYPE_DHT22, DHT_PIN, &humidity, &temp);
    
    if (err == DHT_OK) {
      // Create JSON payload
      char payload[256];
      sprintf(payload, "{\"temperature\": %.1f, \"humidity\": %.1f}", temp, humidity);
      
      // Publish to MQTT
      esp_mqtt_client_publish(
        mqtt_client,
        "farm/farm123/devices/device456/sensor",
        payload,
        strlen(payload),
        1,  // qos
        0   // retain
      );
      
      log_info("Sensor data sent: T=%.1f°C, H=%.1f%%", temp, humidity);
    } else {
      log_error("DHT read failed: %d", err);
    }
    
    // Wait 30 seconds before next reading
    vTaskDelay(30000 / portTICK_PERIOD_MS);
  }
}
```

---

## 📝 Code Guidelines

### ✅ DO:
- Use FreeRTOS tasks for concurrent operations
- Implement error handling with logging
- Set status LEDs for visual feedback
- Publish heartbeat every 30 seconds
- Use JSON for MQTT payloads
- Handle WiFi disconnects gracefully
- Store device config in NVS (non-volatile storage)
- Test hardware connections before firmware

### ❌ DON'T:
- Block main app_main() function (use tasks)
- Hardcode GPIO pins (use macros from header)
- Ignore error codes from API calls
- Allocate large buffers on stack (use heap)
- Use blocking MQTT calls
- Assume WiFi is always connected
- Publish every sensor reading (buffer & aggregate)
- Ignore watchdog timer (implement kicking)

---

## 🔒 Security Checklist

- ✅ WiFi credentials in NVS (encrypted storage), not hardcode
- ✅ MQTT password protected (username/password auth)
- ✅ Device identification via MAC address or certificate
- ✅ OTA updates signed/verified before flashing
- ✅ Input validation on MQTT commands (JSON parsing safe)
- ✅ Idle timeout (disconnect if no heartbeat)
- ✅ Button input debouncing (prevent accidental triggers)

---

## 🛠️ Firmware Update (OTA)

### Process

```c
void ota_updater_task(void *pvParameters) {
  // 1. Check for firmware update from server
  esp_https_ota_config_t config = {
    .http_config = &http_config,
    .ota_buf_size = 4096
  };
  
  // 2. Perform OTA
  esp_err_t ret = esp_https_ota(&config);
  
  if (ret == ESP_OK) {
    log_info("OTA update successful!");
    esp_restart();  // Reboot
  } else {
    log_error("OTA failed: %s", esp_err_to_name(ret));
  }
}
```

---

## 🆘 Common Issues & Solutions

### Issue: Device won't connect to WiFi
```
Error: WiFi connection timeout
```
**Fix**: 
- Check SSID/password correct in NVS
- Verify WiFi signal strength
- Reboot device

### Issue: MQTT disconnects frequently
```
Error: MQTT disconnected
```
**Fix**: 
- Check MQTT broker connection (Raspberry Pi running?)
- Verify network connectivity (ping broker IP)
- Check firewall allows port 1883

### Issue: DHT22 sensor fails
```
Error: DHT read failed: -1
```
**Fix**: 
- Check 4.7kΩ pull-up resistor on data line
- Verify GPIO pin assignment
- Try slower reading interval (not too frequent)

### Issue: Relay stuck/won't respond
```
Relay pin not changing state
```
**Fix**: 
- Verify relay module power (3.3V)
- Check GPIO pin wired correctly
- Verify relay module responds to logic levels
- Test with simple GPIO blink test

### Issue: OTA update fails
```
Error: OTA checksum verification failed
```
**Fix**: 
- Verify firmware binary is signed correctly
- Check internet connection during OTA
- Ensure device has enough flash storage

---

## 📊 Performance & Resources

### Memory Usage
- Typical firmware binary: 400-500KB
- Runtime memory: 100-200KB heap for buffers
- MQTT publish buffer: 4KB
- Sensor buffer: 1KB

### Power Consumption
- Active (WiFi + MQTT): ~80-100mA
- WiFi idle: ~10-20mA
- Deep sleep: ~10µA

---

## 📚 Key Documents

- `IG_SPECIFICATIONS_EMBEDDED.md` - Full embedded specification
- ESP-IDF Documentation: https://docs.espressif.com/projects/esp-idf/
- MQTT Protocol: https://mqtt.org/

---

## 🧪 Testing Checklist

- ✅ Hardware connections tested (relay, sensor, button)
- ✅ GPIO pins output correct logic levels
- ✅ WiFi connects to local hub
- ✅ MQTT connects and publishes
- ✅ Commands received via MQTT work correctly
- ✅ Sensor readings accurate
- ✅ Status LED indicators work
- ✅ Device reboots gracefully (watchdog)
- ✅ OTA firmware updates work

---

## 🎯 Your Next Tasks

1. **Setup hardware** - Wire GPIO, relays, sensors correctly
2. **Configure WiFi** - Store credentials in NVS
3. **Implement MQTT** - Connect to local hub
4. **Add device control** - GPIO on/off, PWM
5. **Add sensors** - DHT22 temperature/humidity
6. **Test thoroughly** - Each component separately

---

**Happy coding! 🚀 Remember: Test on real hardware, not just in IDE!**
