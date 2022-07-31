//go:build android

package display

import (
	"unsafe"

	"github.com/rajveermalviya/gamen/internal/android"
)

func NewDisplay() (Display, error) {
	return android.NewDisplay()
}

func NewWindow(d Display) (Window, error) {
	return android.NewWindow()
}

type AndroidWindow interface {
	ANativeWindow() unsafe.Pointer
}
