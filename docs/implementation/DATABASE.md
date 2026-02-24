# Tokkatot 2.0: Database Design Specification

**Document Version**: 2.2  
**Last Updated**: February 24, 2026  
**Status**: Final Specification

**Database System**: PostgreSQL only (SQLite removed Feb 2026)  
**Time-Series DB**: PostgreSQL `device_readings` table (no InfluxDB)  
**Cache Layer**: Redis (post-MVP; not yet implemented)

**Design Principle**: Simplified for farmers - no complex RBAC, phone number as registration option

> **v2.2 Change**: SQLite has been **permanently removed**. PostgreSQL is now required for all environments (development, staging, production). The `database/sqlite.go` file is deleted. Start the server with a valid `DATABASE_URL` in `middleware/.env`.

## Overview

The database is designed using normalized relational model (3NF) for transactional data with separate time-series database for sensor metrics. This specification defines all tables, relationships, and data constraints. All design decisions prioritize simplicity for elderly farmers with low digital literacy.

---

## Primary Database Schema (PostgreSQL)

### 1. Users Table

```
users
├── id (UUID, PK)
├── email (TEXT, UNIQUE NULLABLE)
├── phone (TEXT, UNIQUE NULLABLE)
├── phone_country_code (VARCHAR(5), NULLABLE)  -- "+855" for Cambodia
├── password_hash (TEXT, NOT NULL)
├── name (TEXT, NOT NULL)
├── language (VARCHAR(10), DEFAULT 'km')  -- 'km' or 'en'
├── timezone (VARCHAR(40), DEFAULT 'Asia/Phnom_Penh')
├── avatar_url (TEXT, NULL)
├── is_active (BOOLEAN, DEFAULT TRUE)
├── contact_verified (BOOLEAN, DEFAULT FALSE)  -- email or phone verified
├── verification_type (VARCHAR(10), NULL)  -- 'email' or 'phone'
├── last_login (TIMESTAMP, NULL)
├── created_at (TIMESTAMP, DEFAULT NOW())
├── updated_at (TIMESTAMP, DEFAULT NOW())
└── deleted_at (TIMESTAMP, NULL)  -- soft delete
```

**Constraints**:
- Either email OR phone required (one-to-one unique constraint)
- NOT both required, NOT both NULL
- password_hash: bcrypt hash (60+ chars)
- contact_verified: Must be true before farm access

**Indexes**: email (unique nullable), phone (unique nullable), created_at, language  
**Notes**: MFA removed (only for Tokkatot admin, not farmers)

---

### 2. Farms Table

```
farms
├── id (UUID, PK)
├── owner_id (UUID, FK → users.id, NOT NULL)
├── name (TEXT, NOT NULL)
├── location (TEXT, NULL)
├── timezone (VARCHAR(40), DEFAULT 'Asia/Phnom_Penh')
├── latitude (DECIMAL(10,8), NULL)
├── longitude (DECIMAL(11,8), NULL)
├── description (TEXT, NULL)
├── image_url (TEXT, NULL)
├── is_active (BOOLEAN, DEFAULT TRUE)
├── created_at (TIMESTAMP, DEFAULT NOW())
├── updated_at (TIMESTAMP, DEFAULT NOW())
└── deleted_at (TIMESTAMP, NULL)
```

**Indexes**: owner_id, created_at  
**Relationships**: One owner to many farms  

---

### 3. Farm Users (Membership Table)

```
farm_users
├── id (UUID, PK)
├── farm_id (UUID, FK → farms.id, NOT NULL)
├── user_id (UUID, FK → users.id, NOT NULL)
├── role (ENUM: farmer|viewer, NOT NULL)
├── invited_by (UUID, FK → users.id, NOT NULL)

├── is_active (BOOLEAN, DEFAULT TRUE)
├── created_at (TIMESTAMP, DEFAULT NOW())
├── updated_at (TIMESTAMP, DEFAULT NOW())
└── deleted_at (TIMESTAMP, NULL)
```

**Indexes**: farm_id, user_id, farm_id+user_id (composite, unique)  
**Constraints**: farm_id and user_id must exist  
**Purpose**: Many-to-many relationship between users and farms with simplified role assignment

