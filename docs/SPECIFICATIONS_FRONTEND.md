# Tokkatot 2.0: Frontend Specification

**Document Version**: 2.0  
**Last Updated**: February 2026  
**Status**: Final Specification

---

## Overview

Tokkatot 2.0 frontend is designed for elderly farmers with low digital literacy, using low-end smartphones (1-2GB RAM) with 4G connectivity. The UI prioritizes accessibility, simplicity, and offline-first design.

---

## Technology Stack

**Framework**: Vue.js 3 (Composition API)  
**Build Tool**: Vite  
**State Management**: Pinia  
**Component Library**: None (custom components for control)  
**Styling**: TailwindCSS  
**Charts**: Chart.js  
**Real-time**: Socket.io (fallback to polling)  
**Offline Support**: Service Workers  
**Storage**: IndexedDB (local caching)  

---

## Design Principles

### For Elderly Farmers

1. **Simplicity**: 5 main actions maximum per screen
2. **Large Touch Targets**: 48px+ for buttons
3. **High Contrast**: 7:1 WCAG AAA compliance
4. **Clear Feedback**: Immediate visible response to every action
5. **Forgiving**: Undo capability, confirmation on destructive actions
6. **Familiar Navigation**: Facebook-style layout they already know

### Performance Targets

- **Bundle Size**: < 150KB (gzipped)
- **Initial Load**: < 2 seconds on 4G
- **Page Transitions**: < 300ms
- **Real-time Updates**: < 1 second latency
- **Offline Support**: Works with cached data
- **Zoom Support**: Up to 200% magnification

### Multilingual

- **Primary Language**: Khmer
- **Secondary Language**: English
- **Language Toggle**: Visible in settings
- **Date Format**: Per locale
- **Number Format**: Per locale

---

## Navigation Structure

### Bottom Tab Navigation (Facebook-Style)

```
Screen Layout:
┌───────────────────────────────┐
│                               │
│       MAIN CONTENT AREA       │
│       (scrollable)            │
│                               │
└───────────────────────────────┘
│ Home │ Schedules │ Logs │ Settings │
```

**5 Main Tabs**:

**1. Home Tab** (Primary Landing)
- Current status dashboard
- Real-time gauges (temperature, humidity)
- Quick action buttons (top 3 most used)
- Alert banner (if any critical alerts)
- "Control Now" sections for each device

**2. Control Tab** (Device Management)
- List of all devices with on/off toggles
- Organized by type (water, feeder, light, fan)
- Quick status indicators
- Group control (turn all light on/off)

**3. Schedules Tab** (Automation)
- List of active schedules
- Enable/disable toggles
- Next execution time visible
- "Create Schedule" button prominent
- Execution history

**4. Analytics Tab** (Historical Data)
- 24-hour, 7-day, 30-day chart tabs
- Temperature/humidity trends
- Peak time indicators
- Export button

**5. Settings Tab** (Configuration)
- Farm settings
- Language toggle
- Dark/light mode (optional)
- About/Help
- Logout

---

## Accessibility Requirements (WCAG AAA)

### Visual

**Typography**:
- Base font size: 16px (phones), rescalable to 24px
- Headers: 24px (h4), 32px (h3), 40px (h2), 48px (h1)
- Line height: 1.5 (clear spacing)
- Letter spacing: 0.5px for headers

