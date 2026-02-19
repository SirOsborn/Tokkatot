# Tokkatot 2.0: Farmer-Centric Features & Technical Specifications
## Supplementary Document - Critical Additions for Production Release

**Version**: 2.0  
**Date**: February 18, 2026  
**Purpose**: Address specific farmer requirements and synchronization strategy

---

## 1. Remote Patching & Over-The-Air (OTA) Update Strategy

### Problem Statement
Cannot travel to each farm for software updates. Need automated patching mechanism that requires zero farmer intervention.

### Solution Architecture

#### Web App Updates
- Service worker checks for new version at midnight
- Downloads updated bundle in background
- Next app load uses new version
- **Farmer Experience**: Transparent - no action required, no downtime

#### Embedded System Updates (ESP32 + Raspberry Pi)

**Update Flow**:
```
1. Developer → GitHub push
2. CI/CD → Build & sign firmware
3. Upload to S3 with version number
4. Backend API publishes version
5. Device checks for update daily
6. Download at 2-4 AM (off-peak, low bandwidth usage)
7. Verify signature (security check)
8. Flash to OTA partition (safe, can rollback)
9. Device restarts
10. Report status to cloud
11. If offline during update window: will update when internet returns
12. Automatic rollback if new firmware fails to boot
```

#### Update Schedule

| Component | When | Mechanism | Downtime | Notification |
|-----------|------|-----------|----------|--------------|
| Web App | Nightly (midnight) | Service worker cache | ZERO (transparent) | None |
| ESP32 Firmware | Daily (2-4 AM) | MQTT trigger | 2-5 min | None (silent) |
| RPi Local Agent | Daily (3-5 AM) | systemd timer | 30 seconds | None |
| Critical Patches | ASAP | Manual trigger | 5 min max | Alert 1h before |

### API Endpoints for Update Management

```
GET /devices/{device_id}/firmware/latest
  - Check if new firmware available
  - Return: {version, url, signature, size}

GET /devices/{device_id}/firmware/download
  - Download firmware binary
  - Return: Signed binary + checksum

POST /devices/{device_id}/firmware/status
  - Device reports update result
  - Payload: {status: success/failed/rolled_back, version}

GET /web-app/latest-build
  - Web app manifest with new bundle location
  - Return: {js_files, hash, version}
```

### Automatic Rollback Mechanism

If new firmware fails to boot:
1. Watchdog timer triggers (10 seconds)
2. Device reboots
3. OTA bootloader detects previous boot failed
4. Automatically boots previous working version
5. Reports failure to cloud
6. Alert sent to ops team
7. User sees no disruption (app still works)

---

## 2. Seamless App + Embedded System Synchronization

### Problem from Prototype
Last build had synchronization issues causing out-of-sync states. Farmers confused when app showed different status than device.

### Solution: Three-Level Sync Strategy

#### Level 1: Real-Time Feedback (≤ 1 second)

**Scenario: User clicks "Turn ON Water Pump"**

```
Time Event
────────────────────────────────────────
0ms  User taps ON button
10ms App shows optimistic UI: "Water Pump ON"  
15ms App sends: POST /components/1/control {action: on}
50ms Backend sends MQTT command to device
     MQTT: tokkatot/device/esp32-001/commands
     Payload: {component_id: 1, action: on, timestamp: 1234567890}

100ms Device receives MQTT message
150ms Device executes: relay_set(GPIO_12, HIGH)
200ms Device sends confirmation MQTT:
      Topic: tokkatot/device/esp32-001/state_change
      Payload: {component_id: 1, state: on, timestamp: 1234567890}

250ms Backend receives confirmation, logs to database
260ms Backend sends WebSocket to all connected clients:
      {event: "component_state_change", component_id: 1, state: on}

300ms All clients (users viewing dashboard) see real-time update

TIMEOUT CASE (Device doesn't respond within 2 seconds):
2000ms Backend: No response from device
2010ms Backend sends WebSocket error to client
2020ms App shows warning: "⚠️ Device offline or busy"
2030ms App reverts UI to previous state
2040ms User can retry or troubleshoot
```

#### Level 2: Periodic Verification (Every 10 seconds)

Backend performs continuous verification:

```python
every 10 seconds:
    for each device:
        query_device_status()  # Ask device for current state
        
        device_response = {
            component_states: [
                {id: 1, state: on, last_changed: 1234567890},
                {id: 2, state: off, last_changed: 1234567889},
            ],
            timestamp: 1234567999
        }
        
        # Compare to database
        for each component:
            if device_state != db_state:
                # Trust device (source of truth)
                update_database(device_state)
                log_sync_discrepancy()
                broadcast_to_clients(actual_state)
```

