package display

import (
	"time"

	"github.com/rajveermalviya/gamen/cursors"
	"github.com/rajveermalviya/gamen/dpi"
	"github.com/rajveermalviya/gamen/events"
)

// Display is the main event loop handle.
// It provides a way to poll window events from the system
type Display interface {
	// Destroy destroys all the windows associated with this display
	// and the underlying system event loop
	//
	// Can be called multiple times, but subsequent calls are ignored.
	// Can be called from any goroutine.
	Destroy()

	// Wait blocks until there is a new event, it returns after dispatching
	// all the pending events in the queue.
	//
	// Return value of false indicates Display is destroyed and you must
	// exit the loop.
	//
	// Can only be called from main goroutine
	Wait() bool

	// Poll dispatches all the pending events in the queue and returns immediately
	//
	// Return value of false indicates Display is destroyed and you must
	// exit the loop.
	//
	// Can only be called from main goroutine
	Poll() bool

	// WaitTimeout blocks until there is a new event or until the specified timeout
	// whichever happens first
	//
	// Return value of false indicates Display is destroyed and you must
	// exit the loop.
	//
	// Can only be called from main goroutine
	WaitTimeout(time.Duration) bool
}

// Window is a window
//
// Its methods can be called from any goroutine
type Window interface {
	// Destroys the window, it must not be used after this point.
	//
	// Can be called multiple times, but subsequent calls are ignored.
	Destroy()

	// SetTitle sets the title for the window
	//
	// Unsupported backends: android
	SetTitle(string)

	// InnerSize returns the current size of the drawable surface
	//
	InnerSize() dpi.PhysicalSize[uint32]

	// SetInnerSize sets the size of the drawable surface
	//
	// Unsupported backends: android
	SetInnerSize(dpi.Size[uint32])

	// SetMinInnerSize sets the minimum size window can be resized
	//
	// Unsupported backends: android
	SetMinInnerSize(dpi.Size[uint32])

	// SetMaxInnerSize sets the maximum size window can be resized
	//
	// Unsupported backends: android
	SetMaxInnerSize(dpi.Size[uint32])

	// Maximized returns if window is currently maximized
	//
	// Unsupported backends: android, web
	Maximized() bool

	// SetMinimized minimizes or iconifies the window
	//
	// Unsupported backends: android, web
	SetMinimized()

	// SetMinimized maximizes or un-maximizes the window
	//
	// Unsupported backends: android, web
	SetMaximized(bool)

	// SetCursorIcon changes the cursor icon for this window
	//
	// Unsupported backends: android
	SetCursorIcon(cursors.Icon)

	// SetCursorVisible hides or un-hides the cursor for this window
	//
	// Unsupported backends: android
	SetCursorVisible(bool)

	// SetFullscreen sets window to fullscreen or normal mode
	//
	// Unsupported backends: android
	SetFullscreen(bool)

	// Fullscreen returns if window is currently in fullscreen mode
	//
	// Unsupported backends: android
	Fullscreen() bool

	// DragWindow starts an interactive move of window.
	// Window follows the mouse cursor until mouse button is released.
	//
	// Unsupported backends: android, web
	DragWindow()

	// SetDecorations enables or disables window decorations
	//
	// Unsupported backends: android, web
	SetDecorations(bool)

	// Decorated returns if window currently has decorations
	//
	// Unsupported backends: android, web
	Decorated() bool

	// Callbacks

	// SetResizedCallback registers a callback to receive resize events for the window
	// It also fires when scale factor for the window changes
	//
	SetResizedCallback(events.WindowResizedCallback)

	// SetCloseRequestedCallback registers a callback to receive close event for the window
	//
	// Unsupported backends: android
	SetCloseRequestedCallback(events.WindowCloseRequestedCallback)

	// SetCursorEnteredCallback registers a callback to receive an event when cursor enters the window
	//
	// Unsupported backends: android
	SetCursorEnteredCallback(events.WindowCursorEnteredCallback)

	// SetCursorLeftCallback registers a callback to receive an event when cursor leaves the window
	//
	// Unsupported backends: android
	SetCursorLeftCallback(events.WindowCursorLeftCallback)

	// SetCursorLeftCallback registers a callback to receive an event when cursor moves inside the window
	//
	// Unsupported backends: android
	SetCursorMovedCallback(events.WindowCursorMovedCallback)

	// SetMouseScrollCallback registers a callback to receive mouse scroll event for the window
	//
	// Unsupported backends: android
	SetMouseScrollCallback(events.WindowMouseScrollCallback)

	// SetMouseInputCallback registers a callback to receive mouse input event for the window
	//
	// Unsupported backends: android
	SetMouseInputCallback(events.WindowMouseInputCallback)

	// SetMouseInputCallback registers a callback to receive mouse input event for the window
	//
	// Unsupported backends: wayland, web, win32, xcb
	SetTouchInputCallback(events.WindowTouchInputCallback)

	// SetFocusedCallback registers a callback to receive focus event for the window
	//
	SetFocusedCallback(events.WindowFocusedCallback)

	// SetUnfocusedCallback registers a callback to receive un-focus event for the window
	//
	SetUnfocusedCallback(events.WindowUnfocusedCallback)

	// SetModifiersChangedCallback registers a callback to receive an event when modifiers change for the window
	//
	// Unsupported backends: android
	SetModifiersChangedCallback(events.WindowModifiersChangedCallback)

	// SetKeyboardInputCallback registers a callback to receive keyboard events for the window
	//
	SetKeyboardInputCallback(events.WindowKeyboardInputCallback)

	// SetReceivedCharacterCallback registers a callback to receive a UTF-8 character
	// for the corresponding key press inside the window
	//
	SetReceivedCharacterCallback(events.WindowReceivedCharacterCallback)
}

// Extension methods for android backend
type AndroidWindowExt interface {
	// SetSurfaceCreatedCallback registers a callback to receive an event when
	// ANativeWindow is ready for use
	//
	SetSurfaceCreatedCallback(events.WindowSurfaceCreatedCallback)

	// SetSurfaceDestroyedCallback registers a callback to receive an event when
	// ANativeWindow is no longer usable
	//
	SetSurfaceDestroyedCallback(events.WindowSurfaceDestroyedCallback)

	// EnableIme shows the software keyboard
	//
	EnableIme()

	// EnableIme hides the software keyboard
	//
	DisableIme()
}
