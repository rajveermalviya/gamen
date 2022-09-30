//go:build linux && !android

package xcb

/*

#include <stdlib.h>

#include <X11/Xlib-xcb.h>
#include <xcb/randr.h>
#include <xcb/xinput.h>
#include <xcb/xkb.h>
#include <xcb/xcb.h>

#include <xkbcommon/xkbcommon-x11.h>

*/
import "C"

import (
	"errors"
	"log"
	"sync"
	"time"
	"unsafe"

	"github.com/rajveermalviya/gamen/cursors"
	"github.com/rajveermalviya/gamen/dpi"
	"github.com/rajveermalviya/gamen/events"
	"github.com/rajveermalviya/gamen/internal/common/atomicx"
	"github.com/rajveermalviya/gamen/internal/xkbcommon"
	"golang.org/x/sys/unix"
)

type Display struct {
	l *xcb_library

	mu sync.Mutex
	// we allow destroy function to be called multiple
	// times, but in reality we run it once
	destroyRequested atomicx.Bool
	destroyed        atomicx.Bool
	doneFirstLoop    bool

	xlibDisp *C.Display
	xcbConn  *C.struct_xcb_connection_t
	screens  []*Output

	xrandrFirstEvent C.uint8_t

	xiOpcode     C.uint8_t
	xiFirstEvent C.uint8_t

	xkbFirstEvent C.uint8_t

	xkb           *xkbcommon.Xkb
	deviceID      int32
	firstXkbEvent C.uint8_t

	// state
	cursors            map[cursors.Icon]C.xcb_cursor_t // shared mutex
	lastMousePositionX C.xcb_input_fp1616_t            // shared mutex
	lastMousePositionY C.xcb_input_fp1616_t            // shared mutex

	windows          map[C.xcb_window_t]*Window                  // non-shared
	scrollingDevices map[C.xcb_input_device_id_t]scrollingDevice // non-shared
	focus            C.xcb_window_t                              // non-shared
	modifiers        events.ModifiersState                       // non-shared

	// atoms
	wmProtocols             C.xcb_atom_t
	wmDeleteWindow          C.xcb_atom_t
	wmState                 C.xcb_atom_t
	wmChangeState           C.xcb_atom_t
	netWmState              C.xcb_atom_t
	netWmStateMaximizedHorz C.xcb_atom_t
	netWmStateMaximizedVert C.xcb_atom_t
	netWmStateFullscreen    C.xcb_atom_t
	netWmMoveResize         C.xcb_atom_t
	motifWmHints            C.xcb_atom_t
	relHorizWheel           C.xcb_atom_t
	relVertWheel            C.xcb_atom_t
	relHorizScroll          C.xcb_atom_t
	relVertScroll           C.xcb_atom_t
}

