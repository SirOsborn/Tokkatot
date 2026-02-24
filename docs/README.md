# Documentation - Tokkatot 2.0

**Last Updated**: February 23, 2026  
**Structure**: Organized by purpose (guides, implementation, troubleshooting)

---

## 📌 Quick Start

### **New to Tokkatot?** Read these in order:
1. **[ARCHITECTURE.md](ARCHITECTURE.md)** - Coop-centric system design, data hierarchy, user flows
2. **[TECH_STACK.md](TECH_STACK.md)** - Why Go, Vue.js 3, PostgreSQL, PyTorch
3. **[guides/SETUP.md](guides/SETUP.md)** - Install PostgreSQL, build backend, run frontend

### **For AI Agents:**
- **[../AI_INSTRUCTIONS.md](../AI_INSTRUCTIONS.md)** - Master AI agent guide (repository root)
- **[AI_CONTEXT.md](AI_CONTEXT.md)** - AI context for this docs folder

---

## 📂 Documentation Structure

```
docs/
├── ARCHITECTURE.md          ← START HERE! System design
├── TECH_STACK.md            Technology decisions
├── AUTOMATION_USE_CASES.md  🚜 Real-world farmer scenarios
├── README.md                This file
│
├── guides/                  Setup & Installation
│   └── SETUP.md             Complete setup guide
│
├── implementation/          Component Development
│   ├── API.md               Backend API (67 endpoints)
│   ├── DATABASE.md          PostgreSQL schema (14 tables)
│   ├── FRONTEND.md          Vue.js 3 migration
│   ├── AI_SERVICE.md        Disease detection (PyTorch)
│   ├── EMBEDDED.md          ESP32 firmware
│   └── SECURITY.md          JWT auth, registration keys
│
└── troubleshooting/         Problem Solving
    ├── DATABASE.md          Connection issues, schema fixes
    └── API_TESTING.md       Test backend endpoints
```

---

## 📖 Core Documentation

### 🏗️ System Design
| Document | What You'll Learn |
|----------|-------------------|
| **[ARCHITECTURE.md](ARCHITECTURE.md)** | Coop-centric design, physical infrastructure, User→Farm→Coop→Device hierarchy, user journey flows, device control examples |
| **[TECH_STACK.md](TECH_STACK.md)** | Go vs Node.js comparison, Vue.js 3 CDN strategy, PostgreSQL schema design, single VPS deployment (not microservices) |
| **[AUTOMATION_USE_CASES.md](AUTOMATION_USE_CASES.md)** | 🚜 **Real-world farmer automation scenarios**: Conveyor cycling, pulse feeding (multi-step sequences), sensor-driven water pumps, climate control. Detailed JSON examples for all 4 schedule types. **CRITICAL for understanding schedule feature design.** |

---

## 🚀 Guides (Setup & Installation)

| Document | Purpose | When to Use |
|----------|---------|-------------|
| **[guides/SETUP.md](guides/SETUP.md)** | Complete installation guide | First-time setup: PostgreSQL 17, Go 1.23, backend build, `.env` configuration, troubleshooting |

**What's in SETUP.md:**
- Prerequisites (Go, PostgreSQL, Git, VS Code)
- Database setup (create `tokkatot` database)
- Backend configuration (`.env` file)
- Build & run instructions
- Frontend setup (Vue.js 3 CDN)
- AI service (optional)
- Troubleshooting common issues
- Development workflow

---

## 💻 Implementation Guides (Component Development)

### Backend (Go + Fiber)
**[implementation/API.md](implementation/API.md)**
- 66 REST API endpoints (authentication, devices, AI, farms, coops)
- JWT authentication flow
- Request/response formats
- Error handling
- Middleware (auth, logging, CORS)

### Database (PostgreSQL)
**[implementation/DATABASE.md](implementation/DATABASE.md)**
- Schema design (14 tables: users, farms, farm_users, coops, devices, device_configurations, schedules, schedule_executions, device_commands, event_logs, alerts, alert_subscriptions, device_readings, etc.)
- Indexes & performance tuning
- Simplified RBAC (Farmer, Viewer)
- Migrations
- Query examples

