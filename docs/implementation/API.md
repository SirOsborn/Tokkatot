# API Implementation Notes

Source of truth for routing:
- `middleware/main.go`

Key conventions:
- Base path: `/v1`
- JWT auth required for protected routes
- Registration key is the only verification path (no email/SMS verification)
- AI prediction endpoints are not exposed yet (coming in a future patch)
- Automation is **coop-level** (farm is container)

Core endpoints to keep in sync:
- Auth: `/v1/auth/signup`, `/v1/auth/login`, `/v1/auth/refresh`, `/v1/auth/logout`
- Users: `/v1/users/me`, `/v1/users/sessions`
- Farms: `/v1/farms`, `/v1/farms/:id/members`
- Coops: `/v1/farms/:farm_id/coops`, `/v1/farms/:farm_id/coops/:coop_id`
- Devices: `/v1/farms/:id/devices`, `/v1/farms/:id/devices/:id/commands`
- Schedules: `/v1/farms/:id/schedules`
- Telemetry: `/v1/farms/:farm_id/coops/:coop_id/telemetry`
- Device Report: `/v1/farms/:farm_id/coops/:coop_id/devices/report`
- Monitoring Timeline: `/v1/farms/:farm_id/coops/:coop_id/temperature-timeline`

If you change an endpoint, update:
- `middleware/main.go`
- Swagger annotations in `middleware/api/*.go`
- This file (summary or link)