func NewDisplay() (*Display, error) {
	l, err := open_xcb_library()
	if err != nil {
		return nil, err
	}

	l.XInitThreads()
	xlibDisp := l.XOpenDisplay(nil)
	if xlibDisp == nil {
		return nil, errors.New("XOpenDisplay failed")
	}
	xcbConn := l.XGetXCBConnection(xlibDisp)
	l.XSetEventQueueOwner(xlibDisp, C.XCBOwnsEventQueue)

	d := &Display{
		l:                l,
		xlibDisp:         xlibDisp,
		xcbConn:          xcbConn,
		windows:          map[C.xcb_window_t]*Window{},
		scrollingDevices: map[C.xcb_input_device_id_t]scrollingDevice{},
		cursors:          map[cursors.Icon]C.xcb_cursor_t{},
	}

	// xcb-randr
	{
		reply := l.xcb_get_extension_data(xcbConn, (*C.xcb_extension_t)(l.xcb_randr_id))
		if reply == nil || reply.present == 0 {
			return nil, errors.New("xcb-randr not available")
		}

		query := l.xcb_randr_query_version_reply(
			xcbConn,
			C.XCB_RANDR_MAJOR_VERSION,
			C.XCB_RANDR_MINOR_VERSION,
		)
		defer C.free(unsafe.Pointer(query))
		if query == nil || (query.major_version < 1 || (query.major_version == 1 && query.minor_version < 2)) {
			return nil, errors.New("xcb-randr not available")
		}

		d.xrandrFirstEvent = reply.first_event
	}

	// xcb-xinput
	{
		reply := l.xcb_get_extension_data(xcbConn, (*C.xcb_extension_t)(l.xcb_input_id))
		if reply == nil || reply.present == 0 {
			return nil, errors.New("xcb-xinput not available")
		}

		query := l.xcb_input_xi_query_version_reply(
			xcbConn,
			C.XCB_INPUT_MAJOR_VERSION,
			C.XCB_INPUT_MINOR_VERSION,
		)
		defer C.free(unsafe.Pointer(query))
		if query == nil || query.major_version < 2 {
			return nil, errors.New("xcb-xinput not available")
		}
	}

	// xcb-xkb
	{
		reply := l.xcb_get_extension_data(xcbConn, (*C.xcb_extension_t)(l.xcb_xkb_id))
		if reply == nil || reply.present == 0 {
			return nil, errors.New("xcb-xkb not available")
		}

		query := l.xcb_xkb_use_extension_reply(
			xcbConn,
			1,
			0,
		)
		defer C.free(unsafe.Pointer(query))
		if query == nil || query.supported == 0 {
			return nil, errors.New("xcb-xkb not available")
		}
	}

	// atoms
	{
		d.wmProtocols = d.internAtom(true, "WM_PROTOCOLS")
		d.wmDeleteWindow = d.internAtom(false, "WM_DELETE_WINDOW")
		d.wmState = d.internAtom(false, "WM_STATE")
		d.wmChangeState = d.internAtom(false, "WM_CHANGE_STATE")

		d.netWmState = d.internAtom(false, "_NET_WM_STATE")
		d.netWmStateMaximizedHorz = d.internAtom(false, "_NET_WM_STATE_MAXIMIZED_HORZ")
		d.netWmStateMaximizedVert = d.internAtom(false, "_NET_WM_STATE_MAXIMIZED_VERT")
		d.netWmStateFullscreen = d.internAtom(false, "_NET_WM_STATE_FULLSCREEN")
		d.netWmMoveResize = d.internAtom(false, "_NET_WM_MOVERESIZE")

		d.motifWmHints = d.internAtom(false, "_MOTIF_WM_HINTS")

		d.relHorizWheel = d.internAtom(false, "Rel Horiz Wheel")
		d.relVertWheel = d.internAtom(false, "Rel Vert Wheel")
		d.relHorizScroll = d.internAtom(false, "Rel Horiz Scroll")
		d.relVertScroll = d.internAtom(false, "Rel Vert Scroll")
	}

	setup := l.xcb_get_setup(xcbConn)
	d.initializeOutputs(setup)

	if err := d.xiSetupScrollingDevices(C.XCB_INPUT_DEVICE_ALL); err != nil {
		return nil, err
	}

	// xi events
	{
		var eventMask struct {
			deviceid C.xcb_input_device_id_t
			mask_len C.uint16_t
			mask     [8]C.uint8_t
		}

		eventMask.deviceid = C.XCB_INPUT_DEVICE_ALL
		eventMask.mask_len = 1
		setXiMask(&eventMask.mask, C.XCB_INPUT_HIERARCHY)
		setXiMask(&eventMask.mask, C.XCB_INPUT_DEVICE_CHANGED)

		d.l.xcb_input_xi_select_events(xcbConn, 0, 1, (*C.xcb_input_event_mask_t)(unsafe.Pointer(&eventMask)))
	}

	xkb, deviceID, firstXkbEvent, err := xkbcommon.NewFromXcb((*xkbcommon.XcbConnection)(xcbConn))
	if err != nil {
		log.Printf("unable to inititalize xkbcommon: %v\n", err)
	}
	d.xkb = xkb
	d.firstXkbEvent = C.uint8_t(firstXkbEvent)
	d.deviceID = deviceID

	return d, nil
}

