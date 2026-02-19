# Tokkatot 2.0: System Architecture Specification

**Document Version**: 2.0  
**Last Updated**: February 2026  
**Status**: Final Specification  

---

## Overview

Tokkatot 2.0 is a cloud-based smart farming IoT system with edge computing fallback. The system enables remote control and monitoring of farm equipment (water pumps, feeders, lights, fans, heaters, conveyors) with sophisticated scheduling, real-time dashboards, and comprehensive data logging.

### Key Architectural Features
- ✅ **3-Tier Architecture** - Separation of client, API, and data layers
- ✅ **Edge Computing** - Local Raspberry Pi agent for offline operations
- ✅ **Real-Time Communication** - MQTT for devices, WebSocket for clients
- ✅ **Event-Driven** - Device state changes trigger system events
- ✅ **Stateless Services** - Horizontally scalable backend
- ✅ **Offline-First** - Queue-based synchronization when internet returns
- ✅ **Production-Ready** - Enterprise error handling, logging, monitoring

---

## 3-Tier Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                     CLIENT LAYER (User Devices)                 │
│                                                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐         │
│  │  Web Browser │  │  Mobile App  │  │ API Clients  │         │
│  │  (Vue.js 3)  │  │ (React Native)  │ (CLI/Tools)  │         │
│  └──────────────┘  └──────────────┘  └──────────────┘         │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
                     ↕ HTTPS/WSS Secure Connection
┌─────────────────────────────────────────────────────────────────┐
│                   API GATEWAY LAYER (Cloud)                     │
│                                                                 │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │  Load Balancer / Reverse Proxy (Nginx/HAProxy)           │ │
│  │  • Request routing                                        │ │
│  │  • SSL/TLS termination                                   │ │
│  │  • Rate limiting & Auth validation                       │ │
│  │  • Request/Response logging                              │ │
│  └───────────────────────────────────────────────────────────┘ │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
                     ↕ Internal Network/API
┌─────────────────────────────────────────────────────────────────┐
│              APPLICATION LAYER (Microservices)                  │
│                                                                 │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────┐ │
│  │  Auth Service    │  │  Device Service  │  │ Schedule Svc │ │
│  │  • JWT tokens    │  │  • Device state  │  │ • Cron/Tasks │ │
│  │  • User mgmt     │  │  • Life cycle    │  │ • Conditions │ │
│  │  • Permissions   │  │  • Commands      │  │ • History    │ │
│  └──────────────────┘  └──────────────────┘  └──────────────┘ │
│                                                                 │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────┐ │
│  │  Data Service    │  │  Notification    │  │  Config Svc  │ │
│  │  • History logs  │  │  • Alert engine  │  │ • Settings   │ │
│  │  • Queries       │  │  • Push notifs   │  │ • Profiles   │ │
│  │  • Analytics     │  │  • Email alerts  │  │ • Defaults   │ │
│  └──────────────────┘  └──────────────────┘  └──────────────┘ │
│                                                                 │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │  WebSocket Server (Real-Time Streaming)                 │ │
│  │  • Push updates to connected clients                     │ │
│  │  • Bi-directional communication                          │ │
│  │  • Connection management & heartbeat                     │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
                     ↕ Internal API / Message Queue
┌─────────────────────────────────────────────────────────────────┐
│                   DATA LAYER (Persistence)                      │
│                                                                 │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────┐ │
│  │  Primary DB      │  │  Time-Series DB  │  │ Cache Layer  │ │
│  │  (PostgreSQL)    │  │  (InfluxDB)      │  │ (Redis)      │ │
│  │  • Users         │  │  • Sensor data   │  │ • Sessions   │ │
│  │  • Devices       │  │  • Temperature   │  │ • Cache      │ │
│  │  • Farms         │  │  • Humidity      │  │ • Queues     │ │
│  │  • Configs       │  │  • Metrics       │  │              │ │
│  │  • Event logs    │  │                  │  │              │ │
│  └──────────────────┘  └──────────────────┘  └──────────────┘ │
│                                                                 │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │  File Storage (S3 Compatible Object Storage)            │ │
│  │  • Firmware files (ESP32 binaries)                      │ │
│  │  • Device configurations                                │ │
│  │  • Backup files                                         │ │
│  │  • Media/Images (farm photos, etc)                     │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
                     ↕ MQTT / HTTP API
