# Bifrost Game Engine

A modern game development engine built from scratch in Go, supporting both 2D and 3D game development.

## Quick Start

### Running the Editor
```bash
go run main.go editor
```

### Running the GUI Overlay Editor (Recommended)
```bash
# Run the advanced editor with full GUI
go run ./demos/gui_overlay_editor.go
```

### Running Demos
```bash
# List all available demos
go run main.go demos

# Run a specific demo
go run main.go demo ui_editor
go run main.go demo cube_demo
```

## Features

### Core Engine
- Modular architecture with Go workspaces
- OpenGL 4.1 rendering backend
- Component-based scene management
- Comprehensive math library (vectors, matrices, quaternions)
- Input handling system (keyboard, mouse)
- Camera systems (2D/3D, orbit controls)

### Editor
- **GUI Overlay Editor**: Advanced editor with full menu system
- **Interactive 3D Viewport**: Real-time object manipulation
- **Multiple Transform Modes**: Select, Move, Transform with visual indicators
- **Object Creation**: Via GUI menus and keyboard shortcuts (F1-F6)
- **Camera System**: Mouse orbit controls with reset (R key)
- **Grid Visualization**: Toggle via menu or G key
- **Object Management**: Selection, movement, scaling, deletion
- **Real-time Stats**: Object properties display with position/rotation
- **Mode Indicator**: Visual feedback for current editor mode

### Rendering
- **Multi-Shape Support**: Cube, Sphere, Cylinder, Plane, Triangle, Pyramid
- **Proper Geometry**: Each shape renders with correct mesh topology
- **Transform System**: Position, rotation, scale with proper matrix handling
- **Visual Highlighting**: Selected object highlighting
- **Grid Rendering**: Configurable grid overlay
- **Shader Management**: OpenGL 4.1 shader compilation and management

## Project Structure
```
BifrostEngine/
├── main.go           # Main entry point
├── engine/           # Core engine systems
├── renderer/         # Rendering subsystem
│   ├── core/        # Renderer core
│   ├── opengl/      # OpenGL implementation
│   └── window/      # Window management
├── math/            # Mathematics library
├── camera/          # Camera systems
├── scene/           # Scene management
├── input/           # Input handling
├── ui/              # UI framework
├── demos/           # Example applications
└── docs/            # Documentation
```

## Commands

### Main Commands
- `bifrost_engine editor` - Launch the full UI editor
- `bifrost_engine demos` - List available demos
- `bifrost_engine demo <name>` - Run a specific demo
- `bifrost_engine --version` - Show version information
- `bifrost_engine --help` - Display help

### GUI Overlay Editor Controls

#### Camera Controls
- **Mouse**: Left-click + drag to orbit camera around scene
- **Scroll**: Zoom in/out 
- **R**: Reset camera to default position

#### Transform Modes
- **Q**: Select Mode (click objects to select them)
- **M**: Move Mode (drag objects with mouse or use arrow keys)
- **T**: Transform Mode (drag to scale/rotate objects)

#### Object Creation
- **GUI Menu**: Objects → Add [Shape]
- **F1/1**: Add Cube
- **F2/2**: Add Sphere  
- **F3/3**: Add Cylinder
- **F4/4**: Add Plane
- **F5/5**: Add Triangle
- **F6/6**: Add Pyramid

#### Object Manipulation
- **Arrow Keys**: Move selected object (in Move mode)
- **X/Y/Z**: Constrain movement to specific axis
- **Tab**: Select next object
- **Delete**: Remove selected object

#### View Controls
- **View Menu**: Toggle Grid, Stats, etc.
- **G**: Toggle grid visibility

#### Interface Elements
- **Mode Indicator**: Bottom-right corner shows current mode
- **Stats Table**: Bottom-right displays selected object properties
- **Object Menu**: Top menu bar for creation and project management

## Development Status

See [docs/PROGRESS.md](docs/PROGRESS.md) for detailed development progress and roadmap.

## Building from Source

### Prerequisites
- Go 1.19 or later
- OpenGL 4.1 capable graphics card
- C compiler (for CGO dependencies)

### Build Steps
```bash
# Clone the repository
git clone https://github.com/yourusername/BifrostEngine.git
cd BifrostEngine

# Install dependencies
go mod download

# Run the engine directly
go run main.go editor

# Or build the binary
go build -o bifrost_engine main.go

# Then use the binary
./bifrost_engine editor
./bifrost_engine demos
```

## Documentation

- [UI Editor Guide](docs/ui-editor.md) - Comprehensive editor usage guide
- [Progress Tracking](docs/PROGRESS.md) - Development roadmap and status
- [Development Log](docs/DEVELOPMENT_LOG.md) - Recent discoveries and technical solutions
- API documentation (coming soon)

## License

This project is open source and available under the MIT License.

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## Acknowledgments

- Dear ImGui for the immediate mode GUI
- GLFW for window management
- go-gl for OpenGL bindings