func (d *Display) Poll() bool {
	events := make([]*C.xcb_generic_event_t, 0, 2048)

	for {
		ev := d.l.xcb_poll_for_event(d.xcbConn)
		if ev == nil {
			break
		}

		events = append(events, ev)
	}

	for _, ev := range events {
		d.processEvent(ev)
	}

	if d.destroyRequested.Load() && !d.destroyed.Load() {
		d.destroy()
		return false
	}
	return !d.destroyed.Load()
}

func (d *Display) Wait() bool {
	if !d.doneFirstLoop {
		d.doneFirstLoop = true
		return d.Poll()
	}

	fds := []unix.PollFd{{
		Fd:     int32(d.l.xcb_get_file_descriptor(d.xcbConn)),
		Events: unix.POLLIN,
	}}
	if !poll(fds, -1) {
		return false
	}

	return d.Poll()
}

func (d *Display) WaitTimeout(timeout time.Duration) bool {
	if !d.doneFirstLoop {
		d.doneFirstLoop = true
		return d.Poll()
	}

	fds := []unix.PollFd{{
		Fd:     int32(d.l.xcb_get_file_descriptor(d.xcbConn)),
		Events: unix.POLLIN,
	}}
	if !poll(fds, timeout) {
		return false
	}

	return d.Poll()
}

// helper poll function to correctly handle
// timeouts when EINTR occurs
func poll(fds []unix.PollFd, timeout time.Duration) bool {
	switch timeout {
	case -1:
		for {
			result, errno := unix.Ppoll(fds, nil, nil)
			if result > 0 {
				return true
			} else if result == -1 && errno != unix.EINTR && errno != unix.EAGAIN {
				return false
			}
		}

	case 0:
		for {
			result, errno := unix.Ppoll(fds, &unix.Timespec{}, nil)
			if result == -1 && errno != unix.EINTR && errno != unix.EAGAIN {
				return false
			} else {
				return true
			}
		}

	default:
		for {
			start := time.Now()

			ts := unix.NsecToTimespec(int64(timeout))
			result, errno := unix.Ppoll(fds, &ts, nil)

			timeout -= time.Since(start)

			if result > 0 {
				return true
			} else if result == -1 && errno != unix.EINTR && errno != unix.EAGAIN {
				return false
			} else if timeout <= 0 {
				return true
			}
		}
	}
}

func (d *Display) Destroy() {
	d.destroyRequested.Store(true)
}

func (d *Display) destroy() {
	for _, w := range d.windows {
		w.Destroy()
	}

	{
		d.mu.Lock()
		for i, c := range d.cursors {
			d.l.xcb_free_cursor(d.xcbConn, c)
			delete(d.cursors, i)
		}
		d.mu.Unlock()
	}

	if d.xkb != nil {
		d.xkb.Destroy()
		d.xkb = nil
	}

	if d.xlibDisp != nil {
		d.l.XCloseDisplay(d.xlibDisp)
		d.xlibDisp = nil
		d.xcbConn = nil
	}

	if d.l != nil {
		d.l.close()
		d.l = nil
	}

	d.destroyed.Store(true)
}

