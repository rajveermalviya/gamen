//go:build linux && !android

package wayland

/*

#include <stdlib.h>
#include <wayland-client.h>
#include "xdg-shell-client-protocol.h"
#include "xdg-decoration-unstable-v1-client-protocol.h"

extern const struct wl_surface_listener window_surface_listener;
extern const struct xdg_surface_listener xdg_surface_listener;
extern const struct xdg_toplevel_listener xdg_toplevel_listener;
extern const struct zxdg_toplevel_decoration_v1_listener zxdg_toplevel_decoration_v1_listener;

*/
import "C"

import (
	"math"
	"runtime/cgo"
	"sync"
	"unsafe"

	"github.com/rajveermalviya/gamen/cursors"
	"github.com/rajveermalviya/gamen/dpi"
	"github.com/rajveermalviya/gamen/events"
	"github.com/rajveermalviya/gamen/internal/common"
	"github.com/rajveermalviya/gamen/internal/utils"
)

type Window struct {
	// handle for Window to be passed between cgo callbacks
	handle *cgo.Handle
	d      *Display
	// we allow destroy function to be called multiple
	// times, but in reality we run it once
	destroyOnce sync.Once
	mu          sync.Mutex

	// wayland objects
	surface               *C.struct_wl_surface
	xdgSurface            *C.struct_xdg_surface
	xdgToplevel           *C.struct_xdg_toplevel
	xdgToplevelDecoration *C.struct_zxdg_toplevel_decoration_v1

	// state
	scaleFactor        float64
	size               dpi.LogicalSize[uint32]
	outputs            map[*C.struct_wl_output]struct{}
	maximized          bool
	fullscreen         bool
	previousCursorIcon string
	currentCursorIcon  string

	// window callbacks
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
	w := &Window{
		d:                 d,
		outputs:           make(map[*C.struct_wl_output]struct{}),
		scaleFactor:       1,
		currentCursorIcon: "left_ptr",
		maximized:         false,
	}
	handle := cgo.NewHandle(w)
	w.handle = &handle

	w.surface = C.wl_compositor_create_surface(d.compositor)
	C.wl_surface_add_listener(w.surface, &C.window_surface_listener, unsafe.Pointer(w.handle))

	w.xdgSurface = C.xdg_wm_base_get_xdg_surface(d.xdgWmBase, w.surface)
	C.xdg_surface_add_listener(w.xdgSurface, &C.xdg_surface_listener, unsafe.Pointer(w.handle))

	w.xdgToplevel = C.xdg_surface_get_toplevel(w.xdgSurface)
	C.xdg_toplevel_add_listener(w.xdgToplevel, &C.xdg_toplevel_listener, unsafe.Pointer(w.handle))

	if d.xdgDecorationManager != nil {
		w.xdgToplevelDecoration = C.zxdg_decoration_manager_v1_get_toplevel_decoration(d.xdgDecorationManager, w.xdgToplevel)
		C.zxdg_toplevel_decoration_v1_add_listener(w.xdgToplevelDecoration, &C.zxdg_toplevel_decoration_v1_listener, unsafe.Pointer(w.handle))
		C.zxdg_toplevel_decoration_v1_set_mode(w.xdgToplevelDecoration, C.ZXDG_TOPLEVEL_DECORATION_V1_MODE_SERVER_SIDE)
	}

	d.windows[w.surface] = w

	w.size = dpi.LogicalSize[uint32]{
		Width:  640,
		Height: 480,
	}

	return w, nil
}

func (w *Window) WlDisplay() unsafe.Pointer { return unsafe.Pointer(w.d.display) }
func (w *Window) WlSurface() unsafe.Pointer { return unsafe.Pointer(w.surface) }

