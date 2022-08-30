//go:build linux && !android

package xcb

/*

#include <stdlib.h>
#include <X11/Xlib-xcb.h>
#include <xcb/xcb_image.h>
#include <X11/Xcursor/Xcursor.h>

*/
import "C"
import (
	"unsafe"

	"github.com/rajveermalviya/gamen/cursors"
	"github.com/rajveermalviya/gamen/internal/common/xcursor"
)

func createEmptyCursor(xconn *C.xcb_connection_t, root C.xcb_window_t) C.xcb_cursor_t {
	var buf [32]C.uint8_t

	source := C.xcb_create_pixmap_from_bitmap_data(xconn,
		root,
		(*C.uint8_t)(unsafe.Pointer(&buf)),
		16, 16,
		1,
		0, 0,
		nil,
	)
	defer C.xcb_free_pixmap(xconn, source)

	mask := C.xcb_create_pixmap_from_bitmap_data(xconn,
		root,
		(*C.uint8_t)(unsafe.Pointer(&buf)),
		16, 16,
		1,
		0, 0,
		nil,
	)
	defer C.xcb_free_pixmap(xconn, mask)

	cursorId := C.xcb_generate_id(xconn)
	C.xcb_create_cursor(xconn,
		cursorId,
		source, mask,
		0, 0, 0,
		0xFFFF, 0xFFFF, 0xFFFF,
		8, 8,
	)
	return cursorId
}

func loadCursor(dpy *C.Display, name string) C.xcb_cursor_t {
	nameStr := C.CString(name)
	defer C.free(unsafe.Pointer(nameStr))

	cursor := C.XcursorLibraryLoadCursor(dpy, nameStr)
	return C.xcb_cursor_t(cursor)
}

func (d *Display) loadCursorIcon(icon cursors.Icon) C.xcb_cursor_t {
	d.mu.Lock()
	defer d.mu.Unlock()

	c, ok := d.cursors[icon]
	if ok {
		return c
	}

	if icon == 0 {
		c = createEmptyCursor(d.xcbConn, d.screens[0].xcbScreen.root)
		d.cursors[icon] = c
		return c
	}

	for _, name := range xcursor.ToXcursorName(icon) {
		c = loadCursor(d.xlibDisp, name)
		if c != 0 {
			break
		}
	}

	d.cursors[icon] = c
	return c
}
