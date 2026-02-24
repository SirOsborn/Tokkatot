# Frontend Implementation

**Last Updated**: February 24, 2026  
**Status**: ✅ Fully Rebuilt — Vue.js 3 CDN, new design system  
**Tech Stack**: Vue.js 3 CDN (no build step), Mi Sans, Google Material Symbols, global CSS design system

---

## Overview

All pages are static HTML served by the Go Fiber backend. No npm, no bundler — Vue 3 is loaded via CDN. Every page shares a single global design system (`/css/style.css`).

### Served Routes (`middleware/main.go`)

| Route | File | Auth | Notes |
|-------|------|------|-------|
| `/` | `pages/index.html` | ✅ Required | Dashboard, live WebSocket data |
| `/login` | `pages/login.html` | Public | Redirects to `/` if already authed |
| `/register` | `pages/signup.html` | Public | Redirects to `/` if already authed |
| `/profile` | `pages/profile.html` | ✅ Required | View + edit user info |
| `/settings` | `pages/settings.html` | ✅ Required | Language, notifications, logout |
| `/schedules` | `pages/schedules.html` | ✅ Required | Full schedule CRUD + sequence builder |
| `/disease-detection` | `pages/disease-detection.html` | ✅ Required | Coming Soon overlay |
| `/monitoring` | `pages/monitoring.html` | ✅ Required | Apple Weather-style temp timeline |
| 404 fallback | `pages/404.html` | Public | — |

---

## File Structure

```
frontend/
├── css/
│   └── style.css               ← Single global design system (ALL pages use this)
├── js/
│   └── utils/
│       ├── api.js              ← window.API — all HTTP calls + JWT auth
│       ├── i18n.js             ← window.i18n / window.t() — KM/EN translations
│       └── components.js       ← window.loadComponents() — injects header + navbar
├── components/
│   ├── header.html             ← Top bar: logo, farm name, lang toggle, alert bell
│   └── navbar.html             ← Bottom nav: 5 items with Material Symbols icons
├── pages/
│   ├── index.html              ← Dashboard (Vue 3)
│   ├── login.html              ← Login (Vue 3, no header/nav)
│   ├── signup.html             ← Register (Vue 3, no header/nav)
│   ├── profile.html            ← Profile view/edit (Vue 3)
│   ├── settings.html           ← Settings + logout (Vue 3)
│   ├── schedules.html          ← Schedule CRUD + sequence builder (Vue 3)
│   ├── disease-detection.html  ← Coming Soon overlay + upload UI
│   ├── monitoring.html         ← Temperature timeline (vanilla JS)
│   └── 404.html                ← 404 page (no header/nav)
└── assets/
    ├── fonts/
    ├── icons/
    └── images/
```

---

## Design System (`/css/style.css`)

### Color Tokens

```css
--color-primary:        #20a39e;   /* Teal — primary actions */
--color-secondary:      #ffba49;   /* Amber — highlights */
--color-danger:         #ef5b5b;   /* Red — errors, delete */
--color-primary-light:  #e6f7f6;   /* Teal tint — hover states */
--color-bg:             #f4f6f8;   /* Page background */
--color-surface:        #ffffff;   /* Card/input background */
--color-text:           #1a1a1a;   /* Primary text */
--color-text-secondary: #6b7280;   /* Muted text */
--color-border:         #e5e7eb;   /* Borders */
```

### Typography

- **Font**: Mi Sans (jsDelivr CDN) + Noto Sans Khmer fallback
- **Base size**: `--font-size-base: 18px` — large for farmer accessibility
- **Scale**: `--font-size-sm: 14px` / `--font-size-lg: 20px` / `--font-size-xl: 24px`

### Spacing & Touch Targets

```css
--touch-min:     48px;   /* Minimum tap target */
--space-sm:      8px;
--space-md:      16px;
--space-lg:      24px;
--space-xl:      32px;
--header-height: 60px;
--nav-height:    68px;
```

### Body Padding

Pages with shared header + navbar get padding automatically:

```css
body {
  padding-top: var(--header-height);
  padding-bottom: calc(var(--nav-height) + env(safe-area-inset-bottom));
}
```

Pages without nav/header (login, signup, 404) use `body.no-nav.no-header`.

---

## Utility Scripts

### `window.API` (`/js/utils/api.js`)

Central HTTP client. All pages use this instead of raw `fetch`.

