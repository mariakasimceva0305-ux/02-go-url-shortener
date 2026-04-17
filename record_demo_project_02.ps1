$ErrorActionPreference = "Stop"

Set-Location $PSScriptRoot

$baseUrl = "http://localhost:8080"
$goExe = "go"
if (Test-Path "C:\Program Files\Go\bin\go.exe") {
    $goExe = "C:\Program Files\Go\bin\go.exe"
}

Write-Host "=== START PROJECT 02 ==="
Write-Host "Starting server in a separate window..."
Start-Process powershell -ArgumentList @(
    "-NoExit",
    "-Command",
    "cd '$PSScriptRoot'; & '$goExe' mod tidy; & '$goExe' run ./cmd/app"
)

Write-Host "Waiting for server (up to 60 seconds)..."
$isUp = $false
for ($i = 1; $i -le 30; $i++) {
    try {
        $null = Invoke-RestMethod -Method Get -Uri "$baseUrl/health" -TimeoutSec 2
        $isUp = $true
        break
    } catch {
        Start-Sleep -Seconds 2
    }
}
if (-not $isUp) {
    throw "Server did not start in time. Check the second PowerShell window for errors."
}

function Show-Step($title) {
    Write-Host ""
    Write-Host "=== $title ===" -ForegroundColor Cyan
}

try {
    Show-Step "1) GET /health"
    $health = Invoke-RestMethod -Method Get -Uri "$baseUrl/health"
    $health | ConvertTo-Json -Depth 10
    Start-Sleep -Seconds 2

    Show-Step "2) GET /swagger"
    $swagger = Invoke-RestMethod -Method Get -Uri "$baseUrl/swagger"
    $swagger | ConvertTo-Json -Depth 10
    Start-Sleep -Seconds 2

    Show-Step "3) POST /api/v1/links"
    $createBodyObject = @{
        original_url = "https://example.com/article"
        ttl_hours    = 24
    }
    $createBodyJson = $createBodyObject | ConvertTo-Json
    $created = Invoke-RestMethod -Method Post -Uri "$baseUrl/api/v1/links" -ContentType "application/json" -Body $createBodyJson
    $created | ConvertTo-Json -Depth 10

    $id = $created.id
    $code = $created.short_code
    Write-Host "ID: $id"
    Write-Host "CODE: $code"
    Start-Sleep -Seconds 2

    Show-Step "4) GET /api/v1/links"
    $all = Invoke-RestMethod -Method Get -Uri "$baseUrl/api/v1/links"
    $all | ConvertTo-Json -Depth 10
    Start-Sleep -Seconds 2

    Show-Step "5) GET /api/v1/links/:id"
    $one = Invoke-RestMethod -Method Get -Uri "$baseUrl/api/v1/links/$id"
    $one | ConvertTo-Json -Depth 10
    Start-Sleep -Seconds 2

    Show-Step "6) GET /api/v1/links/:id/analytics BEFORE redirect"
    $before = Invoke-RestMethod -Method Get -Uri "$baseUrl/api/v1/links/$id/analytics"
    $before | ConvertTo-Json -Depth 10
    Start-Sleep -Seconds 2

    Show-Step "7) OPEN /r/:code in browser"
    $redirectUrl = "$baseUrl/r/$code"
    Write-Host "Opening: $redirectUrl"
    Start-Process $redirectUrl

    Read-Host "Show redirected page in browser, then press Enter here"

    Show-Step "8) GET /api/v1/links/:id/analytics AFTER redirect"
    $after = Invoke-RestMethod -Method Get -Uri "$baseUrl/api/v1/links/$id/analytics"
    $after | ConvertTo-Json -Depth 10
    Start-Sleep -Seconds 2

    Show-Step "9) OPTIONAL DELETE /api/v1/links/:id"
    Write-Host "If needed, run:"
    Write-Host "Invoke-RestMethod -Method Delete -Uri '$baseUrl/api/v1/links/$id'"

    Write-Host ""
    Write-Host "Demo flow finished."
    Write-Host "Server is still running in the second PowerShell window."
} catch {
    Write-Host ""
    Write-Host "ERROR: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "Open the second PowerShell window and check why server did not start/respond."
}

Read-Host "Press Enter to close this demo window"