**Simplified Role System (Farmer-Centric)**:
- **farmer**: Full farm access — devices, schedules, settings, invite/remove members (farmer or viewer). Multiple farmers can share equal access to the same farm.
- **viewer**: Read-only monitoring + acknowledge alerts. No device control, no settings, no user management.

No complex RBAC — two roles sufficient for all farm operations.

---

### 4. Devices Table

```
devices
├── id (UUID, PK)
├── farm_id (UUID, FK → farms.id, NOT NULL)
├── coop_id (UUID, FK → coops.id, NULL)          -- which coop this device belongs to
├── device_id (VARCHAR(50), UNIQUE, NOT NULL)     -- hardware ID from device
├── name (TEXT, NOT NULL)
├── type (ENUM: gpio|relay|pwm|adc|servo|sensor, NOT NULL)
├── model (TEXT, NULL)  -- e.g., "ESP32-WROOM-32"
├── is_main_controller (BOOLEAN, DEFAULT FALSE)   -- true if this is the Raspberry Pi for the coop
├── firmware_version (VARCHAR(20), NOT NULL)
├── hardware_id (TEXT, UNIQUE, NOT NULL)  -- serial number
├── location (TEXT, NULL)  -- "Room 1", "Water Tank", etc
├── is_active (BOOLEAN, DEFAULT TRUE)
├── is_online (BOOLEAN, DEFAULT FALSE)
├── last_heartbeat (TIMESTAMP, NULL)
├── last_command_status (VARCHAR(50), NULL)  -- success|failed|timeout
├── last_command_at (TIMESTAMP, NULL)
├── created_at (TIMESTAMP, DEFAULT NOW())
└── updated_at (TIMESTAMP, DEFAULT NOW())
```

**Indexes**: farm_id, coop_id, device_id (unique), is_active, is_online, (coop_id, is_main_controller)  
**Relationships**: Many devices per farm; each device optionally scoped to a coop  
**Key field**: `is_main_controller = true` identifies the Raspberry Pi controller for each coop  

---

### 5. Device Configurations Table

> **Implemented**: Table name is `device_configurations` in PostgreSQL schema (see `database/postgres.go`).

```
device_configurations
├── id (UUID, PK)
├── device_id (UUID, FK → devices.id, NOT NULL)
├── parameter_name (VARCHAR(100), NOT NULL)
├── parameter_value (TEXT, NOT NULL)
├── unit (VARCHAR(20), NULL)  -- "°C", "%", "ml/s"
├── min_value (DECIMAL(10,4), NULL)
├── max_value (DECIMAL(10,4), NULL)
├── is_calibrated (BOOLEAN, DEFAULT FALSE)
├── calibrated_at (TIMESTAMP, NULL)
├── created_at (TIMESTAMP, DEFAULT NOW())
├── updated_at (TIMESTAMP, DEFAULT NOW())
UNIQUE(device_id, parameter_name)
```

**Examples**:
- parameter_name="temp_threshold_high", parameter_value="35", unit="°C"
- parameter_name="humidity_threshold_low", parameter_value="40", unit="%"
- parameter_name="polling_interval", parameter_value="30", unit="seconds"

**Upsert pattern**: `ON CONFLICT (device_id, parameter_name) DO UPDATE SET ...` — safe to call repeatedly.

---

### 6. Schedules Table

