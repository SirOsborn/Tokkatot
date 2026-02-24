# 🤖 Tokkatot Automation Use Cases

**Last Updated**: February 24, 2026  
**Purpose**: Real-world automation scenarios for Cambodian poultry farmers  
**Audience**: Developers implementing schedule features, farmers configuring automation

---

## 📋 Overview

Tokkatot supports **4 schedule types** to automate coop operations:

1. **Manual Control** - Direct ON/OFF (no schedule)
2. **time_based** - Turn ON at specific times (cron expressions)
3. **duration_based** - Continuous ON/OFF cycling
4. **condition_based** - Sensor-driven automation (coming soon)

Each type solves specific farmer needs based on equipment characteristics and farming practices.

---

## 🚜 Equipment & Use Cases

### **1. Conveyor Belt (Manure Removal)**

**Equipment Details:**
- **Purpose**: Move chicken feces/manure from under cages to collection area
- **Motor Type**: AC motor with gear reduction, runs along cage rows
- **Power**: 0.5-1 HP, 220V AC
- **Typical Length**: 10-30 meters per coop row
- **Operation**: Scraper/chain mechanism pushes manure to end of coop

#### **Use Case 1A: Always ON (Continuous Operation)**
**Farmer Quote**: _"I just want it to run all the time"_

**Solution**: Manual control - no schedule needed

```json
POST /v1/farms/{farm_id}/devices/{conveyor_device_id}/command
{
  "command_type": "set_relay",
  "command_value": "ON"
}
```

**When to Use**:
- Small coops (< 200 chickens)
- Farmer prefers simplicity
- Equipment has low power consumption
- Manure builds up quickly

**Trade-offs**:
- ❌ Higher electricity costs
- ❌ More wear on motor/belt
- ✅ Simplest setup (no schedule configuration)
- ✅ Always clean, no manual intervention

---

#### **Use Case 1B: Cycling ON/OFF (Electricity Savings)**
**Farmer Quote**: _"Turn ON for 10 minutes, then OFF for 15 minutes, repeat all day to save electricity"_

**Solution**: `duration_based` schedule

```json
POST /v1/farms/{farm_id}/schedules
{
  "name": "Conveyor Auto Cycle",
  "device_id": "{conveyor_device_id}",
  "schedule_type": "duration_based",
  "on_duration": 600,        // 10 minutes = 600 seconds
  "off_duration": 900,       // 15 minutes = 900 seconds
  "action": "set_relay",
  "action_value": "ON",
  "priority": 5,
  "is_active": true
}
```

**How It Works**:
1. Conveyor turns ON for 10 minutes (moves manure)
2. Conveyor turns OFF for 15 minutes (rest period)
3. Repeats indefinitely until schedule is deactivated
4. Next cycle starts immediately after OFF period ends

**When to Use**:
- Medium/large coops (300+ chickens)
- Manure doesn't build up too fast
- Farmer wants to reduce electricity costs
- Equipment can handle frequent starts/stops

**Benefits**:
- ✅ 60% electricity savings (10min ON / 25min total = 40% duty cycle)
- ✅ Less motor wear vs. always-on
- ✅ Set once, runs forever
- ✅ No manual intervention

**Customization Examples**:
- Heavy manure buildup: ON 15min / OFF 10min
- Light usage: ON 5min / OFF 20min
- Night cleaning boost: Create 2 schedules (daytime ON 5/OFF 15, nighttime ON 10/OFF 5)

---

#### **Use Case 1C: Scheduled Cleaning Times**
**Farmer Quote**: _"Turn on conveyor at 6AM, 10AM, 3PM, and 6PM only - runs a bit then stops"_

**Solution**: `time_based` schedule with `action_duration` (auto-turn-off)

```json
POST /v1/farms/{farm_id}/schedules
{
  "name": "Conveyor Scheduled Cleaning",
  "device_id": "{conveyor_device_id}",
  "schedule_type": "time_based",
  "cron_expression": "0 6,10,15,18 * * *",  // 6AM, 10AM, 3PM, 6PM daily
  "action_duration": 1800,   // Auto-turn-off after 30 minutes
  "action": "set_relay",
  "action_value": "ON",
  "priority": 7,
  "is_active": true
}
```

