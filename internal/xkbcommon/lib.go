//go:build linux && !android

package xkbcommon

/*

#cgo CFLAGS: -I${SRCDIR}/include
#cgo LDFLAGS: -ldl

#include <dlfcn.h>

#include <X11/Xlib-xcb.h>
#include <xcb/xkb.h>

#include <xkbcommon/xkbcommon.h>
#include <xkbcommon/xkbcommon-x11.h>
#include <xkbcommon/xkbcommon-compose.h>

// libxkbcommon
struct xkb_context *gamen_xkb_context_new(void *fp, enum xkb_context_flags flags) {
	typedef struct xkb_context *(*proc_xkb_context_new)(enum xkb_context_flags flags);
	return ((proc_xkb_context_new)fp)(flags);
}
struct xkb_compose_table *gamen_xkb_compose_table_new_from_locale(void *fp, struct xkb_context *context, const char *locale, enum xkb_compose_compile_flags flags) {
	typedef struct xkb_compose_table *(*proc_xkb_compose_table_new_from_locale)(struct xkb_context *context, const char *locale, enum xkb_compose_compile_flags flags);
	return ((proc_xkb_compose_table_new_from_locale)fp)(context, locale, flags);
}
struct xkb_compose_state *gamen_xkb_compose_state_new(void *fp, struct xkb_compose_table *table, enum xkb_compose_state_flags flags) {
	typedef struct xkb_compose_state *(*proc_xkb_compose_state_new)(struct xkb_compose_table *table, enum xkb_compose_state_flags flags);
	return ((proc_xkb_compose_state_new)fp)(table, flags);
}
void gamen_xkb_state_unref(void *fp, struct xkb_state *state) {
	typedef void (*proc_xkb_state_unref)(struct xkb_state *state);
	return ((proc_xkb_state_unref)fp)(state);
}
void gamen_xkb_keymap_unref(void *fp, struct xkb_keymap *keymap) {
	typedef void (*proc_xkb_keymap_unref)(struct xkb_keymap *keymap);
	return ((proc_xkb_keymap_unref)fp)(keymap);
}
void gamen_xkb_compose_state_unref(void *fp, struct xkb_compose_state *state) {
	typedef void (*proc_xkb_compose_state_unref)(struct xkb_compose_state *state);
	return ((proc_xkb_compose_state_unref)fp)(state);
}
void gamen_xkb_compose_table_unref(void *fp, struct xkb_compose_table *table) {
	typedef void (*proc_xkb_compose_table_unref)(struct xkb_compose_table *table);
	return ((proc_xkb_compose_table_unref)fp)(table);
}
void gamen_xkb_context_unref(void *fp, struct xkb_context *context) {
	typedef void (*proc_xkb_context_unref)(struct xkb_context *context);
	return ((proc_xkb_context_unref)fp)(context);
}
struct xkb_keymap *gamen_xkb_keymap_new_from_buffer(void *fp, struct xkb_context *context, const char *buffer, size_t length, enum xkb_keymap_format format, enum xkb_keymap_compile_flags flags) {
	typedef struct xkb_keymap *(*proc_xkb_keymap_new_from_buffer)(struct xkb_context *context, const char *buffer, size_t length, enum xkb_keymap_format format, enum xkb_keymap_compile_flags flags);
	return ((proc_xkb_keymap_new_from_buffer)fp)(context, buffer, length, format, flags);
}
struct xkb_state *gamen_xkb_state_new(void *fp, struct xkb_keymap *keymap) {
	typedef struct xkb_state *(*proc_xkb_state_new)(struct xkb_keymap *keymap);
	return ((proc_xkb_state_new)fp)(keymap);
}
int gamen_xkb_keymap_key_repeats(void *fp, struct xkb_keymap *keymap, xkb_keycode_t key) {
	typedef int (*proc_xkb_keymap_key_repeats)(struct xkb_keymap *keymap, xkb_keycode_t key);
	return ((proc_xkb_keymap_key_repeats)fp)(keymap, key);
}
xkb_keysym_t gamen_xkb_state_key_get_one_sym(void *fp, struct xkb_state *state, xkb_keycode_t key) {
	typedef xkb_keysym_t (*proc_xkb_state_key_get_one_sym)(struct xkb_state *state, xkb_keycode_t key);
	return ((proc_xkb_state_key_get_one_sym)fp)(state, key);
}
int gamen_xkb_state_key_get_utf8(void *fp, struct xkb_state *state, xkb_keycode_t key, char *buffer, size_t size) {
	typedef int (*proc_xkb_state_key_get_utf8)(struct xkb_state *state, xkb_keycode_t key, char *buffer, size_t size);
	return ((proc_xkb_state_key_get_utf8)fp)(state, key, buffer, size);
}
enum xkb_compose_feed_result gamen_xkb_compose_state_feed(void *fp, struct xkb_compose_state *state, xkb_keysym_t keysym) {
	typedef enum xkb_compose_feed_result (*proc_xkb_compose_state_feed)(struct xkb_compose_state *state, xkb_keysym_t keysym);
	return ((proc_xkb_compose_state_feed)fp)(state, keysym);
}
enum xkb_compose_status gamen_xkb_compose_state_get_status(void *fp, struct xkb_compose_state *state) {
	typedef enum xkb_compose_status (*proc_xkb_compose_state_get_status)(struct xkb_compose_state *state);
	return ((proc_xkb_compose_state_get_status)fp)(state);
}
int gamen_xkb_compose_state_get_utf8(void *fp, struct xkb_compose_state *state, char *buffer, size_t size) {
	typedef int (*proc_xkb_compose_state_get_utf8)(struct xkb_compose_state *state, char *buffer, size_t size);
	return ((proc_xkb_compose_state_get_utf8)fp)(state, buffer, size);
}
enum xkb_state_component gamen_xkb_state_update_mask(void *fp, struct xkb_state *state, xkb_mod_mask_t depressed_mods, xkb_mod_mask_t latched_mods, xkb_mod_mask_t locked_mods, xkb_layout_index_t depressed_layout, xkb_layout_index_t latched_layout, xkb_layout_index_t locked_layout) {
	typedef enum xkb_state_component (*proc_xkb_state_update_mask)(struct xkb_state *state, xkb_mod_mask_t depressed_mods, xkb_mod_mask_t latched_mods, xkb_mod_mask_t locked_mods, xkb_layout_index_t depressed_layout, xkb_layout_index_t latched_layout, xkb_layout_index_t locked_layout);
	return ((proc_xkb_state_update_mask)fp)(state, depressed_mods, latched_mods, locked_mods, depressed_layout, latched_layout, locked_layout);
}
int gamen_xkb_state_mod_name_is_active(void *fp, struct xkb_state *state, const char *name, enum xkb_state_component type) {
	typedef int (*proc_xkb_state_mod_name_is_active)(struct xkb_state *state, const char *name, enum xkb_state_component type);
	return ((proc_xkb_state_mod_name_is_active)fp)(state, name, type);
}

// libxkbcommon-x11
int gamen_xkb_x11_setup_xkb_extension(void *fp, xcb_connection_t *connection, uint16_t major_xkb_version, uint16_t minor_xkb_version, enum xkb_x11_setup_xkb_extension_flags flags, uint16_t *major_xkb_version_out, uint16_t *minor_xkb_version_out, uint8_t *base_event_out, uint8_t *base_error_out) {
	typedef int (*proc_xkb_x11_setup_xkb_extension)(xcb_connection_t *connection, uint16_t major_xkb_version, uint16_t minor_xkb_version, enum xkb_x11_setup_xkb_extension_flags flags, uint16_t *major_xkb_version_out, uint16_t *minor_xkb_version_out, uint8_t *base_event_out, uint8_t *base_error_out);
	return ((proc_xkb_x11_setup_xkb_extension)fp)(connection, major_xkb_version, minor_xkb_version, flags, major_xkb_version_out, minor_xkb_version_out, base_event_out, base_error_out);
}
int32_t gamen_xkb_x11_get_core_keyboard_device_id(void *fp, xcb_connection_t *connection) {
	typedef int32_t (*proc_xkb_x11_get_core_keyboard_device_id)(xcb_connection_t *connection);
	return ((proc_xkb_x11_get_core_keyboard_device_id)fp)(connection);
}
struct xkb_keymap *gamen_xkb_x11_keymap_new_from_device(void *fp, struct xkb_context *context, xcb_connection_t *connection, int32_t device_id, enum xkb_keymap_compile_flags flags) {
	typedef struct xkb_keymap *(*proc_xkb_x11_keymap_new_from_device)(struct xkb_context *context, xcb_connection_t *connection, int32_t device_id, enum xkb_keymap_compile_flags flags);
	return ((proc_xkb_x11_keymap_new_from_device)fp)(context, connection, device_id, flags);
}
struct xkb_state *gamen_xkb_x11_state_new_from_device(void *fp, struct xkb_keymap *keymap, xcb_connection_t *connection, int32_t device_id) {
	typedef struct xkb_state *(*proc_xkb_x11_state_new_from_device)(struct xkb_keymap *keymap, xcb_connection_t *connection, int32_t device_id);
	return ((proc_xkb_x11_state_new_from_device)fp)(keymap, connection, device_id);
}

// libxcb
xcb_generic_error_t *gamen_xkbcommon_xcb_request_check(void *fp, xcb_connection_t *c, xcb_void_cookie_t cookie) {
	typedef xcb_generic_error_t *(*proc_xcb_request_check)(xcb_connection_t *c, xcb_void_cookie_t cookie);
	return ((proc_xcb_request_check)fp)(c, cookie);
}

// libxcb-xkb
xcb_void_cookie_t gamen_xcb_xkb_select_events_aux_checked(void *fp, xcb_connection_t *c, xcb_xkb_device_spec_t deviceSpec, uint16_t affectWhich, uint16_t clear, uint16_t selectAll, uint16_t affectMap, uint16_t map, const xcb_xkb_select_events_details_t *details) {
	typedef xcb_void_cookie_t (*proc_xcb_xkb_select_events_aux_checked)(xcb_connection_t *c, xcb_xkb_device_spec_t deviceSpec, uint16_t affectWhich, uint16_t clear, uint16_t selectAll, uint16_t affectMap, uint16_t map, const xcb_xkb_select_events_details_t *details);
	return ((proc_xcb_xkb_select_events_aux_checked)fp)(c, deviceSpec, affectWhich, clear, selectAll, affectMap, map, details);
}

*/
import "C"
import (
	"errors"
	"unsafe"
)

