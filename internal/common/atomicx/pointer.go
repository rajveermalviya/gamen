// TODO: remove when we update go.mod to go1.19
package atomicx

import (
	"sync/atomic"
	"unsafe"
)

type Pointer[T any] struct {
	_ noCopy
	v unsafe.Pointer
}

func (x *Pointer[T]) Load() *T {
	return (*T)(atomic.LoadPointer(&x.v))
}
func (x *Pointer[T]) Store(val *T) {
	atomic.StorePointer(&x.v, unsafe.Pointer(val))
}
func (x *Pointer[T]) Swap(new *T) (old *T) {
	return (*T)(atomic.SwapPointer(&x.v, unsafe.Pointer(new)))
}
func (x *Pointer[T]) CompareAndSwap(old, new *T) (swapped bool) {
	return atomic.CompareAndSwapPointer(&x.v, unsafe.Pointer(old), unsafe.Pointer(new))
}
