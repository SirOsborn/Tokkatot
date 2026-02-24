# 🤖 AI Instructions — Tokkatot

**Last Updated**: February 24, 2026  
**Version**: 2.3  
**Purpose**: Single source of truth for all AI agents working on Tokkatot

> **This is the only AI context file you need to read.** It covers the full project: Go middleware, Vue.js frontend, Python AI service, and ESP32 embedded firmware. The component-level `AI_CONTEXT.md` files in each folder are just pointers back here.

---

## 📋 Project Overview

**Tokkatot** is a cloud-based smart farming IoT system designed for Cambodian poultry farmers.

| | |
|---|---|
| **Target users** | Cambodian chicken farmers — mobile-first, Khmer & English, elderly-friendly |
| **Business model** | On-site install by Tokkatot staff; FREE verification via registration key |
| **Hardware** | Raspberry Pi (coop controller), ESP32 sensors, AI camera above manure belt |
| **Backend** | Go 1.23+ / Fiber v2 — 67 REST endpoints + WebSocket |
| **Frontend** | Vue.js 3 (CDN, no build step) — PWA served by Go |
| **Database** | PostgreSQL only (14 tables) — SQLite was permanently removed Feb 2026 |
| **AI service** | Python 3.12 / FastAPI / PyTorch — EfficientNetB0 + DenseNet121 ensemble |
| **Embedded** | ESP32 / ESP-IDF (C) — MQTT to local Raspberry Pi hub |

**Current features (v2.3)**:
- Device control (water pumps, conveyor belts, fans, feeders) per coop
- AI disease detection from chicken feces photos (99% accuracy)
- Real-time monitoring with WebSocket
- Schedule automation: time-based, duration-based, condition-based, manual
  - `action_duration` — auto-turn-off after N seconds
  - `action_sequence` — multi-step pulse patterns (ON 30s, OFF 10s, repeat)
- Alerts + subscriptions (active, history, quiet hours)
- Analytics: dashboard, device-metrics reports, CSV export, event log
- **Temperature timeline dashboard** (`/monitoring`) — Apple Weather-style graph, H/L peaks, hourly strip, dynamic bg gradient, SVG bezier curve

---

## �️ Repository Layout

```
tokkatot/
├── middleware/           Go 1.23 REST API (Fiber v2) — 67 endpoints + WebSocket
│   ├── main.go           Entry point, routes, static file serving
│   ├── api/              Handlers: auth, farm, coop, device, schedule, alert, analytics, user, websocket
│   ├── database/
│   │   └── postgres.go   PostgreSQL schema (14 tables) — ONLY database file
│   ├── models/models.go  Go structs for all tables
│   ├── utils/utils.go    JWT, bcrypt, error helpers
│   └── .env              DATABASE_URL, JWT_SECRET (gitignored)
│
├── frontend/             Vue.js 3 PWA (no build step, CDN)
│   ├── pages/            HTML pages served by Go middleware
│   │   ├── index.html      Dashboard (device grid, WebSocket live updates)
│   │   ├── monitoring.html Temperature timeline (Apple Weather-style)
│   │   ├── disease-detection.html  AI feces upload + prediction result
│   │   ├── login.html / signup.html
│   │   ├── profile.html / settings.html
│   │   └── 404.html
│   ├── js/               Vue app logic (one file per page)
│   │   ├── index.js / scriptHome.js / login.js / signup.js
│   │   ├── profile.js / scriptSettings.js
│   │   ├── header.js / navbar.js
│   │   └── libs/         Third-party JS libraries
│   ├── css/              Styles (one file per page module)
│   │   ├── styleHome.css / loginSignUp.css / styleProfile.css
│   │   ├── styleSettings.css / styleHeader.css / stylenavbar.css
│   ├── components/
│   │   ├── header.html   Shared farm dropdown + alert bell
│   │   └── navbar.html   Bottom tab navigation (5 tabs)
│   └── assets/           Images, fonts (Mi Sans, Khmer), icons
│
├── ai-service/           Python 3.12 / FastAPI / PyTorch
│   ├── app.py            FastAPI server (3 endpoints: /health, /predict, /predict/detailed)
│   ├── inference.py      ChickenDiseaseDetector class (ensemble voting)
│   ├── models.py         EfficientNetB0Wrapper + DenseNet121Wrapper
│   ├── data_utils.py     Image preprocessing, class names, transforms
│   ├── Dockerfile / docker-compose.yml
│   └── outputs/          *.pth model files (NOT in git — proprietary)
│       └── ensemble_model.pth  47.2 MB ensemble weights
│
├── embedded/             ESP32 / ESP-IDF (C)
│   ├── main/src/         main.c, mqtt_handler.c, device_control.c,
│   │                     sensor_reader.c, ota_updater.c
│   ├── main/include/     device_config.h (GPIO pin definitions)
│   └── components/dht/   DHT22 driver
│
├── docs/
│   ├── implementation/   API.md, DATABASE.md, FRONTEND.md, EMBEDDED.md,
│   │                     AI_SERVICE.md, SECURITY.md
│   ├── AUTOMATION_USE_CASES.md   Real farmer scenarios with JSON examples
│   ├── ARCHITECTURE.md / TECH_STACK.md
│   └── troubleshooting/  API_TESTING.md, DATABASE.md
│
├── AI_INSTRUCTIONS.md    ← THIS FILE — single source of truth
├── .github/copilot-instructions.md   VS Code Copilot pointer to this file
└── test_all_endpoints.ps1   Full API test script (PowerShell)
```

