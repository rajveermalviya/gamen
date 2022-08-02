//go:build js

package web

import (
	"sync"
	"sync/atomic"
	"syscall/js"

	"github.com/rajveermalviya/gamen/cursors"
	"github.com/rajveermalviya/gamen/dpi"
	"github.com/rajveermalviya/gamen/events"
)

var windowCounter uint64

type Window struct {
	d           *Display
	destroyOnce sync.Once

	canvas    js.Value
	listeners map[string]js.Func

	currentCursorIcon string

	// callbacks
	resizedCb           events.WindowResizedCallback
	closeRequestedCb    events.WindowCloseRequestedCallback
	focusedCb           events.WindowFocusedCallback
	unfocusedCb         events.WindowUnfocusedCallback
	cursorEnteredCb     events.WindowCursorEnteredCallback
	cursorLeftCb        events.WindowCursorLeftCallback
	cursorMovedCb       events.WindowCursorMovedCallback
	mouseWheelCb        events.WindowMouseScrollCallback
	mouseInputCb        events.WindowMouseInputCallback
	modifiersChangedCb  events.WindowModifiersChangedCallback
	keyboardInputCb     events.WindowKeyboardInputCallback
	receivedCharacterCb events.WindowReceivedCharacterCallback
}

func NewWindow(d *Display) (*Window, error) {
	id := atomic.AddUint64(&windowCounter, 1)

	document := js.Global().Get("document")
	canvas := document.Call("createElement", "canvas")
	canvas.Call("setAttribute", "tabindex", "0")

	setCanvasSize(canvas, dpi.LogicalSize[float64]{
		Width:  640,
		Height: 480,
	})

	w := &Window{
		d:                 d,
		canvas:            canvas,
		currentCursorIcon: "auto",
		listeners:         make(map[string]js.Func, 11),
	}

	setHandlers(w)

	d.windows[id] = w
	return w, nil
}

func (w *Window) Destroy() {
	w.destroyOnce.Do(func() {
		w.resizedCb = nil
		w.closeRequestedCb = nil
		w.focusedCb = nil
		w.unfocusedCb = nil
		w.cursorEnteredCb = nil
		w.cursorLeftCb = nil
		w.cursorMovedCb = nil
		w.mouseWheelCb = nil
		w.mouseInputCb = nil
		w.modifiersChangedCb = nil
		w.keyboardInputCb = nil
		w.receivedCharacterCb = nil

		for event, listener := range w.listeners {
			w.canvas.Call("removeEventListener", event, listener)
		}
		w.listeners = nil

		w.canvas.Call("remove")
	})
}

func (w *Window) WebCanvas() js.Value { return w.canvas }

func (w *Window) SetTitle(title string) {
	js.Global().Get("document").Set("title", title)
}

func (w *Window) InnerSize() dpi.PhysicalSize[uint32] {
	return dpi.PhysicalSize[uint32]{
		Width:  uint32(w.canvas.Get("width").Int()),
		Height: uint32(w.canvas.Get("height").Int()),
	}
}

func (w *Window) SetInnerSize(size dpi.Size[uint32]) {
	old := w.InnerSize()
	setCanvasSize(w.canvas, dpi.CastSize[uint32, float64](size))
	new := w.InnerSize()

	if old != new {
		go func() {
			w.d.eventCallbacksChan <- func() {
				if w.resizedCb != nil {
					w.resizedCb(new.Width, new.Height, scaleFactor())
				}
			}
		}()
	}
}

func (w *Window) SetCursorIcon(icon cursors.Icon) {
	if icon == cursors.Default {
		w.currentCursorIcon = "auto"
	} else {
		w.currentCursorIcon = icon.String()
	}

	w.canvas.Get("style").Call("setProperty", "cursor", w.currentCursorIcon)
}

func (w *Window) SetCursorVisible(visible bool) {
	var icon string
	if visible {
		icon = w.currentCursorIcon
	} else {
		icon = "none"
	}

	w.canvas.Get("style").Call("setProperty", "cursor", icon)
}

func (w *Window) SetFullscreen(fullscreen bool) {
	if fullscreen {
		w.canvas.Call("requestFullscreen")
	} else {
		js.Global().Get("document").Call("exitFullscreen")
	}
}
func (w *Window) Fullscreen() bool {
	el := js.Global().Get("document").Get("fullscreenElement")
	if w.canvas.Equal(el) {
		return true
	}
	return false
}

