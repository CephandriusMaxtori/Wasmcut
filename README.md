# Wasmcut

A DaVinci Resolve-like video editor built with **WASM/WASI** for portable, sandboxed editing logic.

🎬 **[Try the Web UI](https://CephandriusMaxtori.github.io/Wasmcut)** (M2 - Proof of Concept)

---

## Overview

Wasmcut separates the **editing engine** (compiled to WebAssembly) from **I/O-heavy operations** (codec, GPU, UI) in the host application. The Wasm core handles timeline logic, edit commands, and project files - everything portable and sandboxed.

### Key Design

```
Host Application (native)          │  Wasm Core Engine (Go → wasm32-wasip1)
- GPU Rendering (OpenGL/Vulkan)    │  - Timeline model
- FFmpeg decode/encode (cgo)       │  - Edit operations & undo/redo
- Windowing & UI                   │  - Project serialization
```

### Current Status

| Milestone | Status | Features |
|-----------|--------|----------|
| **M0** | ✅ Done | Compile proof: Go→WASM, runs standalone |
| **M1** | ✅ Done | Timeline engine: full edit ops + undo/redo + tests |
| **M2** | ✅ Done | Web UI (GitHub Pages): browser-based editor interface |
| **M3** | 🔄 Next | Media pipeline: FFmpeg decode, frame display |
| **M4** | ⏳ TBD | Multi-track compositing, transitions, preview |
| **M5** | ⏳ TBD | Export: timeline → FFmpeg filter graph → output file |

## Quick Start

### Web UI (M2 - Now Available!)

**[Live Demo](https://CephandriusMaxtori.github.io/Wasmcut)** - No installation needed, runs in browser.

Or run locally:

```bash
cd web
python3 -m http.server 8000
# Open http://localhost:8000
```

### Build from Source

```bash
# M1: Test the timeline engine
go test -v ./shared

# M2: Compile Wasm + run web UI
make build-wasm
cd web && python3 -m http.server 8000
```

## Project Structure

```
.
├── wasm/              # Core Wasm module (Go → wasm32-wasip1)
├── host/              # Host application (Go, loads Wasm)
├── shared/            # Shared data structures & timeline logic
├── web/               # Web UI (HTML/CSS/JS) for GitHub Pages
├── docs/              # Design docs & architecture
└── scripts/           # Build & deployment scripts
```

## Architecture Details

See [DESIGN.md](docs/DESIGN.md) for full architecture overview, including:

- Host ABI (Wasm ⇄ Host boundary)
- Project file format
- Toolchain decisions (Go, TinyGo, wasmtime, FFmpeg)
- Milestone breakdown

## Development

### Building Components

```bash
# Timeline engine tests
make test-wasm

# Compile Wasm to wasm32-wasip1
make build-wasm

# Web UI (requires core.wasm in web/)
cd web && python3 -m http.server
```

### Wasm Module Exports

Currently exported (M1/M2):

- `CreateProject()` → create new timeline
- `GetTimelineState()` → query current state
- `AddClip(trackID, mediaRef, in, out, position)` → add clip
- `Undo()` / `Redo()` → timeline undo/redo

Future (M3+):

- `TrimClip()`, `MoveClip()`, `DeleteClip()`
- Transition/effect management
- Project save/load

## Deployment

GitHub Actions automatically builds and deploys the web UI to GitHub Pages on push to `main`.

Manual deployment:

```bash
bash scripts/deploy.sh
```

## Technical Stack

| Component | Tech | Notes |
|-----------|------|-------|
| Wasm Core | Go + TinyGo | Target: wasm32-wasip1 (WASI Preview 1) |
| Runtime | wasmtime | Embedded in native host (M3+) |
| Web UI | HTML/CSS/JS | No framework, vanilla for simplicity |
| Host (ref) | Go | Reference impl; FFmpeg + GPU rendering (future) |
| Codecs | FFmpeg | Via cgo in host, *not* in Wasm (enforced) |
| Tests | Go testing | 10+ tests, all passing |

## Contributing

Wasmcut is in active development. Areas to help:

- M3: Media pipeline (FFmpeg integration)
- M4: Real-time preview compositing
- M5: Export & encoding
- Web UI improvements (clip dragging, waveforms, effects panel)
- TinyGo evaluation for smaller binary size

---

**Next**: Check out [web/README.md](web/README.md) for web UI details or [docs/DESIGN.md](docs/DESIGN.md) for architecture.
