//go:build js

package web

import (
	"strconv"
	"syscall/js"

	"github.com/rajveermalviya/gamen/dpi"
)

func scaleFactor() float64 {
	return js.Global().Get("devicePixelRatio").Float()
}

func setCanvasSize(canvas js.Value, size dpi.Size[float64]) {
	scaleFactor := scaleFactor()

	physicalSize := size.ToPhysical(scaleFactor)
	logicalSize := size.ToLogical(scaleFactor)

	canvas.Set("width", physicalSize.Width)
	canvas.Set("height", physicalSize.Height)

	canvas.Get("style").Call(
		"setProperty",
		"width",
		strconv.FormatFloat(logicalSize.Width, 'f', 6, 64)+"px",
	)
	canvas.Get("style").Call(
		"setProperty",
		"height",
		strconv.FormatFloat(logicalSize.Height, 'f', 6, 64)+"px",
	)
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > 127 {
			return false
		}
	}
	return true
}
