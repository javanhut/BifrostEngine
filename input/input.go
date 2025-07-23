package input

import (
	"sync"

	"github.com/go-gl/glfw/v3.3/glfw"
)

// KeyState represents the state of a key
type KeyState int

const (
	Released KeyState = iota
	Pressed
	Held
	JustReleased
)

// MouseButton represents mouse buttons
type MouseButton int

const (
	MouseButtonLeft MouseButton = iota
	MouseButtonRight
	MouseButtonMiddle
)

// InputManager handles all input for the engine
type InputManager struct {
	window         *glfw.Window
	keys           map[glfw.Key]KeyState
	prevKeys       map[glfw.Key]KeyState
	mouseButtons   map[MouseButton]KeyState
	prevMouseButtons map[MouseButton]KeyState
	mouseX         float64
	mouseY         float64
	mouseDeltaX    float64
	mouseDeltaY    float64
	scrollX        float64
	scrollY        float64
	callbacks      *InputCallbacks
	mu             sync.RWMutex
}

// InputCallbacks holds user-defined callbacks
type InputCallbacks struct {
	OnKeyPress    func(key glfw.Key)
	OnKeyRelease  func(key glfw.Key)
	OnMousePress  func(button MouseButton)
	OnMouseRelease func(button MouseButton)
	OnMouseMove   func(x, y float64)
	OnScroll      func(xOffset, yOffset float64)
}

// NewInputManager creates a new input manager
func NewInputManager(window *glfw.Window) *InputManager {
	im := &InputManager{
		window:           window,
		keys:             make(map[glfw.Key]KeyState),
		prevKeys:         make(map[glfw.Key]KeyState),
		mouseButtons:     make(map[MouseButton]KeyState),
		prevMouseButtons: make(map[MouseButton]KeyState),
		callbacks:        &InputCallbacks{},
	}
	
	im.setupCallbacks()
	return im
}

// setupCallbacks sets up GLFW callbacks
func (im *InputManager) setupCallbacks() {
	im.window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		im.mu.Lock()
		switch action {
		case glfw.Press:
			im.keys[key] = Pressed
			if im.callbacks.OnKeyPress != nil {
				im.callbacks.OnKeyPress(key)
			}
		case glfw.Release:
			im.keys[key] = JustReleased
			if im.callbacks.OnKeyRelease != nil {
				im.callbacks.OnKeyRelease(key)
			}
		case glfw.Repeat:
			im.keys[key] = Held
		}
		im.mu.Unlock()
	})
	
	im.window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		im.mu.Lock()
		mb := MouseButton(button)
		switch action {
		case glfw.Press:
			im.mouseButtons[mb] = Pressed
			if im.callbacks.OnMousePress != nil {
				im.callbacks.OnMousePress(mb)
			}
		case glfw.Release:
			im.mouseButtons[mb] = JustReleased
			if im.callbacks.OnMouseRelease != nil {
				im.callbacks.OnMouseRelease(mb)
			}
		}
		im.mu.Unlock()
	})
	
	im.window.SetCursorPosCallback(func(w *glfw.Window, xpos, ypos float64) {
		im.mu.Lock()
		im.mouseDeltaX = xpos - im.mouseX
		im.mouseDeltaY = ypos - im.mouseY
		im.mouseX = xpos
		im.mouseY = ypos
		im.mu.Unlock()
		
		if im.callbacks.OnMouseMove != nil {
			im.callbacks.OnMouseMove(xpos, ypos)
		}
	})
	
	im.window.SetScrollCallback(func(w *glfw.Window, xOffset, yOffset float64) {
		im.mu.Lock()
		im.scrollX = xOffset
		im.scrollY = yOffset
		im.mu.Unlock()
		
		if im.callbacks.OnScroll != nil {
			im.callbacks.OnScroll(xOffset, yOffset)
		}
	})
}

// Update updates input states - should be called once per frame
func (im *InputManager) Update() {
	im.mu.Lock()
	defer im.mu.Unlock()
	
	// Copy current states to previous
	for k, v := range im.keys {
		im.prevKeys[k] = v
		if v == Pressed {
			im.keys[k] = Held
		} else if v == JustReleased {
			im.keys[k] = Released
		}
	}
	
	for b, v := range im.mouseButtons {
		im.prevMouseButtons[b] = v
		if v == Pressed {
			im.mouseButtons[b] = Held
		} else if v == JustReleased {
			im.mouseButtons[b] = Released
		}
	}
	
	// Reset deltas
	im.mouseDeltaX = 0
	im.mouseDeltaY = 0
	im.scrollX = 0
	im.scrollY = 0
}

// IsKeyPressed returns true if key was just pressed this frame
func (im *InputManager) IsKeyPressed(key glfw.Key) bool {
	im.mu.RLock()
	defer im.mu.RUnlock()
	return im.keys[key] == Pressed
}

// IsKeyHeld returns true if key is being held down
func (im *InputManager) IsKeyHeld(key glfw.Key) bool {
	im.mu.RLock()
	defer im.mu.RUnlock()
	return im.keys[key] == Held || im.keys[key] == Pressed
}

// IsKeyReleased returns true if key was just released this frame
func (im *InputManager) IsKeyReleased(key glfw.Key) bool {
	im.mu.RLock()
	defer im.mu.RUnlock()
	return im.keys[key] == JustReleased
}

// IsMouseButtonPressed returns true if mouse button was just pressed
func (im *InputManager) IsMouseButtonPressed(button MouseButton) bool {
	im.mu.RLock()
	defer im.mu.RUnlock()
	return im.mouseButtons[button] == Pressed
}

// IsMouseButtonHeld returns true if mouse button is being held
func (im *InputManager) IsMouseButtonHeld(button MouseButton) bool {
	im.mu.RLock()
	defer im.mu.RUnlock()
	return im.mouseButtons[button] == Held || im.mouseButtons[button] == Pressed
}

// IsMouseButtonReleased returns true if mouse button was just released
func (im *InputManager) IsMouseButtonReleased(button MouseButton) bool {
	im.mu.RLock()
	defer im.mu.RUnlock()
	return im.mouseButtons[button] == JustReleased
}

// GetMousePosition returns current mouse position
func (im *InputManager) GetMousePosition() (float64, float64) {
	im.mu.RLock()
	defer im.mu.RUnlock()
	return im.mouseX, im.mouseY
}

// GetMouseDelta returns mouse movement since last frame
func (im *InputManager) GetMouseDelta() (float64, float64) {
	im.mu.RLock()
	defer im.mu.RUnlock()
	return im.mouseDeltaX, im.mouseDeltaY
}

// GetScrollDelta returns scroll wheel delta
func (im *InputManager) GetScrollDelta() (float64, float64) {
	im.mu.RLock()
	defer im.mu.RUnlock()
	return im.scrollX, im.scrollY
}

// SetCallbacks sets the callback functions
func (im *InputManager) SetCallbacks(callbacks *InputCallbacks) {
	im.mu.Lock()
	defer im.mu.Unlock()
	im.callbacks = callbacks
}