package utils

import "math"

const eps = 1e-9

func IsZero(a float32) bool {
	return math.Abs(float64(a)) < eps
}
