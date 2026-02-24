# Frontend Implementation - Vue.js 3 Migration

**Last Updated**: February 24, 2026  
**Status**: MVP Pages Complete (all pages exist; Vue.js migration planned post-MVP)  
**Tech Stack**: Vanilla HTML/CSS/JS (current MVP) → Vue.js 3 CDN (next phase)

---

## MVP Page Inventory (Feb 2026)

All pages are static HTML served by the Go backend via named routes in `middleware/main.go`.

| Page | Route | File | Status |
|------|-------|------|--------|
| Home / Dashboard | `/` and `/index.html` | `pages/index.html` | ✅ Live |
| Login | `/login` | `pages/login.html` | ✅ Live |
| Sign Up | `/register` | `pages/signup.html` | ✅ Live |
| Profile | `/profile` | `pages/profile.html` | ✅ Live |
| Settings | `/settings` | `pages/settings.html` | ✅ Live |
| AI Disease Detection | `/disease-detection` | `pages/disease-detection.html` | 🚧 Coming Soon |
| 404 Not Found | all unmatched | `pages/404.html` | ✅ Live |

### Disease Detection Page — Coming Soon Overlay

`pages/disease-detection.html` has a **full-screen overlay** injected immediately after `<body>` (visible above all page content):

```html
<!-- Remove this div when AI service is ready -->
<div id="coming-soon-overlay" style="position: fixed; inset: 0; z-index: 9999; ...">
  🔬 AI Disease Detection — Coming Soon
</div>
```

**To re-enable the page**: Remove the `<div id="coming-soon-overlay">...</div>` block (marked with comments). The full UI underneath is intact and functional.

---

## Migration Strategy

### Current State
- ✅ Vanilla HTML/CSS/JS files (`frontend/pages/*.html`)
- ✅ No build step required
- ✅ Served directly by Go backend via static routes
- ❌ Code duplication (navbar, header in every page)
- ❌ Manual DOM updates (water_level, temperature)
- ❌ No component reusability

### Target State (Phase by Phase)
- Phase 1: Vue.js 3 CDN (no build step)
- Phase 2: Component system
- Phase 3: Vite + TypeScript (optional)

---

## Phase 1: Add Vue.js 3 CDN

### Step 1: Update HTML Template

**Before (Vanilla):**
```html
<!-- frontend/pages/index.html -->
<!DOCTYPE html>
<html lang="km">
<head>
    <title>Tokkatot - Home</title>
    <link rel="stylesheet" href="../css/styleHome.css">
</head>
<body>
    <!-- Header duplicated in every page -->
    <div id="header-placeholder"></div>
    
    <div class="info-section">
        <h2>ព័ត៌មានបរិយាកាស</h2>
        <div class="info-card temperature">
            <span id="current-temp">--</span>
            <span class="unit">°C</span>
        </div>
        <div class="info-card humidity">
            <span id="current-humidity">--</span>
            <span class="unit">%</span>
        </div>
    </div>
    
    <script src="../js/index.js"></script>  <!-- Manual DOM updates -->
</body>
</html>
```

**After (Vue.js 3 CDN):**
```html
<!-- frontend/pages/index.html -->
<!DOCTYPE html>
<html lang="km">
<head>
    <title>Tokkatot - Home</title>
    <link rel="stylesheet" href="../css/styleHome.css">
    <!-- Add Vue.js 3 via CDN -->
    <script src="https://unpkg.com/vue@3/dist/vue.global.js"></script>
</head>
<body>
    <div id="app">
        <!-- Vue reactive data -->
        <navbar-component></navbar-component>
        
        <div class="info-section">
            <h2>ព័ត៌មានបរិយាកាស</h2>
            <div class="info-card temperature">
                <span>{{ currentTemp }}</span>  <!-- Auto-updates! -->
                <span class="unit">°C</span>
            </div>
            <div class="info-card humidity">
                <span>{{ currentHumidity }}</span>
                <span class="unit">%</span>
            </div>
        </div>
    </div>
    
    <script>
    const { createApp } = Vue
    
    createApp({
        data() {
            return {
                currentTemp: 0,
                currentHumidity: 0,
                coops: [],
                selectedCoop: null
            }
        },
        methods: {
            async fetchCoopData() {
                const token = localStorage.getItem('access_token')
                const res = await fetch('/v1/coops/current', {
                    headers: { 'Authorization': `Bearer ${token}` }
                })
                const data = await res.json()
                
                this.currentTemp = data.temperature
                this.currentHumidity = data.humidity
            },
            
            connectWebSocket() {
                const ws = new WebSocket('ws://localhost:3000/ws')
                
                ws.onmessage = (event) => {
                    const data = JSON.parse(event.data)
                    
                    // Reactive updates - UI changes automatically!
                    if (data.type === 'temperature') {
                        this.currentTemp = data.value
                    } else if (data.type === 'humidity') {
                        this.currentHumidity = data.value
                    }
                }
            }
        },
        mounted() {
            this.fetchCoopData()
            this.connectWebSocket()
        }
    }).mount('#app')
    </script>
</body>
</html>
```

