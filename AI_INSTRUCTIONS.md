# 🤖 AI Agent Instructions for Tokkatot 2.0

**Last Updated**: February 19, 2026  
**Version**: 2.0  
**Purpose**: Guide for AI agents assisting with Tokkatot development

---

## 📋 Project Overview

**Tokkatot 2.0** is a cloud-based smart farming IoT system designed for elderly Cambodian farmers with low digital literacy.

**Key Context:**
- **Target Users**: Elderly Cambodian farmers (Khmer/English, simple UI)
- **Main Hardware**: ESP32 microcontrollers (devices), Raspberry Pi 4B (local hub)
- **Cloud Stack**: Go (API), Vue.js 3 (Frontend), PostgreSQL (Data), InfluxDB (Time-series), Python FastAPI (AI)
- **Primary Feature**: Remote device control + AI disease detection from chicken feces images
- **Architecture**: 3-tier (Client/API/Data) + Edge computing (offline capability)

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
- `IG_SPECIFICATIONS_DATABASE.md` - PostgreSQL schema (13 tables, farmer-centric)
- `IG_SPECIFICATIONS_SECURITY.md` - JWT auth, Email/Phone login, no MFA for farmers
- `IG_SPECIFICATIONS_FRONTEND.md` - Vue.js UI, 48px+ fonts, WCAG AAA accessibility
- `IG_SPECIFICATIONS_EMBEDDED.md` - ESP32 firmware, MQTT protocol
- `IG_SPECIFICATIONS_AI_SERVICE.md` - PyTorch ensemble (99% accuracy)

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
- Pull requests must include documentation updates
- GitHub Actions can run tests/linting on PR

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

## 💬 Communication Style

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
