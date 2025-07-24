# Bifrost Engine UI Editor

The Bifrost Engine UI Editor provides a comprehensive interface for creating and manipulating 3D scenes with an intuitive graphical user interface.

## Features

### Object Creation
- **Quick Add Toolbar**: Easily add primitives (Cube, Sphere, Cylinder, Plane) with one-click buttons
- **Menu System**: Access all object types through organized menus
- **Object Categories**:
  - Primitives: Cube, Sphere, Cylinder, Plane, Cone, Torus, Triangle
  - Lights: Point lights, directional lights (coming soon)
  - Cameras: Additional viewport cameras (coming soon)

### Project Management
- Create new projects with custom names
- Open existing projects from a list
- Save current project state
- Recent projects quick access
- Project metadata tracking (creation date, last modified)

### Scene Hierarchy
- Visual tree view of all objects in the scene
- Click to select objects
- Right-click context menu for object operations:
  - Delete object
  - Duplicate object
- Visual indication of selected object

### Inspector Panel
- Edit selected object properties:
  - Name
  - Transform (Position, Rotation, Scale)
  - Visibility toggle
  - Color/appearance settings
- Real-time updates as you modify values

### Grid System
- Visual reference grid for object placement
- Toggleable visibility
- Customizable size and divisions
- Helps with spatial orientation

### Camera Controls
- **Mouse Controls**:
  - Left-click + drag: Orbit camera around scene
  - Scroll wheel: Zoom in/out
- Smooth camera movement
- Camera position displayed in stats panel

### Toolbar
- Quick access to common tools:
  - Select tool
  - Move tool
  - Rotate tool
  - Scale tool
- Current project display

## Running the UI Editor

```bash
cd demos
go run ui_editor.go
```

## Keyboard Shortcuts

- Tab: Switch between objects (basic editor mode)
- WASD: Move selected object (basic editor mode)
- Q/E: Rotate selected object (basic editor mode)
- R/F: Scale selected object (basic editor mode)
- V: Toggle object visibility (basic editor mode)

## Architecture

The UI system is built with:
- **ImGui**: Immediate mode GUI library for responsive interfaces
- **Modular Design**: Separate components for different UI panels
- **Scene Integration**: Direct manipulation of scene graph
- **Real-time Rendering**: Changes reflected immediately in viewport

## Extending the Editor

To add new object types:
1. Add the type to `ObjectType` enum in `ui/objects.go`
2. Add template to `ObjectTemplates` array
3. Implement rendering in the scene renderer

To add new UI panels:
1. Create render method in `editor.go`
2. Add visibility flag to Editor struct
3. Add menu item to toggle visibility
4. Call render method in main render loop

## Future Enhancements

- Transform gizmos for visual manipulation
- Multi-selection support
- Undo/redo system
- Asset browser
- Material editor
- Animation timeline
- Physics properties panel