---

## 🏗️ System Architecture

### Data Hierarchy

```
User (Farmer Sokha)
  ├─ Phone: 012345678  (login credential; email optional)
  │
  ├─ Farm: "Kandal Province Farm"
  │   ├─ Registration Key: XXXXX-XXXXX-XXXXX (single-use, auto-verifies account)
  │   │
  │   ├─ Coop 1 (500 layer chickens)
  │   │   ├─ Raspberry Pi  — main controller (is_main_controller=true)
  │   │   ├─ AI Camera     — above manure conveyor belt
  │   │   ├─ Water Tank Sensor  — ultrasonic, outside coop
  │   │   ├─ Water Pump Motor   — fills THIS coop's tank
  │   │   ├─ Conveyor Belt Motor
  │   │   └─ Temperature/Humidity Sensor (DHT22 via ESP32)
  │   │
  │   └─ Coop 2 (300 broiler chickens) ← same device set
  │
  └─ Farm 2: "Kampong Cham Farm"  ← farmers can own multiple farms
```

**Rules:**
- Every device belongs to a **coop**, never directly to a farm
- Farmers **cannot** add/remove devices — Tokkatot staff do this on-site
- `farm_id` is for organisational grouping only, not device control

### Communication Flow

```
Farmer's phone
    ↕ HTTPS / WSS
Go Middleware (cloud, port 3000)
    ↕ HTTP REST
FastAPI AI Service (port 8000)
    ↕ WebSocket
Raspberry Pi (local hub per farm)
    ↕ MQTT (WiFi LAN)
ESP32 Devices (per coop)
    ↕ GPIO / PWM
Physical devices (pumps, relays, fans)
```

### Setup Process

```
1. Tokkatot staff visits farm
2. Installs Raspberry Pi + ESP32 devices in each coop
3. Registers farmer account:  POST /v1/auth/signup  with registration_key
4. Account auto-verified (no SMS/email cost) ✅
5. Farmer logs in immediately with phone + password
```

---

## 🔐 Authentication & Roles

### Registration

```json
POST /v1/auth/signup
{
  "phone": "012345678",
  "password": "Farmer123",
  "name": "Sokha",
  "registration_key": "ZKJFA-ZIVMC-HOUGG-XQSRW-ITDYH"
}
→ { "contact_verified": true }   // auto-verified, no SMS needed
```

### Login (phone OR email)

```json
POST /v1/auth/login
{ "phone": "012345678", "password": "Farmer123" }
→ { "access_token": "...(24h)", "refresh_token": "...(30d)", "farms": [...] }
```

### Roles

| Role | Access |
|------|--------|
| **Farmer** | Full access — devices, schedules, farm settings, invite/remove members (farmer or viewer). Multiple farmers can share equal full access to the same farm. |
| **Viewer** | Worker — read-only monitoring + acknowledge alerts only; no device control, no farm settings, no user management |

> **Tokkatot system staff** are **not** a `farm_users` role. They manage registration keys (`generate_reg_key.ps1`), JWT secrets (`.env`), and system bypasses at infrastructure level — outside the normal farm permission model entirely.

JWT claims: `user_id`, `email` (nullable), `phone` (nullable), `farm_id`, `role`. All protected endpoints call `checkFarmAccess(userID, farmID, requiredRole)`.

