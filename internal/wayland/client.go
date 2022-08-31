//go:build linux && !android

package wayland

/*

#cgo LDFLAGS: -ldl

#include <dlfcn.h>
#include "wayland-cursor.h"

struct wl_display;
struct wl_proxy;

struct wl_display *gamen_wl_display_connect(void *fp, const char *name) {
	typedef struct wl_display *(*proc_wl_display_connect)(const char *name);
	return ((proc_wl_display_connect)fp)(name);
}

void gamen_wl_display_disconnect(void *fp, struct wl_display *display) {
	typedef void (*proc_wl_display_disconnect)(struct wl_display *display);
	((proc_wl_display_disconnect)fp)(display);
}

int gamen_wl_display_roundtrip(void *fp, struct wl_display *display) {
	typedef int (*proc_wl_display_roundtrip)(struct wl_display *display);
	return ((proc_wl_display_roundtrip)fp)(display);
}

int gamen_wl_display_prepare_read(void *fp, struct wl_display *display) {
	typedef int (*proc_wl_display_prepare_read)(struct wl_display *display);
	return ((proc_wl_display_prepare_read)fp)(display);
}

int gamen_wl_display_dispatch_pending(void *fp, struct wl_display *display) {
	typedef int (*proc_wl_display_dispatch_pending)(struct wl_display *display);
	return ((proc_wl_display_dispatch_pending)fp)(display);
}

int gamen_wl_display_flush(void *fp, struct wl_display *display) {
	typedef int (*proc_wl_display_flush)(struct wl_display *display);
	return ((proc_wl_display_flush)fp)(display);
}

int gamen_wl_display_get_fd(void *fp, struct wl_display *display) {
	typedef int (*proc_wl_display_get_fd)(struct wl_display *display);
	return ((proc_wl_display_get_fd)fp)(display);
}

void gamen_wl_display_cancel_read(void *fp, struct wl_display *display) {
	typedef void (*proc_wl_display_cancel_read)(struct wl_display *display);
	((proc_wl_display_cancel_read)fp)(display);
}

int gamen_wl_display_read_events(void *fp, struct wl_display *display) {
	typedef int (*proc_wl_display_read_events)(struct wl_display *display);
	return ((proc_wl_display_read_events)fp)(display);
}

void gamen_wl_proxy_destroy(void *fp, struct wl_proxy *proxy) {
	typedef void (*proc_wl_proxy_destroy)(struct wl_proxy *proxy);
	((proc_wl_proxy_destroy)fp)(proxy);
}


struct wl_cursor_theme *gamen_wl_cursor_theme_load(void *fp, const char *name, int size, struct wl_shm *shm) {
	typedef struct wl_cursor_theme *(*proc_wl_cursor_theme_load)(const char *name, int size, struct wl_shm *shm);
	return ((proc_wl_cursor_theme_load)fp)(name, size, shm);
}

struct wl_cursor *gamen_wl_cursor_theme_get_cursor(void *fp, struct wl_cursor_theme *theme, const char *name) {
	typedef struct wl_cursor *(*proc_wl_cursor_theme_get_cursor)(struct wl_cursor_theme *theme, const char *name);
	return ((proc_wl_cursor_theme_get_cursor)fp)(theme, name);
}

void gamen_wl_cursor_theme_destroy(void *fp, struct wl_cursor_theme *theme) {
	typedef void (*proc_wl_cursor_theme_destroy)(struct wl_cursor_theme *theme);
	((proc_wl_cursor_theme_destroy)fp)(theme);
}

struct wl_buffer *gamen_wl_cursor_image_get_buffer(void *fp, struct wl_cursor_image *image) {
	typedef struct wl_buffer *(*proc_wl_cursor_image_get_buffer)(struct wl_cursor_image *image);
	return ((proc_wl_cursor_image_get_buffer)fp)(image);
}

int gamen_wl_cursor_frame_and_duration(void *fp, struct wl_cursor *cursor, uint32_t time, uint32_t *duration) {
	typedef int (*proc_wl_cursor_frame_and_duration)(struct wl_cursor *cursor, uint32_t time, uint32_t *duration);
	return ((proc_wl_cursor_frame_and_duration)fp)(cursor, time, duration);
}

*/
import "C"
import (
	"errors"
	"unsafe"
)

