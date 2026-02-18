# Tokkatot 2.0: Database Design Specification

**Document Version**: 2.0  
**Last Updated**: February 2026  
**Status**: Final Specification

**Database System**: PostgreSQL (primary development)  
**Time-Series DB**: InfluxDB (for sensor data)  
**Cache Layer**: Redis (for sessions and cache)  

---

## Overview

The database is designed using normalized relational model (3NF) for transactional data with separate time-series database for sensor metrics. This specification defines all tables, relationships, and data constraints.

---

## Primary Database Schema (PostgreSQL)

### 1. Users Table

```
users
├── id (UUID, PK)
├── email (TEXT, UNIQUE, NOT NULL)
├── password_hash (TEXT, NOT NULL)
├── name (TEXT, NOT NULL)
├── phone (TEXT, NULL)
├── language (VARCHAR(10), DEFAULT 'km')  -- 'km' or 'en'
├── timezone (VARCHAR(40), DEFAULT 'Asia/Phnom_Penh')
├── role (ENUM: admin|manager|keeper|viewer, DEFAULT 'viewer')
├── avatar_url (TEXT, NULL)
├── is_active (BOOLEAN, DEFAULT TRUE)
├── email_verified (BOOLEAN, DEFAULT FALSE)
├── mfa_enabled (BOOLEAN, DEFAULT FALSE)
├── last_login (TIMESTAMP, NULL)
├── created_at (TIMESTAMP, DEFAULT NOW())
├── updated_at (TIMESTAMP, DEFAULT NOW())
└── deleted_at (TIMESTAMP, NULL)  -- soft delete
```