```js
// HTTP methods — all return parsed JSON
await API.get('/v1/farms')
await API.post('/v1/auth/login', { email, password })
await API.put('/v1/farms/1/schedules/5', payload)
await API.patch('/v1/auth/me', { name })
await API.delete('/v1/farms/1/schedules/5')
await API.upload('/v1/ai/predict-disease', formData)

// Auth helpers
API.isAuthenticated()        // bool
API.requireAuth()            // redirects to /login if not authed; returns false
API.logout()                 // clears storage, redirects to /login

// Farm/coop context
API.getSelectedFarmId()      // localStorage 'selected_farm_id'
API.setSelectedFarmId(id)
API.getSelectedCoopId()
API.setSelectedCoopId(id)
```

**Auto-refresh**: On 401, automatically calls `POST /v1/auth/refresh` once. If that fails, redirects to `/login`.

**localStorage keys**:
- `access_token` — JWT access token
- `refresh_token`
- `user_name`
- `selected_farm_id`
- `selected_coop_id`

### `window.i18n` + `window.t()` (`/js/utils/i18n.js`)

```js
t('nav_home')       // → "ផ្ទះ" (KM) or "Home" (EN)
i18n.setLang('en')
i18n.toggleLang()   // switches KM ↔ EN, saves to localStorage
i18n.getLang()      // → 'km' | 'en'
i18n.applyAll()     // applies all data-i18n attributes on DOM
```

Default language: **Khmer** (`km`). Stored in `localStorage('language')`.

### `window.loadComponents()` (`/js/utils/components.js`)

Fetches `header.html` and `navbar.html`, injects into `#header-placeholder` / `#navbar-placeholder`.

Uses `document.createRange().createContextualFragment()` so `<script>` tags inside components execute correctly (unlike `innerHTML`).

```js
await loadComponents();   // call in every page's mounted() or init()
```

---

## Page Template

Every authenticated page follows this pattern:

```html
<!DOCTYPE html>
<html lang="km">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0, viewport-fit=cover" />
  <title>Page Title – Tokkatot</title>
  <link rel="icon" href="/assets/images/tokkatot logo-02.png" type="image/x-icon" />
  <link rel="stylesheet" href="/css/style.css" />
  <link rel="stylesheet" href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:opsz,wght,FILL,GRAD@20..48,100..700,0..1,-50..200" />
</head>
<body>

<div id="header-placeholder"></div>

<div id="app">
  <!-- Vue 3 content -->
</div>

<div id="navbar-placeholder"></div>

<script src="/js/utils/i18n.js"></script>
<script src="/js/utils/api.js"></script>
<script src="/js/utils/components.js"></script>
<script src="https://unpkg.com/vue@3/dist/vue.global.prod.js"></script>
<script>
const { createApp } = Vue;
createApp({
  data() { return { /* ... */ }; },
  async mounted() {
    if (!window.API.requireAuth()) return;
    await loadComponents();
    // load data...
  },
  methods: { /* ... */ }
}).mount('#app');
</script>
</body>
</html>
```

Public pages (login, signup, 404) — omit placeholders, add `body.no-nav.no-header`.

---

## Page Reference

### `index.html` — Dashboard

**API**:
- `GET /v1/farms` → farm selector
- `GET /v1/farms/:id/coops` → coop pills
- `GET /v1/devices?coop_id=` → device list
- `POST /v1/devices/:id/command` — `{ command_type: 'turn_on' | 'turn_off' }`
- WebSocket: `ws(s)://host/v1/ws?token=&farm_id=&coop_id=` → reactive temp/humidity

**UI**: Farm banner, coop selector pills, 2-col metrics grid, devices list with toggles, quick action grid.

---

### `login.html`

**API**: `POST /v1/auth/login` → stores `access_token`, `refresh_token`, `user_name` → redirects to `/`

---

### `signup.html`

**Browser URL**: `/register`  
**API**: `POST /v1/auth/signup`

```json
{ "name": "...", "email": "...", "password": "...", "registration_key": "...", "phone": "..." }
```

`phone` is optional. On success → shows message → redirects to `/login` after 1.8s.

---

### `profile.html`

**API**:
- `GET /v1/auth/me` → load
- `PATCH /v1/auth/me` → save name + phone

View mode: avatar initials, role badge, read-only fields. Edit mode: name + phone inputs.

---

### `settings.html`

- Language toggle (KM ↔ EN)
- Notifications toggle
- Logout with `v-if` confirm modal (no `window.confirm`)

---

### `schedules.html`