type wl_library struct {
	libWaylandClientHandle unsafe.Pointer
	libWaylandCursorHandle unsafe.Pointer

	wl_display_connect_handle          unsafe.Pointer
	wl_display_roundtrip_handle        unsafe.Pointer
	wl_display_disconnect_handle       unsafe.Pointer
	wl_display_prepare_read_handle     unsafe.Pointer
	wl_display_dispatch_pending_handle unsafe.Pointer
	wl_display_flush_handle            unsafe.Pointer
	wl_display_get_fd_handle           unsafe.Pointer
	wl_display_cancel_read_handle      unsafe.Pointer
	wl_display_read_events_handle      unsafe.Pointer
	wl_proxy_add_listener_handle       unsafe.Pointer
	wl_proxy_destroy_handle            unsafe.Pointer
	wl_proxy_marshal_flags             unsafe.Pointer
	wl_proxy_get_version               unsafe.Pointer

	wl_cursor_theme_load_handle         unsafe.Pointer
	wl_cursor_theme_get_cursor_handle   unsafe.Pointer
	wl_cursor_theme_destroy_handle      unsafe.Pointer
	wl_cursor_image_get_buffer_handle   unsafe.Pointer
	wl_cursor_frame_and_duration_handle unsafe.Pointer
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

func open_wl_library() (*wl_library, error) {
	C.dlerror()

	l := &wl_library{}

	l.libWaylandClientHandle = C.dlopen((*C.char)(unsafe.Pointer(&([]byte("libwayland-client.so.0\x00"))[0])), C.RTLD_LAZY)
	if l.libWaylandClientHandle == nil {
		err := C.dlerror()
		if err != nil {
			l.close()
			return nil, errors.New(C.GoString(err))
		}
	}

	l.libWaylandCursorHandle = C.dlopen((*C.char)(unsafe.Pointer(&([]byte("libwayland-cursor.so.0\x00"))[0])), C.RTLD_LAZY)
	if l.libWaylandCursorHandle == nil {
		err := C.dlerror()
		if err != nil {
			l.close()
			return nil, errors.New(C.GoString(err))
		}
	}

	var err error
	l.wl_display_connect_handle, err = loadSym(l.libWaylandClientHandle, (*C.char)(unsafe.Pointer(&([]byte("wl_display_connect\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.wl_display_roundtrip_handle, err = loadSym(l.libWaylandClientHandle, (*C.char)(unsafe.Pointer(&([]byte("wl_display_roundtrip\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.wl_display_disconnect_handle, err = loadSym(l.libWaylandClientHandle, (*C.char)(unsafe.Pointer(&([]byte("wl_display_disconnect\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.wl_display_prepare_read_handle, err = loadSym(l.libWaylandClientHandle, (*C.char)(unsafe.Pointer(&([]byte("wl_display_prepare_read\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.wl_display_dispatch_pending_handle, err = loadSym(l.libWaylandClientHandle, (*C.char)(unsafe.Pointer(&([]byte("wl_display_dispatch_pending\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.wl_display_flush_handle, err = loadSym(l.libWaylandClientHandle, (*C.char)(unsafe.Pointer(&([]byte("wl_display_flush\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.wl_display_get_fd_handle, err = loadSym(l.libWaylandClientHandle, (*C.char)(unsafe.Pointer(&([]byte("wl_display_get_fd\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.wl_display_cancel_read_handle, err = loadSym(l.libWaylandClientHandle, (*C.char)(unsafe.Pointer(&([]byte("wl_display_cancel_read\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.wl_display_read_events_handle, err = loadSym(l.libWaylandClientHandle, (*C.char)(unsafe.Pointer(&([]byte("wl_display_read_events\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.wl_proxy_add_listener_handle, err = loadSym(l.libWaylandClientHandle, (*C.char)(unsafe.Pointer(&([]byte("wl_proxy_add_listener\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.wl_proxy_destroy_handle, err = loadSym(l.libWaylandClientHandle, (*C.char)(unsafe.Pointer(&([]byte("wl_proxy_destroy\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.wl_proxy_marshal_flags, err = loadSym(l.libWaylandClientHandle, (*C.char)(unsafe.Pointer(&([]byte("wl_proxy_marshal_flags\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.wl_proxy_get_version, err = loadSym(l.libWaylandClientHandle, (*C.char)(unsafe.Pointer(&([]byte("wl_proxy_get_version\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}

	l.wl_cursor_theme_load_handle, err = loadSym(l.libWaylandCursorHandle, (*C.char)(unsafe.Pointer(&([]byte("wl_cursor_theme_load\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.wl_cursor_theme_get_cursor_handle, err = loadSym(l.libWaylandCursorHandle, (*C.char)(unsafe.Pointer(&([]byte("wl_cursor_theme_get_cursor\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.wl_cursor_theme_destroy_handle, err = loadSym(l.libWaylandCursorHandle, (*C.char)(unsafe.Pointer(&([]byte("wl_cursor_theme_destroy\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.wl_cursor_image_get_buffer_handle, err = loadSym(l.libWaylandCursorHandle, (*C.char)(unsafe.Pointer(&([]byte("wl_cursor_image_get_buffer\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}
	l.wl_cursor_frame_and_duration_handle, err = loadSym(l.libWaylandCursorHandle, (*C.char)(unsafe.Pointer(&([]byte("wl_cursor_frame_and_duration\x00"))[0])))
	if err != nil {
		l.close()
		return nil, err
	}

	return l, nil
}

func (l *wl_library) close() {
	if l.libWaylandCursorHandle != nil {
		C.dlclose(l.libWaylandCursorHandle)
		l.libWaylandCursorHandle = nil
	}
	if l.libWaylandClientHandle != nil {
		C.dlclose(l.libWaylandClientHandle)
		l.libWaylandClientHandle = nil
	}
}

func (l *wl_library) wl_display_connect(name *C.char) *C.struct_wl_display {
	return C.gamen_wl_display_connect(l.wl_display_connect_handle, name)
}
func (l *wl_library) wl_display_roundtrip(display *C.struct_wl_display) C.int {
	return C.gamen_wl_display_roundtrip(l.wl_display_roundtrip_handle, display)
}
func (l *wl_library) wl_display_disconnect(display *C.struct_wl_display) {
	C.gamen_wl_display_disconnect(l.wl_display_disconnect_handle, display)
}
func (l *wl_library) wl_display_prepare_read(display *C.struct_wl_display) C.int {
	return C.gamen_wl_display_prepare_read(l.wl_display_prepare_read_handle, display)
}
func (l *wl_library) wl_display_dispatch_pending(display *C.struct_wl_display) C.int {
	return C.gamen_wl_display_dispatch_pending(l.wl_display_dispatch_pending_handle, display)
}
func (l *wl_library) wl_display_flush(display *C.struct_wl_display) (C.int, error) {
	r, err := C.gamen_wl_display_flush(l.wl_display_flush_handle, display)
	return r, err
}
func (l *wl_library) wl_display_get_fd(display *C.struct_wl_display) C.int {
	return C.gamen_wl_display_get_fd(l.wl_display_get_fd_handle, display)
}
func (l *wl_library) wl_display_cancel_read(display *C.struct_wl_display) {
	C.gamen_wl_display_cancel_read(l.wl_display_cancel_read_handle, display)
}
func (l *wl_library) wl_display_read_events(display *C.struct_wl_display) C.int {
	return C.gamen_wl_display_read_events(l.wl_display_read_events_handle, display)
}
func (l *wl_library) wl_proxy_destroy(proxy *C.struct_wl_proxy) {
	C.gamen_wl_proxy_destroy(l.wl_proxy_destroy_handle, proxy)
}

func (l *wl_library) wl_cursor_theme_load(name *C.char, size C.int, shm *C.struct_wl_shm) *C.struct_wl_cursor_theme {
	return C.gamen_wl_cursor_theme_load(l.wl_cursor_theme_load_handle, name, size, shm)
}

func (l *wl_library) wl_cursor_theme_get_cursor(theme *C.struct_wl_cursor_theme, name *C.char) *C.struct_wl_cursor {
	return C.gamen_wl_cursor_theme_get_cursor(l.wl_cursor_theme_get_cursor_handle, theme, name)
}

func (l *wl_library) wl_cursor_theme_destroy(theme *C.struct_wl_cursor_theme) {
	C.gamen_wl_cursor_theme_destroy(l.wl_cursor_theme_destroy_handle, theme)
}

func (l *wl_library) wl_cursor_image_get_buffer(image *C.struct_wl_cursor_image) *C.struct_wl_buffer {
	return C.gamen_wl_cursor_image_get_buffer(l.wl_cursor_image_get_buffer_handle, image)
}

func (l *wl_library) wl_cursor_frame_and_duration(cursor *C.struct_wl_cursor, time C.uint32_t, duration *C.uint32_t) C.int {
	return C.gamen_wl_cursor_frame_and_duration(l.wl_cursor_frame_and_duration_handle, cursor, time, duration)
}