---

## 🖥️ Middleware — Go API (Port 3000)

**Stack**: Go 1.23, Fiber v2.52.6, PostgreSQL (`database/postgres.go` — the only DB file)

### Static Page Routes (served from `frontend/`)

| Route | File |
|---|---|
| `GET /` | `pages/index.html` |
| `GET /login` | `pages/login.html` |
| `GET /register` | `pages/signup.html` |
| `GET /profile` | `pages/profile.html` |
| `GET /settings` | `pages/settings.html` |
| `GET /disease-detection` | `pages/disease-detection.html` |
| `GET /monitoring` | `pages/monitoring.html` |

Static dirs `/assets`, `/components`, `/css`, `/js` served from `frontend/`.

### API Endpoints (67 total, all under `/v1/`)

**Auth** — `api/auth_handler.go` (8 endpoints)
```
POST /v1/auth/signup            POST /v1/auth/login
POST /v1/auth/refresh           POST /v1/auth/logout
POST /v1/auth/verify            POST /v1/auth/forgot-password
POST /v1/auth/reset-password    GET  /v1/auth/me
```

**Farms & Members** — `api/farm_handler.go` (9 endpoints)
```
GET|POST /v1/farms
GET|PUT|DELETE /v1/farms/:farm_id
GET|POST /v1/farms/:farm_id/members
PUT|DELETE /v1/farms/:farm_id/members/:user_id
```

**Coops** — `api/coop_handler.go` (6 endpoints)
```
GET|POST /v1/farms/:farm_id/coops
GET|PUT|DELETE /v1/farms/:farm_id/coops/:coop_id
GET /v1/farms/:farm_id/coops/:coop_id/temperature-timeline?days=7
```

**Devices** — `api/device_handler.go` (17 endpoints)
```
GET|POST /v1/farms/:farm_id/devices
GET|PUT|DELETE /v1/farms/:farm_id/devices/:device_id
POST /v1/farms/:farm_id/devices/:device_id/commands
GET /v1/farms/:farm_id/devices/:device_id/history
GET /v1/farms/:farm_id/devices/:device_id/status
GET|PUT /v1/farms/:farm_id/devices/:device_id/config
POST /v1/farms/:farm_id/devices/:device_id/calibrate
DELETE /v1/farms/:farm_id/devices/:device_id/commands/:id
GET /v1/farms/:farm_id/commands
POST /v1/farms/:farm_id/emergency-stop
POST /v1/farms/:farm_id/devices/batch-command
```

**Schedules** — `api/schedule_handler.go` (7 endpoints)
```
GET|POST /v1/farms/:farm_id/schedules
GET|PUT|DELETE /v1/farms/:farm_id/schedules/:schedule_id
GET /v1/farms/:farm_id/schedules/:schedule_id/executions
POST /v1/farms/:farm_id/schedules/:schedule_id/execute-now
```

**Alerts** — `api/alert_handler.go` (8 endpoints)
```
GET /v1/farms/:farm_id/alerts/history    ← registered BEFORE /:alert_id
GET /v1/farms/:farm_id/alerts
GET|PUT(acknowledge) /v1/farms/:farm_id/alerts/:alert_id
GET|POST /v1/users/alert-subscriptions
GET|PUT|DELETE /v1/users/alert-subscriptions/:id
```

**Analytics** — `api/analytics_handler.go` (6 endpoints)
```
GET /v1/farms/:farm_id/dashboard
GET /v1/farms/:farm_id/reports/device-metrics
GET /v1/farms/:farm_id/reports/device-usage
GET /v1/farms/:farm_id/reports/farm-performance
GET /v1/farms/:farm_id/reports/export         (CSV, Content-Disposition: attachment)
GET /v1/farms/:farm_id/events
```

**Users** — `api/user_handler.go` (6 endpoints)
```
GET|PUT /v1/users/profile
GET /v1/users/sessions
DELETE /v1/users/sessions/:session_id
GET /v1/users/activity-log
```

**WebSocket**: `ws://host/v1/ws?farm_id=&coop_id=`  
Message types: `device_update`, `command_executed`, `alert_triggered`, `schedule_executed`

### Database (PostgreSQL, 14 tables)

**Key tables**: `users`, `farms`, `farm_users`, `coops`, `devices`, `device_commands`, `device_readings`, `device_configurations`, `schedules`, `schedule_executions`, `alerts`, `alert_subscriptions`, `user_sessions`, `registration_keys`

