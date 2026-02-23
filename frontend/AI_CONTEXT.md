# 🤖 AI Context: Vue.js 3 Frontend

**Component**: `frontend/` - Mobile-first PWA for farmers  
**Tech Stack**: Vue.js 3 (CDN-based, no build step), HTML5, CSS3, vanilla JavaScript  
**Purpose**: Accessible device control, disease detection, schedule management for elderly Cambodian farmers  

---

## 📖 Read First

**Before reading this file**, understand the project context:
- **Project overview**: Read [`../AI_INSTRUCTIONS.md`](../AI_INSTRUCTIONS.md) for business model, farmer-centric design, Farm→Coop→Device hierarchy
- **Full frontend spec**: See [`../docs/implementation/FRONTEND.md`](../docs/implementation/FRONTEND.md) for complete UI/UX, accessibility, and page specifications

**This file contains**: Vue.js-specific patterns, common API calls, schedule UI examples, development tasks

---

## 📚 Full Documentation

| Document | Purpose |
|----------|---------|
| [`docs/implementation/FRONTEND.md`](../docs/implementation/FRONTEND.md) | Complete UI/UX specs, page wireframes, accessibility standards |
| [`docs/implementation/API.md`](../docs/implementation/API.md) | All 35 API endpoints this frontend calls |
| [`docs/AUTOMATION_USE_CASES.md`](../docs/AUTOMATION_USE_CASES.md) | Real farmer scenarios for schedule UI (pulse feeding, climate control) |
| [`docs/implementation/SECURITY.md`](../docs/implementation/SECURITY.md) | JWT auth flow, token storage, WebSocket security |
| [`AI_INSTRUCTIONS.md`](../AI_INSTRUCTIONS.md) | Why farmers need simple UI (3-role system, accessibility-first design) |

---

## 📁 Quick File Reference

```
frontend/
├── pages/
│   ├── index.html               # Dashboard (primary page, device status grid)
│   ├── disease-detection.html   # AI disease detection + history
│   ├── schedules.html           # Automation schedules (action_sequence UI)
│   ├── login.html               # Login/signup
│   ├── profile.html             # User profile, sessions
│   └── settings.html            # Settings (language toggle, notifications)
├── components/
│   ├── header.html              # Shared header (farm dropdown, alerts)
│   └── navbar.html              # Bottom tab navigation
├── js/
│   ├── index.js                 # Home page logic (device grid, WebSocket updates)
│   ├── disease-detection.js     # Image upload → AI prediction → display
│   ├── schedules.js             # Create/edit schedules with action_sequence
│   ├── login.js                 # JWT auth + token storage
│   ├── header.js & navbar.js    # Component logic
│   └── utils/api.js             # Helper: fetch with auth, error handling
├── css/
│   ├── styleHome.css            # Dashboard: device cards, sensor gauges
│   ├── styleSchedules.css       # Schedule UI: time pickers, action builders
│   ├── loginSignUp.css          # Auth pages
│   └── (one CSS per page module)
└── assets/
    ├── images/, icons/, fonts/  # Graphics, SVG icons, Khmer fonts
```

---

## 🎯 Farmer Accessibility (CRITICAL)

**Target Users**: Elderly Cambodian farmers, 60+ years old, low digital literacy, 1-2GB RAM phones, 4G network

