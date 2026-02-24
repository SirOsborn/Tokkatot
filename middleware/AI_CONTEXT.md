# 🤖 AI Context: Go API Gateway & Middleware

**Directory**: `middleware/`  
**Your Role**: HTTP API, authentication, device management, schedule automation, alerts, analytics  
**Tech Stack**: Go 1.23+, Fiber v2, PostgreSQL (only), JWT, WebSocket  

**📖 Read First**: `../AI_INSTRUCTIONS.md` (project overview, business model, farmer context)

**📖 Full Documentation**:
- API Spec: `../docs/implementation/API.md` (66 endpoints, all MVP endpoints implemented)
- Database: `../docs/implementation/DATABASE.md` (14 tables, PostgreSQL only)
- Security: `../docs/implementation/SECURITY.md` (JWT, registration keys)
- Automation: `../docs/AUTOMATION_USE_CASES.md` (real-world farmer scenarios)

---

## 🎯 What This Component Does

**RESTful API Gateway** serving 66+ endpoints:
- **Authentication**: Phone/Email login, JWT tokens, registration keys
- **Farms & Coops**: Hierarchical management (Farm → Coop → Device)
- **Farm Members**: Invite, role management, remove members
- **Device Control**: Real-time commands via WebSocket; CRUD + advanced (config, calibrate, history)
- **Schedules**: 4 automation types (time_based, duration_based, condition_based, manual)
- **Alerts**: Farm alert monitoring + user notification subscriptions
- **Analytics**: Dashboard, device metrics reports, CSV export, event log
- **User Sessions**: View and revoke active sessions
- **AI Proxy**: Disease detection (stubbed — Coming Soon overlay; not wired for MVP)
- **WebSocket**: Real-time updates for device status, alerts, commands

**Database**: PostgreSQL only (SQLite permanently removed Feb 2026)

---

## 📁 File Structure

```
middleware/
├── main.go                    Entry point, Fiber server setup, routes (66+ routes)
├── go.mod/go.sum             Go 1.23 dependencies (no sqlite3)
├── .env                      DATABASE_URL, JWT_SECRET (GITIGNORE'D)
│
├── api/                      Endpoint handlers
│   ├── auth_handler.go       Login, signup, token refresh, verify, forgot/reset password (8 endpoints)
│   ├── auth_middleware.go    JWT extraction & farm access checks (checkFarmAccess)
│   ├── farm_handler.go       Farm CRUD + member management (9 endpoints)
│   ├── coop_handler.go       Coop CRUD (5 endpoints)
│   ├── device_handler.go     Device control, commands, CRUD, advanced ops (17 endpoints)
│   ├── schedule_handler.go   Schedule automation (7 endpoints) ⭐
│   ├── alert_handler.go      Farm alerts + user subscriptions (8 endpoints) ⭐ NEW
│   ├── analytics_handler.go  Dashboard, reports, CSV export (6 endpoints) ⭐ NEW
│   ├── user_handler.go       User profile, sessions, activity log (6 endpoints)
│   └── websocket_handler.go  Real-time updates + /ws/stats
│
├── database/
│   └── postgres.go           PostgreSQL schema (14 tables + indexes) — ONLY DB file
│
├── models/
│   └── models.go             Data structures (User, Farm, Coop, Device, Schedule,
│                              Alert, AlertSubscription, UserSession,
│                              DeviceConfiguration, DeviceReading)
│
└── utils/
    └── utils.go              JWT, bcrypt, error responses
```

---

## 🚜 Schedule Automation (CRITICAL - New Feb 2026)

**Purpose**: Automate conveyor belts, feeders, water pumps, climate control for farmers

### 4 Schedule Types

| Type | Use Case | Key Fields | Example |
|------|----------|-----------|---------|
| **time_based** | Trigger at specific times | `cron_expression` + (`action_duration` OR `action_sequence`) | Feeder ON at 6AM, 12PM, 6PM for 15min |
| **duration_based** | Continuous ON/OFF cycling | `on_duration`, `off_duration` | Conveyor ON 10min, OFF 15min, repeat |
| **condition_based** | Sensor-driven | `condition_json` | Pump ON when water < 20%, OFF when > 90% |
| **Manual** | Direct control | None | Farmer manually turns ON/OFF via app |

