# Tokkatot Backend API Test Script
# Runs all endpoint tests sequentially with a single command.
# Usage:
#   .\test_all_endpoints.ps1             # login with email (default)
#   .\test_all_endpoints.ps1 -UsePhone   # login with phone number
#
# Prerequisites:
#   1. Backend running:  cd middleware && .\backend.exe
#   2. Seeded test data: farmer@tokkatot.com, farm 11111111-..., device 33333333-...
#      (Run the seed SQL from AI_INSTRUCTIONS.md or ask Copilot to seed the DB)

param([switch]$UsePhone)

# ── Config ───────────────────────────────────────────────────────────────────
$BASE            = "http://localhost:3000/v1"
$FARMER_EMAIL    = "farmer@tokkatot.com"
$FARMER_PHONE    = "+85512345678"
$FARMER_PASSWORD = "FarmerPass123"
$FARM_ID         = "11111111-1111-1111-1111-111111111111"
$DEVICE_ID       = "33333333-3333-3333-3333-333333333333"
$SEQ_SCHED_ID    = "44444444-4444-4444-4444-444444444444"   # seeded: has action_sequence
$DUR_SCHED_ID    = "55555555-5555-5555-5555-555555555555"   # seeded: has action_duration

# ── Counters & state ─────────────────────────────────────────────────────────
$script:passed   = 0
$script:failed   = 0
$script:skipped  = 0
$script:token    = $null
$script:headers  = @{}
$script:newCoopId     = $null
$script:newScheduleId = $null

# ── Output helpers ───────────────────────────────────────────────────────────
function Write-Section {
    param($title)
    Write-Host "`n══════════════════════════════════════════════════" -ForegroundColor DarkGray
    Write-Host "  $title" -ForegroundColor Yellow
    Write-Host "══════════════════════════════════════════════════" -ForegroundColor DarkGray
}
function Write-Pass   { param($msg) Write-Host "  ✅ PASS  $msg" -ForegroundColor Green }
function Write-Fail   { param($msg) Write-Host "  ❌ FAIL  $msg" -ForegroundColor Red }
function Write-Skip   { param($msg) Write-Host "  ⏭️  SKIP  $msg" -ForegroundColor DarkGray }
function Write-Detail { param($msg) Write-Host "         $msg" -ForegroundColor Gray }

# ── Core test runner ─────────────────────────────────────────────────────────
# Returns the parsed response object on success, $null on failure.
function Invoke-Test {
    param(
        [string]$Method,
        [string]$Url,
        [string]$Description,
        [hashtable]$Headers = @{},
        [string]$Body = $null,
        [scriptblock]$Validate = $null     # optional extra assertions
    )

    try {
        $params = @{
            Uri         = $Url
            Method      = $Method
            Headers     = $Headers
            ErrorAction = "Stop"
        }
        if ($Body) {
            $params.Body        = $Body
            $params.ContentType = "application/json"
        }

        $resp = Invoke-RestMethod @params

        if ($Validate) {
            $validationMsg = & $Validate $resp
            if ($validationMsg) {
                Write-Fail "$Description — assertion failed: $validationMsg"
                $script:failed++
                return $null
            }
        }

        Write-Pass $Description
        $script:passed++
        return $resp
    }
    catch {
        $errDetail = ""
        if ($_.ErrorDetails.Message) {
            try { $errDetail = ($_.ErrorDetails.Message | ConvertFrom-Json).error.message } catch {}
        }
        if (-not $errDetail) { $errDetail = $_.Exception.Message }
        Write-Fail "$Description — $errDetail"
        $script:failed++
        return $null
    }
}

# ── Convenience wrappers ──────────────────────────────────────────────────────
function AuthTest  { param($m,$ep,$desc,$body,$v) Invoke-Test -Method $m -Url "$BASE$ep" -Description $desc -Body $body -Validate $v }
function PubTest   { param($m,$ep,$desc,$body,$v) Invoke-Test -Method $m -Url "$BASE$ep" -Description $desc -Body $body -Validate $v }
function PrivTest  { param($m,$ep,$desc,$body,$v) Invoke-Test -Method $m -Url "$BASE$ep" -Description $desc -Headers $script:headers -Body $body -Validate $v }

