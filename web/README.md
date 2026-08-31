# Wasmcut Web UI

A browser-based video editor interface for Wasmcut, running entirely in WebAssembly.

## Features

- **No Backend Required**: Runs entirely in the browser
- **Wasm Core Engine**: Timeline logic executes in WebAssembly
- **Responsive Design**: Works on desktop and mobile
- **Real-time Editing**: Direct interaction with the Wasm module

## Quick Start

### Local Development

1. Build the Wasm module:
```bash
cd /workspaces/Wasmcut
GOOS=wasip1 GOARCH=wasm go build -o web/core.wasm -ldflags="-s -w" ./wasm
```

2. Serve the web directory locally:
```bash
cd web
python3 -m http.server 8000
```

3. Open http://localhost:8000 in your browser

### GitHub Pages Deployment

The web UI is automatically deployed to GitHub Pages on each push to `main`:

```
https://CephandriusMaxtori.github.io/Wasmcut
```

To enable GitHub Pages:
1. Go to repository Settings > Pages
2. Select "Deploy from a branch"
3. Choose `gh-pages` branch as source

## Architecture

```
Browser (index.html + JavaScript)
    ↓
Wasm Loader (wasm-loader.js)
    ↓
Wasm Module (core.wasm)
    ↓
Timeline Engine (Go + Shared Types)
```

## File Structure

- `index.html` - Main UI page
- `style.css` - Styling (modern dark theme)
- `app.js` - Application logic and UI controls
- `wasm-loader.js` - Wasm module loader and interface
- `core.wasm` - Compiled Wasm module (built from ../wasm/)

## Current Status (M2 - Web UI Proof)

✅ Web UI interface
✅ Wasm module loading in browser
✅ Basic timeline visualization
⏳ Project creation and state management
⏳ Clip editing operations
⏳ Real-time preview

## Development

### Adding New Features

1. Update Wasm exports in `../wasm/main.go`
2. Add corresponding methods in `web/wasm-loader.js`
3. Update UI in `web/app.js` and `web/index.html`
4. Rebuild with `make build-wasm`

### Browser Compatibility

- Chrome/Chromium 74+
- Firefox 79+
- Safari 14.1+
- Edge 79+

Requires WebAssembly support.

## Future Enhancements

- Real-time video preview
- Waveform display
- Keyframe editing
- Effects panel
- Export functionality
- Offline support (PWA)
