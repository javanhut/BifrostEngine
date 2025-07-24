package ui

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type GUISystem struct {
	shader      uint32
	vao         uint32
	vbo         uint32
	windowWidth int
	windowHeight int
	showMainMenu bool
	showObjectMenu bool
	showProjectMenu bool
	showViewMenu bool
	showHelpMenu bool
	mouseX      float64
	mouseY      float64
	leftClickPressed bool
	textRenderer *TextRenderer
	editor      *Editor
	showStatsTable bool
	currentMode string
}

type Button struct {
	X, Y, Width, Height float32
	Text               string
	Clicked            bool
	BackgroundColor    [3]float32
	TextColor          [3]float32
}

type MenuBar struct {
	Height  float32
	Buttons []Button
}

func NewGUISystem(width, height int, editor *Editor) *GUISystem {
	gui := &GUISystem{
		windowWidth:  width,
		windowHeight: height,
		showMainMenu: true,
		textRenderer: NewTextRenderer(),
		editor:      editor,
		showStatsTable: true,
	}
	gui.setupShaders()
	gui.setupBuffers()
	return gui
}

func (gui *GUISystem) setupShaders() {
	vertexShader := `
#version 330 core
layout (location = 0) in vec2 position;
layout (location = 1) in vec3 color;

out vec3 fragColor;

uniform mat4 projection;

void main() {
    gl_Position = projection * vec4(position, 0.0, 1.0);
    fragColor = color;
}
` + "\x00"

	fragmentShader := `
#version 330 core
in vec3 fragColor;
out vec4 FragColor;

void main() {
    FragColor = vec4(fragColor, 1.0);
}
` + "\x00"

	// Compile vertex shader
	vs := gl.CreateShader(gl.VERTEX_SHADER)
	csource, free := gl.Strs(vertexShader)
	gl.ShaderSource(vs, 1, csource, nil)
	free()
	gl.CompileShader(vs)

	// Compile fragment shader
	fs := gl.CreateShader(gl.FRAGMENT_SHADER)
	csource, free = gl.Strs(fragmentShader)
	gl.ShaderSource(fs, 1, csource, nil)
	free()
	gl.CompileShader(fs)

	// Create shader program
	gui.shader = gl.CreateProgram()
	gl.AttachShader(gui.shader, vs)
	gl.AttachShader(gui.shader, fs)
	gl.LinkProgram(gui.shader)

	gl.DeleteShader(vs)
	gl.DeleteShader(fs)
}

func (gui *GUISystem) setupBuffers() {
	gl.GenVertexArrays(1, &gui.vao)
	gl.GenBuffers(1, &gui.vbo)
}

func (gui *GUISystem) Update(mouseX, mouseY float64, leftClick bool) {
	gui.mouseX = mouseX
	gui.mouseY = float64(gui.windowHeight) - mouseY // Flip Y coordinate
	gui.leftClickPressed = leftClick
}

func (gui *GUISystem) SetCurrentMode(mode string) {
	gui.currentMode = mode
}

func (gui *GUISystem) Render() {
	// Enable blending for transparency
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.DEPTH_TEST)

	gl.UseProgram(gui.shader)

	// Set up orthographic projection
	projMatrix := [16]float32{
		2.0 / float32(gui.windowWidth), 0, 0, 0,
		0, 2.0 / float32(gui.windowHeight), 0, 0,
		0, 0, -1, 0,
		-1, -1, 0, 1,
	}
	projLocation := gl.GetUniformLocation(gui.shader, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projLocation, 1, false, &projMatrix[0])

	// Render menu bar
	gui.renderMenuBar()

	// Render dropdown menus if open
	if gui.showObjectMenu {
		gui.renderObjectMenu()
	}
	if gui.showProjectMenu {
		gui.renderProjectMenu()
	}
	if gui.showViewMenu {
		gui.renderViewMenu()
	}
	if gui.showHelpMenu {
		gui.renderHelpMenu()
	}
	
	// Render stats table on top of everything
	if gui.showStatsTable {
		gui.renderStatsTable()
	}
	
	// Render mode indicator in top-right
	gui.renderModeIndicator()

	// Restore OpenGL state
	gl.Enable(gl.DEPTH_TEST)
	gl.Disable(gl.BLEND)
}

