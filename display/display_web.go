//go:build js

package display

import (
	"syscall/js"

	"github.com/rajveermalviya/gamen/internal/web"
)

func NewDisplay() (Display, error) {
	return web.NewDisplay()
}

func NewWindow(d Display) (Window, error) {
	return web.NewWindow(d.(*web.Display))
}

type WebWindow interface {
	WebCanvas() js.Value
}
