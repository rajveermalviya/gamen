//go:build linux && !android

package xcb

/*

#include <stdlib.h>
#include <X11/Xlib-xcb.h>
#include <xcb/xinput.h>
#include <xcb/xcb_icccm.h>
#include <X11/Xcursor/Xcursor.h>

*/
import "C"

import (
	"errors"
	"math"
	"sync"
	"unsafe"

	"github.com/rajveermalviya/gamen/cursors"
	"github.com/rajveermalviya/gamen/dpi"
	"github.com/rajveermalviya/gamen/events"
	"github.com/rajveermalviya/gamen/internal/common/mathx"
)

type Window struct {
	d *Display
	// x window id
	win C.xcb_window_t
	// we allow destroy function to be called multiple
	// times, but in reality we run it once
	destroyOnce sync.Once
	mu          sync.Mutex

	// state
	cursorPos          dpi.PhysicalPosition[float64]
	size               dpi.PhysicalSize[uint32]
	previousCursorIcon cursors.Icon
	currentCursorIcon  cursors.Icon

	// callbacks
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
	win := d.l.xcb_generate_id(d.xcbConn)

	// create window
	{
		var mask C.uint32_t = C.XCB_CW_BACK_PIXMAP |
			C.XCB_CW_BORDER_PIXEL |
			C.XCB_CW_BIT_GRAVITY |
			C.XCB_CW_EVENT_MASK
		values := [...]C.uint32_t{
			C.XCB_BACK_PIXMAP_NONE,             // none background
			d.screens[0].xcbScreen.black_pixel, // black border
			C.XCB_GRAVITY_NORTH_WEST,           // shift inner window from north west
			C.XCB_EVENT_MASK_STRUCTURE_NOTIFY | // listen for resize, keypress, keyrelease
				C.XCB_EVENT_MASK_KEY_PRESS | // (we listen for other input events via xinput)
				C.XCB_EVENT_MASK_KEY_RELEASE,
		}

		cookie := d.l.xcb_create_window_checked(d.xcbConn,
			0,
			win,
			d.screens[0].xcbScreen.root,
			0, 0,
			640, 480,
			0,
			C.XCB_WINDOW_CLASS_INPUT_OUTPUT,
			d.screens[0].xcbScreen.root_visual,
			mask, unsafe.Pointer(&values),
		)
		err := d.l.xcb_request_check(d.xcbConn, cookie)
		if err != nil {
			defer C.free(unsafe.Pointer(err))
			return nil, errors.New("unable to create window")
		}
	}

	// opt into window close event
	{
		wmDeleteWindow := d.wmDeleteWindow
		d.l.xcb_change_property(
			d.xcbConn,
			C.XCB_PROP_MODE_REPLACE,
			win,
			d.wmProtocols,
			C.XCB_ATOM_ATOM,
			32,
			1,
			unsafe.Pointer(&wmDeleteWindow),
		)
	}

	// select xinput events
	{
		var mask struct {
			head C.xcb_input_event_mask_t
			mask C.xcb_input_xi_event_mask_t
		}
		mask.head.deviceid = C.XCB_INPUT_DEVICE_ALL_MASTER
		mask.head.mask_len = C.uint16_t(unsafe.Sizeof(mask.mask) / unsafe.Sizeof(C.uint32_t(0)))

		mask.mask = C.XCB_INPUT_XI_EVENT_MASK_MOTION |
			C.XCB_INPUT_XI_EVENT_MASK_ENTER |
			C.XCB_INPUT_XI_EVENT_MASK_LEAVE |
			C.XCB_INPUT_XI_EVENT_MASK_FOCUS_IN |
			C.XCB_INPUT_XI_EVENT_MASK_FOCUS_OUT |
			C.XCB_INPUT_XI_EVENT_MASK_BUTTON_PRESS |
			C.XCB_INPUT_XI_EVENT_MASK_BUTTON_RELEASE

		d.l.xcb_input_xi_select_events(d.xcbConn, win, 1, &mask.head)
	}

	w := &Window{
		d:                 d,
		win:               win,
		currentCursorIcon: cursors.Default,
	}

	w.setDecorations(d.motifWmHints, true)

	// map window
	{
		cookie := d.l.xcb_map_window_checked(d.xcbConn, win)
		err := d.l.xcb_request_check(d.xcbConn, cookie)
		if err != nil {
			defer C.free(unsafe.Pointer(err))
			return nil, errors.New("unable to map window")
		}
	}

	d.windows[win] = w
	return w, nil
}