# ─────────────────────────────────────────────────────────────────────────────
# 0. CONNECTIVITY CHECK
# ─────────────────────────────────────────────────────────────────────────────
Write-Host "`n🚀 Tokkatot API Test Suite" -ForegroundColor Magenta
Write-Host "   Target: $BASE" -ForegroundColor Gray
Write-Host "   Auth:   $(if ($UsePhone) { 'phone' } else { 'email' })`n" -ForegroundColor Gray

try {
    $null = Invoke-RestMethod -Uri "http://localhost:3000/v1/auth/login" -Method POST `
        -ContentType "application/json" -Body '{"email":"x","password":"x"}' -ErrorAction SilentlyContinue
} catch {
    if ($_.Exception.Message -match "refused|connect") {
        Write-Host "❌ Backend not reachable at localhost:3000" -ForegroundColor Red
        Write-Host "   Run:  cd middleware && .\backend.exe" -ForegroundColor Yellow
        exit 1
    }
}

# ─────────────────────────────────────────────────────────────────────────────
# 1. AUTHENTICATION
# ─────────────────────────────────────────────────────────────────────────────
Write-Section "1. AUTHENTICATION"

# 1-1  Login
$loginBody = if ($UsePhone) {
    "{`"phone`":`"$FARMER_PHONE`",`"password`":`"$FARMER_PASSWORD`"}"
} else {
    "{`"email`":`"$FARMER_EMAIL`",`"password`":`"$FARMER_PASSWORD`"}"
}

$loginResp = Invoke-Test -Method POST -Url "$BASE/auth/login" -Description "Login (POST /auth/login)" `
    -Body $loginBody `
    -Validate {
        param($r)
        if (-not $r.data.access_token) { return "missing data.access_token" }
    }

if (-not $loginResp) {
    Write-Host "`n❌ Cannot continue without a valid token. Aborting." -ForegroundColor Red
    exit 1
}

$script:token   = $loginResp.data.access_token.Trim()
$script:headers = @{ "Authorization" = "Bearer $script:token"; "Content-Type" = "application/json" }
Write-Detail "Token: $($script:token.Substring(0,40))..."

# 1-2  Refresh token
$refreshBody = "{`"refresh_token`":`"$($loginResp.data.refresh_token)`"}"
$refreshResp = Invoke-Test -Method POST -Url "$BASE/auth/refresh" -Description "Refresh token (POST /auth/refresh)" `
    -Body $refreshBody `
    -Validate { param($r) if (-not $r.data.access_token) { return "missing new access_token" } }

# Update token to refreshed one if it worked
if ($refreshResp) {
    $script:token   = $refreshResp.data.access_token.Trim()
    $script:headers = @{ "Authorization" = "Bearer $script:token"; "Content-Type" = "application/json" }
    Write-Detail "Refreshed token accepted"
}

# 1-3  Forgot password (public endpoint, just checks HTTP 200)
Invoke-Test -Method POST -Url "$BASE/auth/forgot-password" `
    -Description "Forgot password (POST /auth/forgot-password)" `
    -Body "{`"email`":`"$FARMER_EMAIL`"}" | Out-Null

# 1-4  Verify contact — skipped (requires external OTP)
Write-Skip "Verify contact (POST /auth/verify) — requires OTP from email/SMS"
$script:skipped++

# 1-5  Reset password — skipped (requires token from email)
Write-Skip "Reset password (POST /auth/reset-password) — requires email token"
$script:skipped++

# ─────────────────────────────────────────────────────────────────────────────
# 2. USER PROFILE
# ─────────────────────────────────────────────────────────────────────────────
Write-Section "2. USER PROFILE"

# 2-1  Get current user
$meResp = PrivTest GET "/users/me" "Get current user (GET /users/me)" -v {
    param($r) if (-not $r.data.id) { return "missing data.id" }
}
if ($meResp) { Write-Detail "User: $($meResp.data.name)  ($($meResp.data.email))" }

# 2-2  Update profile
PrivTest PUT "/users/me" "Update profile (PUT /users/me)" `
    -body "{`"name`":`"Sokha Farmer`",`"language`":`"km`"}" | Out-Null

# 2-3  Revert profile back
PrivTest PUT "/users/me" "Revert profile language (PUT /users/me)" `
    -body "{`"language`":`"en`"}" | Out-Null

# 2-4  Change password — skip (avoid breaking test-user credentials)
Write-Skip "Change password (POST /users/me/change-password) — skipped to preserve test credentials"
$script:skipped++

