//go:build linux && !android

package xcb

/*

#include <xcb/xcb.h>
#include <xcb/randr.h>

*/
import "C"

type Output struct {
	xcbScreen *C.xcb_screen_t
	number    int
}

func (d *Display) initializeOutputs(setup *C.struct_xcb_setup_t) {
	for it := C.xcb_setup_roots_iterator(setup); it.rem > 0; C.xcb_screen_next(&it) {
		d.screens = append(d.screens, &Output{
			xcbScreen: it.data,
			number:    len(d.screens),
		})

		C.xcb_randr_select_input(d.xcbConn,
			it.data.root,
			C.XCB_RANDR_NOTIFY_MASK_SCREEN_CHANGE|
				C.XCB_RANDR_NOTIFY_MASK_OUTPUT_CHANGE|
				C.XCB_RANDR_NOTIFY_MASK_CRTC_CHANGE|
				C.XCB_RANDR_NOTIFY_MASK_OUTPUT_PROPERTY,
		)
	}
}
