# Tokkatot 2.0: Functional & Non-Functional Requirements

**Document Version**: 2.0-FarmerCentric  
**Last Updated**: February 2026  
**Status**: Final Specification

**Design Principle**: Optimized for elderly farmers in Cambodia with low digital literacy

---

## Overview

This document specifies all functional requirements (features to build) and non-functional requirements (performance, security, compatibility targets) for Tokkatot 2.0. All requirements prioritize simplicity and accessibility for farmers with limited technical knowledge.

---

## Functional Requirements

### FR1: User Management & Authentication (Farmer-Centric)

**FR1.1** - User Registration (Email OR Phone)
- Allow new users to create accounts with EITHER email OR phone number (farmer chooses)
- Flexible contact verification: email link OR SMS code (6 digits)
- Password requirement simplified: 8+ characters (no complexity rules needed for farmers)
- Name required for account creation
- Terms of Service acceptance required
- Maximum 3 registration attempts per contact per hour (rate limit)
- New users start as farm owner if first signup

**FR1.2** - User Login
- Support email/phone + password authentication
- Generate JWT tokens valid for 24 hours
- Option to stay logged in on device (extend session to 30 days)
- Failed login lockout: 5 failed attempts = 15 minute lockout
- Log all login attempts (success/failure) for audit
- Support multi-device sessions (logout from specific devices)
- **MFA NOT required** (kept simple for farmers)

**FR1.3** - Password Management
- Password reset via email or SMS code
- Reset codes expire after 30 minutes
- Password change must verify current password
- Automatic logout after 60 minutes of inactivity
- Option to logout all devices at once

**FR1.4** - Simplified Role System (Farmer-Centric)
- 3 roles ONLY: Owner, Manager, Viewer (no complexity)
- Role matrix:
  - **Owner**: Farm owner - full access, invite managers/viewers, manage all settings
  - **Manager**: Can control devices, create schedules, invite viewers
  - **Viewer**: Read-only access - see device status and data only
- No granular permissions, no custom roles (keep it simple)

**FR1.5** - User Profiles
- Profile information: name, language (Khmer/English toggle), timezone
- Avatar upload optional (max 5MB)
- Notification preferences per user (quiet hours for alerts)
- NO API key generation (not for farmers)

**FR1.6** - Activity Audit Logging
- Log all user actions: login, commands, schedule changes
- Timestamp, user ID, action type, IP address
- Retention: 5 years for compliance
- NO admin dashboard for farmers to view logs (done by Tokkatot team)
- Export for Tokkatot team analysis

---

### FR2: Device Management (Tokkatot Team Only)

**FR2.1** - Device Registration (Team Setup ONLY)
- **IMPORTANT: Farmers CANNOT register devices themselves**
- Device registration done by Tokkatot team technicians only
- Devices pre-configured before deployment to farm
- Device metadata: model, firmware version, serial number (set by team)
- Devices assigned to farm/location at setup time
- Farmers only view pre-configured devices

**FR2.2** - Device Monitoring
- Real-time connectivity status (online/offline/error)
- Last heartbeat timestamp
- Current firmware version
- Signal strength indicator (for WiFi)
- Device location/room assignment
- **Farmers CAN customize device name/location only** (e.g., "Water Tank Pump" → "Tank 2 Pump")
- NOT device type, model, hardware settings

**FR2.3** - Device Control (Farmer-Facing)
- Instant ON/OFF commands for each device
- Command acknowledgement from device
- Command timeout (30 seconds) with user notification
- Queuing of commands if device offline (sync when online)
- Command history per device
- Batch commands (multiple devices at once) - simple group control

**FR2.4** - Firmware Management (Team Managed)
- **Tokkatot team ONLY can trigger firmware updates**
- Automatic firmware version checking
- OTA (Over-The-Air) update mechanism
- Scheduled updates (2 AM Cambodia time by default)
- Automatic rollback if boot fails
- Device updates do NOT require farmer interaction
- Farmers receive notification when device reboots for update

**FR2.5** - Device Grouping
- Organize devices by farm/location/type
- Quick control for device groups
- Group-level status dashboard
- Custom device naming

**FR2.6** - Device Configuration
- Customizable device parameters per device
- Sensor calibration adjustments
- Alarm thresholds per device
- Power-saving mode (e.g., lower polling frequency)
- Backup/restore configuration

---

### FR3: Real-Time Control & Monitoring

**FR3.1** - Live Dashboard
- Weather app-style interface
- Current values displayed as gauges:
  - Temperature (°C) with min/max/average
  - Humidity (%) with trend
  - Other sensor readings as relevant
