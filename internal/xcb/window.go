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
	"github.com/rajveermalviya/gamen/internal/utils"
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
	mouseWheelCb        events.WindowMouseWheelCallback
	mouseInputCb        events.WindowMouseInputCallback
	modifiersChangedCb  events.WindowModifiersChangedCallback
	keyboardInputCb     events.WindowKeyboardInputCallback
	receivedCharacterCb events.WindowReceivedCharacterCallback
}

func NewWindow(d *Display) (*Window, error) {
	win := C.xcb_generate_id(d.xcbConn)

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

		cookie := C.xcb_create_window_checked(d.xcbConn,
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
		err := C.xcb_request_check(d.xcbConn, cookie)
		if err != nil {
			defer C.free(unsafe.Pointer(err))
			return nil, errors.New("unable to create window")
		}
	}

	// opt into window close event
	{
		cookie := C.xcb_change_property_checked(
			d.xcbConn,
			C.XCB_PROP_MODE_REPLACE,
			win,
			d.wmProtocols,
			C.XCB_ATOM_ATOM,
			32,
			1,
			unsafe.Pointer(&d.wmDeleteWindow),
		)
		err := C.xcb_request_check(d.xcbConn, cookie)
		if err != nil {
			defer C.free(unsafe.Pointer(err))
			return nil, errors.New("unable to set title")
		}
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

		cookie := C.xcb_input_xi_select_events_checked(d.xcbConn, win, 1, &mask.head)
		err := C.xcb_request_check(d.xcbConn, cookie)
		if err != nil {
			defer C.free(unsafe.Pointer(err))
			return nil, errors.New("unable to select xinput2 events")
		}
	}

	// map window
	{
		cookie := C.xcb_map_window_checked(d.xcbConn, win)
		err := C.xcb_request_check(d.xcbConn, cookie)
		if err != nil {
			defer C.free(unsafe.Pointer(err))
			return nil, errors.New("unable to map window")
		}
	}

	w := &Window{
		d:                 d,
		win:               win,
		currentCursorIcon: cursors.Default,
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

	C.xcb_change_property(
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

		C.xcb_destroy_window(w.d.xcbConn, w.win)
		C.xcb_flush(w.d.xcbConn)
	})
}

func (w *Window) InnerSize() dpi.PhysicalSize[uint32] {
	r := C.xcb_get_geometry_reply(w.d.xcbConn, C.xcb_get_geometry(w.d.xcbConn, w.win), nil)
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
		utils.Max(1, utils.Min(physicalSize.Width, math.MaxInt16)),
		utils.Max(1, utils.Min(physicalSize.Height, math.MaxInt16)),
	}

	C.xcb_configure_window(w.d.xcbConn, w.win, mask, unsafe.Pointer(&values))
	C.xcb_flush(w.d.xcbConn)
}

func (w *Window) SetMinInnerSize(size dpi.Size[uint32]) {
	physicalSize := size.ToPhysical(1)

	var hints C.xcb_size_hints_t

	C.xcb_icccm_get_wm_size_hints_reply(
		w.d.xcbConn,
		C.xcb_icccm_get_wm_normal_hints_unchecked(w.d.xcbConn, w.win),
		&hints,
		nil,
	)

	C.xcb_icccm_size_hints_set_min_size(
		&hints,
		C.int32_t(utils.Min(physicalSize.Width, math.MaxInt16)),
		C.int32_t(utils.Min(physicalSize.Height, math.MaxInt16)),
	)

	C.xcb_icccm_set_wm_normal_hints(w.d.xcbConn, w.win, &hints)
}

func (w *Window) SetMaxInnerSize(size dpi.Size[uint32]) {
	physicalSize := size.ToPhysical(1)

	var hints C.xcb_size_hints_t

	C.xcb_icccm_get_wm_size_hints_reply(
		w.d.xcbConn,
		C.xcb_icccm_get_wm_normal_hints_unchecked(w.d.xcbConn, w.win),
		&hints,
		nil,
	)

	C.xcb_icccm_size_hints_set_max_size(
		&hints,
		C.int32_t(utils.Min(physicalSize.Width, math.MaxInt16)),
		C.int32_t(utils.Min(physicalSize.Height, math.MaxInt16)),
	)

	C.xcb_icccm_set_wm_normal_hints(w.d.xcbConn, w.win, &hints)
}

// func (w *Window) Minimized() bool {
// 	r := C.xcb_get_property_reply(
// 		w.d.xcbConn,
// 		C.xcb_get_property(
// 			w.d.xcbConn,
// 			0,
// 			w.win,
// 			w.d.wmState,
// 			C.XCB_ATOM_ANY,
// 			0,
// 			1,
// 		),
// 		nil,
// 	)
// 	defer C.free(unsafe.Pointer(r))

