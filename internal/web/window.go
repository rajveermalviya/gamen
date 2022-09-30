//go:build js

package web

import (
	"sync"
	"sync/atomic"
	"syscall/js"

	"github.com/rajveermalviya/gamen/cursors"
	"github.com/rajveermalviya/gamen/dpi"
	"github.com/rajveermalviya/gamen/events"
	"github.com/rajveermalviya/gamen/internal/common/atomicx"
)

var windowCounter uint64

type Window struct {
	d           *Display
	destroyOnce sync.Once

	canvas    js.Value
	listeners map[string]js.Func

	currentCursorIcon string

	// callbacks
	resizedCb           atomicx.Pointer[events.WindowResizedCallback]
	closeRequestedCb    atomicx.Pointer[events.WindowCloseRequestedCallback]
	focusedCb           atomicx.Pointer[events.WindowFocusedCallback]
	unfocusedCb         atomicx.Pointer[events.WindowUnfocusedCallback]
	cursorEnteredCb     atomicx.Pointer[events.WindowCursorEnteredCallback]
	cursorLeftCb        atomicx.Pointer[events.WindowCursorLeftCallback]
	cursorMovedCb       atomicx.Pointer[events.WindowCursorMovedCallback]
	mouseWheelCb        atomicx.Pointer[events.WindowMouseScrollCallback]
	mouseInputCb        atomicx.Pointer[events.WindowMouseInputCallback]
	modifiersChangedCb  atomicx.Pointer[events.WindowModifiersChangedCallback]
	keyboardInputCb     atomicx.Pointer[events.WindowKeyboardInputCallback]
	receivedCharacterCb atomicx.Pointer[events.WindowReceivedCharacterCallback]
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
		w.resizedCb.Store(nil)
		w.closeRequestedCb.Store(nil)
		w.focusedCb.Store(nil)
		w.unfocusedCb.Store(nil)
		w.cursorEnteredCb.Store(nil)
		w.cursorLeftCb.Store(nil)
		w.cursorMovedCb.Store(nil)
		w.mouseWheelCb.Store(nil)
		w.mouseInputCb.Store(nil)
		w.modifiersChangedCb.Store(nil)
		w.keyboardInputCb.Store(nil)
		w.receivedCharacterCb.Store(nil)

		for event, listener := range w.listeners {
			w.canvas.Call("removeEventListener", event, listener)
			listener.Release()
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
				if cb := w.resizedCb.Load(); cb != nil {
					if cb := (*cb); cb != nil {
						cb(new.Width, new.Height, scaleFactor())
					}
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

func (*Window) SetMinInnerSize(dpi.Size[uint32]) {}
func (*Window) SetMaxInnerSize(dpi.Size[uint32]) {}
func (*Window) Maximized() bool                  { return false }
func (*Window) SetMinimized()                    {}
func (*Window) SetMaximized(bool)                {}
func (*Window) DragWindow()                      {}
func (*Window) SetDecorations(bool)              {}
func (*Window) Decorated() bool                  { return false }

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
			if cb := w.unfocusedCb.Load(); cb != nil {
				if cb := (*cb); cb != nil {
					cb()
				}
			}
		}
	})

	w.addListener("focus", func(event js.Value) {
		w.d.eventCallbacksChan <- func() {
			if cb := w.focusedCb.Load(); cb != nil {
				if cb := (*cb); cb != nil {
					cb()
				}
			}
		}
	})

	w.addListener("pointerover", func(event js.Value) {
		w.d.eventCallbacksChan <- func() {
			if cb := w.cursorEnteredCb.Load(); cb != nil {
				if cb := (*cb); cb != nil {
					cb()
				}
			}
		}
	})

	w.addListener("pointerout", func(event js.Value) {
		w.d.eventCallbacksChan <- func() {
			if cb := w.cursorLeftCb.Load(); cb != nil {
				if cb := (*cb); cb != nil {
					cb()
				}
			}
		}
	})

	w.addListener("pointermove", func(event js.Value) {
		physicalPosition := dpi.LogicalPosition[float64]{
			X: event.Get("offsetX").Float(),
			Y: event.Get("offsetY").Float(),
		}.ToPhysical(scaleFactor())

		w.d.eventCallbacksChan <- func() {
			if cb := w.cursorMovedCb.Load(); cb != nil {
				if cb := (*cb); cb != nil {
					cb(physicalPosition.X, physicalPosition.Y)
				}
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
			if cb := w.mouseWheelCb.Load(); cb != nil {
				if cb := (*cb); cb != nil {
					cb(delta, axis, value)
				}
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
			if cb := w.cursorMovedCb.Load(); cb != nil {
				if cb := (*cb); cb != nil {
					cb(physicalPosition.X, physicalPosition.Y)
				}
			}
			if cb := w.mouseInputCb.Load(); cb != nil {
				if cb := (*cb); cb != nil {
					cb(events.ButtonStatePressed, button)
				}
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
			if cb := w.mouseInputCb.Load(); cb != nil {
				if cb := (*cb); cb != nil {
					cb(events.ButtonStateReleased, button)
				}
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

			if cb := w.keyboardInputCb.Load(); cb != nil {
				if cb := (*cb); cb != nil {
					cb(
						events.ButtonStatePressed,
						events.ScanCode(scanCode),
						vKey,
					)
				}
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

			if cb := w.keyboardInputCb.Load(); cb != nil {
				if cb := (*cb); cb != nil {
					cb(
						events.ButtonStateReleased,
						events.ScanCode(scanCode),
						vKey,
					)
				}
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
			if cb := w.receivedCharacterCb.Load(); cb != nil {
				if cb := (*cb); cb != nil {
					cb(char)
				}
			}
		}
	})
}

func (w *Window) SetCloseRequestedCallback(cb events.WindowCloseRequestedCallback) {
	w.closeRequestedCb.Store(&cb)
}
func (w *Window) SetResizedCallback(cb events.WindowResizedCallback) {
	w.resizedCb.Store(&cb)
}
func (w *Window) SetFocusedCallback(cb events.WindowFocusedCallback) {
	w.focusedCb.Store(&cb)
}
func (w *Window) SetUnfocusedCallback(cb events.WindowUnfocusedCallback) {
	w.unfocusedCb.Store(&cb)
}
func (w *Window) SetCursorEnteredCallback(cb events.WindowCursorEnteredCallback) {
	w.cursorEnteredCb.Store(&cb)
}
func (w *Window) SetCursorLeftCallback(cb events.WindowCursorLeftCallback) {
	w.cursorLeftCb.Store(&cb)
}
func (w *Window) SetCursorMovedCallback(cb events.WindowCursorMovedCallback) {
	w.cursorMovedCb.Store(&cb)
}
func (w *Window) SetMouseScrollCallback(cb events.WindowMouseScrollCallback) {
	w.mouseWheelCb.Store(&cb)
}
func (w *Window) SetMouseInputCallback(cb events.WindowMouseInputCallback) {
	w.mouseInputCb.Store(&cb)
}
func (w *Window) SetTouchInputCallback(cb events.WindowTouchInputCallback) {
	// TODO:
}
func (w *Window) SetModifiersChangedCallback(cb events.WindowModifiersChangedCallback) {
	w.modifiersChangedCb.Store(&cb)
}
func (w *Window) SetKeyboardInputCallback(cb events.WindowKeyboardInputCallback) {
	w.keyboardInputCb.Store(&cb)
}
func (w *Window) SetReceivedCharacterCallback(cb events.WindowReceivedCharacterCallback) {
	w.receivedCharacterCb.Store(&cb)
}