**Mismatch Detection Algorithm**:
```
IF device says ON but DB says OFF:
    reason = unexpected_state_change
    source = device_is_authoritative
    action = update_db_and_notify_users
    alert_level = INFO
    
IF device says OFF but DB says ON:
    reason = device_stopped_unexpectedly  
    source = device_is_authoritative
    action = update_db_and_notify
    alert_level = WARNING
```

#### Level 3: Hourly Full Audit

```
Hourly reconciliation process:
1. Load all state change events from database (past hour)
2. Query device: Send full event log request
3. Device responds with: All state changes it recorded
4. Compare sequences:
   - Event by event
   - Check timestamps match (within ±2 seconds)
   - Verify no missing events
   - Detect duplicate events
5. If discrepancies found:
   - Alert ops team
   - Generate quality report
   - Schedule data reconciliation
   - Device logs take priority for ground truth
```

### Conflict Resolution Rules

| Scenario | Action | Priority | Reason |
|----------|--------|----------|--------|
| **App says OFF, Device says ON** | Accept device state | Device | Device is physical source of truth |
| **Cloud DB says 3 min runtime, Device says 5 min** | Use device time | Device | Device has accurate system time |
| **Multiple users issue conflict** | Use timestamp | Latest | Most recent action wins |
| **Network delayed command** | Check device timestamp first | Device | Verify timing before accepting |
| **User clicks OFF twice** | No-op second click | UI | Show "already OFF" feedback |
| **Schedule triggers but device offline** | Local device executes | Device | Schedules stored locally on device |

### Offline Mode Synchronization

**Device Offline**:
```
1. User opens app on phone (with local RPi connection)
2. Taps "Turn ON Water Pump"
3. App sends to local MQTT broker on RPi
4. RPi executes immediately (no cloud latency)
5. RPi queues the action: {queued_actions: [...]}
6. Device status shows immediately to user
7. When cloud internet returns:
   - RPi sends queued events to cloud
   - Cloud verifies sequence and timestamps
   - Updates if any conflicts
   - Confirms all users  
8. User sees: ✅ "Synced at 2:34 PM"
```

**App Offline**:
```
1. User actions cached in IndexedDB
2. UI shows: "Queued - waiting for connection"
3. When internet returns:
   - App sends queued commands in order
   - Server confirms execution
   - App updates UI with actual result
   - If mismatch: Show warning "Device state different, refresh?"
4. Automatic retry: Exponential backoff
```

### Sync Error Handling

| Error Type | User Experience | Recovery |
|------------|-----------------|----------|
| Device timeout (> 2s) | Button grays out, shows "⚠️ Device offline" | Retry automatically in 5 seconds |
| Network latency | "Sending..." loading spinner | Auto-retry with exponential backoff |
| Device error | Shows actual error from device (e.g. "Jam detected") | Display error code, suggest restart |
| State mismatch | Alert: "Device state different" | [Sync Now] button to reconcile |
| Schedule missed | Notification + reason shown | Query device logs for explanation |
| Command queued (offline) | "Queued - will send when online" | Automatic send when connectivity returns |

---

## 3. Performance Optimization for Low-End Mobile Devices

### Target Devices

- **CPU**: Single/dual core, 1.2-1.8 GHz
- **RAM**: 1-2 GB (may have other apps running)
- **Storage**: 32-64 GB (may be full)
- **Display**: 5-5.5" at 720p resolution
- **Network**: 4G LTE only (no 5G)
- **Browser**: Stock Chrome, Samsung Browser
- **OS**: Android 6-10 (not latest)
- **Users**: Elderly, poor eyesight, limited digital literacy

### Load Time Targets

```
FCP (First Contentful Paint):     < 2 seconds on 4G
TTI (Time to Interactive):        < 4 seconds  
DOM Ready:                        < 3 seconds
Total Bundle Size (gzipped):      < 150 KB
Memory Peak (on 1GB device):      < 50 MB
First Interaction Response Time:  < 100 ms
```

### Frontend Performance Checklist

**Bundle Size Optimization**:
- [ ] Vue.js 3 minified: < 60 KB gzipped (vs React >130 KB)
- [ ] Chart.js for graphs: < 30 KB gzipped
- [ ] Custom CSS: < 30 KB gzipped
- [ ] Total JS: < 150 KB gzipped
- [ ] Zero images on app (emoji only)
- [ ] Khmer font subset: 40 KB WOFF2

