package mathx

import (
	"math"

	"golang.org/x/exp/constraints"
)

func Min[T constraints.Integer | constraints.Float](a, b T) T {
	return T(math.Min(float64(a), float64(b)))
}

func Max[T constraints.Integer | constraints.Float](a, b T) T {
	return T(math.Max(float64(a), float64(b)))
}

func Abs[T constraints.Integer | constraints.Float](a T) T {
	return T(math.Abs(float64(a)))
}