┌─────────────────────────────────────────────────────────────────┐
│                   EDGE LAYER (On-Farm)                          │
│                                                                 │
│  ┌────────────────────────────────────────────────────────────┐│
│  │  Local Hub: Raspberry Pi 4B                              ││
│  │  ┌──────────────────────────────────────────────────────┐ ││
│  │  │ Local Agent Service                                  │ ││
│  │  │ • MQTT broker / client                              │ ││
│  │  │ • Device communication                              │ ││
│  │  │ • Offline command queueing                          │ ││
│  │  │ • Cloud sync coordination                           │ ││
│  │  │ • Fallback scheduling                               │ ││
│  │  │ • Status/health reporting                           │ ││
│  │  └──────────────────────────────────────────────────────┘ ││
│  └────────────────────────────────────────────────────────────┘│
│                                                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐┌──────────┐│
│  │ ESP32 Node 1│  │ ESP32 Node 2│  │ESP32 Node N││ Sensors  ││
│  └─────────────┘  └─────────────┘  └─────────────┘└──────────┘│
│  • Water Pump   • Feeder         • Heater        • DHT22    │
│  • Light        • Conveyor       • Fan           • Relay    │
│  • Motion       • Weight Scale   • Alarm         • Buttons  │
│  • Status LEDs  • Emergency Stop • Servo         • Displays │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## Architecture Components

### 1. Client Layer

**Purpose**: User interface and device interaction  
**Technologies**: Vue.js 3, React Native, Web APIs  
**Deployment**: Browser, Mobile App Store, Local Storage

**Components:**
- **Web Application** (Vue.js 3)
  - Responsive design for phones, tablets, desktops
  - Facebook-style navigation (bottom tab bar)
  - Home page as primary landing
  - Real-time dashboard with WebSocket updates
  - Scheduling interface
  - Offline support via Service Workers
  
- **Mobile App** (React Native - optional for future)
  - Native iOS/Android implementation
  - In-app notifications (no push)
  - Home screen widgets
  - Offline queueing
  
- **API Clients** (CLI, Tools)
  - For developers and system integration
  - OpenAPI/REST-based

**Performance Targets:**
- Web app bundle: < 150KB (gzipped)
- Initial load: < 2 seconds on 4G
- React re-render: < 60ms
- WebSocket reconnection: < 3 seconds

---

### 2. API Gateway Layer

**Purpose**: Request routing, security, rate limiting, load distribution  
**Deployment**: Cloud provider (DigitalOcean, AWS)  
**Technologies**: Nginx, HAProxy, Kong

**Responsibilities:**
- TLS/SSL termination (HTTPS, WSS)
- Request routing to appropriate services
- Rate limiting (API abuse prevention)
- Request/response logging
- Authentication validation
- Compression (gzip, brotli)
- Caching headers management
- CORS policy enforcement
- DDoS protection

**Constraints:**
- Max request size: 10MB
- Max timeout: 30 seconds
- Rate limit: 1000 req/min per IP (default), 10 req/min unauthenticated
- Connection pooling: 1000 concurrent connections

---

### 3. Application Layer

**Purpose**: Business logic, data processing, real-time communication  
**Technologies**: Go (Fiber), Node.js Express, Python FastAPI  
**Deployment**: Docker containers (Kubernetes or Docker Compose)

**Microservices:**

**A. Authentication Service (Farmer-Centric)**
- User registration (email OR phone number, not both required)
- Simple login with JWT tokens (no MFA for farmers)
- Password reset via email or SMS
- Simplified role-based access (Owner, Manager, Viewer only)
- Multi-device session management
- Activity audit logging

**B. Device Service (Tokkatot Team Managed)**
- Device lifecycle management (added by Tokkatot team only)
- Device state tracking (online/offline/error)
- Command execution to devices (via MQTT)
- Firmware version management and OTA updates
- Device grouping and organization
- Health monitoring and diagnostics
- Connection status broadcasts

**C. Schedule Service**
- CRUD operations for schedules
- Cron expression parsing and execution
- Duration-based cycle management
- Condition-based trigger evaluation
- Schedule versioning and rollback
- Execution history and audit
- Disabled schedule management
- Multi-device coordination