func (*Window) SetMinInnerSize(size dpi.Size[uint32]) {}
func (*Window) SetMaxInnerSize(size dpi.Size[uint32]) {}
func (*Window) Maximized() bool                       { return false }
func (*Window) SetMinimized()                         {}
func (*Window) SetMaximized(maximized bool)           {}

func (w *Window) addListener(eventName string, f func(event js.Value)) {
	listener := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) == 0 {
			return js.Undefined()
		}

		event := args[0]
		event.Call("stopPropagation")

		f(event)

		return js.Undefined()
	})

	w.canvas.Call("addEventListener", eventName, listener)

	w.listeners[eventName] = listener
}

func setHandlers(w *Window) {
	w.addListener("blur", func(event js.Value) {
		w.d.eventCallbacksChan <- func() {
			if w.unfocusedCb != nil {
				w.unfocusedCb()
			}
		}
	})

	w.addListener("focus", func(event js.Value) {
		w.d.eventCallbacksChan <- func() {
			if w.focusedCb != nil {
				w.focusedCb()
			}
		}
	})

	w.addListener("pointerover", func(event js.Value) {
		w.d.eventCallbacksChan <- func() {
			if w.cursorEnteredCb != nil {
				w.cursorEnteredCb()
			}
		}
	})

	w.addListener("pointerout", func(event js.Value) {
		w.d.eventCallbacksChan <- func() {
			if w.cursorLeftCb != nil {
				w.cursorLeftCb()
			}
		}
	})

	w.addListener("pointermove", func(event js.Value) {
		physicalPosition := dpi.LogicalPosition[float64]{
			X: event.Get("offsetX").Float(),
			Y: event.Get("offsetY").Float(),
		}.ToPhysical(scaleFactor())

		w.d.eventCallbacksChan <- func() {
			if w.cursorMovedCb != nil {
				w.cursorMovedCb(physicalPosition.X, physicalPosition.Y)
			}
		}
	})

	w.addListener("wheel", func(event js.Value) {
		event.Call("preventDefault")

		const (
			DOM_DELTA_PIXEL = 0x00
			DOM_DELTA_LINE  = 0x01
			DOM_DELTA_PAGE  = 0x02
		)

		var delta events.MouseScrollDelta
		var axis events.MouseScrollAxis
		var value float64

		switch event.Get("deltaMode").Int() {
		case DOM_DELTA_PIXEL:
			delta = events.MouseScrollDeltaPixel

			x := event.Get("deltaX").Float()
			y := event.Get("deltaY").Float()

			if x != 0 {
				axis = events.MouseScrollAxisHorizontal
				value = scaleFactor() * -x
			} else if y != 0 {
				axis = events.MouseScrollAxisVertical
				value = scaleFactor() * -y
			}

		case DOM_DELTA_LINE:
			delta = events.MouseScrollDeltaLine

			x := event.Get("deltaX").Float()
			y := event.Get("deltaY").Float()

			if x != 0 {
				axis = events.MouseScrollAxisHorizontal
				value = -x
			} else if y != 0 {
				axis = events.MouseScrollAxisVertical
				value = -y
			}

		default:
			return
		}

		w.d.eventCallbacksChan <- func() {
			if w.mouseWheelCb != nil {
				w.mouseWheelCb(delta, axis, value)
			}
		}
	})

	w.addListener("pointerdown", func(event js.Value) {
		physicalPosition := dpi.LogicalPosition[float64]{
			X: event.Get("offsetX").Float(),
			Y: event.Get("offsetY").Float(),
		}.ToPhysical(scaleFactor())

		var button events.MouseButton

		i := event.Get("button").Int()
		switch i {
		case 0:
			button = events.MouseButtonLeft
		case 1:
			button = events.MouseButtonMiddle
		case 2:
			button = events.MouseButtonRight

		default:
			button = events.MouseButton(i - 3)
		}

		w.d.eventCallbacksChan <- func() {
			if w.cursorMovedCb != nil {
				w.cursorMovedCb(physicalPosition.X, physicalPosition.Y)
			}
			if w.mouseInputCb != nil {
				w.mouseInputCb(events.ButtonStatePressed, button)
			}
		}
	})

	w.addListener("pointerup", func(event js.Value) {
		var button events.MouseButton

		i := event.Get("button").Int()
		switch i {
		case 0:
			button = events.MouseButtonLeft
		case 1:
			button = events.MouseButtonMiddle
		case 2:
			button = events.MouseButtonRight

		default:
			button = events.MouseButton(i - 3)
		}

		w.d.eventCallbacksChan <- func() {
			if w.mouseInputCb != nil {
				w.mouseInputCb(events.ButtonStateReleased, button)
			}
		}
	})

	w.addListener("keydown", func(event js.Value) {
		eventKey := event.Get("key").String()
		isKeyString := len(eventKey) == 1 || !isASCII(eventKey)
		isShortcutModifiers := (event.Get("ctrlKey").Bool() || event.Get("altKey").Bool()) &&
			!event.Call("getModifierState", "AltGr").Bool()

		if !isKeyString || isShortcutModifiers {
			event.Call("preventDefault")
		}

		scanCode := event.Get("keyCode").Int()
		if scanCode == 0 {
			scanCode = event.Get("charCode").Int()
		}

		w.d.eventCallbacksChan <- func() {
			vKey, ok := mapKeyCode(event.Get("code").String())
			if !ok {
				vKey = events.VirtualKey(scanCode)
			}

			if w.keyboardInputCb != nil {
				w.keyboardInputCb(
					events.ButtonStatePressed,
					events.ScanCode(scanCode),
					vKey,
				)
			}
		}
	})

	w.addListener("keyup", func(event js.Value) {
		event.Call("preventDefault")

		scanCode := event.Get("keyCode").Int()
		if scanCode == 0 {
			scanCode = event.Get("charCode").Int()
		}

		w.d.eventCallbacksChan <- func() {
			vKey, ok := mapKeyCode(event.Get("code").String())
			if !ok {
				vKey = events.VirtualKey(scanCode)
			}

			if w.keyboardInputCb != nil {
				w.keyboardInputCb(
					events.ButtonStateReleased,
					events.ScanCode(scanCode),
					vKey,
				)
			}
		}
	})

	w.addListener("keypress", func(event js.Value) {
		event.Call("preventDefault")

		key := []rune(event.Get("key").String())
		if len(key) == 0 {
			return
		}
		char := key[0]

		w.d.eventCallbacksChan <- func() {
			if w.receivedCharacterCb != nil {
				w.receivedCharacterCb(char)
			}
		}
	})
}

