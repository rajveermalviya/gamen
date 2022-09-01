//go:build linux && !android

package wayland

/*

#include "wayland-client-protocol.h"

*/
import "C"

import (
	"runtime/cgo"
	"unsafe"
)

type Output struct {
	output *C.struct_wl_output
	name   uint32

	// from geometry event
	x, y                          int32
	physicalWidth, physicalHeight int32
	subpixel                      enum_wl_output_subpixel
	make, model                   string
	transform                     enum_wl_output_transform

	// from mode event
	flags         enum_wl_output_mode
	width, height int32
	refresh       int32

	scaleFactor int32
}

//export outputHandleGeometry
func outputHandleGeometry(data unsafe.Pointer, wl_output *C.struct_wl_output,
	x, y C.int32_t,
	physical_width, physical_height C.int32_t,
	subpixel enum_wl_output_subpixel,
	make, model *C.char,
	transform enum_wl_output_transform,
) {
	d, ok := (*cgo.Handle)(data).Value().(*Display)
	if !ok {
		return
	}

	output, ok := d.outputs[wl_output]
	if !ok {
		return
	}

	output.x = int32(x)
	output.y = int32(y)
	output.physicalWidth = int32(physical_width)
	output.physicalHeight = int32(physical_height)
	output.subpixel = subpixel
	output.make = C.GoString(make)
	output.model = C.GoString(model)
	output.transform = transform
}

//export outputHandleMode
func outputHandleMode(data unsafe.Pointer, wl_output *C.struct_wl_output,
	flags enum_wl_output_mode,
	width, height C.int32_t,
	refresh C.int32_t,
) {
	d, ok := (*cgo.Handle)(data).Value().(*Display)
	if !ok {
		return
	}

	output, ok := d.outputs[wl_output]
	if !ok {
		return
	}

	output.flags = flags
	output.width = int32(width)
	output.height = int32(height)
	output.refresh = int32(refresh)
}

//export outputHandleScale
func outputHandleScale(data unsafe.Pointer, wl_output *C.struct_wl_output, factor C.int32_t) {
	d, ok := (*cgo.Handle)(data).Value().(*Display)
	if !ok {
		return
	}

	output, ok := d.outputs[wl_output]
	if !ok {
		return
	}

	output.scaleFactor = int32(factor)
}

//export outputHandleDone
func outputHandleDone(data unsafe.Pointer, wl_output *C.struct_wl_output) {
	d, ok := (*cgo.Handle)(data).Value().(*Display)
	if !ok {
		return
	}

	for _, w := range d.windows {
		func(w *Window) {
			w.mu.Lock()
			defer w.mu.Unlock()
			w.updateScaleFactor()
		}(w)
	}
}