```
schedules
├── id (UUID, PK)
├── farm_id (UUID, FK → farms.id, NOT NULL)
├── coop_id (UUID, FK → coops.id, NULL)
├── device_id (UUID, FK → devices.id, NOT NULL)
├── name (TEXT, NOT NULL)
├── schedule_type (ENUM: time_based|duration_based|condition_based, NOT NULL)
├── cron_expression (TEXT, NULL)  -- "0 6,12,18 * * *" for 6am, 12pm, 6pm daily
├── on_duration (INTEGER, NULL)  -- For duration_based: seconds to stay ON
├── off_duration (INTEGER, NULL)  -- For duration_based: seconds to stay OFF
├── condition_json (JSONB, NULL)  -- For condition_based: {"sensor":"temperature","operator":">","threshold":30}
├── action (VARCHAR(20), NOT NULL)  -- "on", "off", "set_value"
├── action_value (TEXT, NULL)  -- Optional value (e.g., PWM duty cycle)
├── action_duration (INTEGER, NULL)  -- For time_based: auto-turn-off after X seconds
├── action_sequence (JSONB, NULL)  -- For time_based: multi-step pattern [{"action":"ON","duration":30},{"action":"OFF","duration":10}]
├── priority (INTEGER, DEFAULT 0)  -- 0-10, higher = more important (conflict resolution)
├── is_active (BOOLEAN, DEFAULT TRUE)
├── next_execution (TIMESTAMP, NULL)  -- Calculated next run time for time_based
├── last_execution (TIMESTAMP, NULL)  -- Last successful execution
├── execution_count (INTEGER, DEFAULT 0)  -- Total times executed
├── created_by (UUID, FK → users.id, NOT NULL)
├── created_at (TIMESTAMP, DEFAULT NOW())
├── updated_at (TIMESTAMP, DEFAULT NOW())
└── deleted_at (TIMESTAMP, NULL)
```

**Field Purpose**:
- **time_based**: Uses `cron_expression` + optionally `action_duration` (simple auto-off) OR `action_sequence` (multi-step pattern)
- **duration_based**: Uses `on_duration` and `off_duration` for continuous cycling (e.g., conveyor ON 10min, OFF 15min, repeat)
- **condition_based**: Uses `condition_json` for sensor-driven automation (e.g., pump ON when water < 20%)

**New Fields (Feb 2026)**:
- `action_duration`: Time-based schedules can auto-turn-off after X seconds (e.g., feeder ON at 6AM, auto-off after 15 minutes)
- `action_sequence`: Multi-step patterns for pulse operations (e.g., feeder: ON 30sec, pause 10sec, ON 30sec, pause 10sec)
  - Format: `[{"action":"ON","duration":30},{"action":"OFF","duration":10},{"action":"ON","duration":30},{"action":"OFF","duration":10}]`
  - Maximum 20 steps per sequence
  - After sequence completes, device returns to OFF state until next cron trigger

**Indexes**: farm_id, device_id, is_active, next_execution, schedule_type  
**Relationships**: Many schedules per device  
**Use Cases**: See `docs/AUTOMATION_USE_CASES.md` for detailed farmer scenarios  

---

### 7. Schedule Executions (Audit Log) Table

```
schedule_executions
├── id (UUID, PK)
├── schedule_id (UUID, FK → schedules.id, NOT NULL)
├── device_id (UUID, FK → devices.id, NOT NULL)
├── scheduled_time (TIMESTAMP, NOT NULL)
├── actual_execution_time (TIMESTAMP, NULL)
├── status (ENUM: executed|failed|skipped, NOT NULL)
├── execution_duration_ms (INTEGER, NULL)  -- milliseconds
├── device_response (JSONB, NULL)  -- raw response from device
├── error_message (TEXT, NULL)
└── created_at (TIMESTAMP, DEFAULT NOW())
```

**Indexes**: schedule_id, (scheduled_time DESC), status  
**Purpose**: Audit trail of all schedule executions  
**Retention**: 5 years  

---

### 8. Device Commands Table

```
device_commands
├── id (UUID, PK)
├── device_id (UUID, FK → devices.id, NOT NULL)
├── farm_id (UUID, FK → farms.id, NOT NULL)
├── coop_id (UUID, FK → coops.id, NULL)            -- coop the command targets
├── issued_by (UUID, FK → users.id, NOT NULL)       -- who sent the command
├── command_type (VARCHAR(50), NOT NULL)             -- "on", "off", "set_value", "status", "reboot"
├── command_value (TEXT, NULL)                       -- optional value (e.g., PWM duty cycle)
├── status (ENUM: pending|success|failed|timeout, NOT NULL, DEFAULT 'pending')
├── response (TEXT, NULL)                            -- raw device response
├── issued_at (TIMESTAMP, DEFAULT NOW())
├── executed_at (TIMESTAMP, NULL)
└── created_at (TIMESTAMP, DEFAULT NOW())
```