---

## Phase 2: Component System

### Create Reusable Components

**File: `frontend/components/navbar.js`**
```javascript
// Navbar component (reusable across all pages)
app.component('navbar-component', {
    data() {
        return {
            user: null,
            currentFarm: null
        }
    },
    template: `
        <nav class="navbar">
            <div class="navbar-brand">
                <img src="/assets/images/tokkatot-logo.png" alt="Tokkatot">
            </div>
            <div class="navbar-menu">
                <a href="/dashboard">ផ្ទះ</a>
                <a href="/coops">សំបុក</a>
                <a href="/settings">ការកំណត់</a>
                <a href="#" @click="logout">ចេញ</a>
            </div>
            <div class="navbar-user">
                <span>{{ user?.name }}</span>
                <span class="farm-name">{{ currentFarm?.name }}</span>
            </div>
        </nav>
    `,
    methods: {
        logout() {
            localStorage.clear()
            window.location.href = '/login'
        }
    },
    mounted() {
        // Load user data
        const token = localStorage.getItem('access_token')
        if (token) {
            fetch('/v1/auth/me', {
                headers: { 'Authorization': `Bearer ${token}` }
            })
            .then(res => res.json())
            .then(data => {
                this.user = data.user
                this.currentFarm = data.farms[0]  // Default farm
            })
        }
    }
})
```

**File: `frontend/components/coop-card.js`**
```javascript
app.component('coop-card', {
    props: {
        coop: {
            type: Object,
            required: true
        }
    },
    template: `
        <div class="coop-card" @click="selectCoop">
            <h3>{{ coop.name }}</h3>
            <div class="coop-stats">
                <div class="stat">
                    <span class="icon">🐔</span>
                    <span class="value">{{ coop.current_count }}/{{ coop.capacity }}</span>
                </div>
                <div class="stat" :class="waterLevelClass">
                    <span class="icon">💧</span>
                    <span class="value">{{ coop.waterLevel }}%</span>
                </div>
                <div class="stat">
                    <span class="icon">🌡️</span>
                    <span class="value">{{ coop.temperature }}°C</span>
                </div>
            </div>
            <div v-if="coop.hasAlert" class="alert">
                ⚠️ Disease detected
            </div>
        </div>
    `,
    computed: {
        waterLevelClass() {
            return {
                'low-water': this.coop.waterLevel < 30,
                'ok-water': this.coop.waterLevel >= 30
            }
        }
    },
    methods: {
        selectCoop() {
            this.$router.push(`/coops/${this.coop.id}`)
        }
    }
})
```

**Usage in pages:**
```html
<!-- frontend/pages/coops.html -->
<div id="app">
    <navbar-component></navbar-component>
    
    <h1>ជ្រើសរើសសំបុកមាន់</h1>
    
    <div class="coop-grid">
        <coop-card 
            v-for="coop in coops" 
            :key="coop.id" 
            :coop="coop">
        </coop-card>
    </div>
</div>

<script src="../components/navbar.js"></script>
<script src="../components/coop-card.js"></script>
<script>
const { createApp } = Vue

createApp({
    data() {
        return {
            coops: []
        }
    },
    async mounted() {
        const res = await fetch('/v1/farms/current/coops')
        const data = await res.json()
        this.coops = data.coops
    }
}).mount('#app')
</script>
```

---

## Phase 3: Vite Build (Optional - Later)

### When to Add Build Step?
✅ **Add Vite when:**
- Need TypeScript
- Want HMR (Hot Module Replacement)
- Component count > 10
- Team size > 2 developers

### Setup Vite
```powershell
cd frontend

# Initialize Vue 3 project with Vite
npm create vite@latest . -- --template vue

# Install dependencies
npm install

# Install router, pinia (state management)
npm install vue-router pinia axios

# Dev server with HMR
npm run dev  # http://localhost:5173

# Production build
npm run build  # Creates frontend/dist/
```