# ─────────────────────────────────────────────────────────────────────────────
# 3. FARM MANAGEMENT
# ─────────────────────────────────────────────────────────────────────────────
Write-Section "3. FARM MANAGEMENT"

# 3-1  List farms
$farmsResp = PrivTest GET "/farms?limit=20" "List farms (GET /farms)" -v {
    param($r) if ($r.data -eq $null) { return "missing data" }
}
if ($farmsResp) { Write-Detail "Total farms: $($farmsResp.data.total)" }

# 3-2  Get seeded farm
$farmResp = PrivTest GET "/farms/$FARM_ID" "Get farm (GET /farms/{farm_id})" -v {
    param($r) if (-not $r.data.id) { return "missing data.id" }
}
if ($farmResp) { Write-Detail "Farm: $($farmResp.data.name)" }

# 3-3  Update farm
PrivTest PUT "/farms/$FARM_ID" "Update farm name (PUT /farms/{farm_id})" `
    -body "{`"name`":`"Sokha Poultry Farm`"}" | Out-Null

# 3-4  Revert farm name
PrivTest PUT "/farms/$FARM_ID" "Revert farm name (PUT /farms/{farm_id})" `
    -body "{`"name`":`"Test Farm`"}" | Out-Null

# ─────────────────────────────────────────────────────────────────────────────
# 4. COOP MANAGEMENT
# ─────────────────────────────────────────────────────────────────────────────
Write-Section "4. COOP MANAGEMENT"

# 4-1  Create coop
$newCoopBody = @{
    number        = 99
    name          = "Test Coop (automated)"
    capacity      = 200
    current_count = 150
    chicken_type  = "broiler"
} | ConvertTo-Json
$coopCreateResp = PrivTest POST "/farms/$FARM_ID/coops" "Create coop (POST /farms/{farm_id}/coops)" `
    -body $newCoopBody -v {
        param($r) if (-not $r.data.id) { return "missing data.id" }
    }
if ($coopCreateResp) {
    $script:newCoopId = $coopCreateResp.data.id
    Write-Detail "Coop ID: $script:newCoopId"
}

# 4-2  List coops
$coopsResp = PrivTest GET "/farms/$FARM_ID/coops" "List coops (GET /farms/{farm_id}/coops)" -v {
    param($r) if ($r.data -eq $null) { return "missing data" }
}
if ($coopsResp) { Write-Detail "Total coops: $($coopsResp.data.total)" }

# 4-3  Get the newly created coop
if ($script:newCoopId) {
    PrivTest GET "/farms/$FARM_ID/coops/$script:newCoopId" `
        "Get coop (GET /farms/{farm_id}/coops/{coop_id})" | Out-Null
}

# 4-4  Update coop
if ($script:newCoopId) {
    PrivTest PUT "/farms/$FARM_ID/coops/$script:newCoopId" `
        "Update coop (PUT /farms/{farm_id}/coops/{coop_id})" `
        -body "{`"current_count`":175}" | Out-Null
}

# 4-5  Delete coop (cleanup)
if ($script:newCoopId) {
    PrivTest DELETE "/farms/$FARM_ID/coops/$script:newCoopId" `
        "Delete coop (DELETE /farms/{farm_id}/coops/{coop_id})" | Out-Null
}

# ─────────────────────────────────────────────────────────────────────────────
# 5. DEVICE MANAGEMENT
# ─────────────────────────────────────────────────────────────────────────────
Write-Section "5. DEVICE MANAGEMENT"

# 5-1  List devices
$devicesResp = PrivTest GET "/farms/$FARM_ID/devices" "List devices (GET /farms/{farm_id}/devices)" -v {
    param($r) if ($r.data -eq $null) { return "missing data" }
}
if ($devicesResp) { Write-Detail "Total devices: $($devicesResp.data.total)" }

# 5-2  Get seeded device
$devResp = PrivTest GET "/farms/$FARM_ID/devices/$DEVICE_ID" `
    "Get device (GET /farms/{farm_id}/devices/{device_id})" -v {
        param($r) if (-not $r.data.id) { return "missing data.id" }
    }
if ($devResp) { Write-Detail "Device: $($devResp.data.name)  status=$($devResp.data.status)" }

# 5-3  Send command to device
$cmdBody = "{`"command`":`"ping`",`"parameters`":{}}"
$cmdResp = PrivTest POST "/farms/$FARM_ID/devices/$DEVICE_ID/commands" `
    "Send command (POST /farms/{farm_id}/devices/{device_id}/commands)" -body $cmdBody -v {
        param($r) if (-not $r.data.id -and -not $r.data.command_id) { return "missing command id" }
    }
