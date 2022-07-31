package dpi

import "golang.org/x/exp/constraints"

type Position[T constraints.Integer | constraints.Float] interface{ implementsPosition(T) }

func (LogicalPosition[T]) implementsPosition(T)  {}
func (PhysicalPosition[T]) implementsPosition(T) {}

type Size[T constraints.Integer | constraints.Float] interface {
	ToLogical(scaleFactor float64) LogicalSize[T]
	ToPhysical(scaleFactor float64) PhysicalSize[T]
	implementsSize(T)
}

func (LogicalSize[T]) implementsSize(T)  {}
func (PhysicalSize[T]) implementsSize(T) {}

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