**Vite config:**
```javascript
// vite.config.js
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
    plugins: [vue()],
    server: {
        port: 5173,
        proxy: {
            '/v1': 'http://localhost:3000'  // Proxy API to backend
        }
    },
    build: {
        outDir: 'dist',
        assetsDir: 'assets'
    }
})
```

---

## Key Patterns

### Authentication
```javascript
// frontend/js/auth.js
export const auth = {
    install(app) {
        app.config.globalProperties.$auth = {
            async login(phone, password) {
                const res = await fetch('/v1/auth/login', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ phone, password })
                })
                
                const data = await res.json()
                if (data.success) {
                    localStorage.setItem('access_token', data.data.access_token)
                    localStorage.setItem('refresh_token', data.data.refresh_token)
                    return true
                }
                return false
            },
            
            logout() {
                localStorage.clear()
                window.location.href = '/login'
            },
            
            getToken() {
                return localStorage.getItem ('access_token')
            },
            
            isAuthenticated() {
                return !!this.getToken()
            }
        }
    }
}

// Usage in components
this.$auth.login('012345678', 'password')
```

### API Helper
```javascript
// frontend/js/api.js
export const api = {
    async request(url, options = {}) {
        const token = localStorage.getItem('access_token')
        
        const res = await fetch(url, {
            ...options,
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
                ...options.headers
            }
        })
        
        const data = await res.json()
        
        // Auto-refresh token if expired
        if (data.error?.code === 'token_expired') {
            await this.refreshToken()
            return this.request(url, options)  // Retry
        }
        
        return data
    },
    
    async refreshToken() {
        const refresh = localStorage.getItem('refresh_token')
        const res = await fetch('/v1/auth/refresh', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ refresh_token: refresh })
        })
        
        const data = await res.json()
        if (data.success) {
            localStorage.setItem('access_token', data.data.access_token)
        }
    },
    
    // Shortcuts
    get(url) {
        return this.request(url)
    },
    
    post(url, body) {
        return this.request(url, {
            method: 'POST',
            body: JSON.stringify(body)
        })
    },
    
    patch(url, body) {
        return this.request(url, {
            method: 'PATCH',
            body: JSON.stringify(body)
        })
    }
}

// Usage
const coops = await api.get('/v1/farms/1/coops')
await api.post('/v1/devices/123/command', { action: 'on' })
```

### WebSocket Real-Time Updates
```javascript
// frontend/js/websocket.js
export class WebSocketManager {
    constructor(url) {
        this.url = url
        this.ws = null
        this.listeners = {}
        this.reconnectInterval = 5000
    }
    
    connect() {
        this.ws = new WebSocket(this.url)
        
        this.ws.onopen = () => {
            console.log('✅ WebSocket connected')
            this.emit('connected')
        }
        
        this.ws.onmessage = (event) => {
            const data = JSON.parse(event.data)
            this.emit(data.type, data)
        }
        
        this.ws.onclose = () => {
            console.log('❌ WebSocket disconnected, reconnecting...')
            setTimeout(() => this.connect(), this.reconnectInterval)
        }
    }
    
    on(event, callback) {
        if (!this.listeners[event]) {
            this.listeners[event] = []
        }
        this.listeners[event].push(callback)
    }
    
    emit(event, data) {
        if (this.listeners[event]) {
            this.listeners[event].forEach(cb => cb(data))
        }
    }
    
    send(type, data) {
        this.ws.send(JSON.stringify({ type, ...data }))
    }
}

// Usage in Vue app
import { WebSocketManager } from './websocket.js'

const ws = new WebSocketManager('ws://localhost:3000/ws')
ws.connect()

// Listen for sensor updates
ws.on('water_level', (data) => {
    app.waterLevel = data.value
})

ws.on('temperature', (data) => {
    app.temperature = data.value
})
```

---

## File Structure

