# Goxfer

A simple and fast file transfer tool that allows you to share files across your local network with a beautiful CLI interface.

## Features

- Fast file transfers over local network
- QR code generation for easy access
- Beautiful terminal UI with progress tracking
- Support for both uploads and downloads
- Web-based interface accessible from any device

## Installation

### Quick Install (Recommended)

Run the install script which automatically detects your platform:

```bash
curl -fsSL https://raw.githubusercontent.com/Tony-ArtZ/goxfer/main/install.sh | bash
```

Or with wget:

```bash
wget -qO- https://raw.githubusercontent.com/Tony-ArtZ/goxfer/main/install.sh | bash
```

For Windows (PowerShell):

```powershell
irm https://raw.githubusercontent.com/Tony-ArtZ/goxfer/main/install.ps1 | iex
```

### Manual Installation

1. Download the latest release for your platform from the [releases page](https://github.com/Tony-ArtZ/goxfer/releases)
2. Extract the archive
3. Move the binary to a directory in your PATH

### Building from Source

If you have Go installed:

```bash
git clone https://github.com/Tony-ArtZ/goxfer.git
cd goxfer
go build -o goxfer
```

## Usage

Run `goxfer` to start the application:

```bash
goxfer
```

The application will:

1. Start a local web server on an available port (starting from 8080)
2. Display your local IP address and a QR code
3. Show an interactive CLI for managing file transfers

### Sharing Files

Once the application is running, you can share files by entering their path in the CLI:

1. Enter the path to the file or directory you want to share
2. The file will be queued and available for download
3. Users on the network can download the file via the web interface

Access the web interface by:

- Scanning the QR code with your phone
- Opening the displayed URL in any browser on your network

## How It Works

Goxfer creates a local HTTP server that allows devices on your network to:

- Download files that you explicitly queue from the CLI
- Upload files to the server directory
- View transfer progress in real-time

All transfers happen directly over your local network - no internet connection required!

## Requirements

- No external dependencies required for the binary
- Works on Linux, macOS, and Windows
- Requires devices to be on the same local network

## Troubleshooting

If you are unable to access the file transfer interface from another device:

1. Ensure both devices are connected to the same local network
2. Check your firewall settings to ensure the application is allowed to accept incoming connections
3. Verify that the port (default 8080) is not being blocked

## License

MIT License - see LICENSE file for details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Author

Created by Tony-ArtZ