### Frontend (Vue.js 3)
**[implementation/FRONTEND.md](implementation/FRONTEND.md)**
- **3-Phase Migration Strategy:**
  - Phase 1: CDN setup (no build step)
  - Phase 2: Component system (navbar, coop-card)
  - Phase 3: Vite build (optional, for production)
- Component patterns (authentication, API helpers, WebSocket)
- Accessibility for farmers (48px+ touch targets, high contrast, Khmer language)
- Mobile-first design (PWA)

### AI Service (Python + PyTorch)
**[implementation/AI_SERVICE.md](implementation/AI_SERVICE.md)**
- FastAPI server (port 8000)
- Ensemble model (EfficientNetB0 + DenseNet121)
- Disease detection (5 classes: Healthy, Coccidiosis, Salmonella, E.coli, Newcastle)
- `/predict` endpoint (image upload → disease diagnosis)
- Model training & evaluation
- Docker deployment

### Embedded (ESP32 / Raspberry Pi)
**[implementation/EMBEDDED.md](implementation/EMBEDDED.md)**
- ESP32 firmware (C/ESP-IDF)
- Sensor drivers (DHT22 temperature/humidity)
- MQTT communication
- Raspberry Pi controller
- OTA updates
- Device registration

### Security (Authentication & Authorization)
**[implementation/SECURITY.md](implementation/SECURITY.md)**
- JWT authentication (access tokens, refresh tokens)
- Registration key system (FREE verification, no SMS costs)
- RBAC (Farmer/Viewer)
- Password hashing (bcrypt)
- TLS/SSL certificates
- API security best practices

---

## 🔧 Troubleshooting

| Document | Common Issues Solved |
|----------|----------------------|
| **[troubleshooting/DATABASE.md](troubleshooting/DATABASE.md)** | PostgreSQL connection errors, schema sync issues, migration failures, `tokkatot` database not found |
| **[troubleshooting/API_TESTING.md](troubleshooting/API_TESTING.md)** | Test authentication endpoints, debug JWT tokens, PowerShell test scripts, common API errors |

---

## 🎯 Task-Based Navigation

### **I want to...**

#### 🛠️ Build & Setup
- **Set up the backend** → [guides/SETUP.md](guides/SETUP.md)
- **Fix database connection** → [troubleshooting/DATABASE.md](troubleshooting/DATABASE.md)
- **Test API endpoints** → [troubleshooting/API_TESTING.md](troubleshooting/API_TESTING.md)
- **Configure .env file** → [guides/SETUP.md](guides/SETUP.md) (Backend Configuration section)

#### 🧠 Understand the System
- **Understand the full system** → [ARCHITECTURE.md](ARCHITECTURE.md)
- **See the database schema** → [implementation/DATABASE.md](implementation/DATABASE.md)
- **Understand coop-device relationship** → [ARCHITECTURE.md](ARCHITECTURE.md) (Physical Infrastructure & Data Hierarchy sections)
- **Know the tech choices** → [TECH_STACK.md](TECH_STACK.md)
- **🚜 Understand farmer automation (schedules)** → [AUTOMATION_USE_CASES.md](AUTOMATION_USE_CASES.md) - **Conveyor cycles, pulse feeding, multi-step sequences, sensor-driven pumps**

#### 🔨 Implement Features
- **Add new API endpoint** → [implementation/API.md](implementation/API.md)
- **Add new database table** → [implementation/DATABASE.md](implementation/DATABASE.md)
- **Work on schedules/automation** → [AUTOMATION_USE_CASES.md](AUTOMATION_USE_CASES.md) + [implementation/DATABASE.md](implementation/DATABASE.md) (schedules table)
- **Work on disease detection** → [implementation/AI_SERVICE.md](implementation/AI_SERVICE.md)
- **Build frontend pages** → [implementation/FRONTEND.md](implementation/FRONTEND.md)
- **Program ESP32/Raspberry Pi** → [implementation/EMBEDDED.md](implementation/EMBEDDED.md)
- **Migrate to Vue.js 3** → [implementation/FRONTEND.md](implementation/FRONTEND.md) (3-phase strategy)