type xkbcommon_library struct {
	libxkbcommonHandle    unsafe.Pointer
	libxkbcommonX11Handle unsafe.Pointer
	libxcbHandle          unsafe.Pointer
	libxcbXkbHandle       unsafe.Pointer

	// libxkbcommon
	xkb_context_new_handle                   unsafe.Pointer
	xkb_compose_table_new_from_locale_handle unsafe.Pointer
	xkb_compose_state_new_handle             unsafe.Pointer
	xkb_state_unref_handle                   unsafe.Pointer
	xkb_keymap_unref_handle                  unsafe.Pointer
	xkb_compose_state_unref_handle           unsafe.Pointer
	xkb_compose_table_unref_handle           unsafe.Pointer
	xkb_context_unref_handle                 unsafe.Pointer
	xkb_keymap_new_from_buffer_handle        unsafe.Pointer
	xkb_state_new_handle                     unsafe.Pointer
	xkb_keymap_key_repeats_handle            unsafe.Pointer
	xkb_state_key_get_one_sym_handle         unsafe.Pointer
	xkb_state_key_get_utf8_handle            unsafe.Pointer
	xkb_compose_state_feed_handle            unsafe.Pointer
	xkb_compose_state_get_status_handle      unsafe.Pointer
	xkb_compose_state_get_utf8_handle        unsafe.Pointer
	xkb_state_update_mask_handle             unsafe.Pointer
	xkb_state_mod_name_is_active_handle      unsafe.Pointer

	// libxkbcommon-x11
	xkb_x11_setup_xkb_extension_handle         unsafe.Pointer
	xkb_x11_get_core_keyboard_device_id_handle unsafe.Pointer
	xkb_x11_keymap_new_from_device_handle      unsafe.Pointer
	xkb_x11_state_new_from_device_handle       unsafe.Pointer

	// libxcb
	xcb_request_check_handle unsafe.Pointer

	// libxcb-xkb
	xcb_xkb_select_events_aux_checked_handle unsafe.Pointer
}

