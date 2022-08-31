//go:build linux && !android

package wayland

/*

#include "wayland-util.h"

*/
import "C"
import "unsafe"

func castWlArrayToSlice[T any](array *C.struct_wl_array) []T {
	var out T
	return unsafe.Slice((*T)(array.data), uintptr(array.size)/unsafe.Sizeof(out))
}
