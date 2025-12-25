#!/bin/bash
set -e

# Goxfer Installation Script
# This script installs the latest version of goxfer

REPO="Tony-ArtZ/goxfer"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
BINARY_NAME="goxfer"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case "$OS" in
        linux*)
            OS="linux"
            ;;
        darwin*)
            OS="darwin"
            ;;
        mingw*|msys*|cygwin*)
            OS="windows"
            ;;
        *)
            error "Unsupported operating system: $OS"
            ;;
    esac
    
    case "$ARCH" in
        x86_64|amd64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        armv7l)
            ARCH="arm"
            ;;
        i386|i686)
            ARCH="386"
            ;;
        *)
            error "Unsupported architecture: $ARCH"
            ;;
    esac
    
    info "Detected platform: $OS-$ARCH"
}

# Get the latest release version
get_latest_version() {
    info "Fetching latest release..."
    LATEST_VERSION=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    
    if [ -z "$LATEST_VERSION" ]; then
        error "Failed to fetch latest version"
    fi
    
    info "Latest version: $LATEST_VERSION"
}

# Download and install binary
install_binary() {
    DOWNLOAD_URL="https://github.com/$REPO/releases/download/$LATEST_VERSION/goxfer-$OS-$ARCH"
    
    if [ "$OS" = "windows" ]; then
        DOWNLOAD_URL="${DOWNLOAD_URL}.exe"
        BINARY_NAME="${BINARY_NAME}.exe"
    fi
    
    info "Downloading from: $DOWNLOAD_URL"
    
    TEMP_FILE=$(mktemp)
    
    if ! curl -fsSL "$DOWNLOAD_URL" -o "$TEMP_FILE"; then
        error "Failed to download binary. Please check if the release exists for your platform."
    fi
    
    # Create install directory if it doesn't exist
    mkdir -p "$INSTALL_DIR"
    
    # Move binary to install directory
    mv "$TEMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
    
    info "Installed to: $INSTALL_DIR/$BINARY_NAME"
}

# Check if install directory is in PATH
check_path() {
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        warn "Install directory ($INSTALL_DIR) is not in your PATH"
        warn "Add the following line to your shell configuration file:"
        echo ""
        echo "  export PATH=\"\$PATH:$INSTALL_DIR\""
        echo ""
    fi
}

# Verify installation
verify_installation() {
    if [ -f "$INSTALL_DIR/$BINARY_NAME" ]; then
        info "Installation successful!"
        info "Run '$BINARY_NAME' to get started"
    else
        error "Installation verification failed"
    fi
}

# Main installation process
main() {
    echo "╔═══════════════════════════════════╗"
    echo "║   Goxfer Installation Script      ║"
    echo "╚═══════════════════════════════════╝"
    echo ""
    
    detect_platform
    get_latest_version
    install_binary
    check_path
    verify_installation
}

main
