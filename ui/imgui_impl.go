package ui

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/inkyblackness/imgui-go/v4"
)

type ImGuiContext struct {
	context   *imgui.Context
	io        imgui.IO
	renderer  *OpenGL3
	platform  *GLFW
}

func NewImGuiContext(window *glfw.Window) *ImGuiContext {
	context := imgui.CreateContext(nil)
	io := imgui.CurrentIO()
	
	// Setup ImGui style
	imgui.StyleColorsDark()
	
	// Setup platform and renderer
	platform := NewGLFW(io, window)
	renderer := NewOpenGL3()
	
	return &ImGuiContext{
		context:  context,
		io:       io,
		renderer: renderer,
		platform: platform,
	}
}

func (ctx *ImGuiContext) NewFrame() {
	ctx.platform.NewFrame()
	imgui.NewFrame()
}

func (ctx *ImGuiContext) Render() {
	imgui.Render()
	ctx.renderer.Render(imgui.RenderedDrawData())
}

func (ctx *ImGuiContext) Destroy() {
	ctx.renderer.Destroy()
	ctx.platform.Destroy()
	ctx.context.Destroy()
}

// GLFW Platform implementation
type GLFW struct {
	io              imgui.IO
	window          *glfw.Window
	time            float64
	mouseJustPressed [3]bool
}

func NewGLFW(io imgui.IO, window *glfw.Window) *GLFW {
	platform := &GLFW{
		io:     io,
		window: window,
	}
	
	platform.setKeyMapping()
	platform.installCallbacks()
	
	return platform
}

func (platform *GLFW) setKeyMapping() {
	// Keyboard mapping
	keys := map[int]int{
		imgui.KeyTab:        int(glfw.KeyTab),
		imgui.KeyLeftArrow:  int(glfw.KeyLeft),
		imgui.KeyRightArrow: int(glfw.KeyRight),
		imgui.KeyUpArrow:    int(glfw.KeyUp),
		imgui.KeyDownArrow:  int(glfw.KeyDown),
		imgui.KeyPageUp:     int(glfw.KeyPageUp),
		imgui.KeyPageDown:   int(glfw.KeyPageDown),
		imgui.KeyHome:       int(glfw.KeyHome),
		imgui.KeyEnd:        int(glfw.KeyEnd),
		imgui.KeyInsert:     int(glfw.KeyInsert),
		imgui.KeyDelete:     int(glfw.KeyDelete),
		imgui.KeyBackspace:  int(glfw.KeyBackspace),
		imgui.KeySpace:      int(glfw.KeySpace),
		imgui.KeyEnter:      int(glfw.KeyEnter),
		imgui.KeyEscape:     int(glfw.KeyEscape),
		imgui.KeyA:          int(glfw.KeyA),
		imgui.KeyC:          int(glfw.KeyC),
		imgui.KeyV:          int(glfw.KeyV),
		imgui.KeyX:          int(glfw.KeyX),
		imgui.KeyY:          int(glfw.KeyY),
		imgui.KeyZ:          int(glfw.KeyZ),
	}
	
	for imguiKey, glfwKey := range keys {
		platform.io.KeyMap(imguiKey, glfwKey)
	}
}

func (platform *GLFW) installCallbacks() {
	platform.window.SetMouseButtonCallback(platform.mouseButtonCallback)
	platform.window.SetScrollCallback(platform.scrollCallback)
	platform.window.SetKeyCallback(platform.keyCallback)
	platform.window.SetCharCallback(platform.charCallback)
}

func (platform *GLFW) mouseButtonCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press && button >= 0 && button < 3 {
		platform.mouseJustPressed[button] = true
	}
}

func (platform *GLFW) scrollCallback(window *glfw.Window, xoff, yoff float64) {
	platform.io.AddMouseWheelDelta(float32(xoff), float32(yoff))
}

func (platform *GLFW) keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Press {
		platform.io.KeyPress(int(key))
	}
	if action == glfw.Release {
		platform.io.KeyRelease(int(key))
	}
	
	// Modifiers
	platform.io.KeyCtrl(int(mods&glfw.ModControl) != 0)
	platform.io.KeyShift(int(mods&glfw.ModShift) != 0)
	platform.io.KeyAlt(int(mods&glfw.ModAlt) != 0)
	platform.io.KeySuper(int(mods&glfw.ModSuper) != 0)
}

