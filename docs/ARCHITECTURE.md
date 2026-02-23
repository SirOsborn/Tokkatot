# Architecture - Tokkatot 2.0

**Last Updated**: February 23, 2026  
**Type**: System Design

---

## System Overview

**Tokkatot** = IoT dashboard for Cambodian chicken farmers  
**Core Principle**: **Coop-centric** - All devices belong to specific coops  
**Deployment**: Single VPS (not microservices)  
**Users**: Budget farmers (simple UI, phone-based)

---

## Physical Infrastructure

### Chicken Coop Layout
```
┌─────────────────────────────────────────┐
│  Coop Building                          │
│  ┌─────┐ ┌─────┐   ┌─────┐ ┌─────┐      │  ← Chicken cages
│  └──┬──┘ └──┬──┘   └──┬──┘ └──┬──┘      │
│  🍴─────────────────────────────→       │  ← feeds pipes
│  💧─────────────────────────────→       │  ← Water pipes   
│  ═══╧═══════╧═════════╧═══════╧═══      │  ← Manure belt
│                            📹 AI Camera │  ← Disease detection
└─────────────────────────────────────────┘
         │
    ┌────┴─────┐
    │Water Tank│  ← Ultrasonic sensor
    │  (95%)   │
    └────┬─────┘
    ⚙️ Pump      ← Dedicated per coop
```

### Devices Per Coop
| Device | Purpose | Connection |
|--------|---------|------------|
| Raspberry Pi | Main controller | WiFi → Cloud |
| AI Camera | Feces monitoring | USB → Raspberry Pi |
| Ultrasonic Sensor | Water level (0-100%) | GPIO → Raspberry Pi |
| Water Pump | Auto-fill tank | Relay → Raspberry Pi |
| Conveyor Motor | Manure removal | Relay → Raspberry Pi |

**Key**: Each coop = Independent unit (Coop 1's pump ≠ Coop 2's pump)

---

## Data Hierarchy

```
User (Farmer)
  📱 phone: "012345678"  ← Login ID
  🔑 password: hashed
  🌐 language: "km" (Khmer)
  │
  ├─ Farm 1 (Kandal Province)
  │   ├─ Coop 1 (500 chickens)
  │   │   ├─ Raspberry Pi (main_controller)
  │   │   ├─ AI Camera
  │   │   ├─ Water Sensor
  │   │   ├─ Water Pump
  │   │   └─ Conveyor Motor
  │   │
  │   └─ Coop 2 (300 chickens)
  │       └─ (own devices)
  │
  └─ Farm 2 (Kampong Cham)
      └─ Coop 1
          └─ (own devices)
```

### Database Schema (8 Tables)
```sql
users              -- Farmers (phone login)
  ↓
farms              -- Locations (Kandal, Kampong Cham)
  ↓
coops              -- Chicken houses (number, capacity, chicken_type)
  ↓
devices            -- IoT hardware (coop_id, is_main_controller)
  ↓
device_commands    -- Control actions (on/off, coop_id)
schedules          -- Automated tasks (coop_id)
event_logs         -- History (device actions, alerts)
registration_keys  -- On-site verification (farm_name, expires_at)
```

**Indexes**: 16+ for fast queries on `(farm_id, coop_id)`, `(coop_id, is_main_controller)`

---

## System Architecture

### Deployment (Single VPS)
```
┌──────────────────────────────────────┐
│  VPS ($10/month)                     │
│  Ubuntu 22.04 LTS                    │
│                                      │
│  ┌────────────────────────────────┐  │
│  │ Go Backend (systemd service)   │  │  ← Port 3000
│  │ - Authentication (JWT)         │  │
│  │ - Device API                   │  │
│  │ - WebSocket (real-time)        │  │
│  └────────────────────────────────┘  │
│                                      │
│  ┌────────────────────────────────┐  │
│  │ Python AI Service (Docker)     │  │  ← Port 8000
│  │ - PyTorch ensemble model       │  │
│  │ - Disease detection            │  │
│  └────────────────────────────────┘  │
│                                      │
│  ┌────────────────────────────────┐  │
│  │ PostgreSQL 17                  │  │  ← Port 5432
│  │ - Main database                │  │
│  │ - JSONB, full-text search      │  │
│  └────────────────────────────────┘  │
│                                      │
│  ┌────────────────────────────────┐  │
│  │ Caddy (Reverse Proxy)          │  │  ← Port 443 (HTTPS)
│  │ - Auto SSL/TLS                 │  │
│  │ - Static files (Vue.js)        │  │
│  └────────────────────────────────┘  │
└──────────────────────────────────────┘
         ↕ HTTPS/WSS
┌──────────────────────────────────────┐
│  Farmer's Phone (Vue.js 3 PWA)       │
│  - Offline-capable                   │
│  - Khmer language                    │
│  - Touch-optimized (48px targets)    │
└──────────────────────────────────────┘
         ↕ MQTT over WiFi
┌──────────────────────────────────────┐
│  Raspberry Pi (per coop)             │
│  - Python controller                 │
│  - MQTT broker                       │
│  - Sensor aggregation                │
└──────────────────────────────────────┘
```

