# API Implementation Notes

Source of truth for routing:
- `middleware/main.go`

Key conventions:
- Base path: `/v1`
- JWT auth required for protected routes
- Registration key is the only verification path (no email/SMS verification)
- AI prediction endpoints are not exposed yet (coming in a future patch)

If you change an endpoint, update:
- `middleware/main.go`
- Swagger annotations in `middleware/api/*.go`
- This file (summary or link)
