//go:build windows

package display

import "github.com/rajveermalviya/gamen/internal/win32"

func NewDisplay() (Display, error) {
	return win32.NewDisplay()
}

func NewWindow(d Display) (Window, error) {
	return win32.NewWindow(d.(*win32.Display))
}

type Win32Window interface {
	Win32Hinstance() uintptr
	Win32Hwnd() uintptr
}