**No SQLite** — `database/sqlite.go` was permanently deleted Feb 2026. If `DATABASE_URL` is missing the server exits immediately.

**Common patterns**:
```go
// PostgreSQL only
database.DB.Query(`SELECT * FROM farms WHERE id = $1`, farmID)

// Upsert config (safe to call repeatedly)
db.Exec(`INSERT INTO device_configurations (...) VALUES (...)
  ON CONFLICT (device_id, parameter_name) DO UPDATE SET ...`)

// Get sensor readings (last N hours)
db.Query(`SELECT sensor_type, value, unit, timestamp FROM device_readings
  WHERE device_id = $1 AND timestamp > CURRENT_TIMESTAMP - ($2 * INTERVAL '1 hour')`,
  deviceID, hours)
```

### Temperature Timeline (v2.3)

**Handler**: `TemperatureTimelineHandler` in `api/coop_handler.go`  
**Route**: `GET /v1/farms/:farm_id/coops/:coop_id/temperature-timeline?days=7`

```json
{
  "sensor_found": true,
  "current_temp": 34.2,
  "bg_hint": "hot",
  "today": {
    "date": "2026-02-24",
    "hourly": [{"hour": "14:00", "temp": 34.2}],
    "high": {"temp": 34.5, "time": "14:00"},
    "low":  {"temp": 24.1, "time": "05:00"}
  },
  "history": [
    {"date": "2026-02-24", "label": "Today",     "high": {...}, "low": {...}},
    {"date": "2026-02-23", "label": "Yesterday", "high": {...}, "low": {...}}
  ]
}
```

`bg_hint` values: `scorching` ≥35°C, `hot` ≥32°C, `warm` ≥28°C, `neutral` ≥24°C, `cool` ≥20°C, `cold` <20°C.  
Returns `sensor_found: false` (HTTP 200) when coop has no active temperature sensor. Temperature only — humidity excluded everywhere.

### Schedule Automation

**4 types**:

| Type | Key fields | Use case |
|------|-----------|---------|
| `time_based` | `cron_expression` + `action_duration` OR `action_sequence` | Feeder ON at 6AM for 15 min |
| `duration_based` | `on_duration`, `off_duration` | Conveyor ON 10min, OFF 15min, repeat |
| `condition_based` | `condition_json` | Pump ON when water < 20% |
| `manual` | — | Farmer direct control |

**`action_duration`** (auto-turn-off):
```json
{ "cron_expression": "0 6,12,18 * * *", "action_duration": 900, "action": "set_relay", "action_value": "ON" }
```

**`action_sequence`** (pulse pattern):
```json
{
  "cron_expression": "0 6,12,18 * * *",
  "action_sequence": "[{\"action\":\"ON\",\"duration\":30},{\"action\":\"OFF\",\"duration\":10},{\"action\":\"ON\",\"duration\":30},{\"action\":\"OFF\",\"duration\":10}]"
}
```

### Build & Test (Windows)

```powershell
cd middleware
go build -o backend.exe
.\backend.exe          # starts on http://localhost:3000

# Full API test suite from repo root:
.\test_all_endpoints.ps1             # email login
.\test_all_endpoints.ps1 -UsePhone   # phone login
```

Test data seeded: email `farmer@tokkatot.com` / `FarmerPass123`, farm `11111111-...`, device `33333333-...`

### Critical Rules

1. Never commit `.env` — contains `DATABASE_URL`, `JWT_SECRET`
2. Always call `checkFarmAccess()` before device control
3. Log every device command to `device_commands` table
4. Broadcast WebSocket after state changes
5. Use `utils.BadRequest()`, `utils.NotFound()` for consistent errors

---

## 🌐 Frontend — Vue.js 3 PWA

**Stack**: Vue.js 3 via CDN (no npm, no build step), Mi Sans + Material Symbols, vanilla JS

**Target users**: Elderly Cambodian farmers, 60+, low digital literacy, 1–2 GB RAM phones, 4G

### Accessibility (non-negotiable)

- Buttons: min 48×48 px
- Body text: min 16 px; headings 24 px+; important numbers 48 px+
- Contrast: WCAG AAA (black on white — no grey text)
- Icons: always paired with text labels
- Max 3 main actions per screen
- Khmer/English toggle reachable from every page (top-right header)

### File Map (what actually exists)

