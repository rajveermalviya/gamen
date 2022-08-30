package main

import (
	"fmt"
	"runtime"

	"github.com/rajveermalviya/gamen/display"
	"github.com/rajveermalviya/gamen/dpi"
	"github.com/rajveermalviya/gamen/events"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	d, err := display.NewDisplay()
	if err != nil {
		panic(err)
	}
	defer d.Destroy()

	w, err := display.NewWindow(d)
	if err != nil {
		panic(err)
	}
	defer w.Destroy()

	// On web, we don't append the canvas created by NewWindow
	// to body, users have explicitly do it.
	//
	// So this helper function does that
	initWebCanvas(w)

	w.SetTitle("gamen hello example")

	if w, ok := w.(display.AndroidWindowExt); ok {
		w.SetSurfaceCreatedCallback(func() { logf("SurfaceCreated") })
		w.SetSurfaceDestroyedCallback(func() { logf("SurfaceDestroyed") })
	}

	w.SetResizedCallback(func(physicalWidth, physicalHeight uint32, scaleFactor float64) {
		logf("Resized: physicalWidth=%v physicalHeight=%v scaleFactor=%v", physicalWidth, physicalHeight, scaleFactor)
	})

	w.SetCursorEnteredCallback(func() { logf("CursorEntered") })
	w.SetCursorLeftCallback(func() { logf("CursorLeft") })

	w.SetCursorMovedCallback(func(physicalX, physicalY float64) {
		logf("CursorMoved: physicalX=%v physicalY=%v", physicalX, physicalY)
	})

	w.SetMouseScrollCallback(func(delta events.MouseScrollDelta, axis events.MouseScrollAxis, value float64) {
		logf("MouseWheel: delta=%v axis=%v value=%v", delta, axis, value)
	})
	w.SetMouseInputCallback(func(state events.ButtonState, button events.MouseButton) {
		logf("MouseInput: state=%v button=%v", state, button)
	})

	w.SetTouchInputCallback(func(phase events.TouchPhase, location dpi.PhysicalPosition[float64], id events.TouchPointerID) {
		logf("TouchInput: phase=%v location=%v id=%v", phase, location, id)
	})

	w.SetFocusedCallback(func() { logf("Focused") })
	w.SetUnfocusedCallback(func() { logf("Unfocused") })

	w.SetModifiersChangedCallback(func(state events.ModifiersState) {
		logf("ModifiersChanged: state=%v", state)
	})
	w.SetKeyboardInputCallback(func(state events.ButtonState, scanCode events.ScanCode, virtualKeyCode events.VirtualKey) {
		logf("KeyboardInput: state=%v scanCode=%v virtualKeyCode=%v", state, scanCode, virtualKeyCode)
	})
	w.SetReceivedCharacterCallback(func(char rune) {
		logf("ReceivedCharacter: %#U", char)
	})

	w.SetCloseRequestedCallback(func() {
		logf("CloseRequested")
		d.Destroy()
	})

	for {
		if !d.Wait() {
			break
		}

		// render here
	}
}

func logf(format string, a ...any) {
	// we use println instead of fmt.Print because
	// android doesn't pipe stdout/stderr to logcat
	println(fmt.Sprintf(format, a...))
}
