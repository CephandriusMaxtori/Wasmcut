#!/bin/bash
# Build script for compiling Wasm module to wasm32-wasip1

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed"
    exit 1
fi

# Create output directory
mkdir -p "$PROJECT_ROOT/dist"

echo "Building Wasmcut core to wasm32-wasip1..."

cd "$PROJECT_ROOT/wasm"

# Compile to wasm32-wasip1 (WASI Preview 1)
GOOS=wasip1 GOARCH=wasm go build \
    -o "$PROJECT_ROOT/dist/core.wasm" \
    -ldflags="-s -w" \
    main.go

echo "✓ Build complete: dist/core.wasm"
ls -lh "$PROJECT_ROOT/dist/core.wasm"
