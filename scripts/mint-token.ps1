# Mints a MunkeHub member token and writes helper/config.json so the
# helper can run with zero arguments.
#
# Usage (you type the admin token yourself when prompted — it never
# appears in chat or in shell history):
#   .\scripts\mint-token.ps1 -Server http://<server-ip>:4711
param(
    [Parameter(Mandatory = $true)]
    [string]$Server,          # e.g. http://192.168.1.50:4711
    [string]$Label = "emil"   # your roster label
)

$ErrorActionPreference = "Stop"

$secure = Read-Host "Admin token (MUNKEHUB_ADMIN_TOKEN)" -AsSecureString
$bstr = [Runtime.InteropServices.Marshal]::SecureStringToBSTR($secure)
try {
    $admin = [Runtime.InteropServices.Marshal]::PtrToStringBSTR($bstr)
} finally {
    [Runtime.InteropServices.Marshal]::ZeroFreeBSTR($bstr)
}

try {
    $resp = Invoke-RestMethod -Method Post -Uri "$Server/api/v1/tokens" `
        -Headers @{ "X-Admin-Token" = $admin } `
        -ContentType "application/json" `
        -Body (@{ label = $Label } | ConvertTo-Json)
} catch {
    Write-Host "Mint failed: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "Check the server address and admin token, then retry." -ForegroundColor Yellow
    exit 1
}

$cfgPath = Join-Path $PSScriptRoot "..\helper\config.json"
@{ api = $Server; token = $resp.token } | ConvertTo-Json | Set-Content -Path $cfgPath

Write-Host "OK - token minted for '$Label'" -ForegroundColor Green
Write-Host "Config written to helper\config.json (gitignored - keep it private)"