func (w *Window) Destroy() {
	w.destroyOnce.Do(func() {
		w.mu.Lock()
		defer w.mu.Unlock()

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

		if _, ok := w.d.windows[w.surface]; ok {
			w.d.windows[w.surface] = nil
			delete(w.d.windows, w.surface)
		}

		if w.xdgToplevelDecoration != nil {
			C.zxdg_toplevel_decoration_v1_destroy(w.xdgToplevelDecoration)
			w.xdgToplevelDecoration = nil
		}

		if w.xdgToplevel != nil {
			C.xdg_toplevel_destroy(w.xdgToplevel)
			w.xdgToplevel = nil
		}

		if w.xdgSurface != nil {
			C.xdg_surface_destroy(w.xdgSurface)
			w.xdgSurface = nil
		}

		if w.surface != nil {
			C.wl_surface_destroy(w.surface)
			w.surface = nil
		}

		if w.handle != nil {
			w.handle.Delete()
			w.handle = nil
		}
	})
}

func (w *Window) SetTitle(title string) {
	titlePtr := C.CString(title)
	defer C.free(unsafe.Pointer(titlePtr))

	C.xdg_toplevel_set_title(w.xdgToplevel, titlePtr)
}

func (w *Window) InnerSize() dpi.PhysicalSize[uint32] {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.size.ToPhysical(w.scaleFactor)
}

func (w *Window) SetInnerSize(size dpi.Size[uint32]) {
	w.mu.Lock()
	defer w.mu.Unlock()

	scaleFactor := w.scaleFactor
	physicalSize := size.ToPhysical(scaleFactor)
	logicalSize := size.ToLogical(scaleFactor)

	width := utils.Max(1, physicalSize.Width)
	height := utils.Max(1, physicalSize.Height)

	w.size = logicalSize

	var resizedCb events.WindowResizedCallback
	if w.resizedCb != nil {
		resizedCb = w.resizedCb
	}

	w.d.scheduleCallback(func() {
		resizedCb(width, height, scaleFactor)
	})
}

func (w *Window) SetMinInnerSize(size dpi.Size[uint32]) {
	w.mu.Lock()
	scaleFactor := w.scaleFactor
	w.mu.Unlock()

	logicalSize := size.ToLogical(scaleFactor)

	C.xdg_toplevel_set_min_size(
		w.xdgToplevel,
		C.int32_t(logicalSize.Width),
		C.int32_t(logicalSize.Height),
	)
}

func (w *Window) SetMaxInnerSize(size dpi.Size[uint32]) {
	w.mu.Lock()
	scaleFactor := w.scaleFactor
	w.mu.Unlock()

	logicalSize := size.ToLogical(scaleFactor)

	C.xdg_toplevel_set_max_size(
		w.xdgToplevel,
		C.int32_t(logicalSize.Width),
		C.int32_t(logicalSize.Height),
	)
}

func (w *Window) Maximized() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.maximized
}
func (w *Window) SetMinimized() {
	C.xdg_toplevel_set_minimized(w.xdgToplevel)
}
func (w *Window) SetMaximized(maximized bool) {
	if maximized {
		C.xdg_toplevel_set_maximized(w.xdgToplevel)
	} else {
		C.xdg_toplevel_unset_maximized(w.xdgToplevel)
	}
}

func (w *Window) SetCursorIcon(icon cursors.Icon) {
	if icon == 0 {
		// 0 is internally used to hide cursor,
		// users should instead use SetCursorVisible()
		// so make this no-op
		return
	}

	var cursor *C.struct_wl_cursor
	var name string
	for _, n := range common.ToXcursorName(icon) {
		name = n
		cursor = w.d.pointer.loadCursor(n, 24, w.scaleFactor)
		if cursor != nil {
			break
		}
	}

	// couldn't find the specified cursor, so no-op
	if cursor == nil {
		return
	}

	w.d.pointer.mu.Lock()
	// if not current window, don't change cursor for pointer
	// as doing so can incorrectly set cursor for other window
	// if application has multiple windows
	if w.d.pointer.focus != w.surface {
		w.d.pointer.mu.Unlock()

		// save it so that when pointer enters this window,
		// pointer will set this cursor
		w.mu.Lock()
		w.currentCursorIcon = name
		w.mu.Unlock()
		return
	} else {
		w.d.pointer.mu.Unlock()
	}

	w.mu.Lock()
	w.currentCursorIcon = name
	w.mu.Unlock()
	w.d.pointer.setCursor(cursor, name, w.scaleFactor)
}

