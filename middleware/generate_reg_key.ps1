# Generate Registration Key for On-Site Setup
param(
    [string]$FarmName,
    [string]$CustomerName,
    [string]$CustomerPhone,
    [string]$Location,
    [int]$ExpiryDays = 90
)

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Registration Key Generator" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

# Generate random registration key (format: XXXXX-XXXXX-XXXXX-XXXXX-XXXXX)
function Generate-RegKey {
    $chars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"  # Exclude confusing chars (0,O,1,I)
    $key = ""
    for ($i = 0; $i -lt 5; $i++) {
        if ($i -gt 0) { $key += "-" }
        $segment = ""
        for ($j = 0; $j -lt 5; $j++) {
            $segment += $chars[(Get-Random -Minimum 0 -Maximum $chars.Length)]
        }
        $key += $segment
    }
    return $key
}

$regKey = Generate-RegKey
$keyId = [guid]::NewGuid().ToString()
$expiresAt = (Get-Date).AddDays($ExpiryDays).ToString("yyyy-MM-dd HH:mm:ss")

Write-Host "`nGenerated Registration Key:" -ForegroundColor Yellow
Write-Host "  Key: $regKey" -ForegroundColor Green
Write-Host "  ID: $keyId" -ForegroundColor Gray
Write-Host "  Expires: $expiresAt" -ForegroundColor Gray

if ($FarmName) {
    Write-Host "`nCustomer Details:" -ForegroundColor Yellow
    Write-Host "  Farm: $FarmName" -ForegroundColor Gray
    if ($CustomerName) { Write-Host "  Name: $CustomerName" -ForegroundColor Gray }
    if ($CustomerPhone) { Write-Host "  Phone: $CustomerPhone" -ForegroundColor Gray }
    if ($Location) { Write-Host "  Location: $Location" -ForegroundColor Gray }
}

Write-Host "`nSQL to insert into database:" -ForegroundColor Cyan
$sql = @"
INSERT INTO registration_keys 
(id, key_code, farm_name, customer_name, customer_phone, farm_location, expires_at, created_by, created_at)
VALUES 
('$keyId', '$regKey', $(if($FarmName){"'$FarmName'"}else{"NULL"}), $(if($CustomerName){"'$CustomerName'"}else{"NULL"}), $(if($CustomerPhone){"'$CustomerPhone'"}else{"NULL"}), $(if($Location){"'$Location'"}else{"NULL"}), '$expiresAt', 'admin', CURRENT_TIMESTAMP);
"@

Write-Host $sql -ForegroundColor White

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "Give this key to field staff for on-site registration" -ForegroundColor Yellow
Write-Host "========================================" -ForegroundColor Cyan
