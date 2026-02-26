# SnapMark

SnapMark is a cross-platform screenshot + annotation utility built with Go and Fyne.

## Features

- Capture modes:
  - Full screen (2-second delay)
  - Region selection (drag to capture)
  - Window capture (select from open windows)
- Annotation tools:
  - Rectangle
  - Ellipse
  - Arrow
  - Text placement
  - Pixelate (10x10 mosaic blocks, non-destructive overlay baked on save)
- Configurable:
  - Colour picker
  - Stroke width
  - Font size (text tool)
- Undo stack (10 levels)
- Output:
  - Save As PNG/JPG
  - Copy to clipboard (platform-specific)

## Build

Requires Go 1.22+

```bash
cd snapmark
go mod tidy
make build
```

### Platform builds

```bash
make build-linux
make build-windows
make build-mac
```

### Run locally

```bash
make run
```

## Linux dependencies

For clipboard copy support on Linux, install one of:

- `xclip` (preferred)
- `xsel` (fallback)

For window list capture on Linux, install:

- `wmctrl`

## Project layout

```text
snapmark/
├── main.go
├── go.mod
├── go.sum
├── Makefile
└── internal/
    ├── capture/
    ├── editor/
    └── actions/
```
