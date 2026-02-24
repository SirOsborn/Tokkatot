# Tokkatot 2.0: API Specification

**Document Version**: 2.3  
**Last Updated**: February 24, 2026  
**Status**: MVP Complete (all non-AI endpoints implemented) + Temperature Timeline + Full Frontend Rebuild  
**Base URL**: `https://api.tokkatot.local/v1` (production), `http://localhost:3000/v1` (development)

> **Implementation Status (v2.3)**: All 67 non-AI API endpoints are implemented and build-verified. SQLite removed — PostgreSQL required. AI disease detection endpoint is stubbed; disease-detection UI shows "Coming Soon" overlay. Temperature timeline endpoint live (`GET /farms/:farm_id/coops/:coop_id/temperature-timeline`). Frontend fully rebuilt (Vue.js 3 CDN, Mi Sans, Material Symbols, new design system). New static page route added: `GET /schedules` → `pages/schedules.html` (schedule CRUD + multi-step sequence builder).

### Frontend Static Page Routes (`middleware/main.go`)

The Go middleware serves all frontend pages directly. These are **not** API endpoints — they return HTML files.

| Browser Route | File Served | Notes |
|---|---|---|
| `GET /` | `pages/index.html` | Dashboard |
| `GET /login` | `pages/login.html` | — |
| `GET /register` | `pages/signup.html` | Calls `POST /v1/auth/signup` |
| `GET /profile` | `pages/profile.html` | — |
| `GET /settings` | `pages/settings.html` | — |
| `GET /disease-detection` | `pages/disease-detection.html` | Coming Soon overlay |
| `GET /monitoring` | `pages/monitoring.html` | Temperature timeline |
| `GET /schedules` | `pages/schedules.html` | ✅ Live — schedule CRUD + sequence builder |

Static asset directories also served: `/assets`, `/components`, `/css`, `/js`.

---

## Overview

The Tokkatot API is a RESTful service with real-time capabilities for managing smart chicken farming operations. Designed for elderly farmers with low digital literacy in Cambodia.

**Key Characteristics**:
- RESTful design with JSON payloads
- JWT token-based authentication (simple, no MFA required)
- **Role System**: Farmer (full control), Viewer (read-only + ack alerts)
- Pagination support for list endpoints
- Real-time updates via WebSocket (Socket.io)
- Device communication via MQTT
- **Device Setup by Tokkatot Team Only** (farmers cannot add devices)

**Farmer-Centric Design**:
- Registration: Email OR phone number (farmers choose one)
- Simple language: Khmer/English toggle
- No complex permission system
- Tokkatot team manages all device setup and configuration

---

## Authentication & Authorization

### JWT Token Flow

**Token Issuance**:
```
POST /auth/login
Body: { email, password }
Response: { access_token, refresh_token, expires_in }
```

**Token Structure**:
```json
{
  "sub": "user_id_uuid",
  "email": "farmer@example.com",
  "phone": "+855012345678",
  "farm_id": "farm_id_uuid",
  "role": "farmer",
  "iat": 1708324800,
  "exp": 1708411200
}
```

> `email` and `phone` are nullable — only the field used at signup/login is populated.

**Token Expiry & Refresh**:
- Access token: 24 hours
- Refresh token: 30 days
- New tokens issued simultaneously on successful refresh
- MFA: **Not required** for farmers (optional for admin users only)

### Role System

| Role | Permissions | Notes |
|------|-------------|-------|
| **Farmer** | Full farm control: devices, schedules, coops, farm settings, invite/remove members | Created at farm registration; multiple farmers can share a farm |
| **Viewer** | Read-only monitoring + acknowledge alerts (maintenance workers for large farms) | No device control, no settings |

**No complex RBAC** — two farm roles sufficient for operations.  
**Tokkatot system staff** are not a farm role — they manage registration keys, JWT secrets, and system-level access outside of `farm_users` entirely.  
**Device Management**: Only Tokkatot team can add/remove devices (not farmers).

### Permission Checks

All endpoints check:
1. Valid JWT token present
2. User role matches minimum required role
3. User has access to requested farm (via farm_users table)

---

## API Endpoints

### Authentication Endpoints (8)

#### 1. User Login
```
POST /auth/login
Content-Type: application/json

Request:
{
  "email": "farmer@example.com",
  "password": "secure_password",
  "device_name": "Samsung A12"  // optional
}

Response (200):
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "expires_in": 86400,
  "user": {
    "id": "uuid",
    "email": "farmer@example.com",
    "name": "Neath",
    "role": "farmer",
    "language": "km"
  }
}

Errors:
- 401: Invalid credentials
- 429: Too many login attempts (rate limited: 10/minute)
```

#### 2. User Registration (Email OR Phone)
```
POST /auth/signup
Content-Type: application/json

Request (Option A - Email):
{
  "email": "farmer@example.com",
  "password": "Password123",
  "name": "Neath",
  "language": "km"
}

Request (Option B - Phone):
{
  "phone": "+855987654321",
  "password": "Password123",
  "name": "Neath",
  "language": "km"
}

Response (201):
{
  "user_id": "uuid",
  "contact": "farmer@example.com",  // or phone number
  "verification_required": true,
  "message": "Verification code sent"
}

Validation:
- Email format OR phone format (E.164) - one required, not both
- Password: min 8 chars (for farmer simplicity, not 10)
- Language: 'km' or 'en'
- Name: required

Verification:
- Email: Link in email
- Phone: SMS code (6 digits, valid 24 hours)

Notes:
- Farmers choose registration method (email or phone)
- No complex password requirements
- Simple verification process
```

#### 3. Refresh Token
```
POST /auth/refresh
Authorization: Bearer {refresh_token}

Response (200):
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "expires_in": 86400
}

Errors:
- 401: Refresh token expired or invalid
```

#### 4. Logout
```
POST /auth/logout
Authorization: Bearer {access_token}

Response (200):
{
  "message": "Successfully logged out"
}
```

#### 5. Verify Email
```
POST /auth/verify-email
Content-Type: application/json

Request:
{
  "token": "verification_token_from_email"
}

Response (200):
{
  "message": "Email verified successfully"
}

Errors:
- 400: Invalid or expired token
```

#### 6. Request Password Reset
```
POST /auth/forgot-password
Content-Type: application/json

Request:
{
  "email": "farmer@example.com"
}

Response (200):
{
  "message": "Password reset email sent"
}

Rate Limit: 5 per hour per email
```

#### 7. Reset Password
```
POST /auth/reset-password
Content-Type: application/json

Request:
{
  "token": "reset_token_from_email",
  "password": "new_secure_password"
}

Response (200):
{
  "message": "Password reset successfully"
}
```

