package glfw

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/rajveermalviya/gamen/cursors"
	"github.com/rajveermalviya/gamen/dpi"
	"github.com/rajveermalviya/gamen/events"
)

type Window struct {
	w *glfw.Window
}

func (w *Window) GlfwWindow() *glfw.Window {
	return w.w
}

func (w *Window) Destroy() {
	println("glfw.DestroyWindow()")
	w.w.Destroy()
}

func (w *Window) SetTitle(title string) {
	w.w.SetTitle(title)
}

func (w *Window) InnerSize() dpi.PhysicalSize[uint32] {
	width, height := w.w.GetSize()
	return dpi.PhysicalSize[uint32]{Width: uint32(width), Height: uint32(height)}
}

func (w *Window) SetInnerSize(size dpi.Size[uint32]) {
	psize := size.ToPhysical(1)
	w.w.SetSize(int(psize.Width), int(psize.Height))
}

func (w *Window) SetMinInnerSize(size dpi.Size[uint32]) {
	psize := size.ToPhysical(1)
	w.w.SetSizeLimits(int(psize.Width), int(psize.Height), glfw.DontCare, glfw.DontCare)
}

func (w *Window) SetMaxInnerSize(size dpi.Size[uint32]) {
	psize := size.ToPhysical(1)
	w.w.SetSizeLimits(glfw.DontCare, glfw.DontCare, int(psize.Width), int(psize.Height))
}

func (w *Window) Maximized() bool {
	return w.w.GetAttrib(glfw.Maximized) != 0
}

func (w *Window) SetMinimized() {
	w.w.Iconify()
}

func (w *Window) SetMaximized(value bool) {
	if value {
		w.w.Maximize()
	} else {
		w.w.Restore()
	}
}

func (w *Window) SetCursorIcon(cursor cursors.Icon) {
	c := glfw.CreateStandardCursor(glfw.ArrowCursor)
	w.w.SetCursor(c)
}

func (w *Window) SetCursorVisible(value bool) {
	if value {
		w.w.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	} else {
		w.w.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
	}
}

func (w *Window) SetFullscreen(value bool) {
	if value {
		w.w.SetMonitor(glfw.GetPrimaryMonitor(), 0, 0, 0, 0, 0)
	} else {
		w.w.SetMonitor(nil, 0, 0, 0, 0, 0)
	}
}

func (w *Window) Fullscreen() bool {
	return w.w.GetMonitor() != nil
}

func (w *Window) DragWindow() {

}

func (w *Window) SetDecorations(value bool) {
	if value {
		w.w.SetAttrib(glfw.Decorated, glfw.True)
	} else {
		w.w.SetAttrib(glfw.Decorated, glfw.False)
	}
}

func (w *Window) Decorated() bool {
	return w.w.GetAttrib(glfw.Decorated) != 0
}

func (w *Window) SetResizedCallback(cb events.WindowResizedCallback) {
	w.w.SetSizeCallback(func(w *glfw.Window, width int, height int) {
		cb(uint32(width), uint32(height), 1.0)
	})
}

func (w *Window) SetCloseRequestedCallback(cb events.WindowCloseRequestedCallback) {
	w.w.SetCloseCallback(func(w *glfw.Window) {
		cb()
	})
}

func (w *Window) SetCursorEnteredCallback(cb events.WindowCursorEnteredCallback) {
	w.w.SetCursorEnterCallback(func(w *glfw.Window, entered bool) {
		if entered {
			cb()
		}
	})
}

func (w *Window) SetCursorLeftCallback(cb events.WindowCursorLeftCallback) {
	w.w.SetCursorEnterCallback(func(w *glfw.Window, entered bool) {
		if !entered {
			cb()
		}
	})
}

func (w *Window) SetCursorMovedCallback(cb events.WindowCursorMovedCallback) {
	w.w.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		cb(xpos, ypos)
	})
}

func (w *Window) SetMouseScrollCallback(cb events.WindowMouseScrollCallback) {
	w.w.SetScrollCallback(func(w *glfw.Window, xoff float64, yoff float64) {
		if xoff != 0 {
			cb(events.MouseScrollDeltaPixel, events.MouseScrollAxisHorizontal, xoff)
		}
		if yoff != 0 {
			cb(events.MouseScrollDeltaPixel, events.MouseScrollAxisVertical, yoff)
		}
	})
}

func (w *Window) SetMouseInputCallback(cb events.WindowMouseInputCallback) {
	w.w.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		var state events.ButtonState
		switch action {
		case glfw.Press:
			state = events.ButtonStatePressed
		case glfw.Release:
			state = events.ButtonStateReleased
		default:
			return
		}

		var btn events.MouseButton
		switch button {
		case glfw.MouseButtonLeft:
			btn = events.MouseButtonLeft
		case glfw.MouseButtonRight:
			btn = events.MouseButtonRight
		case glfw.MouseButtonMiddle:
			btn = events.MouseButtonMiddle
		default:
			return
		}

		cb(state, btn)
	})
}

func (w *Window) SetTouchInputCallback(cb events.WindowTouchInputCallback) {

}

func (w *Window) SetFocusedCallback(cb events.WindowFocusedCallback) {

}

func (w *Window) SetUnfocusedCallback(cb events.WindowUnfocusedCallback) {

}

func (w *Window) SetModifiersChangedCallback(cb events.WindowModifiersChangedCallback) {
	w.w.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		var mod events.ModifiersState
		if mods&glfw.ModShift != 0 {
			mod |= events.ModifiersStateShift
		}
		if mods&glfw.ModControl != 0 {
			mod |= events.ModifiersStateCtrl
		}
		if mods&glfw.ModAlt != 0 {
			mod |= events.ModifiersStateAlt
		}
		if mods&glfw.ModSuper != 0 {
			mod |= events.ModifiersStateLogo
		}
		if mod != 0 {
			cb(mod)
		}
	})
}

func (w *Window) SetKeyboardInputCallback(cb events.WindowKeyboardInputCallback) {
	w.w.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		virtualKey := glfwKeyToVirtualKey(key)
		var state events.ButtonState
		switch action {
		case glfw.Press:
			state = events.ButtonStatePressed
		case glfw.Release:
			state = events.ButtonStateReleased
		default:
			return
		}
		cb(state, events.ScanCode(scancode), virtualKey)
	})
}

func (w *Window) SetReceivedCharacterCallback(cb events.WindowReceivedCharacterCallback) {
	w.w.SetCharCallback(func(w *glfw.Window, char rune) {
		cb(char)
	})
}