func (platform *GLFW) charCallback(window *glfw.Window, char rune) {
	platform.io.AddInputCharacter(char)
}

func (platform *GLFW) NewFrame() {
	// Setup display size
	w, h := platform.window.GetSize()
	platform.io.SetDisplaySize(imgui.Vec2{X: float32(w), Y: float32(h)})
	
	// Setup time step
	currentTime := glfw.GetTime()
	if platform.time > 0 {
		platform.io.SetDeltaTime(float32(currentTime - platform.time))
	}
	platform.time = currentTime
	
	// Setup mouse
	if platform.window.GetAttrib(glfw.Focused) == glfw.True {
		x, y := platform.window.GetCursorPos()
		platform.io.SetMousePosition(imgui.Vec2{X: float32(x), Y: float32(y)})
	} else {
		platform.io.SetMousePosition(imgui.Vec2{X: -1, Y: -1})
	}
	
	for i := 0; i < len(platform.mouseJustPressed); i++ {
		down := platform.mouseJustPressed[i] || platform.window.GetMouseButton(glfw.MouseButton(i)) == glfw.Press
		platform.io.SetMouseButtonDown(i, down)
		platform.mouseJustPressed[i] = false
	}
}

func (platform *GLFW) Destroy() {
	// Cleanup callbacks would go here
}

// OpenGL3 Renderer implementation
type OpenGL3 struct {
	shaderHandle      uint32
	vertHandle        uint32
	fragHandle        uint32
	attribLocationTex int32
	attribLocationProjMtx int32
	attribLocationVtxPos int32
	attribLocationVtxUV int32
	attribLocationVtxColor int32
	vboHandle         uint32
	elementsHandle    uint32
	fontTexture       uint32
}

func NewOpenGL3() *OpenGL3 {
	renderer := &OpenGL3{}
	renderer.createDeviceObjects()
	return renderer
}

func (renderer *OpenGL3) createDeviceObjects() {
	// Vertex shader
	vertexShaderSource := `
#version 130
uniform mat4 ProjMtx;
in vec2 Position;
in vec2 UV;
in vec4 Color;
out vec2 Frag_UV;
out vec4 Frag_Color;
void main() {
    Frag_UV = UV;
    Frag_Color = Color;
    gl_Position = ProjMtx * vec4(Position.xy, 0, 1);
}
` + "\x00"

	// Fragment shader
	fragmentShaderSource := `
#version 130
uniform sampler2D Texture;
in vec2 Frag_UV;
in vec4 Frag_Color;
out vec4 Out_Color;
void main() {
    Out_Color = Frag_Color * texture(Texture, Frag_UV.st);
}
` + "\x00"

	// Create shaders
	renderer.vertHandle = renderer.createShader(gl.VERTEX_SHADER, vertexShaderSource)
	renderer.fragHandle = renderer.createShader(gl.FRAGMENT_SHADER, fragmentShaderSource)
	
	// Create program
	renderer.shaderHandle = gl.CreateProgram()
	gl.AttachShader(renderer.shaderHandle, renderer.vertHandle)
	gl.AttachShader(renderer.shaderHandle, renderer.fragHandle)
	gl.LinkProgram(renderer.shaderHandle)
	
	// Get uniform locations
	renderer.attribLocationTex = gl.GetUniformLocation(renderer.shaderHandle, gl.Str("Texture\x00"))
	renderer.attribLocationProjMtx = gl.GetUniformLocation(renderer.shaderHandle, gl.Str("ProjMtx\x00"))
	
	// Get attribute locations
	renderer.attribLocationVtxPos = gl.GetAttribLocation(renderer.shaderHandle, gl.Str("Position\x00"))
	renderer.attribLocationVtxUV = gl.GetAttribLocation(renderer.shaderHandle, gl.Str("UV\x00"))
	renderer.attribLocationVtxColor = gl.GetAttribLocation(renderer.shaderHandle, gl.Str("Color\x00"))
	
	// Create buffers
	gl.GenBuffers(1, &renderer.vboHandle)
	gl.GenBuffers(1, &renderer.elementsHandle)
	
	renderer.createFontTexture()
}

func (renderer *OpenGL3) createShader(shaderType uint32, source string) uint32 {
	shader := gl.CreateShader(shaderType)
	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)
	return shader
}

