//go:build linux && !android

package wayland

/*

#cgo linux pkg-config: wayland-client wayland-cursor

#include <stdlib.h>
#include <wayland-client.h>
#include <wayland-cursor.h>
#include "xdg-shell-client-protocol.h"
#include "xdg-decoration-unstable-v1-client-protocol.h"

extern const struct wl_registry_listener wl_registry_listener;
extern const struct wl_output_listener wl_output_listener;
extern const struct xdg_wm_base_listener xdg_wm_base_listener;
extern const struct wl_seat_listener wl_seat_listener;
extern const struct wl_pointer_listener wl_pointer_listener;
extern const struct wl_keyboard_listener wl_keyboard_listener;
extern const struct wl_callback_listener go_wl_callback_listener;

*/
import "C"

import (
	"errors"
	"log"
	"runtime/cgo"
	"sync"
	"time"
	"unsafe"

	"github.com/rajveermalviya/gamen/internal/xkbcommon"
	"golang.org/x/sys/unix"
)

type Display struct {
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
	// connect to wayland server
	display := C.wl_display_connect(
		/* name of socket */ nil, // use default path
	)
	if display == nil {
		return nil, errors.New("failed to connect to wayland server")
	}

	d := &Display{
		display: display,
		windows: make(map[*C.struct_wl_surface]*Window),
		outputs: make(map[*C.struct_wl_output]*Output),
	}
	handle := cgo.NewHandle(d)
	d.handle = &handle

	// register all interfaces
	d.registry = C.wl_display_get_registry(d.display)
	C.wl_registry_add_listener(d.registry, &C.wl_registry_listener, unsafe.Pointer(d.handle))

	// wait for interface register callbacks
	C.wl_display_roundtrip(d.display)
	// wait for initial interface events
	C.wl_display_roundtrip(d.display)

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
			C.wl_seat_destroy(d.seat)
			d.seat = nil
		}

		if d.xdgDecorationManager != nil {
			C.zxdg_decoration_manager_v1_destroy(d.xdgDecorationManager)
			d.xdgDecorationManager = nil
		}

		if d.xdgWmBase != nil {
			C.xdg_wm_base_destroy(d.xdgWmBase)
			d.xdgWmBase = nil
		}

		for output := range d.outputs {
			C.wl_output_destroy(output)
			d.outputs[output] = nil
			delete(d.outputs, output)
		}

		if d.shm != nil {
			C.wl_shm_destroy(d.shm)
			d.shm = nil
		}

		if d.compositor != nil {
			C.wl_compositor_destroy(d.compositor)
			d.compositor = nil
		}

		if d.registry != nil {
			C.wl_registry_destroy(d.registry)
			d.registry = nil
		}

		if d.display != nil {
			C.wl_display_disconnect(d.display)
			d.display = nil
		}

		if d.handle != nil {
			d.handle.Delete()
			d.handle = nil
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
		k.handleKeyEvent(C.uint32_t(k.repeatKey), C.WL_KEYBOARD_KEY_STATE_PRESSED)
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
	cb := C.wl_display_sync(d.display)

	fnHandle := cgo.NewHandle(fn)
	C.wl_callback_add_listener(cb, &C.go_wl_callback_listener, unsafe.Pointer(&fnHandle))
}

