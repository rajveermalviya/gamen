//go:build linux && !android

package wayland

/*

#include <stdlib.h>
#include "wayland-util.h"
#include "wayland-cursor.h"

*/
import "C"

import (
	"runtime/cgo"
	"sync"
	"time"
	"unsafe"

	"github.com/rajveermalviya/gamen/events"
)

type Pointer struct {
	d  *Display
	mu sync.Mutex

	pointer *C.struct_wl_pointer

	serial uint32
	focus  *C.struct_wl_surface

	pixelDeltaVertical   float64
	pixelDeltaHorizontal float64

	lineDeltaVertical   float64
	lineDeltaHorizontal float64

	currentCursor                   *C.struct_wl_cursor
	currentCursorAnimationStartTime time.Time
	cursorThemes                    map[uint32]*C.struct_wl_cursor_theme
	cursorSurface                   *C.struct_wl_surface
	cursorSurfaceFrameCallback      *C.struct_wl_callback
}

func (p *Pointer) destroy() {
	if p.currentCursor != nil {
		p.currentCursor = nil
	}

	if p.cursorSurfaceFrameCallback != nil {
		p.d.l.wl_callback_destroy(p.cursorSurfaceFrameCallback)
		p.cursorSurfaceFrameCallback = nil
	}

	if p.cursorSurface != nil {
		p.d.l.wl_surface_destroy(p.cursorSurface)
		p.cursorSurface = nil
	}

	if p.cursorThemes != nil {
		for _, theme := range p.cursorThemes {
			p.d.l.wl_cursor_theme_destroy(theme)
		}
		p.cursorThemes = nil
	}

	if p.pointer != nil {
		p.d.l.wl_pointer_destroy(p.pointer)
		p.pointer = nil
	}
}

func (p *Pointer) loadCursor(name string, size uint32, scaleFactor float64) *C.struct_wl_cursor {
	p.mu.Lock()
	defer p.mu.Unlock()

	size = size * uint32(scaleFactor)

	theme, ok := p.cursorThemes[size]
	if !ok {
		theme = p.d.l.wl_cursor_theme_load(nil, C.int(size), p.d.shm)
		p.cursorThemes[size] = theme
	}

	nameStr := C.CString(name)
	defer C.free(unsafe.Pointer(nameStr))

	cursor := p.d.l.wl_cursor_theme_get_cursor(theme, nameStr)
	if cursor == nil {
		return nil
	}
	if cursor.image_count == 0 {
		return nil
	}

	return cursor
}

func (p *Pointer) setCursor(cursor *C.struct_wl_cursor, name string, scaleFactor float64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// hide cursor
	if cursor == nil {
		p.d.l.wl_pointer_set_cursor(p.pointer, C.uint32_t(p.serial), nil, 0, 0)
		p.currentCursor = nil
		return
	}

	if p.cursorSurfaceFrameCallback != nil {
		p.d.l.wl_callback_destroy(p.cursorSurfaceFrameCallback)
		p.cursorSurfaceFrameCallback = nil
	}

	imageSlice := unsafe.Slice(cursor.images, cursor.image_count)
	image := imageSlice[0]
	cursorBuffer := p.d.l.wl_cursor_image_get_buffer(image)

	if p.cursorSurface == nil {
		p.cursorSurface = p.d.l.wl_compositor_create_surface(p.d.compositor)
	}

	p.d.l.wl_surface_set_buffer_scale(p.cursorSurface, C.int32_t(scaleFactor))
	p.d.l.wl_surface_attach(p.cursorSurface, cursorBuffer, 0, 0)
	p.d.l.wl_surface_damage_buffer(p.cursorSurface, 0, 0, C.int32_t(image.width), C.int32_t(image.height))
	p.d.l.wl_surface_commit(p.cursorSurface)

	p.d.l.wl_pointer_set_cursor(
		p.pointer,
		C.uint32_t(p.serial),
		p.cursorSurface,
		C.int32_t(float64(image.hotspot_x)/scaleFactor),
		C.int32_t(float64(image.hotspot_y)/scaleFactor),
	)
	p.currentCursor = cursor

	if cursor.image_count > 1 {
		p.currentCursorAnimationStartTime = time.Now()
		p.startAnimatingCursor()
	}
}

