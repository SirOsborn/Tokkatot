# Tokkatot 2.0: Embedded Systems Specification

**Document Version**: 2.0  
**Last Updated**: February 2026  
**Status**: Final Specification

**Primary Microcontroller**: ESP32  
**Development Framework**: ESP-IDF 4.4+  
**Language**: C (with FreeRTOS)

---

## Overview

Tokkatot 2.0 firmware runs on ESP32 microcontrollers deployed on-farm. Each device controls hardware (relays, sensors, PWM outputs) and communicates with the cloud backend and local Raspberry Pi agent.

---

## Firmware Architecture

### High-Level Architecture

```
┌──────────────────────────────────────────┐
│          Application Layer               │
│  ┌──────────┐ ┌──────────┐ ┌─────────┐ │
│  │ MQTT Task│ │GPIO Tasks│ │Sensor   │ │
│  │          │ │          │ │ Tasks   │ │
│  └──────────┘ └──────────┘ └─────────┘ │
├──────────────────────────────────────────┤
│          FreeRTOS Scheduler               │
│  (Handles task switching, 100ms tick)    │
├──────────────────────────────────────────┤
│          Hardware Drivers (ESP-IDF)       │
│  GPIO, UART, SPI, I2C, ADC, PWM, FLASH  │
├──────────────────────────────────────────┤
│          Bootloader & OTA Support         │
│  (Handles firmware updates)               │
└──────────────────────────────────────────┘
```

### Firmware Modules

**1. WiFi Manager**
- Connects to farm WiFi
- Automatic reconnection with exponential backoff
- Stores credentials in secure EEPROM
- Signal strength monitoring
- Fallback to offline mode if connection lost

**2. MQTT Client**
- Connects to local RPi agent or cloud broker
- Publishes device status, sensor readings
- Subscribes to command topic
- TLS support for secure communication
- QoS levels: 1 (at-least-once for critical)
- Connection keep-alive (30 second heartbeat)

**3. Device Control (GPIO, Relays, PWM)**
- GPIO pins for relay control (on/off switching)
- PWM for dimmable devices (fans, lights at variable speeds)
- Status feedback (actual state of device)
- Safety limits (e.g., max ON duration)
- Error detection (short circuit, overload)

**4. Sensor Manager**
- DHT22 temperature/humidity sensor
- ADC for analog sensors (soil moisture, water level)
- Sensor reading with error checking
- Outlier detection and filtering
- Configurable sampling interval (30 sec default)
- Circular buffer for averaging readings

**5. State Machine**
- Offline mode: execute schedules locally
- Online mode: take commands from cloud
- Error mode: safe shutdown if critical error
- Update mode: standby during OTA update

**6. Non-Volatile Storage (NVS)**
- Device configuration (MQTT broker, WiFi credentials)
- Schedule data (for offline execution)
- Firmware version and update status
- Device ID and hardware MAC address
- Uptime counter, error logs

**7. OTA (Over-The-Air) Update Handler**
- Receive firmware binary via MQTT
- Verify signature before installation
- Automatic rollback if boot fails
- Staged update support
- Update notification to cloud

---

## Hardware Interface Specification

### GPIO Pins (ESP32-WROOM-32)

**Reserved Pins**:
```
GPIO 1   → UART TX (debugging, do not use)
GPIO 3   → UART RX (debugging, do not use)
GPIO 6-11 → Flash SPI (do not use)
GPIO 17  → PSRAM SPI (do not use)
GPIO 36 → ADC0 (battery monitoring, optional)
```

**Available Pins for Devices** (suggested allocation):
```
GPIO 12 → Relay 1 (Water Pump)
GPIO 13 → Relay 2 (Feeder Motor)
GPIO 14 → Relay 3 (Light Control)
GPIO 15 → Relay 4 (Fan Control)
GPIO 16 → Relay 5 (Heater Control)
GPIO 17 → Relay 6 (Conveyor/Spare)

GPIO 32 → ADC for analog sensor 1
GPIO 33 → ADC for analog sensor 2
GPIO 34 → ADC for emergency stop / safety

GPIO 2  → Status LED (red, errors)
GPIO 4  → Status LED (green, online)
GPIO 5  → Status LED (blue, activity)

GPIO 19 → I2C SDA (sensors using I2C)
GPIO 23 → I2C SCL (sensors using I2C)

GPIO 21 → OneWire DS18B20 (temperature, optional)
GPIO 22 → DHT22 data line (humidity + temp)

GPIO 25 → PWM out Signal (for servo/dimmer, optional)
GPIO 26 → PWM out 2 (additional pwm)
```