#### 8. Change Password
```
POST /auth/change-password
Authorization: Bearer {access_token}
Content-Type: application/json

Request:
{
  "current_password": "old_password",
  "new_password": "new_password"
}

Response (200):
{
  "message": "Password changed successfully"
}

Errors:
- 401: Current password incorrect
```

---

**Note**: MFA is optional for Tokkatot admin users only. Farmers do not use MFA (simplified authentication for low-literacy users).

---

### User Management Endpoints (5)

#### 11. Get Current User Profile
```
GET /users/me
Authorization: Bearer {access_token}

Response (200):
{
  "id": "uuid",
  "email": "farmer@example.com",
  "name": "Neath",
  "phone": "+855987654321",
  "language": "km",
  "timezone": "Asia/Phnom_Penh",
  "role": "farmer",
  "avatar_url": "https://...",
  "mfa_enabled": false,
  "last_login": "2026-02-18T15:30:00Z",
  "created_at": "2026-01-15T10:00:00Z"
}
```

#### 12. Update User Profile
```
PUT /users/me
Authorization: Bearer {access_token}
Content-Type: application/json

Request:
{
  "name": "Neath Updated",
  "phone": "+855987654321",
  "language": "en",
  "timezone": "Asia/Phnom_Penh",
  "avatar_url": "https://..."
}

Response (200):
{
  "id": "uuid",
  "name": "Neath Updated",
  ...fields updated...
}
```

#### 13. Get User Sessions
```
GET /users/sessions
Authorization: Bearer {access_token}

Response (200):
{
  "sessions": [
    {
      "id": "uuid",
      "device_name": "Samsung A12",
      "ip_address": "192.168.1.1",
      "user_agent": "Mozilla/5.0...",
      "last_activity": "2026-02-19T10:30:00Z",
      "expires_at": "2026-03-21T10:30:00Z",
      "is_current": true
    }
  ]
}
```

#### 14. Revoke Session
```
DELETE /users/sessions/{session_id}
Authorization: Bearer {access_token}

Response (200):
{
  "message": "Session revoked"
}
```

#### 15. Get User Activity Log
```
GET /users/activity-log?limit=50&offset=0&days=30
Authorization: Bearer {access_token}

Response (200):
{
  "activities": [
    {
      "id": "uuid",
      "event_type": "login",
      "ip_address": "192.168.1.1",
      "device_name": "Samsung A12",
      "timestamp": "2026-02-19T10:30:00Z"
    }
  ],
  "total": 150,
  "limit": 50,
  "offset": 0
}
```

---

### Farm Management Endpoints (8)

#### 16. Create Farm
```
POST /farms
Authorization: Bearer {access_token}
Content-Type: application/json

Request:
{
  "name": "Main Chicken Farm",
  "location": "Siem Reap Province",
  "latitude": 13.3671,
  "longitude": 103.8448,
  "timezone": "Asia/Phnom_Penh",
  "description": "Growing 500 chickens",
  "image_url": "https://..."
}

Response (201):
{
  "id": "farm_uuid",
  "owner_id": "user_uuid",
  "name": "Main Chicken Farm",
  "location": "Siem Reap Province",
  "created_at": "2026-02-19T10:30:00Z"
}

Permission: farmer
```
```
GET /farms?limit=20&offset=0
Authorization: Bearer {access_token}

Response (200):
{
  "farms": [
    {
      "id": "farm_uuid",
      "name": "Main Chicken Farm",
      "owner_id": "user_uuid",
      "location": "Siem Reap Province",
      "is_active": true,
      "device_count": 6,
      "online_devices": 5,
      "created_at": "2026-02-19T10:30:00Z"
    }
  ],
  "total": 2,
  "limit": 20,
  "offset": 0
}
```

#### 18. Get Farm Details
```
GET /farms/{farm_id}
Authorization: Bearer {access_token}

Response (200):
{
  "id": "farm_uuid",
  "name": "Main Chicken Farm",
  "owner_id": "user_uuid",
  "location": "Siem Reap Province",
  "latitude": 13.3671,
  "longitude": 103.8448,
  "timezone": "Asia/Phnom_Penh",
  "description": "Growing 500 chickens",
  "image_url": "https://...",
  "is_active": true,
  "statistics": {
    "total_devices": 6,
    "online_devices": 5,
    "total_schedules": 12,
    "active_alerts": 2
  },
  "created_at": "2026-02-19T10:30:00Z",
  "updated_at": "2026-02-19T10:30:00Z"
}

Permission: User must have access to farm via farm_users
```

#### 19. Update Farm
```
PUT /farms/{farm_id}
Authorization: Bearer {access_token}
Content-Type: application/json

Request:
{
  "name": "Updated Farm Name",
  "location": "New Location",
  "timezone": "Asia/Phnom_Penh",
  "description": "Updated description"
}

Response (200):
{
  "id": "farm_uuid",
  "name": "Updated Farm Name",
  ...updated fields...
}

Permission: farmer
```

#### 20. Delete Farm
```
DELETE /farms/{farm_id}
Authorization: Bearer {access_token}

Response (200):
{
  "message": "Farm deleted successfully"
}

Permission: farmer
Effect: Soft delete (deleted_at timestamp set)
```

#### 21. Get Farm Members
```
GET /farms/{farm_id}/members?limit=20&offset=0
Authorization: Bearer {access_token}

Response (200):
{
  "members": [
    {
      "user_id": "uuid",
      "name": "Neath",
      "email": "farmer@example.com",
      "role": "farmer",
      "invited_by": "inviter_uuid",
      "joined_at": "2026-02-19T10:30:00Z"
    }
  ],
  "total": 3
}

Permission: farmer
```

#### 22. Invite Member to Farm
```
POST /farms/{farm_id}/members
Authorization: Bearer {access_token}
Content-Type: application/json

Request:
{
  "email": "newfarmer@example.com",
  "role": "viewer"
}

Response (201):
{
  "invitation_id": "uuid",
  "email": "newfarmer@example.com",
  "role": "viewer",
  "status": "pending",
  "expires_at": "2026-02-26T10:30:00Z"
}

Permission: farmer
Actions: Send invitation email with acceptance link
```

#### 23. Update Member Role
```
PUT /farms/{farm_id}/members/{user_id}
Authorization: Bearer {access_token}
Content-Type: application/json

Request:
{
  "role": "viewer"
}

Response (200):
{
  "user_id": "uuid",
  "name": "Farmer Name",
  "email": "farmer@example.com",
  "role": "viewer",
  "updated_at": "2026-02-19T10:30:00Z"
}

Permission: User must be farm admin
```

---

### Coop Management Endpoints (5 + 1 temperature timeline)

