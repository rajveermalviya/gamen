//go:build linux && !android

package wayland

/*

#include <stdlib.h>
#include "wayland-client-protocol.h"
#include "xdg-shell-client-protocol.h"
#include "xdg-decoration-unstable-v1-client-protocol.h"

extern const struct wl_registry_listener gamen_wl_registry_listener;
extern const struct wl_output_listener gamen_wl_output_listener;
extern const struct xdg_wm_base_listener gamen_xdg_wm_base_listener;
extern const struct wl_seat_listener gamen_wl_seat_listener;
extern const struct wl_pointer_listener gamen_wl_pointer_listener;
extern const struct wl_keyboard_listener gamen_wl_keyboard_listener;
extern const struct wl_callback_listener gamen_wl_callback_listener;

*/
import "C"

import (
	"errors"
	"log"
	"runtime/cgo"
	"sync"
	"time"
	"unsafe"

	"github.com/rajveermalviya/gamen/internal/common/mathx"
	"github.com/rajveermalviya/gamen/internal/xkbcommon"
	"golang.org/x/sys/unix"
)

type Display struct {
	l *wl_library

	// handle for Display to be passed between cgo callbacks
	handle *cgo.Handle
	// we allow destroy function to be called multiple
	// times, but in reality we run it once
	destroyOnce sync.Once

	// wayland objects
	display              *C.struct_wl_display
	registry             *C.struct_wl_registry
	compositor           *C.struct_wl_compositor
	shm                  *C.struct_wl_shm
	xdgWmBase            *C.struct_xdg_wm_base
	seat                 *C.struct_wl_seat
	xdgDecorationManager *C.struct_zxdg_decoration_manager_v1

	// wayland seats
	pointer  *Pointer
	keyboard *Keyboard

	// we use xkbcommon to parse keymaps and
	// handle compose sequences
	xkb *xkbcommon.Xkb

	outputs map[*C.struct_wl_output]*Output
	windows map[*C.struct_wl_surface]*Window
}

func NewDisplay() (*Display, error) {
	l, err := open_wl_library()
	if err != nil {
		return nil, err
	}

	// connect to wayland server
	display := l.wl_display_connect(
		/* name of socket */ nil, // use default path
	)
	if display == nil {
		return nil, errors.New("failed to connect to wayland server")
	}

	d := &Display{
		l:       l,
		display: display,
		windows: make(map[*C.struct_wl_surface]*Window),
		outputs: make(map[*C.struct_wl_output]*Output),
	}
	handle := cgo.NewHandle(d)
	d.handle = &handle

	// register all interfaces
	d.registry = l.wl_display_get_registry(d.display)
	l.wl_registry_add_listener(d.registry, &C.gamen_wl_registry_listener, unsafe.Pointer(d.handle))

	// wait for interface register callbacks
	l.wl_display_roundtrip(d.display)
	// wait for initial interface events
	l.wl_display_roundtrip(d.display)

	// initialize xkbcommon
	xkb, err := xkbcommon.New()
	if err != nil {
		log.Printf("unable to inititalize xkbcommon: %v\n", err)
	}
	d.xkb = xkb

	return d, nil
}

func (d *Display) Destroy() {
	d.destroyOnce.Do(func() {
		// destroy all the windows
		for s, w := range d.windows {
			w.Destroy()
			d.windows[s] = nil
			delete(d.windows, s)
		}

		if d.keyboard != nil {
			d.keyboard.destroy()
			d.keyboard = nil
		}

		if d.pointer != nil && d.pointer.pointer != nil {
			d.pointer.destroy()
			d.pointer = nil
		}

		if d.xkb != nil {
			d.xkb.Destroy()
			d.xkb = nil
		}

		if d.seat != nil {
			d.l.wl_seat_destroy(d.seat)
			d.seat = nil
		}

		if d.xdgDecorationManager != nil {
			d.l.zxdg_decoration_manager_v1_destroy(d.xdgDecorationManager)
			d.xdgDecorationManager = nil
		}

		if d.xdgWmBase != nil {
			d.l.xdg_wm_base_destroy(d.xdgWmBase)
			d.xdgWmBase = nil
		}

		for output := range d.outputs {
			d.l.wl_output_destroy(output)
			d.outputs[output] = nil
			delete(d.outputs, output)
		}

		if d.shm != nil {
			d.l.wl_shm_destroy(d.shm)
			d.shm = nil
		}

		if d.compositor != nil {
			d.l.wl_compositor_destroy(d.compositor)
			d.compositor = nil
		}

		if d.registry != nil {
			d.l.wl_registry_destroy(d.registry)
			d.registry = nil
		}

		if d.display != nil {
			d.l.wl_display_disconnect(d.display)
			d.display = nil
		}

		if d.handle != nil {
			d.handle.Delete()
			d.handle = nil
		}

		if d.l != nil {
			d.l.close()
			d.l = nil
		}
	})
}