### UI Constraints:
1. **Buttons**: Minimum 48px height/width (touch-friendly for older hands)
2. **Fonts**: Minimum 16px body text, 24px+ headings, 48px+ important numbers
3. **Colors**: WCAG AAA contrast (black text on white background, not gray)
4. **Language**: Khmer & English toggle (easy to find, top-right header)
5. **Icons**: Always pair with text labels (don't rely on icons alone)
6. **Actions**: Maximum 3 main actions per screen (avoid complexity)

### No Build System:
- Vue.js 3 loaded via CDN (`https://unpkg.com/vue@3/dist/vue.global.js`)
- No npm, no webpack, no transpilation
- Files served directly by Go middleware (static file server)
- Opening `index.html` in browser = instant preview

**Why**: Farmer devices can't handle complex SPA bundles, need fast load on 4G.

---

## 🔧 Common Vue.js Patterns

### Page Initialization (Typical `js/index.js` structure)

```javascript
const { createApp } = Vue;

createApp({
  data() {
    return {
      farm: null,            // Current farm object
      devices: [],           // Array of device objects
      sensors: [],           // Array of sensor readings
      ws: null,              // WebSocket connection
      loading: true,
      error: null
    };
  },
  
  async mounted() {
    await this.validateJWT();      // Check token validity
    await this.loadFarmData();     // GET /api/farms/{id}
    await this.loadDevices();      // GET /api/farms/{id}/devices
    this.connectWebSocket();       // ws://server/ws
  },
  
  methods: {
    async validateJWT() {
      const token = localStorage.getItem('access_token');
      if (!token) window.location.href = '/pages/login.html';
      // Optionally: POST /auth/verify to check token
    },
    
    async loadFarmData() {
      const farmId = localStorage.getItem('selected_farm_id');
      const response = await fetch(`/api/farms/${farmId}`, {
        headers: { 'Authorization': `Bearer ${localStorage.getItem('access_token')}` }
      });
      this.farm = await response.json();
    },
    
    async toggleDevice(deviceId, newState) {
      // Optimistic update
      const device = this.devices.find(d => d.id === deviceId);
      device.state = newState;
      
      // Send command
      await fetch(`/api/devices/${deviceId}/commands`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ state: newState })
      });
      
      // Wait for WebSocket confirmation (server updates device.state)
    },
    
    connectWebSocket() {
      const token = localStorage.getItem('access_token');
      this.ws = new WebSocket(`wss://${window.location.host}/ws?token=${token}`);
      
      this.ws.onmessage = (event) => {
        const msg = JSON.parse(event.data);
        
        if (msg.type === 'device:update') {
          const device = this.devices.find(d => d.id === msg.device_id);
          if (device) device.state = msg.state;
        }
        
        if (msg.type === 'sensor:reading') {
          const sensor = this.sensors.find(s => s.id === msg.sensor_id);
          if (sensor) sensor.value = msg.value;
        }
      };
    }
  }
}).mount('#app');
```

---

## 📅 Schedule Automation UI (NEW in v2.0)

**Context**: Farmers need to create multi-step sequences for feeders, conveyor belts (see [`docs/AUTOMATION_USE_CASES.md`](../docs/AUTOMATION_USE_CASES.md))

### Example: Pulse Feeding Schedule UI

**User Flow**:
1. Click "Create Schedule" on dashboard
2. Select device: "Feeder 1"
3. Choose type: "Time-based" (runs at specific time)
4. Set time: "6:00 AM"
5. **Action Sequence Builder**:
   - Row 1: [ON] for [30] seconds
   - Row 2: [OFF] for [10] seconds
   - Row 3: [ON] for [30] seconds
   - Row 4: [OFF] for [10] seconds
   - Row 5: [OFF] until next schedule

**HTML Template** (`pages/schedules.html` - Action Sequence Builder):

```html
<div v-if="scheduleType === 'time_based'">
  <h3>Action Sequence</h3>
  <p class="help-text">Create multi-step patterns (e.g., pulse feeding: ON 30s, pause 10s, repeat)</p>
  
  <div v-for="(step, index) in actionSequence" :key="index" class="action-step">
    <label>Step {{ index + 1 }}</label>
    <select v-model="step.action">
      <option value="ON">Turn ON</option>
      <option value="OFF">Turn OFF</option>
    </select>
    <input type="number" v-model="step.duration" min="1" placeholder="Seconds">
    <button @click="removeStep(index)" class="btn-remove">Remove</button>
  </div>
  
  <button @click="addStep" class="btn-add">+ Add Step</button>
</div>

<button @click="createSchedule" class="btn-primary">Save Schedule</button>
```

**JavaScript** (`js/schedules.js`):

```javascript
data() {
  return {
    scheduleType: 'time_based',
    actionSequence: [
      { action: 'ON', duration: 30 },
      { action: 'OFF', duration: 10 }
    ]
  };
},

