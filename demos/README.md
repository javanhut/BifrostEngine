# Bifrost Engine Demos

This directory contains various demos and examples showcasing different features of the Bifrost Engine.

## Available Demos

### Interactive Demos
- **`basic_editor.go`** - Full interactive 3D editor with mouse controls
  - Mouse drag to orbit camera
  - WASD to move objects
  - Tab to switch objects
  - Number keys to add objects

- **`shape_demo.go`** - Interactive shape viewer
  - Press 1-6 to switch between different rendering modes
  - Shows static/rotating triangles and cubes

### Visual Demos
- **`cube_demo.go`** - Animated rotating 3D cube with colored faces
- **`camera_demo.go`** - Camera orbiting around a 3D object
- **`main_demo.go`** - Basic static 3D cube display

### Debug/Test Demos
- **`camera_debug.go`** - Camera matrix debugging tool
- **`simple_render.go`** - Minimal OpenGL rendering test
- **`perspective_check.go`** - Matrix and perspective testing

## How to Run

```bash
# Navigate to demos directory
cd demos

# Run any demo
go run basic_editor.go
go run shape_demo.go
go run cube_demo.go
```

## Recommended Starting Point

Start with **`basic_editor.go`** - it provides the most comprehensive demonstration of the engine's capabilities with full 3D navigation and object manipulation.

## Controls

### Basic Editor
- **Mouse**: Orbit camera around scene
- **Mouse Wheel**: Zoom in/out
- **Tab**: Switch selected object
- **WASD**: Move selected object
- **QE**: Rotate selected object
- **RF**: Scale selected object
- **V**: Toggle object visibility
- **1-2**: Add cube/triangle

### Shape Demo
- **1-6**: Switch rendering modes
- **ESC**: Exit

Enjoy exploring the Bifrost Engine!