# 🐔 Tokkatot - Smart Poultry Farm Management System

![Tokkatot Logo](frontend/assets/images/tokkatot%20logo-02.png)

**Advanced IoT-Based Poultry Disease Detection & Farm Automation**

Tokkatot is a comprehensive smart poultry management system designed for Cambodian farmers. It combines IoT sensor technology, AI-powered disease detection, and automated farm controls to improve poultry health monitoring and farm productivity.

---

## 🚀 Quick Start (Local Setup)

### Prerequisites
- **Go 1.23+**: [Install Go](https://go.dev/doc/install)
- **PostgreSQL 17+**: [Install PostgreSQL](https://www.postgresql.org/download/)
- **Python 3.12+**: [Install Python](https://www.python.org/downloads/) (for AI Service)

### Installation
1.  **Clone the Repository**:
    ```bash
    git clone https://github.com/SirOsbornOjr/tokkatot.git
    cd tokkatot
    ```
2.  **Backend Setup**:
    ```bash
    cd middleware
    cp .env.example .env  # Configure your DATABASE_URL
    go mod download
    go run main.go
    ```
3.  **Frontend**:
    The frontend is served directly by the Go backend at `http://localhost:3000`.

---

## 🛠️ Technology Stack

| Component | Technology | Description |
|---|---|---|
| **Backend** | Go 1.23 / Fiber v2 | High-performance REST API with JWT Auth |
| **Database** | PostgreSQL 17 | Unified schema for Users, Farms, and Devices |
| **Frontend** | Vue.js 3 / PWA | Mobile-first, Khmer-language, zero-build PWA |
| **AI Service** | FastAPI / PyTorch | EfficientNetB0 + DenseNet121 Ensemble |
| **IoT/Embedded** | ESP32 / ESP-IDF | C-based firmware with MQTT communication |

---

## 📂 Project Structure

- `middleware/`: Go backend source code and database migrations.
- `frontend/`: Vue.js 3 PWA (HTML/JS/CSS templates).
- `ai-service/`: Python-based disease detection service.
- `embedded/`: ESP32 firmware and hardware configuration.
- `docs/`: Technical architectural details and user guides.

---

## 📖 Documentation

- **[ARCHITECHTURE.md](docs/ARCHITECTURE.md)**: Technical deep dive into the system design.
- **[CONTRIBUTING.md](CONTRIBUTING.md)**: Development guidelines and AI Agent instructions.
- **[USER_GUIDE.md](docs/USER_GUIDE.md)** (Coming Soon): Guide for farm owners and workers.

---

## 🔐 Security & Compliance

Tokkatot uses **JWT-based authentication** and a **Registration Key system** to ensure zero SMS costs for farmers while maintaining high security. The system implements a **Unified Schema** where user profiles and farm data are tightly synchronized.

---

**Proprietary Software - Tokkatot Startup**
*For internal use only. Unauthorized copying or distribution is prohibited.*
