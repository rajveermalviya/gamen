//go:build linux && !android

package xkbcommon

import (
	"errors"
	"fmt"
	"os"
	"unsafe"
)

/*

#cgo linux pkg-config: x11-xcb xcb-xkb xkbcommon xkbcommon-x11

#include <stdlib.h>

#include <X11/Xlib-xcb.h>
#include <xcb/xkb.h>

#include <xkbcommon/xkbcommon.h>
#include <xkbcommon/xkbcommon-x11.h>
#include <xkbcommon/xkbcommon-compose.h>

*/
import "C"

// https://github.com/xkbcommon/libxkbcommon/blob/master/tools/interactive-x11.c
// https://github.com/xkbcommon/libxkbcommon/blob/master/tools/interactive-wayland.c

var locale = func() string {
	locale := os.Getenv("LC_ALL")
	if locale == "" {
		locale = os.Getenv("LC_CTYPE")
	}
	if locale == "" {
		locale = os.Getenv("LANG")
	}
	if locale == "" {
		locale = "C"
	}
	return locale
}()

type Xkb struct {
	context *C.struct_xkb_context
	keymap  *C.struct_xkb_keymap
	state   *C.struct_xkb_state

	composeTable *C.struct_xkb_compose_table
	composeState *C.struct_xkb_compose_state
}

func New() (xkb *Xkb, err error) {
	context := C.xkb_context_new(C.XKB_CONTEXT_NO_FLAGS)
	if context == nil {
		return nil, errors.New("failed to create xkb context")
	}

	xkb = &Xkb{context: context}

	localStr := C.CString(locale)
	defer C.free(unsafe.Pointer(localStr))

	xkb.composeTable = C.xkb_compose_table_new_from_locale(context, localStr, C.XKB_COMPOSE_COMPILE_NO_FLAGS)
	if xkb.composeTable == nil {
		return
	}

	xkb.composeState = C.xkb_compose_state_new(xkb.composeTable, C.XKB_COMPOSE_STATE_NO_FLAGS)
	if xkb.composeState == nil {
		return
	}

	return
}

type XcbConnection = C.xcb_connection_t

func NewFromXcb(conn *XcbConnection) (xkb *Xkb, deviceId int32, firstEvent uint8, err error) {
	var firstXkbEvent C.uint8_t

	ret := C.xkb_x11_setup_xkb_extension(conn,
		C.XKB_X11_MIN_MAJOR_XKB_VERSION,
		C.XKB_X11_MIN_MINOR_XKB_VERSION,
		C.XKB_X11_SETUP_XKB_EXTENSION_NO_FLAGS,
		nil, nil, &firstXkbEvent, nil)
	if ret == 0 {
		return nil, 0, 0, errors.New("failed to setup xkb extension")
	}
	firstEvent = uint8(firstXkbEvent)

	deviceId = int32(C.xkb_x11_get_core_keyboard_device_id((*C.xcb_connection_t)(conn)))
	if deviceId == -1 {
		return nil, 0, 0, errors.New("unable to find core keyboard device")
	}

	context := C.xkb_context_new(C.XKB_CONTEXT_NO_FLAGS)
	if context == nil {
		return nil, 0, 0, errors.New("failed to create xkb context")
	}

	xkb = &Xkb{context: context}

	err = xkb.UpdateKeymap(conn, deviceId)
	if err != nil {
		C.xkb_context_unref(context)
		return nil, 0, 0, fmt.Errorf("failed to create keymap: %w", err)
	}

	err = selectXkbEvents((*C.xcb_connection_t)(conn), deviceId)
	if err != nil {
		C.xkb_context_unref(context)
		return nil, 0, 0, fmt.Errorf("failed to select xcb-xkb events: %w", err)
	}

	localStr := C.CString(locale)
	defer C.free(unsafe.Pointer(localStr))

	xkb.composeTable = C.xkb_compose_table_new_from_locale(context, localStr, C.XKB_COMPOSE_COMPILE_NO_FLAGS)
	if xkb.composeTable == nil {
		return
	}

	xkb.composeState = C.xkb_compose_state_new(xkb.composeTable, C.XKB_COMPOSE_STATE_NO_FLAGS)
	if xkb.composeState == nil {
		return
	}

	return
}

