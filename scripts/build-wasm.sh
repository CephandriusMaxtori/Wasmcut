#!/bin/bash
# Build script for compiling the browser Wasm module

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

echo "Building Wasmcut core for the browser..."

cd "$PROJECT_ROOT/wasm"

# Compile to the Go JavaScript/Wasm runtime target.
GOOS=js GOARCH=wasm go build \
    -o "$PROJECT_ROOT/dist/core.wasm" \
    -ldflags="-s -w" \
    main.go

echo "✓ Build complete: dist/core.wasm"
ls -lh "$PROJECT_ROOT/dist/core.wasm"

cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" "$PROJECT_ROOT/dist/wasm_exec.js"
echo "✓ Copied Go browser runtime: dist/wasm_exec.js"
