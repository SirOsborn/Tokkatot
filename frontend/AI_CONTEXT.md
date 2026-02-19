# 🤖 AI Context: Vue.js 3 Frontend

**Directory**: `frontend/`  
**Your Role**: User interface, real-time updates, device control, offline support  
**Tech Stack**: Vue.js 3, HTML5, CSS3, JavaScript (vanilla, no build system)  

---

## 🎯 What You're Building

**Mobile-First Web Application** (browser-based)
- **Target**: Elderly farmers in Cambodia with 1-2GB RAM phones, 4G network
- **Language**: Khmer & English toggle
- **Accessibility**: 48px+ fonts, high contrast colors, WCAG AAA compliant
- **Connectivity**: Real-time WebSocket + offline queue support (Service Workers)
- **Responsiveness**: Works on phones (320px), tablets (768px), desktops (1024px+)

**Key Features**:
- Dashboard: Live device status, temperature/humidity
- Disease Detection: Upload images, see AI predictions
- Device Control: Toggle lights, pumps, fans, feeders
- Scheduling: Create automation rules (if temp > 30°C, turn on fan)
- Profile: Edit user info, manage sessions
- Monitoring: View alerts, event history

---

## 📁 File Structure

```
frontend/
├── pages/
│   ├── index.html          # Home/Dashboard (primary page)
│   ├── disease-detection.html  # AI disease detection
│   ├── login.html          # Login/Signup
│   ├── profile.html        # User profile
│   ├── settings.html       # App settings
│   └── 404.html            # Error page
├── components/
│   ├── header.html         # Shared header
│   └── navbar.html         # Bottom tab navigation
├── js/
│   ├── index.js            # Home page logic
│   ├── disease-detection.js # Disease detection logic
│   ├── login.js            # Login/signup logic
│   ├── profile.js          # Profile page logic
│   ├── header.js           # Header component logic
│   ├── navbar.js           # Navbar component logic
│   ├── scriptHome.js       # Additional home page scripts
│   ├── scriptSettings.js   # Settings page scripts
│   └── libs/               # External libraries (if any)
├── css/
│   ├── styleHeader.css     # Header styles
│   ├── styleHome.css       # Home page styles
│   ├── stylenavbar.css     # Navbar styles
│   ├── styleProfile.css    # Profile page styles
│   ├── styleSettings.css   # Settings page styles
│   ├── loginSignUp.css     # Login/signup styles
│   └── (one CSS per page module)
├── assets/
│   ├── images/             # PNG, JPG images
│   ├── icons/              # SVG icons
│   └── fonts/              # Khmer, English fonts
└── AI_CONTEXT.md           # This file
```

---

## 🚀 Getting Started

### Local Development

```bash
cd frontend

# Simply open in browser (no build step!)
# Option 1: Open directly
file:///path/to/tokkatot/frontend/pages/index.html

# Option 2: Use live server
npx http-server . -p 8080

# Then visit http://localhost:8080/pages/index.html
```

### Testing on Mobile

```bash
# If using local http-server
npx http-server . -p 8080

# Access from phone: http://<LAPTOP_IP>:8080/pages/index.html
# Example: http://192.168.1.100:8080/pages/index.html
```

---

## 🎨 Design Philosophy

### Farmer-Centric UI Principles

1. **Simplicity Over Features**
   - Max 3 main actions per screen
   - Avoid: Complex menus, too many options
   - Prefer: Single large buttons, clear labels in Khmer

2. **Accessibility First**
   - Minimum 48px buttons (touch-friendly for older hands)
   - High contrast: Black (text) on white (background)
   - Sans-serif fonts only (easier to read)
   - Language toggle easy to find

3. **Mobile First**
   - Assume 4G network (not always stable)
   - Support offline mode (queue commands locally)
   - Minimize data transfer (lazy load images)
   - Test on low-end phones (2GB RAM)