```
frontend/
├── pages/                      ← HTML pages (Vue apps)
│   ├── login.html
│   ├── index.html              ← Dashboard
│   ├── coops.html              ← Coop selection
│   ├── coop-detail.html        ← Single coop view
│   ├── settings.html
│   └── profile.html
│
├── components/                 ← Vue components (JS files)
│   ├── navbar.js
│   ├── header.js
│   ├── coop-card.js
│   ├── device-control.js
│   └── alert-banner.js
│
├── js/                         ← Utilities
│   ├── auth.js                 ← Auth plugin
│   ├── api.js                  ← API helper
│   ├── websocket.js            ← WebSocket manager
│   └── utils.js                ← Formatters, validators
│
├── css/                        ← Styles (keep existing)
│   ├── styleHome.css
│   ├── stylenavbar.css
│   └── styleProfile.css
│
└── assets/                     ← Static files
    ├── images/
    ├── icons/
    └── fonts/
```

---

## Accessibility (Farmer-Friendly UI)

### Design Principles
```css
/* Large touch targets */
.button {
    min-width: 48px;
    min-height: 48px;
    font-size: 18px;
}

/* High contrast */
.text-primary {
    color: #000000;
    background: #FFFFFF;
}

.alert-danger {
    color: #FFFFFF;
    background: #DC2626;
}

/* Large fonts for readability */
body {
    font-size: 16px;  /* Base */
}

h1 {
    font-size: 32px;  /* Khmer readable */
}

.value {
    font-size: 48px;  /* Sensor values big! */
}
```

### Khmer Language Support
```javascript
// frontend/js/i18n.js
export const translations = {
    km: {
        home: 'ផ្ទះ',
        coops: 'សំបុក',
        settings: 'ការកំណត់',
        waterLevel: 'កម្រិតទឹក',
        temperature: 'សីតុណ្ហភាព',
        humidity: 'សំណើម',
        turnOn: 'បើក',
        turnOff: 'បិទ',
        alert: 'ការព្រមានព័ត៌មាន'
    },
    en: {
        home: 'Home',
        coops: 'Coops',
        settings: 'Settings',
        waterLevel: 'Water Level',
        temperature: 'Temperature',
        humidity: 'Humidity',
        turnOn: 'Turn On',
        turnOff: 'Turn Off',
        alert: 'Alert'
    }
}

// Usage
const t = (key) => translations[currentLang][key]
```

---

## Schedule/Automation UI (NEW in v2.0)

**Purpose**: Allow farmers to create multi-step automation sequences for feeders, conveyors, pumps.

**Real Farmer Need**: "I want my feeder to pulse: ON 30 seconds, pause 10 seconds, repeat 2 times, then turn OFF until next feeding time" (prevents chicken crowding at feed bowl)

### Schedule Pages

**`frontend/pages/schedules.html`** - Schedule management
- List all schedules with current state (enabled/disabled, next execution time)
- Create new schedule button (leads to create form)
- Edit/delete existing schedules
- Manual "Execute Now" button

**`frontend/pages/schedule-create.html`** - Create/Edit schedule form

### Action Sequence Builder Component

**Visual Multi-Step Pattern Editor** (inspired by automation tools like IFTTT)

```html
<!-- frontend/components/action-sequence-builder.js -->
<template>
  <div class="sequence-builder">
    <h3>បង្កើតលំដាប់សកម្មភាព (Action Sequence)</h3>
    <p class="help-text">
      បង្កើតលំដាប់ជាបន្តបន្ទាប់: បើក → ផ្អាក → បើក → ផ្អាក → បិទ<br>
      Create multi-step pattern: ON → pause → ON → pause → OFF
    </p>
    
    <!-- List of steps -->
    <div v-for="(step, index) in actionSequence" :key="index" class="step-row">
      <span class="step-number">{{ index + 1 }}</span>
      
      <!-- Action type dropdown -->
      <select v-model="step.action" class="action-select">
        <option value="ON">បើក (ON)</option>
        <option value="OFF">បិទ (OFF)</option>
      </select>
      
      <!-- Duration input -->
      <input 
        type="number" 
        v-model="step.duration" 
        min="0" 
        placeholder="រយៈពេល (វិនាទី)"
        class="duration-input"
      />
      <span class="unit">វិនាទី (seconds)</span>
      
      <!-- Remove step button -->
      <button @click="removeStep(index)" class="btn-remove" :disabled="actionSequence.length <= 1">
        ✕
      </button>
    </div>
    
    <!-- Add step button -->
    <button @click="addStep" class="btn-add-step">
      + បន្ថែមមួយជំហានទៀត (Add Step)
    </button>
    
    <!-- Visual timeline preview -->
    <div class="timeline-preview">
      <h4>មើលមុន (Preview)</h4>
      <div class="timeline">
        <div v-for="(step, index) in actionSequence" 
             :key="index"
             :class="['timeline-block', step.action === 'ON' ? 'on' : 'off']"
             :style="{ width: getStepWidth(step.duration) }">
          <span class="timeline-label">{{ step.action }}</span>
          <span class="timeline-duration">{{ step.duration }}s</span>
        </div>
      </div>
      <p class="total-time">សរុប (Total): {{ totalDuration }}s ({{ formatMinutes(totalDuration) }})</p>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      actionSequence: [
        { action: 'ON', duration: 30 },
        { action: 'OFF', duration: 10 }
      ]
    }
  },
  
  computed: {
    totalDuration() {
      return this.actionSequence.reduce((sum, step) => {
        return sum + (step.duration === 0 ? 0 : parseInt(step.duration))
      }, 0)
    }
  },
  
  methods: {
    addStep() {
      this.actionSequence.push({ action: 'ON', duration: 10 })
    },
    
    removeStep(index) {
      if (this.actionSequence.length > 1) {
        this.actionSequence.splice(index, 1)
      }
    },
    
    getStepWidth(duration) {
      if (duration === 0) return '50px'  // "Until next schedule" step
      const maxDuration = Math.max(...this.actionSequence.map(s => s.duration))
      return `${(duration / maxDuration) * 200}px`
    },
    
    formatMinutes(seconds) {
      const mins = Math.floor(seconds / 60)
      const secs = seconds % 60
      return mins > 0 ? `${mins}នាទី ${secs}វិនាទី` : `${secs}វិនាទី`
    },
    
    getSequenceJSON() {
      // Return JSON for API submission
      return JSON.stringify(this.actionSequence)
    }
  }
}
</script>
```