**API**:
- `GET /v1/farms/:id/schedules`
- `POST /v1/farms/:id/schedules`
- `PUT /v1/farms/:id/schedules/:sid`
- `DELETE /v1/farms/:id/schedules/:sid`
- `POST /v1/farms/:id/schedules/:sid/execute-now`

**Features**:
- Schedule cards with enabled toggle, Run Now, Edit, Delete
- Bottom-sheet modal form (Vue `v-if`) for create/edit
- **Multi-step sequence builder**: toggleable, add/remove ON/OFF steps with per-step duration (seconds)
- **Auto-off duration**: `action_duration` field (0 = disabled)
- FAB `+` button
- Cron expression field with human-readable hint

**Sequence payload** (stored in `action_sequence` as JSON string):
```json
[{"action":"ON","duration":30},{"action":"OFF","duration":10}]
```

---

### `disease-detection.html`

**Status**: Coming Soon overlay active.

To enable: remove `<div id="coming-soon-overlay">...</div>` (marked with comments in file).

**Full UI underneath**:
- Upload zone (camera capture + file picker + drag & drop)
- Image preview
- Confidence bar, healthy/disease badge, recommendation text
- `POST /v1/ai/predict-disease` with `FormData` (field: `image`)

---

### `monitoring.html`

Apple Weather-style. Uses `window.API` (vanilla JS, no Vue).

**API**: `GET /v1/farms/:fid/coops/:cid/temperature-timeline?days=7`

**UI**:
1. Coop picker `<select>`
2. Hero — large temp, H/L peak times, bg_hint badge
3. Hourly scroll strip
4. SVG bezier curve graph with gradient fill + H/L markers
5. Daily history with proportional range bars

**Dynamic background** (`bg_hint` from API):

| Class | Condition | Gradient |
|---|---|---|
| `bg-scorching` | ≥ 35°C | `#7f0000 → #e85d04` |
| `bg-hot` | ≥ 32°C | `#c1121f → #f48c06` |
| `bg-warm` | ≥ 28°C | `#e85d04 → #faa307` |
| `bg-neutral` | ≥ 24°C | `#2d6a4f → #74c69d` |
| `bg-cool` | ≥ 20°C | `#023e8a → #0096c7` |
| `bg-cold` | < 20°C | `#03045e → #0077b6` |

---

### `404.html`

No header/navbar. Khmer "រកទំព័រមិនឃើញ", teal 404, home button to `/`.

---

## Components

### `header.html`

```
[Logo + Farm Name]              [🌐 lang] [🔔 alerts] [👤]
```

- `#header-farm-name` — set by `loadComponents()` from `GET /v1/farms`
- `#lang-toggle-btn` → calls `window.headerToggleLang()`
- Alert bell with `#header-alert-badge`
- Profile icon links to `/profile`

### `navbar.html`

| Icon (Material) | Key | Route |
|---|---|---|
| `home` | `nav_home` | `/` |
| `monitor_heart` | `nav_monitoring` | `/monitoring` |
| `biotech` | `nav_disease` | `/disease-detection` |
| `calendar_clock` | `nav_schedules` | `/schedules` |
| `settings` | `nav_settings` | `/settings` |

Active item highlighted via `data-nav` attribute matched by `components.js`.

---

## i18n Key Reference

| Key | Khmer | English |
|-----|-------|---------|
| `nav_home` | ផ្ទះ | Home |
| `nav_monitoring` | ត្រួតពិនិត្យ | Monitoring |
| `nav_disease` | ជំងឺ | Disease |
| `nav_schedules` | កាលវិភាគ | Schedules |
| `nav_settings` | ការកំណត់ | Settings |
| `login` | ចូលគណនី | Login |
| `logout` | ចេញ | Logout |
| `save` | រក្សាទុក | Save |
| `cancel` | បោះបង់ | Cancel |
| `delete` | លុប | Delete |
| `enabled` | បើកដំណើរការ | Enabled |
| `disabled` | បិទ | Disabled |
| `temperature` | សីតុណ្ហភាព | Temperature |
| `humidity` | សំណើម | Humidity |
| `devices` | ឧបករណ៍ | Devices |
| `schedules` | កាលវិភាគ | Schedules |
| `add_schedule` | បន្ថែមកាលវិភាគ | Add Schedule |
| `edit_schedule` | កែកាលវិភាគ | Edit Schedule |
| `no_data` | មិនមានទិន្នន័យ | No data |
| `error` | កំហុស | Error |
