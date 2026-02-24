# Technology Stack - Tokkatot 2.0

**Version**: 2.1 (Finalized)  
**Last Updated**: February 23, 2026

---

## Quick Overview

| Layer | Technology | Version | Why |
|-------|-----------|---------|-----|
| **Backend** | Go + Fiber | 1.23 + v2.52 | Single binary, 30k req/sec, low memory |
| **Frontend** | Vue.js 3 (CDN) | 3.4+ | Reactive UI, 40KB, no build step initially |
| **Database** | PostgreSQL | 17+ | ACID, jsonb, full-text search |
| **AI Service** | Python + PyTorch | 3.12 + 2.0 | Disease detection, ensemble model |
| **Embedded** | ESP32 + C | ESP-IDF 5.0 | IoT sensors, MQTT, low power |
| **Deployment** | Docker + systemd | - | Single-server, VPS-ready |

---

## Backend: Go 1.23 + Fiber v2

### Why Go?
✅ **Single binary deployment** - Copy `backend.exe`, done  
✅ **Low memory** - 20-50MB vs Node.js 100-200MB  
✅ **Fast** - 30k requests/sec (3x faster than Node.js)  
✅ **Concurrency** - Goroutines handle 1000+ IoT device connections  
✅ **Type safety** - Catch errors at compile time  

### Core Dependencies
```go
// go.mod
go 1.23

require (
    github.com/gofiber/fiber/v2 v2.52.6      // Web framework
    github.com/golang-jwt/jwt/v4 v4.5.0      // JWT auth
    github.com/lib/pq v1.10.9                // PostgreSQL driver
    github.com/joho/godotenv v1.5.1          // .env file
    github.com/google/uuid v1.6.0            // UUID generation
    golang.org/x/crypto v0.19.0              // Bcrypt
)
```

### API Framework: Fiber v2
- Express.js-like syntax (easy to learn)
- Zero memory allocations (fast)
- Built-in WebSocket support (IoT real-time)

**Example:**
```go
app := fiber.New()
app.Get("/health", func(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{"status": "ok"})
})
```

---

## Frontend: Vue.js 3 (Progressive Enhancement)

### Migration Strategy: Vanilla → Vue.js 3 CDN

**Phase 1: CDN Setup (No Build Step)**
```html
<!-- Add to every page -->
<script src="https://unpkg.com/vue@3/dist/vue.global.js"></script>
<script>
const { createApp } = Vue
const app = createApp({
  data() {
    return {
      waterLevel: 0,
      coops: [],
      currentCoop: null
    }
  },
  methods: {
    async fetchCoops() {
      const res = await fetch('/v1/coops')
      this.coops = await res.json()
    }
  },
  mounted() {
    this.fetchCoops()
    // WebSocket for real-time updates
    const ws = new WebSocket('ws://localhost:3000/ws')
    ws.onmessage = (e) => {
      const data = JSON.parse(e.data)
      if (data.type === 'water_level') {
        this.waterLevel = data.value
      }
    }
  }
}).mount('#app')
</script>
```

**Phase 2: Component System (Later)**
```javascript
// components/CoopCard.js
app.component('coop-card', {
  props: ['coop'],
  template: `
    <div class="coop-card">
      <h3>{{ coop.name }}</h3>
      <p>🐔 {{ coop.current_count }}/{{ coop.capacity }}</p>
      <p>💧 {{ coop.waterLevel }}%</p>
      <p>🌡️ {{ coop.temperature }}°C</p>
    </div>
  `
})
```

**Phase 3: Build Optimization (Optional)**
- Add Vite for production builds
- TypeScript for type safety
- Tailwind CSS for styling

### Why Vue.js 3 (not React)?
✅ **Smaller** - 40KB vs React 130KB  
✅ **CDN-friendly** - No build step required  
✅ **Reactive** - Perfect for IoT dashboards (auto-update water levels)  
✅ **Progressive** - Migrate page by page (not all-or-nothing)  
✅ **Mobile performance** - Better on low-end Android devices  

---

## Database: PostgreSQL 17

### Why PostgreSQL?
✅ **ACID compliance** - Data integrity for IoT readings  
✅ **JSONB** - Store flexible device metadata  
✅ **Performance** - 10k+ writes/sec with proper indexing  
✅ **Full-text search** - Search diseases, alerts (Khmer + English)  
✅ **Free & Open Source** - No licensing costs  

### Schema Design
```sql
-- 14 tables, coop-centric design
users → farms → coops → devices
                  ↓
              schedules, device_commands, event_logs
                  
registration_keys (on-site verification)
```

**Performance:**
- 16+ indexes for fast queries
- Composite indexes on (farm_id, coop_id)
- Connection pooling (25 max)

---

## AI Service: Python 3.12 + PyTorch 2.0