**Indexes**: device_id, farm_id, coop_id, status, (created_at DESC)  
**Purpose**: Track all user-initiated commands to IoT devices  
**Valid command_types**: `on`, `off`, `set_value`, `status`, `reboot`  

---

### 9. Event Logs Table

```
event_logs
├── id (UUID, PK)
├── farm_id (UUID, FK → farms.id, NOT NULL)
├── user_id (UUID, FK → users.id, NOT NULL)
├── event_type (VARCHAR(50), NOT NULL)  -- "login", "device_control", "schedule_update"
├── resource_id (UUID, NULL)             -- the affected farm/coop/device/schedule UUID
├── old_value (JSONB, NULL)              -- state before the change
├── new_value (JSONB, NULL)              -- state after the change
├── ip_address (VARCHAR(45), NULL)
└── created_at (TIMESTAMP, DEFAULT NOW())
```

**Indexes**: (farm_id, user_id), (created_at DESC)  
**Purpose**: Lightweight audit trail of user actions (device control, schedule changes, logins)  
**Pattern**: `old_value` / `new_value` JSONB diff pattern — store the before/after JSON for any resource change  
**Purpose**: Immutable audit trail  
**Retention**: 5 years, immutable, searchable  

---

### 10. Alerts Table

```
alerts
├── id (UUID, PK)
├── farm_id (UUID, FK → farms.id, NOT NULL)
├── device_id (UUID, FK → devices.id, NULL)
├── alert_type (VARCHAR(50), NOT NULL)  -- "temp_high", "device_offline", "low_battery"
├── severity (ENUM: info|warning|critical, NOT NULL)
├── message (TEXT, NOT NULL)
├── threshold_value (DECIMAL(10,2), NULL)
├── actual_value (DECIMAL(10,2), NULL)
├── is_active (BOOLEAN, DEFAULT TRUE)
├── acknowledged_by (UUID, FK → users.id, NULL)
├── acknowledged_at (TIMESTAMP, NULL)
├── triggered_at (TIMESTAMP, DEFAULT NOW())
├── resolved_at (TIMESTAMP, NULL)
└── deleted_at (TIMESTAMP, NULL)
```

**Indexes**: farm_id, device_id, is_active, triggered_at  
**Purpose**: Alert history and acknowledgement tracking  

---

### 11. Alert Subscriptions Table

```
alert_subscriptions
├── id (UUID, PK)
├── user_id (UUID, FK → users.id, NOT NULL)
├── alert_type (VARCHAR(50), NOT NULL)  -- which alerts to receive
├── channel (ENUM: email|sms|push|telegram, NOT NULL)  -- how to receive
├── is_enabled (BOOLEAN, DEFAULT TRUE)
├── quiet_hours_start (TIME, NULL)  -- e.g., 20:00 (8pm, no alerts)
├── quiet_hours_end (TIME, NULL)  -- e.g., 06:00 (6am, resume alerts)
├── created_at (TIMESTAMP, DEFAULT NOW())
├── updated_at (TIMESTAMP, DEFAULT NOW())
└── deleted_at (TIMESTAMP, NULL)
```

**Indexes**: user_id, alert_type  
**Purpose**: User notification preferences  

---

### 12. Notification Log Table (In-App Dashboard Only)

```
notification_log
├── id (UUID, PK)
├── user_id (UUID, FK → users.id, NOT NULL)
├── farm_id (UUID, FK → farms.id, NOT NULL)
├── alert_id (UUID, FK → alerts.id, NULL)
├── notification_type (VARCHAR(50), NOT NULL)  -- device_offline|temperature_alert|etc
├── severity (ENUM: urgent|important|info, NOT NULL)
├── title (VARCHAR(100), NOT NULL)  -- short display title
├── message (TEXT, NOT NULL)  -- detailed message text
├── is_read (BOOLEAN, DEFAULT FALSE)  -- farmer marks as read
├── metadata (JSONB, NULL)  -- device_id, temperature value, etc
├── created_at (TIMESTAMP, DEFAULT NOW())
├── read_at (TIMESTAMP, NULL)
└── deleted_at (TIMESTAMP, NULL)
```

