# 🤖 AGENT.md — Tokkatot AI Master Guide

This file is the primary context and instruction set for AI agents (Claude, GPT, Gemini) working on the Tokkatot project. It defines the system's purpose, architecture, and the technical "hooks" required for effective autonomous coding.

---

## 🧐 What is Tokkatot?
Tokkatot is an **AI-powered Smart Poultry Farm Management System** specifically designed for Cambodian farmers. It integrates IoT sensor hardware with deep learning disease detection to improve farm productivity.

## 🎯 Why Tokkatot? (The "Why")
- **Accessibility**: UI must be high-contrast, large-font, and support Khmer natively.
- **Maintainability**: One-page, one-stylesheet philosophy. Avoid long, monolithic files.
- **Disease Prevention**: Early detection from manure photos prevents farm-wide outbreaks.
- **Zero Cost**: Registration Key system ensures no SMS/Email verification costs for users.

## 🏗️ How it Works (The "How")
- **Frontend**: Vue.js 3 (CDN-based PWA) with a **Modular CSS Architecture**.
- **Middleware**: Go 1.23 + Fiber v2. REST API + WebSocket + PostgreSQL.
- **AI Service**: Python FastAPI + PyTorch Ensemble.
- **Embedded**: ESP32 (C/ESP-IDF) + MQTT local hub.

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
