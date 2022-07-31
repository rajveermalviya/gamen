//go:build linux && !android

package xcb

import (
	"errors"
	"unsafe"

	"golang.org/x/exp/constraints"
)

/*

#include <stdlib.h>
#include <xcb/xcb.h>
#include <xcb/xinput.h>
#include <xcb/xproto.h>

*/
import "C"

func setXiMask(mask *[8]C.uint8_t, bit C.int) {
	mask[bit>>3] |= 1 << (bit & 7)
}

func checkXInputVersion(conn *C.xcb_connection_t) error {
	cookie := C.xcb_input_xi_query_version(conn, C.XCB_INPUT_MAJOR_VERSION, C.XCB_INPUT_MINOR_VERSION)
	reply := C.xcb_input_xi_query_version_reply(conn, cookie, nil)
	defer C.free(unsafe.Pointer(reply))

	if reply.major_version >= 2 && reply.minor_version >= 1 {
		return nil
	}
	return errors.New("xinput version not supported")
}

func internAtom(conn *C.xcb_connection_t, onlyIfExists bool, name string) C.xcb_atom_t {
	nameStr := C.CString(name)
	defer C.free(unsafe.Pointer(nameStr))
	cookie := C.xcb_intern_atom(conn, Cbool[C.uint8_t](onlyIfExists), C.uint16_t(len(name)), nameStr)
	reply := C.xcb_intern_atom_reply(conn, cookie, nil)
	defer C.free(unsafe.Pointer(reply))
	return reply.atom
}

func Cbool[T constraints.Integer | constraints.Float](v bool) T {
	if v {
		return 1
	}
	return 0
}

// https://gitlab.freedesktop.org/xorg/lib/libxi/-/blob/bca3474a8622fde5815260461784282f78a4efb5/src/XExtInt.c#L74
func fixed1616ToFloat64(v C.xcb_input_fp1616_t) float64 {
	return float64(v) * 1.0 / (1 << 16)
}

// https://gitlab.freedesktop.org/xorg/lib/libxi/-/blob/bca3474a8622fde5815260461784282f78a4efb5/src/XExtInt.c#L1700
func fixed3232ToFloat64(v C.xcb_input_fp3232_t) float64 {
	return float64(v.integral) + float64(v.frac)/(1<<32)
}
