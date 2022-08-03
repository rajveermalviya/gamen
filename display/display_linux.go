//go:build linux && !android

package display

import (
	"os"
	"unsafe"

	"github.com/rajveermalviya/gamen/internal/wayland"
	"github.com/rajveermalviya/gamen/internal/xcb"
)

// NewDisplay initializes the event loop and returns
// a handle to manage it.
//
// Must only be called from main goroutine.
func NewDisplay() (Display, error) {
	switch os.Getenv("GAMEN_DISPLAY_BACKEND") {
	case "wayland":
		return wayland.NewDisplay()

	case "xcb":
		return xcb.NewDisplay()

	case "":
		d, err := wayland.NewDisplay()
		if err == nil {
			return d, nil
		}

		return xcb.NewDisplay()

	default:
		panic("invalid GAMEN_DISPLAY_BACKEND")
	}
}

// NewWindow creates a new window for the provided
// display event loop.
//
// To receive events you must set individual callbacks
// via Set[event]Callback methods.
//
// Must only be called from main goroutine.
func NewWindow(d Display) (Window, error) {
	switch d := d.(type) {
	case *wayland.Display:
		return wayland.NewWindow(d)
	case *xcb.Display:
		return xcb.NewWindow(d)

	default:
		panic("invalid Display")
	}
}

type WaylandWindow interface {
	WlDisplay() unsafe.Pointer
	WlSurface() unsafe.Pointer
}

type XcbWindow interface {
	XcbConnection() unsafe.Pointer
	XcbWindow() uint32
}

type XlibWindow interface {
	XlibDisplay() unsafe.Pointer
	XlibWindow() uint32
}