- Device status indicators (on/off/error)
- Multi-component layout (customizable per user preference)
- Auto-refresh: 2-5 second intervals

**FR3.2** - Real-Time Data Streaming
- Latest sensor data available at server
- WebSocket push to all connected clients
- Fallback to HTTP polling (3 second intervals) if WebSocket fails
- Data format: JSON with timestamp
- Latency target: ≤2 seconds from device to dashboard

**FR3.3** - Historical Data Visualization
- 24-hour view: hourly data with line chart
- 7-day view: daily data with line chart
- 30-day view: daily aggregated data
- Custom date range: select any start/end date
- Peak/min/average indicators on chart
- Zoom in/out on charts
- Export data as CSV/JSON

**FR3.4** - System Health Overview
- Quick status: 🟢 All Good / 🟡 Warnings / 🔴 Errors
- List of warnings (device offline, sensor malfunction, etc)
- List of errors (critical failures)
- Recommended actions for each alert

---

### FR4: Scheduling & Automation

**FR4.1** - Time-Based Scheduling
- Schedule format: Day/Time specification
  - Daily: 6:00 AM - 18:00 PM (example: feeder only 6am-6pm)
  - Weekly: Monday-Friday only, or custom days
  - Custom date ranges: Jan 1 - Feb 28 only
- Cron-like syntax support for advanced users
- Multiple schedules per device
- Schedule enable/disable toggle
- Schedule priority (if multiple schedules conflict)

**FR4.2** - Duration-Based Cycles
- ON duration: how long system runs (seconds/minutes/hours)
- OFF duration: how long to wait before next cycle
- Repeat count: 1x to unlimited
- Window: start/end times within which cycles repeat
- Example: Water pump - 5 min ON, 2 hour OFF, repeat 6 AM - 6 PM daily
- Cycle history and next execution time visible

**FR4.3** - Condition-Based Triggers
- Temperature thresholds: if > 30°C turn ON fan
- Humidity thresholds: if < 40% turn ON humidifier
- Multi-condition logic: (temp > 30) AND (humidity < 60) = turn on both
- Boolean operators: AND, OR, NOT
- Condition history: when was condition last met?

**FR4.4** - Schedule Management
- Create/edit/delete schedules with UI form
- Duplicate existing schedule
- Bulk edit multiple schedules
- Schedule conflict warnings
- Schedule dry-run: preview what would execute
- View upcoming schedules (next 7 days)
- Automatic schedule execution with confirmation in logs

**FR4.5** - Execution History
- Log each schedule execution with:
  - Scheduled time vs actual execution time
  - Command sent to device
  - Device response/result
  - Any errors or warnings
- Query history by device, schedule, date range
- Export execution history as CSV

---

### FR5: Data Logging & History

**FR5.1** - Event Logging
- All device state changes logged:
  - Timestamp, device ID, old state, new state
  - User who triggered change (if manual)
  - Schedule that triggered change (if automated)
- All user actions logged:
  - Login/logout, commands, settings changes
  - Schedule creation/modification
  - User invitations
- User identification and IP address per log
- Immutable log records (cannot be deleted)

**FR5.2** - Sensor Data Logging
- Continuous sensor readings per device
- Configurable collection interval: 30 sec, 1 min, 5 min, etc
- Data aggregation options:
  - Detailed: keep raw 1-minute readings for 30 days
  - Aggregated: hourly average for 1 year
  - Summarized: daily average for 5 years
- Automatic compression/archival after retention period
- Query sensor data by device/time range
- Statistical analysis: min, max, average, percentiles

**FR5.3** - Maintenance Events
- Equipment uptime tracking
- Usage counters (hours run per device)
- Component lifetime estimates
- Maintenance alerts based on usage
- Maintenance log entries (manual)
- Predicted replacement dates

**FR5.4** - Data Export
- Export data as CSV, JSON, XML
- Date range selection
- Device/farm filtering
- Scheduled exports (daily/weekly/monthly)
- Email export automatically
- Downloadable from web interface

---

### FR6: Notifications & Alerts (In-App Message Log Only)

**FR6.1** - Alert Triggers
- Temperature out of range
- Humidity out of range
- Device offline > 30 min
- Firmware update available
- Schedule failed 3x
- Low battery
- Connection errors (local only)
- Anomaly detection

**FR6.2** - Message Log
- All alerts logged to dashboard message center
- Alert severity: urgent (device critical), important (value changes), info (schedule events)
- Each message shows: title, timestamp, device name, alert details
- Unread indicator (red dot)
- Click to see full details and related device