**D. Data Service**
- Historical data retrieval
- Time-range queries (24h, 7d, 30d, custom)
- Data aggregation (average, min, max, sum)
- Anomaly detection
- Trend analysis
- Export functionality (CSV, JSON)
- Cache management

**E. Notification/Message Service**
- **In-app message log** - all alerts visible on dashboard message center
- Real-time dashboard updates via WebSocket push
- Message severity: urgent, important, info
- Read/unread status tracking
- Metadata attached: device info, sensor values, timestamps
- Message retention: 90 days in dashboard, 2 years in database
- No external notifications (no SMS, email, push)

**F. Configuration Service**
- Farm settings management
- User preferences
- Notification thresholds
- Default values
- Feature flags
- System configuration

**G. WebSocket Server**
- Real-time subscriptions
- Live data streaming
- Connection lifecycle management
- Heartbeat/keepalive
- Message queuing for offline clients
- Broadcast to multiple clients

---

### 4. Data Layer

**Purpose**: Persistent storage, caching, object storage  
**Deployment**: Cloud services (managed databases)

**Databases and Storage:**

**A. Primary Database (PostgreSQL)**
- Users and authentication
- Farms and devices
- Schedules and automation rules
- Event audit logs
- Notifications and alerts
- Configuration and settings
- User profiles and preferences

**B. Time-Series Database (InfluxDB)**
- Sensor readings (temperature, humidity, etc)
- Device state timelines
- System metrics
- Performance data
- Retention: 5 years of data

**C. Cache Layer (Redis)**
- Session storage (JWT auth state)
- Rate limiting counters
- Real-time state caching
- Message queuing
- Pub/sub for real-time updates
- TTL-based expiration

**D. Object Storage (S3-Compatible)**
- Firmware files for OTA updates
- Device configuration backups
- Media files (farm photos)
- Log archives
- Encrypted backups

---

### 5. Edge/Local Layer

**Purpose**: Local autonomy, offline operations, device communication  
**Deployment**: On-farm hardware (Raspberry Pi 4B, ESP32)

**A. Local Hub (Raspberry Pi 4B)**
- Runs local agent service
- Network connectivity (WiFi/Ethernet)
- MQTT broker/client
- Offline command queue
- Local scheduling execution
- Health monitoring
- Cloud synchronization coordinator

**B. Device Nodes (ESP32)**
- Hardware controllers (relays, PWM, ADC)
- Sensor data collection
- Real-time feedback loops
- Command execution
- Status LED indicators
- Battery backup (optional)
- Over-the-Air (OTA) firmware updates

---

## Data Flow Patterns

### Pattern 1: Device Control Command Flow

```
User App
  ↓
[Client sends HTTPS POST to /api/devices/[ID]/commands]
  ↓
API Gateway (Rate limit, Auth validate)
  ↓
Device Service
  ↓
[Publish to MQTT: farm/devices/[ID]/command]
  ↓
Local Hub (RPi)
  ↓
[Route to target ESP32 or execute locally]
  ↓
ESP32 Device
  ↓
[Execute hardware command]
  ↓
[Publish state change to MQTT: farm/devices/[ID]/status]
  ↓
Local Hub (capture state change)
  ↓
Device Service (via API/webhook)
  ↓
[Update DB, broadcast via WebSocket to all clients]
  ↓
User App (real-time notification)
```

**Timeline:**
- User to gateway: ~20ms
- Gateway to service: ~5ms
- MQTT publish: ~10ms
- Device execution: varies (100ms-5s)
- Status update broadcast: ~50ms
- Client receives update: ~100ms (if connected)
- **Total latency: 300-500ms** (offline) or **100-200ms** (online)

---

### Pattern 2: Sensor Data Collection

```
ESP32 Device
  ↓
[Read sensor value every 30 seconds]
  ↓
[Publish to MQTT: farm/sensors/[ID]/data]
  ↓
Local Hub (aggregate/buffer)
  ↓
[Every 5 minutes, send to cloud API]
  ↓
API Gateway
  ↓
Data Service
  ↓
[Write to TimeSeries DB (InfluxDB)]
  ↓
[Calculate rolling averages]
  ↓
[Cache latest value in Redis]
  ↓
[Publish to WebSocket subscribers]
  ↓
User App (real-time dashboard update)
```