**Indexes**: email (unique), created_at, role  
**Constraints**: email format validation, password_hash not null  

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
├── role (ENUM: admin|manager|keeper|viewer, NOT NULL)
├── invited_by (UUID, FK → users.id, NOT NULL)
├── is_active (BOOLEAN, DEFAULT TRUE)
├── created_at (TIMESTAMP, DEFAULT NOW())
├── updated_at (TIMESTAMP, DEFAULT NOW())
└── deleted_at (TIMESTAMP, NULL)
```

**Indexes**: farm_id, user_id, farm_id+user_id (composite, unique)  
**Constraints**: farm_id and user_id must exist  
**Purpose**: Many-to-many relationship between users and farms with role assignment  

---

### 4. Devices Table

```
devices
├── id (UUID, PK)
├── farm_id (UUID, FK → farms.id, NOT NULL)
├── device_id (VARCHAR(50), UNIQUE, NOT NULL)  -- hardware ID from device
├── name (TEXT, NOT NULL)
├── type (ENUM: gpio|relay|pwm|adc|servo|sensor, NOT NULL)
├── model (TEXT, NULL)  -- e.g., "ESP32-WROOM-32"
├── firmware_version (VARCHAR(20), NOT NULL)
├── hardware_id (TEXT, UNIQUE, NOT NULL)  -- serial number
├── location (TEXT, NULL)  -- "Room 1", "Water Tank", etc
├── is_active (BOOLEAN, DEFAULT TRUE)
├── is_online (BOOLEAN, DEFAULT FALSE)
├── last_heartbeat (TIMESTAMP, NULL)
├── last_command_status (VARCHAR(50), NULL)  -- success|failed|timeout
├── last_command_at (TIMESTAMP, NULL)
├── created_at (TIMESTAMP, DEFAULT NOW())
├── updated_at (TIMESTAMP, DEFAULT NOW())
└── deleted_at (TIMESTAMP, NULL)
```

**Indexes**: farm_id, device_id (unique), is_active, is_online  
**Relationships**: Many devices per farm  

---

### 5. Device Configs Table

```
device_configs
├── id (UUID, PK)
├── device_id (UUID, FK → devices.id, NOT NULL)
├── parameter_name (TEXT, NOT NULL)
├── parameter_value (JSONB, NOT NULL)
├── unit (VARCHAR(20), NULL)  -- "°C", "%", "ml/s"
├── min_value (DECIMAL(10,2), NULL)
├── max_value (DECIMAL(10,2), NULL)
├── is_calibrated (BOOLEAN, DEFAULT FALSE)
├── created_at (TIMESTAMP, DEFAULT NOW())
├── updated_at (TIMESTAMP, DEFAULT NOW())
└── deleted_at (TIMESTAMP, NULL)
```

**Examples**:
- parameter_name="temp_threshold_high", parameter_value=35, unit="°C"
- parameter_name="humidity_threshold_low", parameter_value=40, unit="%"
- parameter_name="polling_interval", parameter_value=30, unit="seconds"

---

### 6. Schedules Table

```
schedules
├── id (UUID, PK)
├── farm_id (UUID, FK → farms.id, NOT NULL)
├── device_id (UUID, FK → devices.id, NOT NULL)
├── name (TEXT, NOT NULL)
├── schedule_type (ENUM: time_based|duration_based|condition_based, NOT NULL)
├── cron_expression (TEXT, NULL)  -- "0 6 * * *" for 6am daily
├── start_time (TIME, NULL)  -- for simpler time-based schedules
├── end_time (TIME, NULL)
├── on_duration (INTEGER, NULL)  -- seconds
├── off_duration (INTEGER, NULL)  -- seconds
├── repeat_count (INTEGER, NULL)  -- NULL = infinite
├── condition_json (JSONB, NULL)  -- {"type":"temperature","operator":">","value":30}
├── action (ENUM: on|off|toggle, NOT NULL)
├── is_enabled (BOOLEAN, DEFAULT TRUE)
├── priority (INTEGER, DEFAULT 0)  -- for conflict resolution: higher = higher priority
├── created_by (UUID, FK → users.id, NOT NULL)
├── created_at (TIMESTAMP, DEFAULT NOW())
├── updated_at (TIMESTAMP, DEFAULT NOW())
└── deleted_at (TIMESTAMP, NULL)
```

**Indexes**: farm_id, device_id, is_enabled  
**Relationships**: Many schedules per device  

---

### 7. Schedule Executions (Audit Log) Table

```
schedule_executions
├── id (UUID, PK)
├── schedule_id (UUID, FK → schedules.id, NOT NULL)
├── device_id (UUID, FK → devices.id, NOT NULL)
├── scheduled_time (TIMESTAMP, NOT NULL)
├── actual_execution_time (TIMESTAMP, NULL)
├── status (ENUM: pending|executed|failed|timeout|skipped, NOT NULL)
├── command_sent (JSONB, NOT NULL)  -- what was sent to device
├── device_response (JSONB, NULL)  -- device's response
├── error_message (TEXT, NULL)
├── execution_duration_ms (INTEGER, NULL)  -- milliseconds
├── created_at (TIMESTAMP, DEFAULT NOW())
└── deleted_at (TIMESTAMP, NULL)
```

**Indexes**: schedule_id, device_id, scheduled_time  
**Purpose**: Audit trail of all schedule executions  
**Retention**: 5 years  

---

### 8. Device Commands Table

```
device_commands
├── id (UUID, PK)
├── farm_id (UUID, FK → farms.id, NOT NULL)
├── device_id (UUID, FK → devices.id, NOT NULL)
├── command_type (VARCHAR(50), NOT NULL)  -- "on", "off", "set_value"
├── command_params (JSONB, NULL)
├── requested_by (UUID, FK → users.id, NOT NULL)
├── status (ENUM: queued|sent|executing|completed|failed|timeout, NOT NULL)
├── response_data (JSONB, NULL)
├── error_message (TEXT, NULL)
├── requested_at (TIMESTAMP, DEFAULT NOW())
├── sent_at (TIMESTAMP, NULL)
├── completed_at (TIMESTAMP, NULL)
├── created_at (TIMESTAMP, DEFAULT NOW())
└── deleted_at (TIMESTAMP, NULL)
```

**Indexes**: device_id, status, requested_by, created_at  
**Purpose**: Track all user-initiated commands to devices  

---

### 9. Event Logs Table

```
event_logs
├── id (BIGSERIAL, PK)  -- use BIGSERIAL for very high volume
├── farm_id (UUID, FK → farms.id, NOT NULL)
├── device_id (UUID, FK → devices.id, NULL)
├── user_id (UUID, FK → users.id, NULL)
├── event_type (VARCHAR(50), NOT NULL)  -- "device_online", "command_executed", "schedule_triggered"
├── event_category (VARCHAR(20), NOT NULL)  -- "device", "user", "system"
├── severity (ENUM: info|warning|error|critical, NOT NULL)
├── message (TEXT, NOT NULL)
├── event_data (JSONB, NULL)
├── ip_address (INET, NULL)
├── user_agent (TEXT, NULL)
├── created_at (TIMESTAMP, DEFAULT NOW())
└── deleted_at (TIMESTAMP, NULL)
```

**Indexes**: farm_id, device_id, user_id, event_type, created_at  
**Partitioning**: Consider partition by farm_id or date for massive farms  
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

### 13. User Sessions Table (Optional, if using session-based auth)

```
user_sessions
├── id (UUID, PK)
├── user_id (UUID, FK → users.id, NOT NULL)
├── jwt_token_hash (TEXT, NOT NULL)
├── device_name (TEXT, NULL)  -- "Neath's Phone", "Farm Computer"
├── ip_address (INET, NOT NULL)
├── user_agent (TEXT, NOT NULL)
├── is_active (BOOLEAN, DEFAULT TRUE)
├── last_activity (TIMESTAMP, DEFAULT NOW())
├── expires_at (TIMESTAMP, NOT NULL)
├── created_at (TIMESTAMP, DEFAULT NOW())
└── deleted_at (TIMESTAMP, NULL)
```

**Indexes**: user_id, jwt_token_hash (unique), expires_at  
**Purpose**: Track active sessions across devices  

---

## Time-Series Database Schema (InfluxDB)

**Used for**: Sensor data, device metrics, system performance metrics  
**Retention**: 5 years with automatic downsampling

### Measurement: sensor_readings

```
sensor_readings
├── timestamp (required, auto-generated by InfluxDB)
├── farm_id (tag)
├── device_id (tag)
├── sensor_type (tag)  -- "temperature", "humidity", "flow_rate"
├── unit (tag)  -- "°C", "%", "ml/s"
├── value (field, float)
└── quality (field, enum_integer)  -- 0=good, 1=estimated, 2=interpolated
```

**Data Points per Day**: ~40K (1 reading per 2-3 seconds × 100 devices)  
**Query Examples**:
- Last hour temperature readings
- Daily average temperature
- Humidity trends over 30 days
- Peak temperature times

---

### Measurement: device_metrics

```
device_metrics
├── timestamp
├── farm_id (tag)
├── device_id (tag)
├── metric_type (tag)  -- "runtime_hours", "on_cycles", "errors"
├── value (field, float)
└── status (field, string)
```

---

### Measurement: system_metrics

```
system_metrics
├── timestamp
├── metric_type (tag)  -- "api_response_time", "error_rate", "uptime"
├── service_name (tag)  -- "device_service", "schedule_service"
├── value (field, float)
└── percentile (tag, optional)  -- "p50", "p95", "p99"
```

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

### From v1.0 SQLite to v2.0 PostgreSQL

**Migration Path**:
1. Export v1.0 SQLite tables to CSV
2. Transform data to v2.0 schema
3. Import into PostgreSQL staging environment
4. Validate data integrity
5. Run production sync (requires downtime: ~1 hour)
6. Keep SQLite as backup for 30 days

**Transformation Rules**:
- SQLite INTEGER IDs → PostgreSQL UUIDs
- SQLite DATETIME → PostgreSQL TIMESTAMP
- Device configs (flat) → JSONB (hierarchical)
- User roles (text) → ENUM type

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

**Related Documents**
- SPECIFICATIONS_ARCHITECTURE.md
- SPECIFICATIONS_REQUIREMENTS.md
- SPECIFICATIONS_API.md
