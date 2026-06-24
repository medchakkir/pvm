# PVM Installer for Windows
# Usage: irm https://github.com/medchakkir/pvm/releases/latest/download/install.ps1 | iex

param(
    [string]$Version = "latest"
)

$ErrorActionPreference = "Stop"

$Repo    = "medchakkir/pvm"
$InstallDir = "$env:USERPROFILE\.pvm\bin"

Write-Host "PVM Installer" -ForegroundColor Cyan
Write-Host "=============" -ForegroundColor Cyan

# Resolve version
if ($Version -eq "latest") {
    Write-Host "Fetching latest release..." -ForegroundColor Gray
    $release = Invoke-RestMethod "https://api.github.com/repos/$Repo/releases/latest"
    $Version = $release.tag_name -replace '^v', ''
}

Write-Host "Installing PVM v$Version..." -ForegroundColor White

# Detect architecture
$arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else {
    Write-Error "PVM requires a 64-bit Windows system."
    exit 1
}

# Build download URL
$zipName  = "pvm-v${Version}-windows-${arch}.zip"
$url      = "https://github.com/$Repo/releases/download/v$Version/$zipName"
$tmpZip   = "$env:TEMP\$zipName"

# Download
Write-Host "Downloading $zipName..." -ForegroundColor Gray
$ProgressPreference = "SilentlyContinue" # speeds up Invoke-WebRequest significantly
Invoke-WebRequest -Uri $url -OutFile $tmpZip -UseBasicParsing

# Extract
Write-Host "Extracting..." -ForegroundColor Gray
if (-not (Test-Path $InstallDir)) {
    New-Item -ItemType Directory -Path $InstallDir | Out-Null
}
Expand-Archive -Path $tmpZip -DestinationPath $InstallDir -Force
Remove-Item $tmpZip

# Add to PATH if not already there
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -notlike "*$InstallDir*") {
    Write-Host "Adding $InstallDir to your PATH..." -ForegroundColor Gray
    [Environment]::SetEnvironmentVariable("Path", "$userPath;$InstallDir", "User")
    $env:Path += ";$InstallDir"
}

Write-Host ""
Write-Host "PVM v$Version installed successfully!" -ForegroundColor Green
Write-Host ""
Write-Host "Get started:" -ForegroundColor White
Write-Host "  pvm list-remote        # see available PHP versions"
Write-Host "  pvm install 8.3        # install PHP 8.3"
Write-Host "  pvm use 8.3            # switch to PHP 8.3"
Write-Host "  pvm bin                # show PATH setup info"
Write-Host ""
Write-Host "Restart your terminal for PATH changes to take effect." -ForegroundColor Yellow
