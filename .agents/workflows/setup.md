---
description: How to set up the Tokkatot development environment
---

1.  **Database (PostgreSQL)**:
    *   Install PostgreSQL 17+.
    *   Create a user and database.
    *   Set `DATABASE_URL` in `middleware/.env`. (Format: `postgres://user:pass@host:port/db?sslmode=disable`).
2.  **Backend (Go)**:
    *   Install Go 1.23+.
    *   `cd middleware && go mod download`.
    *   Run `go run main.go`.
3.  **AI Service (Python)**:
    *   Install Python 3.12+.
    *   `cd ai-service && pip install -r requirements.txt`.
    *   Run `python app.py`. (Requires `outputs/*.pth` local ensemble weights).
4.  **Frontend**:
    *   The frontend is served directly by Go at `http://localhost:3000`.
    *   No build step is required (CDN-based).