```
frontend/pages/
  index.html              Dashboard — device status grid, WebSocket live updates
  monitoring.html         Temperature timeline (Apple Weather-style)
  disease-detection.html  AI feces photo upload + prediction
  login.html / signup.html / profile.html / settings.html / 404.html

frontend/js/
  index.js                Dashboard Vue app
  scriptHome.js           Home page helpers
  login.js / signup.js    Auth pages
  profile.js              Profile Vue app
  scriptSettings.js       Settings page
  header.js / navbar.js   Shared component logic
  libs/                   Third-party JS

frontend/css/
  styleHome.css / loginSignUp.css / styleProfile.css
  styleSettings.css / styleHeader.css / stylenavbar.css

frontend/components/
  header.html / navbar.html
```

> **Note**: `disease-detection.js`, `utils/api.js` do **not exist yet**. `schedules.html` ✅ exists as a self-contained page (CSS + JS inline, same pattern as `monitoring.html`).

### Vue App Pattern

```javascript
const { createApp } = Vue;
createApp({
  data() {
    return { farm: null, devices: [], ws: null, loading: true };
  },
  async mounted() {
    const token = localStorage.getItem('access_token');
    if (!token) { window.location.href = '/login'; return; }
    await this.loadFarmData();
    this.connectWebSocket();
  },
  methods: {
    async loadFarmData() {
      const farmId = localStorage.getItem('selected_farm_id');
      const res = await fetch(`/v1/farms/${farmId}`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('access_token')}` }
      });
      this.farm = await res.json();
    },
    connectWebSocket() {
      const token = localStorage.getItem('access_token');
      this.ws = new WebSocket(`wss://${window.location.host}/v1/ws?token=${token}`);
      this.ws.onmessage = (e) => {
        const msg = JSON.parse(e.data);
        if (msg.type === 'device_update') {
          const d = this.devices.find(d => d.id === msg.device_id);
          if (d) d.state = msg.state;
        }
      };
    }
  }
}).mount('#app');
```

**API URL pattern**: always `/v1/` prefix (e.g. `/v1/farms/${id}/devices`).

### Responsive CSS Breakpoints

```css
/* Mobile first ≤480px */
button { font-size: 48px; min-height: 60px; }
/* Tablet 481–768px */
@media (min-width: 481px) { .device-grid { grid-template-columns: repeat(2, 1fr); } }
/* Desktop >768px */
@media (min-width: 769px) { .device-grid { grid-template-columns: repeat(3, 1fr); } }
```

### Internationalisation

```javascript
const I18N = {
  en: { 'device.on': 'On', 'schedule.create': 'Create Schedule' },
  km: { 'device.on': 'បើក', 'schedule.create': 'បង្កើតកាលវិភាគ' }
};
const lang = localStorage.getItem('language') || 'km';
document.querySelectorAll('[data-i18n]').forEach(el => {
  el.textContent = I18N[lang][el.dataset.i18n];
});
```

### Adding a New Page

1. Create `frontend/pages/new-page.html` with `<script src="https://unpkg.com/vue@3/dist/vue.global.js">` and `<div id="app">`
2. Create `frontend/js/new-page.js` with `createApp({...}).mount('#app')`
3. Create `frontend/css/styleNewPage.css` (mobile-first)
4. Add tab to `frontend/components/navbar.html`
5. Register static route in `middleware/main.go`

---

## 🤖 AI Service — PyTorch Disease Detection (Port 8000)

**Stack**: Python 3.12, PyTorch 2.0+, FastAPI, Uvicorn  
**Model file**: `outputs/ensemble_model.pth` (47.2 MB, **not in git**, must exist locally)

### Why It Exists

Cambodian farmers can't afford vets for every sick bird. AI model diagnoses disease from a feces photo in <3 seconds with 99% accuracy, enabling early treatment before an outbreak spreads.

### Ensemble Architecture

Two models vote:
- **EfficientNetB0** (98.05% recall) — lightweight, fast, 224×224 RGB input
- **DenseNet121** (96.69% recall) — dense feature reuse, catches subtle patterns

```python
avg_probs = (efficientnet_probs + densenet_probs) / 2
predicted_class = argmax(avg_probs)
if max(avg_probs) < 0.50:
    return "uncertain", confidence, "Please retake photo in better lighting"