func (w *Window) XcbConnection() unsafe.Pointer {
	return unsafe.Pointer(w.d.xcbConn)
}
func (w *Window) XcbWindow() uint32 {
	return uint32(w.win)
}

func (w *Window) XlibDisplay() unsafe.Pointer {
	return unsafe.Pointer(w.d.xlibDisp)
}
func (w *Window) XlibWindow() uint32 {
	return uint32(w.win)
}

func (w *Window) SetTitle(title string) {
	titlePtr := C.CString(title)
	defer C.free(unsafe.Pointer(titlePtr))

	w.d.l.xcb_change_property(
		w.d.xcbConn,
		C.XCB_PROP_MODE_REPLACE,
		w.win,
		C.XCB_ATOM_WM_NAME,
		C.XCB_ATOM_STRING, 8,
		C.uint32_t(len(title)+1),
		unsafe.Pointer(titlePtr),
	)
}

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

		if _, ok := w.d.windows[w.win]; ok {
			w.d.windows[w.win] = nil
			delete(w.d.windows, w.win)
		}

		w.d.l.xcb_destroy_window(w.d.xcbConn, w.win)
		w.d.l.xcb_flush(w.d.xcbConn)
	})
}

func (w *Window) InnerSize() dpi.PhysicalSize[uint32] {
	r := w.d.l.xcb_get_geometry_reply(w.d.xcbConn, w.win)
	if r == nil {
		return dpi.PhysicalSize[uint32]{}
	}

	defer C.free(unsafe.Pointer(r))
	return dpi.PhysicalSize[uint32]{
		Width:  uint32(r.width),
		Height: uint32(r.height),
	}
}

func (w *Window) SetInnerSize(size dpi.Size[uint32]) {
	physicalSize := size.ToPhysical(1)

	var mask C.uint16_t = C.XCB_CONFIG_WINDOW_WIDTH | C.XCB_CONFIG_WINDOW_HEIGHT
	values := [...]uint32{
		mathx.Max(1, mathx.Min(physicalSize.Width, math.MaxInt16)),
		mathx.Max(1, mathx.Min(physicalSize.Height, math.MaxInt16)),
	}

	w.d.l.xcb_configure_window(w.d.xcbConn, w.win, mask, unsafe.Pointer(&values))
	w.d.l.xcb_flush(w.d.xcbConn)
}

func (w *Window) SetMinInnerSize(size dpi.Size[uint32]) {
	physicalSize := size.ToPhysical(1)

	var hints C.xcb_size_hints_t
	w.d.l.xcb_icccm_get_wm_normal_hints_reply(
		w.d.xcbConn,
		w.win,
		&hints,
	)

	w.d.l.xcb_icccm_size_hints_set_min_size(
		&hints,
		C.int32_t(mathx.Min(physicalSize.Width, math.MaxInt16)),
		C.int32_t(mathx.Min(physicalSize.Height, math.MaxInt16)),
	)

	w.d.l.xcb_icccm_set_wm_normal_hints(w.d.xcbConn, w.win, &hints)
}

func (w *Window) SetMaxInnerSize(size dpi.Size[uint32]) {
	physicalSize := size.ToPhysical(1)

	var hints C.xcb_size_hints_t
	w.d.l.xcb_icccm_get_wm_normal_hints_reply(
		w.d.xcbConn,
		w.win,
		&hints,
	)

	w.d.l.xcb_icccm_size_hints_set_max_size(
		&hints,
		C.int32_t(mathx.Min(physicalSize.Width, math.MaxInt16)),
		C.int32_t(mathx.Min(physicalSize.Height, math.MaxInt16)),
	)

	w.d.l.xcb_icccm_set_wm_normal_hints(w.d.xcbConn, w.win, &hints)
}

func (w *Window) Maximized() bool {
	r := w.d.l.xcb_get_property_reply(
		w.d.xcbConn,
		0,
		w.win,
		w.d.netWmState,
		C.XCB_ATOM_ATOM,
		0,
		1024,
	)
	defer C.free(unsafe.Pointer(r))

	dataSlice := unsafe.Slice(
		(*C.xcb_atom_t)(w.d.l.xcb_get_property_value(r)),
		uintptr(w.d.l.xcb_get_property_value_length(r))/unsafe.Sizeof(C.xcb_atom_t(0)),
	)

	var hasMaximizedHorz, hasMaximizedVert bool
	for _, atom := range dataSlice {
		if !hasMaximizedHorz && atom == w.d.netWmStateMaximizedHorz {
			hasMaximizedHorz = true
		}
		if !hasMaximizedVert && atom == w.d.netWmStateMaximizedVert {
			hasMaximizedVert = true
		}
		if hasMaximizedHorz && hasMaximizedVert {
			return true
		}
	}

	return hasMaximizedHorz && hasMaximizedVert
}

