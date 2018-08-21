package types

import "math"

const (
	SmallestNormal32 = 1.17549435E-38          // 2**-126
	SmallestNormal64 = 2.2250738585072014e-308 // 2**-1022
	FLoat64Epsilon   = 1e-9
)

// http://floating-point-gui.de/errors/comparison/
func nearlyEqualImpl(a, b, epsilon, smallest, max float64) bool {
	// Calculate the difference.
	absA := math.Abs(a)
	absB := math.Abs(b)
	diff := math.Abs(a - b)

	if a == b { // shortcut, handles infinities
		return true
	} else if a == 0 || b == 0 || diff < smallest {
		// a or b is zero or both are extremely close to it
		// relative error is less meaningful here
		return diff < (epsilon * smallest)
	} else { // use relative error
		return diff/math.Min(absA+absB, max) < epsilon
	}
	return false
}

func NearlyEqual32(a, b, epsilon float32) bool {
	return nearlyEqualImpl(float64(a), float64(b), float64(epsilon), SmallestNormal32, math.MaxFloat32)
}

func NearlyEqual64(a, b, epsilon float64) bool {
	return nearlyEqualImpl(float64(a), float64(b), float64(epsilon), SmallestNormal64, math.MaxFloat64)
}
