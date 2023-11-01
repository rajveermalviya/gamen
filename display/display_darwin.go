//go:build darwin

package display

import (
	glfwraw "github.com/go-gl/glfw/v3.3/glfw"
	"github.com/rajveermalviya/gamen/internal/glfw"
)

// NewDisplay initializes the event loop and returns
// a handle to manage it.
//
// Must only be called from main goroutine.
func NewDisplay() (Display, error) {
	return glfw.NewDisplay()
}

// NewWindow creates a new window for the provided
// display event loop.
//
// To receive events you must set individual callbacks
// via Set[event]Callback methods.
//
// Must only be called from main goroutine.
func NewWindow(d Display) (Window, error) {
	return d.(*glfw.Display).NewWindow()
}

type GlfwWindow interface {
	GlfwWindow() *glfwraw.Window
}
