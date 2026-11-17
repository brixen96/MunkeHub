# Syncs the MunkeHub addon source from the repo into the WoW: Forever AddOns folder.
$src = Join-Path $PSScriptRoot "..\addon\MunkeHub"
$dstRoot = "C:\Program Files (x86)\World of Warcraft\_forever_\Interface\AddOns"

Copy-Item -Path $src -Destination $dstRoot -Recurse -Force
Write-Host "Synced MunkeHub addon -> $dstRoot\MunkeHub"