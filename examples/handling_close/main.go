package main

import (
	"runtime"

	"github.com/rajveermalviya/gamen/display"
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

	w.SetTitle("gamen handling_close example")

	if _, ok := w.(display.AndroidWindowExt); ok {
		// on Android exit when back button is pressed
		w.SetKeyboardInputCallback(func(state events.ButtonState, scanCode events.ScanCode, virtualKeyCode events.VirtualKey) {
			if state == events.ButtonStatePressed && virtualKeyCode == 4 { // TODO: add this to events.VirtualKey constants
				d.Destroy()
			}
		})
	} else {
		w.SetCloseRequestedCallback(func() { d.Destroy() })
	}

	for {
		// render here

		if !d.Wait() {
			break
		}
	}
}
