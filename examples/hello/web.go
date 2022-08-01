//go:build js

package main

import (
	"syscall/js"

	"github.com/rajveermalviya/gamen/display"
)

func initWebCanvas(w display.Window) {
	if w, ok := w.(display.WebWindow); ok {
		canvas := w.WebCanvas()
		// set a color to make it visible
		canvas.Get("style").Set("background-color", "blue")
		js.Global().Get("document").Get("body").Call("appendChild", canvas)
	}
}