**Load Time Optimization**:
- [ ] Lazy load pages: Data & Settings load on demand
- [ ] Service worker pre-caches app shell on first load
- [ ] API responses cached: 5 min for sensor data, 1 hour for config
- [ ] Precompile Khmer font to only include 1000 most common characters
- [ ] No render-blocking resources above the fold
- [ ] Critical CSS inline in HTML

**Runtime Performance**:
- [ ] NO animations on control buttons (just color change)
- [ ] Virtual scrolling for long lists (Activity log, device list)
- [ ] Debounced search input (minimum 500 ms)
- [ ] No heavy computation on main thread
- [ ] WebSocket keepalive to prevent reconnects
- [ ] Connection pooling for HTTP requests

**Network Optimization**:
- [ ] GZIP compression on all text responses
- [ ] API response payloads < 50 KB each
- [ ] Batch multiple data queries into single API call
- [ ] MQTT QoS level 1 (balance reliability vs overhead)
- [ ] WebSocket binary mode when possible
- [ ] Image lazy loading with lqip (low quality image placeholder)

**Memory Optimization**:
- [ ] Limit DOM nodes visible in viewport (< 500)
- [ ] Clear event listeners when leaving page
- [ ] Data-driven templates (no DOM duplication)
- [ ] Memory profile < 50 MB on 1 GB device
- [ ] Monitor for memory leaks in WebSocket handlers

### Embedded System Performance

**ESP32 Optimization**:
- Core 0: Device control operations (highest priority, responsive)
- Core 1: Background tasks (sensor reading, WiFi, scheduling)
- Dedicated task priorities to prevent priority inversion
- Minimal JSON payloads (no pretty printing)
- MQTT QoS 1 for reliability without overhead

**Raspberry Pi Local Agent**:
- Keep total RAM usage < 256 MB
- Use local MQTT broker (non-blocking)
- SQLite for local cache (not PostgreSQL)
- Cron jobs for periodic cloud sync
- Memory monitoring: Restart if > 400 MB (safety limit)

### Network Resilience

```python
# Retry logic for unstable 4G
def send_command_with_retry(command):
    for attempt in range(1, 6):  # Up to 5 attempts
        try:
            result = send(command, timeout=5s)
            return result  # Success
        except Timeout:
            wait = exponential_backoff(attempt)  # 1s, 2s, 4s, 8s, 16s
            if attempt < 5:
                retry(after=wait)
        except Error:
            if recoverable:
                retry(after=exponential_backoff(attempt))
            else:
                raise
    
    # If all retries fail: Queue action locally
    queue_action_locally(command)
    return "Queued - will send when online"
```

---

## 4. Navigation UI Design (Facebook-Style)

### Bottom Navigation Bar

**Position**: Fixed at bottom of screen (sticky), never scrolls away

**Visual**:
```
┌─────────────────────────────────────────┐
│                                         │
│          Main Page Content              │
│                                         │
├─────────────────────────────────────────┤
│  🏠         📊         🔔        ⚙️       👤  │
│ Home      Data      Alerts    Settings Profile │
│                                         │
└─────────────────────────────────────────┘
```

### Tab Order & Priority

1. **🏠 Home** (Control)
   - **PRIMARY LANDING PAGE**
   - Farm system controls (most important)
   - Real-time status
   - Quick actions
   - RED badge for errors

2. **📊 Data** (Analytics)
   - Historical trends
   - Sensor charts
   - 24h / 7d / 30d views
   - Export functionality

3. **🔔 Alerts** (Notifications)
   - System events
   - Alerts and warnings
   - RED badge for unread
   - History of all events

4. **⚙️ Settings** (Configuration)
   - Schedules (create/edit)
   - Farm settings
   - Language selection
   - Theme (light/dark mode)
   - Notification preferences

5. **👤 Profile** (Account)
   - User information
   - Logout button
   - Account security

### Tab Design Specifications

