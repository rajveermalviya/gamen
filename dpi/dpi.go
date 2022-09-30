package dpi

import "golang.org/x/exp/constraints"

type Position[T constraints.Integer | constraints.Float] interface {
	implementsPosition()
}

func (LogicalPosition[T]) implementsPosition()  {}
func (PhysicalPosition[T]) implementsPosition() {}

type Size[T constraints.Integer | constraints.Float] interface {
	ToLogical(scaleFactor float64) LogicalSize[T]
	ToPhysical(scaleFactor float64) PhysicalSize[T]
	implementsSize()
}

func (LogicalSize[T]) implementsSize()  {}
func (PhysicalSize[T]) implementsSize() {}

func CastSize[I, O constraints.Integer | constraints.Float](size Size[I]) Size[O] {
	switch size := size.(type) {
	case PhysicalSize[I]:
		return PhysicalSize[O]{
			Width:  O(size.Width),
			Height: O(size.Height),
		}

	case LogicalSize[I]:
		return LogicalSize[O]{
			Width:  O(size.Width),
			Height: O(size.Height),
		}

	default:
		panic("unreachable")
	}
}