### New Features (Feb 2026)

**1. `action_duration` field** - Auto-turn-off after X seconds
```json
{
  "schedule_type": "time_based",
  "cron_expression": "0 6,12,18 * * *",  // 6AM, 12PM, 6PM
  "action_duration": 900,  // Auto-OFF after 15 minutes
  "action": "set_relay",
  "action_value": "ON"
}
```
**Farmer benefit**: "Feed chickens at 6AM for 15 minutes, then auto-stop"

**2. `action_sequence` field** - Multi-step patterns
```json
{
  "schedule_type": "time_based",
  "cron_expression": "0 6,12,18 * * *",
  "action_sequence": "[
    {\"action\":\"ON\",\"duration\":30},
    {\"action\":\"OFF\",\"duration\":10},
    {\"action\":\"ON\",\"duration\":30},
    {\"action\":\"OFF\",\"duration\":10}
  ]",
  "action": "set_relay"
}
```
**Farmer benefit**: "Pulse feeding - motor ON 30sec, pause 10sec (chickens approach bowls), ON 30sec, pause 10sec"

### Real-World Examples

See `../docs/AUTOMATION_USE_CASES.md` for detailed JSON examples:
- **Conveyor Cycling**: ON 10min, OFF 15min, repeat (60-75% electricity savings)
- **Pulse Feeding**: Multi-step sequences prevent chicken congestion at feed bowls
- **Sensor Pumps**: Auto refill water tank based on ultrasonic sensor readings
- **Climate Control**: Fan ON when temp > 32°C, heater ON when temp < 18°C

### Schedule Endpoints

**File**: `api/schedule_handler.go` (7 endpoints)
```
POST   /v1/farms/:farm_id/schedules                          Create schedule
GET    /v1/farms/:farm_id/schedules                          List schedules
GET    /v1/farms/:farm_id/schedules/:schedule_id             Get single schedule
PUT    /v1/farms/:farm_id/schedules/:schedule_id             Update schedule
DELETE /v1/farms/:farm_id/schedules/:schedule_id             Delete schedule
GET    /v1/farms/:farm_id/schedules/:schedule_id/executions  Execution history
POST   /v1/farms/:farm_id/schedules/:schedule_id/execute-now Manual trigger
```

### Database Schema

**Table**: `schedules` (in `database/postgres.go` and `database/sqlite.go`)

Key fields:
- `schedule_type`: time_based | duration_based | condition_based
- `cron_expression`: "0 6,12,18 * * *" (minute hour day month weekday)
- `on_duration` / `off_duration`: Seconds for duration_based cycling
- `action_duration`: Seconds before auto-turn-off (time_based only)
- `action_sequence`: JSONB array of {action, duration} steps
- `condition_json`: Sensor rules like {"sensor":"water_level","operator":"<","threshold":20}
- `priority`: 0-10, higher = more important (for conflict resolution)

**See**: `../docs/implementation/DATABASE.md` (schedules table section)

## 🖐️ Farm Member Endpoints (NEW Feb 2026)

**File**: `api/farm_handler.go` (4 member endpoints appended after farm CRUD)
```
GET    /v1/farms/:farm_id/members                      List all farm members + roles
POST   /v1/farms/:farm_id/members                      Invite user by email/phone (manager+)
PUT    /v1/farms/:farm_id/members/:user_id             Change member role (owner only)
DELETE /v1/farms/:farm_id/members/:user_id             Remove member (owner only, can't remove owners)
```

**Rules**:
- Cannot change your own role
- Cannot set anyone to `owner` role (only one via farm creation)
- Cannot remove a user with role `owner`
- Invite by email OR phone; user must already exist in system

---

## 🔔 Alert & Subscription Endpoints (NEW Feb 2026)

**File**: `api/alert_handler.go` (8 endpoints)
```
GET /v1/farms/:farm_id/alerts/history        Alert history (days param, includes duration_minutes)
GET /v1/farms/:farm_id/alerts                Active alerts (filter: is_active, severity)
GET /v1/farms/:farm_id/alerts/:alert_id      Single alert
PUT /v1/farms/:farm_id/alerts/:alert_id/acknowledge  Mark alert as acknowledged

POST   /v1/users/alert-subscriptions         Create/upsert notification preference
GET    /v1/users/alert-subscriptions         List user's subscriptions
PUT    /v1/users/alert-subscriptions/:id     Update subscription (quiet hours, enabled)
DELETE /v1/users/alert-subscriptions/:id     Delete subscription
```