func (renderer *OpenGL3) createFontTexture() {
	io := imgui.CurrentIO()
	image := io.Fonts().TextureDataRGBA32()
	
	gl.GenTextures(1, &renderer.fontTexture)
	gl.BindTexture(gl.TEXTURE_2D, renderer.fontTexture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.PixelStorei(gl.UNPACK_ROW_LENGTH, 0)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(image.Width), int32(image.Height), 0, gl.RGBA, gl.UNSIGNED_BYTE, image.Pixels)
	
	io.Fonts().SetTextureID(imgui.TextureID(renderer.fontTexture))
}

func (renderer *OpenGL3) Render(drawData imgui.DrawData) {
	// Backup GL state
	var lastActiveTexture int32
	gl.GetIntegerv(gl.ACTIVE_TEXTURE, &lastActiveTexture)
	gl.ActiveTexture(gl.TEXTURE0)
	
	var lastProgram int32
	gl.GetIntegerv(gl.CURRENT_PROGRAM, &lastProgram)
	
	var lastTexture int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &lastTexture)
	
	var lastArrayBuffer int32
	gl.GetIntegerv(gl.ARRAY_BUFFER_BINDING, &lastArrayBuffer)
	
	var lastElementArrayBuffer int32
	gl.GetIntegerv(gl.ELEMENT_ARRAY_BUFFER_BINDING, &lastElementArrayBuffer)
	
	var lastViewport [4]int32
	gl.GetIntegerv(gl.VIEWPORT, &lastViewport[0])
	
	var lastScissorBox [4]int32
	gl.GetIntegerv(gl.SCISSOR_BOX, &lastScissorBox[0])
	
	var lastBlendSrcRgb int32
	gl.GetIntegerv(gl.BLEND_SRC_RGB, &lastBlendSrcRgb)
	
	var lastBlendDstRgb int32
	gl.GetIntegerv(gl.BLEND_DST_RGB, &lastBlendDstRgb)
	
	var lastBlendSrcAlpha int32
	gl.GetIntegerv(gl.BLEND_SRC_ALPHA, &lastBlendSrcAlpha)
	
	var lastBlendDstAlpha int32
	gl.GetIntegerv(gl.BLEND_DST_ALPHA, &lastBlendDstAlpha)
	
	var lastBlendEquationRgb int32
	gl.GetIntegerv(gl.BLEND_EQUATION_RGB, &lastBlendEquationRgb)
	
	var lastBlendEquationAlpha int32
	gl.GetIntegerv(gl.BLEND_EQUATION_ALPHA, &lastBlendEquationAlpha)
	
	lastEnableBlend := gl.IsEnabled(gl.BLEND)
	lastEnableCullFace := gl.IsEnabled(gl.CULL_FACE)
	lastEnableDepthTest := gl.IsEnabled(gl.DEPTH_TEST)
	lastEnableScissorTest := gl.IsEnabled(gl.SCISSOR_TEST)
	
	// Setup render state
	gl.Enable(gl.BLEND)
	gl.BlendEquation(gl.FUNC_ADD)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.CULL_FACE)
	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.SCISSOR_TEST)
	
	// Setup viewport and projection matrix
	displaySize := drawData.DisplaySize()
	gl.Viewport(0, 0, int32(displaySize.X), int32(displaySize.Y))
	orthoProjection := [4][4]float32{
		{2.0 / displaySize.X, 0.0, 0.0, 0.0},
		{0.0, 2.0 / -displaySize.Y, 0.0, 0.0},
		{0.0, 0.0, -1.0, 0.0},
		{-1.0, 1.0, 0.0, 1.0},
	}
	
	gl.UseProgram(renderer.shaderHandle)
	gl.Uniform1i(renderer.attribLocationTex, 0)
	gl.UniformMatrix4fv(renderer.attribLocationProjMtx, 1, false, &orthoProjection[0][0])
	
	gl.BindBuffer(gl.ARRAY_BUFFER, renderer.vboHandle)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, renderer.elementsHandle)
	
	gl.EnableVertexAttribArray(uint32(renderer.attribLocationVtxPos))
	gl.EnableVertexAttribArray(uint32(renderer.attribLocationVtxUV))
	gl.EnableVertexAttribArray(uint32(renderer.attribLocationVtxColor))
	
	vertexSize, vertexOffsetPos, vertexOffsetUv, vertexOffsetCol := imgui.VertexBufferLayout()
	gl.VertexAttribPointer(uint32(renderer.attribLocationVtxPos), 2, gl.FLOAT, false, int32(vertexSize), gl.PtrOffset(int(vertexOffsetPos)))
	gl.VertexAttribPointer(uint32(renderer.attribLocationVtxUV), 2, gl.FLOAT, false, int32(vertexSize), gl.PtrOffset(int(vertexOffsetUv)))
	gl.VertexAttribPointer(uint32(renderer.attribLocationVtxColor), 4, gl.UNSIGNED_BYTE, true, int32(vertexSize), gl.PtrOffset(int(vertexOffsetCol)))
	
	// Render command lists
	for _, list := range drawData.CommandLists() {
		gl.BufferData(gl.ARRAY_BUFFER, list.VertexBufferSize(), list.VertexBufferData(), gl.STREAM_DRAW)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, list.IndexBufferSize(), list.IndexBufferData(), gl.STREAM_DRAW)
		
		for _, cmd := range list.Commands() {
			if cmd.HasUserCallback() {
				cmd.CallUserCallback(list)
			} else {
				gl.BindTexture(gl.TEXTURE_2D, uint32(cmd.TextureID()))
				clipRect := cmd.ClipRect()
				gl.Scissor(int32(clipRect.X), int32(displaySize.Y-clipRect.W), int32(clipRect.Z-clipRect.X), int32(clipRect.W-clipRect.Y))
				gl.DrawElementsBaseVertex(gl.TRIANGLES, int32(cmd.ElementCount()), gl.UNSIGNED_SHORT, gl.PtrOffset(int(cmd.IndexOffset()*2)), int32(cmd.VertexOffset()))
			}
		}
	}
	
	// Restore GL state
	gl.UseProgram(uint32(lastProgram))
	gl.BindTexture(gl.TEXTURE_2D, uint32(lastTexture))
	gl.ActiveTexture(uint32(lastActiveTexture))
	gl.BindBuffer(gl.ARRAY_BUFFER, uint32(lastArrayBuffer))
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, uint32(lastElementArrayBuffer))
	gl.BlendEquationSeparate(uint32(lastBlendEquationRgb), uint32(lastBlendEquationAlpha))
	gl.BlendFuncSeparate(uint32(lastBlendSrcRgb), uint32(lastBlendDstRgb), uint32(lastBlendSrcAlpha), uint32(lastBlendDstAlpha))
	
	if lastEnableBlend {
		gl.Enable(gl.BLEND)
	} else {
		gl.Disable(gl.BLEND)
	}
	
	if lastEnableCullFace {
		gl.Enable(gl.CULL_FACE)
	} else {
		gl.Disable(gl.CULL_FACE)
	}
	
	if lastEnableDepthTest {
		gl.Enable(gl.DEPTH_TEST)
	} else {
		gl.Disable(gl.DEPTH_TEST)
	}
	
	if lastEnableScissorTest {
		gl.Enable(gl.SCISSOR_TEST)
	} else {
		gl.Disable(gl.SCISSOR_TEST)
	}
	
	gl.Viewport(lastViewport[0], lastViewport[1], lastViewport[2], lastViewport[3])
	gl.Scissor(lastScissorBox[0], lastScissorBox[1], lastScissorBox[2], lastScissorBox[3])
}

func (renderer *OpenGL3) Destroy() {
	if renderer.vboHandle != 0 {
		gl.DeleteBuffers(1, &renderer.vboHandle)
	}
	if renderer.elementsHandle != 0 {
		gl.DeleteBuffers(1, &renderer.elementsHandle)
	}
	if renderer.shaderHandle != 0 && renderer.vertHandle != 0 {
		gl.DetachShader(renderer.shaderHandle, renderer.vertHandle)
	}
	if renderer.shaderHandle != 0 && renderer.fragHandle != 0 {
		gl.DetachShader(renderer.shaderHandle, renderer.fragHandle)
	}
	if renderer.vertHandle != 0 {
		gl.DeleteShader(renderer.vertHandle)
	}
	if renderer.fragHandle != 0 {
		gl.DeleteShader(renderer.fragHandle)
	}
	if renderer.shaderHandle != 0 {
		gl.DeleteProgram(renderer.shaderHandle)
	}
	if renderer.fontTexture != 0 {
		gl.DeleteTextures(1, &renderer.fontTexture)
	}
}