**How It Works**:
1. At 6:00 AM: Conveyor turns ON
2. At 6:30 AM: Conveyor auto-turns-off (1800 seconds later)
3. Waits until 10:00 AM
4. At 10:00 AM: Conveyor turns ON again
5. At 10:30 AM: Auto-turns-off
6. Repeats for 3PM and 6PM

**Cron Expression Guide**:
```
"0 6,10,15,18 * * *"
 │ │              
 │ └─ Hours: 6AM, 10AM, 3PM (15), 6PM (18)
 └─── Minute: 0 (exactly on the hour)
```

**When to Use**:
- Farmer follows strict feeding/cleaning schedule
- Manure removal timed with other activities
- Nighttime noise concerns (no overnight operation)
- Coordination with farm staff shifts

**Benefits**:
- ✅ Predictable operation times
- ✅ Aligns with feeding schedule
- ✅ No operation during farmer's rest time
- ✅ Easy to explain to farm workers

**Customization Examples**:
- Shorter run time: `"action_duration": 600` (10 minutes)
- More frequent: `"cron_expression": "0 */2 * * *"` (every 2 hours)
- Weekends only: `"cron_expression": "0 8,16 * * 0,6"` (Sunday/Saturday, 8AM & 4PM)

---

### **2. Feeder Motor (Spiral Feed Dispenser)**

**Equipment Details:**
- **Type**: Auger/spiral feeder (cylindrical tube with internal screw spiral)
- **Material**: Stainless steel spiral inside PVC pipe
- **Motor**: Torque motor (low RPM, high torque) to rotate spiral
- **Length**: 10-50 meters (runs along entire coop length)
- **Operation**: Spiral rotates → pushes feed pellets from storage bin → drops through holes into feeding troughs
- **Feed Distribution**: Holes every 1-2 meters allow feed to drop into bowls below chicken cages

**Farmer's Description**: _"We have a tube with feed inside, it has holes for the feed to fall into the bowls for chickens. Inside the tube there's a spiral made of stainless steel connected to a motor - when it turns, the feed moves across all the cages till the end where leftover feed drops out."_

#### **Use Case 2A: Simple Timed Feeding**
**Farmer Quote**: _"Turn feeder ON at 6AM, 12PM, 6PM - runs for 15 minutes each time, then stops until next feeding"_

**Solution**: `time_based` schedule with `action_duration`

```json
POST /v1/farms/{farm_id}/schedules
{
  "name": "Daily Feeding Schedule",
  "device_id": "{feeder_device_id}",
  "schedule_type": "time_based",
  "cron_expression": "0 6,12,18 * * *",  // 6AM, 12PM (noon), 6PM daily
  "action_duration": 900,    // 15 minutes = 900 seconds
  "action": "set_relay",
  "action_value": "ON",
  "priority": 10,            // High priority - food is critical
  "is_active": true
}
```

**How It Works**:
1. **6:00 AM**: Motor starts, spiral rotates
2. **6:00-6:15 AM**: Feed pellets move through tube, drop into bowls
3. **6:15 AM**: Motor auto-stops (900 seconds elapsed)
4. **6:15 AM - 12:00 PM**: Motor OFF (chickens eat, no new feed dispensed)
5. **12:00 PM**: Motor starts again
6. **12:15 PM**: Auto-stops
7. Repeats at 6PM

**When to Use**:
- Standard 3x daily feeding schedule (morning, noon, evening)
- Fixed portion sizes (15min rotation = specific amount of feed)
- Chickens eat between feeding times (no continuous grazing)

**Benefits**:
- ✅ Consistent feeding times (chickens adapt to routine)
- ✅ Prevents overfeeding (portion control)
- ✅ Farmer can leave farm between feedings
- ✅ No wasted feed (leftover drops out at tube end)

---

