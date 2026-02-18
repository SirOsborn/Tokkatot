# Tokkatot 2.0: Technology Stack Specification

**Document Version**: 2.0  
**Last Updated**: February 2026  
**Status**: Final Specification

---

## Technology Stack Overview

Tokkatot 2.0 uses modern, production-ready technologies selected for reliability, scalability, and team expertise.

---

## Frontend Technologies

### Framework: Vue.js 3

**Why Vue.js 3?**
- Small bundle size (~40KB gzipped) vs React (~130KB)
- Performance on low-end devices
- Composition API for better code organization
- Reactive data binding
- Large ecosystem
- Easy to learn for new developers

**Versions & Requirements**:
- **Vue.js**: 3.3+
- **Vite**: 4.0+ (build tool)
- **Node.js**: 16+ (development)
- **npm**: 8+

**Core Dependencies**:
```json
{
  "vue": "^3.3.0",
  "vue-router": "^4.2.0",
  "pinia": "^2.1.0",
  "axios": "^1.4.0",
  "socket.io-client": "^4.6.0",
  "chart.js": "^3.9.0",
  "tailwindcss": "^3.3.0",
  "vite": "^4.3.0"
}
```

### Styling: TailwindCSS

**Why Tailwind?**
- Utility-first CSS framework
- Small final CSS size (purging unused)
- Consistent design system
- Dark mode support built-in
- Responsive design helpers

**Configuration**:
```javascript
// tailwind.config.js
module.exports = {
  content: ['./index.html', './src/**/*.vue'],
  theme: {
    fontSize: {
      'xs': '12px',
      'sm': '14px',
      'base': '16px',
      'lg': '18px',
      'xl': '20px',
      '2xl': '24px',
      '3xl': '32px',
      '4xl': '40px',
      '5xl': '48px'
    },
    colors: {
      white: '#FFFFFF',
      black: '#000000',
      green: '#10B981',
      red: '#EF4444',
      yellow: '#F59E0B',
      gray: {
        50, 100, 200, 300, 400, 500, 600, 700, 800, 900
      }
    }
  }
}
```

### Charts: Chart.js

**Why Chart.js?**
- Lightweight (8KB gzipped)
- Responsive charts
- Multiple chart types (line, bar, gauge)
- Easy customization
- Good mobile support

**Supported Chart Types**:
- Line charts (sensor trends)
- Bar charts (daily totals)
- Gauge charts (current values)
- Pie charts (distribution, optional)

### State Management: Pinia

**Why Pinia?**
- Modern replacement for Vuex
- Smaller API surface
- Better TypeScript support
- Devtools integration
- No boilerplate

**Stores**:
```javascript
// stores/user.js
export const useUserStore = defineStore('user', () => {
  const user = ref(null)
  const farm = ref(null)
  
  const logout = () => { /* ... */ }
  
  return { user, farm, logout }
})
```

---

## Backend Technologies

### Runtime: Go (Go 1.21+)

**Why Go?**
- Fast compilation and startup
- Built-in concurrency (goroutines)
- Statically typed (catch errors early)
- Single binary deployment
- Great for API servers
- Growing cloud-native ecosystem

**Version**: 1.21+

### Framework: Fiber

**Why Fiber?**
- Lightweight web framework
- Similar to Express.js API
- Processing speed comparable to FastHTTP
- Easy middleware system
- Good documentation

**Core Dependencies**:
```go
import (
  "github.com/gofiber/fiber/v2"
  "github.com/gofiber/fiber/v2/middleware/cors"
  "github.com/gofiber/fiber/v2/middleware/logger"
  "github.com/gofiber/fiber/v2/middleware/compress"
  "github.com/gofiber/contrib/jwt"
)
```

### Authentication: JWT (jsonwebtoken)

**Implementation**:
```go
import "github.com/golang-jwt/jwt/v5"

// Token generation
token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
  "user_id": userId,
  "email": userEmail,
  "farm_id": farmId,
  "exp": time.Now().Add(24 * time.Hour).Unix(),
})
```

### Password Hashing: bcrypt

**Implementation**:
```go
import "golang.org/x/crypto/bcrypt"

// Hash password
hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)

// Verify password
err := bcrypt.CompareHashAndPassword(hash, []byte(password))
```

### Database Driver: pgx

**Why pgx?**
- High performance PostgreSQL driver
- Connection pooling
- Support for JSON/JSONB
- Prepared statements
- Native array support

