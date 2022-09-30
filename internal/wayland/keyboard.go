//go:build linux && !android

package wayland

/*

#include <stdlib.h>
#include "wayland-client-protocol.h"

*/
import "C"

import (
	"log"
	"runtime/cgo"
	"time"
	"unsafe"

	"github.com/rajveermalviya/gamen/events"
	"github.com/rajveermalviya/gamen/internal/xkbcommon"
	"golang.org/x/sys/unix"
)

type Keyboard struct {
	d        *Display
	keyboard *C.struct_wl_keyboard

	// from repeat_info event
	haveServerRepeat  bool
	serverRepeatRate  uint32
	serverRepeatDelay time.Duration

	// some state to handle key repeats
	// key that was pressed
	repeatKey uint32
	// time when key was first pressed
	repeatKeyStartTime time.Time
	// time when we sent the last synthetic key press event
	repeatKeyLastSendTime time.Time

	// which surface has keyboard focus
	focus *C.struct_wl_surface

	modifiers events.ModifiersState
}

func (k *Keyboard) destroy() {
	if k.keyboard != nil {
		k.d.l.wl_keyboard_destroy(k.keyboard)
		k.keyboard = nil
	}
}

//export keyboardHandleKeymap
func keyboardHandleKeymap(data unsafe.Pointer, wl_keyboard *C.struct_wl_keyboard, format C.uint32_t, fd C.int32_t, size C.uint32_t) {
	defer unix.Close(int(fd))

	ev, ok := (*cgo.Handle)(data).Value().(*Display)
	if !ok {
		return
	}

	if ev.xkb == nil {
		return
	}

	buf, err := unix.Mmap(int(fd), 0, int(size), unix.PROT_READ, unix.MAP_SHARED)
	if err != nil {
		log.Printf("failed to mmap keymap: %v\n", err)
		return
	}
	defer unix.Munmap(buf)

	err = ev.xkb.KeymapFromBuffer(buf)
	if err != nil {
		log.Printf("failed to create keymap from buffer: %v\n", err)
		return
	}
}

//export keyboardHandleEnter
func keyboardHandleEnter(data unsafe.Pointer, wl_keyboard *C.struct_wl_keyboard, serial C.uint32_t, surface *C.struct_wl_surface, keys *C.struct_wl_array) {
	d, ok := (*cgo.Handle)(data).Value().(*Display)
	if !ok {
		return
	}
	k := d.keyboard

	k.focus = surface

	w, ok := d.windows[surface]
	if !ok {
		return
	}

	if cb := w.focusedCb.Load(); cb != nil {
		if cb := (*cb); cb != nil {
			cb()
		}
	}

	// send modifiers that are already pressed
	if k.modifiers != 0 {
		if cb := w.modifiersChangedCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				cb(k.modifiers)
			}
		}
	}
}

//export keyboardHandleLeave
func keyboardHandleLeave(data unsafe.Pointer, wl_keyboard *C.struct_wl_keyboard, serial C.uint32_t, surface *C.struct_wl_surface) {
	d, ok := (*cgo.Handle)(data).Value().(*Display)
	if !ok {
		return
	}
	k := d.keyboard

	w, ok := d.windows[surface]
	if !ok {
		return
	}

	// remove modifiers if pressed
	if k.modifiers != 0 {
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

	k.repeatKey = 0
	k.focus = nil
}

//export keyboardHandleKey
func keyboardHandleKey(data unsafe.Pointer, wl_keyboard *C.struct_wl_keyboard, serial C.uint32_t, _time C.uint32_t, key C.uint32_t, state enum_wl_keyboard_key_state) {
	d, ok := (*cgo.Handle)(data).Value().(*Display)
	if !ok {
		return
	}

	xkb := d.xkb
	if xkb == nil {
		return
	}

	k := d.keyboard
	if k.focus == nil {
		return
	}

	now := time.Now()
	k.repeatKeyStartTime = now
	k.repeatKeyLastSendTime = now
	if state == WL_KEYBOARD_KEY_STATE_RELEASED && k.repeatKey == uint32(key) {
		k.repeatKey = 0
	}

	k.handleKeyEvent(key, state)

	if state == WL_KEYBOARD_KEY_STATE_PRESSED && xkb.KeyRepeats(xkbcommon.KeyCode(key+8)) {
		k.repeatKey = uint32(key)
		k.repeatKeyStartTime = time.Now()
		k.repeatKeyLastSendTime = k.repeatKeyStartTime
	}
}

func (k *Keyboard) handleKeyEvent(key C.uint32_t, state enum_wl_keyboard_key_state) {
	xkb := k.d.xkb
	w, ok := k.d.windows[k.focus]
	if !ok {
		return
	}

	key += 8
	sym := xkb.GetOneSym(xkbcommon.KeyCode(key))

	if cb := w.keyboardInputCb.Load(); cb != nil {
		if cb := (*cb); cb != nil {
			var buttonState events.ButtonState
			switch state {
			case WL_KEYBOARD_KEY_STATE_PRESSED:
				buttonState = events.ButtonStatePressed
			case WL_KEYBOARD_KEY_STATE_RELEASED:
				buttonState = events.ButtonStateReleased
			}

			cb(
				buttonState,
				events.ScanCode(key),
				xkbcommon.KeySymToVirtualKey(sym),
			)
		}
	}

	if state == WL_KEYBOARD_KEY_STATE_PRESSED {
		if cb := w.receivedCharacterCb.Load(); cb != nil {
			if cb := (*cb); cb != nil {
				utf8 := xkb.GetUtf8(xkbcommon.KeyCode(key), xkbcommon.KeySym(sym))
				for _, char := range utf8 {
					cb(char)
				}
			}
		}
	}
}

//export keyboardHandleModifiers
func keyboardHandleModifiers(data unsafe.Pointer, wl_keyboard *C.struct_wl_keyboard, serial C.uint32_t, mods_depressed C.uint32_t, mods_latched C.uint32_t, mods_locked C.uint32_t, group C.uint32_t) {
	d, ok := (*cgo.Handle)(data).Value().(*Display)
	if !ok {
		return
	}
	k := d.keyboard

	if d.xkb == nil {
		return
	}

	if d.xkb.UpdateMask(
		xkbcommon.ModMask(mods_depressed),
		xkbcommon.ModMask(mods_latched),
		xkbcommon.ModMask(mods_locked),
		0,
		0,
		xkbcommon.LayoutIndex(group),
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

	k.modifiers = m

	if k.focus == nil {
		return
	}

	w, ok := d.windows[k.focus]
	if !ok {
		return
	}

	if cb := w.modifiersChangedCb.Load(); cb != nil {
		if cb := (*cb); cb != nil {
			cb(m)
		}
	}
}

//export keyboardHandleRepeatInfo
func keyboardHandleRepeatInfo(data unsafe.Pointer, wl_keyboard *C.struct_wl_keyboard, rate C.int32_t, delay C.int32_t) {
	d, ok := (*cgo.Handle)(data).Value().(*Display)
	if !ok {
		return
	}

	k := d.keyboard

	k.haveServerRepeat = true
	k.serverRepeatRate = uint32(rate)
	k.serverRepeatDelay = time.Duration(delay) * time.Millisecond
}