func (p *Pointer) startAnimatingCursor() {
	var fn func()
	fn = func() {
		p.mu.Lock()
		defer p.mu.Unlock()

		if p.cursorSurfaceFrameCallback != nil {
			p.d.l.wl_callback_destroy(p.cursorSurfaceFrameCallback)
			p.cursorSurfaceFrameCallback = nil
		}

		if p.currentCursor == nil {
			return
		}

		imageIdx := p.d.l.wl_cursor_frame_and_duration(
			p.currentCursor,
			C.uint32_t(p.currentCursorAnimationStartTime.UnixMilli()),
			nil,
		)

		imageSlice := unsafe.Slice(p.currentCursor.images, p.currentCursor.image_count)
		image := imageSlice[imageIdx]
		cursorBuffer := p.d.l.wl_cursor_image_get_buffer(image)

		p.d.l.wl_surface_attach(p.cursorSurface, cursorBuffer, 0, 0)
		p.d.l.wl_surface_damage_buffer(p.cursorSurface, 0, 0, C.int32_t(image.width), C.int32_t(image.height))

		p.cursorSurfaceFrameCallback = p.d.l.wl_surface_frame(p.cursorSurface)
		p.d.setCallbackListener(p.cursorSurfaceFrameCallback, fn)
		p.d.l.wl_surface_commit(p.cursorSurface)

		p.currentCursorAnimationStartTime = time.Now()
	}

	p.cursorSurfaceFrameCallback = p.d.l.wl_surface_frame(p.cursorSurface)
	p.d.setCallbackListener(p.cursorSurfaceFrameCallback, fn)
	p.d.l.wl_surface_commit(p.cursorSurface)
}

//export pointerHandleEnter
func pointerHandleEnter(data unsafe.Pointer, wl_pointer *C.struct_wl_pointer, serial C.uint32_t, surface *C.struct_wl_surface, surface_x C.wl_fixed_t, surface_y C.wl_fixed_t) {
	if surface == nil {
		return
	}

	d, ok := (*cgo.Handle)(data).Value().(*Display)
	if !ok {
		return
	}

	d.pointer.mu.Lock()
	d.pointer.serial = uint32(serial)
	d.pointer.focus = surface
	d.pointer.mu.Unlock()

	w, ok := d.windows[surface]
	if !ok {
		return
	}

	w.mu.Lock()
	currentCursorIcon := w.currentCursorIcon
	scaleFactor := w.scaleFactor
	w.mu.Unlock()

	// user can call window.SetCursor when window is not in focus
	// so we just store the state so when pointer enters window
	// we set cursor to how the window's state requires it
	if currentCursorIcon == "" {
		d.pointer.setCursor(nil, "", 0)
	} else {
		cursor := d.pointer.loadCursor(currentCursorIcon, 24, scaleFactor)
		if cursor != nil {
			d.pointer.setCursor(cursor, currentCursorIcon, scaleFactor)
		}
	}

	if cb := w.cursorEnteredCb.Load(); cb != nil {
		if cb := (*cb); cb != nil {
			cb()
		}
	}
}

//export pointerHandleLeave
func pointerHandleLeave(data unsafe.Pointer, wl_pointer *C.struct_wl_pointer, serial C.uint32_t, surface *C.struct_wl_surface) {
	if surface == nil {
		return
	}

	d, ok := (*cgo.Handle)(data).Value().(*Display)
	if !ok {
		return
	}

	d.pointer.mu.Lock()
	d.pointer.serial = uint32(serial)
	d.pointer.focus = nil
	d.pointer.mu.Unlock()

	d.pointer.setCursor(nil, "", 0)

	w, ok := d.windows[surface]
	if !ok {
		return
	}

	if cb := w.cursorLeftCb.Load(); cb != nil {
		if cb := (*cb); cb != nil {
			cb()
		}
	}
}

//export pointerHandleMotion
func pointerHandleMotion(data unsafe.Pointer, wl_pointer *C.struct_wl_pointer, time C.uint32_t, surface_x C.double, surface_y C.double) {
	d, ok := (*cgo.Handle)(data).Value().(*Display)
	if !ok {
		return
	}

	d.pointer.mu.Lock()
	focus := d.pointer.focus
	d.pointer.mu.Unlock()

	if focus == nil {
		return
	}

	w, ok := d.windows[focus]
	if !ok {
		return
	}

	if cb := w.cursorMovedCb.Load(); cb != nil {
		if cb := (*cb); cb != nil {
			cb(float64(surface_x), float64(surface_y))
		}
	}
}

