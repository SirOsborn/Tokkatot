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

- **Current Release:** v1.0 (Prototype)
- **In Development:** v2.0 (Production) - See [docs/00_SPECIFICATIONS_INDEX.md](docs/00_SPECIFICATIONS_INDEX.md)

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

**Status:** In development - See [docs/00_SPECIFICATIONS_INDEX.md](docs/00_SPECIFICATIONS_INDEX.md)

v2.0 will add:
- ☁️ Cloud connectivity (DigitalOcean)
- 🌍 Multi-farm support  
- 📱 Mobile app + PWA
- 🔔 Push notifications
- 📊 Advanced analytics
- 🇰🇭 Khmer/English language support
- ♿ Accessibility for elderly farmers

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
├── docs/                        # Documentation & Specifications (v2.0)
│   ├── 00_SPECIFICATIONS_INDEX.md           # Start here - navigation hub
│   ├── 01_SPECIFICATIONS_ARCHITECTURE.md    # System design & architecture
│   ├── 02_SPECIFICATIONS_REQUIREMENTS.md    # Functional & non-functional requirements
│   ├── IG_SPECIFICATIONS_DATABASE.md        # Database schema (PostgreSQL)
│   ├── IG_SPECIFICATIONS_API.md             # Backend API (58 endpoints)
│   ├── IG_SPECIFICATIONS_FRONTEND.md        # Frontend UI/UX for farmers
│   ├── IG_SPECIFICATIONS_EMBEDDED.md        # ESP32 firmware architecture
│   ├── IG_SPECIFICATIONS_SECURITY.md        # Authentication & security
│   ├── IG_TOKKATOT_2.0_FARMER_CENTRIC_SPECIFICATIONS.md  # Farmer accessibility
│   ├── OG_SPECIFICATIONS_TECHNOLOGY_STACK.md # Technology selections
│   ├── OG_SPECIFICATIONS_DEPLOYMENT.md      # Cloud infrastructure
│   ├── OG_PROJECT_TIMELINE.md               # Development phases & milestones
│   ├── OG_TEAM_STRUCTURE.md                 # Team roles & responsibilities
│   └── OG_RISK_MANAGEMENT.md                # Risk analysis & mitigation
│
├── frontend/                    # Progressive Web App (v1.0)
│   ├── pages/                   # HTML pages
│   ├── js/                      # JavaScript modules
│   ├── css/                     # Stylesheets
│   ├── components/              # Reusable components
│   ├── assets/                  # Images, fonts, icons
│   ├── manifest.json            # PWA manifest
│   └── sw.js                    # Service worker
│
├── middleware/                  # Go backend server (v1.0)
│   ├── main.go                  # Entry point
│   ├── api/                     # API handlers
│   │   ├── authentication.go    # Auth logic
│   │   ├── profiles.go          # User profiles
│   │   ├── data-handler.go      # IoT data proxy
│   │   └── disease-detection.go # AI integration
│   ├── database/                # Database layer
│   │   └── sqlite3_db.go        # SQLite operations
│   ├── utils/                   # Utilities
│   ├── go.mod                   # Go dependencies
│   └── .env                     # Configuration
│
├── ai-service/                  # Python AI service (v1.0)
│   ├── app.py                   # Flask application
│   ├── model/                   # Trained models
│   │   ├── *.h5                 # Keras model
│   │   └── *.pkl                # Label encoder
│   └── requirements.txt         # Python dependencies
│
├── embedded/                    # ESP32 firmware (v2.0 in progress)
│   ├── main/                    # ESP-IDF main component
│   ├── components/              # Custom components (DHT, servo)
│   ├── CMakeLists.txt           # Build config
│   └── sdkconfig                # ESP-IDF settings
│
├── scripts/                     # Deployment automation (v1.0)
│   ├── deploy-all.sh            # Master deployment
│   ├── setup-access-point.sh    # WiFi AP setup
│   ├── setup-middleware-service.sh
│   ├── setup-github-runner.sh
│   ├── verify-system.sh         # System health check
│   └── README.md                # Script documentation
│
├── certs/                       # SSL certificates
├── generate-cert.sh             # Certificate generation script
├── LICENSE                      # MIT License
└── README.md                    # This file
```

---

## 📚 Documentation

### Version 2.0 Specifications (Production Release - In Development)

**Farmer-Centric Smart Poultry System** - Designed for elderly Cambodian farmers with low digital literacy

**Start here:** [docs/00_SPECIFICATIONS_INDEX.md](docs/00_SPECIFICATIONS_INDEX.md) - Complete navigation guide

#### Core Specifications (Read in Order)
| Document | Purpose | Audience |
|----------|---------|----------|
| [01_SPECIFICATIONS_ARCHITECTURE.md](docs/01_SPECIFICATIONS_ARCHITECTURE.md) | System design, 3-tier architecture, data flow | Tech Lead, Backend, DevOps |
| [02_SPECIFICATIONS_REQUIREMENTS.md](docs/02_SPECIFICATIONS_REQUIREMENTS.md) | Functional requirements, farmer-centric design | All team members |

#### Implementation Guides (IG_*)
| Document | Purpose | Audience |
|----------|---------|----------|
| [IG_SPECIFICATIONS_DATABASE.md](docs/IG_SPECIFICATIONS_DATABASE.md) | PostgreSQL schema (13 tables), simplified roles | Backend, DevOps |
| [IG_SPECIFICATIONS_API.md](docs/IG_SPECIFICATIONS_API.md) | Backend API (58 endpoints), simplified for farmers | Backend, Frontend |
| [IG_SPECIFICATIONS_FRONTEND.md](docs/IG_SPECIFICATIONS_FRONTEND.md) | UI/UX for farmers (48px+ fonts, WCAG AAA, Khmer/English) | Frontend, Design |
| [IG_SPECIFICATIONS_EMBEDDED.md](docs/IG_SPECIFICATIONS_EMBEDDED.md) | ESP32 firmware, Tokkatot team manages setup | Embedded |
| [IG_SPECIFICATIONS_SECURITY.md](docs/IG_SPECIFICATIONS_SECURITY.md) | Authentication (email/phone), simplified roles, encryption | Security, Backend |
| [IG_TOKKATOT_2.0_FARMER_CENTRIC_SPECIFICATIONS.md](docs/IG_TOKKATOT_2.0_FARMER_CENTRIC_SPECIFICATIONS.md) | Phone/Email registration, accessibility for elderly farmers | All team members |

#### Operational Guides (OG_*)
| Document | Purpose | Audience |
|----------|---------|----------|
| [OG_SPECIFICATIONS_TECHNOLOGY_STACK.md](docs/OG_SPECIFICATIONS_TECHNOLOGY_STACK.md) | Tech selections (Go, Python, PostgreSQL, DigitalOcean) | Tech Lead, all devs |
| [OG_SPECIFICATIONS_DEPLOYMENT.md](docs/OG_SPECIFICATIONS_DEPLOYMENT.md) | Cloud infrastructure, Docker, CI/CD pipelines | DevOps |
| [OG_PROJECT_TIMELINE.md](docs/OG_PROJECT_TIMELINE.md) | 10 phases, 27-35 weeks, milestones | Project Manager |
| [OG_TEAM_STRUCTURE.md](docs/OG_TEAM_STRUCTURE.md) | Team roles, responsibilities, handoff procedures | Management |
| [OG_RISK_MANAGEMENT.md](docs/OG_RISK_MANAGEMENT.md) | 10 identified risks, mitigation strategies | Project Manager, Tech Lead |
 
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
- 🔄 **Phase 1-2:** Backend API architecture (Go + PostgreSQL + InfluxDB)
- 🔄 **Phase 3:** Frontend v2 (Vue.js 3, WCAG AAA accessibility, offline support)
- 🔄 **Phase 4:** Embedded v2 (ESP32 OTA updates, MQTT, better reliability)
- 🔄 **Phase 5:** Cloud integration (DigitalOcean, Kubernetes)
- 🔄 **Phase 6-8:** Testing, deployment, rollout

**Changes from v1.0 → v2.0:**
- Cloud-connected (not just local)
- Multi-farm support (not single farm only)
- Real-time monitoring (WebSocket + MQTT)
- Remote device control
- OTA firmware updates (no farm visits)
- In-app alerts & message log (dashboard only)
- Farmer-centric UI (large text, Khmer language, accessibility)
- 99.5% uptime target
- 5-year data retention

**Timeline:** 27-35 weeks (6-8 months) - See [docs/OG_PROJECT_TIMELINE.md](docs/OG_PROJECT_TIMELINE.md)

---

## 🎨 Design & Resources

### UI Design Pages (v1.0)
Mockups and design specifications for version 1.0 interface:
- Dashboard Layout
- Device Control Interface
- Farm Settings Page
- Disease Detection UI

See `design/` folder for detailed page designs and Figma exports.

### QR Codes
Quick access codes for demo and documentation:
- **App Access QR** - Direct link to v1.0 prototype
- **Website QR** - [https://tokkatot.aztrolabe.com](https://tokkatot.aztrolabe.com)
- **Documentation QR** - Specifications & technical docs
- **Support QR** - Email & contact information
---

<div align="center">

**Proprietary Software - Tokkatot Startup**

For internal use only. Unauthorized copying, modification, or distribution is prohibited.

</div>