**FR6.3** - Message Management
- Customize which alerts are enabled per user
- Set alert thresholds per alert type
- Mute/snooze alerts for 1 hour / 4 hours / 24 hours
- Message history visible in dashboard (90 days)
- Mark messages as read/unread
- Delete old messages

**FR6.4** - In-App Notifications Only
- Real-time dashboard push (WebSocket)
- No SMS, email, or push notifications
- Farmers check dashboard to see all alerts
- Messages stored in database for 2 years audit trail

---

### FR7: Device Components

**FR7.1** - Water Pump
- ON/OFF control
- Duration-based schedules (5-120 minute cycles)
- Flow rate monitoring (if device has flow meter)
- Pressure monitoring (if device has pressure sensor)
- History: how many times ran, total duration

**FR7.2** - Feeder
- ON/OFF control
- Duration control (1-60 minutes)
- Feed count history
- Jam detection alerts
- Automatic low feed alerts

**FR7.3** - Light
- ON/OFF control with scheduled dimming (if supported)
- Brightness level control (0-100%)
- Lighting schedule (6 AM - 6 PM standard, customizable)
- Power consumption tracking
- Bulb lifespan tracking

**FR7.4** - Fan
- ON/OFF control
- Speed control if multi-speed model
- Temperature-triggered auto-on (if sensor available)
- Runtime tracking
- Maintenance alerts

**FR7.5** - Heater
- ON/OFF control
- Temperature setpoint configuration
- Thermostat mode: manual or auto
- heating schedule (winter season typically)
- Energy consumption tracking
- Overheat protection

**FR7.6** - Conveyor (manure removal, etc)
- ON/OFF control
- Duration control
- Run count tracking
- Jammed status detection
- Maintenance schedules

---

## Non-Functional Requirements

### NFR1: Performance

| Metric | Target | Tolerance |
|--------|--------|-----------|
| Dashboard load time | < 2 seconds | 95th percentile on 4G |
| API response time | < 200ms | 95th percentile |
| Chart rendering | < 500ms | initial render |
| Real-time update latency | < 1 second | device to client |
| Schedule execution accuracy | ±1 second | from scheduled time |
| Device command latency | < 500ms | from user to device |
| Page transitions | < 300ms |  |
| Search response | < 500ms | for 5 years of data |
| File upload | Support 20 MB | simultaneous multi-file |
| Concurrent users per farm | Support 10+ | simultaneous |

### NFR2: Availability

- **Uptime Target**: 99.5% (max 3.6 hours downtime per month)
- **SLA**: 99.9% for critical services (auth, device control)
- **Graceful Degradation**: UI works offline with cached data
- **Local Fallback**: 72+ hours without cloud connectivity
- **Schedule continuation**: Schedules execute even if cloud down
- **Automatic Failover**: Database replicas, service instances

### NFR3: Scalability

**Startup Phase (Month 1-3)**
- 10-20 farms
- 50-100 devices
- 20-30 concurrent users
- 10K-50K API calls/day

**Growth Phase (Month 4-12)**
- 50-100 farms
- 500-1000 devices
- 100-200 concurrent users
- 100K-500K API calls/day

**Production Phase (Year 2+)**
- 500-5000 farms
- 5K-50K devices
- 1000-5000 concurrent users
- 1M-10M API calls/day

Architecture supports horizontal scaling via:
- Stateless service instances
- Message queue for async tasks
- Database read replicas
- Caching layer (Redis)
- CDN for static assets

### NFR4: Security

**Authentication**
- JWT tokens, 24-hour expiration
- Password: 8+ characters, mixed case, number, symbol
- Rate limiting: 5 failed login attempts = 15 min lockout
- Multi-factor authentication (MFA) optional for Admins

**Authorization**
- Role-based access control (4 roles)
- Granular permissions per role
- Principle of least privilege
- Regular permission audits

**Encryption**
- TLS 1.3+ for all communications
- Data at rest: encryption for sensitive fields
- Key management: AWS KMS or equivalent
- Certificate pinning for mobile app

**Data Protection**
- GDPR/CCPA compliance
- Data anonymization for exports
- Automated backups, 30-day retention
- Disaster recovery plan

**Monitoring**
- Security event logging
- Intrusion detection system (IDS)
- Regular security audits (quarterly)
- Penetration testing (yearly)

### NFR5: Compatibility

**Browsers**
- Chrome 80+ (primary)
- Safari 13+ (iOS)
- Firefox 75+
- Samsung Browser 12+
- Edge 80+

**Operating Systems**
- Android 6.0+ (1-2GB RAM minimum)
- iOS 13+
- Windows 10+
- Linux (backend deployments)

**Devices**
- ESP32 (primary IoT device)
- Raspberry Pi 4B (local agent)
- Generic MQTT clients
- Any device with MQTT support (future)

