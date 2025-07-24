package ui

import (
	"github.com/inkyblackness/imgui-go/v4"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type SimpleImGuiContext struct {
	context   *imgui.Context
	platform  *SimplePlatform
	renderer  *SimpleRenderer
}

type SimplePlatform struct {
	window *glfw.Window
	time   float64
}

type SimpleRenderer struct {
	program uint32
	texture uint32
	vao     uint32
	vbo     uint32
	ebo     uint32
}

func NewSimpleImGuiContext(window *glfw.Window) *SimpleImGuiContext {
	ctx := imgui.CreateContext(nil)
	imgui.StyleColorsDark()
	
	platform := &SimplePlatform{window: window}
	renderer := &SimpleRenderer{}
	
	// Setup key mapping
	io := imgui.CurrentIO()
	io.KeyMap(imgui.KeyTab, int(glfw.KeyTab))
	io.KeyMap(imgui.KeyLeftArrow, int(glfw.KeyLeft))
	io.KeyMap(imgui.KeyRightArrow, int(glfw.KeyRight))
	io.KeyMap(imgui.KeyUpArrow, int(glfw.KeyUp))
	io.KeyMap(imgui.KeyDownArrow, int(glfw.KeyDown))
	io.KeyMap(imgui.KeyEnter, int(glfw.KeyEnter))
	io.KeyMap(imgui.KeyEscape, int(glfw.KeyEscape))
	io.KeyMap(imgui.KeyBackspace, int(glfw.KeyBackspace))
	io.KeyMap(imgui.KeyDelete, int(glfw.KeyDelete))
	io.KeyMap(imgui.KeySpace, int(glfw.KeySpace))
	
	renderer.createDeviceObjects()
	
	return &SimpleImGuiContext{
		context:  ctx,
		platform: platform,
		renderer: renderer,
	}
}

func (ctx *SimpleImGuiContext) NewFrame() {
	// Update display size
	w, h := ctx.platform.window.GetSize()
	imgui.CurrentIO().SetDisplaySize(imgui.Vec2{X: float32(w), Y: float32(h)})
	
	// Update time
	currentTime := glfw.GetTime()
	if ctx.platform.time > 0 {
		imgui.CurrentIO().SetDeltaTime(float32(currentTime - ctx.platform.time))
	}
	ctx.platform.time = currentTime
	
	imgui.NewFrame()
}

func (ctx *SimpleImGuiContext) Render() {
	imgui.Render()
	ctx.renderer.renderDrawData(imgui.RenderedDrawData())
}

func (ctx *SimpleImGuiContext) Destroy() {
	ctx.renderer.destroy()
	ctx.context.Destroy()
}

func (r *SimpleRenderer) createDeviceObjects() {
	// Create a minimal shader for ImGui
	vertexShader := `
#version 330 core
layout (location = 0) in vec2 Position;
layout (location = 1) in vec2 UV;
layout (location = 2) in vec4 Color;

uniform mat4 ProjMtx;
out vec2 Frag_UV;
out vec4 Frag_Color;

void main() {
	Frag_UV = UV;
	Frag_Color = Color;
	gl_Position = ProjMtx * vec4(Position.xy, 0, 1);
}
`
	
	fragmentShader := `
#version 330 core
in vec2 Frag_UV;
in vec4 Frag_Color;

uniform sampler2D Texture;
out vec4 Out_Color;

void main() {
	Out_Color = Frag_Color * texture(Texture, Frag_UV.st);
}
`
	
	// Compile shaders
	vs := gl.CreateShader(gl.VERTEX_SHADER)
	csource, free := gl.Strs(vertexShader + "\x00")
	gl.ShaderSource(vs, 1, csource, nil)
	free()
	gl.CompileShader(vs)
	
	fs := gl.CreateShader(gl.FRAGMENT_SHADER)
	csource, free = gl.Strs(fragmentShader + "\x00")
	gl.ShaderSource(fs, 1, csource, nil)
	free()
	gl.CompileShader(fs)
	
	r.program = gl.CreateProgram()
	gl.AttachShader(r.program, vs)
	gl.AttachShader(r.program, fs)
	gl.LinkProgram(r.program)
	gl.DeleteShader(vs)
	gl.DeleteShader(fs)
	
	// Create buffers
	gl.GenBuffers(1, &r.vbo)
	gl.GenBuffers(1, &r.ebo)
	gl.GenVertexArrays(1, &r.vao)
	
	// Create font texture
	io := imgui.CurrentIO()
	image := io.Fonts().TextureDataRGBA32()
	
	gl.GenTextures(1, &r.texture)
	gl.BindTexture(gl.TEXTURE_2D, r.texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(image.Width), int32(image.Height), 0, gl.RGBA, gl.UNSIGNED_BYTE, image.Pixels)
	
	io.Fonts().SetTextureID(imgui.TextureID(r.texture))
}

func (r *SimpleRenderer) renderDrawData(drawData imgui.DrawData) {
	// This is a minimal implementation
	// For now, we'll just enable blending for ImGui
	gl.Enable(gl.BLEND)
	gl.BlendEquation(gl.FUNC_ADD)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.CULL_FACE)
	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.SCISSOR_TEST)
	
	// Restore state after
	defer func() {
		gl.Disable(gl.BLEND)
		gl.Enable(gl.CULL_FACE)
		gl.Enable(gl.DEPTH_TEST)
		gl.Disable(gl.SCISSOR_TEST)
	}()
}

func (r *SimpleRenderer) destroy() {
	if r.vbo != 0 {
		gl.DeleteBuffers(1, &r.vbo)
	}
	if r.ebo != 0 {
		gl.DeleteBuffers(1, &r.ebo)
	}
	if r.vao != 0 {
		gl.DeleteVertexArrays(1, &r.vao)
	}
	if r.texture != 0 {
		gl.DeleteTextures(1, &r.texture)
	}
	if r.program != 0 {
		gl.DeleteProgram(r.program)
	}
}