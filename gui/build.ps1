<#
.SYNOPSIS
    Builds a GUI application for the Release Package Manager.

.DESCRIPTION
    This script sets up the build environment, configures GCC/CGO, and builds 
    the Fyne GUI application. It can optionally clean previous builds, update 
    dependencies, and run tests.

.PARAMETER Clean
    Remove previous build artifacts and clean Go cache before building.

.PARAMETER Update
    Update Go modules using 'go mod tidy' before building.

.PARAMETER Test
    Run all tests using 'go test ./...' before building.

.PARAMETER Output
    Specify the output executable filename. Default is "release_package_manager_gui.exe".

.EXAMPLE
    .\build_gui.ps1
    Basic build with default settings.

.EXAMPLE
    .\build_gui.ps1 -Clean -Update -Test
    Clean build with dependency updates and tests.

.EXAMPLE
    .\build_gui.ps1 -Output "myapp.exe" -Clean
    Clean build with custom output filename.

.NOTES
    Requires:
    - Go programming language
    - GCC compiler (TDM-GCC or MinGW-w64)
    - CGO enabled environment
#>

# To see help, type on PowerShell prompt: Get-Help .\build_gui.ps1

param(
    [switch]$Clean, # Defaults to $false
    [switch]$Update,
    [switch]$Test,
    [string]$Output = "..\exe\rpmg.exe"
)

Write-Host "Building Fyne GUI Application..." -ForegroundColor Green

# Check if Go is installed
try {
    $goVersion = go version
    Write-Host "Using Go: $goVersion" -ForegroundColor Blue
}
catch {
    Write-Error "Go is not installed or not in PATH. Please install Go first."
    exit 1
}

# Check for GCC, needed for CGO below. Clang did not work for some reason.
if (-not (Get-Command gcc -ErrorAction SilentlyContinue)) {
    Write-Host "GCC not found." -ForegroundColor Yellow
    Write-Host "Please download and install TDM-GCC from: https://jmeubank.github.io/tdm-gcc/" -ForegroundColor Yellow
    Write-Host "Or install MinGW-w64 via: choco install mingw" -ForegroundColor Yellow
    exit 1
}

# Set CGO environment
$env:CGO_ENABLED = "1"

# Test GCC
try {
    $gccVersion = gcc --version
    Write-Host "Using GCC: $($gccVersion[0])" -ForegroundColor Green
}
catch {
    Write-Error "GCC is not working properly"
    exit 1
}

# Clean build if requested
if ($Clean) {
    Write-Host "Cleaning previous build artifacts..." -ForegroundColor Yellow
    if (Test-Path $Output) {
        Remove-Item $Output -Force
        Write-Host "Removed: $Output" -ForegroundColor Green
    }
    go clean -cache
    Write-Host "Go cache cleaned" -ForegroundColor Green
}

# Update dependencies if requested
if ($Update) {
    Write-Host "Updating Go modules..." -ForegroundColor Yellow
    go mod tidy
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Failed to tidy Go modules"
        exit 1
    }
}

# Run tests if requested
if ($Test) {
    Write-Host "Running tests..." -ForegroundColor Yellow
    go test ./...
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Tests failed"
        exit 1
    }
    Write-Host "Tests passed" -ForegroundColor Green
}

# Build the application
Write-Host "Building application: $Output" -ForegroundColor Yellow
go build -o $Output .
if ($LASTEXITCODE -ne 0) {
    Write-Error "Build failed"
    exit 1
}

if (Test-Path $Output) {
    $fileSize = (Get-Item $Output).Length
    Write-Host "Build successful! Output: $Output ($([math]::Round($fileSize/1MB, 2)) MB)" -ForegroundColor Green
}
else {
    Write-Error "Build completed but output file not found: $Output"
    exit 1
}

Write-Host "Build process completed successfully!" -ForegroundColor Green