**Color Contrast**:
- Text on background: 7:1 ratio (AAA standard)
- Use: Black (#000000) text on white (#FFFFFF)
- Or: Dark gray on light gray (adjust gray shades)
- Not: Pure gray text (insufficient contrast)
- Test: WebAIM Contrast Checker

**Color Usage**:
- Don't rely on color alone to convey meaning
- Use icons + color
- Example: Red icon + text "Error" (not just red background)

### Mobile

**Touch Targets**:
- Minimum: 48x48 pixels (AAA standard)
- Recommended: 60x60 pixels for critical actions
- Spacing: 8px minimum between touch targets

**Orientation**:
- Support portrait and landscape
- No forced orientation
- Responsive layout adapts to both

**Zoom Support**:
- Allow zoom to 200% without overflow
- Text remains readable at 200%
- No horizontal scroll at 200%

### Users with Low Vision

**Font Scaling**:
- Respect system font size settings
- Use `rem` units (scale with system setting)
- Support 200% magnification before horizontal scroll

**Visual Indicators**:
- Status indicators: Check mark, X, warning icon
- Add text labels + icons (not icons alone)

**Contrast Checker**:
All UI elements tested with WebAIM contrast tool:
- Text: 7:1 minimum
- UI components: 3:1 minimum
- Graphics: 3:1 minimum

---

## UI Component Specification

### Buttons

**Primary Button** (most common action)
```
Appearance:
- Size: 60x56px minimum
- Background: Dark color (#1F2937)
- Text: White (#FFFFFF)
- Font: Bold, 16px
- Border: None
- Border-radius: 8px
- Shadow: Subtle drop shadow

States:
- Default: As above
- Hover: Lighter background
- Active/Pressed: Darker background
- Disabled: Grayed out, cursor not-allowed
```

**Secondary Button** (less common)
```
Appearance:
- Size: 60x56px minimum
- Border: 2px solid dark
- Background: White
- Text: Dark color
- Font: 16px

States:
- Default: As above
- Hover: Light gray background
- Active: Dark text, light gray background
- Disabled: Grayed out
```

**Toggle (ON/OFF)**
```
Appearance:
- Size: 60pt height, variable width
- Background: Color change indicates state
- ON: Green (#10B981)
- OFF: Gray (#D1D5DB)
- Text: "ON" or "OFF" inside toggle
- Animation: Smooth 300ms transition
```

### Input Fields

**Text Input**
```
- Font: 16px
- Height: 48px
- Border: 2px solid (#E5E7EB)
- Border on focus: Blue (#3B82F6)
- Padding: 12px

Label (above field):
- Font: 16px
- Margin bottom: 8px
- Bold weight
```

**Select Dropdown**
```
- Font: 16px
- Height: 48px
- Shows current selection
- Dropdown opens on tap
- Each option: 48px tall (for touch)
```

**Number Input (for durations, thresholds)**
```
- Font: 16px
- Height: 48px
- Plus/minus buttons: 40x40px
- Keyboard: Numeric only
- Value displayed prominently
```

### Information Display

**Gauge (Current Value)**
```
Layout:
- Value: Large (40px text)
- Unit: Medium (20px text)
- Min/Max shown below
- Color: Green (normal), Yellow (warning), Red (alarm)

Example:
  ╭─────────────╮
  │    28°C     │  ← Value (40px)
  │ Temperature │  ← Label (16px)
  │ Min: 20 Max: 32 │ ← Range
  ╰─────────────╯
```

**Status Indicator**
```
- Icon (32x32px): Checkmark / X / Warning
- Text (16px): Status message
- Color: Green / Red / Orange
- Updates in real-time if connected

Example:
  ✓ Online    (green)
  ✗ Offline   (red)
  ⚠ Error     (orange)
```

**Chart (Time-Series Data)**
```
- Height: 200px minimum
- Width: Full screen width
- Labels: Bottom (time), Left (value + unit)
- Line thickness: 3px
- Tap to zoom/explore data
- Show legend (max 4 series)
```

---

## Page Layouts

### Home Page (Primary Landing)

```
┌─────────────────────────────┐
│  Farm: Neath's Farm  ✓   ⟳  │ ← Farm selector, refresh button
├─────────────────────────────┤
│                             │
│   TEMPERATURE GAUGE         │
│     ╭──────────────╮        │
│     │    28°C     │         │
│     │ Comfortable │         │
│     ╰──────────────╯        │
│                             │
│   QUICK ACTIONS             │
│   ┌──────────┬──────────┐   │
│   │ Water ON │Feeder ON │   │
│   └──────────┴──────────┘   │
│   ┌──────────┬──────────┐   │
│   │ Light ON │ Fan OFF  │   │
│   └──────────┴──────────┘   │
│                             │
│   ALERTS (if any)           │
│   ⚠ Temperature high (28°C) │
│                             │
└─────────────────────────────┘
 Home│Ctrl│Sched│Data│Setting
```

### Control Page

```
┌─────────────────────────────┐
│ FARM EQUIPMENT              │
├─────────────────────────────┤
│                             │
│ 💧 WATER                    │
│   ┌───────────────────────┐ │
│   │ Pump 1                │ │ ← Toggle
│   │ [  ⊙ ON ] Last: 2h ago │ │
│   └───────────────────────┘ │
│                             │
│ 🌾 FEEDER                   │
│   ┌───────────────────────┐ │
│   │ Feeder 1              │ │
│   │ [ • OFF ] Try now →   │ │
│   └───────────────────────┘ │
│                             │
│ 💡 LIGHTING                 │
│   ┌───────────────────────┐ │
│   │ Main Light            │ │
│   │ [  ⊙ ON ] Scheduled   │ │
│   └───────────────────────┘ │
│                             │
│ ❄ CLIMATE                   │
│   Temperature: 28°C (⟳)     │
│   [  • HEAT OFF    ]        │
│   [  ⊙ FAN ON     ]        │
│                             │
└─────────────────────────────┘
 Home│Ctrl│Sched│Data│Setting
```

### Schedules Page

```
┌─────────────────────────────┐
│ SCHEDULES                   │
├─────────────────────────────┤
│ [+ Create Schedule] button  │
├─────────────────────────────┤
│                             │
│ ✓ Water Pump (Enabled)      │
│   6:00 AM - 6:00 PM daily   │
│   Next run: Today 6:00 AM   │
│   [Edit] [Delete]           │
│                             │
│ • Feeder (Disabled)         │
│   9:00 AM - 5:00 PM daily   │
│   Next run: Disabled        │
│   [Edit] [Delete]           │
│                             │
│ ✓ Fan (if temp > 30°C)     │
│   Automatic mode            │
│   Next run: Check now       │
│   [Edit] [Delete]           │
│                             │
└─────────────────────────────┘
 Home│Ctrl│Sched│Data│Setting
```

### Data/Analytics Page

```
┌─────────────────────────────┐
│ TEMPERATURE ANALYTICS       │
├──────────────────────────────┤
│ [24h] [7d] [30d] [More ▼]  │
├──────────────────────────────┤
│                             │
│      CHART                  │
│      ╱╲      ╱╲            │
│     ╱  ╲    ╱  ╲           │
│    ╱    ╲  ╱    ╲          │
│ 0h    6h    12h    18h 24h  │
│                             │
│ 📊 Statistics               │
│ Max:  32°C @ 14:30          │
│ Min:  20°C @ 03:15          │
│ Avg:  26°C                  │
│                             │
│ [📥 Download CSV]           │
│                             │
└─────────────────────────────┘
 Home│Ctrl│Sched│Data│Setting
```

### Settings Page

```
┌─────────────────────────────┐
│ SETTINGS                    │
├─────────────────────────────┤
│                             │
│ FARM                        │
│   Farm Name: [text input]   │
│   Location: [text]          │
│   Timezone: [dropdown]      │
│                             │
│ INTERFACE                   │
│   Language: [Khmer ✓ English] │
│   Dark Mode: [OFF]          │
│   Text Size: [Normal ▼]     │
│                             │
│ NOTIFICATIONS (In-App Only) │
│   App Notifications: [ON]   │
│   Sound Alerts: [ON]        │
│   Quiet Hours: 8 PM - 6 AM  │
│                             │
│ ACCOUNT                     │
│   Email: farmer@example.com │
│   Password: [Change]        │
│   Sessions: [2 devices]     │
│   [Logout All]              │
│                             │
│ ABOUT                       │
│   Version: 2.0.0            │
│   Help & Support            │
│   [Logout]                  │
│                             │
└─────────────────────────────┘
 Home│Ctrl│Sched│Data│Setting
```

---

## Offline Support

**Service Worker Implementation**:
- Cache API responses during online
- Serve cached data when offline
- Queue commands for later sync
- Background sync when reconnected
- Show offline indicator in UI

**Cacheable Data**:
- Device states (updated every 30 seconds)
- Schedules (updated when modified)
- Historical data (cached charts)
- UI static assets (HTML, CSS, JS)

**Non-Cacheable**:
- Real-time sensor readings (> 30 seconds old)
- User commands (can queue)
- Logs and analytics (use cached if available)

---

## Real-Time Updates

**WebSocket Connection**:
- Auto-connect on app load
- Auto-reconnect with exponential backoff
- Fallback to polling if WebSocket fails
- Polling interval: 3 seconds on 4G

**Subscriptions**:
- Subscribe to device status
- Subscribe to device alerts
- Subscribe to schedule updates
- Unsubscribe when leaving page

**Update Format**:
```json
{
  "type": "device_status",
  "device_id": "esp32-001",
  "status": "on",
  "timestamp": "2026-02-18T14:30:00Z",
  "value": 28.5
}
```

---

## Performance Optimization

**Bundle Size Reduction**:
- Tree-shaking (remove unused code)
- Code splitting (lazy load routes)
- Image optimization (WebP format)
- CSS purging (remove unused styles)
- No unnecessary dependencies

**Target: < 150KB gzipped**

**Rendering Optimization**:
- Virtual scrolling (large lists)
- Debounce chart redraws
- Throttle resize handlers
- Minimize DOM updates
- Use CSS animations (not JS)

**Caching Strategy**:
- Browser cache: 1 year for static assets
- Service Worker cache: 30 days
- API response cache: 5 minutes (configurable)
- Local storage: User preferences

---

## Testing

**Browser Support**:
- Chrome 80+ (primary)
- Safari 13+ (iOS)
- Samsung Browser 12+
- Firefox 75+

**Device Testing**:
- Test on actual 1-2GB RAM phones
- Test on slow 4G networks (throttle to 1 Mbps)
- Test on unreliable networks (add packet loss)
- Test zoom/magnification up to 200%

**Test Scenarios**:
- Offline mode: no internet connection
- Slow network: 1 Mbps throttling
- Device offline: backend unreachable
- WebSocket failure: fallback to polling
- Theme switching: dark/light mode
- Language switching: instant translation
- Zoom levels: 100%, 150%, 200%

---

## Accessibility Testing Checklist

- [ ] Keyboard navigation works (all interactive elements)
- [ ] Screen reader text present (ARIA labels)
- [ ] Touch targets at least 48x48px
- [ ] Color contrast 7:1 (no pure gray text)
- [ ] Form labels present and associated
- [ ] Focus indicators visible
- [ ] Zoom support to 200%
- [ ] No timeouts that can't be extended
- [ ] Errors identified clearly
- [ ] Animations can be disabled
- [ ] Language setting works
- [ ] Date/time formatted per locale

---

**Version History**
| Version | Date | Changes |
|---------|------|---------|
| 2.0 | Feb 2026 | Farmer-centric design specification |

**Related Documents**
- TOKKATOT_2.0_FARMER_CENTRIC_ADDITIONS.md
- SPECIFICATIONS_TECHNOLOGY_STACK.md
- SPECIFICATIONS_ARCHITECTURE.md
