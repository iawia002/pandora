package math

import (
	"golang.org/x/exp/constraints"
)

// Number is a constraint that permits any numeric type.
type Number interface {
	constraints.Integer | constraints.Float
}

var minusOne = -1

// Abs returns the absolute value of x.
func Abs[X Number](x X) X {
	if x < 0 {
		x *= X(minusOne)
	}
	return x
}