#### **Use Case 2B: Multi-Step Feeding Sequence**
**Farmer Quote**: _"At feeding time, I want: motor ON 30 seconds, PAUSE 10 seconds, motor ON 30 seconds, PAUSE 10 seconds, then big OFF waiting for next scheduled time. This gives chickens time to eat in between feed drops."_

**Solution**: `time_based` schedule with `action_sequence` (multi-step pattern)

```json
POST /v1/farms/{farm_id}/schedules
{
  "name": "Pulse Feeding Sequence",
  "device_id": "{feeder_device_id}",
  "schedule_type": "time_based",
  "cron_expression": "0 6,12,18 * * *",  // 6AM, 12PM, 6PM daily
  "action_sequence": "[{\"action\":\"ON\",\"duration\":30},{\"action\":\"OFF\",\"duration\":10},{\"action\":\"ON\",\"duration\":30},{\"action\":\"OFF\",\"duration\":10}]",
  "action": "set_relay",     // Fallback action (required)
  "action_value": "ON",
  "priority": 10,
  "is_active": true
}
```

**Action Sequence Breakdown**:
```json
[
  {"action": "ON",  "duration": 30},  // Step 1: Feed dispenses for 30 seconds
  {"action": "OFF", "duration": 10},  // Step 2: Pause 10 seconds (chickens approach bowls)
  {"action": "ON",  "duration": 30},  // Step 3: More feed for 30 seconds
  {"action": "OFF", "duration": 10}   // Step 4: Final pause 10 seconds
]
// Total sequence time: 80 seconds (30+10+30+10)
// After sequence completes: Motor stays OFF until next scheduled time (6AM → 12PM)
```

**How It Works**:
1. **6:00:00 AM**: Motor ON (spiral rotates)
2. **6:00:30 AM**: Motor OFF (pause for chickens to approach)
3. **6:00:40 AM**: Motor ON again (second feed burst)
4. **6:01:10 AM**: Motor OFF (final pause)
5. **6:01:20 AM**: Sequence complete → Motor stays OFF
6. **6:01:20 AM - 12:00:00 PM**: 10+ hour rest period
7. **12:00:00 PM**: Sequence repeats

**When to Use**:
- Chickens rush to bowls (need pause between feed bursts)
- Prevent feed congestion at tube holes
- Aggressive eaters dominate first feeding burst
- Better feed distribution across all cages

**Benefits**:
- ✅ Chickens have time to reach feeding bowls
- ✅ More even distribution (slower eaters get 2nd wave)
- ✅ Prevents feed pile-up at first few holes
- ✅ Better utilization of feed (less waste)

**Customization Examples**:
- Longer pauses: `{"action": "OFF", "duration": 20}` (20 seconds to eat)
- More bursts: Add 3rd/4th ON/OFF steps for gradual feeding
- Variable timing:
  ```json
  [
    {"action": "ON",  "duration": 60},   // Initial burst (60 sec)
    {"action": "OFF", "duration": 30},   // Longer pause
    {"action": "ON",  "duration": 30},   // Top-up
    {"action": "OFF", "duration": 15}    // Final pause
  ]
  ```

**Why This Matters**:
- **Normal Feeding**: All feed comes at once → chickens at far end get less → uneven growth
- **Pulse Feeding**: Feed comes in waves → gives slower chickens time → better flock health

---

#### **Use Case 2C: Grazing/Continuous Feed (Layer Chickens)**
**Farmer Quote**: _"My layer chickens eat all day, I want feed to move slowly and continuously during daylight hours"_

**Solution**: `duration_based` schedule with short ON/long OFF cycles

```json
POST /v1/farms/{farm_id}/schedules
{
  "name": "Continuous Grazing Feed",
  "device_id": "{feeder_device_id}",
  "schedule_type": "duration_based",
  "on_duration": 10,         // 10 seconds of feed movement
  "off_duration": 300,       // 5 minutes pause (300 seconds)
  "action": "set_relay",
  "action_value": "ON",
  "priority": 8,
  "is_active": true          // Runs 24/7 unless manually deactivated
}
```

