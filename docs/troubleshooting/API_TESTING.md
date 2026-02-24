# Testing Tokkatot Backend - Quick Start

## Option 1: Using Docker PostgreSQL (Recommended)

### Step 1: Start Docker Desktop
Make sure Docker Desktop is running on your machine.

### Step 2: Start PostgreSQL Container
```powershell
docker run -d `
  --name tokkatot-postgres `
  -e POSTGRES_PASSWORD=postgres `
  -e POSTGRES_DB=tokkatot `
  -p 5432:5432 `
  postgres:14
```

### Step 3: Verify PostgreSQL is Running
```powershell
docker ps
# Should show tokkatot-postgres container
```

---

## Option 2: Using Local PostgreSQL Installation

If you have PostgreSQL installed locally:
1. Ensure PostgreSQL service is running
2. Create database: `CREATE DATABASE tokkatot;`
3. Update `.env` file with your credentials

---

## Start the Backend Server

```powershell
cd c:\Users\PureGoat\tokkatot\middleware
.\backend.exe
```

Expected output:
```
✅ Configuration loaded - Environment: development
✅ Database connection established
✅ Database schema created/updated
✅ Server starting on 0.0.0.0:3000
```

---

## Test API Endpoints

### Test 1: Health Check (Server Running)
```powershell
curl http://localhost:3000/v1/auth/login
# Should return 400 Bad Request (expected - no credentials provided)
```

### Test 2: User Signup
```powershell
$headers = @{
    "Content-Type" = "application/json"
}

$body = @{
    email = "farmer@tokkatot.com"
    password = "SecureFarmer2026!"
    name = "Neath Farmer"
    language = "km"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:3000/v1/auth/signup" `
  -Method POST `
  -Headers $headers `
  -Body $body
```

Expected response:
```json
{
  "success": true,
  "data": {
    "user_id": "uuid-here",
    "message": "Verification code sent to farmer@tokkatot.com"
  },
  "message": "User registered successfully"
}
```

### Test 3: User Login
```powershell
$body = @{
    email = "farmer@tokkatot.com"
    password = "SecureFarmer2026!"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://localhost:3000/v1/auth/login" `
  -Method POST `
  -Headers $headers `
  -Body $body

# Save token for next tests
$token = $response.data.access_token
Write-Host "Token: $token"
```

Expected response:
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
    "expires_in": 86400,
    "user": {
      "id": "uuid",
      "name": "Neath Farmer",
      "email": "farmer@tokkatot.com",
      "role": "farmer",
      "language": "km"
    }
  },
  "message": "Login successful"
}
```

### Test 4: Token Refresh
```powershell
$body = @{
    refresh_token = $response.data.refresh_token
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:3000/v1/auth/refresh" `
  -Method POST `
  -Headers $headers `
  -Body $body
```

### Test 5: Protected Endpoint (Will add in Phase 2)
```powershell
# Test with Authorization header
$authHeaders = @{
    "Authorization" = "Bearer $token"
    "Content-Type" = "application/json"
}

# This will be available in Phase 2:
# Invoke-RestMethod -Uri "http://localhost:3000/v1/users/me" -Headers $authHeaders
```

---

## Test with Phone Number Instead of Email

```powershell
$body = @{
    phone = "+855123456789"
    password = "SecureFarmer2026!"
    name = "Sophal Farmer"
    language = "km"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:3000/v1/auth/signup" `
  -Method POST `
  -Headers $headers `
  -Body $body
```

---

## Troubleshooting

### Error: "Failed to initialize database"
- Check PostgreSQL is running: `docker ps` or check local PostgreSQL service
- Verify DB credentials in `.env` file
- Test connection: `psql -h localhost -U postgres -d tokkatot`

### Error: "JWT_SECRET environment variable not set"
- Ensure `.env` file exists in `middleware/` directory
- Check `.env` has `JWT_SECRET` value

### Error: "Port 3000 already in use"
- Change `SERVER_PORT=3001` in `.env`
- Or stop conflicting process: `netstat -ano | findstr :3000`

### Error: "Database connection refused"
- Docker: `docker start tokkatot-postgres`
- Local: Start PostgreSQL service

---

## Database Inspection

Connect to PostgreSQL and verify data:
```bash
docker exec -it tokkatot-postgres psql -U postgres -d tokkatot

# List tables
\dt

# View users
SELECT id, email, phone, name, language, created_at FROM users;

# View farms
SELECT id, owner_id, name, created_at FROM farms;

# View farm_users
SELECT farm_id, user_id, role FROM farm_users;

# Exit
\q
```

---

## Next Steps After Successful Tests

Once authentication tests pass:
1. ✅ Phase 1 Complete: Authentication & Database
2. 🔄 Phase 2: User Management, Farm Management, Device Management
3. 🔄 Phase 3: Scheduling, Event Logging, Real-time WebSocket

