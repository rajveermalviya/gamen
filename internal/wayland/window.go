//go:build linux && !android

package wayland

/*

#include <stdlib.h>
#include "wayland-client-protocol.h"
#include "xdg-shell-client-protocol.h"
#include "xdg-decoration-unstable-v1-client-protocol.h"

extern const struct wl_surface_listener gamen_wl_surface_listener;
extern const struct xdg_surface_listener gamen_xdg_surface_listener;
extern const struct xdg_toplevel_listener gamen_xdg_toplevel_listener;
extern const struct zxdg_toplevel_decoration_v1_listener gamen_zxdg_toplevel_decoration_v1_listener;

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
	"github.com/rajveermalviya/gamen/internal/common/atomicx"
	"github.com/rajveermalviya/gamen/internal/common/mathx"
	"github.com/rajveermalviya/gamen/internal/common/xcursor"
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
	scaleFactor        float64                          // shared mutex
	size               dpi.LogicalSize[uint32]          // shared mutex
	outputs            map[*C.struct_wl_output]struct{} // shared mutex
	previousCursorIcon string                           // shared mutex
	currentCursorIcon  string                           // shared mutex

	maximized  atomicx.Bool // shared atomic
	fullscreen atomicx.Bool // shared atomic

	// window callbacks
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
	w := &Window{
		d:                 d,
		outputs:           make(map[*C.struct_wl_output]struct{}),
		scaleFactor:       1,
		currentCursorIcon: "left_ptr",
		size: dpi.LogicalSize[uint32]{
			Width:  640,
			Height: 480,
		},
	}
	handle := cgo.NewHandle(w)
	w.handle = &handle

	w.surface = d.l.wl_compositor_create_surface(d.compositor)
	d.l.wl_surface_add_listener(w.surface, &C.gamen_wl_surface_listener, unsafe.Pointer(w.handle))

	w.xdgSurface = d.l.xdg_wm_base_get_xdg_surface(d.xdgWmBase, w.surface)
	d.l.xdg_surface_add_listener(w.xdgSurface, &C.gamen_xdg_surface_listener, unsafe.Pointer(w.handle))

	w.xdgToplevel = d.l.xdg_surface_get_toplevel(w.xdgSurface)
	d.l.xdg_toplevel_add_listener(w.xdgToplevel, &C.gamen_xdg_toplevel_listener, unsafe.Pointer(w.handle))

	if d.xdgDecorationManager != nil {
		w.xdgToplevelDecoration = d.l.zxdg_decoration_manager_v1_get_toplevel_decoration(d.xdgDecorationManager, w.xdgToplevel)
		d.l.zxdg_toplevel_decoration_v1_add_listener(w.xdgToplevelDecoration, &C.gamen_zxdg_toplevel_decoration_v1_listener, unsafe.Pointer(w.handle))
		d.l.zxdg_toplevel_decoration_v1_set_mode(w.xdgToplevelDecoration, ZXDG_TOPLEVEL_DECORATION_V_1_MODE_SERVER_SIDE)
	}

	d.l.wl_surface_commit(w.surface)

	d.windows[w.surface] = w
	return w, nil
}

func (w *Window) WlDisplay() unsafe.Pointer { return unsafe.Pointer(w.d.display) }
func (w *Window) WlSurface() unsafe.Pointer { return unsafe.Pointer(w.surface) }

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

		if _, ok := w.d.windows[w.surface]; ok {
			w.d.windows[w.surface] = nil
			delete(w.d.windows, w.surface)
		}

		if w.xdgToplevelDecoration != nil {
			w.d.l.zxdg_toplevel_decoration_v1_destroy(w.xdgToplevelDecoration)
			w.xdgToplevelDecoration = nil
		}

		if w.xdgToplevel != nil {
			w.d.l.xdg_toplevel_destroy(w.xdgToplevel)
			w.xdgToplevel = nil
		}

		if w.xdgSurface != nil {
			w.d.l.xdg_surface_destroy(w.xdgSurface)
			w.xdgSurface = nil
		}

		if w.surface != nil {
			w.d.l.wl_surface_destroy(w.surface)
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

	w.d.l.xdg_toplevel_set_title(w.xdgToplevel, titlePtr)
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

	width := mathx.Max(1, physicalSize.Width)
	height := mathx.Max(1, physicalSize.Height)

	w.size = logicalSize

	w.d.scheduleCallback(func() {
		if cb := w.resizedCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				cb(width, height, scaleFactor)
			}
		}
	})
}

func (w *Window) SetMinInnerSize(size dpi.Size[uint32]) {
	w.mu.Lock()
	scaleFactor := w.scaleFactor
	w.mu.Unlock()

	logicalSize := size.ToLogical(scaleFactor)

	w.d.l.xdg_toplevel_set_min_size(
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

	w.d.l.xdg_toplevel_set_max_size(
		w.xdgToplevel,
		C.int32_t(logicalSize.Width),
		C.int32_t(logicalSize.Height),
	)
}

func (w *Window) Maximized() bool {
	return w.maximized.Load()
}
func (w *Window) SetMinimized() {
	w.d.l.xdg_toplevel_set_minimized(w.xdgToplevel)
}
func (w *Window) SetMaximized(maximized bool) {
	if maximized {
		w.d.l.xdg_toplevel_set_maximized(w.xdgToplevel)
	} else {
		w.d.l.xdg_toplevel_unset_maximized(w.xdgToplevel)
	}
}

func (w *Window) SetCursorIcon(icon cursors.Icon) {
	if icon == 0 {
		// 0 is internally used to hide cursor,
		// users should instead use SetCursorVisible()
		// so make this no-op
		return
	}

	w.mu.Lock()
	scaleFactor := w.scaleFactor
	w.mu.Unlock()

	var cursor *C.struct_wl_cursor
	var name string
	for _, n := range xcursor.ToXcursorName(icon) {
		name = n
		cursor = w.d.pointer.loadCursor(n, 24, scaleFactor)
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
	scaleFactor = w.scaleFactor
	w.mu.Unlock()
	w.d.pointer.setCursor(cursor, name, scaleFactor)
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
		w.d.l.xdg_toplevel_set_fullscreen(w.xdgToplevel, nil)
	} else {
		w.d.l.xdg_toplevel_unset_fullscreen(w.xdgToplevel)
	}
}

func (w *Window) Fullscreen() bool {
	return w.fullscreen.Load()
}

func (w *Window) DragWindow() {
	w.d.pointer.mu.Lock()
	serial := w.d.pointer.serial
	w.d.pointer.mu.Unlock()

	w.d.l.xdg_toplevel_move(w.xdgToplevel, w.d.seat, C.uint32_t(serial))
}

func (w *Window) SetDecorations(decorate bool) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if decorate {
		if w.d.xdgDecorationManager != nil && w.xdgToplevelDecoration == nil {
			w.xdgToplevelDecoration = w.d.l.zxdg_decoration_manager_v1_get_toplevel_decoration(w.d.xdgDecorationManager, w.xdgToplevel)
			w.d.l.zxdg_toplevel_decoration_v1_add_listener(w.xdgToplevelDecoration, &C.gamen_zxdg_toplevel_decoration_v1_listener, unsafe.Pointer(w.handle))
			w.d.l.zxdg_toplevel_decoration_v1_set_mode(w.xdgToplevelDecoration, ZXDG_TOPLEVEL_DECORATION_V_1_MODE_SERVER_SIDE)
			w.d.l.wl_surface_commit(w.surface)
		}
	} else {
		if w.d.xdgDecorationManager != nil && w.xdgToplevelDecoration != nil {
			w.d.l.zxdg_toplevel_decoration_v1_set_mode(w.xdgToplevelDecoration, ZXDG_TOPLEVEL_DECORATION_V_1_MODE_CLIENT_SIDE)
			w.d.l.zxdg_toplevel_decoration_v1_destroy(w.xdgToplevelDecoration)
			w.xdgToplevelDecoration = nil
			w.d.l.wl_surface_commit(w.surface)
		}
	}
}

func (w *Window) Decorated() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.xdgToplevelDecoration != nil
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

func (w *Window) updateScaleFactor() {
	var scaleFactor float64 = 1

	w.mu.Lock()
	for output := range w.outputs {
		o, ok := w.d.outputs[output]
		if ok {
			scaleFactor = math.Max(float64(o.scaleFactor), scaleFactor)
		}
	}

	if w.scaleFactor == scaleFactor {
		w.mu.Unlock()
		return
	}
	w.scaleFactor = scaleFactor
	logicalSize := w.size
	physicalSize := logicalSize.ToPhysical(scaleFactor)
	w.mu.Unlock()

	w.d.l.wl_surface_set_buffer_scale(w.surface, C.int32_t(scaleFactor))
	w.d.l.wl_surface_commit(w.surface)

	w.d.scheduleCallback(func() {
		if cb := w.resizedCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				cb(
					physicalSize.Width,
					physicalSize.Height,
					scaleFactor,
				)
			}
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
	w.outputs[output] = struct{}{}
	w.mu.Unlock()

	w.updateScaleFactor()
}

//export windowSurfaceHandleLeave
func windowSurfaceHandleLeave(data unsafe.Pointer, wl_surface *C.struct_wl_surface, output *C.struct_wl_output) {
	w, ok := (*cgo.Handle)(data).Value().(*Window)
	if !ok {
		return
	}

	w.mu.Lock()
	delete(w.outputs, output)
	w.mu.Unlock()

	w.updateScaleFactor()
}

//export xdgSurfaceHandleConfigure
func xdgSurfaceHandleConfigure(data unsafe.Pointer, xdg_surface *C.struct_xdg_surface, serial C.uint32_t) {
	w, ok := (*cgo.Handle)(data).Value().(*Window)
	if !ok {
		return
	}

	w.d.l.xdg_surface_ack_configure(xdg_surface, serial)
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

	for _, state := range castWlArrayToSlice[enum_xdg_toplevel_state](states) {
		switch state {
		case XDG_TOPLEVEL_STATE_MAXIMIZED:
			maximized = true
		case XDG_TOPLEVEL_STATE_FULLSCREEN:
			fullscreen = true
		}
	}

	logicalSize := dpi.LogicalSize[uint32]{
		Width:  uint32(width),
		Height: uint32(height),
	}

	w.maximized.Store(maximized)
	w.fullscreen.Store(fullscreen)

	w.mu.Lock()
	w.size = logicalSize
	scaleFactor := w.scaleFactor
	w.mu.Unlock()

	physicalSize := logicalSize.ToPhysical(scaleFactor)

	if cb := w.resizedCb.Load(); cb != nil {
		if cb := (*cb); cb != nil {
			cb(
				physicalSize.Width,
				physicalSize.Height,
				scaleFactor,
			)
		}
	}
}

//export xdgToplevelHandleClose
func xdgToplevelHandleClose(data unsafe.Pointer, xdg_toplevel *C.struct_xdg_toplevel) {
	w, ok := (*cgo.Handle)(data).Value().(*Window)
	if !ok {
		return
	}

	if cb := w.closeRequestedCb.Load(); cb != nil {
		if cb := (*cb); cb != nil {
			cb()
		}
	}
}