**Active Tab Indicator**:
- Underline or highlight (2-3 px thick)
- Use primary color (#0D7377)
- Text bold when active

**Badge Notifications**:
- Red circular badge (🔴)
- White number inside
- Show count of unread/errors
- Remove when all viewed

**Size & Spacing**:
- Each tab: 60-80px wide
- Icon: 24x24 px
- Text: 12px
- Padding: 8-12 px

**Touch Target**:
- Minimum 44px height (better: 48px)
- Full touch area for each tab

---

## 5. Home Page Design for Elderly Farmers

### Layout Philosophy

**Goal**: Simple, easy-to-understand interface with minimal cognitive load

**Principles**:
- Largest numbers up front (not buried)
- Color feedback (green = good, red = problem)
- Clear status text (not ambiguous icons)
- Large, easy-to-tap buttons
- One action per card
- Confirmation dialogs for risky actions
- Progress feedback for all operations

### Visual Hierarchy

```
┌─────────────────────────────────────────┐
│ Farm: Tokkatot #1      ⚠️ 2 Alerts      │ (Header)
├─────────────────────────────────────────┤
│                                         │
│  ╔═══════════════════════════════════╗ │
│  ║ 🌡️  TEMPERATURE                ║ │
│  ║                                 ║ │
│  ║      28.5 °C                    ║ │ (Large: 48px)
│  ║   Good (24°-29°C)               ║ │ (Green color)
│  ╚═══════════════════════════════════╝ │
│                                         │
│  ╔═══════════════════════════════════╗ │
│  ║ 💧 HUMIDITY                      ║ │
│  ║                                 ║ │
│  ║        65 %                     ║ │ (Large: 42px)
│  ║   Good (55%-75%)                ║ │ (Green color)
│  ╚═══════════════════════════════════╝ │
│                                         │
│ ─────────────────────────────────────── │
│         SYSTEM CONTROLS               │
│ ─────────────────────────────────────── │
│                                         │
│  ╔═══════════════════════════════════╗ │
│  ║ 💧 WATER PUMP                    ║ │
│  ║ Status: ● ON                     ║ │ (Green status)
│  ║ Running: 3:42 / 5:00 minutes     ║ │
│  ║                                 ║ │
│  ║  [ OFF ]      [ STOP ]           ║ │ (Large buttons: 48px)
│  ╚═══════════════════════════════════╝ │
│                                         │
│  ╔═══════════════════════════════════╗ │
│  ║ 🐔 FEEDER                        ║ │
│  ║ Status: ● OFF                    ║ │ (Gray status)
│  ║ Last fed: 2 hours ago            ║ │
│  ║                                 ║ │
│  ║  [ ON ]     [ FEED 500G ]        ║ │
│  ╚═══════════════════════════════════╝ │
│                                         │
│  [Scroll for more systems...]          │
│                                         │
├─────────────────────────────────────────┤
│ 🏠 Home  │  📊  │ 🔔  │ ⚙️  │ 👤      │
└─────────────────────────────────────────┘
```

### Component Card Design

**Each System Card**:
```
┌─────────────────────────────────────┐
│ [Icon] System Name                  │ (16px text)
│                                     │
│ Status: ● ON/OFF                    │ (Color indicator)
│ Info: Running 3:42 min / Last 2h    │ (Status info)
│                                     │
│ [ OFF ]  [ STOP/ACTION ]            │ (Large tap targets)
│                                     │
│ Yellow bar: 75% (if warning)        │
└─────────────────────────────────────┘
```

### Font & Readability

**Sizes**:
- Farm name header: 18-20px
- Section titles: 16px
- Component names: 16px
- Large numbers (temp): 48px
- Status text: 14-16px
- Button text: 16-18px

**Line Height**: 
- All text: 1.6 or higher

**Font Family**:
- System font (iOS: San Francisco, Android: Roboto)
- Khmer font fallback
- No decorative fonts

**Color Contrast**:
- Regular text: 7:1 minimum
- Disabled text: 4.5:1 minimum
- Recommended: 10:1 for elderly users

### Interaction Design

**Button Behavior**:
- Instant visual feedback (color change)
- NO animation (instant)
- Clear disabled state (grayed out)
- Confirmation required for risky actions

**Confirmation Dialogs**:
```
┌─────────────────────────────────┐
│ Confirm Action                  │
├─────────────────────────────────┤
│                                 │
│ Stop Water Pump?               │
│                                 │
│          [ Cancel ] [ OK ]      │
│                                 │
└─────────────────────────────────┘
```

**Feedback Messages**:
- ✅ "Water pump turned on"
- ⚠️ "Temperature too high"
- ❌ "Device offline"
- ⏳ "Sending command..."
- ❓ "Device not responding"

---

## 6. Multilingual Implementation (Khmer + English)

### Language Architecture

**Supported Languages**:
- 🇰🇭 **Khmer (ភាសាខ្មែរ)** - Primary (Cambodia default)
- 🇬🇧 **English** - Secondary (international expansion)

### Toggle UI

**Location**: Settings tab (⚙️) → Language section

**Design**:
```
┌─────────────────────────────────┐
│ LANGUAGE (ភាសា)                 │
├─────────────────────────────────┤
│                                 │
│  🇰🇭 Khmer (ភាសាខ្មែរ)        │
│  🇬🇧 English                   │
│                                 │
│  Selected: ● Khmer              │
│                                 │
└─────────────────────────────────┘
```

### Implementation Details

**Frontend i18n Framework**:
- i18next or similar
- All text in JSON translation files
- Client-side rendering (no API translation calls)
- Fallback: English if Khmer translation missing

**Translation File Structure**:
```json
{
  "khmer": {
    "controls": "ឧបករណ៍គ្រប់គ្រង",
    "temperature": "សីតុណ្ហភាព",
    "status_on": "បើក",
    "status_off": "បិទ",
    "alerts": "ការជូនដំណឹង",
    "settings": "ការកំណត់",
    "water_pump": "ម៉ាស៊ីនបូមទឹក",
    "feeder": "ឧបករណ៍급여",
    "good_range": "ជួរល្អ",
    "warning": "ប្រយ័ត្ន",
    "error": "កំហុស"
  },
  "english": {
    "controls": "Controls",
    "temperature": "Temperature",
    "status_on": "On",
    "status_off": "Off",
    ...
  }
}
```

**Khmer Font Support**:
- Embedded font file: Khmer UI or Khmer OS (40 KB WOFF2)
- Character subset: 1000 most common Khmer characters
- No web font API requests (save bandwidth)
- Fallback: System Khmer font if available

**Persistence**:
- Save language choice in user profile on cloud
- Also save in localStorage (for offline)
- Applied immediately (no page reload needed)

### System Messages

**Principles**:
- Simple, action-oriented language
- Avoid idioms or phrases
- Use emoji for visual context
- One idea per message

**Examples**:

Khmer:
- "⚠️ សីតុណ្ហភាពលើសពីរង្វាន់" (Temperature exceeds limit)
- "✅ ម៉ាស៊ីនបូមទឹកបើក" (Water pump turned on)
- "❌ ឧបករណ៍ออフលាញ" (Device offline)
- "⏳ ដំណើរការផ្ញើទិន្នន័យ..." (Sending data...)

English:
- "⚠️ Temperature too high"
- "✅ Water pump turned on"
- "❌ Device offline"
- "⏳ Sending data..."

### API & Backend Localization

**Strategy**: 
- Backend returns language codes, NOT translated text
- Frontend handles all translation

**Example API Response**:
```json
{
  "event": {
    "type": "error",
    "code": "DEVICE_OFFLINE",
    "timestamp": "2026-02-18T14:35:22Z"
  }
}
```

**Frontend Translation**:
```javascript
const message = i18n.t("errors.DEVICE_OFFLINE")
// 🇰🇭 Returns: "ឧបករណ៍ផ្តាច់ពីបណ្តាញ"
// 🇬🇧 Returns: "Device offline"
```

---

## Database & Synchronization Technical Specs

### Real-Time Event Stream

**Event Types Captured**:
1. `SENSOR_READING` - New sensor measurement
2. `STATE_CHANGE` - Device ON/OFF transition
3. `COMMAND_ISSUED` - User command received
4. `COMMAND_EXECUTED` - Device executed command
5. `SCHEDULE_TRIGGERED` - Automatic schedule fired
6. `ERROR_OCCURRED` - System or device error
7. `USER_ACTION` - User login/button click
8. `SYNC_COMPLETED` - Cloud-device sync finished

**Event Storage Requirements**:
- Immutable (never update/delete after creation)
- Microsecond precision timestamps
- Distributed across time-series database
- Partitioned by date for query efficiency
- 5-year retention for compliance

### State Change Validation

**Before Accepting Any State Change**:
1. Verify device ID exists in system
2. Verify device is in expected state
3. Validate timestamp (> previous change)
4. Check valid state transition
5. Detect and reject duplicates
6. Verify digital signature (if from device)
7. Log with reason (manual/schedule/error)

**Rejection Criteria**:
- Stale timestamp (> 1 minute old)
- Invalid device ID
- Out-of-sequence events
- Duplicate detection
- Blacklisted user/device

---

## Success Metrics for Production

**Farmer Satisfaction**:
- ✅ First-time user can control farm within 2 minutes
- ✅ 95% uptime for critical operations
- ✅ Zero data loss events
- ✅ Average response time < 500ms

**Technical Metrics**:
- ✅ App load time < 2 seconds (4G)
- ✅ State sync latency < 1 second
- ✅ 99.9% message delivery (MQTT)
- ✅ Zero unplanned downtime after launch

**Update Metrics**:
- ✅ 100% auto-update success rate
- ✅ Zero app crashes after auto-update
- ✅ Firmware updates complete without farmer interaction
- ✅ Rollback triggers < 0.1% of updates

---

**Document Status**: Complete Technical Specification for Implementation  
**Version**: 2.0  
**Last Updated**: February 18, 2026  
**Ready for**: Development Team Review & Project Kickoff

---
