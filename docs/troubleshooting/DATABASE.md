# Quick PostgreSQL Database Fix for Tokkatot Backend

## Issue
Backend cannot connect to PostgreSQL due to authentication failure.

## Solution Options

### Option 1: Use Existing eventms Database (Quickest)

Update your PostgreSQL password or verify access:

```powershell
# Test PostgreSQL connection manually
cd "C:\Program Files\PostgreSQL\17\bin"
.\psql.exe -U eventms -d eventms

# If prompted for password, enter it
# If successful, create tokkatot schema:
CREATE SCHEMA IF NOT EXISTS tokkatot;
```

### Option 2: Create New tokkatot Database

```powershell
# Connect as postgres superuser
cd "C:\Program Files\PostgreSQL\17\bin"
.\psql.exe -U postgres

# In psql prompt:
CREATE DATABASE tokkatot;
CREATE USER tokkatot_user WITH ENCRYPTED PASSWORD 'Tokkatot2026!';
GRANT ALL PRIVILEGES ON DATABASE tokkatot TO tokkatot_user;
\q
```

Then update `middleware\.env`:
```env
DB_USER=tokkatot_user
DB_PASSWORD=Tokkatot2026!
DB_NAME=tokkatot
```

### Option 3: Remove System Environment Variables

Your system has DB_USER and DB_PASSWORD set globally. To use .env file instead:

```powershell
# Remove user environment variables
[Environment]::SetEnvironmentVariable("DB_USER", $null, "User")
[Environment]::SetEnvironmentVariable("DB_PASSWORD", $null, "User")

# Restart PowerShell terminal
# Then set in .env file only
```

### Option 4: Override with PowerShell Session Variables

```powershell
# Set for current session only
$env:DB_USER = "postgres"
$env:DB_PASSWORD = "your_postgres_password"
$env:DB_NAME = "tokkatot"

# Then run backend
cd c:\Users\PureGoat\tokkatot\middleware
.\backend.exe
```

## Steps After Database is Accessible

1. **Start Backend**:
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

2. **Test in Another Terminal**:
   ```powershell
   # Test 1: Server responding
   curl http://localhost:3000/v1/auth/login

   # Should get 400 Bad Request (expected - no credentials)
   ```

3. **Create Test User**:
   ```powershell
   $body = @{
       email = "test@tokkatot.com"
       password = "TestFarmer123!"
       name = "Test Farmer"
       language = "km"
   } | ConvertTo-Json

   Invoke-RestMethod -Uri "http://localhost:3000/v1/auth/signup" `
       -Method POST `
       -ContentType "application/json" `
       -Body $body
   ```

4. **Login**:
   ```powershell
   $loginBody = @{
       email = "test@tokkatot.com"
       password = "TestFarmer123!"
   } | ConvertTo-Json

   $response = Invoke-RestMethod -Uri "http://localhost:3000/v1/auth/login" `
       -Method POST `
       -ContentType "application/json" `
       -Body $loginBody

   # View response
   $response | ConvertTo-Json -Depth 10

   # Save token
   $token = $response.data.access_token
   Write-Host "Access Token: $token"
   ```

5. **Verify Database**:
   ```powershell
   # Connect to PostgreSQL
   cd "C:\Program Files\PostgreSQL\17\bin"
   .\psql.exe -U eventms -d tokkatot  # or your user/db

   # Check tables
   \dt

   # View users
   SELECT id, name, email, created_at FROM users;

   # View farms
   SELECT id, name, owner_id, created_at FROM farms;
   ```

## Quick Diagnostic

Run this PowerShell script to check configuration:

```powershell
Write-Host "=== Tokkatot Backend Diagnostics ===" -ForegroundColor Cyan

# Check PostgreSQL service
Write-Host "`n1. PostgreSQL Services:" -ForegroundColor Yellow
Get-Service | Where-Object {$_.Name -like "*postgres*"} | Format-Table Name, Status

# Check environment variables
Write-Host "`n2. Environment Variables:" -ForegroundColor Yellow
"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME" | ForEach-Object {
    $val = [Environment]::GetEnvironmentVariable($_, "User")
    if ($val) { Write-Host "  $_=$val" }
}

# Check .env file
Write-Host "`n3. .env File:" -ForegroundColor Yellow
if (Test-Path "c:\Users\PureGoat\tokkatot\middleware\.env") {
    Get-Content "c:\Users\PureGoat\tokkatot\middleware\.env" | Where-Object { $_ -match "^DB_" }
} else {
    Write-Host "  .env file not found!" -ForegroundColor Red
}

# Check backend binary
Write-Host "`n4. Backend Binary:" -ForegroundColor Yellow
if (Test-Path "c:\Users\PureGoat\tokkatot\middleware\backend.exe") {
    $file = Get-Item "c:\Users\PureGoat\tokkatot\middleware\backend.exe"
    Write-Host "  ✅ Exists ($([math]::Round($file.Length/1MB, 2)) MB)"
} else {
    Write-Host "  ❌ Not found" -ForegroundColor Red
}

# Test PostgreSQL port
Write-Host "`n5. PostgreSQL Port 5432:" -ForegroundColor Yellow
try {
    $conn = Test-NetConnection -ComputerName localhost -Port 5432 -WarningAction SilentlyContinue
    if ($conn.TcpTestSucceeded) {
        Write-Host "  ✅ Port 5432 is open"
    } else {
        Write-Host "  ❌ Port 5432 is not accessible" -ForegroundColor Red
    }
} catch {
    Write-Host "  ⚠️  Cannot test port" -ForegroundColor Yellow
}

Write-Host "`n=====================================" -ForegroundColor Cyan
```

## Need Help?

If still having issues, provide:
1. PostgreSQL version: `psql --version`
2. Can you connect manually? `psql -U postgres -d postgres`
3. Error logs from backend

