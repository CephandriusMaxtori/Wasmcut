#!/bin/bash
# Script to build and deploy to GitHub Pages

set -e

echo "🔨 Building Wasmcut Web..."

# Build Wasm module
mkdir -p dist
GOOS=wasip1 GOARCH=wasm go build -o dist/core.wasm -ldflags="-s -w" ./wasm
echo "✓ Wasm compiled"

# Copy web files to dist
mkdir -p dist/web
cp web/index.html dist/web/
cp web/style.css dist/web/
cp web/app.js dist/web/
cp web/wasm-loader.js dist/web/
cp dist/core.wasm dist/web/

echo "✓ Web files copied to dist/web"
echo ""
echo "🚀 Ready for deployment!"
echo ""
echo "To deploy to GitHub Pages:"
echo "  1. Push to GitHub"
echo "  2. Go to Settings > Pages"
echo "  3. Select 'Deploy from a branch'"
echo "  4. Choose 'main' and '/dist' as source"
