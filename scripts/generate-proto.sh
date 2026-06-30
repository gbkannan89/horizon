#!/usr/bin/env bash
set -euo pipefail

cd proto

# Find all .proto files
PROTO_FILES=$(find . -name '*.proto' -type f)

if [ -z "$PROTO_FILES" ]; then
  echo "No proto files found."
  exit 0
fi

echo "Generating code for proto files..."
buf generate

echo "Proto generation complete."