4. **Visual Feedback**
   - Loading indicators (user knows it's working)
   - Success/error messages (user knows what happened)
   - Real-time status updates (WebSocket pushes)

---

## 🎭 Page Structure

### Home / Dashboard (`pages/index.html`)

**Displays**:
- Top card: Farm name, weather, time
- Device status grid: Water Pump (on/off), Light (on/off), etc
- Temperature/Humidity gauge (from sensors)
- Alert summary: Recent issues
- Quick action buttons

**Interactions**:
- Tap device card to open control
- Pull-down refresh (reload sensor data)
- Click alert to view details

### Disease Detection (`pages/disease-detection.html`)

**Features**:
- Camera button: Open phone camera
- Upload button: Select from photos
- AI prediction display:
  - Disease name (large, clear)
  - Confidence percentage
  - Treatment recommendations (numbered steps)
  - "Take another photo" if uncertain

**Real-time flow**:
```
1. User takes photo
2. Image sent to Go API: POST /api/ai/predict
3. Go API forwards to FastAPI: POST http://localhost:8000/predict
4. FastAPI returns disease + confidence
5. Frontend displays result with recommendations
6. Result saved to database
```

### Login / Signup (`pages/login.html`)

**Fields**:
- Email OR Phone (user's choice, not both required)
- Password (minimum 8 characters)
- Device name (optional: "Samsung A12")

**Simple flow**:
```
Existing user? → Login → Dashboard
New user? → Signup → Verify email → Set password → Dashboard
```

### Profile (`pages/profile.html`)

**Shows**:
- User name
- Email / Phone (whichever was used)
- Farms (list owned/managed farms)
- Active sessions (logout from devices)
- Change password button

### Settings (`pages/settings.html`)

**Options**:
- Language: Khmer / English
- Notifications: On/Off
- Theme: Light / Dark (optional)
- Logout button

---

## 🔧 Key JavaScript Functions

### `js/index.js` (Dashboard)

```javascript
// Initialize page
document.addEventListener('DOMContentLoaded', () => {
  validateJWT();           // Check token
  loadFarmData();          // Get dashboard data
  connectWebSocket();      // Real-time updates
  setupDeviceControls();   // Button listeners
});

// Functions
async function loadFarmData() {
  // GET /api/farms/{farm_id} → Display farm info
}

function connectWebSocket() {
  // Connect to ws://server/ws?token=<jwt>
  // Listen for device:update events → Update UI
}

function toggleDevice(deviceId, state) {
  // Send: POST /api/devices/{id}/commands
  // Update: UI immediately (optimistic)
  // Wait: Confirmation from server
}
```

### `js/disease-detection.js`

```javascript
async function predictDisease(imageFile) {
  // 1. Validate: File size < 5MB, type PNG/JPEG
  // 2. Upload: POST /api/ai/predict with image
  // 3. Display: Disease name, confidence, treatment
  // 4. Log: Save prediction history
}

function displayPrediction(result) {
  // Show result.disease in large text
  // Show result.confidence as percentage
  // Show result.treatment_options as numbered list
  // Show recommendation as orange alert box
}
```

### `js/login.js`

```javascript
async function handleLogin(event) {
  // Validate: email/phone not empty, password ≥ 8 chars
  // POST: /auth/login with credentials
  // Store: JWT tokens in localStorage
  // Redirect: to dashboard
}

async function handleSignup(event) {
  // Validate: confirm passwords match
  // POST: /auth/signup
  // Redirect: to email verification
}
```

---

## 🌐 API Integration

### Authentication

```javascript
// Store tokens
localStorage.setItem('access_token', response.access_token);
localStorage.setItem('refresh_token', response.refresh_token);

// Use in requests
fetch('/api/farms', {
  headers: {
    'Authorization': `Bearer ${localStorage.getItem('access_token')}`
  }
});

// Refresh token
async function refreshToken() {
  const response = await fetch('/auth/refresh', {
    method: 'POST',
    body: JSON.stringify({
      refresh_token: localStorage.getItem('refresh_token')
    })
  });
  // Update tokens...
}
```

### WebSocket (Real-time Updates)

```javascript
let ws;

function connectWebSocket() {
  const token = localStorage.getItem('access_token');
  ws = new WebSocket(`wss://server/ws?token=${token}`);

  ws.onmessage = (event) => {
    const message = JSON.parse(event.data);
    
    if (message.type === 'device:update') {
      // Update device status in UI
      updateDeviceCard(message.device_id, message.state);
    }
    
    if (message.type === 'alert:triggered') {
      showNotification(message.message);
    }
  };
}

// Send command
function sendDeviceCommand(deviceId, command) {
  ws.send(JSON.stringify({
    type: 'device:command',
    device_id: deviceId,
    command: command
  }));
}
```

### Disease Detection

```javascript
async function uploadImageForPrediction(imageFile) {
  const formData = new FormData();
  formData.append('image', imageFile);
  
  const response = await fetch('/api/ai/predict', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${localStorage.getItem('access_token')}`
    },
    body: formData
  });
  
  const result = await response.json();
  // {
  //   "disease": "Coccidiosis",
  //   "confidence": 0.99,
  //   "recommendation": "Isolate affected birds...",
  //   "treatment_options": [...]
  // }
  
  displayResult(result);
}
```

---

## 📝 Code Guidelines

### ✅ DO:
- Use semantic HTML (`<button>`, `<input>`, `<form>`)
- Add `aria-label` for accessibility
- Test font sizes (48px minimum for buttons)
- Use clear color contrast (WCAG AAA)
- Validate user input on client-side
- Show loading indicators while fetching
- Handle network errors gracefully
- Store tokens securely in localStorage
- Support both languages (Khmer/English)
- Lazy load images over slow networks

### ❌ DON'T:
- Use `eval()` or dynamic code execution
- Trust API responses without validation
- Hardcode API URLs (use config/env)
- Store passwords in localStorage
- Make unnecessary API calls
- Use tiny fonts (< 48px for buttons)
- Forget to validate image uploads
- Ignore WebSocket disconnects
- Use color alone to convey meaning (add text)
- Assume fast internet everywhere

---

## 🔒 Security Checklist

- ✅ JWT tokens in localStorage (with XSS prevention)
- ✅ HTTPS only in production
- ✅ No sensitive data in localStorage
- ✅ Validate image uploads on client (before sending)
- ✅ CSRF tokens in forms (if applicable)
- ✅ Sanitize user input (prevent XSS)
- ✅ Rate limit file uploads (max 5MB)
- ✅ Graceful error messages (no sensitive info leakage)

---

## 📊 Responsive Design

### Breakpoints

```css
/* Mobile: <= 480px */
@media (max-width: 480px) {
  button { font-size: 48px; }
  .card { margin: 8px; }
}

/* Tablet: 481px - 768px */
@media (min-width: 481px) and (max-width: 768px) {
  .grid { grid-template-columns: repeat(2, 1fr); }
}

/* Desktop: > 768px */
@media (min-width: 769px) {
  .grid { grid-template-columns: repeat(3, 1fr); }
}
```

### Layout Patterns

```html
<!-- Header -->
<header>
  <h1>Tokkatot</h1>
  <button class="menu-toggle">☰</button>
</header>

<!-- Main Content -->
<main>
  <div class="device-grid">
    <div class="device-card">...</div>
  </div>
</main>

<!-- Bottom Navigation -->
<nav class="navbar">
  <a href="/">Home</a>
  <a href="/disease">Disease</a>
  <a href="/profile">Profile</a>
</nav>
```

---

## 🌍 Internationalization (i18n)

### Language Toggle

```javascript
const translations = {
  en: {
    'home.title': 'Home',
    'device.on': 'On',
    'device.off': 'Off'
  },
  km: {
    'home.title': 'ដើម',
    'device.on': 'បើក',
    'device.off': 'បិទ'
  }
};

function setLanguage(lang) {
  localStorage.setItem('language', lang);
  document.querySelectorAll('[data-i18n]').forEach(el => {
    const key = el.getAttribute('data-i18n');
    el.textContent = translations[lang][key];
  });
}

// In HTML
<h1 data-i18n="home.title"></h1>
```

---

## 🌐 Offline Support (Service Workers)

```javascript
// Register service worker for offline
if ('serviceWorker' in navigator) {
  navigator.serviceWorker.register('/sw.js');
}

// Queue commands when offline
const offlineQueue = [];

window.addEventListener('offline', () => {
  console.log('Device is offline');
  // Queue future commands
});

window.addEventListener('online', () => {
  console.log('Device is back online');
  // Sync queued commands
  syncOfflineQueue();
});
```

---

## 🆘 Common Issues & Solutions

### Issue: WebSocket disconnects
**Fix**: Implement auto-reconnect with exponential backoff

### Issue: Images don't load on 4G
**Fix**: Lazy load images, use webp format, add loading spinner

### Issue: Touch not registering on buttons
**Fix**: Increase button size to 48px minimum, add touch event handlers

### Issue: Khmer text renders as boxes
**Fix**: Include Khmer font in assets (e.g., Battambang.ttf), specify in CSS

---

## 📚 Key Documents

- `IG_SPECIFICATIONS_FRONTEND.md` - UI/UX specs, accessibility standards
- `IG_SPECIFICATIONS_API.md` - API endpoints this frontend calls
- `01_SPECIFICATIONS_ARCHITECTURE.md` - How frontend fits in system

---

## 🧪 Testing Checklist

- ✅ Test on small screen (320px width)
- ✅ Test on phone with 4G (simulate slow network)
- ✅ Test with large fonts (zoom 150%)
- ✅ Test Khmer/English toggle
- ✅ Test offline mode (disconnect network)
- ✅ Test WebSocket disconnect/reconnect
- ✅ Test image upload (max 5MB)
- ✅ Test device control (real device)

---

## 🎯 Your Next Tasks

1. **Create page layouts** - HTML structure for each page
2. **Implement navigation** - Navbar, header routing
3. **Add styling** - Responsive CSS, accessibility colors
4. **Implement API calls** - Connect to Go API endpoints
5. **Add real-time** - WebSocket integration
6. **Test thoroughly** - Mobile, accessibility, offline

---

**Happy coding! 🚀 Remember: If elderly farmers can't use it, it's not done yet.**