### Relay Module Specifications

**Types Supported**:
-  5V relay module (common, cheap)
- 12V relay module (industrial grade)
- Solid-state relay (SSR, for inductive loads)
- Optocoupler isolation recommended

**Electrical Safety**:
- Free-wheeling diode on relay coil (back-EMF protection)
- Current limiting resistor on GPIO pin
- Max GPIO current: 12mA per pin
- Use MOSFET or NPN transistor for high-current relays

### Sensor Specifications

**DHT22** (Temperature/Humidity):
- Operating range: -40°C to +80°C (±2°C accuracy)
- Humidity: 0-100% RH (±2% accuracy)
- Sampling: Every 30 seconds (2-second timeout)
- Protocol: Single-wire digital

**ADC for Analog Sensors**:
- 12-bit resolution (0-4095 counts)
- Voltage range: 0-3.3V
- Example: Water level float switch
- Example: Soil moisture sensor
- Calibrated with min/max reference values

**Optional Sensors**:
- Weight scale via HX711 (load cell for pellet counter)
- Water flow meter (pulse counting)
- Motion sensor (PIR for intruder detection)
- pH probe (water quality, via ADC or external chip)

---

## Firmware Configuration

### Config Structure (NVS Storage)

```c
typedef struct {
  char device_id[32];           // Unique device ID
  char wifi_ssid[32];           // WiFi network name
  char wifi_password[64];       // WiFi password
  char mqtt_broker_ip[32];      // RPi or cloud IP
  uint16_t mqtt_port;           // 1883 or 8883 for TLS
  uint16_t polling_interval;    // Sensor read interval (sec)
  uint16_t heartbeat_interval;  // MQTT keep-alive (sec)
  uint32_t uptime_seconds;      // Total uptime
  uint8_t online_mode;          // 1=online, 0=offline
  uint8_t fw_version_major;     // Major version
  uint8_t fw_version_minor;     // Minor version
  uint8_t fw_version_patch;     // Patch version
} device_config_t;
```

### MQTT Topics and Payloads

**Device Publishing** (Device → Cloud/RPi):

```
Topic: farm/{farmId}/devices/{deviceId}/heartbeat
Payload: {
  "device_id": "esp32-001",
  "timestamp": 1677000000,
  "online": true,
  "firmware_version": "2.0.0",
  "uptime_seconds": 86400,
  "signal_strength": -65,
  "heap_free": 102400,
  "last_update": "02/18/2026"
}

Topic: farm/{farmId}/devices/{deviceId}/status
Payload: {
  "device_id": "esp32-001",
  "timestamp": 1677000000,
  "state": "on",
  "previous_state": "off",
  "triggered_by": "schedule|user_command|automation",
  "duration_ms": 5000
}

Topic: farm/{farmId}/sensors/{deviceId}/data
Payload: {
  "device_id": "temp-sensor-01",
  "timestamp": 1677000000,
  "readings": {
    "temperature": 28.5,
    "humidity": 65.2,
    "unit": "celsius|percent"
  }
}
```

**Cloud/RPi Publishing** (Cloud → Device):

```
Topic: farm/{farmId}/devices/{deviceId}/command
Payload: {
  "command_id": "cmd-uuid-123",
  "device_id": "esp32-001",
  "action": "on|off|toggle|set_value",
  "duration_ms": 300000,
  "timestamp": 1677000000
}

Topic: farm/{farmId}/devices/{deviceId}/fw_update
Payload: {
  "fw_url": "https://cdn.example.com/firmware/v2.0.0.bin",
  "fw_size": 512000,
  "fw_hash": "sha256:abc123...",
  "fw_version": "2.0.0",
  "timestamp": 1677000000
}
```

**Device Response** (Device → Cloud/RPi):

```
Topic: farm/{farmId}/devices/{deviceId}/status
Payload (after command):
{
  "command_id": "cmd-uuid-123",
  "status": "success|failed|timeout",
  "new_state": "on|off",
  "timestamp": 1677000000,
  "error": "none|timeout|invalid_command|hardware_error"
}
```

