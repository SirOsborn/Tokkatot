# 🤖 AI Agent Instructions for Tokkatot 2.0

**Last Updated**: February 23, 2026  
**Version**: 2.2 (Updated with accurate DB schema: devices coop_id/is_main_controller, device_commands field names, event_logs structure, schedule_executions status enum)  
**Purpose**: Guide for AI agents assisting with Tokkatot development

---

## 📋 Project Overview

**Tokkatot 2.0** is a cloud-based smart farming IoT system designed for Cambodian poultry farmers.

**Key Context:**
- **Target Users**: Cambodian chicken farmers (Khmer/English, mobile-first UI)
- **Business Model**: On-site installation by Tokkatot staff, FREE verification using registration keys
- **Main Hardware**: Raspberry Pi (main controller per coop), ESP32 sensors, AI cameras
- **Cloud Stack**: Go 1.23+ (API), Vue.js 3 (PWA), PostgreSQL (Data), Python FastAPI (AI disease detection)
- **Primary Features**: 
  - Coop-specific device control (water pumps, conveyor belts, sensors)
  - AI-powered disease detection from chicken feces images (camera above manure conveyor belt)
  - Real-time monitoring (water levels, temperature, humidity)
  - Automated alerts (disease outbreak, low water, equipment failure)
- **Architecture**: Farm → Coops → Devices (hierarchical, farmer can own multiple farms)

---

## 🐔 **CRITICAL: System Architecture Understanding**

### **Business Model & Setup Process**
```
1. Farmer purchases Tokkatot service
2. Tokkatot team visits farm location (on-site installation)
3. Team installs hardware (Raspberry Pi, sensors, cameras, pumps)
4. Team registers farmer's account using:
   - Phone number (required)
   - Password (set on-site)
   - Registration key (pre-generated, tied to farm location)
5. Account is AUTO-VERIFIED (no email/SMS costs) ✅
6. Farmer logs in immediately
```

**Key Point**: NO EMAIL/SMS VERIFICATION NEEDED! Registration key proves legitimacy.

---

## 🚨 **MANDATORY: Self-Documentation After Building Features**

**CRITICAL REQUIREMENT**: Whenever you complete building a significant feature, you MUST update documentation AND AI context files immediately.

### Why This Matters

Future AI sessions need to know:
- ✅ What features exist (e.g., `action_sequence` for multi-step automation)
- ✅ How things were implemented (Go patterns, Vue.js components, ESP32 code)
- ✅ Why decisions were made (farmer needs, technical constraints)
- ✅ Where to find examples (API requests, database queries, UI components)

**Without documentation updates, knowledge is lost between sessions!**

### When You MUST Update Docs

**Update docs when you complete**:
- ✅ New feature (e.g., pulse feeding automation, disease detection API)
- ✅ Database schema change (added table, added fields like `action_sequence`)
- ✅ New API endpoint or modified endpoint behavior
- ✅ New automation pattern (schedule types, sensor triggers)
- ✅ New UI component (schedule builder, device control panel)
- ✅ Architecture decision (PostgreSQL-only migration, JWT auth flow)
- ✅ Integration work (Go → FastAPI → PyTorch, MQTT protocol)

**Don't update for**:
- ❌ Minor bug fixes (typo corrections, null checks)
- ❌ Code refactoring without functional changes
- ❌ Variable/function renames
- ❌ Comment additions
- ❌ Formatting/linting changes

### What You MUST Update

**Every significant feature requires updating**:

1. **Implementation specs** (`docs/implementation/*.md`):
   - `API.md` - If you added/modified endpoints (show full request/response examples)
   - `DATABASE.md` - If you changed schema (show CREATE TABLE, explain new fields)
   - `FRONTEND.md` - If you added UI (show Vue.js code, CSS, user flow)
   - `EMBEDDED.md` - If you changed firmware (show C code, MQTT topics)
   - `AI_SERVICE.md` - If you changed AI model (show PyTorch code, FastAPI endpoints)

2. **AI context files** (teach future AI sessions):
   - `middleware/AI_CONTEXT.md` - If you added Go patterns, database queries, endpoints
   - `frontend/AI_CONTEXT.md` - If you added Vue.js components, API calls, UI patterns
   - `ai-service/AI_CONTEXT.md` - If you changed model architecture, preprocessing
   - `embedded/AI_CONTEXT.md` - If you changed MQTT protocol, GPIO control, sensors
   - `AI_INSTRUCTIONS.md` (this file) - If you added major system concept (rarely)

