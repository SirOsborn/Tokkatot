# Setup Guide - Tokkatot Development Environment

**Last Updated**: February 23, 2026  
**Time Required**: 30 minutes  
**For**: New developers, fresh installations

---

## Prerequisites

### Windows 11 (Recommended)
```powershell
# Install all tools via winget
winget install GoLang.Go.1.23
winget install PostgreSQL.PostgreSQL.17
winget install Git.Git
winget install Microsoft.VisualStudioCode
winget install Docker.DockerDesktop  # Optional - for AI service
winget install OpenJS.NodeJS  # Optional - for Vue.js build tools
```

### Verify Installations
```powershell
go version        # Should show: go1.23.x
psql --version    # Should show: psql (PostgreSQL) 17.x
git --version     # Should show: git version 2.x
code --version    # Should show VS Code version
```

---

## 1. Clone Repository

```powershell
cd C:\Users\YourName
git clone https://github.com/your-org/tokkatot.git
cd tokkatot
```

**Folder Structure:**
```
tokkatot/
├── middleware/       ← Go backend
├── frontend/         ← Vue.js (to migrate)
├── ai-service/       ← Python AI
├── embedded/         ← ESP32 firmware
└── docs/             ← This file!
```

---

## 2. Database Setup (PostgreSQL)

### Quick Setup
```powershell
# 1. PostgreSQL should be running (check Services)
Get-Service postgresql*

# 2. Connect as postgres user
psql -U postgres

# 3. Create database and user (in psql shell)
CREATE DATABASE tokkatot;
CREATE USER tokkatot_user WITH ENCRYPTED PASSWORD 'your_secure_password';
GRANT ALL PRIVILEGES ON DATABASE tokkatot TO tokkatot_user;
\q

# 4. Test connection
psql -U tokkatot_user -d tokkatot -h localhost
```

### Alternative: Use Default postgres User
```powershell
# Simpler for development
psql -U postgres
CREATE DATABASE tokkatot;
\q
```

---

## 3. Backend Setup (Go)

### Configure Environment
```powershell
cd middleware

# Copy example .env file
cp .env.example .env

# Edit .env file (use VS Code or notepad)
code .env
```

**Update `.env` values:**
```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres              # Or tokkatot_user
DB_PASSWORD=YOUR_PASSWORD     # ← CHANGE THIS!
DB_NAME=tokkatot
DB_SSLMODE=disable

# Server
SERVER_PORT=3000
SERVER_HOST=0.0.0.0

# JWT Secret (generate random 32+ chars)
JWT_SECRET=YOUR_RANDOM_SECRET_HERE  # ← CHANGE THIS!
# Example generation: openssl rand -hex 32

# Registration Key (for testing)
REG_KEY=XXXXX-XXXXX-XXXXX-XXXXX-XXXXX

# Environment
ENVIRONMENT=development
```

### Build & Run
```powershell
# Install dependencies
go mod download

# Build binary
go build -o backend.exe

# Run backend
.\backend.exe
```

**Expected Output:**
```
✅ Configuration loaded - Environment: development
✅ Database connection established
✅ Database schema created/updated
✅ Server starting on 0.0.0.0:3000
```

### Test API
```powershell
# In new terminal
cd middleware
.\test_api.ps1
```

**Successful test shows:**
```
✅ Health check passed
✅ Signup successful
✅ Login successful
✅ Access token valid
```

---

## 4. Frontend Setup (Vue.js Migration)

**Current State:**  
Frontend uses vanilla HTML/CSS/JS (no build step)

**Migration Plan:**  
Phase 1: Add Vue.js 3 via CDN (no build)  
Phase 2: Components (later)  
Phase 3: Vite build (optional)

### Quick Test (Serve Frontend)
```powershell
# Backend already serves frontend files

# Open browser
start http://localhost:3000/login
```

**Files served:**
- `frontend/pages/login.html` → `/login`
- `frontend/pages/index.html` → `/` (dashboard)
- `frontend/css/`, `frontend/js/`, `frontend/assets/` → Static assets

---

## 5. AI Service Setup (Optional - Python)

### Prerequisites
```powershell
# Install Python
winget install Python.Python.3.12

# Verify
python --version  # Should show: Python 3.12.x
```

### Setup AI Service
```powershell
cd ai-service

# Create virtual environment
python -m venv env

# Activate (PowerShell)
.\env\Scripts\Activate

# Install dependencies
pip install --upgrade pip
pip install -r requirements.txt

# Run service
python app.py
```

**Expected Output:**
```
INFO:     Uvicorn running on http://0.0.0.0:8000
INFO:     Application startup complete
```