#### 24. List Coops
```
GET /farms/{farm_id}/coops
Authorization: Bearer {access_token}

Response (200): Array of coop objects (id, name, capacity, created_at)
Permission: viewer+
```

#### 25. Create Coop
```
POST /farms/{farm_id}/coops
Authorization: Bearer {access_token}
Body: { "name": "Coop A", "capacity": 200 }

Response (201): Coop object
Permission: farmer
```

#### 26. Get Coop
```
GET /farms/{farm_id}/coops/{coop_id}
Authorization: Bearer {access_token}

Response (200): Coop object
Permission: viewer+
```

#### 27. Update Coop
```
PUT /farms/{farm_id}/coops/{coop_id}
Authorization: Bearer {access_token}
Body: { "name": "Coop B", "capacity": 250 }

Response (200): Updated coop object
Permission: farmer
```

#### 28. Delete Coop
```
DELETE /farms/{farm_id}/coops/{coop_id}
Authorization: Bearer {access_token}

Response (200): { "message": "Coop deleted" }
Permission: farmer
Effect: Soft delete
```

#### 29. Temperature Timeline (Apple Weather-style)
```
GET /farms/{farm_id}/coops/{coop_id}/temperature-timeline?days=7
Authorization: Bearer {access_token}

Response (200):
{
  "success": true,
  "data": {
    "coop_id": "uuid",
    "coop_name": "Coop A",
    "farm_id": "uuid",
    "device_id": "uuid",
    "sensor_found": true,
    "current_temp": 34.2,
    "bg_hint": "hot",
    "today": {
      "date": "2026-02-24",
      "hourly": [
        { "hour": "06:00", "temp": 27.4 },
        { "hour": "07:00", "temp": 28.1 },
        ...
      ],
      "high": { "temp": 34.5, "time": "14:00" },
      "low":  { "temp": 24.1, "time": "05:00" }
    },
    "history": [
      {
        "date": "2026-02-24",
        "label": "Today",
        "high": { "temp": 34.5, "time": "14:00" },
        "low":  { "temp": 24.1, "time": "05:00" }
      },
      {
        "date": "2026-02-23",
        "label": "Yesterday",
        "high": { "temp": 35.1, "time": "14:00" },
        "low":  { "temp": 24.0, "time": "04:00" }
      },
      { "date": "2026-02-22", "label": "Fri", "high": {...}, "low": {...} }
    ]
  },
  "message": "Temperature timeline fetched"
}

bg_hint values:
  scorching  >= 35°C  (deep red gradient)
  hot        >= 32°C  (red-orange gradient)
  warm       >= 28°C  (orange gradient)
  neutral    >= 24°C  (green gradient)
  cool       >= 20°C  (blue gradient)
  cold        < 20°C  (dark blue gradient)

If sensor_found = false: returns 200 with empty today/history and no device_id.
Only temperature sensor_type readings are queried — humidity is excluded.

Permission: viewer+
Frontend page: /monitoring
```

---

### Device Management Endpoints (10)

#### 24. List Farm Devices
```
GET /farms/{farm_id}/devices?limit=20&offset=0&filter=online
Authorization: Bearer {access_token}

Query Parameters:
- limit: 1-100 (default 20)
- offset: pagination offset
- filter: 'all'|'online'|'offline' (default 'all')
- type: filter by device type

Response (200):
{
  "devices": [
    {
      "id": "device_uuid",
      "device_id": "ESP32_001",
      "name": "Water Pump",
      "type": "relay",
      "model": "ESP32-WROOM-32",
      "firmware_version": "1.2.3",
      "location": "Water Tank",
      "is_online": true,
      "last_heartbeat": "2026-02-19T10:28:00Z",
      "battery_level": null,
      "last_command_status": "success",
      "created_at": "2026-01-15T10:00:00Z"
    }
  ],
  "total": 6,
  "online_count": 5
}
```

#### 25. Get Device Details
```
GET /farms/{farm_id}/devices/{device_id}
Authorization: Bearer {access_token}

Response (200):
{
  "id": "device_uuid",
  "device_id": "ESP32_001",
  "name": "Water Pump",
  "type": "relay",
  "model": "ESP32-WROOM-32",
  "firmware_version": "1.2.3",
  "hardware_id": "SN123456789",
  "location": "Water Tank",
  "is_active": true,
  "is_online": true,
  "last_heartbeat": "2026-02-19T10:28:00Z",
  "last_command_status": "success",
  "last_command_at": "2026-02-19T10:15:00Z",
  "configuration": [
    {
      "parameter_name": "polling_interval",
      "parameter_value": 30,
      "unit": "seconds"
    }
  ],
  "created_at": "2026-01-15T10:00:00Z",
  "updated_at": "2026-02-19T10:30:00Z"
}

Permission: User must have access to farm
```

#### 26. Add Device (Tokkatot Team Only)
```
POST /farms/{farm_id}/devices
Authorization: Bearer {access_token}
Content-Type: application/json

IMPORTANT: This endpoint is for Tokkatot team administrators only.
Farmers CANNOT add devices themselves. Device setup is done by Tokkatot team.

Request:
{
  "device_id": "ESP32_001",
  "name": "Water Pump",
  "type": "relay",
  "model": "ESP32-WROOM-32",
  "firmware_version": "1.2.3",
  "hardware_id": "SN123456789",
  "location": "Water Tank"
}

Response (201):
{
  "id": "device_uuid",
  "device_id": "ESP32_001",
  ...device data...
}

Validation:
- device_id: unique within farm
- type: in ('gpio', 'relay', 'pwm', 'adc', 'servo', 'sensor')
- firmware_version: semver format

Permission: Tokkatot admin accounts only - NOT farmers
Notes:
- Farmers cannot use this endpoint
- UI does not expose this to farmers
- Tokkatot team uses this to configure farms
```

#### 27. Update Device (Name/Location Only)
```
PUT /farms/{farm_id}/devices/{device_id}
Authorization: Bearer {access_token}
Content-Type: application/json

Request:
{
  "name": "Updated Water Pump",
  "location": "Tank Area 2"
}

Response (200):
{
  "id": "device_uuid",
  "name": "Updated Water Pump",
  ...updated fields...
}

Permission: farmer
Notes:
- Farmers can ONLY update device name and location
- Device configuration/type cannot be changed by farmers
- Configuration changes must go through Tokkatot team
```

#### 28. Delete Device (Tokkatot Team Only)
```
DELETE /farms/{farm_id}/devices/{device_id}
Authorization: Bearer {access_token}

IMPORTANT: This endpoint is for Tokkatot team administrators only.
Farmers CANNOT delete devices.

Response (200):
{
  "message": "Device deleted successfully"
}

Permission: Tokkatot admin accounts only
Effect: Soft delete
Notes:
- Farmers cannot use this endpoint
- Device removal requires Tokkatot team coordination
```

