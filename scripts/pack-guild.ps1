# Packs a ready-to-run MunkeHub zip for one guild member:
#   - addon folder (MunkeHub)
#   - helper exe + their personal config.json
#   - README.txt with setup steps
#
# Usage:
#   .\scripts\pack-guild.ps1 -Label sven
#   .\scripts\pack-guild.ps1 -Label sven -ApiUrl http://1.2.3.4:4712
#
# Prompts for: admin token (never stored), and the API address friends
# use (defaults to the port-forwarded web port — the helper posts through
# nginx's /api proxy, so friends never need the LAN-only API port).
param(
    [Parameter(Mandatory = $true)]
    [string]$Label,                       # friend's roster label, e.g. "sven"

    [string]$ApiUrl = "",                 # what THEIR helper should call
    [string]$MintServer = "http://brixmedia.local:4711"  # LAN API for minting
)

$ErrorActionPreference = "Stop"
$repo = Split-Path $PSScriptRoot -Parent

# --- 1. API address for the friend (must be reachable from outside) -------
if ($ApiUrl -eq "") {
    $ApiUrl = Read-Host "API address for $Label (e.g. http://<your-public-ip>:4712)"
}
$ApiUrl = $ApiUrl.TrimEnd("/")

# --- 2. Admin token (typed by you, never written anywhere) ----------------
$secure = Read-Host "Admin token (MUNKEHUB_ADMIN_TOKEN)" -AsSecureString
$bstr = [Runtime.InteropServices.Marshal]::SecureStringToBSTR($secure)
try { $admin = [Runtime.InteropServices.Marshal]::PtrToStringBSTR($bstr) }
finally { [Runtime.InteropServices.Marshal]::ZeroFreeBSTR($bstr) }

# --- 3. Mint the member's token -------------------------------------------
try {
    $resp = Invoke-RestMethod -Method Post -Uri "$MintServer/api/v1/tokens" `
        -Headers @{ "X-Admin-Token" = $admin } `
        -ContentType "application/json" `
        -Body (@{ label = $Label } | ConvertTo-Json)
} catch {
    Write-Host "Mint failed: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}
Write-Host "Token minted for '$Label'" -ForegroundColor Green

# --- 4. Build the helper (fresh binary) -----------------------------------
Push-Location (Join-Path $repo "helper")
go build -o munkehub-helper.exe .
if ($LASTEXITCODE -ne 0) { Pop-Location; Write-Host "helper build failed" -ForegroundColor Red; exit 1 }
Pop-Location

# --- 5. Stage the package -------------------------------------------------
$stage = Join-Path $repo "dist\stage-$Label"
if (Test-Path $stage) { Remove-Item $stage -Recurse -Force }
New-Item -ItemType Directory -Path $stage -Force | Out-Null

Copy-Item (Join-Path $repo "addon\MunkeHub") (Join-Path $stage "MunkeHub") -Recurse
Copy-Item (Join-Path $repo "helper\munkehub-helper.exe") $stage
@{ api = $ApiUrl; token = $resp.token } | ConvertTo-Json | Set-Content (Join-Path $stage "config.json")

@'
MunkeHub setup - Munkeklubben
=============================

1. INSTALL THE ADDON
   Copy the "MunkeHub" folder into:
     World of Warcraft\_forever_\Interface\AddOns\

2. START WORLD OF WARCRAFT
   At character select, click "AddOns" (bottom left) and make sure
   MunkeHub is ticked.

3. LOG IN, then type:  /reload
   This saves your character stats to disk.

4. RUN THE HELPER
   Double-click munkehub-helper.exe and leave it running (you can
   minimize the window). It uploads your stats automatically.

5. CHECK THE WEBSITE - your character should appear within ~10 seconds.

NOTES
- Your personal token is inside config.json. Don't share the file.
- The helper must keep running while you play for stats to update.
- Problems? Ask Emil.
'@ | Set-Content (Join-Path $stage "README.txt")

# --- 6. Zip it -----------------------------------------------------------
$zip = Join-Path $repo "dist\MunkeHub-$Label.zip"
if (Test-Path $zip) { Remove-Item $zip -Force }
Compress-Archive -Path "$stage\*" -DestinationPath $zip
Remove-Item $stage -Recurse -Force

Write-Host ""
Write-Host "Package ready: $zip" -ForegroundColor Green
Write-Host "Send it to $Label - they just extract and follow README.txt."