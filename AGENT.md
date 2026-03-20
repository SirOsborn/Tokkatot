# 🤖 AGENT.md — Tokkatot AI Master Guide

This file is the primary context and instruction set for AI agents (Claude, GPT, Gemini) working on the Tokkatot project. It defines the system's purpose, architecture, and the technical "hooks" required for effective autonomous coding.

---

## 🧐 What is Tokkatot?
Tokkatot is a **smart poultry farm management system** designed for Cambodian farmers. It integrates IoT sensor hardware, automated coop controls, and (later) AI disease detection to improve farm productivity.

## 🎯 Why Tokkatot? (The "Why")
- **Accessibility**: UI must be high-contrast, large-font, and support Khmer natively.
- **Maintainability**: One-page, one-stylesheet philosophy. Avoid long, monolithic files.
- **Disease Prevention**: AI detection is planned for a later patch.
- **Zero Cost**: Registration Key system ensures no SMS/Email verification costs for users.

## 🏗️ How it Works (The "How")
- **Frontend**: Vue.js 3 (CDN-based PWA) with a **Modular CSS Architecture**.
- **Middleware**: Go 1.23 + Fiber v2. REST API + WebSocket + PostgreSQL.
- **AI Service**: Python FastAPI + PyTorch Ensemble (not integrated yet).
- **Embedded**: ESP32 (ESP-IDF) controlled by a Raspberry Pi 4B gateway.

## 📊 Data Hierarchy
1.  **User**: Farmer or worker (role-based access).
2.  **Farm**: Physical location owned by a User.
3.  **Coop**: Primary control unit. All automation is **coop-level**.
4.  **Device**: Sensors and actuators assigned to a coop. Missing devices are allowed and marked inactive.

## 🌐 API & Communication (Cloud)
- **Auth**: `/v1/auth/signup`, `/v1/auth/login`, `/v1/auth/refresh`, `/v1/auth/logout` (no email/SMS verification).
- **Farms**: `/v1/farms`, `/v1/farms/:id/members`.
- **Coops**: `/v1/farms/:farm_id/coops`, `/v1/farms/:farm_id/coops/:coop_id`.
- **Devices**: `/v1/farms/:id/devices`, `/v1/farms/:id/devices/:id/commands`.
- **Schedules**: `/v1/farms/:id/schedules` (coop-level).
- **Telemetry**: `/v1/farms/:farm_id/coops/:coop_id/telemetry` (gateway → cloud).
- **Device Report**: `/v1/farms/:farm_id/coops/:coop_id/devices/report` (gateway → cloud).
- **Monitoring Timeline**: `/v1/farms/:farm_id/coops/:coop_id/temperature-timeline`.

## 📡 Embedded & IoT
- **Platform**: ESP32 (ESP-IDF) + Raspberry Pi 4B gateway.
- **Protocol**: ESP32 exposes local HTTPS endpoints; Pi polls sensors and executes commands.
- **Logic**: ON/OFF relays + schedule sequences executed by gateway.
- **Telemetry**: Pi posts temperature/humidity/water level to cloud.
- **Water System**: floating valve only (no pump). Water sensor is monitor-only.
- **Feeding System**: high-torque feeder motor (relay on/off).
- **Water Alert Rule**: Water below half threshold for 1 minute triggers alert.

---

## 🔗 Agent Skills & Hooks
Use these "hooks" as entry points for your research and tasks:

### 🗄️ Database Hook: `middleware/database/schema.go`
- **Schema**: Single source of truth for the PostgreSQL DDL.

### 🎨 Frontend Styling Hook: `frontend/css/`
- **Theme**: `theme.css` (Design tokens: colors, spacing, typography).
- **UI Elements**: `components.css` (Standardized buttons, cards, badges, inputs).
- **Layout**: `layout.css` (Shared headers, navbars, page containers).
- **Page Styles**: Modular CSS files (e.g., `dashboard.css`, `monitoring.css`, `auth.css`).

### 📱 Frontend Logic Hook: `frontend/js/utils/api.js`
- **API Client**: Standard fetch wrapper for the Vue app.
- **I18N**: `frontend/js/utils/i18n.js` for Khmer translations.

### 🌐 API Hook: `middleware/main.go`
- **Routing**: All REST and WebSocket routes are registered here.

### 🛠️ Script Hook: `middleware/scripts/`
- **Utilities**: `migrate_fresh.go` (schema reset) and other diagnostics.

### 📋 Workflow Hook: `.agents/workflows/`
- **Procedures**: Guides for `setup`, `test`, `add-api`, and `add-page`.

---

## 🛠️ Agent Implementation Rules
- **Modular Frontend**: When adding a new page, create a dedicated `.css` file in `frontend/css/` and link it after `theme.css` and `components.css`.
- **No Monoliths**: Keep Go handlers and JS logic in small, focused files.
- **Khmer First**: Always provide Khmer translations for UI text.
- **Type Safety**: Use `GetUserIDFromContext(c)` for UUIDs in Go.
- **Environment Driven**: Use `.env` for all secrets and configuration.

---

**Proprietary Software - Tokkatot Startup**
*Be proactive, verify your changes, and prioritize the farmer's experience.*
