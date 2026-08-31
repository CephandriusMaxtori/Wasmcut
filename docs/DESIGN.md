# Wasmcut — WASM/WASI Video Editor Design Doc

## 1. Overview

A video editor whose **core logic** (timeline model, edit operations, project
format) is compiled to WebAssembly and run under a WASI runtime, embedded in
a native host application that handles GPU rendering, codecs, and UI.

**Goal:** portable, sandboxed, embeddable editing engine — not a from-scratch
GPU/codec stack.

**Non-goals (for v1):** real-time GPU compositing inside Wasm, in-Wasm codec
decode, WASI Preview 2 / Component Model (not yet mature enough for this).

---

## 2. Architecture

```
┌─────────────────────────────────────────────────────┐
│                     Host Application (Go)             │
│                                                         │
│  ┌───────────────┐  ┌───────────────┐  ┌────────────┐ │
│  │   Windowing/UI │  │  FFmpeg (cgo) │  │  GPU render │ │
│  │  (native, cgo) │  │ decode/encode │  │ (OpenGL/    │ │
│  │                │  │               │  │  Vulkan)    │ │
│  └───────┬────────┘  └───────┬───────┘  └─────┬──────┘ │
│          │                    │                 │        │
│          └──────────┬─────────┴────────┬────────┘        │
│                      │  Host ABI calls  │                 │
│              ┌───────▼──────────────────▼───────┐         │
│              │     wasmtime-go (or wasmer-go)     │        │
│              │  loads & runs the Wasm module      │        │
│              └───────────────┬────────────────────┘        │
└──────────────────────────────┼──────────────────────────────┘
                                │
                  ┌─────────────▼─────────────┐
                  │   Wasm Core Engine (Go/    │
                  │   TinyGo → wasm32-wasip1)  │
                  │                             │
                  │  - Timeline data model      │
                  │  - Edit commands            │
                  │  - Project (de)serialization│
                  │  - Effect/parameter graph   │
                  └─────────────────────────────┘
```

The Wasm module is the **portable brain**: timeline logic, no GPU, no codecs,
no windowing. The host does everything I/O- and hardware-heavy and calls into
the module for edit-state logic.

---

## 3. Components

### 3.1 Core Engine (Wasm/WASI)

- Language: Go or TinyGo, target `wasm32-wasip1`.
- Responsibilities:

  - Timeline/track/clip data model
  - Edit operations: cut, trim, reorder, split, ripple-delete
  - Transition & effect parameter graph (data only — no rendering)
  - Project file format: load/save (JSON or binary)
  - Undo/redo command stack
  - Explicitly excluded: codec access, GPU calls, file I/O beyond project data
    (cgo is unavailable under the Wasm target, so this is enforced by the
    toolchain, not just convention).

### 3.2 Host Application (native Go)

- Embeds the Wasm module via `wasmtime-go` (or `wasmer-go`).
- Owns:

  - Windowing/UI (timeline widget, preview window, waveform display)
  - FFmpeg integration (via cgo bindings) for decode/encode
  - GPU compositing/preview rendering (OpenGL/Vulkan/Metal, chosen per platform)
  - Media pool / file management on disk
  - Talks to the Wasm module through a defined host ABI (see §4).

### 3.3 Rendering/Compositing

- Lives host-side, not in Wasm, for v1.
- Preview: host decodes frames (FFmpeg) → GPU shader compositing per the
  timeline state fetched from the Wasm module → draw to screen.
- Export: host walks the timeline (via Wasm queries), builds an FFmpeg
  filter graph (concat/overlay/xfade) or does frame-by-frame compositing,
  then encodes.

---

## 4. Host ABI (Wasm ⇄ Host boundary)

Minimal command/query interface, versioned from the start.

**Host → Wasm (commands/queries):**