---

## Task Scheduling (FreeRTOS)

### Tasks and Priorities

```c
Task Priorities (0=lowest, 24=highest):

Priority 24: WiFi/MQTT reconnection (critical)
Priority 20: OTA firmware update handler
Priority 15: Real-time relay control (user commands)
Priority 14: Sensor reading (DHT22, ADC)
Priority 12: MQTT publish (non-blocking)
Priority 10: Schedule execution (local automation)
Priority 8:  Status LED updates
Priority 5:  Health monitoring
Priority 2:  Debug logging
Priority 1:  Idle task (FreeRTOS default)
```

### Task Stack Sizes

```c
#define WIFI_TASK_STACK_SIZE    (4096)
#define MQTT_TASK_STACK_SIZE    (4096)
#define SENSOR_TASK_STACK_SIZE  (2048)
#define LED_TASK_STACK_SIZE     (1024)
#define COMMAND_TASK_STACK_SIZE (2048)
#define OTA_TASK_STACK_SIZE     (8192)
```

---

## State Diagram

```
         ┌─────────────┐
         │   Booting   │
         └──────┬──────┘
                ↓
    ┌───────────────────────┐
    │ Load config from NVS  │
    └───────────┬───────────┘
                ↓
    ┌───────────────────────┐
    │ Initialize hardware   │
    │ (GPIO, UART, I2C)     │
    └───────────┬───────────┘
                ↓
    ┌─────────────────────────────────┐
    │  Attempt WiFi connection        │◄─────────┐
    │  (with exponential backoff)     │          │
    └───┬──────────────┬──────────────┘          │
        │              │                         │
   Success         Failure                      │
   (3+ hours)      (timeout)                    │
        │              │                         │
        ↓              ↓                         │
    ┌──────────┐    ┌────────────────────────┐ │
    │ Online   │    │ Offline Mode           │ │
    │ Mode     │    │ (execute local         │ │
    │          │    │  schedules)            │ │
    └────┬─────┘    └───────────┬────────────┘ │
         │                      │              │
         ├──→ ┌────────────────┐◄──────────────┤
         │    │ Sensor Loop    │              │
         │    │ (30s interval) │              │
         │    └────────────────┘              │
         │                                    │
         ├──→ ┌────────────────────┐          │
         │    │ MQTT Publish       │          │
         │    │ Status & sensors   │          │
         │    └────────────────────┘          │
         │                                    │
         ├──→ ┌────────────────────┐          │
         │    │ Check for commands │          │
         │    │ (subscribe topic)  │          │
         │    └────────────────────┘          │
         │                                    │
         ├──→ ┌────────────────────┐          │
         │    │ Execute command    │          │
         │    │ (relay control)    │          │
         │    └────────────────────┘          │
         │                                    │
         └────→ Retry WiFi ───────────────────┘
              (on connection loss)
```

---

## Power Management

### Sleep Modes

**Active Mode** (Normal operation):
- WiFi on, MQTT connected
- Current: ~80mA @ 3.3V
- Used for: Always-on operation

**Light Sleep** (WiFi off, CPU active):
- Current: ~15mA @ 3.3V
- Used for: Battery operation (future)
- Wake on: GPIO interrupt, timer

**Deep Sleep** (CPU off, RTC memory on):
- Current: ~10µA @ 3.3V
- Used for: Extended battery mode
- Wake on: RTC timer or external GPIO

**Current Implementation**: Active mode always (assuming constant power)

---

## Watchdog & Recovery

### Watchdog Timer

```c
esp_task_wdt_init(10, true);  // 10-second timeout
esp_task_wdt_add(NULL);        // Watch current task

// Within each task:
esp_task_wdt_reset();          // Reset timer every loop
```

**On Watchdog Timeout**:
1. Automatic restart (chip resets)
2. Bootloader checks firmware integrity
3. If bad: rollback to previous version
4. If good: continue normal boot
5. Report restart reason to cloud

---

## OTA Update Procedure

### Update Flow

