//go:build android

package android

import (
	"unsafe"

	"github.com/rajveermalviya/gamen/cursors"
	"github.com/rajveermalviya/gamen/dpi"
	"github.com/rajveermalviya/gamen/events"
)

/*

#include <game-activity/native_app_glue/android_native_app_glue.h>

*/
import "C"

type Window struct{}

func NewWindow() (*Window, error) { return &Window{}, nil }

func (*Window) ANativeWindow() unsafe.Pointer {
	if app := androidApp.Load(); app != nil {
		return unsafe.Pointer(app.window)
	}
	return nil
}

func (*Window) InnerSize() dpi.PhysicalSize[uint32] {
	if app := androidApp.Load(); app != nil && app.window != nil {
		return dpi.PhysicalSize[uint32]{
			Width:  uint32(C.ANativeWindow_getWidth(app.window)),
			Height: uint32(C.ANativeWindow_getHeight(app.window)),
		}
	}

	return dpi.PhysicalSize[uint32]{}
}

func (*Window) EnableIme() {
	if app := androidApp.Load(); app != nil && app.activity != nil {
		C.GameActivity_showSoftInput(app.activity, 0)
	}
}

func (*Window) DisableIme() {
	if app := androidApp.Load(); app != nil && app.activity != nil {
		C.GameActivity_hideSoftInput(app.activity, 0)
	}
}

func (*Window) SetTitle(string)                  {}
func (*Window) Destroy()                         {}
func (*Window) SetInnerSize(dpi.Size[uint32])    {}
func (*Window) SetMinInnerSize(dpi.Size[uint32]) {}
func (*Window) SetMaxInnerSize(dpi.Size[uint32]) {}
func (*Window) Maximized() bool                  { return false }
func (*Window) SetMinimized()                    {}
func (*Window) SetMaximized(bool)                {}
func (*Window) SetCursorIcon(cursors.Icon)       {}
func (*Window) SetCursorVisible(bool)            {}
func (*Window) SetFullscreen(bool)               {}
func (*Window) Fullscreen() bool                 { return false }
func (*Window) DragWindow()                      {}
func (*Window) SetDecorations(bool)              {}
func (*Window) Decorated() bool                  { return false }

func (w *Window) SetSurfaceCreatedCallback(cb events.WindowSurfaceCreatedCallback) {
	windowSurfaceCreatedCb.Store(&cb)
}
func (w *Window) SetSurfaceDestroyedCallback(cb events.WindowSurfaceDestroyedCallback) {
	windowSurfaceDestroyedCb.Store(&cb)
}
func (w *Window) SetResizedCallback(cb events.WindowResizedCallback) {
	windowResizedCallback.Store(&cb)
}
func (w *Window) SetFocusedCallback(cb events.WindowFocusedCallback) {
	windowFocusedCb.Store(&cb)
}
func (w *Window) SetUnfocusedCallback(cb events.WindowUnfocusedCallback) {
	windowUnfocusedCb.Store(&cb)
}
func (w *Window) SetTouchInputCallback(cb events.WindowTouchInputCallback) {
	windowTouchInputCb.Store(&cb)
}
func (w *Window) SetKeyboardInputCallback(cb events.WindowKeyboardInputCallback) {
	windowKeyboardInputCb.Store(&cb)
}
func (w *Window) SetReceivedCharacterCallback(cb events.WindowReceivedCharacterCallback) {
	windowReceivedCharacterCallback.Store(&cb)
}

func (w *Window) SetCloseRequestedCallback(cb events.WindowCloseRequestedCallback)     {}
func (w *Window) SetCursorEnteredCallback(cb events.WindowCursorEnteredCallback)       {}
func (w *Window) SetCursorLeftCallback(cb events.WindowCursorLeftCallback)             {}
func (w *Window) SetCursorMovedCallback(cb events.WindowCursorMovedCallback)           {}
func (w *Window) SetMouseScrollCallback(cb events.WindowMouseScrollCallback)           {}
func (w *Window) SetMouseInputCallback(cb events.WindowMouseInputCallback)             {}
func (w *Window) SetModifiersChangedCallback(cb events.WindowModifiersChangedCallback) {}