**Network**
- 4G LTE (primary)
- WiFi 802.11ac/ax
- Minimum bandwidth: 1 Mbps
- Support mobile with intermittent connectivity

### NFR6: Reliability

**Error Handling**
- Graceful error messages to users
- No data loss on system restart
- Automatic recovery from transient failures
- Retry mechanisms with exponential backoff

**Data Integrity**
- ACID compliance for transactions
- Foreign key constraints
- Data validation at input
- Write-ahead logging for durability

**Monitoring & Alerting**
- Monitor CPU, memory, disk usage
- Alert on service health degradation
- Health check endpoints for all services
- Error rate monitoring and alerting

### NFR7: Maintainability

**Code Quality**
- Standardized code style (linters)
- Unit test coverage: 70%+
- Integration test coverage: 50%+
- Code review before merge

**Documentation**
- API documentation (OpenAPI/Swagger)
- Architecture documentation
- Deployment procedures
- Troubleshooting guides

**Logging & Debugging**
- Structured logging (JSON format)
- Log aggregation (ELK Stack or Loki)
- Trace ID for request tracking
- Debug mode for development

### NFR8: Accessibility (WCAG AAA)

**Visual**
- Minimum contrast ratio: 7:1 (text), 3:1 (graphics)
- Font sizes: minimum 16px on mobile
- Color not the only indicator
- Icon + text labels for all buttons

**Mobile**
- Touch targets: minimum 48x48 pixels
- Orientation support: portrait and landscape
- Responsive design: mobile-first
- Function available in portrait and landscape

**Users with Low Literacy**
- Simple navigation: Facebook-style bottom tabs
- Large buttons: 60px+ for primary actions
- High contrast: white on dark backgrounds
- Minimal text: use icons + short labels
- Khmer language as primary (not English)

**Users with Low Vision**
- 24px+ base font size
- Typography scale: 24, 32, 40, 48px headers
- High contrast: 7:1 WCAG AAA
- No pure gray text (use black/dark colors)
- Zoom support to 200%

### NFR9: Localization

**Language Support**
- Primary: Khmer (all interface)
- Secondary: English (farm names, device names)
- Language toggle in settings
- Date/time format per locale

**Multilingual Database**
- Interface strings in translation table
- Dynamic language switching
- RTL support if needed for Khmer script (optional)
- Timezone support per farm

### NFR10: Farmer-Centric Design

**Target Users**
- Age 40-70 years (mostly elderly farmers)
- Low digital literacy (use Facebook but not much else)
- Poor eyesight (presbyopia common)
- Use 1-2GB RAM Android phones
- 4G network with frequent disruptions

**Design Principles**
- **Simplicity**: 3-5 main actions, nothing more
- **Affordance**: Buttons look clickable, obvious actions
- **Feedback**: Instant visible response to every action
- **Error Prevention**: Confirm destructive actions
- **Offline First**: Works without internet, syncs when available

### NFR11: Data Retention & Privacy

- **Detailed logs**: 30 days
- **Sensor data**: 5 years (aggregated by month)
- **Event logs**: 5 years (audit trail)
- **User data**: Deleted on request (GDPR/CCPA compliant)
- **Automated backups**: Daily, retained 30 days
- **Disaster recovery**: RPO < 1 hour, RTO < 4 hours

---

## Constraint Matrix

| Constraint | Value | Rationale |
|-----------|-------|-----------|
| Max devices per farm | 100 | Hardware limitation |
| Max schedules per device | 20 | UI/UX complexity |
| Max users per farm | 50 | Access control |
| Max data points per query | 100K | Memory/performance |
| Max file size upload | 20MB | Network efficiency |
| API rate limit | 1000 req/min | Fair usage |
| WebSocket connections/server | 10K | Memory constraint |
| Database connections | 500 | Connection pooling |
| Cache size | 10GB | Memory allocation |
| Log retention | 5 years | Storage cost |

---

## Compliance & Standards

- **IoT Security**: OWASP IoT Top 10
- **Data Protection**: GDPR, CCPA
- **Accessibility**: WCAG 2.1 Level AAA
- **API**: OpenAPI 3.0 specification
- **Code Quality**: Linting with industry standards
- **Testing**: TDD/BDD where practical

---

**Version History**
| Version | Date | Changes |
|---------|------|---------|
| 2.0 | Feb 2026 | Initial production requirements |

**Related Documents**
- SPECIFICATIONS_ARCHITECTURE.md
- SPECIFICATIONS_API.md
- SPECIFICATIONS_FRONTEND.md
- SPECIFICATIONS_DATABASE.md
