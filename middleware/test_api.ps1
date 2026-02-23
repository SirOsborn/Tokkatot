# Test Tokkatot API Endpoints
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Tokkatot API Testing" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

$baseUrl = "http://localhost:3000"
$testEmail = "farmer$(Get-Random -Minimum 100 -Maximum 999)@tokkatot.com"

# OPTION 1: Signup with Registration Key (On-Site Setup)
Write-Host "`n[TEST 1A] POST /v1/auth/signup (With Registration Key)" -ForegroundColor Yellow
Write-Host "Simulating on-site account creation by staff..." -ForegroundColor Gray

$signupWithKeyBody = @{
    email = $testEmail
    password = "Farmer123"
    name = "Sokha Farmer"
    registration_key = "ZKJFA-ZIVMC-HOUGG-XQSRW-ITDYH"  # From .env REG_KEY
} | ConvertTo-Json

try {
    $signupResponse = Invoke-RestMethod -Uri "$baseUrl/v1/auth/signup" -Method POST -ContentType "application/json" -Body $signupWithKeyBody
    Write-Host "✅ Signup with registration key successful!" -ForegroundColor Green
    Write-Host "   Auto-verified: YES (no email/SMS needed)" -ForegroundColor Green
    Write-Host "Email: $testEmail" -ForegroundColor Gray
    $userId = $signupResponse.data.user_id
    Write-Host "User ID: $userId" -ForegroundColor Gray
    $hasRegKey = $true
} catch {
    Write-Host "⚠️  Registration key signup failed (key may not exist in DB)" -ForegroundColor Yellow
    if ($_.ErrorDetails) {
        $_.ErrorDetails.Message
    }
    $hasRegKey = $false
}

# OPTION 2: Normal Signup (Without Registration Key)
if (-not $hasRegKey) {
    Write-Host "`n[TEST 1B] POST /v1/auth/signup (Without Registration Key)" -ForegroundColor Yellow
    $signupBody = @{
        email = $testEmail
        password = "Farmer123"
        name = "Sokha Farmer"
    } | ConvertTo-Json

    try {
        $signupResponse = Invoke-RestMethod -Uri "$baseUrl/v1/auth/signup" -Method POST -ContentType "application/json" -Body $signupBody
        Write-Host "✅ Normal signup successful!" -ForegroundColor Green
        Write-Host "   Auto-verified: $(if($signupResponse.data.verified){'YES (dev mode)'}else{'NO (needs verification)'})" -ForegroundColor Yellow
        $userId = $signupResponse.data.user_id
    } catch {
        Write-Host "❌ Signup failed:" -ForegroundColor Red
        if ($_.ErrorDetails) {
            $_.ErrorDetails.Message
        }
    }
}

# Test 2: Login (should work if verified)
Write-Host "`n[TEST 2] POST /v1/auth/login" -ForegroundColor Yellow
$loginBody = @{
    email = $testEmail
    password = "Farmer123"
} | ConvertTo-Json

try {
    $loginResponse = Invoke-RestMethod -Uri "$baseUrl/v1/auth/login" -Method POST -ContentType "application/json" -Body $loginBody
    Write-Host "✅ Login successful!" -ForegroundColor Green
    $accessToken = $loginResponse.data.access_token
    $refreshToken = $loginResponse.data.refresh_token
    Write-Host "Access Token: $($accessToken.Substring(0,50))..." -ForegroundColor Gray
    Write-Host "User ID: $($loginResponse.data.user.id)" -ForegroundColor Gray
    Write-Host "Farm ID: $($loginResponse.data.farm_id)" -ForegroundColor Gray
    Write-Host "Role: $($loginResponse.data.role)" -ForegroundColor Gray
} catch {
    Write-Host "❌ Login failed:" -ForegroundColor Red
    $_.Exception.Message
    if ($_.ErrorDetails) {
        $_.ErrorDetails.Message
    }
}

# Test 3: Refresh Token
if ($refreshToken) {
    Write-Host "`n[TEST 3] POST /v1/auth/refresh" -ForegroundColor Yellow
    $refreshBody = @{
        refresh_token = $refreshToken
    } | ConvertTo-Json

    try {
        $newTokenResponse = Invoke-RestMethod -Uri "$baseUrl/v1/auth/refresh" -Method POST -ContentType "application/json" -Body $refreshBody
        Write-Host "✅ Token refresh successful!" -ForegroundColor Green
        Write-Host "New Access Token: $($newTokenResponse.data.access_token.Substring(0,50))..." -ForegroundColor Gray
    } catch {
        Write-Host "❌ Token refresh failed:" -ForegroundColor Red
        if ($_.ErrorDetails) {
            $_.ErrorDetails.Message
        }
    }
}

# Test 4: Using Protected Endpoint (Example - will fail until implemented)
if ($accessToken) {
    Write-Host "`n[TEST 4] GET /v1/users/me (Protected)" -ForegroundColor Yellow
    try {
        $headers = @{
            Authorization = "Bearer $accessToken"
        }
        $meResponse = Invoke-RestMethod -Uri "$baseUrl/v1/users/me" -Method GET -Headers $headers
        Write-Host "✅ Protected endpoint works!" -ForegroundColor Green
        $meResponse | ConvertTo-Json -Depth 3
    } catch {
        Write-Host "⚠️  Protected endpoint not implemented yet" -ForegroundColor Yellow
    }
}

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "  Tests Complete" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