**Data retention:**
- Detailed (1-minute): 30 days
- Aggregated (hourly): 1 year
- Summarized (daily): 5 years

---

### Pattern 3: Schedule Execution

```
Schedule
  ↓
[Next execution time calculated]
  ↓
Cron triggers at scheduled time
  ↓
Schedule Service:
  - Evaluate conditions (temperature, humidity, etc)
  - Prepare command payload
  - Publish to device queue
  ↓
[Send via MQTT to Local Hub]
  ↓ (if cloud available) or (queue locally if offline)
  ↓
Local Hub
  ↓
[Execute command on ESP32]
  ↓
[Record execution and result]
  ↓
[Sync results to cloud when available]
  ↓
User logs updated
```

**Execution guarantee:**
- Expected deviation: ±1 second (from scheduled time)
- Retry: 3 times (if device offline)
- Fallback: Local RPi maintains schedule if cloud fails

---

### Pattern 4: Real-Time Sync (Offline → Online)

```
User (Offline)
  ↓
[Commands queued locally in Service Worker]
  ↓
[Device receives commands from local RPi]
  ↓
[Status updates cached locally]
  ↓
  
Internet connection returns
  ↓
[Client connects to cloud again]
  ↓
Sync Service:
  - Send: [list of pending commands]
  - Receive: [latest cloud state]
  - Merge: device state is source of truth
  - Update: local cache with cloud data
  ↓
[Verify no conflicts]
  ↓
[Broadcast merged state to all clients]
  ↓
Dashboard updated
```

---

## Integration Points

### 1. Device ↔ Cloud Integration

**Protocol**: MQTT + REST API  
**Security**: TLS 1.3, certificate pinning  
**Payload Format**: JSON  

**Topics:**
```
farm/{farmId}/devices/{deviceId}/
├── command        # Cloud → Device (commands)
├── status         # Device → Cloud (state changes)
├── sensor         # Device → Cloud (sensor readings)
├── heartbeat      # Device → Cloud (keep-alive)
└── fw_update      # Cloud → Device (OTA updates)
```

---

### 2. Client ↔ Cloud Integration

**Protocols**: 
- REST: For data queries and state mutations
- WebSocket: For real-time updates
- HTTP Polling: Fallback for poor 4G connections (3-second intervals)

**Authentication**: JWT in Authorization header  
**Rate Limiting**: Token bucket algorithm  
**Compression**: gzip for responses > 1KB  

---

### 3. Service-to-Service Integration

**Protocol**: Internal REST API + gRPC (optional)  
**Discovery**: Service registry (Consul, Kubernetes DNS)  
**Circuit Breaker**: Prevent cascading failures  
**Retry Logic**: Exponential backoff  

---

## Constraints & Requirements

### Network Constraints

- **Expected bandwidth**: 1-2 Mbps per device
- **Latency tolerance**: up to 5 seconds (schedule execution still works)
- **Polling interval**: 30 seconds (device heartbeat)
- **Offline duration**: support 72+ hours
- **Connection type**: 4G LTE (primary), WiFi (fallback)

### Performance Targets

| Metric | Target | Tolerance |
|--------|--------|-----------|
| Page load | < 2 seconds | ± 500ms |
| Dashboard refresh | < 1 second | ± 200ms |
| Command latency | < 500ms | ± 100ms |
| Schedule deviation | ±1 second | Max 5 seconds |
| API responsiveness | < 200ms | 95th percentile |
| Uptime | 99.5% | Max 3.6 hours downtime/month |
| Data sync | < 30 seconds | After internet returns |

### Scalability Requirements

- **Current**: 10-50 farms
- **Target Year 1**: 100-500 farms
- **Target Year 2**: 1000+ farms
- **Devices per farm**: 10-50 devices
- **Concurrent users**: 2-5 per farm
- **API calls per day**: ~100K-500K

### Security Requirements (See IG_SPECIFICATIONS_SECURITY.md)

- TLS 1.3+ for all communications
- JWT tokens with 24-hour expiration (no MFA for farmers)
- Simplified role system (Owner, Manager, Viewer)
- Rate limiting and DDoS protection
- Encrypted storage of sensitive data
- Audit logs for all actions
- Phone/Email registration support

---

## Fallback & Recovery Mechanisms