### Create Schedule Form

**Full schedule creation UI** with action sequence builder:

```javascript
// frontend/pages/schedule-create.html
createApp({
  data() {
    return {
      scheduleType: 'time_based',
      deviceId: null,
      scheduleName: '',
      cronExpression: '0 6 * * *',  // 6 AM daily
      actionValue: 'on',
      actionDuration: null,       // Simple auto-off (if not using sequence)
      actionSequence: null,       // Multi-step pattern (if using)
      useSequence: false,         // Toggle between simple/sequence
      priority: 5,
      devices: []
    }
  },
  
  async mounted() {
    // Load available devices
    const farmId = localStorage.getItem('selected_farm_id')
    const response = await fetch(`/api/farms/${farmId}/devices`, {
      headers: { 'Authorization': `Bearer ${localStorage.getItem('access_token')}` }
    })
    this.devices = await response.json()
  },
  
  methods: {
    async createSchedule() {
      const payload = {
        device_id: this.deviceId,
        name: this.scheduleName,
        schedule_type: this.scheduleType,
        cron_expression: this.cronExpression,
        action_value: this.actionValue,
        priority: this.priority,
        is_enabled: true
      }
      
      // Add simple auto-off OR multi-step sequence (mutually exclusive)
      if (!this.useSequence && this.actionDuration) {
        payload.action_duration = parseInt(this.actionDuration)
      } else if (this.useSequence && this.$refs.sequenceBuilder) {
        payload.action_sequence = this.$refs.sequenceBuilder.getSequenceJSON()
      }
      
      const farmId = localStorage.getItem('selected_farm_id')
      const response = await fetch(`/api/farms/${farmId}/schedules`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('access_token')}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(payload)
      })
      
      if (response.ok) {
        alert('កាលវិភាគត្រូវបានបង្កើតដោយជោគជ័យ! (Schedule created successfully!)')
        window.location.href = '/pages/schedules.html'
      } else {
        const error = await response.json()
        alert(`កំហុស: ${error.message}`)
      }
    }
  }
})
```

### CSS for Schedule UI

