//go:build android

package display

import (
	"unsafe"

	"github.com/rajveermalviya/gamen/internal/android"
)

// NewDisplay initializes the event loop and returns
// a handle to manage it.
//
// Must only be called from main goroutine.
func NewDisplay() (Display, error) {
	return android.NewDisplay()
}

// NewWindow creates a new window for the provided
// display event loop.
//
// To receive events you must set individual callbacks
// via Set[event]Callback methods.
//
// Must only be called from main goroutine.
func NewWindow(d Display) (Window, error) {
	return android.NewWindow()
}

type AndroidWindow interface {
	ANativeWindow() unsafe.Pointer
}