### ORM (Optional): sqlc or GORM

**Recommendation**: Use sqlc for strongly-typed queries  
**Alternative**: Raw SQL with pgx (if team prefers)

### Time-Series Client: InfluxDB Client

```go
import ifluxdb "github.com/influxdata/influxdb-client-go/v2"
```

### Real-Time: Socket.io (go-socket.io)

```go
import "github.com/googollee/go-socket.io"
```

### Message Queue: Redis (Optional)

```go
import "github.com/redis/go-redis/v9"
```

### Configuration: Viper

```go
import "github.com/spf13/viper"
```

---

## Database Technologies

### Primary: PostgreSQL 15+

**Why PostgreSQL?**
- ACID compliance (transactions)
- Advanced features (JSON, arrays, full-text search)
- Perfect for relational data
- Excellent performance
- Mature and stable

**Extensions**:
- `pgcrypto` (encryption functions)
- `uuid-ossp` (UUID generation)
- `pg_trgm` (similarity searches)

**Version**: 15+

### Time-Series: InfluxDB

**Why InfluxDB?**
- Purpose-built for time-series metrics
- Excellent compression
- Built-in downsampling
- Fast queries on large datasets
- 5-year data retention ready

**Version**: 2.7+

### Cache: Redis

**Why Redis?**
- In-memory speed
- Support for data structures
- Pub/Sub for real-time
- TTL-based expiration
- Session storage

**Version**: 7.0+

### Object Storage: S3-Compatible

**Implementation**: DigitalOcean Spaces or AWS S3

**Client Library**: 
```go
import "github.com/aws/aws-sdk-go/aws/s3"
```

---

## IoT & Device Layer

### Embedded OS: ESP-IDF

**ESP32 Development Framework**:
- Official Espressif development framework
- C/C++ based
- Real-time operating system (FreeRTOS)
- WiFi and Bluetooth support
- OTA update capability

**Version**: 4.4+

**Core Components**:
```c
#include "esp_wifi.h"     // WiFi management
#include "esp_mqtt.h"     // MQTT client
#include "driver/gpio.h"  // GPIO control
#include "driver/adc.h"   // ADC for sensors
#include "driver/pwm.h"   // PWM for dimmers
#include "nvs.h"          // Non-volatile storage
#include "spiffs.h"       // File system
#include "esp_ota_ops.h"  // OTA updates
#include "esp_tls.h"      // TLS security
```

### Message Protocol: MQTT

**Why MQTT?**
- Lightweight (perfect for IoT)
- Publish/subscribe model
- Quality of Service (QoS) levels
- Last Will and Testament (offline detection)
- Still works over 2G/3G/4G

**Implementation**:
- **ESP32**: esp-mqtt library
- **Raspberry Pi**: mosquitto-clients or paho-mqtt
- **Cloud**: Eclipse Mosquitto or managed MQTT

**Topics**:
```
farm/{farmId}/devices/{deviceId}/
├── command        # Cloud → Device
├── status         # Device → Cloud
├── sensor         # Periodic readings
├── heartbeat      # Keep-alive
└── fw_update      # OTA firmware
```

### Local Hub: Raspberry Pi 4B

**Specifications**:
- Quad-core Cortex-A72 (1.5 GHz)
- 4GB RAM
- 32GB microSD card
- Gigabit Ethernet
- Dual-band WiFi

**Operating System**: Raspberry Pi OS (lite or full)

**Local Agent Runtime**: Python, Go, or Node.js

---

## DevOps & Infrastructure

### Containerization: Docker

**Dockerfiles**: Alpine-based for smaller sizes

**Example Backend Dockerfile**:
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o backend .