**Note**: In-app only. No SMS, email, or push notifications.

**Indexes**: user_id, status, created_at  
**Purpose**: Notification delivery tracking  

---

### 13. User Sessions Table

> **Implemented**: Table exists in PostgreSQL schema (see `database/postgres.go`).

```
user_sessions
├── id (UUID, PK)
├── user_id (UUID, FK → users.id ON DELETE CASCADE, NOT NULL)
├── device_name (VARCHAR(255), NULL)  -- "Neath's Phone", "Farm Computer"
├── ip_address (VARCHAR(45), NULL)
├── user_agent (TEXT, NULL)
├── refresh_token (TEXT, NOT NULL, UNIQUE)
├── last_activity (TIMESTAMP, DEFAULT NOW())
├── expires_at (TIMESTAMP, NOT NULL)
└── created_at (TIMESTAMP, DEFAULT NOW())
```

**Indexes**: user_id, expires_at, refresh_token (unique)  
**Purpose**: Track active sessions across devices; used by session management endpoints  
**Endpoints**: `GET /users/sessions`, `DELETE /users/sessions/:session_id`

---

### 14. Device Readings Table

```
device_readings
├── id (UUID, PK)
├── device_id (UUID, FK → devices.id, NOT NULL)
├── sensor_type (VARCHAR(50), NOT NULL)  -- "temperature", "humidity", "water_level"
├── value (DECIMAL(10,4), NOT NULL)
├── unit (VARCHAR(20), NOT NULL, DEFAULT '')  -- "°C", "%", "cm"
├── quality (VARCHAR(20), NOT NULL, DEFAULT 'good')  -- "good", "degraded", "bad"
└── timestamp (TIMESTAMP, DEFAULT NOW())
```

**Indexes**: `(device_id)`, `(timestamp DESC)`, `(device_id, sensor_type)`  
**Purpose**: PostgreSQL time-series sensor data storage  
**Endpoints**: `GET /farms/:farm_id/devices/:device_id/history`, `GET /farms/:farm_id/reports/device-metrics`

---

## Cache Layer (Redis)

**Used for**: Sessions, rate limiting counters, real-time cache  
**TTL Strategy**: Automatic expiration per value type

### Data Structures

```
sessions:
  Key: session:{sessionId}
  Value: {userId, farmId, role, expiresAt}
  TTL: 24 hours

rate_limit:
  Key: rate_limit:{userId}:{endpoint}
  Value: request count
  TTL: 1 minute

device_state:
  Key: device:{deviceId}:state
  Value: {status, last_update, online}
  TTL: 5 minutes (refresh on each heartbeat)

schedule_queue:
  Key: schedule:queue
  Value: Sorted set of next 100 schedules to execute
  TTL: none (persistent)

device_command_queue:
  Key: device:{deviceId}:commands
  Value: Queue of pending commands
  TTL: none (persistent until executed)

pub_sub channels:
  farm:{farmId}:updates  # broadcasts of device updates
  farm:{farmId}:alerts   # broadcasts of alerts
  device:{deviceId}:commands  # device command delivery
```

---

## Data Relationships (ERD)

```
users (1) ──→ (many) farm_users (many) ──→ (1) farms
users (1) ──→ (many) event_logs
users (1) ──→ (many) device_commands
users (1) ──→ (many) alert_subscriptions
users (1) ──→ (many) user_sessions

farms (1) ──→ (many) devices
farms (1) ──→ (many) schedules
farms (1) ──→ (many) device_commands
farms (1) ──→ (many) event_logs
farms (1) ──→ (many) alerts

devices (1) ──→ (many) device_configs
devices (1) ──→ (many) schedules
devices (1) ──→ (many) schedule_executions
devices (1) ──→ (many) device_commands
devices (1) ──→ (many) event_logs
devices (1) ──→ (many) alerts

schedules (1) ──→ (many) schedule_executions

alerts (1) ──→ (many) notifications
```

---

## Constraints & Validations

### Column Constraints