```

4 classes: `Coccidiosis`, `Healthy`, `Newcastle Disease`, `Salmonella`

### Endpoints

```
GET  /health         Service status + model_loaded, device, timestamp
POST /predict        { disease, confidence, recommendation, treatment_options }
POST /predict/detailed   Per-model confidence scores (research/debug)
```

Image validation: PNG/JPEG only, max 5 MB. Check magic bytes, not just extension.

### Key Patterns (inference.py)

```python
class ChickenDiseaseDetector:
    def __init__(self, model_path='outputs/ensemble_model.pth'):
        checkpoint = torch.load(model_path, map_location='cpu')
        self.model = create_ensemble()
        self.model.load_state_dict(checkpoint['model_state_dict'])
        self.model.eval()        # CRITICAL: disable dropout/batch norm
    
    def predict(self, pil_image):
        with torch.no_grad():    # CRITICAL: no gradient computation
            tensor = self.preprocess(pil_image)
            probs = torch.softmax(self.model(tensor), dim=1)
            conf, cls = torch.max(probs, dim=1)
            return CLASS_NAMES[cls.item()], conf.item()
    
    def preprocess(self, pil_image):
        return transforms.Compose([
            transforms.Resize((224, 224)),
            transforms.ToTensor(),
            transforms.Normalize([0.485,0.456,0.406],[0.229,0.224,0.225])
        ])(pil_image).unsqueeze(0)
```

### Go ↔ AI Integration

Go middleware calls `http://localhost:8000/predict` (or Docker network).  
Timeout: 3 s (CPU), target <500 ms (GPU). See `middleware/api/disease-detection.go`.

### Docker

```bash
cd ai-service
docker build -t tokkatot-ai:latest .   # outputs/ must exist with .pth files
docker-compose up -d tokkatot-ai       # port 8000, 2 CPU / 4 GB RAM
```

### Local Dev

```bash
cd ai-service
python3 -m venv env && source env/bin/activate
pip install -r requirements.txt
python3 app.py     # http://localhost:8000 (needs outputs/ensemble_model.pth)
```

### Critical Rules

- Never expose `outputs/ensemble_model.pth` path in error messages
- Never commit `*.pth` files to git
- Always call `.eval()` before inference
- `outputs/` must exist locally before `docker build` or `python app.py`

---

## 📡 Embedded — ESP32 Firmware

**Stack**: ESP-IDF (C), FreeRTOS, MQTT, GPIO/PWM  
**Comms**: MQTT → Raspberry Pi local hub → WebSocket → Go Middleware

### GPIO Pinout (device_config.h)

```c
#define RELAY_PIN_1    GPIO_NUM_12   // Water pump, light, etc.
#define RELAY_PIN_2    GPIO_NUM_13
#define PWM_FAN_PIN    GPIO_NUM_14   // 0–100% speed
#define PWM_LIGHT_PIN  GPIO_NUM_15
#define DHT22_PIN      GPIO_NUM_26   // Temp + humidity
#define BUTTON_PIN     GPIO_NUM_27
#define LED_GREEN_PIN  GPIO_NUM_4    // Online indicator
#define LED_RED_PIN    GPIO_NUM_2    // Error indicator
```

### MQTT Topics

**Subscribe (receive commands)**:
```
farm/{farm_id}/devices/{device_id}/command
    payload: { "command": "on"|"off"|"pwm", "pwm_value": 0-100, "duration_seconds": N }

farm/{farm_id}/devices/{device_id}/sequence
    payload: { "action_sequence": [{"action":"ON","duration":30}, ...] }

farm/{farm_id}/devices/{device_id}/config
farm/{farm_id}/devices/{device_id}/ota
```

**Publish (send status)**:
```
farm/{farm_id}/devices/{device_id}/status   (every 30 s — heartbeat)
    payload: { "is_online": true, "command_state": "on", "uptime_seconds": N, "rssi": -65 }

farm/{farm_id}/devices/{device_id}/sensor   (every 30 s)
    payload: { "temperature": 28.5, "humidity": 65.2 }

farm/{farm_id}/devices/{device_id}/error
```

### Multi-Step Sequence Execution (device_control.c)

```c
void execute_action_sequence(cJSON *steps) {
    for (int i = 0; i < cJSON_GetArraySize(steps); i++) {
        cJSON *step = cJSON_GetArrayItem(steps, i);
        const char *action = cJSON_GetStringValue(cJSON_GetObjectItem(step, "action"));
        int duration = (int)cJSON_GetNumberValue(cJSON_GetObjectItem(step, "duration"));
        if (strcmp(action, "ON") == 0)  turn_relay_on(RELAY_PIN_1);
        else                             turn_relay_off(RELAY_PIN_1);
        if (duration > 0) vTaskDelay((duration * 1000) / portTICK_PERIOD_MS);
        else break;   // duration=0 means "until next schedule"
    }
}
```