func selectXkbEvents(conn *C.xcb_connection_t, deviceID int32) error {
	requiredEvents := C.XCB_XKB_EVENT_TYPE_NEW_KEYBOARD_NOTIFY |
		C.XCB_XKB_EVENT_TYPE_MAP_NOTIFY |
		C.XCB_XKB_EVENT_TYPE_STATE_NOTIFY

	requiredNknDetails := C.XCB_XKB_NKN_DETAIL_KEYCODES

	requiredMapParts := C.XCB_XKB_MAP_PART_KEY_TYPES |
		C.XCB_XKB_MAP_PART_KEY_SYMS |
		C.XCB_XKB_MAP_PART_MODIFIER_MAP |
		C.XCB_XKB_MAP_PART_EXPLICIT_COMPONENTS |
		C.XCB_XKB_MAP_PART_KEY_ACTIONS |
		C.XCB_XKB_MAP_PART_VIRTUAL_MODS |
		C.XCB_XKB_MAP_PART_VIRTUAL_MOD_MAP

	requiredStateDetails := C.XCB_XKB_STATE_PART_MODIFIER_BASE |
		C.XCB_XKB_STATE_PART_MODIFIER_LATCH |
		C.XCB_XKB_STATE_PART_MODIFIER_LOCK |
		C.XCB_XKB_STATE_PART_GROUP_BASE |
		C.XCB_XKB_STATE_PART_GROUP_LATCH |
		C.XCB_XKB_STATE_PART_GROUP_LOCK

	details := &C.xcb_xkb_select_events_details_t{
		affectNewKeyboard:  C.uint16_t(requiredNknDetails),
		newKeyboardDetails: C.uint16_t(requiredNknDetails),
		affectState:        C.uint16_t(requiredStateDetails),
		stateDetails:       C.uint16_t(requiredStateDetails),
	}

	cookie := C.xcb_xkb_select_events_aux_checked(
		conn,
		C.xcb_xkb_device_spec_t(deviceID),
		C.uint16_t(requiredEvents),
		0,
		0,
		C.uint16_t(requiredMapParts),
		C.uint16_t(requiredMapParts),
		details,
	)

	err := C.xcb_request_check(conn, cookie)
	if err != nil {
		C.free(unsafe.Pointer(err))
		return errors.New("unable to bind events")
	}
	return nil
}

func (xkb *Xkb) KeymapFromBuffer(buf []byte) error {
	keymap := C.xkb_keymap_new_from_buffer(
		xkb.context,
		(*C.char)(unsafe.Pointer(&buf[0])),
		C.size_t(len(buf)-1),
		C.XKB_KEYMAP_FORMAT_TEXT_V1,
		C.XKB_KEYMAP_COMPILE_NO_FLAGS,
	)
	if keymap == nil {
		return errors.New("unable to create keymap from buffer")
	}

	state := C.xkb_state_new(keymap)
	if state == nil {
		C.xkb_keymap_unref(keymap)
		return errors.New("unable to create new state")
	}

	xkb.keymap = keymap
	xkb.state = state
	return nil
}

func (xkb *Xkb) UpdateKeymap(conn *XcbConnection, deviceID int32) error {
	keymap := C.xkb_x11_keymap_new_from_device(
		xkb.context,
		conn,
		C.int32_t(deviceID),
		C.XKB_KEYMAP_COMPILE_NO_FLAGS)
	if keymap == nil {
		return errors.New("unable to create keymap from device")
	}

	state := C.xkb_x11_state_new_from_device(keymap, conn, C.int32_t(deviceID))
	if state == nil {
		C.xkb_keymap_unref(keymap)
		return errors.New("unable to create state from device")
	}

	xkb.keymap = keymap
	xkb.state = state
	return nil
}

type (
	KeyCode = C.xkb_keycode_t
	KeySym  = C.xkb_keysym_t
)

func (xkb *Xkb) KeyRepeats(key KeyCode) bool {
	if xkb.keymap == nil {
		return false
	}
	return C.xkb_keymap_key_repeats(xkb.keymap, key) != 0
}

func (xkb *Xkb) GetOneSym(key KeyCode) C.xkb_keysym_t {
	return C.xkb_state_key_get_one_sym(xkb.state, key)
}