**How It Works**:
1. Motor ON for 10 seconds (small amount of feed moves)
2. Motor OFF for 5 minutes (chickens eat)
3. Repeats continuously (12 feed bursts per hour)

**When to Use**:
- Layer chickens (egg production requires constant nutrition)
- Free-range style feeding (chickens eat when hungry)
- Prevent feed spoilage (small amounts dispensed frequently)

**Benefits**:
- ✅ Mimics natural grazing behavior
- ✅ Fresh feed always available
- ✅ Prevents gorging (slow, steady intake)
- ✅ Better digestion for laying hens

**To Stop Overnight**:
Create 2 schedules:
1. **Daytime Grazing** (6AM-6PM): Same as above, but add cron filter (advanced feature)
2. **Manual Stop**: Deactivate schedule at night, reactivate in morning

---

### **3. Water Pump (Tank Refill)**

**Equipment Details:**
- **Type**: Submersible/centrifugal pump
- **Source**: Well, river, or main water line
- **Destination**: Elevated water tank (gravity-fed to coop pipes)
- **Sensor**: Ultrasonic distance sensor at top of tank (measures water level)
- **Capacity**: 500-2000 liters per tank
- **Fill Time**: 5-30 minutes depending on pump size

#### **Use Case 3A: Sensor-Driven Auto-Refill**
**Farmer Quote**: _"Turn pump ON when water level is low, turn OFF when tank is full"_

**Solution**: `condition_based` schedule (sensor-driven)

```json
POST /v1/farms/{farm_id}/schedules
{
  "name": "Auto Water Refill",
  "device_id": "{water_pump_device_id}",
  "schedule_type": "condition_based",
  "condition_json": "{\"sensor\":\"water_level\",\"operator\":\"<\",\"threshold\":20}",
  "action": "set_relay",
  "action_value": "ON",
  "priority": 10,            // Critical - chickens need water
  "is_active": true
}

// Companion "turn OFF" schedule:
POST /v1/farms/{farm_id}/schedules
{
  "name": "Auto Water Stop",
  "device_id": "{water_pump_device_id}",
  "schedule_type": "condition_based",
  "condition_json": "{\"sensor\":\"water_level\",\"operator\":\">\",\"threshold\":90}",
  "action": "set_relay",
  "action_value": "OFF",
  "priority": 10,
  "is_active": true
}
```

**Condition JSON Explained**:
```json
{
  "sensor": "water_level",       // Device sensor type
  "operator": "<",               // Less than (also supports: >, <=, >=, ==)
  "threshold": 20                // 20% full (low level alarm)
}
```

**How It Works**:
1. Water level drops to 19% → "Auto Water Refill" schedule triggers
2. Pump turns ON
3. Tank fills gradually (5-30 minutes)
4. Water level reaches 90% → "Auto Water Stop" schedule triggers
5. Pump turns OFF
6. System waits for next low-level event

**When to Use**:
- Ultrasonic sensor installed on tank
- Chickens have variable water consumption (weather-dependent)
- Farmer not present to monitor water levels
- Backup automation (prevents empty tank emergencies)

**Benefits**:
- ✅ Zero manual intervention
- ✅ Prevents water waste (no overflow)
- ✅ Never runs dry (chickens always have water)
- ✅ Adapts to consumption patterns (hot days = more refills)

**Safety Features** (to be implemented):
- Max runtime limit: `"action_duration": 1800` (30min timeout prevents pump damage if sensor fails)
- Minimum OFF time: Prevent rapid cycling (schedule cooldown logic)

---

#### **Use Case 3B: Scheduled Refill (No Sensor)**
**Farmer Quote**: _"I don't have a sensor yet, just fill the tank at 6AM and 6PM for 20 minutes"_

**Solution**: `time_based` schedule with `action_duration`

```json
POST /v1/farms/{farm_id}/schedules
{
  "name": "Morning & Evening Tank Fill",
  "device_id": "{water_pump_device_id}",
  "schedule_type": "time_based",
  "cron_expression": "0 6,18 * * *",  // 6AM and 6PM daily
  "action_duration": 1200,   // 20 minutes = 1200 seconds
  "action": "set_relay",
  "action_value": "ON",
  "priority": 9,
  "is_active": true
}
```

