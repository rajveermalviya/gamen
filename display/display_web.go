//go:build js

package display

import (
	"syscall/js"

	"github.com/rajveermalviya/gamen/internal/web"
)

// NewDisplay initializes the event loop and returns
// a handle to manage it.
//
// Must only be called from main goroutine.
func NewDisplay() (Display, error) {
	return web.NewDisplay()
}

// NewWindow creates a new window for the provided
// display event loop.
//
// To receive events you must set individual callbacks
// via Set[event]Callback methods.
//
// Must only be called from main goroutine.
func NewWindow(d Display) (Window, error) {
	return web.NewWindow(d.(*web.Display))
}

type WebWindow interface {
	WebCanvas() js.Value
}
