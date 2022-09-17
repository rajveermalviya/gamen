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

func (d *Display) createEmptyCursor() C.xcb_cursor_t {
	var buf [32]C.uint8_t

	root := d.screens[0].xcbScreen.root

	source := d.l.xcb_create_pixmap_from_bitmap_data(d.xcbConn,
		root,
		(*C.uint8_t)(unsafe.Pointer(&buf)),
		16, 16,
		1,
		0, 0,
		nil,
	)
	defer d.l.xcb_free_pixmap(d.xcbConn, source)

	mask := d.l.xcb_create_pixmap_from_bitmap_data(d.xcbConn,
		root,
		(*C.uint8_t)(unsafe.Pointer(&buf)),
		16, 16,
		1,
		0, 0,
		nil,
	)
	defer d.l.xcb_free_pixmap(d.xcbConn, mask)

	cursorId := d.l.xcb_generate_id(d.xcbConn)
	d.l.xcb_create_cursor(d.xcbConn,
		cursorId,
		source, mask,
		0, 0, 0,
		0xFFFF, 0xFFFF, 0xFFFF,
		8, 8,
	)
	return cursorId
}

func (d *Display) loadXCursor(name string) C.xcb_cursor_t {
	nameStr := C.CString(name)
	defer C.free(unsafe.Pointer(nameStr))

	cursor := d.l.XcursorLibraryLoadCursor(d.xlibDisp, nameStr)
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
		c = d.createEmptyCursor()
		d.cursors[icon] = c
		return c
	}

	for _, name := range xcursor.ToXcursorName(icon) {
		c = d.loadXCursor(name)
		if c != 0 {
			break
		}
	}

	d.cursors[icon] = c
	return c
}