#### 29. Get Device History
```
GET /farms/{farm_id}/devices/{device_id}/history?hours=24&limit=1000
Authorization: Bearer {access_token}

Query Parameters:
- hours: 1-720 (default 24)
- limit: 1-10000 (default 1000)
- metric: 'temperature'|'humidity'|'all' (default 'all')

Response (200):
{
  "device_id": "device_uuid",
  "readings": [
    {
      "timestamp": "2026-02-19T10:30:00Z",
      "sensor_type": "temperature",
      "value": 28.5,
      "unit": "°C",
      "quality": "good"
    }
  ],
  "total": 48,
  "from_timestamp": "2026-02-18T10:30:00Z",
  "to_timestamp": "2026-02-19T10:30:00Z"
}

Data Source: `device_readings` table (PostgreSQL)
```

#### 30. Calibrate Device
```
POST /farms/{farm_id}/devices/{device_id}/calibrate
Authorization: Bearer {access_token}
Content-Type: application/json

Request:
{
  "parameter_name": "temperature_offset",
  "parameter_value": 2.5,
  "unit": "°C"
}

Response (200):
{
  "parameter_name": "temperature_offset",
  "parameter_value": 2.5,
  "unit": "°C",
  "is_calibrated": true,
  "calibrated_at": "2026-02-19T10:30:00Z"
}

Permission: farmer
```

#### 31. Get Device Configuration
```
GET /farms/{farm_id}/devices/{device_id}/config
Authorization: Bearer {access_token}

Response (200):
{
  "device_id": "device_uuid",
  "configurations": [
    {
      "id": "config_uuid",
      "parameter_name": "temperature_threshold_high",
      "parameter_value": 35,
      "unit": "°C",
      "min_value": 20,
      "max_value": 50,
      "is_calibrated": true
    }
  ]
}
```

#### 32. Update Device Configuration
```
PUT /farms/{farm_id}/devices/{device_id}/config
Authorization: Bearer {access_token}
Content-Type: application/json

Request:
{
  "configurations": [
    {
      "parameter_name": "temperature_threshold_high",
      "parameter_value": 36,
      "unit": "°C"
    }
  ]
}

Response (200):
{
  "device_id": "device_uuid",
  "configurations": [...updated configs...]
}

Permission: farmer
Effect: Queues command to device via MQTT
```

#### 33. Get Device Firmware Update Status
```
GET /farms/{farm_id}/devices/{device_id}/firmware
Authorization: Bearer {access_token}

Response (200):
{
  "device_id": "device_uuid",
  "current_version": "1.2.3",
  "available_version": "1.3.0",
  "update_available": true,
  "release_notes": "Bug fixes and performance improvements",
  "update_status": "idle",
  "last_checked": "2026-02-19T10:30:00Z"
}
```

#### 34. Trigger Firmware Update
```
POST /farms/{farm_id}/devices/{device_id}/firmware/update
Authorization: Bearer {access_token}
Content-Type: application/json

Request:
{
  "target_version": "1.3.0",
  "schedule_at": "2026-02-20T02:00:00Z"  // optional
}

Response (202):
{
  "update_id": "uuid",
  "device_id": "device_uuid",
  "target_version": "1.3.0",
  "status": "scheduled",
  "scheduled_at": "2026-02-20T02:00:00Z"
}

Permission: farmer
Effect: Sends OTA update command to device
```

---

### Device Control Endpoints (8)

#### 35. Send Command to Device
```
POST /farms/{farm_id}/devices/{device_id}/commands
Authorization: Bearer {access_token}
Content-Type: application/json

Request:
{
  "command_type": "on",
  "command_params": {},
  "priority": "normal"
}

Response (202):
{
  "command_id": "uuid",
  "device_id": "device_uuid",
  "command_type": "on",
  "status": "queued",
  "requested_by": "user_uuid",
  "requested_at": "2026-02-19T10:30:00Z"
}

Command Types:
- 'on': Turn device on
- 'off': Turn device off
- 'toggle': Toggle state
- 'set_value': Set analog value {value: 0-255}
- 'set_angle': Set servo angle {angle: 0-180}
- 'set_speed': Set motor speed {speed: 0-100}

Permission: farmer
Priority: 'high'|'normal'|'low' (high priority commands skip queue)
```

#### 36. Get Command Status
```
GET /farms/{farm_id}/devices/{device_id}/commands/{command_id}
Authorization: Bearer {access_token}

Response (200):
{
  "command_id": "uuid",
  "device_id": "device_uuid",
  "command_type": "on",
  "status": "completed",
  "requested_by": "user_uuid",
  "requested_at": "2026-02-19T10:30:00Z",
  "sent_at": "2026-02-19T10:30:01Z",
  "completed_at": "2026-02-19T10:30:02Z",
  "device_response": {
    "status": "success",
    "message": "Device turned on"
  }
}

Status Values:
- 'queued': Awaiting transmission
- 'sent': Transmitted to device
- 'executing': Device processing
- 'completed': Success
- 'failed': Command failed
- 'timeout': No device response
```

#### 37. List Device Commands
```
GET /farms/{farm_id}/devices/{device_id}/commands?limit=50&offset=0&status=all
Authorization: Bearer {access_token}

Query Parameters:
- limit: 1-100 (default 50)
- offset: pagination offset
- status: 'all'|'queued'|'completed'|'failed'

Response (200):
{
  "commands": [
    {
      "command_id": "uuid",
      "command_type": "on",
      "status": "completed",
      "requested_at": "2026-02-19T10:30:00Z"
    }
  ],
  "total": 156
}
```

#### 38. Cancel Command
```
DELETE /farms/{farm_id}/devices/{device_id}/commands/{command_id}
Authorization: Bearer {access_token}

Response (200):
{
  "message": "Command cancelled"
}

Effect: Only works if status is 'queued'
Errors: 409 if command already sent to device
```

#### 39. Get Command History (Batch View)
```
GET /farms/{farm_id}/commands?hours=24&device_id=optional&limit=100
Authorization: Bearer {access_token}

Response (200):
{
  "commands": [
    {
      "command_id": "uuid",
      "device_id": "device_uuid",
      "device_name": "Water Pump",
      "command_type": "on",
      "status": "completed",
      "requested_by": "user_uuid",
      "requested_at": "2026-02-19T10:30:00Z"
    }
  ],
  "total": 248
}
```