//export pointerHandleButton
func pointerHandleButton(data unsafe.Pointer, wl_pointer *C.struct_wl_pointer, serial C.uint32_t, time C.uint32_t, button C.uint32_t, state enum_wl_pointer_button_state) {
	const (
		BTN_LEFT   = 272
		BTN_RIGHT  = 273
		BTN_MIDDLE = 274
	)

	d, ok := (*cgo.Handle)(data).Value().(*Display)
	if !ok {
		return
	}

	d.pointer.mu.Lock()
	d.pointer.serial = uint32(serial)
	focus := d.pointer.focus
	d.pointer.mu.Unlock()

	if focus == nil {
		return
	}

	w, ok := d.windows[focus]
	if !ok {
		return
	}

	if cb := w.mouseInputCb.Load(); cb != nil {
		if cb := (*cb); cb != nil {
			var s events.ButtonState
			switch state {
			case WL_POINTER_BUTTON_STATE_PRESSED:
				s = events.ButtonStatePressed
			case WL_POINTER_BUTTON_STATE_RELEASED:
				s = events.ButtonStateReleased
			}

			var b events.MouseButton
			switch button {
			case BTN_LEFT:
				b = events.MouseButtonLeft
			case BTN_RIGHT:
				b = events.MouseButtonRight
			case BTN_MIDDLE:
				b = events.MouseButtonMiddle
			default:
				b = events.MouseButton(button)
			}

			cb(s, b)
		}
	}
}

//export pointerHandleAxis
func pointerHandleAxis(data unsafe.Pointer, wl_pointer *C.struct_wl_pointer, time C.uint32_t, axis enum_wl_pointer_axis, value C.double) {
	d, ok := (*cgo.Handle)(data).Value().(*Display)
	if !ok {
		return
	}

	// we call callbacks on frame event
	switch axis {
	case WL_POINTER_AXIS_VERTICAL_SCROLL:
		d.pointer.pixelDeltaVertical -= float64(value)

	case WL_POINTER_AXIS_HORIZONTAL_SCROLL:
		d.pointer.pixelDeltaHorizontal -= float64(value)
	}
}

//export pointerHandleAxisDiscrete
func pointerHandleAxisDiscrete(data unsafe.Pointer, wl_pointer *C.struct_wl_pointer, axis enum_wl_pointer_axis, discrete C.int32_t) {
	d, ok := (*cgo.Handle)(data).Value().(*Display)
	if !ok {
		return
	}

	// we call callbacks on frame event
	switch axis {
	case WL_POINTER_AXIS_VERTICAL_SCROLL:
		d.pointer.lineDeltaVertical -= float64(discrete)

	case WL_POINTER_AXIS_HORIZONTAL_SCROLL:
		d.pointer.lineDeltaHorizontal -= float64(discrete)
	}
}

//export pointerHandleFrame
func pointerHandleFrame(data unsafe.Pointer, wl_pointer *C.struct_wl_pointer) {
	d, ok := (*cgo.Handle)(data).Value().(*Display)
	if !ok {
		return
	}

	d.pointer.mu.Lock()
	focus := d.pointer.focus
	d.pointer.mu.Unlock()

	if focus == nil {
		return
	}

	w, ok := d.windows[focus]
	if !ok || w == nil {
		return
	}

	if d.pointer.lineDeltaVertical != 0 {
		if cb := w.mouseWheelCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				cb(
					events.MouseScrollDeltaLine,
					events.MouseScrollAxisVertical,
					d.pointer.lineDeltaVertical,
				)
			}
		}

		d.pointer.lineDeltaVertical = 0
	} else if d.pointer.pixelDeltaVertical != 0 {
		if cb := w.mouseWheelCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				cb(
					events.MouseScrollDeltaPixel,
					events.MouseScrollAxisVertical,
					d.pointer.pixelDeltaVertical,
				)
			}
		}

		d.pointer.pixelDeltaVertical = 0
	} else if d.pointer.lineDeltaHorizontal != 0 {
		if cb := w.mouseWheelCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				cb(
					events.MouseScrollDeltaLine,
					events.MouseScrollAxisHorizontal,
					d.pointer.lineDeltaHorizontal,
				)
			}
		}

		d.pointer.lineDeltaHorizontal = 0
	} else if d.pointer.pixelDeltaHorizontal != 0 {
		if cb := w.mouseWheelCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				cb(
					events.MouseScrollDeltaPixel,
					events.MouseScrollAxisHorizontal,
					d.pointer.pixelDeltaHorizontal,
				)
			}
		}

		d.pointer.pixelDeltaHorizontal = 0
	}
}