// wayland defers key repeat to clients
func (d *Display) handleRepeatKeyFromPoll() {
	k := d.keyboard
	if k == nil {
		return
	}

	if k.repeatKey == 0 {
		// there is no key pressed
		return
	}

	if k.haveServerRepeat && k.serverRepeatRate == 0 {
		// server prefers no repeat
		return
	}

	rate := k.serverRepeatRate
	delay := k.serverRepeatDelay
	if !k.haveServerRepeat {
		// some default values
		rate = 33
		delay = 500 * time.Millisecond
	}

	// we have to wait for 'delay' duration
	// until we can start sending key repeat events
	if time.Since(k.repeatKeyStartTime) < delay {
		return
	}

	// interval between two key repeat events
	interval := time.Second / time.Duration(rate)

	if time.Since(k.repeatKeyLastSendTime) > interval {
		// send the event as interval has passed
		k.handleKeyEvent(C.uint32_t(k.repeatKey), WL_KEYBOARD_KEY_STATE_PRESSED)
		k.repeatKeyLastSendTime = time.Now()
	}
}

func (d *Display) Poll() bool {
	d.handleRepeatKeyFromPoll()
	return d.pollAndDispatchEvents(0) != -1
}

func (d *Display) Wait() bool {
	// TODO: find a better way to do this
	//
	// we switch to Poll if a key is pressed
	// to handle key repeats
	k := d.keyboard
	if k != nil {
		if k.repeatKey != 0 {
			return d.Poll()
		}
	}

	return d.pollAndDispatchEvents(-1) != -1
}

func (d *Display) WaitTimeout(timeout time.Duration) bool {
	// TODO: find a better way to do this
	//
	// we switch to Poll if a key is pressed
	// to handle key repeats
	k := d.keyboard
	if k != nil {
		if k.repeatKey != 0 {
			return d.Poll()
		}
	}

	return d.pollAndDispatchEvents(timeout) != -1
}

// schedule's a callback to run on main eventqueue and main thread
func (d *Display) scheduleCallback(fn func()) {
	cb := d.l.wl_display_sync(d.display)
	d.setCallbackListener(cb, func() {
		d.l.wl_callback_destroy(cb)
		fn()
	})
}

func (d *Display) setCallbackListener(cb *C.struct_wl_callback, fn func()) {
	fnHandle := cgo.NewHandle(fn)
	d.l.wl_callback_add_listener(cb, &C.gamen_wl_callback_listener, unsafe.Pointer(&fnHandle))
}

//export goWlCallbackDone
func goWlCallbackDone(data unsafe.Pointer, wl_callback *C.struct_wl_callback, callback_data C.uint32_t) {
	fnHandle := (*cgo.Handle)(data)
	defer fnHandle.Delete()

	fn, ok := fnHandle.Value().(func())
	if !ok {
		return
	}

	fn()
}

//export registryHandleGlobal
func registryHandleGlobal(data unsafe.Pointer, wl_registry *C.struct_wl_registry, name C.uint32_t, iface *C.char, version C.uint32_t) {
	d, ok := (*cgo.Handle)(data).Value().(*Display)
	if !ok {
		return
	}

	switch C.GoString(iface) {
	case C.GoString(C.wl_compositor_interface.name):
		d.compositor = (*C.struct_wl_compositor)(d.l.wl_registry_bind(wl_registry, name, &C.wl_compositor_interface, mathx.Min(5, version)))

	case C.GoString(C.wl_shm_interface.name):
		d.shm = (*C.struct_wl_shm)(d.l.wl_registry_bind(wl_registry, name, &C.wl_shm_interface, mathx.Min(1, version)))

	case C.GoString(C.zxdg_decoration_manager_v1_interface.name):
		d.xdgDecorationManager = (*C.struct_zxdg_decoration_manager_v1)(d.l.wl_registry_bind(wl_registry, name, &C.zxdg_decoration_manager_v1_interface, mathx.Min(1, version)))

	case C.GoString(C.wl_output_interface.name):
		output := (*C.struct_wl_output)(d.l.wl_registry_bind(wl_registry, name, &C.wl_output_interface, mathx.Min(2, version)))
		d.outputs[output] = &Output{
			output:      output,
			name:        uint32(name),
			scaleFactor: 1,
		}
		d.l.wl_output_add_listener(output, &C.gamen_wl_output_listener, unsafe.Pointer(d.handle))

	case C.GoString(C.xdg_wm_base_interface.name):
		d.xdgWmBase = (*C.struct_xdg_wm_base)(d.l.wl_registry_bind(wl_registry, name, &C.xdg_wm_base_interface, mathx.Min(4, version)))
		d.l.xdg_wm_base_add_listener(d.xdgWmBase, &C.gamen_xdg_wm_base_listener, unsafe.Pointer(d.handle))

	case C.GoString(C.wl_seat_interface.name):
		d.seat = (*C.struct_wl_seat)(d.l.wl_registry_bind(wl_registry, name, &C.wl_seat_interface, mathx.Min(5, version)))
		d.l.wl_seat_add_listener(d.seat, &C.gamen_wl_seat_listener, unsafe.Pointer(d.handle))
	}
}