methods: {
  addStep() {
    this.actionSequence.push({ action: 'ON', duration: 10 });
  },
  
  removeStep(index) {
    this.actionSequence.splice(index, 1);
  },
  
  async createSchedule() {
    const payload = {
      farm_id: parseInt(localStorage.getItem('selected_farm_id')),
      coop_id: this.selectedCoopId,
      device_id: this.selectedDeviceId,
      schedule_type: this.scheduleType,
      time_value: this.timeValue,                          // e.g., "06:00"
      action_value: 'ON',
      action_duration: 3600,                               // Auto-turn-off after 1 hour
      action_sequence: JSON.stringify(this.actionSequence), // NEW FIELD
      priority: 5
    };
    
    await fetch(`/api/farms/${payload.farm_id}/schedules`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(payload)
    });
    
    // Redirect to schedules list
  }
}
```

**API Endpoint Used**: `POST /api/farms/{farm_id}/schedules` (see [`docs/implementation/API.md`](../docs/implementation/API.md#post-apifarmsfarm_idschedules))

**Backend Validation**:
- `action_sequence` must be valid JSON array
- Each step needs `action` (ON/OFF) and `duration` (seconds)
- Maximum 20 steps per sequence (prevent farmer confusion)

---

## 🌐 API Integration Patterns

### Authentication (JWT)

**Login Flow**:
```javascript
// POST /auth/login
const response = await fetch('/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ email, password })
});

const { access_token, refresh_token } = await response.json();
localStorage.setItem('access_token', access_token);
localStorage.setItem('refresh_token', refresh_token);

// Redirect to dashboard
window.location.href = '/pages/index.html';
```

**Protected Requests**:
```javascript
async function apiCall(endpoint, method = 'GET', body = null) {
  const headers = {
    'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
    'Content-Type': 'application/json'
  };
  
  const response = await fetch(endpoint, { method, headers, body: JSON.stringify(body) });
  
  if (response.status === 401) {
    // Token expired, refresh or redirect to login
    await refreshToken();
    return apiCall(endpoint, method, body); // Retry
  }
  
  return response.json();
}
```

**Token Refresh** (when access_token expires):
```javascript
async function refreshToken() {
  const response = await fetch('/auth/refresh', {
    method: 'POST',
    body: JSON.stringify({ refresh_token: localStorage.getItem('refresh_token') })
  });
  
  const { access_token } = await response.json();
  localStorage.setItem('access_token', access_token);
}
```

### Disease Detection

**Upload Image → Get AI Prediction**:
```javascript
async function uploadImage(imageFile) {
  const formData = new FormData();
  formData.append('image', imageFile);
  
  const response = await fetch('/api/ai/predict', {
    method: 'POST',
    headers: { 'Authorization': `Bearer ${localStorage.getItem('access_token')}` },
    body: formData
  });
  
  const result = await response.json();
  // {
  //   "disease": "Coccidiosis",
  //   "confidence": 0.99,
  //   "recommendation": "Isolate affected birds, administer treatment...",
  //   "treatment_options": ["Amprolium 20% solution", "Sulfonamides"]
  // }
  
  displayResult(result);
}
```

**Display Result** (Large font for farmer readability):
```html
<div class="prediction-result">
  <h2 style="font-size: 48px; color: {{result.confidence > 0.8 ? 'red' : 'orange'}}">
    {{ result.disease }}
  </h2>
  <p style="font-size: 32px;">Confidence: {{ (result.confidence * 100).toFixed(0) }}%</p>
  <p style="font-size: 24px;">{{ result.recommendation }}</p>
</div>
```

---

## 🔒 Security Best Practices

- ✅ **JWT Storage**: Use `localStorage` (XSS risk acceptable for this use case - elderly farmers won't have browser extensions)
- ✅ **HTTPS Only**: Production backend enforces HTTPS (see `middleware/.env` - `TLS_CERT`, `TLS_KEY`)
- ✅ **No Hardcoded URLs**: API base URL from `window.location.host` (works in dev & production)
- ✅ **Input Validation**: Validate image file size (< 5MB), type (PNG/JPEG) before upload
- ✅ **Error Handling**: Never show raw error messages to farmers (use friendly Khmer text)
- ✅ **WebSocket Auth**: Token passed in query param (`ws://server/ws?token=<jwt>`)

---

## 🌍 Internationalization (i18n)

**Language Toggle** (Top-right header):

```javascript
const translations = {
  en: {
    'dashboard.title': 'Dashboard',
    'device.on': 'On',
    'device.off': 'Off',
    'schedule.create': 'Create Schedule'
  },
  km: {
    'dashboard.title': 'ផ្ទាំងគ្រប់គ្រង',
    'device.on': 'បើក',
    'device.off': 'បិទ',
    'schedule.create': 'បង្កើតកាលវិភាគ'
  }
};

function setLanguage(lang) {
  localStorage.setItem('language', lang);
  document.querySelectorAll('[data-i18n]').forEach(el => {
    const key = el.getAttribute('data-i18n');
    el.textContent = translations[lang][key];
  });
}

// On page load
document.addEventListener('DOMContentLoaded', () => {
  const lang = localStorage.getItem('language') || 'km'; // Default: Khmer
  setLanguage(lang);
});
```