### Communication Flow
```
1. Sensor Reading:
   Ultrasonic Sensor → Raspberry Pi → Cloud Backend → Database
                                    ↓
                              WebSocket → Farmer's Phone UI

2. Device Command:
   Farmer's Phone → Backend API → Database (log command)
                                ↓
                          Raspberry Pi (MQTT)
                                ↓
                          Water Pump ON

3. AI Detection:
   AI Camera → Raspberry Pi → Backend → AI Service (PyTorch)
                            ↓               ↓
                        Database ←─── Disease: "Coccidiosis" (87%)
                            ↓
                        Alert → Farmer's Phone
```

---

## User Journey

### 1. Login
```
┌──────────────────┐
│  Login Screen    │
│ Phone: 012...    │
│ Password: ****   │
│ [  ចូល / Login ] │
└──────────────────┘
```

### 2. Select Farm (if multiple)
```
┌─────────────────┐
│ Kandal Farm     │  ← Tap
│ 📍 Kandal       │
│ 🐔 2 Coops      │
├─────────────────┤
│ Kampong Cham    │
│ 📍 Kampong Cham │
│ 🐔 1 Coop       │
└─────────────────┘
```

### 3. Select Coop
```
┌─────────────────┐
│ Coop 1          │  ← Tap
│ 🐔 480/500      │
│ 💧 95% ✅      │
│ 🌡️ 28°C         │
├─────────────────┤
│ Coop 2          │
│ 🐔 290/300     │
│ 💧 20% ⚠️      │  ← Low water alert!
│ 🌡️ 29°C        │
└─────────────────┘
```

### 4. Coop Dashboard (Real-Time)
```
┌──────────────────────┐
│ Coop 1 - Kandal Farm │
├──────────────────────┤
│ 📸 Live Feed         │
│ [feces image]        │
├──────────────────────┤
│ 💧 Water: 95% ✅    │
│ ████████████████▒▒   │
├──────────────────────┤
│ ⚙️ Water Pump: OFF   │
│ [ Turn ON ]          │  ← Manual control
├──────────────────────┤
│ 🌡️ Temp: 28°C        │
│ 💨 Humidity: 65%     │
├──────────────────────┤
│ 🚨 Alerts            │
│ • Disease detected ⚠️│
│   (2 hours ago)      │
└──────────────────────┘
```

---

## Device Control Logic

### Example: Water Level Monitoring
```python
# Raspberry Pi (runs 24/7)
while True:
    water_level = read_ultrasonic_sensor()  # 0-100%
    
    if water_level < 20:  # Low!
        # Find this coop's pump in database
        pump = find_device(coop_id=current_coop, type='pump')
        
        # Turn on pump
        send_command(device_id=pump.id, action='on')
        
        # Wait until tank full
        while water_level < 90:
            sleep(30)
            water_level = read_ultrasonic_sensor()
        
        # Stop pump
        send_command(device_id=pump.id, action='off')
        
        # Log event
        log_event(coop_id, 'water_fill', '20% → 90%')
    
    sleep(60)  # Check every minute
```

### Example: Disease Detection
```python
# Every 30 minutes
if time.now().minute in [0, 30]:
    # Capture feces image
    image = ai_camera.capture()
    
    # Send to AI service
    result = requests.post(
        'http://localhost:8000/predict',
        files={'image': image}
    )
    
    if result['confidence'] > 0.8:
        # High confidence disease detected!
        create_alert(
            coop_id=current_coop,
            type='disease',
            disease=result['disease'],
            confidence=result['confidence'],
            severity='high'
        )
        
        # Send push notification
        send_notification(
            user_id=owner_id,
            title=f"⚠️ Alert: {coop.name}",
            message=f"Potential {result['disease']} detected"
        )
```