3. **Use case documentation**:
   - `docs/AUTOMATION_USE_CASES.md` - If you solved a real farmer problem (show scenario, JSON example, benefit)

### How to Judge "It's Time to Update"

**Ask yourself**:
1. ✅ "Would a future AI be confused about what I just built?" → Update now
2. ✅ "Did I just solve a real farmer problem?" → Document the solution
3. ✅ "Did I create a reusable pattern?" → Add to AI_CONTEXT.md
4. ✅ "Have I made 3-5 related changes?" → Time to consolidate into docs
5. ✅ "Am I switching to a different component?" → Document current work first

**Timing Rule (Goldilocks Zone)**:
- ⏱️ **Too fast**: Don't update after every single function (causes noise)
- ⏱️ **Too slow**: Don't wait weeks for "perfect time" (knowledge gets lost)
- ✅ **Just right**: Update after 30-60 minutes of significant work OR when feature is complete OR before switching components

### Update Checklist (Copy This)

```markdown
After completing significant work, verify:

[ ] Did I change database schema? → Update docs/implementation/DATABASE.md
[ ] Did I add/modify API endpoints? → Update docs/implementation/API.md  
[ ] Did I add UI components? → Update docs/implementation/FRONTEND.md
[ ] Did I change firmware behavior? → Update docs/implementation/EMBEDDED.md
[ ] Did I solve a farmer problem? → Update docs/AUTOMATION_USE_CASES.md
[ ] Did I add reusable patterns? → Update component AI_CONTEXT.md files
[ ] Did I test all examples? → Verify JSON/SQL/code compiles and runs
[ ] Did I add cross-references? → Link related docs together
```

### Real Example: action_sequence Feature

**What was built**: Multi-step automation for pulse feeding (ON 30s, pause 10s, repeat)

**Documentation updated** (same session):
1. ✅ `docs/implementation/DATABASE.md` - Added `action_sequence JSONB` field spec
2. ✅ `docs/implementation/API.md` - Updated 4 schedule endpoints with field examples
3. ✅ `docs/implementation/FRONTEND.md` - Added Action Sequence Builder UI (300+ lines)
4. ✅ `docs/implementation/EMBEDDED.md` - Added ESP32 execution code (200+ lines)
5. ✅ `docs/AUTOMATION_USE_CASES.md` - Created 500+ line guide with farmer scenarios
6. ✅ `middleware/AI_CONTEXT.md` - Added schedule automation section
7. ✅ `AI_INSTRUCTIONS.md` - Added automation & schedules overview

**Result**: Future AI sessions know this feature exists, how to use it, and can build on it.

**If not documented**: Future AI would reinvent the feature or never discover it exists.

### Enforcement

