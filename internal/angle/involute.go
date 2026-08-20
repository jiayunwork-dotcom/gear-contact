package angle

import "math"

// Inv computes the involute function inv(x) = tan(x) - x for x in radians.
//
// The involute of the standard pressure angle appears on the right hand side
// of the no-backlash working pressure angle equation, and its inverse has to
// be solved numerically. Keeping the forward and the inverse in one package
// guarantees that both sides of that equation use the same definition.
func Inv(x float64) float64 {
	return math.Tan(x) - x
}

// InvDeriv returns the derivative of Inv with respect to its argument:
// d/dx (tan(x) - x) = tan^2(x). Newton iterations on the involute equation
// use this value to steer the update.
func InvDeriv(x float64) float64 {
	t := math.Tan(x)
	return t * t
}

// InvDomain reports whether x is inside the domain on which the involute
// function is single valued and monotone: (0, pi/2). Outside this interval
// Inv is either negative (for x <= 0) or undefined (at x = pi/2).
func InvDomain(x float64) bool {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return false
	}
	return x > 0 && x < RightAngleRad
}

// InvIncreasing reports whether the involute is monotone increasing between
// the two arguments. The numerical solve relies on monotonicity; when it is
// lost the solver falls back to a bracketed bisection that does not need a
// derivative but still needs a sign change.
func InvIncreasing(a, b float64) bool {
	return InvDomain(a) && InvDomain(b) && a < b
}

// InvOfPressureAngle is a convenience wrapper used by the mesh package: it
// evaluates the involute at the standard pressure angle in degrees. Keeping
// the degree-to-radian conversion here means every call site uses the exact
// same alpha instead of maintaining its own conversion.
func InvOfPressureAngle(alphaRad float64) float64 {
	return Inv(alphaRad)
}