- `create_project() -> project_id`
- `load_project(bytes) -> project_id`
- `save_project(project_id) -> bytes`
- `add_clip(project_id, media_ref, track, in, out, position)`
- `trim_clip(project_id, clip_id, new_in, new_out)`
- `move_clip(project_id, clip_id, new_track, new_position)`
- `delete_clip(project_id, clip_id)`
- `set_transition(project_id, clip_a, clip_b, type, duration)`
- `get_timeline_state(project_id) -> serialized_state`
- `undo(project_id)` / `redo(project_id)`

**Wasm → Host (imports, if needed):**

- Kept minimal in v1 — ideally the Wasm module is pure logic with no
  outbound calls, communicating only via return values. Avoids relying on
  immature WASI import/threading behavior.

Data crosses the boundary as serialized bytes (JSON to start; consider
FlatBuffers/CBOR later if perf matters).

---

## 5. Project File Format (draft)

```
{
  "version": 1,
  "tracks": [
    {
      "id": "track-1",
      "type": "video",
      "clips": [
        {
          "id": "clip-1",
          "media_ref": "media://source1.mp4",
          "in": 0.0,
          "out": 5.2,
          "position": 0.0
        }
      ]
    }
  ],
  "transitions": [],
  "effects": []
}
```

---

## 6. Toolchain Decisions

| Concern | Decision | Notes |
|---------|----------|-------|
| Wasm core language | Go, evaluate TinyGo | TinyGo for smaller binary/startup; confirm stdlib subset covers needs |
| Wasm target | `wasm32-wasip1` | Preview 2/Component Model not yet used |
| Wasm runtime (host) | wasmtime (via wasmtime-go) | wasmer as fallback if wasmtime gaps found |
| Codec handling | Host-side FFmpeg via cgo | No cgo inside Wasm target — enforced constraint |
| Rendering | Host-side GPU (OpenGL/Vulkan) | Not attempted inside Wasm for v1 |
| Serialization | JSON for v1 | Revisit if perf becomes an issue |

---

## 7. Milestones

1. **M0 — Compile loop proof:** minimal Go/TinyGo module compiled to
   `wasm32-wasip1`, run standalone via `wasmtime` CLI, returns a hardcoded
   timeline JSON.
2. **M1 — Timeline engine:** full data model + edit commands + undo/redo,
   tested via the `wasmtime` CLI or a Go test harness calling the module.
3. **M2 — Host embedding:** native Go host loads the module via
   `wasmtime-go`, drives it through the ABI, prints/serializes state.
4. **M3 — Media pipeline:** host-side FFmpeg decode of a real file, frames
   drawn to a window (no compositing yet — single track playback).
5. **M4 — Basic compositing:** multi-track overlay, simple cut transitions,
   real-time-ish preview.
6. **M5 — Export:** timeline → FFmpeg filter graph → encoded output file.

---

## 8. Status Update (As of M2)

### ✅ Completed

- **M0:** Wasm module compiles to wasm32-wasip1 (2.2MB binary)
- **M1:** Full timeline engine with 10 passing unit tests
  - Command pattern for all edit operations
  - Undo/redo stack with proper state management
  - JSON serialization/deserialization
  - Track/clip/transition/effect data models
- **M2:** Web UI for GitHub Pages
  - Browser-based editor interface
  - Modern dark theme with responsive design
  - Wasm module loader (no WASI required)
  - GitHub Actions workflow for auto-deployment
  - Timeline visualization and controls

### ⏳ Next (M3)

- FFmpeg integration for media decode
- Frame display and preview
- Host-side rendering pipeline

---

## 9. Open Questions

- TinyGo vs standard Go: need to confirm TinyGo's stdlib supports full timeline/undo-stack needs.
- Memory limits: Wasm linear memory is 32-bit addressable; need to validate large project handling.
- Serialization format: JSON may get slow for complex timelines with many clips/keyframes.
- Whether/when to revisit WASI Preview 2 / Component Model as Go tooling matures.
- How much of the effect/transition *parameter* model to put in Wasm vs. leave as opaque data the host interprets.