func (w *Window) SetCloseRequestedCallback(cb events.WindowCloseRequestedCallback) {
	w.closeRequestedCb = cb
}
func (w *Window) SetResizedCallback(cb events.WindowResizedCallback)     { w.resizedCb = cb }
func (w *Window) SetFocusedCallback(cb events.WindowFocusedCallback)     { w.focusedCb = cb }
func (w *Window) SetUnfocusedCallback(cb events.WindowUnfocusedCallback) { w.unfocusedCb = cb }
func (w *Window) SetCursorEnteredCallback(cb events.WindowCursorEnteredCallback) {
	w.cursorEnteredCb = cb
}
func (w *Window) SetCursorLeftCallback(cb events.WindowCursorLeftCallback)   { w.cursorLeftCb = cb }
func (w *Window) SetCursorMovedCallback(cb events.WindowCursorMovedCallback) { w.cursorMovedCb = cb }
func (w *Window) SetMouseScrollCallback(cb events.WindowMouseScrollCallback) { w.mouseWheelCb = cb }
func (w *Window) SetMouseInputCallback(cb events.WindowMouseInputCallback)   { w.mouseInputCb = cb }
func (w *Window) SetTouchInputCallback(cb events.WindowTouchInputCallback) {
	// TODO:
}
func (w *Window) SetModifiersChangedCallback(cb events.WindowModifiersChangedCallback) {
	w.modifiersChangedCb = cb
}
func (w *Window) SetKeyboardInputCallback(cb events.WindowKeyboardInputCallback) {
	w.keyboardInputCb = cb
}
func (w *Window) SetReceivedCharacterCallback(cb events.WindowReceivedCharacterCallback) {
	w.receivedCharacterCb = cb
}
