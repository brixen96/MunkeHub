# Packs a ready-to-run MunkeHub zip for one guild member:
#   - addon folder (MunkeHub)
#   - helper exe + their personal config.json
#   - README.txt with setup steps
#
# First run asks for API address + admin token once and saves them to
# scripts\pack-config.json (gitignored). After that:
#
#   .\scripts\pack-guild.ps1 -Label sven
#
# is all it takes. Use -Reconfigure to change the saved settings.
param(
    [Parameter(Mandatory = $true)]
    [string]$Label,                        # friend's roster label, e.g. "sven"

    [switch]$Reconfigure,                  # re-prompt and overwrite saved config
    [string]$ApiUrl = "",                  # override: what THEIR helper calls
    [string]$MintServer = ""               # override: LAN API for minting
)

$ErrorActionPreference = "Stop"
$repo = Split-Path $PSScriptRoot -Parent
$cfgPath = Join-Path $PSScriptRoot "pack-config.json"

# --- 1. Load or create saved config ---------------------------------------
$cfg = @{}
if (Test-Path $cfgPath) {
    $cfg = Get-Content $cfgPath -Raw | ConvertFrom-Json -AsHashtable
}

if ($Reconfigure -or -not $cfg.apiUrl -or -not $cfg.adminToken) {
    if (-not $ApiUrl) {
        $ApiUrl = Read-Host "API address friends use (e.g. http://<your-public-ip>:4712)"
    }
    if (-not $MintServer) {
        $MintServer = Read-Host "LAN API for minting (default http://brixmedia.local:4711)"
        if (-not $MintServer) { $MintServer = "http://brixmedia.local:4711" }
    }
    $secure = Read-Host "Admin token (MUNKEHUB_ADMIN_TOKEN)" -AsSecureString
    $bstr = [Runtime.InteropServices.Marshal]::SecureStringToBSTR($secure)
    try { $admin = [Runtime.InteropServices.Marshal]::PtrToStringBSTR($bstr) }
    finally { [Runtime.InteropServices.Marshal]::ZeroFreeBSTR($bstr) }

    $cfg = @{
        apiUrl     = $ApiUrl.TrimEnd("/")
        mintServer = $MintServer.TrimEnd("/")
        adminToken = $admin
    }
    $cfg | ConvertTo-Json | Set-Content $cfgPath
    Write-Host "Settings saved to scripts\pack-config.json (gitignored)" -ForegroundColor Green
}

# Command-line overrides beat saved config
if ($ApiUrl)     { $cfg.apiUrl     = $ApiUrl.TrimEnd("/") }
if ($MintServer) { $cfg.mintServer = $MintServer.TrimEnd("/") }

# --- 2. Mint the member's token -------------------------------------------
try {
    $resp = Invoke-RestMethod -Method Post -Uri "$($cfg.mintServer)/api/v1/tokens" `
        -Headers @{ "X-Admin-Token" = $cfg.adminToken } `
        -ContentType "application/json" `
        -Body (@{ label = $Label } | ConvertTo-Json)
} catch {
    Write-Host "Mint failed: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "Tip: run with -Reconfigure if the admin token or server changed." -ForegroundColor Yellow
    exit 1
}
Write-Host "Token minted for '$Label'" -ForegroundColor Green

# --- 3. Build the helper (fresh binary) -----------------------------------
Push-Location (Join-Path $repo "helper")
go build -o munkehub-helper.exe .
if ($LASTEXITCODE -ne 0) { Pop-Location; Write-Host "helper build failed" -ForegroundColor Red; exit 1 }
Pop-Location

# --- 4. Stage the package -------------------------------------------------
$stage = Join-Path $repo "dist\stage-$Label"
if (Test-Path $stage) { Remove-Item $stage -Recurse -Force }
New-Item -ItemType Directory -Path $stage -Force | Out-Null

Copy-Item (Join-Path $repo "addon\MunkeHub") (Join-Path $stage "MunkeHub") -Recurse
Copy-Item (Join-Path $repo "helper\munkehub-helper.exe") $stage
@{ api = $cfg.apiUrl; token = $resp.token } | ConvertTo-Json | Set-Content (Join-Path $stage "config.json")

@'
MunkeHub setup - Munkeklubben
=============================

1. INSTALL THE ADDON
   Copy the "MunkeHub" folder into:
     World of Warcraft\_forever_\Interface\AddOns\

2. START WORLD OF WARCRAFT
   At character select, click "AddOns" (bottom left) and make sure
   MunkeHub is ticked.

3. SET YOUR LABEL (once)
   Log in and type:  /mh label <yourname>
   Then type:  /reload
   Your stats now broadcast to the guild automatically.

4. RUN THE HELPER (optional - for instant sync)
   Double-click munkehub-helper.exe and leave it running.
   Without it, your stats sync when Emil is online in-game.

5. CHECK THE WEBSITE - your character should appear shortly.

NOTES
- Your personal token is inside config.json. Don't share the file.
- Step 3 alone is enough - nothing else to run.
- Problems? Ask Emil.
'@ | Set-Content (Join-Path $stage "README.txt")

# --- 5. Zip it -----------------------------------------------------------
$zip = Join-Path $repo "dist\MunkeHub-$Label.zip"
if (Test-Path $zip) { Remove-Item $zip -Force }
Compress-Archive -Path "$stage\*" -DestinationPath $zip
Remove-Item $stage -Recurse -Force

Write-Host ""
Write-Host "Package ready: $zip" -ForegroundColor Green
Write-Host "Send it to $Label - they just extract and follow README.txt."
explorer.exe (Join-Path $repo "dist") | Out-Null  # open the folder for you