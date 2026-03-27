# 🐔 Tokkatot - Smart Poultry Farm Management System

![Tokkatot Logo](frontend/assets/images/tokkatot%20logo-02.png)

**Digitalizing Cambodian Poultry Farming with IoT Automation & Cloud Intelligence**

Tokkatot is a production-ready smart poultry management platform. It combines IoT sensor technology, coop-level automation, and a robust cloud backend to improve farm productivity while ensuring 24/7 poultry health monitoring.

---

## 🚀 Quick Start (Production Setup)

The fastest way to launch Tokkatot is using **Docker Compose**.

### 1. Prerequisites
- **Docker & Docker Compose**: [Install Docker](https://docs.docker.com/get-docker/)
- **Go 1.23+**: (Only required if running outside Docker)

### 2. Configuration
```bash
# Clone the repository
git clone https://github.com/SirOsbornOjr/tokkatot.git
cd tokkatot

# Populate your .env (Use the provided generator for VAPID keys)
cp .env.example .env
go run middleware/scripts/generate_vapid/main.go
```

### 3. Launch
```bash
docker-compose up -d
```
The app will be available at `https://app.tokkatot.com` (if SSL is configured) or `http://localhost:3000` (locally).

---

## 🛠️ Technology Stack

| Layer | Technology | Description |
|---|---|---|
| **Infrastructure** | Docker / Nginx / Certbot | Containerized stack with SSL automation |
| **Backend** | Go 1.23 / Fiber v2 | High-performance REST API with JWT Auth |
| **CI/CD** | GitHub Actions / GHCR | Automated build/push/deploy pipeline |
| **Database** | PostgreSQL 15 | Persistent storage with named volumes |
| **Frontend** | Vue.js 3 / PWA | Mobile-first, Khmer-language, zero-build PWA |
| **IoT/Embedded** | ESP32-IDF / RPi 4B | Local HTTPS control & Telemetry ingestion |

---

## 🔄 CI/CD Pipeline

Tokkatot features a professional 3-tier deployment pipeline via **GitHub Actions**:

1.  **Dev**: Pushes to `dev` branch build and deploy to the development server.
2.  **Stage**: Pushes to `stage` are for final cloud validation.
3.  **Prod**: Pushes to `main` trigger a production rebuild and deployment to `app.tokkatot.com`.

---

## 🔐 Security & Compliance

- **Pre-render Auth**: Blocking scripts prevent unauthorized UI access.
- **Docker Lockdown**: Middleware is isolated from the public internet; all traffic must pass through the Nginx gateway.
- **Registration Keys**: Enforced signup flow requiring physical keys or admin invites.
- **Health Monitoring**: Standardized `/v1/health` probes for cloud uptime.

---

## 📖 Documentation

- **[DEPLOYMENT.md](DEPLOYMENT.md)**: Step-by-step guide for AWS EC2 & SSL setup.
- **[ARCHITECTURE.md](docs/ARCHITECTURE.md)**: Technical system design and API surface.
- **[CONTRIBUTING.md](CONTRIBUTING.md)**: Developer guidelines and AI Agent instructions.

---
**Proprietary Software - Tokkatot Startup**
*Designed for reliability, accessibility, and high impact in Cambodian agriculture.*
