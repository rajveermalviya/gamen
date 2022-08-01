package display

import (
	"time"

	"github.com/rajveermalviya/gamen/cursors"
	"github.com/rajveermalviya/gamen/dpi"
	"github.com/rajveermalviya/gamen/events"
)

type Display interface {
	Destroy()
	Wait() bool
	Poll() bool
	WaitTimeout(time.Duration) bool
}

type Window interface {
	Destroy()

	SetTitle(string)

	InnerSize() dpi.PhysicalSize[uint32]

	SetInnerSize(dpi.Size[uint32])
	SetMinInnerSize(dpi.Size[uint32])
	SetMaxInnerSize(dpi.Size[uint32])

	Maximized() bool

	SetMinimized()
	SetMaximized(maximized bool)

	SetCursorIcon(cursors.Icon)
	SetCursorVisible(visible bool)

	SetFullscreen(fullscreen bool)
	Fullscreen() bool

	SetResizedCallback(events.WindowResizedCallback)
	SetCloseRequestedCallback(events.WindowCloseRequestedCallback)

	SetCursorEnteredCallback(events.WindowCursorEnteredCallback)
	SetCursorLeftCallback(events.WindowCursorLeftCallback)
	SetCursorMovedCallback(events.WindowCursorMovedCallback)

	SetMouseWheelCallback(events.WindowMouseWheelCallback)
	SetMouseInputCallback(events.WindowMouseInputCallback)

	SetTouchInputCallback(events.WindowTouchInputCallback)

	SetFocusedCallback(events.WindowFocusedCallback)
	SetUnfocusedCallback(events.WindowUnfocusedCallback)

	SetModifiersChangedCallback(events.WindowModifiersChangedCallback)
	SetKeyboardInputCallback(events.WindowKeyboardInputCallback)
	SetReceivedCharacterCallback(events.WindowReceivedCharacterCallback)
}

type AndroidWindowExt interface {
	SetSurfaceCreatedCallback(events.WindowSurfaceCreatedCallback)
	SetSurfaceDestroyedCallback(events.WindowSurfaceDestroyedCallback)

	EnableIme()
	DisableIme()
}
