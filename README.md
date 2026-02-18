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
- **In Development:** v2.0 (Production) - See [docs/SPECIFICATIONS_INDEX.md](docs/SPECIFICATIONS_INDEX.md)

---

### Option 1: One-Command Setup (Recommended)

```bash
# Clone repository
git clone https://github.com/SirOsborn/Tokkatot.git
cd Tokkatot

# Run complete deployment
sudo bash scripts/deploy-all.sh
```

### Option 2: Manual Setup

```bash
# 1. Setup WiFi Access Point
sudo bash scripts/setup-access-point.sh

# 2. Build and install middleware
sudo bash scripts/setup-middleware-service.sh

# 3. (Optional) Setup GitHub Actions runner
sudo bash scripts/setup-github-runner.sh

# 4. Reboot
sudo reboot
```

### Access the System

1. **Connect to WiFi**
   - SSID: `Smart Poultry 1.0.0-0001`
   - Password: `skibiditoilet168`

2. **Open Browser**
   - Go to: `http://10.0.0.1:4000`

3. **Create Account**
   - Sign up and start managing your farm!

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
│   ├── SPECIFICATIONS_INDEX.md  # Start here - navigation hub
│   ├── SPECIFICATIONS_ARCHITECTURE.md
│   ├── SPECIFICATIONS_REQUIREMENTS.md
│   ├── SPECIFICATIONS_DATABASE.md
│   ├── SPECIFICATIONS_TECHNOLOGY_STACK.md
│   ├── SPECIFICATIONS_FRONTEND.md
│   ├── SPECIFICATIONS_EMBEDDED.md
│   ├── SPECIFICATIONS_DEPLOYMENT.md
│   ├── SPECIFICATIONS_SECURITY.md
│   ├── PROJECT_TIMELINE.md
│   ├── TEAM_STRUCTURE.md
│   ├── RISK_MANAGEMENT.md
│   └── TOKKATOT_2.0_FARMER_CENTRIC_ADDITIONS.md
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

**Start here:** [docs/SPECIFICATIONS_INDEX.md](docs/SPECIFICATIONS_INDEX.md) - Complete navigation guide

| Document | Purpose | Audience |
|----------|---------|----------|
| [SPECIFICATIONS_ARCHITECTURE.md](docs/SPECIFICATIONS_ARCHITECTURE.md) | System design & data flow | Tech Lead, Backend, DevOps |
| [SPECIFICATIONS_REQUIREMENTS.md](docs/SPECIFICATIONS_REQUIREMENTS.md) | Functional requirements (FR1-FR8) | All team members |
| [SPECIFICATIONS_DATABASE.md](docs/SPECIFICATIONS_DATABASE.md) | Database schema (13 PostgreSQL tables) | Backend, DevOps |
| [SPECIFICATIONS_TECHNOLOGY_STACK.md](docs/SPECIFICATIONS_TECHNOLOGY_STACK.md) | Tech selection & justification | Tech Lead, all devs |
| [SPECIFICATIONS_FRONTEND.md](docs/SPECIFICATIONS_FRONTEND.md) | UI/UX & accessibility (WCAG AAA) | Frontend, Design |
| [SPECIFICATIONS_EMBEDDED.md](docs/SPECIFICATIONS_EMBEDDED.md) | ESP32 firmware architecture | Embedded |
| [SPECIFICATIONS_DEPLOYMENT.md](docs/SPECIFICATIONS_DEPLOYMENT.md) | DigitalOcean, Kubernetes, CI/CD | DevOps |
| [SPECIFICATIONS_SECURITY.md](docs/SPECIFICATIONS_SECURITY.md) | JWT auth, RBAC, encryption, compliance | Security, DevOps |
| [PROJECT_TIMELINE.md](docs/PROJECT_TIMELINE.md) | 10 phases, 27-35 weeks, milestones | Project Manager |
| [TEAM_STRUCTURE.md](docs/TEAM_STRUCTURE.md) | Team roles & responsibilities | Management |
| [RISK_MANAGEMENT.md](docs/RISK_MANAGEMENT.md) | 10 identified risks & mitigation | Project Manager, Tech Lead |
| [TOKKATOT_2.0_FARMER_CENTRIC_ADDITIONS.md](docs/TOKKATOT_2.0_FARMER_CENTRIC_ADDITIONS.md) | Farmer-first design details | All team members |

### Version 1.0 Guides (Current Prototype)

| Document | Description |
|----------|-------------|
| [QUICKSTART.md](QUICKSTART.md) | 5-minute setup guide |
| [DEPLOYMENT.md](DEPLOYMENT.md) | Comprehensive deployment guide |
| [VIRTUALBOX_SETUP.md](VIRTUALBOX_SETUP.md) | Testing with VirtualBox |
| [scripts/README.md](scripts/README.md) | Script documentation |

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

## 🛠️ Development Setup

### Prerequisites

```bash
# Install Go
wget https://go.dev/dl/go1.23.6.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.23.6.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Install Python
sudo apt install python3 python3-pip python3-venv

# Install Git
sudo apt install git
```

### Clone and Build

```bash
# Clone repository
git clone https://github.com/SirOsborn/Tokkatot.git
cd Tokkatot

# Setup middleware
cd middleware
go mod download
go build -o middleware main.go

# Setup AI service
cd ../ai-service
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt

# Run development server
cd ../middleware
./middleware
```

---

## 🧪 Testing

### System Verification

```bash
# Run complete system check
sudo bash scripts/verify-system.sh
```

### Manual Testing

```bash
# Test middleware
curl http://localhost:4000

# Test AI service
curl http://localhost:5000/health

# Test with VirtualBox
# See VIRTUALBOX_SETUP.md for detailed instructions
```

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

## 🎓 Usage Examples

### Disease Detection (Python)

```python
import requests

url = "http://10.0.0.1:4000/api/ai/predict-disease"
files = {"image": open("chicken_sample.jpg", "rb")}
headers = {"Authorization": "Bearer YOUR_JWT_TOKEN"}

response = requests.post(url, files=files, headers=headers)
result = response.json()

print(f"Disease: {result['prediction']['predicted_disease']}")
print(f"Confidence: {result['prediction']['confidence']}")
```

### Get Sensor Data (JavaScript)

```javascript
fetch('http://10.0.0.1:4000/api/get-current-data', {
    headers: {
        'Authorization': 'Bearer YOUR_JWT_TOKEN'
    }
})
.then(res => res.json())
.then(data => {
    console.log('Temperature:', data.temperature);
    console.log('Humidity:', data.humidity);
});
```

### Toggle Device (cURL)

```bash
# Toggle fan
curl -X GET http://10.0.0.1:4000/api/toggle-fan \
  -H "Cookie: token=YOUR_JWT_TOKEN"
```

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

**Timeline:** 27-35 weeks (6-8 months) - See [docs/PROJECT_TIMELINE.md](docs/PROJECT_TIMELINE.md)
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

**Timeline:** 27-35 weeks (6-8 months) - See [docs/PROJECT_TIMELINE.md](docs/PROJECT_TIMELINE.md)

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

See `qr-codes/` folder for printable assets.

---

<div align="center">

**Proprietary Software - Tokkatot Startup**

For internal use only. Unauthorized copying, modification, or distribution is prohibited.

</div>