func (d *Display) processEvent(e *C.xcb_generic_event_t) {
	defer C.free(unsafe.Pointer(e))

	switch e.response_type & ^C.uint8_t(0x80) {
	case C.XCB_CONFIGURE_NOTIFY:
		ev := (*C.xcb_configure_notify_event_t)(unsafe.Pointer(e))

		w, ok := d.windows[ev.event]
		if !ok {
			return
		}

		size := dpi.PhysicalSize[uint32]{
			Width:  uint32(ev.width),
			Height: uint32(ev.height),
		}

		w.mu.Lock()
		if w.size != size {
			w.size = size
			w.mu.Unlock()

			if cb := w.resizedCb.Load(); cb != nil {
				if cb := (*cb); cb != nil {
					cb(size.Width, size.Height, 1)
				}
			}
		} else {
			w.mu.Unlock()
		}

	case C.XCB_CLIENT_MESSAGE:
		ev := (*C.xcb_client_message_event_t)(unsafe.Pointer(e))

		w, ok := d.windows[ev.window]
		if !ok {
			return
		}

		data32 := unsafe.Slice((*C.xcb_atom_t)(unsafe.Pointer(&ev.data)), 5)
		if data32[0] == d.wmDeleteWindow {

			if cb := w.closeRequestedCb.Load(); cb != nil {
				if cb := (*cb); cb != nil {
					cb()
				}
			}
		}

	case C.XCB_KEY_PRESS:
		ev := (*C.xcb_key_press_event_t)(unsafe.Pointer(e))

		w, ok := d.windows[ev.event]
		if !ok {
			return
		}

		sym := d.xkb.GetOneSym(xkbcommon.KeyCode(ev.detail))

		if cb := w.keyboardInputCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				cb(
					events.ButtonStatePressed,
					events.ScanCode(ev.detail),
					xkbcommon.KeySymToVirtualKey(sym),
				)
			}
		}

		if cb := w.receivedCharacterCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				utf8 := d.xkb.GetUtf8(xkbcommon.KeyCode(ev.detail), xkbcommon.KeySym(sym))
				for _, char := range utf8 {
					cb(char)
				}
			}
		}

	case C.XCB_KEY_RELEASE:
		ev := (*C.xcb_key_release_event_t)(unsafe.Pointer(e))

		w, ok := d.windows[ev.event]
		if !ok {
			return
		}

		if cb := w.keyboardInputCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				sym := d.xkb.GetOneSym(xkbcommon.KeyCode(ev.detail))

				cb(
					events.ButtonStateReleased,
					events.ScanCode(ev.detail),
					xkbcommon.KeySymToVirtualKey(sym),
				)
			}
		}

	case C.XCB_GE_GENERIC:
		d.processXIEvents(e)

	case d.firstXkbEvent:
		d.processXkbEvent(e)
	}
}