#### 40. Emergency Stop All Devices
```
POST /farms/{farm_id}/emergency-stop
Authorization: Bearer {access_token}

Response (200):
{
  "stopped_devices": 6,
  "message": "All devices stopped"
}

Permission: farmer
Effect: Immediately stops all active devices on farm
Actions: Sends high-priority STOP command to all devices
```

#### 41. Get Device Real-Time Status
```
GET /farms/{farm_id}/devices/{device_id}/status
Authorization: Bearer {access_token}

Response (200):
{
  "device_id": "device_uuid",
  "is_online": true,
  "last_heartbeat": "2026-02-19T10:29:55Z",
  "current_state": "on",
  "current_value": 28.5,
  "unit": "°C",
  "signal_strength": -65,
  "battery_level": null,
  "uptime_hours": 1248
}

Update Frequency: Real-time via WebSocket
Fallback: HTTP polling every 5 seconds
```

#### 42. Batch Device Control
```
POST /farms/{farm_id}/devices/batch-command
Authorization: Bearer {access_token}
Content-Type: application/json

Request:
{
  "device_ids": ["device_uuid_1", "device_uuid_2"],
  "command_type": "on",
  "command_params": {}
}

Response (202):
{
  "command_batch_id": "uuid",
  "commands": [
    {
      "command_id": "uuid",
      "device_id": "device_uuid",
      "status": "queued"
    }
  ],
  "total": 2
}

Permission: farmer
```

---

### Scheduling Endpoints (7)

#### 43. Create Schedule
```
POST /farms/{farm_id}/schedules
Authorization: Bearer {access_token}
Content-Type: application/json

Request (Simple time-based):
{
  "device_id": "device_uuid",
  "name": "Morning Water",
  "schedule_type": "time_based",
  "cron_expression": "0 6 * * *",
  "action": "on",
  "action_duration": 900,  // Auto-turn-off after 15 minutes (optional)
  "priority": 0,
  "is_enabled": true
}

Request (Multi-step pulse sequence - NEW in v2.0):
{
  "device_id": "feeder_uuid",
  "name": "Pulse Feeding - Morning",
  "schedule_type": "time_based",
  "cron_expression": "0 6 * * *",
  "action": "on",
  "action_sequence": [
    {"action": "ON", "duration": 30},   // Motor ON 30 seconds
    {"action": "OFF", "duration": 10},  // Pause 10 seconds
    {"action": "ON", "duration": 30},   // Motor ON 30 seconds
    {"action": "OFF", "duration": 10},  // Pause 10 seconds
    {"action": "OFF", "duration": 0}    // Stay OFF until next schedule
  ],
  "priority": 5,
  "is_enabled": true
}

Response (201):
{
  "id": "schedule_uuid",
  "device_id": "device_uuid",
  "name": "Morning Water",
  "schedule_type": "time_based",
  "cron_expression": "0 6 * * *",
  "action": "on",
  "action_duration": 900,
  "action_sequence": null,  // Or array if multi-step
  "next_execution": "2026-02-20T06:00:00Z",
  "created_by": "user_uuid",
  "created_at": "2026-02-19T10:30:00Z"
}

Schedule Types:
- 'time_based': Uses cron_expression (0 6 * * * = 6am daily) + optional action_duration OR action_sequence
- 'duration_based': Uses on_duration/off_duration (repeating pattern)
- 'condition_based': Uses condition_json (temp > 30)
- 'manual': User-triggered only

New Fields (v2.0):
- action_duration (int): Auto-turn-off after X seconds for simple schedules
- action_sequence (JSON array): Multi-step patterns for pulse operations (feeders, conveyors)

See: docs/AUTOMATION_USE_CASES.md for real farmer scenarios

Permission: farmer
```

#### 44. List Schedules
```
GET /farms/{farm_id}/schedules?device_id=optional&limit=50
Authorization: Bearer {access_token}

Response (200):
{
  "schedules": [
    {
      "id": "schedule_uuid",
      "device_id": "device_uuid",
      "device_name": "Water Pump",
      "name": "Morning Water",
      "schedule_type": "time_based",
      "cron_expression": "0 6 * * *",
      "action": "on",
      "action_duration": 900,
      "action_sequence": null,
      "is_enabled": true,
      "next_execution": "2026-02-20T06:00:00Z",
      "last_execution": "2026-02-19T06:00:00Z",
      "execution_count": 45
    },
    {
      "id": "schedule_uuid_2",
      "device_id": "feeder_uuid",
      "device_name": "Feeder 1",
      "name": "Pulse Feeding",
      "schedule_type": "time_based",
      "cron_expression": "0 6,12,18 * * *",
      "action": "on",
      "action_duration": null,
      "action_sequence": [{"action":"ON","duration":30},{"action":"OFF","duration":10},{"action":"ON","duration":30}],
      "is_enabled": true,
      "next_execution": "2026-02-20T06:00:00Z",
      "execution_count": 132
    }
  ],
  "total": 12
}
```

#### 45. Get Schedule Details
```
GET /farms/{farm_id}/schedules/{schedule_id}
Authorization: Bearer {access_token}

Response (200):
{
  "id": "schedule_uuid",
  "device_id": "device_uuid",
  "name": "Morning Water",
  "schedule_type": "time_based",
  "cron_expression": "0 6 * * *",
  "action": "on",
  "action_duration": 900,
  "action_sequence": null,
  "is_enabled": true,
  "priority": 0,
  "next_execution": "2026-02-20T06:00:00Z",
  "last_execution": "2026-02-19T06:00:00Z",
  "execution_count": 45,
  "created_by": "user_uuid",
  "created_at": "2026-01-15T10:00:00Z",
  "updated_at": "2026-02-19T10:30:00Z"
}
```

#### 46. Update Schedule
```
PUT /farms/{farm_id}/schedules/{schedule_id}
Authorization: Bearer {access_token}
Content-Type: application/json

Request (Update timing):
{
  "name": "Early Morning Water",
  "cron_expression": "0 5 * * *",
  "is_enabled": false
}

Request (Update to multi-step sequence):
{
  "action_sequence": [
    {"action": "ON", "duration": 45},
    {"action": "OFF", "duration": 15},
    {"action": "ON", "duration": 45}
  ]
}

Response (200):
{
  "id": "schedule_uuid",
  "name": "Early Morning Water",
  "cron_expression": "0 5 * * *",
  "action_sequence": [{"action":"ON","duration":45},{"action":"OFF","duration":15},{"action":"ON","duration":45}],
  "is_enabled": false,
  ...other fields...
}

Permission: farmer
```

#### 47. Delete Schedule
```
DELETE /farms/{farm_id}/schedules/{schedule_id}
Authorization: Bearer {access_token}

Response (200):
{
  "message": "Schedule deleted successfully"
}

Permission: farmer
Effect: Soft delete
```

