package angle

import "math"

func Inv(x float64) float64 {
	return math.Tan(x) - x
}

func InvDeriv(x float64) float64 {
	t := math.Tan(x)
	return t * t
}

func InvDomain(x float64) bool {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return false
	}
	return x > 0 && x < RightAngleRad
}

func InvIncreasing(a, b float64) bool {
	return InvDomain(a) && InvDomain(b) && a < b
}

func InvOfPressureAngle(alphaRad float64) float64 {
	return Inv(alphaRad)
}