func (d *Display) processXIEvents(e *C.xcb_generic_event_t) {
	ev := (*C.xcb_ge_generic_event_t)(unsafe.Pointer(e))

	switch ev.event_type {
	case C.XCB_INPUT_DEVICE_CHANGED:
		ev := (*C.xcb_input_device_changed_event_t)(unsafe.Pointer(ev))

		switch ev.reason {
		case C.XCB_INPUT_CHANGE_REASON_DEVICE_CHANGE:
			// reset all devices
			d.xiSetupScrollingDevices(ev.sourceid)

		case C.XCB_INPUT_CHANGE_REASON_SLAVE_SWITCH:
			// only reset current device
			d.resetScrollPosition(ev.sourceid)
		}

	case C.XCB_INPUT_HIERARCHY:
		ev := (*C.xcb_input_hierarchy_event_t)(unsafe.Pointer(ev))

		// ignore other events
		if ev.flags&(C.XCB_INPUT_HIERARCHY_MASK_SLAVE_REMOVED|C.XCB_INPUT_HIERARCHY_MASK_SLAVE_ADDED) == 0 {
			return
		}

		d.xiSetupScrollingDevices(C.XCB_INPUT_DEVICE_ALL)

	case C.XCB_INPUT_ENTER:
		ev := (*C.xcb_input_enter_event_t)(unsafe.Pointer(e))

		d.resetScrollPosition(ev.sourceid)

		w, ok := d.windows[ev.event]
		if !ok {
			return
		}

		if cb := w.cursorEnteredCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				cb()
			}
		}

	case C.XCB_INPUT_LEAVE:
		ev := (*C.xcb_input_leave_event_t)(unsafe.Pointer(e))

		w, ok := d.windows[ev.event]
		if !ok {
			return
		}

		if cb := w.cursorLeftCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				cb()
			}
		}

	case C.XCB_INPUT_FOCUS_IN:
		ev := (*C.xcb_input_focus_in_event_t)(unsafe.Pointer(e))

		d.focus = ev.event

		w, ok := d.windows[ev.event]
		if !ok {
			return
		}

		if cb := w.focusedCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				cb()
			}
		}

	case C.XCB_INPUT_FOCUS_OUT:
		ev := (*C.xcb_input_focus_out_event_t)(unsafe.Pointer(e))

		w, ok := d.windows[ev.event]
		if ok {
			if d.modifiers != 0 {
				if cb := w.modifiersChangedCb.Load(); cb != nil {
					if cb := (*cb); cb != nil {
						cb(0)
					}
				}
			}

			if cb := w.unfocusedCb.Load(); cb != nil {
				if cb := (*cb); cb != nil {
					cb()
				}
			}
		}

		d.focus = 0

	case C.XCB_INPUT_BUTTON_PRESS, C.XCB_INPUT_BUTTON_RELEASE:
		ev := (*C.xcb_input_button_press_event_t)(unsafe.Pointer(e))

		d.mu.Lock()
		d.lastMousePositionX = ev.event_x
		d.lastMousePositionY = ev.event_y
		d.mu.Unlock()

		// ignore emulated touch & mouse wheel events
		if ev.flags&C.XCB_INPUT_POINTER_EVENT_FLAGS_POINTER_EMULATED != 0 {
			return
		}

		w, ok := d.windows[ev.event]
		if !ok {
			return
		}

		if cb := w.mouseInputCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				var state events.ButtonState
				if ev.event_type == C.XCB_INPUT_BUTTON_PRESS {
					state = events.ButtonStatePressed
				} else {
					state = events.ButtonStateReleased
				}

				switch ev.detail {
				case C.XCB_BUTTON_INDEX_1:
					cb(state, events.MouseButtonLeft)

				case C.XCB_BUTTON_INDEX_2:
					cb(state, events.MouseButtonMiddle)

				case C.XCB_BUTTON_INDEX_3:
					cb(state, events.MouseButtonRight)

				case 4, 5, 6, 7:
					// ignore, handled below

				default:
					cb(state, events.MouseButton(ev.detail))
				}
			}
		}

		if cb := w.mouseWheelCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				switch ev.detail {
				case 4:
					cb(events.MouseScrollDeltaLine, events.MouseScrollAxisVertical, 1)
				case 5:
					cb(events.MouseScrollDeltaLine, events.MouseScrollAxisVertical, -1)
				case 6:
					cb(events.MouseScrollDeltaLine, events.MouseScrollAxisHorizontal, 1)
				case 7:
					cb(events.MouseScrollDeltaLine, events.MouseScrollAxisVertical, -1)
				}
			}
		}

	case C.XCB_INPUT_MOTION:
		ev := (*C.xcb_input_motion_event_t)(unsafe.Pointer(e))

		d.mu.Lock()
		d.lastMousePositionX = ev.event_x
		d.lastMousePositionY = ev.event_y
		d.mu.Unlock()

		w, ok := d.windows[ev.event]
		if ok {
			newCursorPos := dpi.PhysicalPosition[float64]{
				X: fixed1616ToFloat64(ev.event_x),
				Y: fixed1616ToFloat64(ev.event_y),
			}

			w.mu.Lock()
			if w.cursorPos != newCursorPos {
				w.cursorPos = newCursorPos
				w.mu.Unlock()

				if cb := w.cursorMovedCb.Load(); cb != nil {
					if cb := (*cb); cb != nil {
						cb(newCursorPos.X, newCursorPos.Y)
					}
				}
			} else {
				w.mu.Unlock()
			}
		}

		dev, ok := d.scrollingDevices[ev.sourceid]
		if !ok {
			return
		}

		maskLen := d.l.xcb_input_button_press_valuator_mask_length(ev)
		mask := unsafe.Slice(d.l.xcb_input_button_press_valuator_mask(ev), maskLen)

		axisValues := unsafe.Slice(
			d.l.xcb_input_button_press_axisvalues(ev),
			d.l.xcb_input_button_press_axisvalues_length(ev),
		)

		axisValuesIndex := 0

		for i := C.uint16_t(0); i < C.uint16_t(maskLen)*8; i++ {
			if hasXiMask(mask, i) {
				if dev.horizontalScroll.index == i {
					x := fixed3232ToFloat64(axisValues[axisValuesIndex])
					axisValuesIndex++

					delta := (x - dev.horizontalScroll.position) / dev.horizontalScroll.increment
					dev.horizontalScroll.position = x
					d.scrollingDevices[ev.sourceid] = dev

					if cb := w.mouseWheelCb.Load(); cb != nil {
						if cb := (*cb); cb != nil {
							cb(
								events.MouseScrollDeltaLine,
								events.MouseScrollAxisHorizontal,
								-float64(delta),
							)
						}
					}
				} else if dev.verticalScroll.index == i {
					x := fixed3232ToFloat64(axisValues[axisValuesIndex])
					axisValuesIndex++

					delta := (x - dev.verticalScroll.position) / dev.verticalScroll.increment
					dev.verticalScroll.position = x
					d.scrollingDevices[ev.sourceid] = dev

					if cb := w.mouseWheelCb.Load(); cb != nil {
						if cb := (*cb); cb != nil {
							cb(
								events.MouseScrollDeltaLine,
								events.MouseScrollAxisVertical,
								-float64(delta),
							)
						}
					}
				}
			}
		}
	}
}