FROM alpine:3.18
COPY --from=builder /app/backend /app/backend
EXPOSE 8080
CMD ["/app/backend"]
```

**Image Size Target**: < 100MB per service

### Orchestration: Docker Compose (MVP) or Kubernetes (Scale)

**Recommendation**: Start with Docker Compose, migrate to Kubernetes as load increases

### CI/CD: GitHub Actions

**Pre-commit Checks**:
- Linting (golangci-lint, ESLint)
- Unit tests
- SAST (Semgrep, SonarQube)
- Dependency scanning

**Build & Deploy**:
- Build Docker images
- Push to container registry
- Deploy to staging (automated)
- Deploy to production (manual approval)

---

## Monitoring & Logging

### Metrics: Prometheus

**Scrape Interval**: 15 seconds

**Key Metrics**:
- API request duration (histogram)
- Request count (counter)
- Error rate (counter)
- Active connections (gauge)
- Database query performance

### Visualization: Grafana

**Dashboards**:
- System health (CPU, memory, disk)
- API performance
- Database metrics
- Business metrics (active farms, users)

### Logging: ELK Stack or Loki

**Log Levels**: DEBUG, INFO, WARN, ERROR, CRITICAL

**Format**: JSON for easy parsing
```json
{
  "timestamp": "2026-02-18T14:30:00Z",
  "level": "ERROR",
  "service": "device-service",
  "message": "Failed to send command",
  "device_id": "esp32-001",
  "error": "timeout"
}
```

---

## Development Tools

### Version Control: Git & GitHub

- Hosting: GitHub
- Branch strategy: Git Flow or GitHub Flow
- Code review: Pull Requests (mandatory 2 approvals)
- CI/CD: GitHub Actions

### API Documentation: OpenAPI/Swagger

- Framework: Swagger/OpenAPI 3.0
- Generation: From code annotations
- UI: Swagger UI for testing

### Environment Management: .env files

**Development**: `.env.local` (never commit)  
**Staging**: `.env.staging` (in git, encrypted values)  
**Production**: GitHub Secrets  

### Task Runner: Make

**Makefile targets**:
```makefile
make build       # Build all services
make test        # Run tests
make lint        # Lint code
make start       # Docker compose up
make stop        # Docker compose down
make migrate     # Run database migrations
```

---

## Performance Benchmarks

### API Response Times

| Endpoint | Target | Notes |
|----------|--------|-------|
| POST /login | < 200ms | Database + crypto |
| GET /devices | < 100ms | Cached query |
| POST /commands | < 50ms | Queue operation |
| WebSocket connect | < 300ms | TLS handshake |

### Database Performance

| Query | Target | Data Size |
|-------|--------|-----------|
| Get user farms | < 50ms | 100 farms |
| Query sensor data (24h) | < 500ms | 1M data points |
| Get device state | < 20ms | Indexed, cached |
| Insert sensor reading | < 10ms | Batch write |

### Frontend Performance

| Metric | Target | Tolerance |
|--------|--------|-----------|
| Page load | < 2s | 95th percentile |
| Re-render | < 60ms | React/Vue |
| Chart render | < 500ms | 10K data points |
| Offline sync | < 30s | After internet |

---

## Dependency Management

### Vulnerability Scanning

- **Tool**: Snyk, Dependabot, npm audit
- **Frequency**: On every commit
- **Policy**: No known CVEs in production

### Dependency Updates

- **Patch updates** (1.0.0 → 1.0.1): Automatic
- **Minor updates** (1.0.0 → 1.1.0): Weekly review
- **Major updates** (1.0.0 → 2.0.0): Manual with testing

### Lock Files

- **Frontend**: package-lock.json (committed)
- **Backend**: go.sum (committed)
- **Ensures**: Reproducible builds

---

## Recommended IDE & Tools

**Development**:
- **IDE**: VS Code (with Go, Vue extensions)
- **Database GUI**: pgAdmin, DBeaver
- **API Testing**: Postman, Insomnia
- **Database Migration**: Flyway, Liquibase

**Deployment**:
- **Kubernetes**: kubectl, Helm charts
- **IaC**: Terraform or Pulumi
- **Container Registry**: DigitalOcean Registry (DOCR)

---

## Version Lock & Release Strategy

**Frontend** (Vue.js):
- Vue.js: 3.3+ (LTS support)
- Compatible dependencies: pinned to major.minor
- Breaking changes: quarterly reviews

**Backend** (Go):
- Go: 1.21+ (latest stable)
- Supported versions: current and previous
- Module versioning: SemVer

**Release Cycle**:
- **Security patches**: ASAP (hot-fix)
- **Bug fixes**: Weekly releases
- **Features**: Bi-weekly releases
- **Major versions**: Quarterly planning

---

**Version History**
| Version | Date | Changes |
|---------|------|---------|
| 2.0 | Feb 2026 | Initial technology stack specification |

**Related Documents**
- SPECIFICATIONS_ARCHITECTURE.md
- SPECIFICATIONS_DEPLOYMENT.md
- SPECIFICATIONS_EMBEDDED.md