func (gui *GUISystem) renderMenuBar() {
	// Menu bar background
	gui.drawRect(0, float32(gui.windowHeight-30), float32(gui.windowWidth), 30, [3]float32{0.2, 0.2, 0.2})

	// Menu buttons - increased widths to accommodate bitmap text
	buttons := []struct {
		x, w float32
		text string
		action func()
	}{
		{10, 65, "File", func() { gui.showProjectMenu = !gui.showProjectMenu; gui.closeOtherMenus("project") }},
		{85, 100, "Objects", func() { gui.showObjectMenu = !gui.showObjectMenu; gui.closeOtherMenus("objects") }},
		{195, 65, "View", func() { gui.showViewMenu = !gui.showViewMenu; gui.closeOtherMenus("view") }},
		{270, 65, "Help", func() { gui.showHelpMenu = !gui.showHelpMenu; gui.closeOtherMenus("help") }},
	}

	for _, btn := range buttons {
		y := float32(gui.windowHeight - 30)
		
		// Check if button is hovered
		hovered := gui.mouseX >= float64(btn.x) && gui.mouseX <= float64(btn.x+btn.w) &&
			gui.mouseY >= float64(y) && gui.mouseY <= float64(y+30)
		
		// Button background
		if hovered {
			gui.drawRect(btn.x, y, btn.w, 30, [3]float32{0.4, 0.4, 0.4})
			if gui.leftClickPressed {
				btn.action()
			}
		} else {
			gui.drawRect(btn.x, y, btn.w, 30, [3]float32{0.3, 0.3, 0.3})
		}

		// Button border
		gui.drawRectOutline(btn.x, y, btn.w, 30, [3]float32{0.6, 0.6, 0.6})

		// Render actual text - better centered in wider buttons
		gui.renderText(btn.x+8, y+6, btn.text, 1.5, [3]float32{1.0, 1.0, 1.0})
	}
}

func (gui *GUISystem) renderObjectMenu() {
	// Dropdown background - wider to accommodate text
	x, y := float32(85), float32(gui.windowHeight-30-175)
	width, height := float32(160), float32(175)
	
	gui.drawRect(x, y, width, height, [3]float32{0.15, 0.15, 0.15})
	gui.drawRectOutline(x, y, width, height, [3]float32{0.5, 0.5, 0.5})

	// Menu items
	items := []string{"Add Cube", "Add Sphere", "Add Cylinder", "Add Plane", "Add Triangle", "Add Pyramid", "Add Light"}
	itemHeight := float32(25)

	for i, item := range items {
		itemY := y + float32(i)*itemHeight
		
		// Check if item is hovered
		hovered := gui.mouseX >= float64(x) && gui.mouseX <= float64(x+width) &&
			gui.mouseY >= float64(itemY) && gui.mouseY <= float64(itemY+itemHeight)

		if hovered {
			gui.drawRect(x, itemY, width, itemHeight, [3]float32{0.3, 0.5, 0.7})
			if gui.leftClickPressed {
				// Convert menu item name to object type
				var objectType string
				switch item {
				case "Add Cube":
					objectType = "cube"
				case "Add Sphere":
					objectType = "sphere"
				case "Add Cylinder":
					objectType = "cylinder"
				case "Add Plane":
					objectType = "plane"
				case "Add Triangle":
					objectType = "triangle"
				case "Add Pyramid":
					objectType = "pyramid"
				case "Add Light":
					objectType = "light"
				default:
					objectType = "cube"
				}
				
				// Create object in editor
				gui.editor.AddObjectWithType(objectType)
				fmt.Printf("Created %s in viewport\n", objectType)
				gui.showObjectMenu = false
			}
		}

		// Render actual text - better centered in wider menus
		gui.renderText(x+8, itemY+4, item, 1.2, [3]float32{0.9, 0.9, 0.9})
	}
}