### Build & Flash (Windows)

```powershell
cd embedded
idf.py menuconfig         # set WiFi SSID/password, MQTT broker IP, device ID
idf.py build
idf.py -p COM3 flash      # check Device Manager for correct port
idf.py -p COM3 monitor    # Ctrl+] to exit
```

### Critical Rules

- WiFi credentials stored in NVS (encrypted), never hardcoded
- MQTT requires username/password auth (no anonymous)
- Heartbeat every 30 s so backend marks device online/offline
- All commands must be idempotent (safe to execute twice)
- Enable hardware watchdog timer

---

## 🔐 Security Rules (All Components)

| Rule | Where |
|------|-------|
| No secrets in code | All services use `.env` (gitignored) |
| Parameterised queries only | Go: `$1, $2`; Python: SQLAlchemy params |
| bcrypt passwords | Go: `golang.org/x/crypto/bcrypt` |
| JWT expiry | Access 24 h, refresh 30 d |
| Image validation | Size ≤5 MB, magic bytes check (PNG/JPEG) |
| No sensitive data in errors | Generic "Internal Server Error" to client |
| HTTPS in production | TLS_CERT + TLS_KEY in middleware/.env |
| MQTT auth | Username/password, not open broker |

---

## 📝 Documentation Update Protocol

**Update documentation when you complete a significant feature.** Knowledge lost between AI sessions = bugs and wasted work.

### What to update

| Changed | Update |
|---------|--------|
| Database schema | `docs/implementation/DATABASE.md` |
| API endpoints | `docs/implementation/API.md` |
| Frontend UI | `docs/implementation/FRONTEND.md` |
| Firmware | `docs/implementation/EMBEDDED.md` |
| AI model/service | `docs/implementation/AI_SERVICE.md` |
| Farmer use case | `docs/AUTOMATION_USE_CASES.md` |
| Major system concept | This file (`AI_INSTRUCTIONS.md`) |

### Checklist (copy-paste)

```
[ ] Schema changed?         → docs/implementation/DATABASE.md
[ ] Endpoints added/changed? → docs/implementation/API.md
[ ] New UI component?        → docs/implementation/FRONTEND.md
[ ] Firmware changed?        → docs/implementation/EMBEDDED.md
[ ] Farmer problem solved?   → docs/AUTOMATION_USE_CASES.md
[ ] All code examples compile/run?
[ ] All JSON examples valid?
```

### Timing rule

- Too soon: don't update after every function (noise)
- Too late: knowledge lost for next session
- **Just right**: after 30–60 min of significant work, or when a feature is fully working

---

## 🧪 Development Checklist

### Add a new API endpoint

1. Write handler in `middleware/api/`
2. Register route in `middleware/main.go`
3. Add `checkFarmAccess()` if it touches farm data
4. Log any device command to `device_commands`
5. Update `docs/implementation/API.md`

### Add a new database table

1. Add Go struct to `middleware/models/models.go`
2. Add `CREATE TABLE` in `middleware/database/postgres.go`
3. Update `docs/implementation/DATABASE.md`

### Add a new frontend page

1. `frontend/pages/new.html` — Vue CDN script + `<div id="app">`
2. `frontend/js/new.js` — createApp + mounted() data load
3. `frontend/css/styleNew.css` — mobile-first, 48 px buttons
4. `frontend/components/navbar.html` — add tab
5. `middleware/main.go` — register static route

---

## ✅ Code Review Checklist

- [ ] No `.env` values committed
- [ ] No `*.pth` model files committed
- [ ] All protected endpoints check JWT + farm access
- [ ] Input validated (size, type, SQL injection)
- [ ] Errors don't expose internal paths or DB details
- [ ] Documentation updated
- [ ] No hardcoded URLs, passwords, or secrets

---

## 🎯 Final Reminder

You are building for **elderly Cambodian farmers**. Every decision: *can a 65-year-old farmer use this on a budget phone with 4G?*

- Simplicity beats cleverness
- Reliability beats features
- Accessibility beats aesthetics

When in doubt — choose the simpler solution.
