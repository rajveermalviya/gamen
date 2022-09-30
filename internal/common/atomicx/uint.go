// TODO: remove when we update go.mod to go1.19
package atomicx

import (
	"sync/atomic"
)

type Uint[T ~uint8 | ~uint16 | ~uint32] struct {
	_ noCopy
	v uint32
}

func (x *Uint[T]) Load() T {
	return T(atomic.LoadUint32(&x.v))
}
func (x *Uint[T]) Store(val T) {
	atomic.StoreUint32(&x.v, uint32(val))
}
func (x *Uint[T]) Swap(new T) (old T) {
	return T(atomic.SwapUint32(&x.v, uint32(new)))
}
func (x *Uint[T]) CompareAndSwap(old, new T) (swapped bool) {
	return atomic.CompareAndSwapUint32(&x.v, uint32(old), uint32(new))
}
