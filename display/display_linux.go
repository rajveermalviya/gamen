//go:build linux && !android

package display

import (
	"os"
	"unsafe"

	"github.com/rajveermalviya/gamen/internal/wayland"
	"github.com/rajveermalviya/gamen/internal/xcb"
)

var backend = os.Getenv("GAMEN_DISPLAY_BACKEND")

func NewDisplay() (Display, error) {
	switch backend {
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