func (w *Window) SetCursorVisible(visible bool) {
	if visible {
		w.mu.Lock()
		if w.currentCursorIcon == "" {
			w.currentCursorIcon = w.previousCursorIcon
			currentCursor := w.currentCursorIcon
			scaleFactor := w.scaleFactor
			w.mu.Unlock()

			w.d.pointer.mu.Lock()
			if w.d.pointer.focus == w.surface {
				w.d.pointer.mu.Unlock()

				cursor := w.d.pointer.loadCursor(currentCursor, 24, scaleFactor)
				if cursor != nil {
					w.d.pointer.setCursor(cursor, currentCursor, scaleFactor)
				}
			} else {
				w.d.pointer.mu.Unlock()
			}
		} else {
			w.mu.Unlock()
		}
	} else {
		w.mu.Lock()
		if w.currentCursorIcon != "" {
			w.previousCursorIcon = w.currentCursorIcon
			w.currentCursorIcon = ""
			w.mu.Unlock()

			w.d.pointer.mu.Lock()
			if w.d.pointer.focus == w.surface {
				w.d.pointer.mu.Unlock()

				w.d.pointer.setCursor(nil, "", 0)
			} else {
				w.d.pointer.mu.Unlock()
			}
		} else {
			w.mu.Unlock()
		}
	}
}

func (w *Window) SetFullscreen(fullscreen bool) {
	if fullscreen {
		C.xdg_toplevel_set_fullscreen(w.xdgToplevel, nil)
	} else {
		C.xdg_toplevel_unset_fullscreen(w.xdgToplevel)
	}
}

func (w *Window) Fullscreen() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.fullscreen
}

func (w *Window) SetCloseRequestedCallback(cb events.WindowCloseRequestedCallback) {
	w.mu.Lock()
	w.closeRequestedCb = cb
	w.mu.Unlock()
}
func (w *Window) SetResizedCallback(cb events.WindowResizedCallback) {
	w.mu.Lock()
	w.resizedCb = cb
	w.mu.Unlock()
}
func (w *Window) SetFocusedCallback(cb events.WindowFocusedCallback) {
	w.mu.Lock()
	w.focusedCb = cb
	w.mu.Unlock()
}
func (w *Window) SetUnfocusedCallback(cb events.WindowUnfocusedCallback) {
	w.mu.Lock()
	w.unfocusedCb = cb
	w.mu.Unlock()
}
func (w *Window) SetCursorEnteredCallback(cb events.WindowCursorEnteredCallback) {
	w.mu.Lock()
	w.cursorEnteredCb = cb
	w.mu.Unlock()
}
func (w *Window) SetCursorLeftCallback(cb events.WindowCursorLeftCallback) {
	w.mu.Lock()
	w.cursorLeftCb = cb
	w.mu.Unlock()
}
func (w *Window) SetCursorMovedCallback(cb events.WindowCursorMovedCallback) {
	w.mu.Lock()
	w.cursorMovedCb = cb
	w.mu.Unlock()
}
func (w *Window) SetMouseScrollCallback(cb events.WindowMouseScrollCallback) {
	w.mu.Lock()
	w.mouseWheelCb = cb
	w.mu.Unlock()
}
func (w *Window) SetMouseInputCallback(cb events.WindowMouseInputCallback) {
	w.mu.Lock()
	w.mouseInputCb = cb
	w.mu.Unlock()
}
func (w *Window) SetTouchInputCallback(cb events.WindowTouchInputCallback) {
	// TODO:
}
func (w *Window) SetModifiersChangedCallback(cb events.WindowModifiersChangedCallback) {
	w.mu.Lock()
	w.modifiersChangedCb = cb
	w.mu.Unlock()
}
func (w *Window) SetKeyboardInputCallback(cb events.WindowKeyboardInputCallback) {
	w.mu.Lock()
	w.keyboardInputCb = cb
	w.mu.Unlock()
}
func (w *Window) SetReceivedCharacterCallback(cb events.WindowReceivedCharacterCallback) {
	w.mu.Lock()
	w.receivedCharacterCb = cb
	w.mu.Unlock()
}