func (d *Display) processXkbEvent(e *C.xcb_generic_event_t) {
	type xkbEventBase struct {
		response_type C.uint8_t
		xkbType       C.uint8_t
		sequence      C.uint16_t
		time          C.xcb_timestamp_t
		deviceID      C.uint8_t
		_             [3]byte
	}

	ev := (*xkbEventBase)(unsafe.Pointer(e))

	if ev.deviceID != C.uint8_t(d.deviceID) {
		return
	}

	switch ev.xkbType {
	case C.XCB_XKB_NEW_KEYBOARD_NOTIFY:
		ev := (*C.xcb_xkb_new_keyboard_notify_event_t)(unsafe.Pointer(e))

		if ev.changed&C.XCB_XKB_NKN_DETAIL_KEYCODES != 0 {
			d.xkb.UpdateKeymap((*xkbcommon.XcbConnection)(d.xcbConn), d.deviceID)
		}

	case C.XCB_XKB_MAP_NOTIFY:
		d.xkb.UpdateKeymap((*xkbcommon.XcbConnection)(d.xcbConn), d.deviceID)

	case C.XCB_XKB_STATE_NOTIFY:
		ev := (*C.xcb_xkb_state_notify_event_t)(unsafe.Pointer(e))

		if d.xkb.UpdateMask(
			xkbcommon.ModMask(ev.baseMods),
			xkbcommon.ModMask(ev.latchedMods),
			xkbcommon.ModMask(ev.lockedMods),
			xkbcommon.LayoutIndex(ev.baseGroup),
			xkbcommon.LayoutIndex(ev.latchedGroup),
			xkbcommon.LayoutIndex(ev.lockedGroup),
		) {
			return
		}

		var m events.ModifiersState

		if d.xkb.ModIsShift() {
			m |= events.ModifiersStateShift
		}
		if d.xkb.ModIsCtrl() {
			m |= events.ModifiersStateCtrl
		}
		if d.xkb.ModIsAlt() {
			m |= events.ModifiersStateAlt
		}
		if d.xkb.ModIsLogo() {
			m |= events.ModifiersStateLogo
		}

		d.modifiers = m

		if d.focus == 0 {
			return
		}

		w, ok := d.windows[d.focus]
		if !ok {
			return
		}

		if cb := w.modifiersChangedCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				cb(m)
			}
		}
	}
}
