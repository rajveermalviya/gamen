//go:build windows

package display

import "github.com/rajveermalviya/gamen/internal/win32"

// NewDisplay initializes the event loop and returns
// a handle to manage it.
//
// Must only be called from main goroutine.
func NewDisplay() (Display, error) {
	return win32.NewDisplay()
}

// NewWindow creates a new window for the provided
// display event loop.
//
// To receive events you must set individual callbacks
// via Set[event]Callback methods.
//
// Must only be called from main goroutine.
func NewWindow(d Display) (Window, error) {
	return win32.NewWindow(d.(*win32.Display))
}

type Win32Window interface {
	Win32Hinstance() uintptr
	Win32Hwnd() uintptr
}