**Route ordering important**: `/alerts/history` is registered BEFORE `/alerts/:alert_id` in `main.go` to prevent Fiber from treating "history" as an alert_id.

---

## 📊 Analytics & Reporting Endpoints (NEW Feb 2026)

**File**: `api/analytics_handler.go` (6 endpoints)
```
GET /v1/farms/:farm_id/dashboard                          Overview: devices, alerts, recent events
GET /v1/farms/:farm_id/reports/device-metrics             Aggregated sensor data (avg/min/max by day)
GET /v1/farms/:farm_id/reports/device-usage               Device reliability & command success rate
GET /v1/farms/:farm_id/reports/farm-performance           Farm-wide uptime & automation efficiency
GET /v1/farms/:farm_id/reports/export                     CSV download (device_metrics or farm_events)
GET /v1/farms/:farm_id/events                             Audit event log
```

**Export handler**: Sets `Content-Type: text/csv` and `Content-Disposition: attachment` directly on the response. Two types: `device_metrics` (requires `device_id`) and `farm_events`.

---

## 🔧 Device Advanced Endpoints (NEW Feb 2026)

**File**: `api/device_handler.go` (appended after existing handlers)
```
POST   /v1/farms/:farm_id/devices                          Add device (owner only)
PUT    /v1/farms/:farm_id/devices/:device_id               Update name/location (manager+)
DELETE /v1/farms/:farm_id/devices/:device_id               Soft-delete (owner only)

GET    /v1/farms/:farm_id/devices/:device_id/history       Sensor readings (hours, metric, limit)
GET    /v1/farms/:farm_id/devices/:device_id/status        Real-time status + latest reading
GET    /v1/farms/:farm_id/devices/:device_id/config        Config parameters
PUT    /v1/farms/:farm_id/devices/:device_id/config        Upsert config (ON CONFLICT)
POST   /v1/farms/:farm_id/devices/:device_id/calibrate     Set calibrated value
DELETE /v1/farms/:farm_id/devices/:device_id/commands/:id  Cancel pending command

GET    /v1/farms/:farm_id/commands                         Command history across all devices
POST   /v1/farms/:farm_id/emergency-stop                   Stop all online devices
POST   /v1/farms/:farm_id/devices/batch-command            Send same command to multiple devices
```

---

## 👤 User Session Endpoints (NEW Feb 2026)

**File**: `api/user_handler.go` (3 handlers appended)
```
GET    /v1/users/sessions                  List active sessions (not expired)
DELETE /v1/users/sessions/:session_id      Revoke a session
GET    /v1/users/activity-log              Event log for current user (days, limit, offset)
```

---


**Login** (`POST /v1/auth/login`):
```json
{
  "phone": "012345678",  // OR "email": "user@example.com"
  "password": "Farmer123"
}
```

**JWT Token Validation** (all protected endpoints):
1. Extract token from `Authorization: Bearer <token>` header
2. Validate signature, expiry
3. Set `c.Locals("user_id")` and `c.Locals("role")` for handlers
4. Check farm access: `checkFarmAccess(userID, farmID, requiredRole)`

**Roles**:
- **Owner**: Full farm access (owns the farm)
- **Manager**: Device control, schedules, can invite Viewers
- **Viewer**: Read-only monitoring data

---

## 🗄️ Database Patterns

**Connection** (PostgreSQL only):
```go
// Production & development: PostgreSQL via DATABASE_URL env var
database.DB.Query(`SELECT * FROM farms WHERE id = $1`, farmID)
```

**No SQLite**: `database/sqlite.go` has been permanently deleted. If `DATABASE_URL` is missing or PostgreSQL is unreachable, the server exits immediately with a fatal error.