#### 48. Get Schedule Execution History
```
GET /farms/{farm_id}/schedules/{schedule_id}/executions?limit=100&days=30
Authorization: Bearer {access_token}

Response (200):
{
  "executions": [
    {
      "id": "execution_uuid",
      "schedule_id": "schedule_uuid",
      "scheduled_time": "2026-02-19T06:00:00Z",
      "actual_execution_time": "2026-02-19T06:00:02Z",
      "status": "executed",
      "execution_duration_ms": 245,
      "device_response": {
        "status": "success"
      }
    }
  ],
  "total": 45,
  "success_rate": 97.8
}

Data Source: schedule_executions table
```

#### 49. Manually Execute Schedule
```
POST /farms/{farm_id}/schedules/{schedule_id}/execute-now
Authorization: Bearer {access_token}

Response (202):
{
  "execution_id": "uuid",
  "schedule_id": "schedule_uuid",
  "status": "queued",
  "message": "Schedule queued for immediate execution"
}

Permission: farmer
Effect: Skips queue, highest priority
```

---

### Monitoring & Alerts Endpoints (8)

#### 50. Get Farm Alerts
```
GET /farms/{farm_id}/alerts?limit=50&is_active=true&severity=all
Authorization: Bearer {access_token}

Query Parameters:
- limit: 1-100 (default 50)
- is_active: true|false|all
- severity: 'info'|'warning'|'critical'|'all'

Response (200):
{
  "alerts": [
    {
      "id": "alert_uuid",
      "device_id": "device_uuid",
      "device_name": "Temperature Sensor",
      "alert_type": "temp_high",
      "severity": "critical",
      "message": "Temperature: 38°C (threshold: 35°C)",
      "threshold_value": 35,
      "actual_value": 38,
      "is_active": true,
      "triggered_at": "2026-02-19T14:30:00Z",
      "acknowledged_by": null
    }
  ],
  "total": 3,
  "active_count": 3,
  "critical_count": 1
}
```

#### 51. Get Alert Details
```
GET /farms/{farm_id}/alerts/{alert_id}
Authorization: Bearer {access_token}

Response (200):
{
  "id": "alert_uuid",
  "farm_id": "farm_uuid",
  "device_id": "device_uuid",
  "alert_type": "temp_high",
  "severity": "critical",
  "message": "Temperature: 38°C (threshold: 35°C)",
  "threshold_value": 35,
  "actual_value": 38,
  "is_active": true,
  "triggered_at": "2026-02-19T14:30:00Z",
  "acknowledged_by": "user_uuid",
  "acknowledged_at": "2026-02-19T14:35:00Z",
  "resolved_at": null
}
```

#### 52. Acknowledge Alert
```
PUT /farms/{farm_id}/alerts/{alert_id}/acknowledge
Authorization: Bearer {access_token}

Response (200):
{
  "id": "alert_uuid",
  "is_active": false,
  "acknowledged_by": "user_uuid",
  "acknowledged_at": "2026-02-19T14:35:00Z"
}

Effect: User confirms they've seen the alert
```

#### 53. Create Alert Subscription
```
POST /users/alert-subscriptions
Authorization: Bearer {access_token}
Content-Type: application/json

Request:
{
  "alert_type": "temp_high",
  "channel": "push",
  "is_enabled": true,
  "quiet_hours_start": "20:00",
  "quiet_hours_end": "06:00"
}

Response (201):
{
  "id": "subscription_uuid",
  "user_id": "user_uuid",
  "alert_type": "temp_high",
  "channel": "push",
  "is_enabled": true
}

Channels:
- 'push': In-app push notification
- No SMS/email (farmer-centric approach)
```

#### 54. Get Alert Subscriptions
```
GET /users/alert-subscriptions
Authorization: Bearer {access_token}

Response (200):
{
  "subscriptions": [
    {
      "id": "subscription_uuid",
      "alert_type": "temp_high",
      "channel": "push",
      "is_enabled": true,
      "quiet_hours_start": "20:00",
      "quiet_hours_end": "06:00"
    }
  ]
}
```

#### 55. Update Alert Subscription
```
PUT /users/alert-subscriptions/{subscription_id}
Authorization: Bearer {access_token}
Content-Type: application/json

Request:
{
  "is_enabled": false,
  "quiet_hours_start": "21:00"
}

Response (200):
{
  "id": "subscription_uuid",
  ...updated fields...
}
```

#### 56. Delete Alert Subscription
```
DELETE /users/alert-subscriptions/{subscription_id}
Authorization: Bearer {access_token}

Response (200):
{
  "message": "Subscription deleted"
}
```

#### 57. Get Alert History (Time-Range)
```
GET /farms/{farm_id}/alerts/history?days=30&limit=200
Authorization: Bearer {access_token}

Response (200):
{
  "alerts": [
    {
      "id": "alert_uuid",
      "alert_type": "temp_high",
      "severity": "critical",
      "triggered_at": "2026-02-19T14:30:00Z",
      "resolved_at": "2026-02-19T15:45:00Z",
      "duration_minutes": 75
    }
  ],
  "total": 24,
  "critical_count": 3,
  "warning_count": 18,
  "info_count": 3
}

Data Source: alerts table with deleted_at IS NULL
```

#### 58. Get Dashboard Overview
```
GET /farms/{farm_id}/dashboard
Authorization: Bearer {access_token}

Response (200):
{
  "farm": {
    "id": "farm_uuid",
    "name": "Main Chicken Farm"
  },
  "device_status": {
    "total": 6,
    "online": 5,
    "offline": 1,
    "error": 0
  },
  "alerts": {
    "active": 2,
    "critical": 1,
    "warning": 1
  },
  "recent_events": [
    {
      "id": "uuid",
      "event_type": "device_online",
      "device_name": "Temperature Sensor",
      "timestamp": "2026-02-19T10:30:00Z"
    }
  ],
  "quick_stats": {
    "avg_temperature": 28.5,
    "avg_humidity": 65,
    "last_24h_commands": 48,
    "last_24h_alerts": 5
  }
}

Performance: Cached for 30 seconds
```

---

### Analytics & Reporting Endpoints (5)

#### 59. Get Device Metrics Report
```
GET /farms/{farm_id}/reports/device-metrics?device_id=uuid&from=2026-02-01&to=2026-02-19&metric=temperature
Authorization: Bearer {access_token}

Query Parameters:
- device_id: Required
- from/to: Date range (YYYY-MM-DD)
- metric: 'temperature'|'humidity'|'all'
- interval: '1h'|'1d'|'1w' (default '1d')

Response (200):
{
  "device_id": "device_uuid",
  "metric": "temperature",
  "from_date": "2026-02-01",
  "to_date": "2026-02-19",
  "data_points": [
    {
      "timestamp": "2026-02-01T00:00:00Z",
      "value": 28.5,
      "min": 26.2,
      "max": 31.8,
      "avg": 28.9
    }
  ],
  "summary": {
    "min": 24.5,
    "max": 35.2,
    "avg": 28.9,
    "stddev": 2.1
  }
}

Data Source: `device_readings` table (PostgreSQL)
```