//export registryHandleGlobalRemove
func registryHandleGlobalRemove(data unsafe.Pointer, wl_registry *C.struct_wl_registry, name C.uint32_t) {
	d, ok := (*cgo.Handle)(data).Value().(*Display)
	if !ok {
		return
	}

	for _, output := range d.outputs {
		if output.name == uint32(name) {
			d.l.wl_output_destroy(output.output)
			d.outputs[output.output] = nil
			delete(d.outputs, output.output)
		}
	}
}

//export xdgWmBaseHandlePing
func xdgWmBaseHandlePing(data unsafe.Pointer, xdg_wm_base *C.struct_xdg_wm_base, serial C.uint32_t) {
	d, ok := (*cgo.Handle)(data).Value().(*Display)
	if !ok {
		return
	}

	d.l.xdg_wm_base_pong(xdg_wm_base, serial)
}

//export seatHandleCapabilities
func seatHandleCapabilities(data unsafe.Pointer, wl_seat *C.struct_wl_seat, capabilities enum_wl_seat_capability) {
	d, ok := (*cgo.Handle)(data).Value().(*Display)
	if !ok {
		return
	}

	if (capabilities&WL_SEAT_CAPABILITY_POINTER) != 0 && d.pointer == nil {
		pointer := d.l.wl_seat_get_pointer(wl_seat)
		d.pointer = &Pointer{
			d:            d,
			pointer:      pointer,
			cursorThemes: make(map[uint32]*C.struct_wl_cursor_theme),
		}

		d.l.wl_pointer_add_listener(pointer, &C.gamen_wl_pointer_listener, unsafe.Pointer(d.handle))
	} else if (capabilities&WL_SEAT_CAPABILITY_POINTER) == 0 && d.pointer != nil {
		d.pointer.destroy()
		d.pointer = nil
	}

	if (capabilities&WL_SEAT_CAPABILITY_KEYBOARD) != 0 && d.keyboard == nil {
		keyboard := d.l.wl_seat_get_keyboard(wl_seat)
		d.keyboard = &Keyboard{
			d:        d,
			keyboard: keyboard,
		}

		d.l.wl_keyboard_add_listener(keyboard, &C.gamen_wl_keyboard_listener, unsafe.Pointer(d.handle))
	} else if (capabilities&WL_SEAT_CAPABILITY_KEYBOARD) == 0 && d.keyboard != nil {
		d.keyboard.destroy()
		d.keyboard = nil
	}
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

// loose port of wl_display_dispatch to handle timeouts
// TODO: move this to C
func (d *Display) pollAndDispatchEvents(timeout time.Duration) (ret C.int) {
	if d.display == nil {
		return -1
	}

	var errno error
	if ret = d.l.wl_display_prepare_read(d.display); ret == -1 {
		ret = d.l.wl_display_dispatch_pending(d.display)
		return
	}

	for {
		ret, errno = d.l.wl_display_flush(d.display)
		if ret != -1 || errno != unix.EAGAIN {
			break
		}

		fds := []unix.PollFd{{
			Fd:     int32(d.l.wl_display_get_fd(d.display)),
			Events: unix.POLLOUT,
		}}

		if r, _ := unix.Ppoll(fds, nil, nil); r == -1 {
			d.l.wl_display_cancel_read(d.display)
			return -1
		}
	}

	if ret < 0 && errno != unix.EPIPE {
		d.l.wl_display_cancel_read(d.display)
		return -1
	}

	fds := []unix.PollFd{{
		Fd:     int32(d.l.wl_display_get_fd(d.display)),
		Events: unix.POLLIN,
	}}
	if !poll(fds, timeout) {
		d.l.wl_display_cancel_read(d.display)
		return -1
	}

	if d.l.wl_display_read_events(d.display) == -1 {
		return -1
	}

	return d.l.wl_display_dispatch_pending(d.display)
}