**Common Queries**:
```go
// Get user's farms
db.Query(`SELECT f.* FROM farms f 
  INNER JOIN farm_users fu ON f.id = fu.farm_id 
  WHERE fu.user_id = $1`, userID)

// Get active schedules for device
db.Query(`SELECT * FROM schedules 
  WHERE device_id = $1 AND is_active = true 
  ORDER BY priority DESC, next_execution ASC`, deviceID)

// Upsert device configuration (safe to call repeatedly)
db.Exec(`INSERT INTO device_configurations (...) VALUES (...)
  ON CONFLICT (device_id, parameter_name) DO UPDATE SET ...`)

// Get device readings (last 24h)
db.Query(`SELECT sensor_type, value, unit, timestamp FROM device_readings
  WHERE device_id = $1 AND timestamp > CURRENT_TIMESTAMP - ($2 * INTERVAL '1 hour')`,
  deviceID, hours)
```

---

## 🌐 WebSocket Real-Time Updates

**Connection**: `ws://localhost:3000/v1/ws?farm_id={id}&coop_id={id}`

**Hub Manager** (`api/websocket_handler.go`):
```go
// Broadcast after state change
BroadcastDeviceUpdate(farmID, coopID, device)
BroadcastScheduleExecution(farmID, coopID, scheduleExecution)
```

**Message Types**:
- `device_update`: Device state changed
- `command_executed`: Device command completed
- `alert_triggered`: Sensor alert (low water, high temp)
- `schedule_executed`: Schedule ran successfully

---

## 🧪 Common Development Tasks

### Add New Endpoint
1. Create handler in `api/` (e.g., `CreateFarmHandler`)
2. Add route in `main.go`: `v1.Post("/farms", api.CreateFarmHandler)`
3. Add auth middleware if needed
4. Update `../docs/implementation/API.md`

### Add New Database Table
1. Add struct to `models/models.go`
2. Add CREATE TABLE in `database/postgres.go` (only PostgreSQL file)
3. Update `../docs/implementation/DATABASE.md`

### Test Locally
```powershell
# Build and start backend (Windows)
cd middleware
go build -o backend.exe
.\backend.exe   # Starts on http://localhost:3000
```

**Run all API tests** (single command, from repo root):
```powershell
.\test_all_endpoints.ps1             # login with email (default)
.\test_all_endpoints.ps1 -UsePhone   # login with phone number
```

The script tests all sections sequentially: Auth → Profile → Farm → Coop → Device → Schedules (incl. `action_sequence`/`action_duration`) → WebSocket → Logout and prints a pass/fail/skip summary.

**Seeded test-data constants used by the script**:
- Email: `farmer@tokkatot.com` / Password: `FarmerPass123`
- Farm ID: `11111111-1111-1111-1111-111111111111`
- Device ID: `33333333-3333-3333-3333-333333333333`
- Schedule with `action_sequence`: `44444444-4444-4444-4444-444444444444`
- Schedule with `action_duration`: `55555555-5555-5555-5555-555555555555`

---

## ⚠️ Critical Rules

1. **Never commit `.env`** - contains secrets
2. **Always validate input** - prevent SQL injection
3. **Check farm access** - `checkFarmAccess()` before device control
4. **Log device commands** - insert into `device_commands` table
5. **Broadcast WebSocket** - after state changes
6. **Use UTC timestamps** - timezone consistency
7. **Consistent errors** - use `utils.BadRequest()`, `utils.NotFound()`

---

## 📚 Documentation Map

**Component AI Context Files** (read for specific tech stack):
- `../middleware/AI_CONTEXT.md` ← YOU ARE HERE (Go API)
- `../frontend/AI_CONTEXT.md` (Vue.js 3 PWA, UI components)
- `../ai-service/AI_CONTEXT.md` (PyTorch disease detection)
- `../embedded/AI_CONTEXT.md` (ESP32 firmware, MQTT)

**Implementation Guides** (read before coding):
- `../docs/implementation/API.md` - Complete API spec
- `../docs/implementation/DATABASE.md` - Full schema
- `../docs/implementation/SECURITY.md` - Auth & authorization
- `../docs/AUTOMATION_USE_CASES.md` - Real farmer scenarios ⭐

**Project Context** (read first):
- `../AI_INSTRUCTIONS.md` - Master guide (business model, farmer-centric design)

**End of middleware/AI_CONTEXT.md**
