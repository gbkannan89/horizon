#!/usr/bin/env bash
set -euo pipefail

echo "=== Horizon Development Setup ==="

# Prerequisites
command -v go >/dev/null 2>&1 || { echo "Error: Go 1.24+ is required"; exit 1; }
command -v docker >/dev/null 2>&1 || { echo "Error: Docker is required"; exit 1; }
command -v docker compose >/dev/null 2>&1 || { echo "Error: Docker Compose is required"; exit 1; }

GO_VERSION=$(go version | grep -oP 'go\K[0-9]+\.[0-9]+')
echo "Go version: $GO_VERSION"

# Download Go dependencies
echo "Downloading Go dependencies..."
go mod download

# Install Go tools
echo "Installing Go tools..."
cd tools && go mod download && cd ..

# Generate protos
echo "Generating protobuf code..."
make proto 2>/dev/null || echo "Proto generation skipped (buf not installed)"

echo ""
echo "=== Setup Complete ==="
echo "Run 'make dev' to start the development environment."
echo "Run 'make test' to run tests."