#### 60. Get Device Usage Report
```
GET /farms/{farm_id}/reports/device-usage?device_id=uuid&days=30
Authorization: Bearer {access_token}

Response (200):
{
  "device_id": "device_uuid",
  "device_name": "Water Pump",
  "from_date": "2026-01-20",
  "to_date": "2026-02-19",
  "total_on_time_hours": 168,
  "total_cycles": 532,
  "avg_uptime_percent": 97.2,
  "reliability_score": 9.7,
  "on_time_by_day": [
    {
      "date": "2026-02-19",
      "hours": 8.5,
      "cycles": 45
    }
  ]
}

Data Source: `device_readings` table (PostgreSQL)
```

#### 61. Get Farm Performance Report
```
GET /farms/{farm_id}/reports/farm-performance?days=30
Authorization: Bearer {access_token}

Response (200):
{
  "farm_id": "farm_uuid",
  "period": "2026-01-20 to 2026-02-19",
  "device_health": {
    "total_devices": 6,
    "healthy": 5,
    "degraded": 1,
    "offline": 0
  },
  "alerts_triggered": 24,
  "critical_alerts": 3,
  "automation_efficiency": {
    "scheduled_commands": 1440,
    "successful": 1401,
    "success_rate": 97.3
  },
  "uptime_percent": 99.1
}
```

#### 62. Export Report (CSV/PDF)
```
GET /farms/{farm_id}/reports/export?type=device_metrics&format=csv&device_id=uuid&from=2026-02-01&to=2026-02-19
Authorization: Bearer {access_token}

Query Parameters:
- type: 'device_metrics'|'device_usage'|'farm_performance'
- format: 'csv'|'pdf' (default 'csv')
- device_id: For device-specific reports

Response (200):
- Content-Type: text/csv or application/pdf
- Body: File download

File Size Limit: 50MB
```

#### 63. Get Farm Event Log
```
GET /farms/{farm_id}/events?limit=100&offset=0&event_type=all&days=30
Authorization: Bearer {access_token}

Query Parameters:
- limit: 1-500 (default 100)
- offset: pagination
- event_type: 'device_online'|'command_executed'|'schedule_triggered'|'all'
- days: 1-365

Response (200):
{
  "events": [
    {
      "id": "uuid",
      "event_type": "device_online",
      "event_category": "device",
      "severity": "info",
      "device_id": "device_uuid",
      "device_name": "Temperature Sensor",
      "message": "Device came online",
      "timestamp": "2026-02-19T10:30:00Z"
    }
  ],
  "total": 1243
}

Data Source: event_logs table
```

### AI/Disease Detection Endpoints (3)

#### 64. Get AI Service Health
```
GET /ai/health
Authorization: Bearer {access_token}

Response (200):
{
  "status": "healthy",
  "model_loaded": true,
  "models": {
    "ensemble_model": "loaded",
    "efficientnetb0": "loaded",
    "densenet121": "loaded"
  },
  "timestamp": "2026-02-19T10:30:00Z",
  "version": "2.0.0"
}

Data Source: AI Service (PyTorch FastAPI)
Requires: User role >= Viewer
```

#### 65. Predict Disease from Image
```
POST /ai/predict
Authorization: Bearer {access_token}
Content-Type: multipart/form-data

Request:
- image: Image file (PNG/JPEG, max 5MB)
- device_id: UUID (optional - for device association)

Response (200):
{
  "prediction_id": "uuid",
  "disease": "Coccidiosis",
  "confidence": 0.98,
  "ensemble_confidence": 0.99,
  "models": {
    "efficientnetb0": {
      "disease": "Coccidiosis",
      "confidence": 0.97
    },
    "densenet121": {
      "disease": "Coccidiosis",
      "confidence": 0.99
    }
  },
  "recommendation": "Isolate affected birds and consult veterinarian",
  "treatment_options": [
    {
      "name": "Amprolium",
      "dosage": "0.0125% in drinking water",
      "duration": "5-7 days"
    }
  ],
  "timestamp": "2026-02-19T10:30:00Z",
  "processing_time_ms": 1200
}

Error Responses:
- 400 Bad Request: Invalid image file (not PNG/JPEG, corrupted, etc.)
- 413 Payload Too Large: Image exceeds 5MB limit
- 503 Service Unavailable: AI model not loaded or AI service down

Data Source: PyTorch Ensemble Model
Requires: User role >= Viewer
Rate Limit: 100 requests/minute
```

#### 66. Predict Disease with Detailed Scores
```
POST /ai/predict/detailed
Authorization: Bearer {access_token}
Content-Type: multipart/form-data

Request:
- image: Image file (PNG/JPEG, max 5MB)
- device_id: UUID (optional)

Response (200):
{
  "prediction_id": "uuid",
  "ensemble_result": {
    "disease": "Coccidiosis",
    "confidence": 0.99
  },
  "per_model_results": {
    "efficientnetb0": {
      "disease": "Coccidiosis",
      "confidence": 0.97,
      "score_per_class": {
        "Healthy": 0.01,
        "Coccidiosis": 0.97,
        "Newcastle_Disease": 0.01,
        "Salmonella": 0.01
      },
      "inference_time_ms": 620
    },
    "densenet121": {
      "disease": "Coccidiosis",
      "confidence": 0.99,
      "score_per_class": {
        "Healthy": 0.00,
        "Coccidiosis": 0.99,
        "Newcastle_Disease": 0.01,
        "Salmonella": 0.00
      },
      "inference_time_ms": 580
    }
  },
  "model_agreement": true,
  "recommendation": "High confidence diagnosis. Isolate affected birds immediately.",
  "risk_level": "critical",
  "timestamp": "2026-02-19T10:30:00Z",
  "total_processing_time_ms": 1200
}

Note: This endpoint provides detailed per-model scores for analysis and debugging.

Data Source: PyTorch Ensemble Model (detailed output)
Requires: farmer
Rate Limit: 50 requests/minute
```

---

## Real-Time Communication

### WebSocket Connections (Socket.io)

**Connection URL**: `wss://api.tokkatot.local/socket.io`

**Authentication**:
```javascript
io('wss://api.tokkatot.local/socket.io', {
  auth: {
    token: 'jwt_access_token'
  }
})
```

