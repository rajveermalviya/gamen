//go:build linux && !android

package xcb

import (
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

func hasXiMask(mask []C.uint32_t, bit C.uint16_t) bool {
	return mask[bit>>3]&(1<<(bit&7)) != 0
}

func (d *Display) internAtom(onlyIfExists bool, name string) C.xcb_atom_t {
	nameStr := C.CString(name)
	defer C.free(unsafe.Pointer(nameStr))
	cookie := d.l.xcb_intern_atom(d.xcbConn, Cbool[C.uint8_t](onlyIfExists), C.uint16_t(len(name)), nameStr)
	reply := d.l.xcb_intern_atom_reply(d.xcbConn, cookie, nil)
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
