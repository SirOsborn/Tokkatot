# 🤖 AI Context: Go API Gateway & Middleware

**Directory**: `middleware/`  
**Your Role**: HTTP API, authentication, authorization, device management, real-time communication  
**Tech Stack**: Go 1.19+, PostgreSQL, MQTT, WebSocket  

---

## 🎯 What You're Building

**RESTful API Gateway** (Port 6060)
- **Authentication**: Email/Phone login, JWT tokens (24h expiry)
- **Authorization**: 3-role system (Owner, Manager, Viewer)
- **Device Management**: Control IoT devices (water pumps, lights, feeders, fans, heaters, conveyors)
- **Scheduling**: Create cron-based automation rules
- **Monitoring**: Real-time sensor data, alerts, event logs
- **AI Integration**: Proxy predictions to FastAPI service (Port 8000)

**Key Statistics**:
- 66 REST endpoints documented in `IG_SPECIFICATIONS_API.md`
- JWT token validation on all protected endpoints
- Device state source of truth (syncs with local RPi agent via MQTT)
- Real-time WebSocket for live updates to connected clients

---

## 📁 File Structure

```
middleware/
├── main.go               # HTTP server setup, port 6060
├── go.mod               # Go module definition
├── go.sum               # Dependency checksums
├── .env                 # Database URL, secret keys (GITIGNORE'D)
├── api/
│   ├── authentication.go   # Login, signup, password reset, token refresh
│   ├── data-handler.go     # Historical data, reports, analytics
│   ├── disease-detection.go # AI service proxy endpoints
│   ├── profiles.go         # User profile management
│   └── ... (other endpoints)
├── database/
│   └── sqlite3_db.go    # PostgreSQL queries, migrations
└── utils/
    └── utils.go         # Helper functions (JWT signing, hashing, etc)
```

---

## 🚀 Getting Started

### Local Development

```bash
cd middleware

# Set up environment
echo "DB_CONNECTION_STRING=postgres://user:pass@localhost:5432/tokkatot" > .env
echo "JWT_SECRET_KEY=your-secret-key-here" >> .env

# Install dependencies
go mod download

# Build
go build -o tokkatot.exe .

# Run
./tokkatot.exe
# Server starts on http://localhost:6060
```

### Testing

```bash
# All tests
go test ./...

# Specific test
go test ./api -v

# Test with coverage
go test -cover ./...
```

---

## 🔐 Authentication & Authorization

### JWT Token Flow

1. **User Login** (`POST /auth/login`)
   - Accept: email/phone + password + device_name
   - Validate credentials against database
   - Generate JWT tokens:
     - Access token (24 hours expiry)
     - Refresh token (30 days expiry)
   - Return both tokens

2. **Protected Endpoints**
   - Extract token from `Authorization: Bearer <token>` header
   - Validate signature and expiry
   - Extract user_id, farm_ids, role from token
   - Check role has permission for resource

3. **Token Refresh** (`POST /auth/refresh`)
   - Accept refresh_token
   - Issue new access_token + refresh_token

### Role System (Farmer-Centric)

| Role | Permissions | Can... | Cannot... |
|------|-------------|--------|-----------|
| **Owner** | Full farm access | Create users, manage devices, view all data | Cannot delete farm |
| **Manager** | Device control + delegation | Control devices, schedule tasks, invite keepers | Cannot invite other managers |
| **Viewer** | Read-only | View monitoring data, see alerts | Cannot make changes |

**Important**: Farmers cannot add/remove devices. Only Tokkatot team can add devices via admin-only endpoints.

---

## 📊 Core Endpoints

### Authentication (8 endpoints)

```go
POST /auth/login                    // Email OR phone
POST /auth/signup                   // Register new user
POST /auth/logout                   // Invalidate session
POST /auth/refresh                  // New tokens
POST /auth/verify-email             // Email verification
POST /auth/forgot-password          // Initiate reset
POST /auth/reset-password           // Complete reset
POST /auth/change-password          // User changes own password
```

### Device Management (10 endpoints)

```go
GET  /farms/{farm_id}/devices               // List all
GET  /farms/{farm_id}/devices/{device_id}   // Single device
POST /admin/devices                         // Add (Tokkatot team only)
POST /farms/{farm_id}/devices/{device_id}/commands  // Control
DELETE /admin/devices/{device_id}           // Remove (Tokkatot team only)
GET  /farms/{farm_id}/devices/{device_id}/state     // Current state
```

