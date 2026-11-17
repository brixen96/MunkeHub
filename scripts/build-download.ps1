# Builds the GENERIC MunkeHub setup package served by the website's
# Download button. Contains NO token — the helper asks for it on first run.
#
# Run before building the API image (the zip is baked into it):
#   .\scripts\build-download.ps1
#   docker buildx build ... api/
$ErrorActionPreference = "Stop"
$repo = Split-Path $PSScriptRoot -Parent

# 1. Build the helper (Windows exe — friends run Windows)
Push-Location (Join-Path $repo "helper")
go build -o munkehub-helper.exe .
if ($LASTEXITCODE -ne 0) { Pop-Location; Write-Host "helper build failed" -ForegroundColor Red; exit 1 }
Pop-Location

# 2. Stage: addon + helper + README (no config.json — no token!)
$stage = Join-Path $repo "api\download\_stage"
$out = Join-Path $repo "api\download\MunkeHub-setup.zip"
if (Test-Path $stage) { Remove-Item $stage -Recurse -Force }
New-Item -ItemType Directory -Path $stage -Force | Out-Null

Copy-Item (Join-Path $repo "addon\MunkeHub") (Join-Path $stage "MunkeHub") -Recurse
Copy-Item (Join-Path $repo "helper\munkehub-helper.exe") $stage

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
   Your stats now broadcast to the guild automatically - you do NOT
   need to run any program.

4. OPTIONAL: INSTANT SYNC (one-time setup)
   a. Double-click munkehub-helper.exe and paste your personal token
      (ask Emil for it).
   b. Then type:  munkehub-helper.exe -install
      This makes it start automatically with Windows, silently in the
      background - you never need to open it again.

NOTES
- Step 3 is enough for most people - nothing else to run.
- To remove auto-start later:  munkehub-helper.exe -uninstall
- Problems? Ask Emil.
'@ | Set-Content (Join-Path $stage "README.txt")

# 3. Zip
if (Test-Path $out) { Remove-Item $out -Force }
Compress-Archive -Path "$stage\*" -DestinationPath $out
Remove-Item $stage -Recurse -Force

Write-Host "Package ready: $out" -ForegroundColor Green
Write-Host "Now build/push the API image - the zip is baked into it."