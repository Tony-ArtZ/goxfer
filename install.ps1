# Goxfer Installation Script for Windows
# Run with: irm https://raw.githubusercontent.com/Tony-ArtZ/goxfer/main/install.ps1 | iex

$ErrorActionPreference = "Stop"

$REPO = "Tony-ArtZ/goxfer"
$BINARY_NAME = "goxfer.exe"
$INSTALL_DIR = "$env:LOCALAPPDATA\goxfer"

function Write-Info {
    param($Message)
    Write-Host "[INFO] $Message" -ForegroundColor Green
}

function Write-Warn {
    param($Message)
    Write-Host "[WARN] $Message" -ForegroundColor Yellow
}

function Write-Error-Custom {
    param($Message)
    Write-Host "[ERROR] $Message" -ForegroundColor Red
    exit 1
}

function Detect-Platform {
    $arch = [System.Environment]::Is64BitOperatingSystem
    if ($arch) {
        return "amd64"
    } else {
        return "386"
    }
}

function Get-LatestVersion {
    Write-Info "Fetching latest release..."
    try {
        $response = Invoke-RestMethod -Uri "https://api.github.com/repos/$REPO/releases/latest"
        $version = $response.tag_name
        Write-Info "Latest version: $version"
        return $version
    } catch {
        Write-Error-Custom "Failed to fetch latest version: $_"
    }
}

function Install-Binary {
    param($Version, $Arch)
    
    $downloadUrl = "https://github.com/$REPO/releases/download/$Version/goxfer-windows-$Arch.exe"
    Write-Info "Downloading from: $downloadUrl"
    
    # Create install directory if it doesn't exist
    if (!(Test-Path $INSTALL_DIR)) {
        New-Item -ItemType Directory -Path $INSTALL_DIR -Force | Out-Null
    }
    
    $binaryPath = Join-Path $INSTALL_DIR $BINARY_NAME
    
    try {
        Invoke-WebRequest -Uri $downloadUrl -OutFile $binaryPath
        Write-Info "Installed to: $binaryPath"
    } catch {
        Write-Error-Custom "Failed to download binary: $_"
    }
}

function Add-ToPath {
    $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
    
    if ($currentPath -notlike "*$INSTALL_DIR*") {
        Write-Info "Adding $INSTALL_DIR to PATH..."
        [Environment]::SetEnvironmentVariable(
            "Path",
            "$currentPath;$INSTALL_DIR",
            "User"
        )
        $env:Path = "$env:Path;$INSTALL_DIR"
        Write-Info "Added to PATH. You may need to restart your terminal."
    }
}

function Verify-Installation {
    $binaryPath = Join-Path $INSTALL_DIR $BINARY_NAME
    if (Test-Path $binaryPath) {
        Write-Info "Installation successful!"
        Write-Info "Run 'goxfer' to get started"
    } else {
        Write-Error-Custom "Installation verification failed"
    }
}

function Main {
    Write-Host "╔═══════════════════════════════════╗"
    Write-Host "║   Goxfer Installation Script      ║"
    Write-Host "╚═══════════════════════════════════╝"
    Write-Host ""
    
    $arch = Detect-Platform
    Write-Info "Detected architecture: windows-$arch"
    
    $version = Get-LatestVersion
    Install-Binary -Version $version -Arch $arch
    Add-ToPath
    Verify-Installation
}

Main