### Internet Connectivity Loss

1. **Local mode activated**
   - RPi becomes primary command executor
   - Commands queued locally
   - Real-time communication switches to local only
   
2. **Schedules continue**
   - Local agent processes schedules
   - No cloud coordination needed
   
3. **User app continues**
   - Read-only access to cached data
   - Commands queued for later sync
   - Service Worker enables offline viewing

### Service Failure

1. **API service crashes**
   - Load balancer detects failure
   - Routes to healthy instances
   - If all instances fail: return 503 error
   
2. **Database failure**
   - Primary fails: automatic failover to replica
   - Cache layer (Redis) provides temporary access
   - Data loss prevention via write-ahead logs

3. **WebSocket connection failure**
   - Automatic fallback to HTTP polling (3 second interval)
   - Reconnection with exponential backoff
   - Message queue prevents lost updates

### Device Firmware Crash

1. **Watchdog timer** triggers automatic restart
2. **Boot validation** checks firmware integrity
3. **If boot fails**: Fall back to previous firmware version
4. **Alert**: Notify cloud of restart event
5. **Retry OTA**: Attempt update again next maintenance window

---

## Architecture Decision Records

| Decision | Rationale | Alternative Considered |
|----------|-----------|----------------------|
| **Vue.js 3** (frontend) | Performance on low-end phones, small bundle size | React (larger), Angular (complex) |
| **Go Fiber** (backend) | Concurrent, simple deployment, fast startup | Node.js (slower startup), Python (slower) |
| **PostgreSQL** (primary DB) | ACID compliance, JSON support, proven at scale | MySQL (less features), SQLite (not scalable) |
| **InfluxDB** (time-series) | Purpose-built for metrics, excellent compression | TimescaleDB (more complex), Prometheus (pull-only) |
| **MQTT** (device protocol) | Lightweight, publish/subscribe, offline support | HTTP (heavier), CoAP (less reliable), custom TCP |
| **Docker** (containerization) | Industry standard, reproducible, easy deployment | VMs (heavier), binaries (less reliable) |
| **DigitalOcean** (cloud) | Cost-effective, simple, good for small teams | AWS (expensive, complex), GCP, Linode |
| **Redis** (caching) | Fast, supports expiration, pub/sub, simple | Memcached (no expiration), Hazelcast (complex) |

---

## Migration Path from v1.0

| Component | v1.0 | v2.0 | Migration |
|-----------|------|------|-----------|
| **Storage** | SQLite | PostgreSQL | Schema mapping, data import script |
| **API** | RESTful | RESTful + GraphQL | New API, v1 deprecated |
| **Real-time** | Polling | MQTT + WebSocket | Gradual rollout |
| **Auth** | Sessions | JWT | Token migration |
| **Devices** | Manual config | Auto-discovery | Backward compatible |
| **Firmware** | Manual USB flash | OTA | Phased rollout |

---

## Key Files & Documentation

- **IG_SPECIFICATIONS_DATABASE.md** - Database schema details
- **IG_SPECIFICATIONS_API.md** - API endpoints and contracts (8 auth, 5 user mgmt, 8 farm mgmt, 9 device mgmt, 8 control, 7 scheduling, 8 alerts, 5 reporting = 58 total)
- **IG_SPECIFICATIONS_FRONTEND.md** - UI/UX for farmers (48px+ fonts, WCAG AAA, Khmer/English)
- **IG_SPECIFICATIONS_EMBEDDED.md** - Device firmware architecture (Tokkatot team manages device setup)
- **IG_SPECIFICATIONS_SECURITY.md** - Security with simplified roles (no complex RBAC)
- **OG_SPECIFICATIONS_DEPLOYMENT.md** - Infrastructure setup
- **IG_TOKKATOT_2.0_FARMER_CENTRIC_SPECIFICATIONS.md** - Phone/Email registration, accessibility for elderly farmers

---

**Version History**
| Version | Date | Changes |
|---------|------|---------|
| 2.0-FarmerCentric | Feb 2026 | Simplified for elderly farmers with low digital literacy |
| | | Phone/Email registration, 3 simple roles, device setup by team |
| 2.0 | Feb 2026 | Initial production specification |

**Next Steps**: Review with tech team, finalize technology selections, begin Phase 1 implementation.