**Before ending your session**:
1. Review what you built today
2. Identify which docs need updates
3. Update ALL relevant docs (don't skip any)
4. Verify examples work (compile code, test JSON, run queries)
5. Add version history entry in `docs/README.md` if major feature

**This is not optional - it's how we maintain institutional knowledge!**

---

### **Data Hierarchy: Farm → Coop → Device**

```
User (Farmer Sokha)
  ├─ Phone: 012345678 (login credential)
  ├─ Password: (set during on-site registration)
  │
  ├─ Farm 1: "Kandal Province Farm"
  │   ├─ Location: Kandal Province, Cambodia
  │   ├─ Registration Key: XXXXX-XXXXX-XXXXX (used once)
  │   │
  │   ├─ Coop 1 (500 layer chickens)
  │   │   ├─ Raspberry Pi (main controller: RaspberryPi-001)
  │   │   ├─ AI Camera (above manure conveyor belt)
  │   │   ├─ Water Tank Sensor (ultrasonic, outside coop)
  │   │   ├─ Water Pump Motor (fills THIS coop's tank)
  │   │   ├─ Conveyor Belt Motor (automated manure removal)
  │   │   └─ Water Pipes (through coop, feeds chickens)
  │   │
  │   └─ Coop 2 (300 broiler chickens)
  │       ├─ Raspberry Pi (main controller: RaspberryPi-002)
  │       ├─ AI Camera
  │       ├─ Water Tank Sensor
  │       ├─ Water Pump Motor
  │       ├─ Conveyor Belt Motor
  │       └─ Water Pipes
  │
  └─ Farm 2: "Kampong Cham Farm"
      └─ Coop 1 (400 mixed chickens)
          ├─ Raspberry Pi (main controller: RaspberryPi-003)
          ├─ AI Camera
          ├─ Water Tank Sensor
          ├─ Water Pump Motor
          ├─ Conveyor Belt Motor
          └─ Water Pipes
```

**IMPORTANT NOTES:**
1. **Each coop is SELF-CONTAINED** - Has its own water pump, camera, sensors
2. **NO farm-level devices** - All devices belong to specific coops
3. **farm_id purpose**: Organizational grouping (multi-province support), NOT device control
4. **Multi-farm support**: Farmers can own farms in different provinces

---

### **Device Details (Per Coop)**

| Device | Location | Purpose | Critical? |
|--------|----------|---------|-----------|
| **AI Camera** | Above manure conveyor belt | Monitors chicken feces texture/pattern for disease detection | ✅ YES |
| **Water Tank Sensor** | Outside coop (on water tank) | Ultrasonic sensor measures water level | ✅ YES |
| **Water Pump Motor** | Connected to farm's main water source | Automatically fills coop's water tank when sensor detects low level | ✅ YES |
| **Conveyor Belt Motor** | Inside coop (under chicken cages) | Automated manure removal system | ⚠️ Important |
| **Water Pipes** | Throughout coop | Delivers water to chickens 24/7 | ✅ YES |
| **Raspberry Pi** | Main controller for THIS coop | Coordinates all coop devices (is_main_controller=true) | ✅ YES |

**System Logic Example:**
```
Coop 1 water tank sensor: "Water level 15%" (LOW!)
  → System checks: coop_id → finds Coop 1's water pump
  → Activates Coop 1's pump motor
  → Pump fills Coop 1's tank from farm's main water source
  → Sensor reads 95% → Pump stops
  → Water flows through Coop 1's pipes to chickens 24/7
```

---

## 🔐 **Authentication & Verification System**

### **Registration (On-Site by Staff)**

**Required Fields:**
- `phone`: Cambodian phone number (e.g., "012345678")
- `password`: Set by staff with farmer present
- `name`: Farmer's name
- `registration_key`: Pre-generated key (e.g., "ZKJFA-ZIVMC-HOUGG-XQSRW-ITDYH")

**Optional Fields:**
- `email`: Not required (farmers may not have email)
- `language`: "km" (Khmer) or "en" (default: "km")

**Registration Flow:**
```json
POST /v1/auth/signup
{
  "phone": "012345678",
  "password": "Farmer123",
  "name": "Sokha",
  "registration_key": "ZKJFA-ZIVMC-HOUGG-XQSRW-ITDYH"
}

Response:
{
  "success": true,
  "data": {
    "user_id": "uuid-here",
    "contact_verified": true  // ✅ Auto-verified by registration key!
  },
  "message": "Account created and verified"
}
```

**Why Registration Keys?**
- ✅ **FREE** - No SMS/email costs
- ✅ **Secure** - Keys expire in 90 days, single-use only
- ✅ **Farmer-friendly** - Staff handles everything on-site
- ✅ **Offline-capable** - Works during installation without internet

### **Login (Farmer Self-Service)**

```json
POST /v1/auth/login
{
  "phone": "012345678",
  "password": "Farmer123"
}

Response:
{
  "success": true,
  "data": {
    "access_token": "jwt-token-24h",
    "refresh_token": "jwt-token-30d",
    "user": {
      "id": "uuid",
      "phone": "012345678",
      "name": "Sokha"
    },
    "farms": [
      {"id": "farm-uuid-1", "name": "Kandal Farm", "coop_count": 2},
      {"id": "farm-uuid-2", "name": "Kampong Cham Farm", "coop_count": 1}
    ]
  }
}
```

**Login → Select Farm → Select Coop → Control Dashboard**

---

## 🎯 Your Role as AI Agent

When working on Tokkatot, you are a **professional software engineer and technical architect**. You should:

1. **Understand farmer-centric design** - Every feature must be simple, accessible, no complex UIs
2. **Maintain consistency** - API specs match database schema, documentation is always in sync
3. **Follow the tech stack** - Go/Vue.js/PostgreSQL/FastAPI (don't suggest alternatives randomly)
4. **Prioritize security** - No secrets in code, use .env files, JWT tokens, proper error handling
5. **Write production-ready code** - Error handling, logging, input validation, no hardcoded values
6. **Update documentation** - If code changes, update relevant spec files
7. **Respect file ownership** - AI/model files are proprietary to Tokkatot (don't push to git)

---

## 📁 Directory Structure & Your Responsibilities

### `/docs/` - Specifications & Architecture
**Your Role**: Keep documentation current with code changes

**File Hierarchy** (Read in this order):
1. `00_SPECIFICATIONS_INDEX.md` - Navigation hub
2. `01_SPECIFICATIONS_ARCHITECTURE.md` - System design, data flows
3. `02_SPECIFICATIONS_REQUIREMENTS.md` - Functional requirements

**Implementation Guides** (IG_*):
- `IG_SPECIFICATIONS_API.md` - 66 REST endpoints, authentication, error handling
- `IG_SPECIFICATIONS_DATABASE.md` - PostgreSQL schema (14 tables, farmer-centric; SQLite removed Feb 2026)
- `IG_SPECIFICATIONS_SECURITY.md` - JWT auth, Email/Phone login, no MFA for farmers
- `IG_SPECIFICATIONS_FRONTEND.md` - Vue.js UI, 48px+ fonts, WCAG AAA accessibility
- `IG_SPECIFICATIONS_EMBEDDED.md` - ESP32 firmware, MQTT protocol
- `IG_SPECIFICATIONS_AI_SERVICE.md` - PyTorch ensemble (99% accuracy)

**Use Case Guides**:
- `AUTOMATION_USE_CASES.md` - **🚜 Real-world farmer automation scenarios** (conveyors, feeders, pumps, climate control)
  - **READ THIS** for understanding schedule types in real farming context
  - Detailed examples: Pulse feeding, cycling conveyors, sensor-driven pumps
  - Multi-step sequences (`action_sequence` field usage)

**Operational Guides** (OG_*):
- `OG_SPECIFICATIONS_DEPLOYMENT.md` - Docker, infrast architecture
- `OG_SPECIFICATIONS_TECHNOLOGY_STACK.md` - Tech selections & versions
- `OG_PROJECT_TIMELINE.md` - 27-35 week schedule
- `OG_TEAM_STRUCTURE.md` - Team roles & responsibilities

**When to Update**:
- ✅ After adding/changing API endpoints → Update IG_SPECIFICATIONS_API.md
- ✅ After modifying database schema → Update IG_SPECIFICATIONS_DATABASE.md
- ✅ After adding features → Update 02_SPECIFICATIONS_REQUIREMENTS.md
- ✅ After changing architecture → Update 01_SPECIFICATIONS_ARCHITECTURE.md & data flow diagrams

---

## 🤖 Automation & Schedules (Critical Farmer Feature)

**Real-World Equipment**:
1. **Conveyor Belt** - Manure removal (motor rotates scraper chain)
2. **Feeder Motor** - Spiral auger feed dispenser (stainless steel spiral in tube pushes feed pellets)
3. **Water Pump** - Tank refill (gravity-fed to coop pipes)
4. **Climate Control** - Fans (cooling) + Heaters (warming)

**4 Schedule Types**:

| Type | Use Case | Key Fields | Example |
|------|----------|-----------|---------|
| **Manual** | Always ON/OFF | None | "I want conveyor running 24/7" |
| **time_based** | Trigger at specific times | `cron_expression`, `action_duration` | "Turn ON feeder at 6AM, 12PM, 6PM for 15min each" |
| **time_based + sequence** | Multi-step pattern at specific times | `cron_expression`, `action_sequence` | "At 6AM: motor ON 30sec, pause 10sec, ON 30sec, pause 10sec" |
| **duration_based** | Continuous ON/OFF cycling | `on_duration`, `off_duration` | "Conveyor ON 10min, OFF 15min, repeat forever" |
| **condition_based** | Sensor-driven automation | `condition_json` | "Pump ON when water < 20%, OFF when > 90%" |

**New Feature - Multi-Step Sequences** (`action_sequence` field):
```json
{
  "cron_expression": "0 6,12,18 * * *",  // Trigger at 6AM, 12PM, 6PM
  "action_sequence": "[
    {\"action\":\"ON\",\"duration\":30},   // Step 1: ON for 30 seconds
    {\"action\":\"OFF\",\"duration\":10},  // Step 2: Pause 10 seconds
    {\"action\":\"ON\",\"duration\":30},   // Step 3: ON for 30 seconds
    {\"action\":\"OFF\",\"duration\":10}   // Step 4: Pause 10 seconds
  ]"
}
// Total: 80 seconds per feeding, then device stays OFF until next trigger
```

**Why This Matters**:
- **Pulse Feeding**: Chickens need time between feed bursts to approach bowls (prevents aggressive eaters from dominating)
- **Electricity Savings**: `duration_based` cycling reduces conveyor runtime by 60-75%
- **Automation**: Farmers can leave farm - schedules handle feeding/cleaning/watering

**👉 See `docs/AUTOMATION_USE_CASES.md` for detailed farmer scenarios with exact JSON examples**

---

### `/middleware/` - Go API Gateway (Port 6060)
**Your Role**: Authentication, authorization, device management, API routing

**Stack**: Go 1.19+, JWT tokens, PostgreSQL driver

**Key Responsibilities**:
- User authentication (Email OR Phone, no MFA for farmers)
- Role-based access control (Owner, Manager, Viewer only)
- Device command routing to ESP32 via MQTT
- Real-time WebSocket updates
- Rate limiting & request validation
- Event logging

**Files**:
- `main.go` - Entry point, HTTP server setup
- `api/` - Endpoint handlers (auth, devices, schedules, etc)
- `database/` - PostgreSQL queries
- `utils/` - Helper functions

**Important Rules**:
- 🔒 Never hardcode database passwords → Use `middleware/.env` (gitignore'd)
- ✅ Always validate user input (size limits, type checks, SQL injection prevention)
- ✅ Always check JWT token & user permissions before granting access
- ✅ All database changes should be reflected in schema migrations
- ✅ Device additions/removals by Tokkatot team only (farmers can't add devices themselves)

---

### `/frontend/` - Vue.js 3 Web Application
**Your Role**: User interface, responsive design, offline support

**Stack**: Vue.js 3, HTML5, CSS3, JavaScript (no build system - vanilla JS)

**Key Responsibilities**:
- Responsive mobile-first design (phones, tablets, desktops)
- Accessibility (48px+ fonts, WCAG AAA, Khmer/English toggle)
- Real-time updates via WebSocket
- Offline capability via Service Workers
- Device control, scheduling, monitoring interfaces

**Files**:
- `pages/` - HTML pages (index, disease-detection, profile, settings, etc)
- `components/` - Reusable header, navbar
- `js/` - JavaScript logic for each page
- `css/` - Styling (one CSS file per page module)
- `assets/` - Images, icons, fonts

**Important Rules**:
- 🌐 No npm packages (vanilla JS only - client-side only)
- ♿ Always test with large fonts (48px minimum for buttons/text)
- 📱 Mobile-first design (test on 1-2GB RAM phones)
- 🔄 Use WebSocket for real-time updates, HTTP for initial data
- 🇰🇭 Support both Khmer & English language toggle

---

### `/ai-service/` - PyTorch Disease Detection (Port 8000)
**Your Role**: AI model service, REST API, disease prediction

**Stack**: Python 3.12, PyTorch 2.0, FastAPI, Uvicorn

**Key Responsibilities**:
- FastAPI endpoints for disease prediction
- PyTorch ensemble model (EfficientNetB0 + DenseNet121)
- Image preprocessing & validation
- 99% accuracy via voting mechanism
- Health checks & error handling

**Files**:
- `app.py` - FastAPI server, endpoint definitions
- `inference.py` - Ensemble model loading & inference logic
- `models.py` - PyTorch model architectures
- `data_utils.py` - Image preprocessing, transforms
- `requirements.txt` - Python dependencies (FastAPI, torch, etc)
- `Dockerfile` - Docker build (Python 3.12-slim, model copying)
- `docker-compose.yml` - Docker Compose config

**Important Rules**:
- 🚫 NEVER commit model files (`*.pth`, `*.h5`) → Use `ai-service/.gitignore`
- ✅ Model inference should complete in 1-3 seconds (CPU) or <500ms (GPU)
- ✅ Always validate image uploads (size max 5MB, format PNG/JPEG only)
- ✅ Return ensemble confidence score + per-model scores in API response
- ✅ Implement safety checks (if confidence < 50%, return "uncertain")
- 🆘 Error responses must NOT expose model paths or sensitive info

**Endpoints** (defined in IG_SPECIFICATIONS_AI_SERVICE.md):
- `GET /health` - Service health & model loading status
- `POST /predict` - Disease prediction (simple response)
- `POST /predict/detailed` - Detailed per-model confidence scores

---

### `/embedded/` - ESP32 Firmware (MQTT Client)
**Your Role**: Device firmware, hardware control, sensor reading

**Stack**: ESP-IDF (C/C++), MQTT, GPIO control

**Key Responsibilities**:
- MQTT communication with local hub (Raspberry Pi)
- GPIO control for actuators (relays, PWM for water pumps, lights, fans, etc)
- Sensor reading (DHT22 for temperature/humidity)
- Firmware OTA updates
- Status LED indicators

**Files**:
- `main/` - Main firmware code
- `components/` - Reusable drivers (DHT sensor, relay control, etc)
- `CMakeLists.txt` - Build configuration
- `sdkconfig` - ESP32 configuration

**Important Rules**:
- 📡 MQTT topics must follow: `farm/{farm_id}/devices/{device_id}/{command_type}`
- 🔒 No hardcoded WiFi passwords → Use provisioning or secure storage
- ✅ Implement command queue for offline operations
- ✅ Send heartbeat every 30 seconds to indicate device is online
- ✅ All commands must be idempotent (safe to execute multiple times)

---

## 🔐 Security Guidelines

**For All AI Agents**:

1. **Never commit secrets**:
   - ✅ Use `.env` files for passwords, API keys, database URLs
   - ✅ All `.env` files are automatically gitignore'd
   - ❌ Never hardcode connection strings or tokens in code

2. **Input validation**:
   - ✅ Validate file types (images are PNG/JPEG only)
   - ✅ Validate file sizes (max 5MB for images)
   - ✅ Validate JSON payloads with schema validation (Pydantic for Python, Go structs)
   - ✅ Rate limit API endpoints to prevent abuse

3. **Database security**:
   - ✅ Use parameterized queries (prepared statements)
   - ❌ Never build SQL strings with string concatenation
   - ✅ Hash passwords with bcrypt (Go: `golang.org/x/crypto/bcrypt`)
   - ✅ Use JWT tokens with expiry (24 hours for access, 30 days for refresh)

4. **Error handling**:
   - ✅ Never expose internal server paths or database details in error messages
   - ✅ Log errors server-side with timestamps and request IDs
   - ✅ Return generic "Internal Server Error" to clients
   - ✅ Use proper HTTP status codes (401 for auth, 403 for permission, 404 for not found, etc)

5. **Network communication**:
   - ✅ Always use HTTPS/TLS in production
   - ✅ WebSocket upgrades must validate JWT token first
   - ✅ MQTT should use username/password authentication (no plaintext)
   - ✅ Rate limit WebSocket connections per user

---

## 📊 Database Philosophy

**Farmer-Centric Schema**:
- **Users**: Email OR phone (not both required), simple 3-role system (Owner/Manager/Viewer)
- **Farms**: One farmer can own multiple farms
- **Devices**: Added by Tokkatot team only (farmers can't add devices)
- **Commands**: Log all device commands for audit trail
- **Sensors**: Store time-series data in InfluxDB, aggregated summaries in PostgreSQL
- **Predictions**: Store AI disease predictions with timestamp, image hash, confidence scores

**Golden Rule**: **Device state is source of truth**. When syncing cloud ↔ local, use device state as authoritative.

---

## 🚀 Deployment & DevOps

**Docker & Containers**:
- Each service has separate Dockerfile (microservices architecture)
- Use Docker Compose for local development
- Model files included in `ai-service` Docker image via `COPY outputs/`
- Health checks on all containers (FastAPI: `/health`, Go: `/api/health`)

**Resource Limits**:
- Go API: 1 CPU, 512MB RAM
- AI Service: 2 CPU, 4GB RAM
- Frontend: Served by Nginx (static) or CDN
- Database: 2 CPU, 4GB RAM (PostgreSQL)

**Git & CI/CD**:
- Model files in `.gitignore` (don't push to GitHub)
- Secrets in `.env` files (don't push to GitHub)
- Build artifacts in `.gitignore`: `middleware/backend.exe`, `middleware/*.exe`, `middleware/*.exe~`
- Temporary test files in `.gitignore`: `test_token.txt`
- Pull requests must include documentation updates
- GitHub Actions can run tests/linting on PR

**Testing** (middleware API — PowerShell, from repo root):
```powershell
.\test_all_endpoints.ps1             # all endpoints, email login
.\test_all_endpoints.ps1 -UsePhone   # all endpoints, phone login
```
Covers: Auth, Profile, Farm, Coop, Device, Schedules (action_sequence + action_duration), WebSocket, Logout.

---

## 📝 Code Review Checklist (GitHub PRs)

When reviewing AI-generated PRs, check:

- ✅ All `.env` files in .gitignore (no secrets leaked)
- ✅ Model files not committed (`*.pth`, `*.h5`)
- ✅ Database schema matches API spec (IG_SPECIFICATIONS_DATABASE.md)
- ✅ API endpoints match spec (IG_SPECIFICATIONS_API.md)
- ✅ Error messages don't expose sensitive info
- ✅ JWT token validation on protected endpoints
- ✅ Input validation on all user inputs
- ✅ Documentation updated to match code changes
- ✅ No hardcoded URLs or passwords
- ✅ Proper use of .env files for config
- ⚠️ For database changes: Is migration included?
- ⚠️ For new endpoints: Is it documented in IG_SPECIFICATIONS_API.md?

---

## 🎓 Learning Resources

**Architecture**:
- Read `01_SPECIFICATIONS_ARCHITECTURE.md` first (understand overall system)
- Then read `02_SPECIFICATIONS_REQUIREMENTS.md` (understand what to build)
- Then read specific IG_* file for your component

**Farmer-Centric Design**:
- See `IG_TOKKATOT_2.0_FARMER_CENTRIC_SPECIFICATIONS.md`
- Rule: If a feature requires more than 2 clicks or shows > 5 options, it's too complex

**Technology Decisions**:
- See `OG_SPECIFICATIONS_TECHNOLOGY_STACK.md` (why Go, Vue, PyTorch, etc)

---

## � AI Context Files (Component-Specific Guides)

**Each service folder has an AI_CONTEXT.md** for component-specific implementation details.

### **Purpose**: Deep-dive into specific tech stack, file structure, common patterns

**Component AI Context Files**:

1. **`middleware/AI_CONTEXT.md`** - Go API Gateway
   - File structure, endpoint handlers, database patterns
   - Schedule automation implementation (action_duration, action_sequence)
   - WebSocket broadcasting, JWT auth flow
   - Common development tasks (add endpoint, add table, test locally)

2. **`frontend/AI_CONTEXT.md`** - Vue.js 3 PWA
   - Component structure, page organization, CSS patterns
   - Farmer accessibility (48px fonts, high contrast, Khmer support)
   - WebSocket client, API integration
   - Mobile-first responsive design

3. **`ai-service/AI_CONTEXT.md`** - PyTorch Disease Detection
   - Model architecture (EfficientNetB0 + DenseNet121 ensemble)
   - FastAPI endpoints (/predict, /health)
   - Training pipeline, evaluation metrics
   - Docker deployment

4. **`embedded/AI_CONTEXT.md`** - ESP32 Firmware
   - Sensor drivers (DHT22, ultrasonic, relay control)
   - MQTT protocol, message formats
   - FreeRTOS tasks, memory management
   - OTA updates

5. **`docs/AI_CONTEXT.md`** - Documentation Maintenance
   - When to update which docs
   - Documentation standards (farmer-first language, concrete examples)
   - Workflow: code change → identify affected docs → update with examples
   
**Reading Order**:
1. **First**: Read `/AI_INSTRUCTIONS.md` (this file) for project overview
2. **Then**: Read component `AI_CONTEXT.md` for specific tech implementation
3. **Reference**: `docs/implementation/*.md` for complete API/DB/Security specs

---
## 📝 Advanced: Documentation Update Details

**This section provides additional context** on the documentation update requirement explained earlier. See "MANDATORY: Self-Documentation After Building Features" section for main instructions.

### Documentation File Structure

**Implementation Specs** (`docs/implementation/*.md`):
- `API.md` - All 35+ REST endpoints with request/response examples
- `DATABASE.md` - Complete schema (10 tables, all fields, indexes)
- `FRONTEND.md` - Vue.js 3 pages, components, routing, API integration
- `AI_SERVICE.md` - PyTorch models, FastAPI endpoints, preprocessing
- `EMBEDDED.md` - ESP32 firmware, MQTT protocol, GPIO pins, sensors
- `SECURITY.md` - JWT auth, input validation, rate limiting

**AI Knowledge Files**:
- `AI_INSTRUCTIONS.md` (this file) - Master guide for entire project
- `middleware/AI_CONTEXT.md` - Go patterns, schedule automation, database queries
- `frontend/AI_CONTEXT.md` - Vue.js components, mobile-first design, API calls
- `ai-service/AI_CONTEXT.md` - PyTorch ensemble, preprocessing, FastAPI health checks
- `embedded/AI_CONTEXT.md` - MQTT command execution, sensor reading, safety logic

**Use Case Documentation**:
- `docs/AUTOMATION_USE_CASES.md` - Real farmer scenarios (Sokha's pulse feeding, Dara's climate control) with JSON examples and benefits

### Cross-Referencing Example

**When documenting `action_sequence` feature**:
1. `DATABASE.md` shows field type: `action_sequence JSONB`
2. `API.md` shows endpoint usage: `POST /farms/{id}/schedules` with JSON example
3. `FRONTEND.md` shows Action Sequence Builder UI component
4. `EMBEDDED.md` shows ESP32 execution code for multi-step patterns
5. `AUTOMATION_USE_CASES.md` shows Farmer Sokha's pulse feeding scenario
6. Each file links to the others using `[docs/AUTOMATION_USE_CASES.md](AUTOMATION_USE_CASES.md)` syntax

### Documentation Verification

**Before marking work complete, verify**:
- ✅ All JSON examples are valid (paste into validator)
- ✅ All SQL queries work (test in psql/sqlite)
- ✅ All code snippets compile (Go: `go build`, Python: `python -m py_compile`)
- ✅ All links work (click every `[text](path.md)` link)
- ✅ Examples use real data (not placeholders like `<farm_id>`)

### Version History Tracking

**After major features, update `docs/README.md`**:
```markdown
## Version History

### v2.0 (Current)
- **Multi-step automation**: action_sequence field for pulse feeding, conveyor operations
- **Auto-turn-off timers**: action_duration for water pumps, lights
- **AI documentation consolidation**: Removed 2,310 lines of duplication across 5 files
```

This provides timeline context for future developers.

---
## �💬 Communication Style

- **Document everything** - If code is unclear, it's not production-ready
- **Update specs first** - Write API spec before coding the endpoint
- **Ask questions** - If a requirement is ambiguous, ask the team (don't guess)
- **Prioritize farmer experience** - "Can a 65-year-old farmer use this?" is the question

---

## 📞 When to Ask for Help

**Red Flags - Ask a Human**:
- ❓ "Should this feature be simple or complex?" → Ask product owner
- ❓ "How should this edge case be handled?" → Ask team architect  
- ❓ "Is this security risk acceptable?" → Ask security team
- ❓ "Do we need a database migration?" → Ask DBA or tech lead
- ❓ Requirements conflict with existing spec → Ask for clarification

**You Can Handle**:
- ✅ Implementing specified endpoints exactly as documented
- ✅ Fixing bugs that violate specifications
- ✅ Refactoring code that works but is messy
- ✅ Writing tests for existing functionality
- ✅ Optimizing performance within spec

---

## ✨ Best Practices

1. **Write tests** - Every new feature should have corresponding tests
2. **Use type safety** - Go struct types, Python type hints, JavaScript PropTypes
3. **Comment complex logic** - Especially business rules for farmers
4. **Version your APIs** - Use `/v1`, `/v2` for API versioning
5. **Monitor errors** - Use structured logging with timestamps & request IDs
6. **Think edge cases** - What if network disconnects? Device offline? Image corrupted?
7. **Performance matters** - Pages should load <2s on 4G, predictions <3s on CPU

---

## 🎯 Summary

You are not just coding—you are **building a system for elderly farmers in Cambodia**. Every decision should prioritize:

1. **Simplicity** - Can a non-technical farmer understand it?
2. **Reliability** - Will it work without internet? What if a device goes offline?
3. **Security** - Are secrets protected? Are farmer photos secure?
4. **Accessibility** - Can someone with poor eyesight use it?

**Golden Rule**: When in doubt, choose the simplest solution that works.

---

**Now go build! 🚀**