**Events - Server to Client**:

```
device:update
  Emitted when device state changes
  Payload: {device_id, is_online, current_state, is_online, last_heartbeat}

device:command-status
  Emitted when command status changes
  Payload: {command_id, status, response_data}

alert:triggered
  Emitted when new alert is triggered
  Payload: {alert_id, alert_type, severity, message, device_id}

alert:resolved
  Emitted when alert is resolved
  Payload: {alert_id, resolved_at}

sensor:reading
  Emitted every 3-5 seconds with latest sensor data
  Payload: {device_id, sensor_type, value, timestamp}

schedule:executed
  Emitted when schedule execution completes
  Payload: {schedule_id, device_id, status, timestamp}

notification:new
  Emitted when new notification arrives
  Payload: {notification_id, title, message, severity}
```

**Events - Client to Server**:

```
farm:subscribe
  Request: {farm_id}
  Response: {subscribed_to: farm_id}

farm:unsubscribe
  Request: {farm_id}

device:command
  Request: {farm_id, device_id, command_type, command_params}
  Response: {command_id, status}
```

**Room Structure**:
- `farm:{farm_id}` - All events for farm
- `device:{device_id}` - Device-specific events
- `user:{user_id}` - User-specific notifications

---

## Error Handling

### HTTP Status Codes

| Status | Meaning | Example |
|--------|---------|---------|
| 200 | Success | Resource retrieved |
| 201 | Created | New resource created |
| 202 | Accepted | Async operation queued |
| 204 | No Content | Successful deletion |
| 400 | Bad Request | Invalid JSON or parameters |
| 401 | Unauthorized | Missing or invalid token |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | Resource doesn't exist |
| 409 | Conflict | State conflict (e.g., device already exists) |
| 429 | Rate Limited | Too many requests |
| 500 | Server Error | Internal error |
| 503 | Unavailable | Service temporarily down |

### Error Response Format

```json
{
  "error": {
    "code": "DEVICE_NOT_FOUND",
    "message": "Device with ID abc123 not found",
    "details": {
      "device_id": "abc123",
      "farm_id": "farm456"
    },
    "timestamp": "2026-02-19T10:30:00Z",
    "request_id": "req_abc123xyz"
  }
}
```

### Common Error Codes

```
INVALID_CREDENTIALS - Email/password incorrect
UNAUTHORIZED - Token missing/expired/invalid
FORBIDDEN - User lacks permission for resource
DEVICE_NOT_FOUND - Device doesn't exist
FARM_NOT_FOUND - Farm doesn't exist
USER_NOT_FOUND - User account not found
INVALID_INPUT - Validation failed
DEVICE_OFFLINE - Device unreachable
COMMAND_TIMEOUT - Device didn't respond
DUPLICATE_RESOURCE - Resource already exists
RATE_LIMIT_EXCEEDED - Too many requests
INTERNAL_ERROR - Unexpected server error
```

---

## Rate Limiting

### Limits per Endpoint Type

| Endpoint Type | Requests/Minute | Burst |
|---------------|-----------------|-------|
| Authentication | 10 | 5 |
| Data Read | 300 | 50 |
| Device Control | 100 | 20 |
| AI Prediction | 100 | 20 |
| Configuration | 50 | 10 |
| Analytics/Reports | 30 | 5 |

### Rate Limit Headers

```
X-RateLimit-Limit: 300
X-RateLimit-Remaining: 287
X-RateLimit-Reset: 1708324860
```

**Response**: 429 Too Many Requests when limit exceeded

---

## Pagination

### Query Parameters

```
limit: 1-100 (default 20)
offset: 0-N (default 0)
```

### Response Format

```json
{
  "data": [...],
  "total": 500,
  "limit": 20,
  "offset": 0,
  "has_more": true
}
```

---

## API Versioning

**Current Version**: v1  
**Base URL**: `/v1`

**Deprecation Policy**:
- Version supported for 2 years
- 6-month deprecation notice before removal
- Headers for identifying deprecated endpoints:
  ```
  Deprecation: true
  Sunset: Wed, 21 September 2028 07:28:00 GMT
  Link: <https://api.tokkatot.local/v2/...>; rel="successor-version"
  ```

---

## Webhooks (Optional Integration)

For external integrations, webhooks can be configured:

```
POST /webhooks/create
Authorization: Bearer {access_token}

Request:
{
  "url": "https://external-service.com/webhook",
  "events": ["device:online", "alert:triggered", "command:completed"],
  "secret": "webhook_secret_key"
}

Response (201):
{
  "webhook_id": "uuid",
  "url": "https://external-service.com/webhook",
  "events": [...],
  "is_active": true
}
```

**Webhook Payload**:
```json
{
  "event": "alert:triggered",
  "timestamp": "2026-02-19T10:30:00Z",
  "data": {
    "alert_id": "uuid",
    "farm_id": "uuid",
    ...
  },
  "signature": "sha256_hash"
}
```

**Retry Policy**: 3 retries with exponential backoff (5s, 25s, 125s)

---

## Security Headers

All responses include:

```
Content-Security-Policy: default-src 'self'
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Strict-Transport-Security: max-age=31536000; includeSubDomains
```

---

## Performance Targets

| Metric | Target |
|--------|--------|
| P50 Response Time | < 200ms |
| P95 Response Time | < 1s |
| P99 Response Time | < 3s |
| Availability | 99.9% |
| Rate Limit: Burst | 50 requests |
| Rate Limit: Sustained | 300/minute |

---

## API Documentation

- **Interactive API Explorer**: `/api/docs` (Swagger UI)
- **OpenAPI Specification**: `/api/openapi.json`
- **Postman Collection**: [Download Link]

---

## Development & Testing

### Test Credentials

```
Email: test@tokkatot.local
Password: TestPassword123!
Farm ID: test_farm_uuid
```

### Mock Device Simulator

For testing without physical devices:
```
POST /dev/mock-device/register
Request: { farm_id, device_id, name }
```

---

**Version History**

| Version | Date | Changes |
|---------|------|---------|
| 2.0 | Feb 2026 | Initial production API specification |
| | | 63 endpoints across 9 categories |
| | | WebSocket real-time support |
| | | Rate limiting and error handling |

---

**Related Documents**:
- [SPECIFICATIONS_ARCHITECTURE.md](SPECIFICATIONS_ARCHITECTURE.md)
- [SPECIFICATIONS_DATABASE.md](OI_SPECIFICATIONS_DATABASE.md)
- [SPECIFICATIONS_SECURITY.md](SPECIFICATIONS_SECURITY.md)
- [SPECIFICATIONS_REQUIREMENTS.md](02_SPECIFICATIONS_REQUIREMENTS.md)