if ($cmdResp) {
    $cmdId = if ($cmdResp.data.id) { $cmdResp.data.id } else { $cmdResp.data.command_id }
    Write-Detail "Command ID: $cmdId"

    # 5-4  Get command status
    if ($cmdId) {
        PrivTest GET "/farms/$FARM_ID/devices/$DEVICE_ID/commands/$cmdId" `
            "Get command status (GET /farms/{farm_id}/devices/{device_id}/commands/{cmd_id})" | Out-Null
    }
}

# 5-5  List commands
PrivTest GET "/farms/$FARM_ID/devices/$DEVICE_ID/commands" `
    "List device commands (GET /farms/{farm_id}/devices/{device_id}/commands)" | Out-Null

# 5-6  Device heartbeat — skipped (requires hardware_id from firmware)
Write-Skip "Device heartbeat (POST /devices/{hardware_id}/heartbeat) — requires device hardware_id"
$script:skipped++

# ─────────────────────────────────────────────────────────────────────────────
# 6. SCHEDULE MANAGEMENT  (includes action_sequence + action_duration)
# ─────────────────────────────────────────────────────────────────────────────
Write-Section "6. SCHEDULE MANAGEMENT  (action_sequence + action_duration)"

# 6-1  List all schedules (shows both seeded schedules)
$schedulesResp = PrivTest GET "/farms/$FARM_ID/schedules" `
    "List schedules (GET /farms/{farm_id}/schedules)" -v {
        param($r) if ($r.data.data -eq $null) { return "missing data.data array" }
    }
if ($schedulesResp) {
    Write-Detail "Total schedules: $($schedulesResp.data.total)"
    foreach ($s in $schedulesResp.data.data) {
        $extra = ""
        if ($s.action_duration) { $extra = "action_duration=$($s.action_duration)s" }
        if ($s.action_sequence)  { $extra = "action_sequence [$($s.action_sequence.Count) steps]" }
        Write-Detail "  📌 $($s.name)  $extra"
    }
}

# 6-2  Get schedule with action_sequence
$seqResp = PrivTest GET "/farms/$FARM_ID/schedules/$SEQ_SCHED_ID" `
    "Get schedule with action_sequence (GET /farms/{farm_id}/schedules/{id})" -v {
        param($r)
        if (-not $r.data.id)              { return "missing data.id" }
        if (-not $r.data.action_sequence) { return "action_sequence should not be null" }
    }
if ($seqResp) { Write-Detail "action_sequence steps: $($seqResp.data.action_sequence.Count)" }

# 6-3  Get schedule with action_duration
$durResp = PrivTest GET "/farms/$FARM_ID/schedules/$DUR_SCHED_ID" `
    "Get schedule with action_duration (GET /farms/{farm_id}/schedules/{id})" -v {
        param($r)
        if (-not $r.data.id)           { return "missing data.id" }
        if (-not $r.data.action_duration) { return "action_duration should not be null" }
    }
if ($durResp) { Write-Detail "action_duration: $($durResp.data.action_duration)s" }

# 6-4  Create schedule with action_sequence
$newSchedBody = @{
    device_id       = $DEVICE_ID
    name            = "Automated Test Schedule"
    schedule_type   = "time_based"
    cron_expression = "0 20 * * *"
    action          = "on"
    action_sequence = @(
        @{ action = "on";  duration_seconds = 30; pause_after_seconds = 5 }
        @{ action = "on";  duration_seconds = 30; pause_after_seconds = 5 }
        @{ action = "off"; duration_seconds = 0;  pause_after_seconds = 0 }
    )
    priority        = 5
    is_active       = $true
} | ConvertTo-Json -Depth 5

$createdResp = PrivTest POST "/farms/$FARM_ID/schedules" `
    "Create schedule with action_sequence (POST /farms/{farm_id}/schedules)" `
    -body $newSchedBody -v {
        param($r)
        if (-not $r.data.id)              { return "missing data.id" }
        if (-not $r.data.action_sequence) { return "action_sequence should be returned" }
    }
if ($createdResp) {
    $script:newScheduleId = $createdResp.data.id
    Write-Detail "New schedule ID: $script:newScheduleId"
    Write-Detail "action_sequence steps: $($createdResp.data.action_sequence.Count)"
}

# 6-5  Update schedule (add action_duration, deactivate)
if ($script:newScheduleId) {
    $updResp = PrivTest PUT "/farms/$FARM_ID/schedules/$script:newScheduleId" `
        "Update schedule action_duration (PUT /farms/{farm_id}/schedules/{id})" `
        -body "{`"action_duration`":1800,`"is_active`":false}" -v {
            param($r)
            if ($r.data.action_duration -ne 1800) { return "action_duration should be 1800" }
        }
    if ($updResp) {
        Write-Detail "action_duration now: $($updResp.data.action_duration)s"
        Write-Detail "is_active now: $($updResp.data.is_active)"
    }
}

# 6-6  Get execution history
if ($script:newScheduleId) {
    PrivTest GET "/farms/$FARM_ID/schedules/$script:newScheduleId/executions" `
        "Get execution history (GET /farms/{farm_id}/schedules/{id}/executions)" | Out-Null
}

# 6-7  Execute-now (trigger manually)
if ($script:newScheduleId) {
    PrivTest POST "/farms/$FARM_ID/schedules/$script:newScheduleId/execute-now" `
        "Execute schedule now (POST /farms/{farm_id}/schedules/{id}/execute-now)" | Out-Null
}

# 6-8  Delete (cleanup)
if ($script:newScheduleId) {
    $delResp = PrivTest DELETE "/farms/$FARM_ID/schedules/$script:newScheduleId" `
        "Delete schedule (DELETE /farms/{farm_id}/schedules/{id})" -v {
            param($r) if (-not $r.message) { return "missing message" }
        }
    if ($delResp) { Write-Detail $delResp.message }
}

# 6-9  Verify count returns to original
$finalResp = PrivTest GET "/farms/$FARM_ID/schedules" `
    "Verify schedule count restored after delete" -v {
        param($r) if ($r.data.total -ne 2) { return "expected 2 schedules, got $($r.data.total)" }
    }
if ($finalResp) { Write-Detail "Schedule count: $($finalResp.data.total) ✓" }

# ─────────────────────────────────────────────────────────────────────────────
# 7. WEBSOCKET
# ─────────────────────────────────────────────────────────────────────────────
Write-Section "7. WEBSOCKET"

# 7-1  WebSocket stats (REST endpoint)
PrivTest GET "/ws/stats" "Get WebSocket stats (GET /ws/stats)" | Out-Null

# 7-2  WebSocket upgrade — skipped (needs a WS client, not Invoke-RestMethod)
Write-Skip "WebSocket upgrade (GET /ws?farm_id=...) — requires WebSocket client"
$script:skipped++

# ─────────────────────────────────────────────────────────────────────────────
# 8. LOGOUT  (at the end so the token stays valid through all tests)
# ─────────────────────────────────────────────────────────────────────────────
Write-Section "8. LOGOUT"

PrivTest POST "/auth/logout" "Logout (POST /auth/logout)" | Out-Null

# ─────────────────────────────────────────────────────────────────────────────
# SUMMARY
# ─────────────────────────────────────────────────────────────────────────────
$total = $script:passed + $script:failed
Write-Host "`n══════════════════════════════════════════════════" -ForegroundColor DarkGray
Write-Host "  TEST SUMMARY" -ForegroundColor Yellow
Write-Host "══════════════════════════════════════════════════" -ForegroundColor DarkGray
Write-Host "  Total:   $total" -ForegroundColor White
Write-Host "  Passed:  $script:passed" -ForegroundColor Green
Write-Host "  Failed:  $script:failed" -ForegroundColor $(if ($script:failed -gt 0) { "Red" } else { "Green" })
Write-Host "  Skipped: $script:skipped" -ForegroundColor DarkGray

if ($script:failed -eq 0) {
    Write-Host "`n  🎉 ALL TESTS PASSED!" -ForegroundColor Green
    exit 0
} else {
    Write-Host "`n  ⚠️  $script:failed TEST(S) FAILED" -ForegroundColor Red
    exit 1
}