---

## Authentication & Security

### Registration (On-Site by Staff)
```
1. Staff visits farm
2. Installs Raspberry Pi, sensors in each coop
3. Generates registration key:
   ./generate_reg_key.ps1 -FarmName "Sokha's Farm" -Phone "012345678"
   → Key: ABCDE-FGHIJ-KLMNO-PQRST-UVWXY

4. Farmer creates account on phone:
   POST /v1/auth/signup
   {
     "phone": "012345678",
     "password": "Farmer123",
     "registration_key": "ABCDE-FGHIJ-KLMNO-PQRST-UVWXY"
   }

5. Backend validates key:
   - Check: key not used
   - Check: key not expired (90 days)
   - If valid: contact_verified = true ✅
   - Mark key as used
```

### Login (Farmer Self-Service)
```
POST /v1/auth/login
{
  "phone": "012345678",
  "password": "Farmer123"
}

Response:
{
  "access_token": "eyJhbG...",  ← 24h expiry
  "refresh_token": "eyJhbG...", ← 30d expiry
  "user": { "id": "...", "name": "Sokha" },
  "farms": [
    {
      "id": "farm-uuid",
      "name": "Kandal Farm",
      "coop_count": 2
    }
  ]
}
```

### Role-Based Access Control (RBAC)
| Role | Permissions |
|------|-------------|
| **Owner** | Full control (create/delete coops, manage users) |
| **Manager** | Control devices, view data, edit schedules |
| **Viewer** | Read-only (dashboards, no device control) |

---

## Performance Requirements

### Response Times (Low-End Android Targets)
- Page load: < 3 seconds
- API response: < 500ms
- Device command: < 2 seconds
- AI prediction: < 3 seconds (CPU fallback)

### Database Optimization
```sql
-- Composite indexes for fast coop queries
CREATE INDEX idx_coops_farm_id ON coops(farm_id);
CREATE INDEX idx_devices_coop_id ON devices(coop_id);
CREATE INDEX idx_devices_main ON devices(coop_id, is_main_controller);

-- Connection pooling
db.SetMaxOpenConns(25)  -- Max concurrent
db.SetMaxIdleConns(5)   -- Keep ready
```

### Scalability Limits
- Current VPS: 100-200 farms (~1000-2000 coops)
- Database: 10k+ IoT readings/second
- WebSocket: 500 concurrent connections
- When to scale: > 500 farms → Add second VPS + load balancer

---

## Offline/Fault Tolerance

### Raspberry Pi Offline Mode
```
Internet Down
  ↓
Raspberry Pi continues:
  ✅ Sensor reading (local)
  ✅ Automated pump control (local rules)
  ✅ Data logging (local SQLite)
  ❌ Disease detection (needs cloud AI)
  ❌ Remote control from phone (needs cloud)

Internet Restored
  ↓
Raspberry Pi syncs:
  1. Upload queued sensor readings
  2. Upload event logs
  3. Download pending commands
  4. Resume disease detection
```

### Farmer Offline Mode (PWA)
```
PWA Service Worker:
  ✅ Cache last dashboard state
  ✅ Show last known values
  ❌ Cannot send device commands
  ❌ Cannot see real-time updates

Internet Restored:
  → Auto-refresh data
  → Enable controls
```

---

## Tech Stack Summary

| Component | Technology | Why |
|-----------|-----------|-----|
| Backend | Go 1.23 + Fiber v2 | Fast, single binary, low memory |
| Frontend | Vue.js 3 (CDN) | Reactive UI, 40KB, progressive |
| Database | PostgreSQL 17 | ACID, JSONB, reliable |
| AI Service | Python 3.12 + PyTorch | Best ML ecosystem |
| Embedded | ESP32 + Raspberry Pi | Cheap ($3-30), WiFi, GPIO |
| Deployment | Docker + systemd | Simple, VPS-friendly |

---

## Next Steps

**For Developers:**
1. Read [TECH_STACK.md](TECH_STACK.md) for detailed technology choices
2. Read [guides/SETUP.md](guides/SETUP.md) for installation
3. Read [implementation/API.md](implementation/API.md) for backend development
4. Read [implementation/FRONTEND.md](implementation/FRONTEND.md) for Vue.js migration

**For AI Agents:**
- This is the **canonical architecture** document
- All implementation must follow coop-centric design
- Phone-based login is non-negotiable (farmers don't use email)
- Registration key system is FREE (no SMS/email costs)