//export goWlCallbackDone
func goWlCallbackDone(data unsafe.Pointer, wl_callback *C.struct_wl_callback, callback_data C.uint32_t) {
	defer C.wl_callback_destroy(wl_callback)

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
		d.compositor = (*C.struct_wl_compositor)(C.wl_registry_bind(wl_registry, name, &C.wl_compositor_interface, version))

	case C.GoString(C.wl_shm_interface.name):
		d.shm = (*C.struct_wl_shm)(C.wl_registry_bind(wl_registry, name, &C.wl_shm_interface, version))

	case C.GoString(C.zxdg_decoration_manager_v1_interface.name):
		d.xdgDecorationManager = (*C.struct_zxdg_decoration_manager_v1)(C.wl_registry_bind(wl_registry, name, &C.zxdg_decoration_manager_v1_interface, version))

	case C.GoString(C.wl_output_interface.name):
		output := (*C.struct_wl_output)(C.wl_registry_bind(wl_registry, name, &C.wl_output_interface, version))
		d.outputs[output] = &Output{
			output:      output,
			name:        uint32(name),
			scaleFactor: 1,
		}
		C.wl_output_add_listener(output, &C.wl_output_listener, unsafe.Pointer(d.handle))

	case C.GoString(C.xdg_wm_base_interface.name):
		d.xdgWmBase = (*C.struct_xdg_wm_base)(C.wl_registry_bind(wl_registry, name, &C.xdg_wm_base_interface, version))
		C.xdg_wm_base_add_listener(d.xdgWmBase, &C.xdg_wm_base_listener, nil)

	case C.GoString(C.wl_seat_interface.name):
		d.seat = (*C.struct_wl_seat)(C.wl_registry_bind(wl_registry, name, &C.wl_seat_interface, version))
		C.wl_seat_add_listener(d.seat, &C.wl_seat_listener, unsafe.Pointer(d.handle))
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
			C.wl_output_destroy(output.output)
			d.outputs[output.output] = nil
			delete(d.outputs, output.output)
		}
	}
}

//export seatHandleCapabilities
func seatHandleCapabilities(data unsafe.Pointer, wl_seat *C.struct_wl_seat, capabilities C.enum_wl_seat_capability) {
	d, ok := (*cgo.Handle)(data).Value().(*Display)
	if !ok {
		return
	}

	if (capabilities&C.WL_SEAT_CAPABILITY_POINTER) != 0 && d.pointer == nil {
		pointer := C.wl_seat_get_pointer(wl_seat)
		d.pointer = &Pointer{
			d:            d,
			pointer:      pointer,
			cursorThemes: make(map[uint32]*C.struct_wl_cursor_theme),
		}

		C.wl_pointer_add_listener(pointer, &C.wl_pointer_listener, unsafe.Pointer(d.handle))
	} else if (capabilities&C.WL_SEAT_CAPABILITY_POINTER) == 0 && d.pointer != nil {
		d.pointer.destroy()
		d.pointer = nil
	}

	if (capabilities&C.WL_SEAT_CAPABILITY_KEYBOARD) != 0 && d.keyboard == nil {
		keyboard := C.wl_seat_get_keyboard(wl_seat)
		d.keyboard = &Keyboard{
			d:        d,
			keyboard: keyboard,
		}

		C.wl_keyboard_add_listener(keyboard, &C.wl_keyboard_listener, unsafe.Pointer(d.handle))
	} else if (capabilities&C.WL_SEAT_CAPABILITY_KEYBOARD) == 0 && d.keyboard != nil {
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
	if ret = C.wl_display_prepare_read(d.display); ret == -1 {
		ret = C.wl_display_dispatch_pending(d.display)
		return
	}

	for {
		ret, errno = C.wl_display_flush(d.display)
		if ret != -1 || errno != unix.EAGAIN {
			break
		}

		fds := []unix.PollFd{{
			Fd:     int32(C.wl_display_get_fd(d.display)),
			Events: unix.POLLOUT,
		}}

		if r, _ := unix.Ppoll(fds, nil, nil); r == -1 {
			C.wl_display_cancel_read(d.display)
			return -1
		}
	}

	if ret < 0 && errno != unix.EPIPE {
		C.wl_display_cancel_read(d.display)
		return -1
	}

	fds := []unix.PollFd{{
		Fd:     int32(C.wl_display_get_fd(d.display)),
		Events: unix.POLLIN,
	}}
	if !poll(fds, timeout) {
		C.wl_display_cancel_read(d.display)
		return -1
	}

	if C.wl_display_read_events(d.display) == -1 {
		return -1
	}

	return C.wl_display_dispatch_pending(d.display)
}