#### 🔐 Security & Auth
- **Understand registration keys** → [implementation/SECURITY.md](implementation/SECURITY.md) + [ARCHITECTURE.md](ARCHITECTURE.md) (Authentication section)
- **Implement role-based access** → [implementation/SECURITY.md](implementation/SECURITY.md) (RBAC section)
- **Debug JWT tokens** → [implementation/API.md](implementation/API.md) (Authentication endpoints)
- **Set up TLS certificates** → [implementation/SECURITY.md](implementation/SECURITY.md) + [guides/SETUP.md](guides/SETUP.md)

#### 🐛 Debug & Troubleshoot
- **Database won't connect** → [troubleshooting/DATABASE.md](troubleshooting/DATABASE.md)
- **API returns 500 error** → [troubleshooting/API_TESTING.md](troubleshooting/API_TESTING.md)
- **JWT token invalid** → [implementation/API.md](implementation/API.md) (Authentication section)
- **Backend won't start** → [guides/SETUP.md](guides/SETUP.md) (Troubleshooting section)

---

## 🔄 Documentation Version History

### v2.0 (February 23, 2026)
- ✅ Reorganized structure (23 files → 14 organized files)
- ✅ Created `guides/`, `implementation/`, `troubleshooting/` folders
- ✅ Consolidated architecture docs (3 files → ARCHITECTURE.md)
- ✅ Consolidated setup guides (5 files → guides/SETUP.md)
- ✅ Created TECH_STACK.md (Go vs Node.js, Vue.js 3 strategy)
- ✅ Created implementation/FRONTEND.md (Vue.js 3 migration guide)
- ✅ Moved troubleshooting guides to dedicated folder
- ✅ **NEW: AUTOMATION_USE_CASES.md** - Real-world farmer scenarios (conveyor cycling, pulse feeding, multi-step sequences, sensor-driven pumps)

### v2.1 (February 24, 2026) - Current
- ✅ **Temperature Monitoring Dashboard** — `monitoring.html` + `TemperatureTimelineHandler`
- ✅ API.md updated to v2.2 — 67 endpoints, full temperature-timeline spec with bg_hint table
- ✅ FRONTEND.md — `/monitoring` page added to page inventory
- ✅ AI_INSTRUCTIONS.md updated to v2.3
- ✅ All component AI_CONTEXT.md files updated

### v1.0 (Legacy)
- Flat file structure (23 files in docs/)
- IG_* naming convention (Implementation Guides)
- OG_* naming convention (Original/Operational Guides)
- 00_SPECIFICATIONS_INDEX.md navigation hub

---

## 📞 Support

### Documentation Issues
- **Missing information?** Check [ARCHITECTURE.md](ARCHITECTURE.md) or [TECH_STACK.md](TECH_STACK.md) first
- **Broken links?** Report in GitHub Issues
- **Need more examples?** See [implementation/](implementation/) folder

### Technical Support
- **Email**: tokkatot.info@gmail.com
- **GitHub Issues**: Bug reports & feature requests
- **Live Chat**: (Coming in v2.1)

---

## 📝 Contributing to Documentation

### Adding New Documentation
1. Determine category: `guides/`, `implementation/`, or `troubleshooting/`
2. Follow existing file structure (clear headings, code examples, tables)
3. Update this README.md with new file reference
4. Keep AI_INSTRUCTIONS.md synchronized

### Documentation Standards
- ✅ Use tables for comparisons (Go vs Node.js, Phase 1 vs Phase 2)
- ✅ Include code examples (before/after, good/bad patterns)
- ✅ Write for farmers & developers (simple language, clear steps)
- ✅ Add troubleshooting sections (common errors, solutions)
- ❌ Avoid long paragraphs (use bullet points, numbered lists)
- ❌ No "yapping" (keep it concise and actionable)

---

**Last Updated**: February 23, 2026  
**Maintained by**: Tokkatot Development Team  
**License**: Proprietary - See [../LICENSE](../LICENSE)