func (gui *GUISystem) renderProjectMenu() {
	// Dropdown background - wider to accommodate text
	x, y := float32(10), float32(gui.windowHeight-30-100)
	width, height := float32(120), float32(100)
	
	gui.drawRect(x, y, width, height, [3]float32{0.15, 0.15, 0.15})
	gui.drawRectOutline(x, y, width, height, [3]float32{0.5, 0.5, 0.5})

	// Menu items
	items := []string{"New Project", "Open", "Save", "Exit"}
	itemHeight := float32(25)

	for i, item := range items {
		itemY := y + float32(i)*itemHeight
		
		// Check if item is hovered
		hovered := gui.mouseX >= float64(x) && gui.mouseX <= float64(x+width) &&
			gui.mouseY >= float64(itemY) && gui.mouseY <= float64(itemY+itemHeight)

		if hovered {
			gui.drawRect(x, itemY, width, itemHeight, [3]float32{0.3, 0.5, 0.7})
			if gui.leftClickPressed {
				fmt.Printf("Project menu item clicked: %s\n", item)
				gui.showProjectMenu = false
			}
		}

		// Render actual text - better centered in wider menus
		gui.renderText(x+8, itemY+4, item, 1.2, [3]float32{0.9, 0.9, 0.9})
	}
}

func (gui *GUISystem) drawRect(x, y, width, height float32, color [3]float32) {
	vertices := []float32{
		// Position, Color
		x, y, color[0], color[1], color[2],
		x + width, y, color[0], color[1], color[2],
		x + width, y + height, color[0], color[1], color[2],
		x, y + height, color[0], color[1], color[2],
	}

	gl.BindVertexArray(gui.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, gui.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.DYNAMIC_DRAW)

	// Position attribute
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// Color attribute
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(2*4))
	gl.EnableVertexAttribArray(1)

	gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)
}

func (gui *GUISystem) drawRectOutline(x, y, width, height float32, color [3]float32) {
	vertices := []float32{
		// Position, Color
		x, y, color[0], color[1], color[2],
		x + width, y, color[0], color[1], color[2],
		x + width, y + height, color[0], color[1], color[2],
		x, y + height, color[0], color[1], color[2],
	}

	gl.BindVertexArray(gui.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, gui.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.DYNAMIC_DRAW)

	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(2*4))
	gl.EnableVertexAttribArray(1)

	gl.DrawArrays(gl.LINE_LOOP, 0, 4)
}

func (gui *GUISystem) renderText(x, y float32, text string, scale float32, color [3]float32) {
	// Get projection matrix
	projMatrix := [16]float32{
		2.0 / float32(gui.windowWidth), 0, 0, 0,
		0, 2.0 / float32(gui.windowHeight), 0, 0,
		0, 0, -1, 0,
		-1, -1, 0, 1,
	}
	
	gui.textRenderer.RenderText(text, x, y, scale, color, projMatrix)
}

func (gui *GUISystem) RenderText(x, y float32, text string, scale float32, color [3]float32) {
	gui.renderText(x, y, text, scale, color)
}

func (gui *GUISystem) DrawRect(x, y, width, height float32, color [3]float32) {
	gui.drawRect(x, y, width, height, color)
}

func (gui *GUISystem) DrawRectOutline(x, y, width, height float32, color [3]float32) {
	gui.drawRectOutline(x, y, width, height, color)
}

func (gui *GUISystem) Cleanup() {
	gl.DeleteVertexArrays(1, &gui.vao)
	gl.DeleteBuffers(1, &gui.vbo)
	gl.DeleteProgram(gui.shader)
	gui.textRenderer.Cleanup()
}

func (gui *GUISystem) GetObjectMenuVisible() bool {
	return gui.showObjectMenu
}

func (gui *GUISystem) GetProjectMenuVisible() bool {
	return gui.showProjectMenu
}

