# Contributing to Tokkatot

Welcome to the Tokkatot project. Join us in building a high-impact, smart poultry poultry automation platform for Cambodia.

## 📦 1. Docker-First Development
Tokkatot is a containerized application. All local development and cloud deployment is handled via `docker-compose`.

-   **Start developing**: `docker-compose up -d`.
-   **Dependencies**: Backend (Go), Frontend (Vue.js 3), Database (PostgreSQL 15), Gateway (Simulator).

## 🚀 2. CI/CD Pipeline & GitHub Secrets
The project uses automated GitHub Actions. To maintain the pipeline, ensure the following secrets are added to your repository's **Actions Secrets**:

| Secret Name | Purpose |
| :--- | :--- |
| `SERVER_HOST` | Production EC2 IP or Domain. |
| `SSH_PRIVATE_KEY` | Key to access the production server. |
| `DB_PASSWORD` | Secure password for PostgreSQL. |
| `JWT_SECRET` | Secret key for JWT auth. |
| `VAPID_PRIVATE_KEY` | Private key for push notifications. |

## 📐 3. AI Instructions & Workflows
If you are an AI assistant helping with this project, review the following:
- `AGENT.md`: Core system architecture and conventions.
- `.agents/workflows/`: Step-by-step guides for common tasks (Setup, Add-page, Testing, etc.).

---
**Proprietary Software - Tokkatot Startup**
*Designed for reliability, accessibility, and high impact in Cambodian agriculture.*