**Note**: AI service requires `outputs/ensemble_model.pth` (47MB, not in repo)

### Alternative: Docker
```powershell
cd ai-service

# Build image
docker build -t tokkatot-ai:latest .

# Run container
docker-compose up -d

# Check logs
docker-compose logs -f
```

---

## 6. Generate Registration Key

```powershell
cd middleware

# Generate registration key for testing
.\generate_reg_key.ps1 -FarmName "Test Farm" -CustomerName "Developer" -CustomerPhone "012345678"
```

**Output:**
```
Registration Key: ABCDE-FGHIJ-KLMNO-PQRST-UVWXY
SQL: INSERT INTO registration_keys (id, key_code, farm_name, ...) VALUES (...);
```

**Add to `.env`:**
```bash
REG_KEY=ABCDE-FGHIJ-KLMNO-PQRST-UVWXY
```

---

## Troubleshooting

### Database Connection Failed
```
Error: pq: password authentication failed for user "postgres"
```

**Fix:**
```powershell
# Option 1: Reset postgres password
psql -U postgres
ALTER USER postgres WITH PASSWORD 'new_password';
\q

# Update .env file
code middleware\.env  # Change DB_PASSWORD

# Option 2: Use peer authentication (localhost only)
# Edit: C:\Program Files\PostgreSQL\17\data\pg_hba.conf
# Change: host all all 127.0.0.1/32 md5
# To:     host all all 127.0.0.1/32 trust
# Restart PostgreSQL service
```

### Port 3000 Already in Use
```
Error: bind: address already in use
```

**Fix:**
```powershell
# Find process using port 3000
Get-Process -Id (Get-NetTCPConnection -LocalPort 3000).OwningProcess

# Stop backend if running
Get-Process backend -ErrorAction SilentlyContinue | Stop-Process -Force

# Or change port in .env
SERVER_PORT=3001
```

### Go Build Fails
```
Error: package middleware/xxx is not in GOROOT
```

**Fix:**
```powershell
# Clean and rebuild
go clean -modcache
go mod tidy
go mod download
go build -o backend.exe
```

### PostgreSQL Service Not Running
```powershell
# Check service status
Get-Service postgresql*

# Start service
Start-Service postgresql-x64-17  # Or your service name

# Set to auto-start
Set-Service postgresql-x64-17 -StartupType Automatic
```

---

## Development Workflow

### Daily Workflow
```powershell
# 1. Start PostgreSQL (if not running)
Get-Service postgresql* | Start-Service

# 2. Open VS Code
code C:\Users\YourName\tokkatot

# 3. Start backend (terminal 1)
cd middleware
.\backend.exe

# 4. Start frontend dev (terminal 2)
# Just open browser: http://localhost:3000

# 5. Make changes, test, commit
git add .
git commit -m "Add feature"
git push
```

### Code Formatting
```powershell
# Go code
cd middleware
go fmt ./...

# Frontend (if using Prettier later)
cd frontend
npx prettier --write "**/*.{js,html,css}"
```

---

## Next Steps

**Backend Development:**
1. Read [../ARCHITECTURE.md](../ARCHITECTURE.md) for system design
2. Read [../implementation/API.md](../implementation/API.md) for API development
3. Read [../implementation/DATABASE.md](../implementation/DATABASE.md) for schema details

**Frontend Development:**
1. Read [../implementation/FRONTEND.md](../implementation/FRONTEND.md) for Vue.js migration
2. Read [../TECH_STACK.md](../TECH_STACK.md) for frontend tech choices

**Embedded Development:**
1. Read [../implementation/EMBEDDED.md](../implementation/EMBEDDED.md) for ESP32 setup
2. Install ESP-IDF toolchain

**Testing:**
1. Read [TEST_BACKEND.md](../TEST_BACKEND.md) for API testing
2. Use `middleware/test_api.ps1` for automated tests

---

## Quick Reference

| Task | Command |
|------|---------|
| Build backend | `cd middleware && go build -o backend.exe` |
| Run backend | `cd middleware && .\backend.exe` |
| Test API | `cd middleware && .\test_api.ps1` |
| Generate reg key | `cd middleware && .\generate_reg_key.ps1 -FarmName "Test"` |
| DB console | `psql -U postgres -d tokkatot` |
| View logs | Backend prints to console (redirect: `.\backend.exe > log.txt`) |
| Stop backend | `Ctrl+C` or `Get-Process backend \| Stop-Process` |

---

**Setup Complete! 🎉**  
Backend should be running on http://localhost:3000  
API docs: Open `docs/implementation/API.md`