**HTML**:
```html
<h1 data-i18n="dashboard.title"></h1>
<button data-i18n="schedule.create"></button>
```

---

## 📱 Responsive Design (Mobile-First)

**Breakpoints**:
```css
/* Mobile: <= 480px (primary) */
@media (max-width: 480px) {
  .device-card { width: 100%; margin: 8px 0; }
  button { font-size: 48px; height: 60px; }
}

/* Tablet: 481px - 768px */
@media (min-width: 481px) and (max-width: 768px) {
  .device-grid { grid-template-columns: repeat(2, 1fr); }
}

/* Desktop: > 768px */
@media (min-width: 769px) {
  .device-grid { grid-template-columns: repeat(3, 1fr); }
}
```

**Accessibility (WCAG AAA)**:
```css
/* High contrast */
body { background: #FFFFFF; color: #000000; }

/* Large buttons */
.btn-primary {
  font-size: 48px;
  padding: 16px 32px;
  min-height: 60px;
  min-width: 200px;
}

/* Clear focus indicators */
button:focus {
  outline: 4px solid #0066CC;
  outline-offset: 2px;
}
```

---

## 🆘 Common Issues & Solutions

| Issue | Solution |
|-------|----------|
| **WebSocket disconnects** | Implement auto-reconnect with exponential backoff (5s, 10s, 20s) |
| **Images don't load on 4G** | Lazy load, use WebP format, show loading spinner |
| **Touch not registering** | Increase button size to 48px minimum, add `touch-action: manipulation` |
| **Khmer text renders as boxes** | Include Khmer font in `assets/fonts/` (e.g., Battambang.ttf), specify in CSS |
| **Token expired mid-session** | Catch 401 errors, call `/auth/refresh`, retry original request |
| **Device state out of sync** | Rely on WebSocket for authoritative state, not optimistic UI updates |

---

## 🧪 Development Tasks

### Add a New Page

1. **Create HTML**: `frontend/pages/new-page.html`
   - Include Vue CDN: `<script src="https://unpkg.com/vue@3/dist/vue.global.js"></script>`
   - Add `<div id="app">` wrapper

2. **Create JS**: `frontend/js/new-page.js`
   - Initialize Vue app with `createApp().mount('#app')`
   - Add `async mounted()` to load data

3. **Create CSS**: `frontend/css/styleNewPage.css`
   - Follow mobile-first breakpoints
   - Use minimum 48px buttons

4. **Add to Navbar**: Update `frontend/components/navbar.html`
   - Add new tab with icon + label

### Add a New API Call

1. **Find endpoint**: See [`docs/implementation/API.md`](../docs/implementation/API.md)
2. **Create helper** in `js/utils/api.js`:
   ```javascript
   export async function getDevices(farmId) {
     return apiCall(`/api/farms/${farmId}/devices`, 'GET');
   }
   ```
3. **Call from page**: `const devices = await getDevices(this.farmId);`

### Test Locally

```bash
cd frontend

# Option 1: Python server
python -m http.server 8080

# Option 2: Node.js http-server
npx http-server . -p 8080

# Visit: http://localhost:8080/pages/index.html
```

---

## 📘 Documentation Map

**AI Context Files** (component-specific guides):
- **This file**: [`frontend/AI_CONTEXT.md`](./AI_CONTEXT.md) - Vue.js patterns, API calls, schedule UI
- [`middleware/AI_CONTEXT.md`](../middleware/AI_CONTEXT.md) - Go API, database queries, WebSocket server
- [`ai-service/AI_CONTEXT.md`](../ai-service/AI_CONTEXT.md) - PyTorch model, FastAPI endpoints
- [`embedded/AI_CONTEXT.md`](../embedded/AI_CONTEXT.md) - ESP32 firmware, MQTT protocol
- [`docs/AI_CONTEXT.md`](../docs/AI_CONTEXT.md) - Documentation maintenance guide

**Master Guide**: [`AI_INSTRUCTIONS.md`](../AI_INSTRUCTIONS.md) - Read first for project overview

---

**Happy coding! 🚀 If elderly farmers can't use it, it's not done yet.**