func (w *Window) SetMinimized() {
	event := C.xcb_client_message_event_t{
		response_type: C.XCB_CLIENT_MESSAGE,
		format:        32,
		sequence:      0,
		window:        w.win,
		_type:         w.d.wmChangeState,
	}

	data := (*[5]uint32)(unsafe.Pointer(&event.data))
	data[0] = C.XCB_ICCCM_WM_STATE_ICONIC

	w.d.l.xcb_send_event(
		w.d.xcbConn,
		0,
		w.d.screens[0].xcbScreen.root,
		C.XCB_EVENT_MASK_STRUCTURE_NOTIFY|C.XCB_EVENT_MASK_SUBSTRUCTURE_REDIRECT,
		(*C.char)(unsafe.Pointer(&event)),
	)

	var hints C.xcb_icccm_wm_hints_t
	if w.d.l.xcb_icccm_get_wm_hints_reply(
		w.d.xcbConn,
		w.win,
		&hints,
	) != 0 {
		w.d.l.xcb_icccm_wm_hints_set_iconic(&hints)
		w.d.l.xcb_icccm_set_wm_hints(w.d.xcbConn, w.win, &hints)
	}
	w.d.l.xcb_flush(w.d.xcbConn)
}

func (w *Window) SetMaximized(maximized bool) {
	event := C.xcb_client_message_event_t{
		response_type: C.XCB_CLIENT_MESSAGE,
		format:        32,
		sequence:      0,
		window:        w.win,
		_type:         w.d.netWmState,
	}

	data := (*[5]uint32)(unsafe.Pointer(&event.data))
	if maximized {
		data[0] = 1
	} else {
		data[0] = 0
	}
	data[1] = uint32(w.d.netWmStateMaximizedHorz)
	data[2] = uint32(w.d.netWmStateMaximizedVert)

	w.d.l.xcb_send_event(
		w.d.xcbConn,
		0,
		w.d.screens[0].xcbScreen.root,
		C.XCB_EVENT_MASK_STRUCTURE_NOTIFY|C.XCB_EVENT_MASK_SUBSTRUCTURE_REDIRECT,
		(*C.char)(unsafe.Pointer(&event)),
	)
	w.d.l.xcb_flush(w.d.xcbConn)
}

func (w *Window) SetCursorIcon(icon cursors.Icon) {
	if icon == 0 {
		// 0 is internally used to hide cursor,
		// users should instead use SetCursorVisible()
		// so make this no-op
		return
	}

	w.mu.Lock()
	w.currentCursorIcon = icon
	w.mu.Unlock()

	cursor := w.d.loadCursorIcon(icon)
	w.d.l.xcb_change_window_attributes(
		w.d.xcbConn,
		w.win,
		C.XCB_CW_CURSOR,
		unsafe.Pointer(&cursor),
	)
	w.d.l.xcb_flush(w.d.xcbConn)
}

func (w *Window) SetCursorVisible(visible bool) {
	if visible {
		w.mu.Lock()
		if w.currentCursorIcon == 0 {
			w.currentCursorIcon = w.previousCursorIcon
			currentCursor := w.currentCursorIcon
			w.mu.Unlock()

			cursor := w.d.loadCursorIcon(currentCursor)
			w.d.l.xcb_change_window_attributes(
				w.d.xcbConn,
				w.win,
				C.XCB_CW_CURSOR,
				unsafe.Pointer(&cursor),
			)
			w.d.l.xcb_flush(w.d.xcbConn)
		} else {
			w.mu.Unlock()
		}
	} else {
		w.mu.Lock()
		if w.currentCursorIcon != 0 {
			w.previousCursorIcon = w.currentCursorIcon
			w.currentCursorIcon = 0
			w.mu.Unlock()

			cursor := w.d.loadCursorIcon(0)
			w.d.l.xcb_change_window_attributes(
				w.d.xcbConn,
				w.win,
				C.XCB_CW_CURSOR,
				unsafe.Pointer(&cursor),
			)
			w.d.l.xcb_flush(w.d.xcbConn)
		} else {
			w.mu.Unlock()
		}
	}
}