func (xkb *Xkb) GetUtf8(key KeyCode, sym KeySym) string {
	if xkb.composeState == nil {
		size := C.xkb_state_key_get_utf8(xkb.state, C.xkb_keycode_t(key), nil, 0) + 1
		if size > 1 {
			buf := (*C.char)(C.malloc(C.size_t(size)))
			defer C.free(unsafe.Pointer(buf))
			C.xkb_state_key_get_utf8(xkb.state, C.xkb_keycode_t(key), buf, C.size_t(size))
			return C.GoString(buf)
		}
		return ""
	}

	feedResult := C.xkb_compose_state_feed(xkb.composeState, C.xkb_keysym_t(sym))
	if feedResult == C.XKB_COMPOSE_FEED_ACCEPTED {
		status := C.xkb_compose_state_get_status(xkb.composeState)
		switch status {
		case C.XKB_COMPOSE_COMPOSED:
			size := C.xkb_compose_state_get_utf8(xkb.composeState, nil, 0) + 1
			if size > 1 {
				buf := (*C.char)(C.malloc(C.size_t(size)))
				defer C.free(unsafe.Pointer(buf))
				C.xkb_compose_state_get_utf8(xkb.composeState, buf, C.size_t(size))
				return C.GoString(buf)
			}
		case C.XKB_COMPOSE_NOTHING:
			size := C.xkb_state_key_get_utf8(xkb.state, C.xkb_keycode_t(key), nil, 0) + 1
			if size > 1 {
				buf := (*C.char)(C.malloc(C.size_t(size)))
				defer C.free(unsafe.Pointer(buf))
				C.xkb_state_key_get_utf8(xkb.state, C.xkb_keycode_t(key), buf, C.size_t(size))
				return C.GoString(buf)
			}
		}
	}

	return ""
}

type (
	ModMask     = C.xkb_mod_mask_t
	LayoutIndex = C.xkb_layout_index_t
)

func (xkb *Xkb) UpdateMask(
	depressed_mods ModMask,
	latched_mods ModMask,
	locked_mods ModMask,
	depressed_layout LayoutIndex,
	latched_layout LayoutIndex,
	locked_layout LayoutIndex,
) bool {
	return C.xkb_state_update_mask(xkb.state,
		depressed_mods,
		latched_mods,
		locked_mods,
		depressed_layout,
		latched_layout,
		locked_layout,
	)&C.XKB_STATE_MODS_EFFECTIVE == 0
}

var (
	XKB_MOD_NAME_SHIFT = (*C.char)(unsafe.Pointer(&[]byte("Shift\x00")[0]))
	XKB_MOD_NAME_CTRL  = (*C.char)(unsafe.Pointer(&[]byte("Control\x00")[0]))
	XKB_MOD_NAME_ALT   = (*C.char)(unsafe.Pointer(&[]byte("Mod1\x00")[0]))
	XKB_MOD_NAME_LOGO  = (*C.char)(unsafe.Pointer(&[]byte("Mod4\x00")[0]))
)

func (xkb *Xkb) ModIsShift() bool {
	return C.xkb_state_mod_name_is_active(xkb.state, XKB_MOD_NAME_SHIFT, C.XKB_STATE_MODS_EFFECTIVE) == 1
}

func (xkb *Xkb) ModIsCtrl() bool {
	return C.xkb_state_mod_name_is_active(xkb.state, XKB_MOD_NAME_CTRL, C.XKB_STATE_MODS_EFFECTIVE) == 1
}

func (xkb *Xkb) ModIsAlt() bool {
	return C.xkb_state_mod_name_is_active(xkb.state, XKB_MOD_NAME_ALT, C.XKB_STATE_MODS_EFFECTIVE) == 1
}

func (xkb *Xkb) ModIsLogo() bool {
	return C.xkb_state_mod_name_is_active(xkb.state, XKB_MOD_NAME_LOGO, C.XKB_STATE_MODS_EFFECTIVE) == 1
}

func (xkb *Xkb) Destroy() {
	if xkb.state != nil {
		C.xkb_state_unref(xkb.state)
		xkb.state = nil
	}
	if xkb.keymap != nil {
		C.xkb_keymap_unref(xkb.keymap)
		xkb.keymap = nil
	}
	if xkb.composeState != nil {
		C.xkb_compose_state_unref(xkb.composeState)
		xkb.composeState = nil
	}
	if xkb.composeTable != nil {
		C.xkb_compose_table_unref(xkb.composeTable)
		xkb.composeTable = nil
	}
	if xkb.context != nil {
		C.xkb_context_unref(xkb.context)
		xkb.context = nil
	}
}
