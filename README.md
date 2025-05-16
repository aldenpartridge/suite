# 🛡️ Cybersecurity Tool Suite

A terminal-based suite of cybersecurity tools built with Go and Bubbletea TUI framework.

## Features

- 🔍 Port Scanner
- 🌐 Subdomain Enumerator
- 🔑 Hash Identifier & Cracker
- 🗺️ Network Mapper
- 🔒 HTTP Header Security Scanner
- 📊 Real-time Progress Display
- 💾 Exportable Results (JSON, TXT)
- 🔌 Plugin-ready Architecture

## Requirements

- Go 1.21+
- Terminal with Unicode support

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/suite
cd suite

# Install dependencies
go mod download

# Build the binary
go build -o cybersuite cmd/cybersuite/main.go
```

## Usage

```bash
# Run interactively
./cybersuite

# Run specific tool directly
./cybersuite --tool portscan --target example.com
```

## Development

```bash
# Run tests
go test ./...

# Run linter
golangci-lint run
```

## Security

This tool suite is designed for authorized security testing only. Always ensure you have permission to scan target systems.

## License

MIT License - See LICENSE file for details 