package dpi

import "golang.org/x/exp/constraints"

type PhysicalPosition[T constraints.Integer | constraints.Float] struct {
	X, Y T
}

type PhysicalSize[T constraints.Integer | constraints.Float] struct {
	Width, Height T
}

func (s PhysicalSize[T]) ToLogical(scaleFactor float64) LogicalSize[T] {
	return LogicalSize[T]{
		Width:  T(float64(s.Width) / scaleFactor),
		Height: T(float64(s.Height) / scaleFactor),
	}
}

func (s PhysicalSize[T]) ToPhysical(scaleFactor float64) PhysicalSize[T] {
	return s
}