| Table | Column | Type | Constraint |
|-------|--------|------|-----------|
| users | email | TEXT | Format: RFC 5322, UNIQUE |
| users | password_hash | TEXT | 60+ chars (bcrypt hash) |
| users | language | VARCHAR(10) | IN ('km', 'en') |
| user_sessions | expires_at | TIMESTAMP | > NOW() on insert |
| farms | timezone | VARCHAR(40) | Valid IANA timezone |
| devices | device_id | VARCHAR(50) | Unique within farm_id |
| device_configs | parameter_value | JSONB | Valid JSON |
| schedules | cron_expression | TEXT | Valid cron syntax |
| schedules | on_duration | INTEGER | > 0 if duration_based |
| event_logs | severity | ENUM | IN (info, warning, error, critical) |

---

## Indexing Strategy

### High-Priority Indexes (Essential)
- users(email) - UNIQUE for login
- devices(farm_id) - JOIN with farms
- schedules(device_id) - JOIN with devices
- event_logs(created_at) - Time-range queries
- device_commands(status) - Command queuing

### Medium-Priority Indexes (Important)
- farms(owner_id) - User's farms lookup
- farm_users(user_id, farm_id) - Composite for multi-farm users
- event_logs(farm_id, created_at) - Range + farm filtering
- alerts(is_active, farm_id) - Active alerts per farm
- device_configs(device_id) - Configuration lookup

### Low-Priority Indexes (Optimization)
- event_logs(event_type) - Filter by event type
- notifications(status) - Retry pending notifications
- device_commands(completed_at) - Archive old commands

---

## Query Patterns

### Query Patterns | Execution Plan
|---|---|
| Get farm's devices | SELECT * FROM devices WHERE farm_id = ? AND deleted_at IS NULL |
| Get user's farms | SELECT f.* FROM farms f JOIN farm_users fu ON f.id = fu.farm_id WHERE fu.user_id = ? |
| Get device history (24h) | SELECT * FROM sensor_readings WHERE device_id = ? AND timestamp > NOW() - INTERVAL '24 hours' |
| Get scheduled jobs (next 1h) | SELECT * FROM schedules WHERE device_id = ? AND is_enabled = TRUE ORDER BY cron_expression |
| Audit trail (by user) | SELECT * FROM event_logs WHERE user_id = ? AND created_at > DATE_TRUNC('month', NOW()) |

---

## Data Migration Strategy

> ⚠️ **v2.0 Note (Feb 2026)**: SQLite was permanently removed. Tokkatot is **PostgreSQL-only**. All installations run directly on PostgreSQL with UUID primary keys, JSONB config columns, and ENUM role types. There is no migration path from a local SQLite database — fresh installs use `database/postgres.go` exclusively.

---

## Backup & Recovery

- **Backup Frequency**: Daily at 2 AM UTC
- **Retention Window**: 30 days rolling
- **Backup Location**: S3-compatible storage
- **Recovery RTO**: < 4 hours
- **Recovery RPO**: < 1 hour (max data loss acceptable)
- **Test Recovery**: Monthly restore drill

---

## Performance Tuning

- **Connection Pooling**: 100 connections per service instance
- **Query Timeout**: 30 seconds (slow queries logged)
- **Batch Inserts**: Use PostgreSQL COPY for bulk operations
- **Partitioning**: Consider partition by farm_id for event_logs if > 10TB
- **Archival**: Move events > 5 years old to cold storage monthly

---

**Version History**
| Version | Date | Changes |
|---------|------|---------|
| 2.0 | Feb 2026 | Initial production schema |
| 2.1 | Feb 23, 2026 | Added action_duration and action_sequence to schedules; added alert_subscriptions |
| 2.2 | Feb 24, 2026 | **SQLite removed** — PostgreSQL only; added device_readings (MVP time-series); renamed device_configs → device_configurations; added refresh_token to user_sessions; all MVP non-AI endpoints implemented |

**Related Documents**
- SPECIFICATIONS_ARCHITECTURE.md
- SPECIFICATIONS_REQUIREMENTS.md
- SPECIFICATIONS_API.md