**How It Works**:
1. **6:00 AM**: Pump ON
2. **6:20 AM**: Pump auto-OFF (1200 seconds elapsed)
3. **6:20 AM - 6:00 PM**: Pump OFF
4. **6:00 PM**: Pump ON again
5. **6:20 PM**: Pump auto-OFF

**When to Use**:
- No water level sensor available
- Fixed consumption patterns (farmer knows tank empties overnight)
- Temporary solution before sensor installation
- Backup schedule (if sensor fails)

**Risks**:
- ⚠️ May overflow if tank not empty
- ⚠️ May underfill if chickens drink more than expected
- ⚠️ Requires farmer to monitor occasionally

**Mitigation**:
- Add safety schedule: Pump should never run > 30 minutes
- Manual override available via app

---

### **4. Temperature Control (Fan + Heater)**

#### **Use Case 4A: Simple Threshold Automation**
**Farmer Quote**: _"Turn fan ON when temperature > 32°C, turn heater ON when < 18°C"_

**Solution**: `condition_based` schedules (2 separate schedules)

```json
// High temperature → Fan ON
POST /v1/farms/{farm_id}/schedules
{
  "name": "Auto Fan Cooling",
  "device_id": "{fan_device_id}",
  "schedule_type": "condition_based",
  "condition_json": "{\"sensor\":\"temperature\",\"operator\":\">\",\"threshold\":32}",
  "action": "set_relay",
  "action_value": "ON",
  "priority": 8,
  "is_active": true
}

// Low temperature → Heater ON
POST /v1/farms/{farm_id}/schedules
{
  "name": "Auto Heater Warming",
  "device_id": "{heater_device_id}",
  "schedule_type": "condition_based",
  "condition_json": "{\"sensor\":\"temperature\",\"operator\":\"<\",\"threshold\":18}",
  "action": "set_relay",
  "action_value": "ON",
  "priority": 9,             // Higher priority (cold is more dangerous)
  "is_active": true
}
```

**How It Works**:
1. Temperature sensor reads 33°C → Fan schedule triggers, fan ON
2. Temperature drops to 31°C → (Fan schedule no longer triggered, but no auto-OFF yet - **to be implemented: hysteresis logic**)
3. Temperature reads 17°C → Heater schedule triggers, heater ON
4. Temperature rises to 19°C → Heater turns OFF

**When to Use**:
- Simple climate control needs
- Day/night temperature swings
- Seasonal heating/cooling (disable schedule off-season)

**Limitations** (Future Enhancement):
- No hysteresis (fan constantly switches ON/OFF at 32°C boundary)
- Solution: Add `"threshold_high": 32, "threshold_low": 30` (ON at 32°C, OFF at 30°C)

---

## 📊 Schedule Type Selection Guide

| **Farmer Wants** | **Schedule Type** | **Key Field** | **Example** |
|---|---|---|---|
| "Always ON/OFF" | Manual control (no schedule) | - | Direct command |
| "ON 10min, OFF 15min, repeat forever" | `duration_based` | `on_duration`, `off_duration` | Conveyor cycling |
| "ON at 6AM, 12PM, 6PM for 15min each" | `time_based` | `cron_expression`, `action_duration` | Feeder schedule |
| "ON at 6AM: 30sec ON, 10sec pause, 30sec ON, 10sec pause, then stop" | `time_based` | `cron_expression`, `action_sequence` | Pulse feeding |
| "ON when temp > 32°C" | `condition_based` | `condition_json` | Auto fan |
| "ON when water < 20%, OFF when > 90%" | `condition_based` (2 schedules) | `condition_json` | Auto pump |

---

## 🛠️ Implementation Notes for Developers

### **action_sequence JSON Format**

```json
{
  "action_sequence": "[{\"action\":\"ON\",\"duration\":30},{\"action\":\"OFF\",\"duration\":10},{\"action\":\"ON\",\"duration\":30},{\"action\":\"OFF\",\"duration\":10}]"
}
```

