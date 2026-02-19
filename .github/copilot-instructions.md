# Copilot Instructions for Tokkatot Repository

## Repository Overview

**Project**: Tokkatot - Smart Poultry Farm Management System  
**Purpose**: IoT-based disease detection, monitoring, and automation for Cambodian poultry farmers  
**Tech Stack**: Go 1.23+ (middleware), Python 3.8+ with PyTorch (AI service), Vue.js 3 (frontend PWA), C/ESP-IDF (embedded), PostgreSQL/InfluxDB (data)  
**Repository Type**: Full-stack web + embedded IoT system with AI inference  
**Key Requirement**: Farmer-centric design (simple 3-role system, mobile-first, Khmer language support)

## Project Structure & Layout

```
tokkatot/
├── middleware/              # Go REST API gateway (66 endpoints, JWT auth)
│   ├── main.go             # Server entry point, loads .env
│   ├── api/                # Endpoint handlers (authentication, devices, AI, etc.)
│   ├── database/           # SQLite wrapper (production uses PostgreSQL)
│   ├── utils/              # Helper functions
│   ├── go.mod              # Go 1.23, Fiber v2.52.6, JWT v4.5.0
│   └── .env                # SECRETS: JWT_SECRET, REG_KEY, TLS_CERT, TLS_KEY (NEVER PUSH)
├── ai-service/             # Python FastAPI ensemble disease detection
│   ├── app.py              # FastAPI server with 3 endpoints (/predict, /health)
│   ├── inference.py        # ChickenDiseaseDetector class (PyTorch ensemble)
│   ├── models.py           # EfficientNetB0 + DenseNet121 architectures
│   ├── data_utils.py       # Image preprocessing, class definitions
│   ├── requirements.txt     # PyTorch 2.0+, FastAPI, Uvicorn, Pydantic
│   ├── pyproject.toml       # Project metadata and dependencies
│   ├── Dockerfile          # Python 3.12-slim, WORKDIR /app, health checks
│   ├── docker-compose.yml  # Port 8000, 2 CPU / 4GB RAM limits, health checks
│   ├── outputs/            # PROPRIETARY: *.pth model files (NEVER PUSH - .gitignore'd)
│   └── .env                # NEVER PUSH - .gitignore'd
├── frontend/               # Vue.js 3 PWA (no build step needed - static files)
│   ├── pages/              # HTML pages (index, login, disease-detection, profile)
│   ├── js/                 # Vue components and handlers
│   ├── css/                # Responsive styles (48px+ fonts for accessibility)
│   ├── components/         # Reusable HTML components (navbar, header)
│   └── assets/             # Images, fonts, icons
├── embedded/               # ESP32 C firmware (IoT device)
│   ├── CMakeLists.txt      # ESP-IDF build config
│   ├── main/
│   │   ├── main.c          # Device boot, MQTT client setup
│   │   └── CMakeLists.txt
│   ├── components/dht/     # DHT22 sensor driver
│   ├── sdkconfig           # ESP-IDF configuration (UART, GPIO, WiFi)
│   └── build/              # Compiled output (do not commit)
├── docs/                   # 14+ specification files (see 00_SPECIFICATIONS_INDEX.md)
├── .github/
│   ├── copilot-instructions.md  # This file
│   ├── CODE_REVIEW_CHECKLIST.md # PR review checklist for Copilot
│   └── workflows/main.yml        # CI/CD pipeline (see Build section)
├── AI_INSTRUCTIONS.md      # Master AI agent guide for all components
├── generate-cert.sh        # Creates self-signed TLS certificates in certs/
├── .gitignore              # Excludes secrets, models, build artifacts
└── LICENSE & README.md
```

## Build & Validation Commands

### Prerequisites & Environment Setup
```bash
# All commands run from workspace root: c:\Users\PureGoat\tokkatot

# Go middleware
go version  # Required: 1.23.6+ (see middleware/go.mod)
cd middleware && go mod download && cd ..

# Python AI service
python3 --version  # Required: 3.8+
cd ai-service && python3 -m venv env && source env/bin/activate && cd ..

# Docker (optional but recommended)
docker --version && docker-compose --version
```

### Middleware Build (Go)
```bash
cd middleware

# Load environment variables (REQUIRED - see main.go line 16-27)
# Create .env file with: JWT_SECRET, REG_KEY, TLS_CERT, TLS_KEY
echo "JWT_SECRET=dev-secret-key-min-32-chars-long" >> .env

# Download dependencies
go mod download

# Build binary
go build -o middleware  # Output: middleware/middleware binary

# Verify build succeeded
ls -la middleware && file middleware

# Clean (if needed)
go clean
cd ..
```

