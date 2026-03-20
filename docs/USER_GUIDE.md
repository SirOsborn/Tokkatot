# Tokkatot User Guide

This guide is a practical overview for farmers and installers. It covers the coop setup flow, Raspberry Pi gateway, and ESP32 device wiring expectations.

---

## 1) Farm + Coop Setup (App)
1. **Create Farm** (Admin provides registration key → farmer signs up).
2. **Create Coops** inside the farm (each coop is a control unit).
3. **Assign Devices** to coops via the gateway (auto‑report).

**Note:** Farms are containers. All automation and monitoring are **coop‑level**.

---

## 2) Device Types (Per Coop)
A coop can have some or all of these:

**Actuators (scheduled)**
- `feeder_motor` (relay ON/OFF)
- `conveyor_belt` (relay ON/OFF)
- `fan` (relay ON/OFF)
- `heater` (relay ON/OFF)

**Sensors (monitor only)**
- `temp_humidity`
- `water_level`

If a device is not installed, it is reported as **inactive** by the gateway.

---

## 3) Raspberry Pi 4B Gateway (Required)
The Pi is the **local gateway** that connects ESP32 to the cloud.

**Responsibilities:**
- Poll ESP32 for sensor data
- Post telemetry to cloud
- Fetch schedules + thresholds from cloud
- Execute commands locally on ESP32

**Required endpoints (cloud):**
- `POST /v1/farms/:farm_id/coops/:coop_id/devices/report`
- `POST /v1/farms/:farm_id/coops/:coop_id/telemetry`

**Telemetry payload (example):**
```json
{
  "hardware_id": "ESP32-ABC",
  "timestamp": "2026-03-20T13:25:00Z",
  "sensors": {
    "temperature_c": 30.4,
    "humidity_pct": 62.0,
    "water_level_raw": 1370
  }
}
```

---

## 4) ESP32 Setup (Per Coop)
ESP32 controls relays and reads sensors locally.

**Expected firmware behavior:**
- Expose local endpoints for:
  - Read sensors (temp/humidity, water level)
  - Turn ON/OFF: feeder motor, conveyor, fan, heater
- Report sensor readings to Pi (via polling)
- Execute ON/OFF commands from Pi

**Water system:**
- No pump. Floating valve only.
- Water sensor is **monitor‑only** and used to alert if stuck valve.

### ESP32 GPIO Wiring (Reference)
> Update these if your pinout changes. Current values match the embedded headers.

**Relays / Actuators**
- Conveyor belt relay → `GPIO25` (`CONVEYER_PIN`)
- Fan relay → `GPIO26` (`FAN_PIN`)
- Heater relay → `GPIO14` (`LIGHTBULB_PIN`)  *(heater now replaces bulb)*
- Feeder motor relay → `GPIO27` (`WATERPUMP_PIN`) *(repurposed for feeder motor)*

**Sensors**
- Temp/Humidity (DHT22) → `GPIO32` (`DHT22_PIN`)
- Water level (ADC) → `ADC1_CH7` (`GPIO35`)

**Notes**
- Relays are active‑low in firmware (ON = 0).
- If you change wiring, update: `embedded/main/include/device_control.h` and `embedded/main/include/sensor_manager.h`.

---

## 5) Temperature Thresholds (Per Coop)
Farmers set thresholds per coop:
- If temp < min → **heater ON**
- If temp > max → **fan ON**
- Else both OFF

Thresholds are stored at coop level and executed by the gateway.

---

## 6) Water Alert Rule (Per Coop)
Alert when:
> water level is below half threshold for **1 minute**

The threshold is calibrated during installation and stored on the coop.

---

## 7) Schedules (Per Coop)
Schedules are only for actuators.
You can configure:
- Start time (daily)
- Run for X minutes
- Off for Y minutes
- Repeat (sequence)

The gateway executes sequence steps locally so it still works offline.

---

## 8) Raspberry Pi 4B Installer Checklist
Use this quick checklist during on‑site setup:

1. **Network**
   - Pi connected to stable Wi‑Fi or LAN
   - Internet access confirmed (for cloud sync)
2. **ESP32 Reachable**
   - Ping or HTTP GET to ESP32 local endpoint works
3. **Device Report**
   - Pi posts device list to:  
     `POST /v1/farms/:farm_id/coops/:coop_id/devices/report`
4. **Telemetry**
   - Pi posts telemetry every 1–5 minutes:
     `POST /v1/farms/:farm_id/coops/:coop_id/telemetry`
5. **Schedules**
   - Ensure schedules execute locally even if internet drops
6. **Water Threshold Calibration**
   - Measure water sensor value at **half‑full**
   - Save to coop as `water_level_half_threshold`
7. **Temp Thresholds**
   - Farmer sets `temp_min` and `temp_max` for the coop

---

## References
- `README.md` for local dev setup
- `ARCHITECTURE.md` for system overview
- `AGENT.md` for developer workflow
