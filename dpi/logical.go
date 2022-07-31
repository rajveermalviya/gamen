package dpi

import "golang.org/x/exp/constraints"

type LogicalPosition[T constraints.Integer | constraints.Float] struct {
	X, Y T
}

func (s LogicalPosition[T]) ToLogical(scaleFactor float64) LogicalPosition[T] {
	return s
}

func (s LogicalPosition[T]) ToPhysical(scaleFactor float64) PhysicalPosition[T] {
	return PhysicalPosition[T]{
		X: T(float64(s.X) * scaleFactor),
		Y: T(float64(s.Y) * scaleFactor),
	}
}

type LogicalSize[T constraints.Integer | constraints.Float] struct {
	Width, Height T
}

func (s LogicalSize[T]) ToLogical(scaleFactor float64) LogicalSize[T] {
	return s
}

func (s LogicalSize[T]) ToPhysical(scaleFactor float64) PhysicalSize[T] {
	return PhysicalSize[T]{
		Width:  T(float64(s.Width) * scaleFactor),
		Height: T(float64(s.Height) * scaleFactor),
	}
}