func (w *Window) SetFullscreen(fullscreen bool) {
	event := C.xcb_client_message_event_t{
		response_type: C.XCB_CLIENT_MESSAGE,
		format:        32,
		sequence:      0,
		window:        w.win,
		_type:         w.d.netWmState,
	}

	data := (*[5]uint32)(unsafe.Pointer(&event.data))
	if fullscreen {
		data[0] = 1
	} else {
		data[0] = 0
	}
	data[1] = uint32(w.d.netWmStateFullscreen)

	w.d.l.xcb_send_event(
		w.d.xcbConn,
		0,
		w.d.screens[0].xcbScreen.root,
		C.XCB_EVENT_MASK_STRUCTURE_NOTIFY|C.XCB_EVENT_MASK_SUBSTRUCTURE_REDIRECT,
		(*C.char)(unsafe.Pointer(&event)),
	)
	w.d.l.xcb_flush(w.d.xcbConn)
}
func (w *Window) Fullscreen() bool {
	r := w.d.l.xcb_get_property_reply(
		w.d.xcbConn,
		0,
		w.win,
		w.d.netWmState,
		C.XCB_ATOM_ATOM,
		0,
		1024,
	)
	defer C.free(unsafe.Pointer(r))

	dataSlice := unsafe.Slice(
		(*C.xcb_atom_t)(w.d.l.xcb_get_property_value(r)),
		uintptr(w.d.l.xcb_get_property_value_length(r))/unsafe.Sizeof(C.xcb_atom_t(0)),
	)

	for _, atom := range dataSlice {
		if atom == w.d.netWmStateFullscreen {
			return true
		}
	}
	return false
}

func (w *Window) DragWindow() {
	const _NET_WM_MOVERESIZE_MOVE = 8

	w.d.mu.Lock()
	mousePosX := w.d.lastMousePositionX
	mousePosY := w.d.lastMousePositionY
	w.d.mu.Unlock()

	r := w.d.l.xcb_translate_coordinates_reply(
		w.d.xcbConn,
		w.win,
		w.d.screens[0].xcbScreen.root,
		C.int16_t(fixed1616ToFloat64(mousePosX)),
		C.int16_t(fixed1616ToFloat64(mousePosY)),
	)

	var posX, posY C.int16_t
	if r != nil {
		defer C.free(unsafe.Pointer(r))

		posX = r.dst_x
		posY = r.dst_y
	}

	event := C.xcb_client_message_event_t{
		response_type: C.XCB_CLIENT_MESSAGE,
		format:        32,
		sequence:      0,
		window:        w.win,
		_type:         w.d.netWmMoveResize,
	}

	data := (*[5]uint32)(unsafe.Pointer(&event.data))
	data[0] = uint32(posX)
	data[1] = uint32(posY)
	data[2] = _NET_WM_MOVERESIZE_MOVE
	data[3] = C.XCB_BUTTON_INDEX_1

	w.d.l.xcb_ungrab_pointer(w.d.xcbConn, C.XCB_CURRENT_TIME)
	w.d.l.xcb_send_event(
		w.d.xcbConn,
		0,
		w.d.screens[0].xcbScreen.root,
		C.XCB_EVENT_MASK_STRUCTURE_NOTIFY|C.XCB_EVENT_MASK_SUBSTRUCTURE_REDIRECT,
		(*C.char)(unsafe.Pointer(&event)),
	)
}

func (w *Window) setDecorations(motifWmHints C.xcb_atom_t, decorate bool) {
	var hints struct {
		flags       uint32
		functions   uint32
		decorations uint32
		inputMode   int32
		status      uint32
	}

	hints.flags = 2 // MWM_HINTS_DECORATIONS
	if decorate {
		hints.decorations = 1
	} else {
		hints.decorations = 0
	}

	w.d.l.xcb_change_property(
		w.d.xcbConn,
		C.XCB_PROP_MODE_REPLACE,
		w.win,
		motifWmHints,
		motifWmHints,
		32,
		5,
		unsafe.Pointer(&hints),
	)
}

func (w *Window) SetDecorations(decorate bool) {
	w.setDecorations(w.d.motifWmHints, decorate)
}

func (w *Window) Decorated() bool {
	type wmHints struct {
		flags       uint32
		functions   uint32
		decorations uint32
		inputMode   int32
		status      uint32
	}

	r := w.d.l.xcb_get_property_reply(
		w.d.xcbConn,
		0,
		w.win,
		w.d.motifWmHints,
		w.d.motifWmHints,
		0,
		C.uint32_t(unsafe.Sizeof(wmHints{})),
	)
	defer C.free(unsafe.Pointer(r))

	var hints wmHints
	if w.d.l.xcb_get_property_value_length(r) == C.int(unsafe.Sizeof(wmHints{})) {
		if v := (*wmHints)(w.d.l.xcb_get_property_value(r)); v != nil {
			hints = *v
		}
	}

	if hints.decorations == 0 {
		return false
	}
	return true
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
