# Copilot Instructions — Tokkatot

> **Read [`AI_INSTRUCTIONS.md`](../AI_INSTRUCTIONS.md) at the repo root. That is the single source of truth for this project.** It covers architecture, all four components, API endpoints, database schema, and development workflow. This file is a quick-reference supplement only.

---

## Project at a Glance

**Tokkatot** — Smart poultry farm management system for Cambodian farmers  
**Stack**: Go 1.23 (Fiber v2) · Vue.js 3 (CDN, no build) · PostgreSQL · Python/FastAPI/PyTorch · ESP32/ESP-IDF  
**Architecture**: Farm → Coop → Devices (hierarchical; all devices belong to coops, never farms)  
**Database**: PostgreSQL only — SQLite was permanently removed Feb 2026  
**API**: 67 REST endpoints under `/v1/` + WebSocket at `/v1/ws`  
**Roles**: `farmer` (full farm control; multiple farmers may share a farm), `viewer` (read-only + acknowledge alerts only)  
**Tokkatot admin**: system-level only (reg keys, JWT secrets) — NOT a `farm_users` role  
**Frontend**: Static files served by Go middleware from `frontend/` directory  
**AI Model**: Proprietary — `ai-service/outputs/*.pth` files are gitignored, must exist locally

---

## Critical Rules

1. **Never commit** `middleware/.env`, `ai-service/.env`, or `outputs/*.pth`
2. **PostgreSQL only** — `database/postgres.go` is the only DB file; there is no SQLite
3. **All API routes** use `/v1/` prefix (not `/api/`)
4. **Always call** `checkFarmAccess()` before any device control
5. **Frontend files** that do not exist yet: `disease-detection.js`, `styleSchedules.css` — they need to be built
6. **Model files** must exist before running `python app.py` or `docker build` in `ai-service/`

---

## Quick Build Commands

```powershell
# Go middleware (Windows)
cd middleware
go build -o backend.exe && .\backend.exe   # starts on http://localhost:3000

# Full API test suite
.\test_all_endpoints.ps1             # email login
.\test_all_endpoints.ps1 -UsePhone   # phone login

# AI service (Python)
cd ai-service
python -m venv env && env\Scripts\activate
pip install -r requirements.txt
python app.py    # http://localhost:8000 (needs outputs/ensemble_model.pth)

# AI service (Docker — recommended)
cd ai-service && docker-compose up -d tokkatot-ai
```

---

## What Exists (Frontend Pages)

| Route | File | Status |
|-------|------|--------|
| `/` | `pages/index.html` | ✅ |
| `/login` | `pages/login.html` | ✅ |
| `/register` | `pages/signup.html` | ✅ |
| `/profile` | `pages/profile.html` | ✅ |
| `/settings` | `pages/settings.html` | ✅ |
| `/disease-detection` | `pages/disease-detection.html` | ✅ (Coming Soon overlay) |
| `/monitoring` | `pages/monitoring.html` | ✅ Temperature timeline |
| `/schedules` | `pages/schedules.html` | ✅ Live |

---

## Current Version (v2.3 — Feb 24, 2026)

- **Temperature timeline** (`/monitoring`) — `TemperatureTimelineHandler` in `api/coop_handler.go`
  - Route: `GET /v1/farms/:farm_id/coops/:coop_id/temperature-timeline?days=7`
  - `bg_hint`: `scorching`≥35°C · `hot`≥32°C · `warm`≥28°C · `neutral`≥24°C · `cool`≥20°C · `cold`<20°C
- **Schedule automation** — `action_duration` (auto-off timer) + `action_sequence` (pulse patterns)
- **Alerts & subscriptions** — `api/alert_handler.go` (8 endpoints)
- **Analytics** — `api/analytics_handler.go` (6 endpoints inc. CSV export)
- **Farm members** — invite, role-change, remove via `api/farm_handler.go`

---

## Documentation Map

| Need | File |
|------|------|
| Full project context | [`AI_INSTRUCTIONS.md`](../AI_INSTRUCTIONS.md) |
| API spec (67 endpoints) | [`docs/implementation/API.md`](../docs/implementation/API.md) |
| Database schema (14 tables) | [`docs/implementation/DATABASE.md`](../docs/implementation/DATABASE.md) |
| Frontend spec | [`docs/implementation/FRONTEND.md`](../docs/implementation/FRONTEND.md) |
| Farmer automation scenarios | [`docs/AUTOMATION_USE_CASES.md`](../docs/AUTOMATION_USE_CASES.md) |
| Security & auth | [`docs/implementation/SECURITY.md`](../docs/implementation/SECURITY.md) |