// caller must lock window
func (w *Window) updateScaleFactor() {
	var scaleFactor float64 = 1
	for output := range w.outputs {
		o, ok := w.d.outputs[output]
		if ok {
			scaleFactor = math.Max(float64(o.scaleFactor), scaleFactor)
		}
	}

	if w.scaleFactor == scaleFactor {
		return
	}

	w.scaleFactor = scaleFactor
	C.wl_surface_set_buffer_scale(w.surface, C.int32_t(scaleFactor))
	C.wl_surface_commit(w.surface)

	physicalSize := w.size.ToPhysical(scaleFactor)

	var resizedCb events.WindowResizedCallback
	if w.resizedCb != nil {
		resizedCb = w.resizedCb
	}

	w.d.scheduleCallback(func() {
		if resizedCb != nil {
			resizedCb(
				physicalSize.Width,
				physicalSize.Height,
				scaleFactor,
			)
		}
	})
}

//export windowSurfaceHandleEnter
func windowSurfaceHandleEnter(data unsafe.Pointer, wl_surface *C.struct_wl_surface, output *C.struct_wl_output) {
	w, ok := (*cgo.Handle)(data).Value().(*Window)
	if !ok {
		return
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	w.outputs[output] = struct{}{}
	w.updateScaleFactor()
}

//export windowSurfaceHandleLeave
func windowSurfaceHandleLeave(data unsafe.Pointer, wl_surface *C.struct_wl_surface, output *C.struct_wl_output) {
	w, ok := (*cgo.Handle)(data).Value().(*Window)
	if !ok {
		return
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	delete(w.outputs, output)
	w.updateScaleFactor()
}

//export xdgToplevelHandleConfigure
func xdgToplevelHandleConfigure(data unsafe.Pointer, xdg_toplevel *C.struct_xdg_toplevel, width C.int32_t, height C.int32_t, states *C.struct_wl_array) {
	if width == 0 || height == 0 {
		return
	}

	w, ok := (*cgo.Handle)(data).Value().(*Window)
	if !ok {
		return
	}

	maximized := false
	fullscreen := false

	for _, state := range castWlArrayToSlice[uint32](states) {
		switch state {
		case C.XDG_TOPLEVEL_STATE_MAXIMIZED:
			maximized = true
		case C.XDG_TOPLEVEL_STATE_FULLSCREEN:
			fullscreen = true
		}
	}

	w.mu.Lock()

	w.maximized = maximized
	w.fullscreen = fullscreen

	w.size = dpi.LogicalSize[uint32]{
		Width:  uint32(width),
		Height: uint32(height),
	}

	scaleFactor := w.scaleFactor
	physicalSize := w.size.ToPhysical(scaleFactor)

	var resizedCb events.WindowResizedCallback
	if w.resizedCb != nil {
		resizedCb = w.resizedCb
	}

	w.mu.Unlock()

	w.d.scheduleCallback(func() {
		if resizedCb != nil {
			resizedCb(
				physicalSize.Width,
				physicalSize.Height,
				scaleFactor,
			)
		}
	})
}

//export xdgToplevelHandleClose
func xdgToplevelHandleClose(data unsafe.Pointer, xdg_toplevel *C.struct_xdg_toplevel) {
	w, ok := (*cgo.Handle)(data).Value().(*Window)
	if !ok {
		return
	}

	w.mu.Lock()
	var closeRequestedCb events.WindowCloseRequestedCallback
	if w.closeRequestedCb != nil {
		closeRequestedCb = w.closeRequestedCb
	}
	w.mu.Unlock()

	if closeRequestedCb != nil {
		closeRequestedCb()
	}
}