### AI/Disease Detection (3 endpoints)

```go
GET  /api/ai/health                 // AI service health check
POST /api/ai/predict                // Proxy prediction request
POST /api/ai/predict/detailed       // Proxy detailed prediction
```

### Scheduling (7 endpoints)

```go
GET    /farms/{farm_id}/schedules           // List
POST   /farms/{farm_id}/schedules           // Create
PUT    /farms/{farm_id}/schedules/{id}      // Update
DELETE /farms/{farm_id}/schedules/{id}      // Delete
GET    /farms/{farm_id}/schedules/{id}/history // View executions
```

### Analytics & Reporting (5 endpoints)

```go
GET /farms/{farm_id}/reports/device-metrics
GET /farms/{farm_id}/reports/device-usage
GET /farms/{farm_id}/reports/farm-performance
GET /farms/{farm_id}/reports/export
GET /farms/{farm_id}/events
```

---

## 🛠️ Key Functions

### `main.go`
- `main()` - Start HTTP server, load config, initialize database
- `setupRoutes()` - Register all endpoint handlers
- `setupMiddleware()` - CORS, logging, rate limiting

### `api/authentication.go`
```go
func Login(w http.ResponseWriter, r *http.Request)        // POST /auth/login
func Signup(w http.ResponseWriter, r *http.Request)       // POST /auth/signup
func RefreshToken(w http.ResponseWriter, r *http.Request) // POST /auth/refresh
func ValidateJWT() MiddlewareFunc                          // JWT validation
```

### `api/disease-detection.go`
```go
func ProxyPrediction(w http.ResponseWriter, r *http.Request)        // POST /api/ai/predict
func ProxyDetailedPrediction(w http.ResponseWriter, r *http.Request) // POST /api/ai/predict/detailed
func ProxyHealthCheck(w http.ResponseWriter, r *http.Request)        // GET /api/ai/health
```

### `database/sqlite3_db.go`
```go
func GetUser(email string) (*User, error)
func CreateUser(user *User) error
func GetDevices(farmID string) ([]*Device, error)
func SavePrediction(prediction *Prediction) error
func LogEvent(event *Event) error
```

---

## 📝 Code Guidelines

### ✅ DO:
- Use PostgreSQL with parameterized queries (prepared statements)
- Validate all user input (size, type, format)
- Hash passwords with bcrypt
- Check JWT token + role on every protected endpoint
- Log all device commands (audit trail)
- Return proper HTTP status codes (200, 201, 400, 401, 403, 404, 500)
- Use `.env` for secrets (DB password, JWT secret)
- Implement error handling with meaningful messages

### ❌ DON'T:
- Build SQL strings with string concatenation (SQL injection risk!)
- Store passwords in plaintext
- Trust headers without validation
- Return internal error details to clients
- Hardcode database URLs or API keys
- Allow farmers to add/remove devices (team only)
- Expose sensitive user information (phone numbers) in responses

---

## 🔒 Security Checklist

- ✅ Passwords hashed with bcrypt (cost ≥ 12)
- ✅ JWT tokens signed with secret key
- ✅ HTTPS enforced (reverse proxy handles TLS)
- ✅ CORS restricted to frontend domain only
- ✅ Rate limiting (10 auth/min, 100 device control/min per user)
- ✅ SQL injection prevention (parameterized queries)
- ✅ CSRF protection on state-changing operations
- ✅ Input validation on all endpoints
- ✅ Access control checks on farm operations
- ✅ Secrets in `.env` (never in code)

---

## 💾 Database Schema

**Key Tables** (See `IG_SPECIFICATIONS_DATABASE.md` for full schema):

```
users
  - id (UUID)
  - email (unique, optional if phone provided)
  - phone (unique, optional if email provided)
  - password_hash
  - role (Owner/Manager/Viewer)
  - created_at

farms
  - id (UUID)
  - owner_id (FK users)
  - name
  - location
  - created_at

devices
  - id (UUID)
  - farm_id (FK farms)
  - name
  - device_type (Pump, Light, Feeder, etc)
  - status (online/offline/error)
  - last_heartbeat
  - added_by (team member, not farmer)

commands
  - id (UUID)
  - device_id (FK devices)
  - command_type
  - parameters
  - executed_at
  - result

sensor_readings (InfluxDB)
  - device_id
  - sensor_type (temperature, humidity)
  - value
  - timestamp

predictions
  - id (UUID)
  - user_id (FK users)
  - device_id (FK devices, optional)
  - disease
  - confidence
  - image_hash
  - created_at

event_logs
  - id (UUID)
  - farm_id (FK farms)
  - event_type
  - message
  - timestamp
```

