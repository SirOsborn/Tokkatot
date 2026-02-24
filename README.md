# 🐔 Tokkatot - Smart Poultry Farm Management System

<div align="center">

![Tokkatot Logo](frontend/assets/images/tokkatot%20logo-02.png)

**Advanced IoT-Based Poultry Disease Detection & Farm Automation**

🌐 **Website:** [https://tokkatot.aztrolabe.com](https://tokkatot.aztrolabe.com)  
📧 **Email:** [tokkatot.info@gmail.com](mailto:tokkatot.info@gmail.com)

[![License](https://img.shields.io/badge/license-Proprietary-red.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)](https://go.dev/)
[![Python](https://img.shields.io/badge/Python-3.8+-3776AB?logo=python)](https://www.python.org/)
[![TensorFlow](https://img.shields.io/badge/TensorFlow-2.x-FF6F00?logo=tensorflow)](https://www.tensorflow.org/)

[Features](#-features) • [Quick Start](#-quick-start) • [Architecture](#-architecture) • [Documentation](#-documentation)

</div>

---

## 📖 About

**Tokkatot** is a comprehensive smart poultry management system designed for Cambodian farmers. It combines IoT sensor technology, AI-powered disease detection, and automated farm controls to improve poultry health monitoring and farm productivity.

### 🎯 Key Capabilities

- 🤖 **AI Disease Detection** - Identify chicken diseases from droppings using EfficientNetB0 deep learning
- 📊 **Real-time Monitoring** - Track temperature, humidity, and environmental conditions
- 🎮 **Remote Control** - Manage lighting, feeding, ventilation, and water systems
- 📱 **Mobile-First** - Progressive Web App accessible from any device
- 🔒 **Secure** - JWT authentication and encrypted IoT communication
- 🌐 **Offline-Ready** - Local network operation, no internet required

---

## ✨ Features

### 🩺 AI Disease Detection System

- **EfficientNetB0 Model** - State-of-the-art CNN architecture
- **5 Disease Classes** - Healthy, Coccidiosis, Salmonella, E.coli, Newcastle
- **Real-time Analysis** - Upload photos, get instant diagnosis
- **Confidence Scoring** - Know how reliable each prediction is
- **Treatment Recommendations** - Actionable advice for each condition

### 🌡️ Environmental Monitoring

- **Temperature Tracking** - Real-time temperature monitoring with history
- **Humidity Control** - Track and manage humidity levels
- **Historical Data** - View trends over time with interactive charts
- **Automated Alerts** - Get notified when conditions are suboptimal

### 🎛️ Smart Farm Controls

- **Automated Lighting** - Schedule and control coop lighting
- **Smart Feeding** - Automated feeder with manual override
- **Ventilation Control** - Automated fan based on temperature
- **Water Management** - Monitor and control water levels
- **Waste Management** - Automated conveyor belt control

### 📱 Progressive Web App

- **Mobile Optimized** - Responsive design for phones and tablets
- **Offline Support** - Service worker caching for reliability
- **Fast Loading** - Optimized assets and critical CSS
- **Khmer Language** - Full support for Cambodian users

---

## � Version

- **Current Release:** v1.0 (Prototype - Raspberry Pi local mode)
- **In Development:** v2.0 (Production - Cloud-based) - See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) and [docs/README.md](docs/README.md)

---

## 🚀 Quick Start

### v1.0 (Current - Local Mode)

**Requirements:**
- Raspberry Pi 4B+ (2GB RAM min) or Ubuntu Server 20.04+
- WiFi connectivity or Ethernet
- ESP32 microcontroller (pre-flashed)
- Docker (recommended) or direct installation

**Installation (Automatic):**
```bash
# Clone repository
git clone https://github.com/SirOsbornOjr/tokkatot.git
cd tokkatot

# Run setup script
bash scripts/deploy-all.sh
```

**Manual Installation:**
1. Set up Raspberry Pi OS or Ubuntu Server 20.04+
2. Configure WiFi Access Point: `bash scripts/setup-access-point.sh`
3. Deploy Middleware: `bash scripts/setup-middleware-service.sh`
4. Deploy AI Service: See `ai-service/README.md`
5. Access: Open browser at `http://10.0.0.1:3000`

**Default Credentials (v1.0 - Demo):**
- Email: `test@tokkatot.local`
- Password: `Password123`

### v2.0 (Production - Cloud-Based)

**Status:** In development - See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) and [docs/guides/SETUP.md](docs/guides/SETUP.md)

v2.0 adds:
- ☁️ Cloud connectivity (DigitalOcean + Kubernetes)
- 🌍 Multi-farm support  
- 📱 Mobile PWA (Vue.js 3 CDN)
- 🔔 Real-time updates (WebSocket)
- 📊 Advanced analytics (InfluxDB)
- 🇰🇭 Khmer/English language support
- ♿ Accessibility for elderly farmers (48px+ touch targets, WCAG AAA)
- 🔑 FREE registration (key system, no SMS costs)

**Tech Stack:** Go 1.23 + Fiber v2, PostgreSQL 17, Vue.js 3 (CDN), Python 3.12 + PyTorch  
**Details:** [docs/TECH_STACK.md](docs/TECH_STACK.md)

---


## 🏗️ Architecture

### System Overview

```
┌─────────────────────────────────────────────┐
│  Raspberry Pi / Ubuntu Server (10.0.0.1)   │
│                                             │
│  ┌─────────────────────────────────────┐  │
│  │ WiFi Access Point                   │  │
│  │ (hostapd + dnsmasq)                │  │
│  └─────────────────────────────────────┘  │
│                                             │
│  ┌─────────────────────────────────────┐  │
│  │ Go Middleware :4000                 │  │
│  │ • Fiber Web Framework               │  │
│  │ • JWT Authentication                │  │
│  │ • SQLite Database                   │  │
│  │ • API Gateway                       │  │
│  └─────────────────────────────────────┘  │
│                                             │
│  ┌─────────────────────────────────────┐  │
│  │ AI Service :5000                    │  │
│  │ • Python Flask                      │  │
│  │ • TensorFlow/Keras                  │  │
│  │ • EfficientNetB0 Model              │  │
│  └─────────────────────────────────────┘  │
└─────────────────────────────────────────────┘
         │                    │
         │ WiFi              │ WiFi
         ▼                    ▼
┌─────────────────┐  ┌─────────────────┐
│  ESP32          │  │  Mobile Device  │
│  (10.0.0.2)     │  │  (10.0.0.X)     │
│                 │  │                 │
│  • DHT22        │  │  • Web Browser  │
│  • Sensors      │  │  • PWA          │
│  • Actuators    │  │  • Dashboard    │
└─────────────────┘  └─────────────────┘
```

### Tech Stack

#### Backend
- **Language:** Go 1.23+
- **Framework:** Fiber v2
- **Database:** SQLite (modernc.org/sqlite)
- **Authentication:** JWT (golang-jwt/jwt)
- **Encryption:** golang.org/x/crypto

#### AI Service
- **Language:** Python 3.8+
- **Framework:** Flask
- **ML Library:** TensorFlow 2.x / Keras
- **Model:** EfficientNetB0
- **Image Processing:** OpenCV, PIL

#### Frontend
- **Languages:** HTML5, CSS3, JavaScript (ES6+)
- **Architecture:** Progressive Web App (PWA)
- **Charts:** Chart.js
- **Icons:** Font Awesome
- **Fonts:** Kantumruy (Khmer)

#### IoT/Embedded
- **Hardware:** ESP32
- **Sensors:** DHT22 (temperature/humidity)
- **Actuators:** Relays, Servos, Motors
- **Communication:** WiFi, Encrypted HTTP

---

## 📁 Project Structure

```
tokkatot/
├── docs/                        # Documentation (v2.0 - Organized)
│   ├── ARCHITECTURE.md          # ← START HERE! System design, coop-centric model
│   ├── TECH_STACK.md            # Technology choices (Go, Vue.js, PostgreSQL)
│   ├── README.md                # Documentation navigation hub
│   │
│   ├── guides/                  # Setup & Installation
│   │   └── SETUP.md             # Complete setup guide (PostgreSQL, Go, frontend)
│   │
│   ├── implementation/          # Component Development
│   │   ├── API.md               # Backend API (Go + Fiber, 67 endpoints)
│   │   ├── DATABASE.md          # Database schema (PostgreSQL, 8 tables)
│   │   ├── FRONTEND.md          # Frontend (Vue.js 3 migration guide)
│   │   ├── AI_SERVICE.md        # AI service (Python + PyTorch, disease detection)
│   │   ├── EMBEDDED.md          # ESP32 firmware (C/ESP-IDF)
│   │   └── SECURITY.md          # Authentication & security (JWT, registration keys)
│   │
│   └── troubleshooting/         # Problem Solving
│       ├── DATABASE.md          # Database connection issues
│       └── API_TESTING.md       # Test backend endpoints
│
├── frontend/                    # Vue.js 3 PWA (Progressive Web App)
│   ├── pages/                   # HTML pages (login, dashboard, coops)
│   ├── components/              # Vue components (navbar, header, coop-card)
│   ├── js/                      # Vue apps, API helpers, WebSocket
│   ├── css/                     # Styles (mobile-first, 48px+ touch targets)
│   ├── assets/                  # Images, fonts, icons
│   ├── manifest.json            # PWA manifest
│   └── sw.js                    # Service worker (offline support)
│
├── middleware/                  # Go backend (REST API + JWT auth)
│   ├── main.go                  # Server entry point
│   ├── api/                     # HTTP handlers
│   │   ├── authentication.go    # Login, signup, registration keys
│   │   ├── profiles.go          # User profiles
│   │   ├── data-handler.go      # IoT sensor data
│   │   └── disease-detection.go # AI integration (calls ai-service)
│   ├── database/                # Database layer
│   │   └── sqlite3_db.go        # SQLite wrapper (production: PostgreSQL)
│   ├── utils/                   # JWT, validation, response helpers
│   ├── go.mod                   # Go 1.23, Fiber v2.52, JWT v4.5
│   ├── .env.example             # Environment template (COMMIT THIS!)
│   └── .env                     # Secrets file (NEVER COMMIT!)
│
├── ai-service/                  # Python AI service (FastAPI + PyTorch)
│   ├── app.py                   # FastAPI server (port 8000)
│   ├── inference.py             # ChickenDiseaseDetector (ensemble model)
│   ├── models.py                # EfficientNetB0 + DenseNet121
│   ├── data_utils.py            # Image preprocessing, class definitions
│   ├── outputs/                 # PROPRIETARY: Model files (*.pth)
│   │   └── ensemble_model.pth   # 47.2 MB trained model (NOT in git)
│   ├── docker-compose.yml       # Docker deployment
│   ├── Dockerfile               # Python 3.12-slim image
│   └── requirements.txt         # PyTorch, FastAPI, Uvicorn
│
├── embedded/                    # ESP32 firmware (C/ESP-IDF)
│   ├── main/                    # Device boot, MQTT client
│   │   └── main.c               # Entry point
│   ├── components/              # Custom components
│   │   └── dht/                 # DHT22 sensor driver
│   ├── CMakeLists.txt           # ESP-IDF build config
│   └── sdkconfig                # ESP-IDF configuration
│
├── certs/                       # SSL certificates (self-signed)
├── generate-cert.sh             # Certificate generation script
├── AI_INSTRUCTIONS.md           # Master AI agent guide
├── LICENSE                      # Proprietary license
└── README.md                    # This file
```

---

## 📚 Documentation

### 🎯 Start Here - New Developer Onboarding

**New to the project?** Read these in order:

1. **[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)** - Understand the coop-centric system design
2. **[docs/TECH_STACK.md](docs/TECH_STACK.md)** - Why Go, Vue.js 3, PostgreSQL, PyTorch
3. **[docs/guides/SETUP.md](docs/guides/SETUP.md)** - Install PostgreSQL, build backend, run frontend

### 📖 Documentation Structure

Documentation is organized by purpose (not file naming conventions):

#### 🔑 Core Concepts
| Document | What You'll Learn |
|----------|-------------------|
| [ARCHITECTURE.md](docs/ARCHITECTURE.md) | Coop-centric design, data hierarchy (User→Farm→Coop→Device), physical infrastructure, user flows |
| [TECH_STACK.md](docs/TECH_STACK.md) | Technology decisions (Go vs Node.js, Vue.js vs React), deployment strategy (single VPS, not microservices) |

#### 🚀 Setup & Installation
| Document | What You'll Learn |
|----------|-------------------|
| [guides/SETUP.md](docs/guides/SETUP.md) | Install prerequisites (Go, PostgreSQL), configure `.env`, build backend, test API, troubleshooting |

#### 💻 Component Development
| Document | What You'll Learn |
|----------|-------------------|
| [implementation/API.md](docs/implementation/API.md) | Backend API endpoints (66 total), JWT authentication, request/response formats |
| [implementation/DATABASE.md](docs/implementation/DATABASE.md) | PostgreSQL schema (8 tables), indexes, migrations, performance tuning |
| [implementation/FRONTEND.md](docs/implementation/FRONTEND.md) | Vue.js 3 migration guide (3 phases), component patterns, WebSocket real-time, accessibility for farmers |
| [implementation/AI_SERVICE.md](docs/implementation/AI_SERVICE.md) | Disease detection service (Python + PyTorch), ensemble model, API integration |
| [implementation/EMBEDDED.md](docs/implementation/EMBEDDED.md) | ESP32 firmware (C/ESP-IDF), sensor drivers (DHT22), MQTT communication |
| [implementation/SECURITY.md](docs/implementation/SECURITY.md) | JWT authentication, registration key system, RBAC (Owner/Manager/Viewer) |

#### 🔧 Troubleshooting
| Document | What You'll Learn |
|----------|-------------------|
| [troubleshooting/DATABASE.md](docs/troubleshooting/DATABASE.md) | Fix connection errors, schema sync issues, migration failures |
| [troubleshooting/API_TESTING.md](docs/troubleshooting/API_TESTING.md) | Test backend endpoints, debug API errors, PowerShell scripts |

### 📋 Complete Documentation Index

**See [docs/README.md](docs/README.md) for full navigation** (includes all guides, legacy docs, project management)
 
---

## 🖥️ System Requirements

### Server (Raspberry Pi / Ubuntu)

- **OS:** Ubuntu Server 20.04+ or Raspberry Pi OS (64-bit)
- **RAM:** 2GB minimum (4GB recommended)
- **Storage:** 20GB minimum
- **Network:** WiFi adapter with AP mode support
- **CPU:** 2 cores recommended

### ESP32 Device

- **Microcontroller:** ESP32 (ESP32-WROOM-32)
- **Sensors:** DHT22, Water level sensor
- **Actuators:** Relays, Servo motor
- **Power:** 5V 2A minimum

### Client Device

- **Browser:** Chrome 90+, Firefox 88+, Safari 14+
- **Connection:** WiFi (2.4GHz)
- **Resolution:** 360px+ width

---

## 🔐 Security

### Authentication

- **JWT Tokens** - Secure session management
- **Password Hashing** - bcrypt with salt
- **Cookie Security** - HttpOnly, SameSite

### IoT Communication

- **AES-128-GCM** - Encrypted data transmission
- **Challenge-Response** - Prevents replay attacks
- **SHA-256** - Message integrity verification

### Best Practices

- Change default JWT secret in production
- Use HTTPS with valid certificates
- Regular security updates
- Firewall configuration
- Rate limiting on API endpoints

---

## 🌐 API Endpoints

### Authentication
- `POST /login` - User login
- `POST /register` - User registration
- `POST /logout` - User logout

### User Management
- `GET /api/profile` - Get user profile
- `POST /api/profile` - Update user profile

### IoT Data
- `GET /api/get-initial-state` - Get all device states
- `GET /api/get-current-data` - Get current sensor data
- `GET /api/get-historical-data` - Get historical data

### Device Control
- `GET /api/toggle-auto` - Toggle automation mode
- `GET /api/toggle-fan` - Toggle ventilation
- `GET /api/toggle-bulb` - Toggle lighting
- `GET /api/toggle-feeder` - Toggle feeder
- `GET /api/toggle-water` - Toggle water pump
- `GET /api/toggle-belt` - Toggle conveyor belt

### AI Disease Detection
- `GET /api/ai/health` - Check AI service health
- `POST /api/ai/predict-disease` - Predict disease from image
- `GET /api/ai/disease-info` - Get disease information

---

## 📝 License & Usage

This project is **proprietary software** developed for Tokkatot Startup. See the [LICENSE](LICENSE) file for complete terms and conditions.

**Unauthorized copying, modification, distribution, or commercial use is strictly prohibited.**

---

## 🗺️ Roadmap

### Version 1.0 ✅ (Current - Prototype/Local)
- ✅ Basic farm monitoring (local WiFi AP)
- ✅ AI disease detection (EfficientNetB0)
- ✅ Manual device control (web UI)
- ✅ SQLite local database
- ✅ ESP32 sensor integration

### Version 2.0 🚀 (In Development - Production)
- 🔄 **Backend:** Go REST API (67 endpoints, JWT auth, PostgreSQL)
- 🔄 **Frontend:** Vue.js 3 (CDN build, WCAG AAA, Khmer language)
- 🔄 **AI Service:** PyTorch ensemble (EfficientNetB0 + DenseNet121)
- 🔄 **Embedded:** ESP32 firmware (MQTT, OTA updates)
- 🔄 **Cloud:** DigitalOcean deployment (Kubernetes, Docker)

**Major v1.0 → v2.0 Changes:**
- Cloud-connected (not just local WiFi AP)
- Multi-farm support (not single farm only)
- Real-time monitoring (WebSocket + MQTT)
- Remote device control (from anywhere)
- OTA firmware updates (no farm visits)
- Farmer-centric UI (48px+ fonts, high contrast, Khmer language)
- Registration key system (FREE verification, no SMS costs)
- 99.5% uptime target
- 5-year data retention

**Documentation:** See [docs/README.md](docs/README.md) for complete v2.0 specifications, project timeline, and team structure

---

## 🎨 Design & Resources

### UI Design Pages (v1.0)
Mockups and design specifications for version 1.0 interface:
- Dashboard Layout
- Device Control Interface
- Farm Settings Page
- Disease Detection UI

See `design/` folder for detailed page designs and Figma exports.

<div align="center">

**Proprietary Software - Tokkatot Startup**

For internal use only. Unauthorized copying, modification, or distribution is prohibited.

</div>