### AI Service Build (Python)
```bash
cd ai-service

# Create Python virtual environment (REQUIRED EVERY TIME - do not reuse across versions)
python3 -m venv env
source env/bin/activate  # On Windows: env\Scripts\activate

# Install dependencies (ALWAYS run this after venv creation)
pip install --upgrade pip setuptools wheel
pip install -r requirements.txt

# Verify installation
python3 -c "import torch; print(f'PyTorch {torch.__version__}')"
python3 -c "import fastapi; print('FastAPI OK')"

# Test app startup (requires outputs/ensemble_model.pth - see NOTE below)
python3 app.py  # Starts on localhost:8000

# Build Docker image (RECOMMENDED over local Python)
docker build -t tokkatot-ai:latest .
docker-compose up -d tokkatot-ai  # Runs on port 8000

# Clean (if needed)
rm -rf env && docker-compose down && docker rmi tokkatot-ai:latest
cd ..
```

**IMPORTANT NOTE**: AI Service model loading:
- `app.py` startup loads ensemble model from `outputs/ensemble_model.pth` (47.2 MB)
- This file is **NOT** in the repository (proprietary, .gitignore'd)
- Model file MUST exist locally for `python app.py` to work
- Docker image COPY includes `outputs/` (implies model present during build)
- If model file missing: "FileNotFoundError: could not find ensemble_model.pth"

### Embedded Build (ESP32)
```bash
# Requires ESP-IDF toolchain installed (see docs/IG_SPECIFICATIONS_EMBEDDED.md)
cd embedded

# Set ESP-IDF environment
export IDF_PATH=/path/to/esp-idf  # e.g., ~/esp/esp-idf

# Build firmware
idf.py -p /dev/ttyUSB0 build  # /dev/ttyUSB0 or COM3 on Windows

# Flash device
idf.py -p /dev/ttyUSB0 flash

# Monitor serial output
idf.py -p /dev/ttyUSB0 monitor

# Clean
idf.py clean
cd ..
```

### Frontend (No Build Required)
- Vue.js 3 with static HTML files
- No npm/build step needed - served directly by middleware
- Files in `frontend/pages/` are served by Go middleware
- Validation: Open browser to `http://localhost:3000` after middleware starts

### CI/CD Pipeline (GitHub Actions)
**File**: `.github/workflows/main.yml`  
**Trigger**: Push to main branch  
**Runner**: self-hosted (specified in workflow)  
**Environment**: Loads secrets from GitHub Actions: JWT_SECRET, REG_KEY, TLS_CERT, TLS_KEY

**Pipeline Steps**:
1. Checkout code
2. Setup Go 1.23.6
3. Generate TLS certificates: `./generate-cert.sh`
4. Create middleware/.env with secrets
5. Build Go binary: `cd middleware && go build`
6. Set Linux capability for port 443: `setcap 'cap_net_bind_service=+ep' ./middleware`
7. Setup Python venv and install AI requirements
8. Restart systemd services: `tokkatot.service`, `tokkatot-ai.service`

**Common CI/CD Failures**:
- ❌ "JWT_SECRET not found" → Secrets not configured in GitHub Actions (Settings → Secrets)
- ❌ "go build: module not found" → Run `go mod download` before build
- ❌ "setcap failed" → Only works on Linux/self-hosted runners, not GitHub cloud
- ❌ "systemctl: permission denied" → Workflow uses `sudo`, self-hosted runner must have passwordless sudo

## Key Implementation Notes

### Critical Files (Do NOT Modify Without Understanding):
- `middleware/main.go` - Loads .env, initializes Fiber server, sets up routes
- `ai-service/app.py` - FastAPI startup event loads ensemble model from disk
- `ai-service/inference.py` - ChickenDiseaseDetector class with ensemble voting logic
- `docs/00_SPECIFICATIONS_INDEX.md` - Read first to understand system design

### Database Connection:
- Middleware uses SQLite in development (`database/sqlite3_db.go`)
- Production target: PostgreSQL (schema in `docs/IG_SPECIFICATIONS_DATABASE.md`)
- No automatic migrations - schema defined in specification, not in code

### Authentication:
- JWT tokens issued by `middleware/api/authentication.go`
- Token claims: user_id, email, role (Owner/Manager/Viewer)
- No MFA for farmers (security tradeoff for accessibility)
- Token validation required on all protected endpoints

### API Integration:
- Middleware and AI service are separate deployments
- Middleware calls AI service at `http://localhost:8000/predict` (see `middleware/api/disease-detection.go`)
- Request format: image file upload, returns `{disease, confidence, recommendations}`
- Timeout: 3 seconds (CPU) or <500ms (GPU) for disease prediction

### File Exclusions (NEVER PUSH):
- `ai-service/outputs/` - Proprietary model files (*.pth)
- `middleware/.env` - JWT_SECRET, credentials
- `ai-service/.env` - Configuration secrets
- `.vscode/settings.json` - Local user preferences
- `*/build/` directories - Build artifacts

**Trust this guide** - use grep/find/search tools ONLY if instructions are incomplete or found to be in error.