**Structure**:
- Array of step objects
- Each step: `{"action": "ON"|"OFF", "duration": <seconds>}`
- Steps execute sequentially
- After last step completes, device returns to initial state (OFF)
- If schedule triggers again before sequence completes, current sequence aborts and restarts

**Validation Rules**:
- Maximum 20 steps per sequence (prevent absurd configs)
- Minimum duration: 1 second
- Maximum duration per step: 3600 seconds (1 hour)
- Total sequence time should be < next cron trigger interval

### **Cron Expression Examples**

```
"0 6 * * *"          // Every day at 6:00 AM
"0 6,12,18 * * *"    // Every day at 6AM, 12PM, 6PM
"0 */2 * * *"        // Every 2 hours (0:00, 2:00, 4:00, ...)
"30 8 * * 1-5"       // Monday-Friday at 8:30 AM (weekday feeding)
"0 7 * * 0,6"        // Sunday & Saturday at 7:00 AM (weekend schedule)
"0 6-18 * * *"       // Every hour from 6AM to 6PM (NOT recommended - use duration_based)
```

**Cron Parser**: Uses `robfig/cron/v3` library (Go)
- Format: `Minute Hour Day Month Weekday`
- Ranges: `1-5` (Monday to Friday)
- Lists: `1,3,5` (Monday, Wednesday, Friday)
- Steps: `*/2` (every 2 hours)

### **Database Schema**

```sql
CREATE TABLE schedules (
  id UUID PRIMARY KEY,
  farm_id UUID NOT NULL,
  coop_id UUID,
  device_id UUID NOT NULL,
  name TEXT NOT NULL,
  schedule_type TEXT CHECK (schedule_type IN ('time_based', 'duration_based', 'condition_based')),
  
  -- time_based fields
  cron_expression TEXT,              -- "0 6,12,18 * * *"
  action_duration INTEGER,           -- Seconds (auto-turn-off after X seconds)
  action_sequence JSONB,             -- Array of {action, duration} steps
  
  -- duration_based fields
  on_duration INTEGER,               -- Seconds ON
  off_duration INTEGER,              -- Seconds OFF
  
  -- condition_based fields
  condition_json JSONB,              -- {"sensor":"temp","operator":">","threshold":32}
  
  -- Common fields
  action TEXT NOT NULL,              -- "on", "off", "set_value"
  action_value TEXT,                 -- Optional value (e.g., PWM duty cycle)
  priority INTEGER DEFAULT 0,        -- 0-10 (higher = more important)
  is_active BOOLEAN DEFAULT true,
  
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);
```

---

## 🎓 Farmer Education (UI Hints)

When building the frontend, provide these guidance texts:

### **duration_based Schedule**
> **💡 Best for**: Equipment that should cycle ON and OFF automatically  
> **Example**: Conveyor belt - ON for 10 minutes, OFF for 15 minutes, repeats all day

### **time_based Schedule with action_duration**
> **💡 Best for**: Equipment that runs at specific times, then stops  
> **Example**: Feeder - Turns ON at 6AM, 12PM, 6PM for 15 minutes each time

### **time_based Schedule with action_sequence**
> **💡 Best for**: Equipment that needs pauses during operation  
> **Example**: Feeder - ON 30 seconds, PAUSE 10 seconds, ON 30 seconds, PAUSE 10 seconds at each feeding time

### **condition_based Schedule**
> **💡 Best for**: Equipment controlled by sensors (temperature, water level)  
> **Example**: Water pump - Turns ON when tank is low, OFF when full

---

## 📞 Farmer Support Scenarios

**Q: "My conveyor runs all the time, electricity bill too high!"**  
**A**: Switch from manual control to `duration_based` schedule:
- ON 5 minutes, OFF 15 minutes = 75% electricity savings

**Q: "Chickens fight over feed when feeder starts"**  
**A**: Use `action_sequence` for pulse feeding:
- ON 30sec, PAUSE 10sec, ON 30sec, PAUSE 10sec
- Gives slower chickens time to reach bowls

