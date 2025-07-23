package window

import (
	"fmt"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
)

type Window struct {
	handle *glfw.Window
	width  int
	height int
	title  string
}

func init() {
	runtime.LockOSThread()
}

func New(width, height int, title string) (*Window, error) {
	if err := glfw.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize GLFW: %w", err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	handle, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create window: %w", err)
	}

	handle.MakeContextCurrent()

	return &Window{
		handle: handle,
		width:  width,
		height: height,
		title:  title,
	}, nil
}

func (w *Window) ShouldClose() bool {
	return w.handle.ShouldClose()
}

func (w *Window) SwapBuffers() {
	w.handle.SwapBuffers()
}

func (w *Window) PollEvents() {
	glfw.PollEvents()
}

func (w *Window) GetSize() (int, int) {
	return w.handle.GetSize()
}

func (w *Window) Destroy() {
	w.handle.Destroy()
	glfw.Terminate()
}

func (w *Window) SetKeyCallback(fn glfw.KeyCallback) glfw.KeyCallback {
	return w.handle.SetKeyCallback(fn)
}

func (w *Window) GetHandle() *glfw.Window {
	return w.handle
}

func (w *Window) SetMouseButtonCallback(fn glfw.MouseButtonCallback) glfw.MouseButtonCallback {
	return w.handle.SetMouseButtonCallback(fn)
}

func (w *Window) SetCursorPosCallback(fn glfw.CursorPosCallback) glfw.CursorPosCallback {
	return w.handle.SetCursorPosCallback(fn)
}

func (w *Window) SetScrollCallback(fn glfw.ScrollCallback) glfw.ScrollCallback {
	return w.handle.SetScrollCallback(fn)
}