---

## 🔗 Integration Points

### With AI Service (FastAPI, Port 8000)

```go
// When farmer submits image for disease detection:
1. Receive image upload at POST /api/ai/predict
2. Validate image (size, format)
3. Forward to FastAPI: POST http://localhost:8000/predict
4. Get response: {disease, confidence, recommendation}
5. Store in PostgreSQL predictions table
6. Broadcast to connected WebSocket clients
7. Return JSON to user mobile app
```

### With Local Hub (Raspberry Pi, MQTT)

```go
// When farmer sends device command:
1. Receive command at POST /devices/{id}/commands
2. Validate user has permission
3. Publish to MQTT: farm/{farm_id}/devices/{device_id}/command
4. Local RPi receives and executes on ESP32
5. ESP32 responds with status
6. Local RPi sends status back via MQTT
7. Go API receives status and broadcasts via WebSocket
```

### With Frontend (WebSocket)

```go
// Real-time updates:
1. Client connects: WebSocket /ws?token=<jwt>
2. Authenticate JWT token
3. Subscribe client to farm updates
4. When device state changes, broadcast to all connected clients
5. Clients receive live updates (no polling needed)
```

---

## 🆘 Common Issues & Solutions

### Issue: Database connection fails
```
Error: connection refused
```
**Fix**: Check `.env` has correct DB_CONNECTION_STRING, PostgreSQL is running

### Issue: JWT validation fails
```
Error: invalid token
```
**Fix**: Ensure JWT_SECRET_KEY in `.env` matches token generation, check token expiry

### Issue: SQL injection warning
```
ERROR: Possible SQL injection detected
```
**Fix**: Always use prepared statements: `db.Prepare("SELECT * FROM users WHERE id = ?")` with args

### Issue: Rate limiting blocks legitimate traffic
```
Error: 429 Too Many Requests
```
**Fix**: Adjust rate limit in middleware, check if client is making duplicate requests

---

## 📚 Key Documents

- `IG_SPECIFICATIONS_API.md` - All 66 endpoints specifications
- `IG_SPECIFICATIONS_DATABASE.md` - PostgreSQL schema
- `IG_SPECIFICATIONS_SECURITY.md` - Authentication/authorization details
- `01_SPECIFICATIONS_ARCHITECTURE.md` - How middleware fits in overall system

---

## 🧪 Testing Checklist

- ✅ Unit tests for authentication logic
- ✅ Integration tests for database operations
- ✅ API endpoint tests (happy path + error cases)
- ✅ JWT token validation tests
- ✅ Authorization tests (role checks)
- ✅ Database migration tests
- ✅ Load tests (rate limiting)

---

## 📋 Environment Variables

```ini
# .env file (GITIGNORE'D)

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=tokkatot
DB_USER=postgres
DB_PASSWORD=secure_password

# JWT
JWT_SECRET_KEY=very_long_random_secret_key_here
JWT_EXPIRY_HOURS=24
REFRESH_TOKEN_EXPIRY_DAYS=30

# Server
SERVER_PORT=6060
ENVIRONMENT=development

# AI Service
AI_SERVICE_URL=http://localhost:8000
AI_SERVICE_TIMEOUT_SECONDS=30

# MQTT (local hub)
MQTT_BROKER=localhost:1883
MQTT_USERNAME=mqtt_user
MQTT_PASSWORD=mqtt_password
```

---

## 🎯 Your Next Tasks

1. **Implement endpoints** according to `IG_SPECIFICATIONS_API.md`
2. **Set up database** - Create schema, run migrations
3. **Implement authentication** - Login, signup, token refresh
4. **Connect to AI service** - Proxy endpoints to FastAPI
5. **Test thoroughly** - Unit tests, integration tests, security tests
6. **Document integration** - How endpoints connect

---

**Happy coding! 🚀 If unexpected issues arise, check the database schema and API spec first.**