func loadSym(handle unsafe.Pointer, symbol *C.char) (unsafe.Pointer, error) {
	C.dlerror()
	fp := C.dlsym(handle, symbol)
	if fp == nil {
		err := C.dlerror()
		if err != nil {
			return nil, errors.New(C.GoString(err))
		}
	}
	return fp, nil
}

func open_xkbcommon_library() (*xkbcommon_library, error) {
	C.dlerror()

	l := &xkbcommon_library{}

	l.libxkbcommonHandle = C.dlopen((*C.char)(unsafe.Pointer(&([]byte("libxkbcommon.so.0\x00"))[0])), C.RTLD_LAZY)
	if l.libxkbcommonHandle == nil {
		err := C.dlerror()
		if err != nil {
			l.close()
			return nil, errors.New(C.GoString(err))
		}
	}
	l.libxkbcommonX11Handle = C.dlopen((*C.char)(unsafe.Pointer(&([]byte("libxkbcommon-x11.so.0\x00"))[0])), C.RTLD_LAZY)
	if l.libxkbcommonX11Handle == nil {
		err := C.dlerror()
		if err != nil {
			l.close()
			return nil, errors.New(C.GoString(err))
		}
	}
	l.libxcbHandle = C.dlopen((*C.char)(unsafe.Pointer(&([]byte("libxcb.so.1\x00"))[0])), C.RTLD_LAZY)
	if l.libxcbHandle == nil {
		err := C.dlerror()
		if err != nil {
			l.close()
			return nil, errors.New(C.GoString(err))
		}
	}
	l.libxcbXkbHandle = C.dlopen((*C.char)(unsafe.Pointer(&([]byte("libxcb-xkb.so.1\x00"))[0])), C.RTLD_LAZY)
	if l.libxcbXkbHandle == nil {
		err := C.dlerror()
		if err != nil {
			l.close()
			return nil, errors.New(C.GoString(err))
		}
	}

	var err error

	// libxkbcommon
	l.xkb_context_new_handle, err = loadSym(l.libxkbcommonHandle, (*C.char)(unsafe.Pointer(&([]byte("xkb_context_new\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xkb_compose_table_new_from_locale_handle, err = loadSym(l.libxkbcommonHandle, (*C.char)(unsafe.Pointer(&([]byte("xkb_compose_table_new_from_locale\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xkb_compose_state_new_handle, err = loadSym(l.libxkbcommonHandle, (*C.char)(unsafe.Pointer(&([]byte("xkb_compose_state_new\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xkb_state_unref_handle, err = loadSym(l.libxkbcommonHandle, (*C.char)(unsafe.Pointer(&([]byte("xkb_state_unref\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xkb_keymap_unref_handle, err = loadSym(l.libxkbcommonHandle, (*C.char)(unsafe.Pointer(&([]byte("xkb_keymap_unref\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xkb_compose_state_unref_handle, err = loadSym(l.libxkbcommonHandle, (*C.char)(unsafe.Pointer(&([]byte("xkb_compose_state_unref\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xkb_compose_table_unref_handle, err = loadSym(l.libxkbcommonHandle, (*C.char)(unsafe.Pointer(&([]byte("xkb_compose_table_unref\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xkb_context_unref_handle, err = loadSym(l.libxkbcommonHandle, (*C.char)(unsafe.Pointer(&([]byte("xkb_context_unref\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xkb_keymap_new_from_buffer_handle, err = loadSym(l.libxkbcommonHandle, (*C.char)(unsafe.Pointer(&([]byte("xkb_keymap_new_from_buffer\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xkb_state_new_handle, err = loadSym(l.libxkbcommonHandle, (*C.char)(unsafe.Pointer(&([]byte("xkb_state_new\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xkb_keymap_key_repeats_handle, err = loadSym(l.libxkbcommonHandle, (*C.char)(unsafe.Pointer(&([]byte("xkb_keymap_key_repeats\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xkb_state_key_get_one_sym_handle, err = loadSym(l.libxkbcommonHandle, (*C.char)(unsafe.Pointer(&([]byte("xkb_state_key_get_one_sym\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xkb_state_key_get_utf8_handle, err = loadSym(l.libxkbcommonHandle, (*C.char)(unsafe.Pointer(&([]byte("xkb_state_key_get_utf8\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xkb_compose_state_feed_handle, err = loadSym(l.libxkbcommonHandle, (*C.char)(unsafe.Pointer(&([]byte("xkb_compose_state_feed\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xkb_compose_state_get_status_handle, err = loadSym(l.libxkbcommonHandle, (*C.char)(unsafe.Pointer(&([]byte("xkb_compose_state_get_status\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xkb_compose_state_get_utf8_handle, err = loadSym(l.libxkbcommonHandle, (*C.char)(unsafe.Pointer(&([]byte("xkb_compose_state_get_utf8\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xkb_state_update_mask_handle, err = loadSym(l.libxkbcommonHandle, (*C.char)(unsafe.Pointer(&([]byte("xkb_state_update_mask\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xkb_state_mod_name_is_active_handle, err = loadSym(l.libxkbcommonHandle, (*C.char)(unsafe.Pointer(&([]byte("xkb_state_mod_name_is_active\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}

	// libxkbcommon-x11
	l.xkb_x11_setup_xkb_extension_handle, err = loadSym(l.libxkbcommonX11Handle, (*C.char)(unsafe.Pointer(&([]byte("xkb_x11_setup_xkb_extension\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xkb_x11_get_core_keyboard_device_id_handle, err = loadSym(l.libxkbcommonX11Handle, (*C.char)(unsafe.Pointer(&([]byte("xkb_x11_get_core_keyboard_device_id\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xkb_x11_keymap_new_from_device_handle, err = loadSym(l.libxkbcommonX11Handle, (*C.char)(unsafe.Pointer(&([]byte("xkb_x11_keymap_new_from_device\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.xkb_x11_state_new_from_device_handle, err = loadSym(l.libxkbcommonX11Handle, (*C.char)(unsafe.Pointer(&([]byte("xkb_x11_state_new_from_device\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}

	// libxcb
	l.xcb_request_check_handle, err = loadSym(l.libxcbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_request_check\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}

	// libxcb-xkb
	l.xcb_xkb_select_events_aux_checked_handle, err = loadSym(l.libxcbXkbHandle, (*C.char)(unsafe.Pointer(&([]byte("xcb_xkb_select_events_aux_checked\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}

	return l, nil
}

func (l *xkbcommon_library) close() {
	if l.libxcbXkbHandle != nil {
		C.dlclose(l.libxcbXkbHandle)
		l.libxcbXkbHandle = nil
	}
	if l.libxcbHandle != nil {
		C.dlclose(l.libxcbHandle)
		l.libxcbHandle = nil
	}
	if l.libxkbcommonX11Handle != nil {
		C.dlclose(l.libxkbcommonX11Handle)
		l.libxkbcommonX11Handle = nil
	}
	if l.libxkbcommonHandle != nil {
		C.dlclose(l.libxkbcommonHandle)
		l.libxkbcommonHandle = nil
	}
}

// libxkbcommon
func (l *xkbcommon_library) xkb_context_new(flags C.enum_xkb_context_flags) *C.struct_xkb_context {
	return C.gamen_xkb_context_new(l.xkb_context_new_handle, flags)
}
func (l *xkbcommon_library) xkb_compose_table_new_from_locale(context *C.struct_xkb_context, locale *C.char, flags C.enum_xkb_compose_compile_flags) *C.struct_xkb_compose_table {
	return C.gamen_xkb_compose_table_new_from_locale(l.xkb_compose_table_new_from_locale_handle, context, locale, flags)
}
func (l *xkbcommon_library) xkb_compose_state_new(table *C.struct_xkb_compose_table, flags C.enum_xkb_compose_state_flags) *C.struct_xkb_compose_state {
	return C.gamen_xkb_compose_state_new(l.xkb_compose_state_new_handle, table, flags)
}
func (l *xkbcommon_library) xkb_state_unref(state *C.struct_xkb_state) {
	C.gamen_xkb_state_unref(l.xkb_state_unref_handle, state)
}
func (l *xkbcommon_library) xkb_keymap_unref(keymap *C.struct_xkb_keymap) {
	C.gamen_xkb_keymap_unref(l.xkb_keymap_unref_handle, keymap)
}
func (l *xkbcommon_library) xkb_compose_state_unref(state *C.struct_xkb_compose_state) {
	C.gamen_xkb_compose_state_unref(l.xkb_compose_state_unref_handle, state)
}
func (l *xkbcommon_library) xkb_compose_table_unref(table *C.struct_xkb_compose_table) {
	C.gamen_xkb_compose_table_unref(l.xkb_compose_table_unref_handle, table)
}
func (l *xkbcommon_library) xkb_context_unref(context *C.struct_xkb_context) {
	C.gamen_xkb_context_unref(l.xkb_context_unref_handle, context)
}
func (l *xkbcommon_library) xkb_keymap_new_from_buffer(context *C.struct_xkb_context, buffer *C.char, length C.size_t, format C.enum_xkb_keymap_format, flags C.enum_xkb_keymap_compile_flags) *C.struct_xkb_keymap {
	return C.gamen_xkb_keymap_new_from_buffer(l.xkb_keymap_new_from_buffer_handle, context, buffer, length, format, flags)
}
func (l *xkbcommon_library) xkb_state_new(keymap *C.struct_xkb_keymap) *C.struct_xkb_state {
	return C.gamen_xkb_state_new(l.xkb_state_new_handle, keymap)
}
func (l *xkbcommon_library) xkb_keymap_key_repeats(keymap *C.struct_xkb_keymap, key C.xkb_keycode_t) C.int {
	return C.gamen_xkb_keymap_key_repeats(l.xkb_keymap_key_repeats_handle, keymap, key)
}
func (l *xkbcommon_library) xkb_state_key_get_one_sym(state *C.struct_xkb_state, key C.xkb_keycode_t) C.xkb_keysym_t {
	return C.gamen_xkb_state_key_get_one_sym(l.xkb_state_key_get_one_sym_handle, state, key)
}
func (l *xkbcommon_library) xkb_state_key_get_utf8(state *C.struct_xkb_state, key C.xkb_keycode_t, buffer *C.char, size C.size_t) C.int {
	return C.gamen_xkb_state_key_get_utf8(l.xkb_state_key_get_utf8_handle, state, key, buffer, size)
}
func (l *xkbcommon_library) xkb_compose_state_feed(state *C.struct_xkb_compose_state, keysym C.xkb_keysym_t) C.enum_xkb_compose_feed_result {
	return C.gamen_xkb_compose_state_feed(l.xkb_compose_state_feed_handle, state, keysym)
}
func (l *xkbcommon_library) xkb_compose_state_get_status(state *C.struct_xkb_compose_state) C.enum_xkb_compose_status {
	return C.gamen_xkb_compose_state_get_status(l.xkb_compose_state_get_status_handle, state)
}
func (l *xkbcommon_library) xkb_compose_state_get_utf8(state *C.struct_xkb_compose_state, buffer *C.char, size C.size_t) C.int {
	return C.gamen_xkb_compose_state_get_utf8(l.xkb_compose_state_get_utf8_handle, state, buffer, size)
}
func (l *xkbcommon_library) xkb_state_update_mask(state *C.struct_xkb_state, depressed_mods C.xkb_mod_mask_t, latched_mods C.xkb_mod_mask_t, locked_mods C.xkb_mod_mask_t, depressed_layout C.xkb_layout_index_t, latched_layout C.xkb_layout_index_t, locked_layout C.xkb_layout_index_t) C.enum_xkb_state_component {
	return C.gamen_xkb_state_update_mask(l.xkb_state_update_mask_handle, state, depressed_mods, latched_mods, locked_mods, depressed_layout, latched_layout, locked_layout)
}
func (l *xkbcommon_library) xkb_state_mod_name_is_active(state *C.struct_xkb_state, name *C.char, _type C.enum_xkb_state_component) C.int {
	return C.gamen_xkb_state_mod_name_is_active(l.xkb_state_mod_name_is_active_handle, state, name, _type)
}

// libxkbcommon-x11
func (l *xkbcommon_library) xkb_x11_setup_xkb_extension(connection *C.xcb_connection_t, major_xkb_version C.uint16_t, minor_xkb_version C.uint16_t, flags C.enum_xkb_x11_setup_xkb_extension_flags, major_xkb_version_out *C.uint16_t, minor_xkb_version_out *C.uint16_t, base_event_out *C.uint8_t, base_error_out *C.uint8_t) C.int {
	return C.gamen_xkb_x11_setup_xkb_extension(l.xkb_x11_setup_xkb_extension_handle, connection, major_xkb_version, minor_xkb_version, flags, major_xkb_version_out, minor_xkb_version_out, base_event_out, base_error_out)
}
func (l *xkbcommon_library) xkb_x11_get_core_keyboard_device_id(connection *C.xcb_connection_t) C.int32_t {
	return C.gamen_xkb_x11_get_core_keyboard_device_id(l.xkb_x11_get_core_keyboard_device_id_handle, connection)
}
func (l *xkbcommon_library) xkb_x11_keymap_new_from_device(context *C.struct_xkb_context, connection *C.xcb_connection_t, device_id C.int32_t, flags C.enum_xkb_keymap_compile_flags) *C.struct_xkb_keymap {
	return C.gamen_xkb_x11_keymap_new_from_device(l.xkb_x11_keymap_new_from_device_handle, context, connection, device_id, flags)
}
func (l *xkbcommon_library) xkb_x11_state_new_from_device(keymap *C.struct_xkb_keymap, connection *C.xcb_connection_t, device_id C.int32_t) *C.struct_xkb_state {
	return C.gamen_xkb_x11_state_new_from_device(l.xkb_x11_state_new_from_device_handle, keymap, connection, device_id)
}

// libxcb
func (l *xkbcommon_library) xcb_request_check(c *C.xcb_connection_t, cookie C.xcb_void_cookie_t) *C.xcb_generic_error_t {
	return C.gamen_xkbcommon_xcb_request_check(l.xcb_request_check_handle, c, cookie)
}

// libxcb-xkb
func (l *xkbcommon_library) xcb_xkb_select_events_aux_checked(c *C.xcb_connection_t, deviceSpec C.xcb_xkb_device_spec_t, affectWhich C.uint16_t, clear C.uint16_t, selectAll C.uint16_t, affectMap C.uint16_t, _map C.uint16_t, details *C.xcb_xkb_select_events_details_t) C.xcb_void_cookie_t {
	return C.gamen_xcb_xkb_select_events_aux_checked(l.xcb_xkb_select_events_aux_checked_handle, c, deviceSpec, affectWhich, clear, selectAll, affectMap, _map, details)
}