### Why Python for AI?
✅ **PyTorch ecosystem** - Pre-trained models, GPU support  
✅ **FastAPI** - Modern async API (similar to Fiber)  
✅ **Separate service** - Isolate heavy AI from backend  

### Architecture
```
Go Backend (Port 3000)
    ↓ HTTP POST /predict
Python AI Service (Port 8000)
    ↓ PyTorch Ensemble Model
EfficientNetB0 + DenseNet121
    ↓ Disease confidence
Return: { disease: "Coccidiosis", confidence: 0.87 }
```

### Docker Deployment
```yaml
# docker-compose.yml
services:
  ai-service:
    image: tokkatot-ai:latest
    ports: ["8000:8000"]
    deploy:
      resources:
        limits: { cpus: '2', memory: 4G }
```

---

## Embedded: ESP32 + Raspberry Pi

### ESP32 (Sensors)
- **Framework**: ESP-IDF 5.0 (C language)
- **Sensors**: DHT22, Ultrasonic HC-SR04, Camera
- **Communication**: MQTT over WiFi
- **Power**: 5V, deep sleep mode

### Raspberry Pi (Main Controller)
- **OS**: Raspberry Pi OS Lite (headless)
- **Language**: Python 3 (simpler than C for farmers' staff)
- **Role**: Coordinate ESP32 sensors, send data to cloud
- **Communication**: MQTT broker, HTTP to backend

**Example:**
```python
# Raspberry Pi controller
import paho.mqtt.client as mqtt
import requests

def on_message(client, userdata, msg):
    # Receive from ESP32 sensor
    data = json.loads(msg.payload)
    
    # Send to cloud backend
    requests.post('http://api.tokkatot.com/v1/devices/readings', 
                  json=data, 
                  headers={'Authorization': f'Bearer {token}'})

client = mqtt.Client()
client.on_message = on_message
client.connect("localhost", 1883)
client.subscribe("coop/+/sensors")
client.loop_forever()
```

---

## Deployment: Single VPS

### Production Setup
```
┌─────────────────────────────────────┐
│  VPS ($5-10/month)                  │
│  Ubuntu 22.04 LTS                   │
│                                      │
│  ┌──────────────────────────────┐  │
│  │ systemd service: backend     │  │ ← Go binary (port 3000)
│  └──────────────────────────────┘  │
│  ┌──────────────────────────────┐  │
│  │ Docker: ai-service           │  │ ← Python (port 8000)
│  └──────────────────────────────┘  │
│  ┌──────────────────────────────┐  │
│  │ PostgreSQL 17                │  │ ← Database (port 5432)
│  └──────────────────────────────┘  │
│  ┌──────────────────────────────┐  │
│  │ Caddy (reverse proxy)        │  │ ← HTTPS (port 443)
│  └──────────────────────────────┘  │
└─────────────────────────────────────┘
```

### Why Not Kubernetes/Complex Setup?
✅ **Budget** - Farmers can't afford $50+/month  
✅ **Simplicity** - Single server easier to maintain  
✅ **Sufficient** - 100-200 farms = ~10k devices (VPS handles this)  

---

## Development Tools

| Tool | Purpose | Install |
|------|---------|---------|
| **Go 1.23** | Backend dev | `winget install GoLang.Go` |
| **Node.js 20** | Frontend (optional Vite) | `winget install OpenJS.NodeJS` |
| **Python 3.12** | AI service | `winget install Python.Python.3.12` |
| **PostgreSQL 17** | Database | `winget install PostgreSQL.PostgreSQL` |
| **Docker Desktop** | AI service container | `winget install Docker.DockerDesktop` |
| **Git** | Version control | `winget install Git.Git` |
| **VS Code** | Code editor | `winget install Microsoft.VisualStudioCode` |

---

## Tech Stack Summary

**Strengths:**
- ✅ Go backend = Fast, lightweight, single binary
- ✅ Vue.js 3 = Reactive UI, progressive enhancement
- ✅ PostgreSQL = Reliable, powerful
- ✅ Python AI = Best ML ecosystem
- ✅ ESP32 = Cheap IoT ($3/device)

**Trade-offs:**
- ⚠️ Multiple languages (Go, Python, JavaScript, C)
- ⚠️ Team needs diverse skills
- ✅ But: Each language is best-in-class for its role

**Alternatives Considered & Rejected:**
- ❌ Node.js backend → Too slow, high memory
- ❌ React frontend → Too heavy (130KB vs 40KB)
- ❌ MySQL → JSONB and array support weaker
- ❌ Arduino (not ESP32) → Less powerful, no WiFi

---

**Next Steps:**
1. Read [guides/SETUP.md](guides/SETUP.md) for installation
2. Read [ARCHITECTURE.md](ARCHITECTURE.md) for system design
3. Start coding with [implementation/](implementation/) guides