```css
/* frontend/css/styleSchedules.css */

.sequence-builder {
  background: #F9FAFB;
  padding: 24px;
  border-radius: 8px;
  margin: 16px 0;
}

.step-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: white;
  border-radius: 4px;
  margin-bottom: 8px;
}

.step-number {
  font-size: 20px;
  font-weight: bold;
  color: #6B7280;
  min-width: 32px;
}

.action-select {
  font-size: 18px;
  padding: 12px;
  min-height: 48px;  /* Farmer accessibility */
  border: 2px solid #D1D5DB;
  border-radius: 4px;
  flex: 1;
}

.duration-input {
  font-size: 18px;
  padding: 12px;
  min-height: 48px;
  width: 120px;
  border: 2px solid #D1D5DB;
  border-radius: 4px;
}

.btn-remove {
  min-width: 48px;
  min-height: 48px;
  background: #EF4444;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 24px;
  cursor: pointer;
}

.btn-add-step {
  min-height: 48px;
  width: 100%;
  background: #10B981;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 18px;
  margin-top: 8px;
  cursor: pointer;
}

/* Timeline preview */
.timeline-preview {
  margin-top: 24px;
  padding: 16px;
  background: white;
  border-radius: 4px;
}

.timeline {
  display: flex;
  gap: 4px;
  margin: 16px 0;
  height: 60px;
  align-items: center;
}

.timeline-block {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  height: 100%;
  border-radius: 4px;
  transition: all 0.3s;
}

.timeline-block.on {
  background: linear-gradient(135deg, #10B981, #059669);
  color: white;
}

.timeline-block.off {
  background: linear-gradient(135deg, #6B7280, #4B5563);
  color: white;
}

.timeline-label {
  font-size: 14px;
  font-weight: bold;
}

.timeline-duration {
  font-size: 12px;
  opacity: 0.8;
}

.total-time {
  font-size: 16px;
  font-weight: bold;
  color: #1F2937;
  margin-top: 8px;
}
```

### Example Use Cases (See `docs/AUTOMATION_USE_CASES.md`)

**Pulse Feeding** (Prevents chicken crowding):
- 6 AM: Motor ON 30s → pause 10s → ON 30s → pause 10s → OFF
- 12 PM: Same pattern
- 6 PM: Same pattern

**Conveyor Belt Cleaning** (Prevents clogging):
- ON 2 minutes → pause 30s → ON 2 minutes → OFF
- Runs every 6 hours

**Climate Control** (Gradual temperature change):
- Fan 25% → wait 5min → 50% → wait 5min → 75% → stay until temp < 30°C

### API Integration

**Submit schedule with action_sequence**:
```javascript
POST /api/farms/{farm_id}/schedules
{
  "device_id": "feeder_uuid",
  "name": "Pulse Feeding - Morning",
  "schedule_type": "time_based",
  "cron_expression": "0 6 * * *",
  "action_value": "on",
  "action_sequence": [
    {"action": "ON", "duration": 30},
    {"action": "OFF", "duration": 10},
    {"action": "ON", "duration": 30},
    {"action": "OFF", "duration": 10},
    {"action": "OFF", "duration": 0}  // Stay OFF until next schedule
  ],
  "priority": 5
}
```

**Backend validates**:
- Maximum 20 steps per sequence
- Each step has `action` (ON/OFF) and `duration` (0-3600 seconds)
- Total duration < 1 hour (prevents runaway sequences)

---

## Migration Checklist

**Phase 1: CDN Setup**
- [ ] Add Vue.js 3 CDN to all pages
- [ ] Convert `index.html` to Vue reactive data
- [ ] Convert `login.html` to Vue form handling
- [ ] Test all pages still work

**Phase 2: Components**
- [ ] Create `navbar.js` component
- [ ] Create `coop-card.js` component
- [ ] Create `device-control.js` component
- [ ] Remove duplicated HTML code

**Phase 3: Real-Time**
- [ ] Implement WebSocket manager
- [ ] Connect to backend WebSocket
- [ ] Live water level updates
- [ ] Live temperature updates
- [ ] Live alerts

**Phase 4: Authentication**
- [ ] Create auth plugin
- [ ] Implement login flow
- [ ] Token auto-refresh
- [ ] Protected routes

**Phase 5: Polish**
- [ ] PWA service worker (offline)
- [ ] Push notifications
- [ ] Loading states
- [ ] Error handling
- [ ] Khmer translations

---

## Testing

```javascript
// Test Vue reactive data
console.log('Initial temp:', app.currentTemp)
app.currentTemp = 30  // Should update UI immediately
console.log('New temp:', app.currentTemp)

// Test WebSocket
ws.emit('temperature', { value: 35 })  // Should update UI

// Test API
await api.post('/v1/devices/123/command', { action: 'on' })
// Should see device turn on in backend logs
```

---

## Next Steps

1. Start with `login.html` (simplest page)
2. Migrate to `index.html` (dashboard with real-time)
3. Create reusable components
4. Add WebSocket for live updates
5. Test on low-end Android devices

**Goal**: Farmer-friendly, reactive, real-time IoT dashboard! 🐔💧📱