func (gui *GUISystem) closeOtherMenus(except string) {
	if except != "objects" {
		gui.showObjectMenu = false
	}
	if except != "project" {
		gui.showProjectMenu = false
	}
	if except != "view" {
		gui.showViewMenu = false
	}
	if except != "help" {
		gui.showHelpMenu = false
	}
}

func (gui *GUISystem) renderViewMenu() {
	// Dropdown background
	x, y := float32(195), float32(gui.windowHeight-30-125)
	width, height := float32(120), float32(125)
	
	gui.drawRect(x, y, width, height, [3]float32{0.15, 0.15, 0.15})
	gui.drawRectOutline(x, y, width, height, [3]float32{0.5, 0.5, 0.5})

	// Menu items
	items := []string{"Toggle Grid", "Wireframe", "Fullscreen", "Reset Camera", "Stats"}
	itemHeight := float32(25)

	for i, item := range items {
		itemY := y + float32(i)*itemHeight
		
		// Check if item is hovered
		hovered := gui.mouseX >= float64(x) && gui.mouseX <= float64(x+width) &&
			gui.mouseY >= float64(itemY) && gui.mouseY <= float64(itemY+itemHeight)

		if hovered {
			gui.drawRect(x, itemY, width, itemHeight, [3]float32{0.3, 0.5, 0.7})
			if gui.leftClickPressed {
				switch item {
				case "Toggle Grid":
					grid := gui.editor.GetGrid()
					grid.Visible = !grid.Visible
					fmt.Printf("Grid: %v\n", grid.Visible)
				case "Wireframe":
					fmt.Println("Wireframe mode not yet implemented")
				case "Fullscreen":
					fmt.Println("Fullscreen mode not yet implemented")
				case "Reset Camera":
					fmt.Println("Reset Camera - use R key or implement camera reset")
				case "Stats":
					gui.showStatsTable = !gui.showStatsTable
					fmt.Printf("Stats table: %v\n", gui.showStatsTable)
				}
				gui.showViewMenu = false
			}
		}

		// Render actual text
		gui.renderText(x+8, itemY+4, item, 1.2, [3]float32{0.9, 0.9, 0.9})
	}
}

func (gui *GUISystem) renderHelpMenu() {
	// Dropdown background
	x, y := float32(270), float32(gui.windowHeight-30-125)
	width, height := float32(140), float32(125)
	
	gui.drawRect(x, y, width, height, [3]float32{0.15, 0.15, 0.15})
	gui.drawRectOutline(x, y, width, height, [3]float32{0.5, 0.5, 0.5})

	// Menu items
	items := []string{"Controls", "Shortcuts", "Documentation", "About", "Report Bug"}
	itemHeight := float32(25)

	for i, item := range items {
		itemY := y + float32(i)*itemHeight
		
		// Check if item is hovered
		hovered := gui.mouseX >= float64(x) && gui.mouseX <= float64(x+width) &&
			gui.mouseY >= float64(itemY) && gui.mouseY <= float64(itemY+itemHeight)

		if hovered {
			gui.drawRect(x, itemY, width, itemHeight, [3]float32{0.3, 0.5, 0.7})
			if gui.leftClickPressed {
				fmt.Printf("Help menu item clicked: %s\n", item)
				gui.showHelpMenu = false
			}
		}

		// Render actual text
		gui.renderText(x+8, itemY+4, item, 1.2, [3]float32{0.9, 0.9, 0.9})
	}
}