// 	dataSlice := unsafe.Slice(
// 		(*uint32)(C.xcb_get_property_value(r)),
// 		uintptr(C.xcb_get_property_value_length(r))/unsafe.Sizeof(uint32(0)),
// 	)
// 	return dataSlice[0] == C.XCB_ICCCM_WM_STATE_ICONIC
// }

func (w *Window) Maximized() bool {
	r := C.xcb_get_property_reply(
		w.d.xcbConn,
		C.xcb_get_property(
			w.d.xcbConn,
			0,
			w.win,
			w.d.netWmState,
			C.XCB_ATOM_ATOM,
			0,
			1024,
		),
		nil,
	)
	defer C.free(unsafe.Pointer(r))

	dataSlice := unsafe.Slice(
		(*C.xcb_atom_t)(C.xcb_get_property_value(r)),
		uintptr(C.xcb_get_property_value_length(r))/unsafe.Sizeof(C.xcb_atom_t(0)),
	)

	var hasMaximizedHorz, hasMaximizedVert bool
	for _, atom := range dataSlice {
		if !hasMaximizedHorz && atom == w.d.netWmStateMaximizedHorz {
			hasMaximizedHorz = true
		}
		if !hasMaximizedVert && atom == w.d.netWmStateMaximizedVert {
			hasMaximizedVert = true
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

	C.xcb_send_event(
		w.d.xcbConn,
		0,
		w.d.screens[0].xcbScreen.root,
		C.XCB_EVENT_MASK_STRUCTURE_NOTIFY|C.XCB_EVENT_MASK_SUBSTRUCTURE_REDIRECT,
		(*C.char)(unsafe.Pointer(&event)),
	)

	var hints C.xcb_icccm_wm_hints_t
	if C.xcb_icccm_get_wm_hints_reply(
		w.d.xcbConn,
		C.xcb_icccm_get_wm_hints(w.d.xcbConn, w.win),
		&hints,
		nil,
	) != 0 {
		C.xcb_icccm_wm_hints_set_iconic(&hints)
		C.xcb_icccm_set_wm_hints(w.d.xcbConn, w.win, &hints)
	}
	C.xcb_flush(w.d.xcbConn)
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

	C.xcb_send_event(
		w.d.xcbConn,
		0,
		w.d.screens[0].xcbScreen.root,
		C.XCB_EVENT_MASK_STRUCTURE_NOTIFY|C.XCB_EVENT_MASK_SUBSTRUCTURE_REDIRECT,
		(*C.char)(unsafe.Pointer(&event)),
	)
	C.xcb_flush(w.d.xcbConn)
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
	C.xcb_change_window_attributes(
		w.d.xcbConn,
		w.win,
		C.XCB_CW_CURSOR,
		unsafe.Pointer(&cursor),
	)
	C.xcb_flush(w.d.xcbConn)
}

func (w *Window) SetCursorVisible(visible bool) {
	if visible {
		w.mu.Lock()
		if w.currentCursorIcon == 0 {
			w.currentCursorIcon = w.previousCursorIcon
			currentCursor := w.currentCursorIcon
			w.mu.Unlock()

			cursor := w.d.loadCursorIcon(currentCursor)
			C.xcb_change_window_attributes(
				w.d.xcbConn,
				w.win,
				C.XCB_CW_CURSOR,
				unsafe.Pointer(&cursor),
			)
			C.xcb_flush(w.d.xcbConn)
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
			C.xcb_change_window_attributes(
				w.d.xcbConn,
				w.win,
				C.XCB_CW_CURSOR,
				unsafe.Pointer(&cursor),
			)
			C.xcb_flush(w.d.xcbConn)
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

	C.xcb_send_event(
		w.d.xcbConn,
		0,
		w.d.screens[0].xcbScreen.root,
		C.XCB_EVENT_MASK_STRUCTURE_NOTIFY|C.XCB_EVENT_MASK_SUBSTRUCTURE_REDIRECT,
		(*C.char)(unsafe.Pointer(&event)),
	)
	C.xcb_flush(w.d.xcbConn)
}
func (w *Window) Fullscreen() bool {
	r := C.xcb_get_property_reply(
		w.d.xcbConn,
		C.xcb_get_property(
			w.d.xcbConn,
			0,
			w.win,
			w.d.netWmState,
			C.XCB_ATOM_ATOM,
			0,
			1024,
		),
		nil,
	)
	defer C.free(unsafe.Pointer(r))

	dataSlice := unsafe.Slice(
		(*C.xcb_atom_t)(C.xcb_get_property_value(r)),
		uintptr(C.xcb_get_property_value_length(r))/unsafe.Sizeof(C.xcb_atom_t(0)),
	)

	for _, atom := range dataSlice {
		if atom == w.d.netWmStateFullscreen {
			return true
		}
	}
	return false
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
func (w *Window) SetMouseWheelCallback(cb events.WindowMouseWheelCallback) {
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
