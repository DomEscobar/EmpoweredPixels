<#
.SYNOPSIS
    Starts the Empowered Pixels development environment.

.DESCRIPTION
    This script launches the Go backend and Vue frontend in separate terminal windows.
    It checks for necessary dependencies and ensures the environment is ready.

.NOTES
    File Name      : run-dev.ps1
    Author         : AI Assistant
    Prerequisite   : Go, Node.js, PostgreSQL
#>

$ErrorActionPreference = "Stop"

function Show-Header {
    Write-Host "==========================================" -ForegroundColor Cyan
    Write-Host "   Empowered Pixels - Dev Launcher" -ForegroundColor Cyan
    Write-Host "==========================================" -ForegroundColor Cyan
    Write-Host ""
}

function Check-Command ($cmd) {
    if (-not (Get-Command $cmd -ErrorAction SilentlyContinue)) {
        Write-Error "$cmd is not installed or not in your PATH."
        exit 1
    }
}

function Start-Service {
    param (
        [string]$Name,
        [string]$Path,
        [string]$Command,
        [string]$Color
    )
    
    Write-Host "Starting $Name..." -ForegroundColor $Color
    $psi = New-Object System.Diagnostics.ProcessStartInfo
    $psi.FileName = "powershell.exe"
    $psi.Arguments = "-NoExit", "-Command", "& { Write-Host '--- $Name ---' -ForegroundColor $Color; cd '$Path'; $Command }"
    $psi.UseShellExecute = $true
    [System.Diagnostics.Process]::Start($psi) | Out-Null
}

Show-Header

# 1. Prerequisite Checks
Write-Host "Checking prerequisites..." -ForegroundColor Yellow
Check-Command "go"
Check-Command "npm"
Write-Host "Prerequisites met." -ForegroundColor Green
Write-Host ""

$root = $PSScriptRoot
$backendDir = Join-Path $root "backend"
$frontendDir = Join-Path $root "frontend"

# 2. Setup Frontend Dependencies
if (-not (Test-Path (Join-Path $frontendDir "node_modules"))) {
    Write-Host "Frontend dependencies missing. Installing..." -ForegroundColor Yellow
    Push-Location $frontendDir
    try {
        npm install
    }
    finally {
        Pop-Location
    }
    Write-Host "Dependencies installed." -ForegroundColor Green
}

# 3. Launch Services
Start-Service -Name "Backend (Go)" -Path $backendDir -Command "go run ./cmd/api/main.go" -Color "Cyan"
Start-Service -Name "Frontend (Vite)" -Path $frontendDir -Command "npm run dev" -Color "Magenta"

Write-Host ""
Write-Host "Done!" -ForegroundColor Green
Write-Host "Backend API: http://localhost:54321" -ForegroundColor Gray
Write-Host "Frontend:    http://localhost:5173 (usually)" -ForegroundColor Gray
Write-Host "Ensure your local PostgreSQL is running on port 5432." -ForegroundColor DarkGray