func (gui *GUISystem) renderStatsTable() {
	objects := gui.editor.GetSceneObjects()
	selectedIndex := gui.editor.GetSelectedObject()
	
	if selectedIndex >= 0 && selectedIndex < len(objects) {
		obj := objects[selectedIndex]
		
		// Table properties - positioned in bottom-right  
		tableWidth := float32(280)
		tableHeight := float32(180) // Increased height to fit all content
		x := float32(gui.windowWidth) - tableWidth - 10
		y := float32(gui.windowHeight) - tableHeight - 10
		
		// Draw table background with higher Z-order
		tableColor := [3]float32{0.05, 0.05, 0.1} // Darker background
		borderColor := [3]float32{0.6, 0.6, 0.7} // Brighter border
		gui.drawRect(x, y, tableWidth, tableHeight, tableColor)
		gui.drawRectOutline(x, y, tableWidth, tableHeight, borderColor)
		
		// Table header with distinct color (at the TOP of the table)
		headerHeight := float32(24)
		headerColor := [3]float32{0.2, 0.3, 0.5} // Blue header
		headerY := y + tableHeight - headerHeight // Position header at top
		gui.drawRect(x, headerY, tableWidth, headerHeight, headerColor)
		gui.renderText(x+5, headerY+4, fmt.Sprintf("Object: %s", obj.Name), 1.3, [3]float32{1.0, 1.0, 1.0})
		
		// Table rows (start from top, go down)
		rowHeight := float32(14) // Reduced row height to fit better
		textScale := float32(1.0) // Slightly smaller text
		labelColor := [3]float32{0.9, 0.9, 0.9} // Bright labels
		valueColor := [3]float32{1.0, 1.0, 0.7} // Bright yellow values
		
		currentY := headerY - 6 // Start closer to header
		
		// Position section
		currentY -= rowHeight
		gui.renderText(x+5, currentY, "Position:", textScale, labelColor)
		currentY -= rowHeight
		gui.renderText(x+15, currentY, fmt.Sprintf("X: %.2f", obj.Position.X), textScale, valueColor)
		currentY -= rowHeight
		gui.renderText(x+15, currentY, fmt.Sprintf("Y: %.2f", obj.Position.Y), textScale, valueColor)
		currentY -= rowHeight
		gui.renderText(x+15, currentY, fmt.Sprintf("Z: %.2f", obj.Position.Z), textScale, valueColor)
		currentY -= 4 // Smaller gap between sections
		
		// Rotation section
		currentY -= rowHeight
		gui.renderText(x+5, currentY, "Rotation:", textScale, labelColor)
		currentY -= rowHeight
		gui.renderText(x+15, currentY, fmt.Sprintf("X: %.1f°", obj.Rotation.X), textScale, valueColor)
		currentY -= rowHeight
		gui.renderText(x+15, currentY, fmt.Sprintf("Y: %.1f°", obj.Rotation.Y), textScale, valueColor)
		currentY -= rowHeight
		gui.renderText(x+15, currentY, fmt.Sprintf("Z: %.1f°", obj.Rotation.Z), textScale, valueColor)
	}
}

func (gui *GUISystem) renderModeIndicator() {
	if gui.currentMode == "" {
		return
	}
	
	// Mode indicator properties - positioned in bottom-right of viewport
	indicatorWidth := float32(160)
	indicatorHeight := float32(32)
	x := float32(gui.windowWidth) - indicatorWidth - 10
	y := float32(10) // Bottom of screen with small margin
	
	// Background colors based on mode
	var bgColor [3]float32
	var textColor [3]float32 = [3]float32{1.0, 1.0, 1.0} // White text
	
	switch gui.currentMode {
	case "select":
		bgColor = [3]float32{0.2, 0.3, 0.8} // Blue for select
	case "move":
		bgColor = [3]float32{0.2, 0.8, 0.3} // Green for move
	case "transform":
		bgColor = [3]float32{0.8, 0.3, 0.2} // Red for transform
	default:
		bgColor = [3]float32{0.5, 0.5, 0.5} // Gray for unknown
	}
	
	// Draw background
	gui.drawRect(x, y, indicatorWidth, indicatorHeight, bgColor)
	gui.drawRectOutline(x, y, indicatorWidth, indicatorHeight, [3]float32{0.8, 0.8, 0.8})
	
	// Format mode text for display
	modeText := ""
	switch gui.currentMode {
	case "select":
		modeText = "SELECT"
	case "move":
		modeText = "MOVE"
	case "transform":
		modeText = "TRANSFORM"
	default:
		modeText = gui.currentMode
	}
	
	// Render mode text centered
	gui.renderText(x+8, y+6, "Mode: " + modeText, 1.1, textColor)
}