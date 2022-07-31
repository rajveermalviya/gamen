//go:build linux && !android

package wayland

/*

#include <stdlib.h>
#include <wayland-client.h>

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
		C.wl_keyboard_destroy(k.keyboard)
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

	w.mu.Lock()
	var focusedCb events.WindowFocusedCallback
	if w.focusedCb != nil {
		focusedCb = w.focusedCb
	}
	w.mu.Unlock()

	if focusedCb != nil {
		focusedCb()
	}

	w.mu.Lock()
	var modifiersChangedCb events.WindowModifiersChangedCallback
	if w.modifiersChangedCb != nil {
		modifiersChangedCb = w.modifiersChangedCb
	}
	w.mu.Unlock()

	// send modifiers that are already pressed
	if modifiersChangedCb != nil && k.modifiers != 0 {
		modifiersChangedCb(k.modifiers)
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

	w.mu.Lock()
	var modifiersChangedCb events.WindowModifiersChangedCallback
	if w.modifiersChangedCb != nil {
		modifiersChangedCb = w.modifiersChangedCb
	}
	w.mu.Unlock()

	// remove modifiers if pressed
	if modifiersChangedCb != nil && k.modifiers != 0 {
		modifiersChangedCb(0)
	}

	w.mu.Lock()
	var unfocusedCb events.WindowUnfocusedCallback
	if w.unfocusedCb != nil {
		unfocusedCb = w.unfocusedCb
	}
	w.mu.Unlock()

	if unfocusedCb != nil {
		unfocusedCb()
	}

	k.repeatKey = 0
	k.focus = nil
}

//export keyboardHandleKey
func keyboardHandleKey(data unsafe.Pointer, wl_keyboard *C.struct_wl_keyboard, serial C.uint32_t, _time C.uint32_t, key C.uint32_t, state C.uint32_t) {
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
	if state == C.WL_KEYBOARD_KEY_STATE_RELEASED && k.repeatKey == uint32(key) {
		k.repeatKey = 0
	}

	k.handleKeyEvent(key, state)

	if state == C.WL_KEYBOARD_KEY_STATE_PRESSED && xkb.KeyRepeats(xkbcommon.KeyCode(key+8)) {
		k.repeatKey = uint32(key)
		k.repeatKeyStartTime = time.Now()
		k.repeatKeyLastSendTime = k.repeatKeyStartTime
	}
}

func (k *Keyboard) handleKeyEvent(key C.uint32_t, state C.uint32_t) {
	xkb := k.d.xkb
	w, ok := k.d.windows[k.focus]
	if !ok {
		return
	}

	key += 8
	sym := xkb.GetOneSym(xkbcommon.KeyCode(key))

	w.mu.Lock()
	var keyboardInputCb events.WindowKeyboardInputCallback
	if w.keyboardInputCb != nil {
		keyboardInputCb = w.keyboardInputCb
	}
	w.mu.Unlock()

	if keyboardInputCb != nil {
		var buttonState events.ButtonState
		switch state {
		case C.WL_KEYBOARD_KEY_STATE_PRESSED:
			buttonState = events.ButtonStatePressed
		case C.WL_KEYBOARD_KEY_STATE_RELEASED:
			buttonState = events.ButtonStateReleased
		}

		keyboardInputCb(
			buttonState,
			events.ScanCode(key),
			xkbcommon.KeySymToVirtualKey(sym),
		)
	}

	w.mu.Lock()
	var receivedCharacterCb events.WindowReceivedCharacterCallback
	if w.receivedCharacterCb != nil {
		receivedCharacterCb = w.receivedCharacterCb
	}
	w.mu.Unlock()

	if state == C.WL_KEYBOARD_KEY_STATE_PRESSED && receivedCharacterCb != nil {
		utf8 := xkb.GetUtf8(xkbcommon.KeyCode(key), xkbcommon.KeySym(sym))
		for _, char := range utf8 {
			receivedCharacterCb(char)
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

	w.mu.Lock()
	var modifiersChangedCb events.WindowModifiersChangedCallback
	if w.modifiersChangedCb != nil {
		modifiersChangedCb = w.modifiersChangedCb
	}
	w.mu.Unlock()

	if modifiersChangedCb != nil {
		modifiersChangedCb(m)
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