**Q: "Water tank overflows sometimes"**  
**A**: Install ultrasonic sensor + `condition_based` schedule:
- Pump ON when < 20% full
- Pump OFF when > 90% full
- No overflow, no manual monitoring

**Q: "I want feeder to run at 7AM, but equipment installation team set it for 6AM"**  
**A**: Edit schedule → Change cron expression from `"0 6 * * *"` to `"0 7 * * *"`

---

## 🔄 Schedule Lifecycle

```
1. Farmer creates schedule via app/web
   ↓
2. API validates fields (cron expression, durations, etc.)
   ↓
3. Schedule saved to database
   ↓
4. Backend scheduler loads active schedules
   ↓
5. Cron/timer triggers schedule execution
   ↓
6. Device command sent to Raspberry Pi via MQTT
   ↓
7. Raspberry Pi controls relay → equipment turns ON/OFF
   ↓
8. Execution logged to schedule_executions table
   ↓
9. Real-time update sent to app via WebSocket
   ↓
10. Farmer sees confirmation on phone
```

---

## ✅ Summary

Tokkatot's scheduling system gives Cambodian farmers **professional-grade automation** with **simple, farmer-friendly configuration**:

- ✅ **4 schedule types** cover all common use cases
- ✅ **No coding required** - simple form inputs
- ✅ **Electricity savings** - duration_based cycling reduces costs
- ✅ **Better animal welfare** - pulse feeding, auto water refill
- ✅ **Set-and-forget** - runs indefinitely until manually stopped
- ✅ **Real-time control** - manual override always available

This document should be referenced when:
- **Designing UI**: Show relevant examples for each equipment type
- **Writing API docs**: Link to use cases for context
- **Training staff**: Explain why farmers need each feature
- **Customer support**: Troubleshoot farmer's specific automation needs

---

**End of Document**

---

## 🌡️ Temperature Monitoring Dashboard

**Last Updated**: February 24, 2026  
**Feature**: Apple Weather-style temperature timeline  
**Endpoint**: `GET /v1/farms/:farm_id/coops/:coop_id/temperature-timeline?days=7`  
**Frontend**: `/monitoring` (`pages/monitoring.html`)

### Farmer Problem

**Farmer Quote**: _"In the hot season, I need to know how hot my coop got today and whether it was dangerous for my chickens — without reading numbers"_

Cambodian farmers deal with extreme heat (30–42°C in dry season). High coop temperatures cause:
- Heat stress (reduced weight gain, egg production drops)
- Sudden death in broilers above 38°C
- Cascading flock illness

Farmers need to **understand temperature at a glance**, not read data tables.

### Solution: Visual Temperature Timeline

**What it shows:**
- Current temperature as a large number (like a weather app)
- Background gradient colour that instinctively signals danger (scorching red = danger, cool blue = fine)
- Today’s highest and lowest temperature **with the exact time they occurred** (e.g. “H: 38.5° at 14:00”)
- Scrollable hourly strip — tap to see what happened hour by hour
- Smooth SVG temperature curve for the full day
- Daily history for the past week (Yesterday, Mon, Tue…)

### bg_hint Farmer Impact

| Colour hint | Temperature | What the farmer should do |
|---|---|---|
| `scorching` | ≥ 35°C | Open vents immediately, turn on fans, reduce feeding |
| `hot`       | ≥ 32°C | Monitor closely, consider fan schedule |
| `warm`      | ≥ 28°C | Normal summer day — no action needed |
| `neutral`   | ≥ 24°C | Ideal temperature range |
| `cool`      | ≥ 20°C | Fine for most breeds |
| `cold`      | < 20°C  | Consider heater for brooder chicks |

### Technical Notes
- **No humidity**: `sensor_type = 'temperature'` filter only — humidity removed from all queries
- **No crash on missing sensor**: Returns `sensor_found: false` with HTTP 200 (farmers with no sensor still see the page)
- **Per-coop scope**: Farmer selects coop via dropdown; each coop has its own timeline
- **Days parameter**: `?days=7` (default), max 30
