package main

import (
	"fmt"
	"runtime"

	"github.com/rajveermalviya/gamen/display"
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

	windowCount := 4
	closedWindows := 0

	for i := 0; i < windowCount; i++ {
		id := i

		w, err := display.NewWindow(d)
		if err != nil {
			panic(err)
		}
		defer w.Destroy()

		w.SetCursorEnteredCallback(func() {
			fmt.Printf("CursorEntered: windowId=%d\n", id)
		})
		w.SetCursorLeftCallback(func() {
			fmt.Printf("CursorLeft: windowId=%d\n", id)
		})
		w.SetCursorMovedCallback(func(physicalX, physicalY float64) {
			fmt.Printf("CursorMoved: windowId=%d physicalX=%v physicalY=%v\n", id, physicalX, physicalY)
		})

		w.SetCloseRequestedCallback(func() {
			w.Destroy()
			closedWindows++

			// when all windows are closed, call Display.Destroy() to exit loop
			if closedWindows == windowCount {
				d.Destroy()
			}
		})
	}

loop:
	for {
		if !d.Wait() {
			break loop
		}
	}
}