```
1. Cloud notifies device via MQTT
2. Device downloads firmware:
   - GET firmware binary from S3/CDN
   - Validate file size
   - Validate signature (RSA 2048 or ECDSA)
   - If invalid: reject and notify
3. Device writes to OTA partition:
   - Write to inactive partition (not running)
   - Save in progress to NVS
4. Verify update:
   - Check CRC of written data
   - If fail: erase and retry
5. Set boot flag:
   - Mark new partition as boot-on-next-restart
   - Save old partition for rollback
6. Restart:
   - Device reboots
   - Bootloader loads new firmware
7. Validation:
   - New firmware starts
   - If boot timeout: automatic rollback
   - If success: send confirmation
```

---

## Error Handling

### Critical Errors (Automatic Restart)

- MQTT client crashed
- Heap corruption detected
- Stack overflow
- Hardware watchdog timeout
- Bootloader detected corruption

### Non-Critical Errors (Retry or Skip)

- DHT22 read timeout (use last value)
- MQTT publish failure (retry next cycle)
- Command timeout (respond with timeout)
- Sensor outlier (filter and discard)

### Error Logging

```c
// Error structure
typedef struct {
  uint32_t timestamp;
  uint8_t error_code;
  char error_message[256];
  uint32_t device_memory_free;
  uint32_t restart_counter;
} error_log_t;

// Store in NVS (circular buffer)
// Send to cloud on successful MQTT connection
```

---

## Memory Management

### Available Memory

- **Flash**: 4MB total (board-dependent)
  - 512KB: Bootloader
  - 1.5MB: App partition
  - 1.5MB: OTA partition
  - 512KB: NVS + Other

- **SRAM**: ~296KB total
  - 160KB: Internal memory
  - 136KB: PSRAM (if equipped)

- **Typical Usage** (2.0 firmware):
  - Code + libraries: 600KB
  - Heap: 100KB available for runtime
  - Stack: 32KB total for all tasks

---

## Debugging & Logging

### Serial Output

```
UART0 @ 115200 baud (8-N-1)
GPIO 1 (TX), GPIO 3 (RX)

Log levels: E (error), W (warn), I (info), D (debug), V (verbose)
```

### Log Format

```
[timestamp] [level] [component] message
[00:01:23.456] [I] [MQTT] Connected to broker
[00:02:45.123] [E] [SENSOR] DHT22 read timeout
```

### Debug Mode

- Enable via menuconfig
- Generates verbose logs
- Increased heap allocation
- Disabled for production builds

---

## Testing Specification

### Unit Testing

**Module Tests**:
- GPIO control (set/get pin state)
- ADC reading (check value range)
- sensor filtering (outlier detection)
- MQTT message parsing
- OTA signature verification

### Integration Testing

- WiFi + MQTT connection
- Command reception and execution
- Sensor publish interval
- Offline schedule execution
- OTA update flow (with rollback)

### Field Testing Scenarios

- Network instability (WiFi dropout)
- Command collision (multiple requests)
- Long-term stability (7-day runtime)
- Firmware update reliability
- Recovery from crashes

---

## Security Considerations

### Certificate Pinning (Optional)

```c
// Pin server certificate to device
esp_tls_cfg_t tls_cfg = {
  .cacert_buf = (uint8_t *)server_cert_pem_start,
  .cacert_len = server_cert_pem_end - server_cert_pem_start,
};
```

### Secure NVS (Future)

```c
// Encrypt sensitive values in NVS
nvs_sec_cfg_t cfg = {...};
nvs_flash_secure_init_partition("nvs", &cfg);
```

---

## Firmware Release Checklist

- [ ] Code compiles with no warnings
- [ ] All unit tests pass
- [ ] Integration tests on hardware pass
- [ ] Memory usage < 80%
- [ ] Heap fragmentation acceptable
- [ ] OTA update tested (including rollback)
- [ ] Watchdog timeout tested
- [ ] Security audit completed
- [ ] Signed with private key
- [ ] Version incremented
- [ ] Release notes documented

---

**Version History**
| Version | Date | Changes |
|---------|------|---------|
| 2.0 | Feb 2026 | Initial embedded spec |

**Related Documents**
- SPECIFICATIONS_ARCHITECTURE.md
- SPECIFICATIONS_DEPLOYMENT.md (OTA section)
- SPECIFICATIONS_SECURITY.md (device security)
