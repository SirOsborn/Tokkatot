---
description: How to set up the Tokkatot development and production environments
---

# Tokkatot Environment Setup

Welcome to the Tokkatot project. Follow these unified instructions to set up your environment using Docker.

## 📦 Prerequisites
- **Docker & Docker Compose**: [Get Docker](https://docs.docker.com/get-docker/)
- **Go 1.23+**: (Required for scripts only). [Install Go](https://go.dev/doc/install)

---

## 🚀 One-Command Launch (Local Development)

1.  **Clone the Repository**:
    ```bash
    git clone https://github.com/SirOsbornOjr/tokkatot.git
    cd tokkatot
    ```

2.  **Populate the Environment**:
    -   Copy the example file: `cp .env.example .env`.
    -   Generate VAPID keys: `go run middleware/scripts/generate_vapid/main.go`.
    -   Fill in the `VAPID_PUBLIC_KEY`, `VAPID_PRIVATE_KEY`, and `JWT_SECRET`.

3.  **Launch the Stack**:
    ```bash
    docker-compose up -d --build
    ```
    -   The app will be available at `http://localhost:3000`.
    -   The database will be initialized automatically from `infra/init.sql`.

---

## ☁️ Production Deployment (AWS EC2)

For cloud deployment, follow these additional steps:

1.  **Server Provisioning**: [Ubuntu 22.04+ EC2 Instance].
2.  **SSH & CI/CD Secrets**: Add your server's `SSH_PRIVATE_KEY` and host settings to GitHub Repository Secrets.
3.  **SSL Configuration**: Follow the **[DEPLOYMENT.md](DEPLOYMENT.md)** guide to set up Certbot and Nginx with Let's Encrypt.

---

## 🛠️ Important Directories
- `/middleware`: Go backend and database logic.
- `/frontend`: Vue.js 3 PWA source.
- `/infra`: Database initialization scripts.
- `/scripts`: Deployment and maintenance utilities.

---
**Proprietary Software - Tokkatot Startup**
*For internal use only. Unauthorized distribution is prohibited.*
