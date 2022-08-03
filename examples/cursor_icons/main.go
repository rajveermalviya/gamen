package main

import (
	"runtime"

	"github.com/rajveermalviya/gamen/cursors"
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

	w.SetTitle("gamen cursor_icons example")

	current := 0
	w.SetKeyboardInputCallback(func(state events.ButtonState, scanCode events.ScanCode, virtualKeyCode events.VirtualKey) {
		if state == events.ButtonStatePressed {
			w.SetCursorIcon(icons[current])

			if current < len(icons)-1 {
				current += 1
			} else {
				current = 0
			}
		}
	})

	w.SetCloseRequestedCallback(func() { d.Destroy() })

	for {
		// render here

		if !d.Wait() {
			break
		}
	}
}

var icons = []cursors.Icon{
	cursors.Default,
	cursors.ContextMenu,
	cursors.Help,
	cursors.Pointer,
	cursors.Progress,
	cursors.Wait,
	cursors.Cell,
	cursors.Crosshair,
	cursors.Text,
	cursors.VerticalText,
	cursors.Alias,
	cursors.Copy,
	cursors.Move,
	cursors.NoDrop,
	cursors.NotAllowed,
	cursors.Grab,
	cursors.Grabbing,
	cursors.AllScroll,
	cursors.ColResize,
	cursors.RowResize,
	cursors.NResize,
	cursors.EResize,
	cursors.SResize,
	cursors.WResize,
	cursors.NEResize,
	cursors.NWResize,
	cursors.SEResize,
	cursors.SWResize,
	cursors.EWResize,
	cursors.NSResize,
	cursors.NESWResize,
	cursors.NWSEResize,
	cursors.ZoomIn,
	cursors.ZoomOut,
}
