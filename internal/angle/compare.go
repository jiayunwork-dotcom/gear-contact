package angle

import "math"

// Floating point comparisons for geometry. The tolerance policy is: absolute
// tolerance for values near zero, relative tolerance for values whose scale
// matters. Both helpers are conservative; a geometry check that passes them
// is correct to well beyond the precision a gear mesh needs.

// DefaultTol is the default absolute tolerance in millimetre-scale geometry
// (diameters and lengths in the order of tens of millimetres).
const DefaultTol = 1e-9

// CloseTo reports whether a and b differ by at most tol in absolute terms.
func CloseTo(a, b, tol float64) bool {
	if math.IsNaN(a) || math.IsNaN(b) {
		return false
	}
	return math.Abs(a-b) <= tol
}

// RelCloseTo reports whether a and b agree to rel relative tolerance. The
// comparison uses the larger magnitude of the two as the scale, so the check
// is symmetric under argument swap.
func RelCloseTo(a, b, rel float64) bool {
	if math.IsNaN(a) || math.IsNaN(b) {
		return false
	}
	scale := math.Max(math.Abs(a), math.Abs(b))
	if scale == 0 {
		return a == b
	}
	return math.Abs(a-b)/scale <= rel
}

// Within reports whether x lies inside the closed interval [lo, hi] with a
// relative tolerance so boundary hits are not rejected by rounding.
func Within(x, lo, hi, rel float64) bool {
	return x+rel*math.Max(math.Abs(lo), math.Abs(hi)) >= lo &&
		x-rel*math.Max(math.Abs(lo), math.Abs(hi)) <= hi
}

// Sign returns -1, 0 or +1 depending on the sign of v. It is used when the
// geometry code has to branch on the sign of a quantity (for example the
// clearance between two circles) and wants an explicit three-way result.
func Sign(v float64) int {
	if math.IsNaN(v) {
		return 0
	}
	if v < 0 {
		return -1
	}
	if v > 0 {
		return 1
	}
	return 0
}

// Ratio computes v/u and returns 0 instead of +Inf when u is numerically zero.
// Dividing a finite length by a numerically zero denominator happens only for
// degenerate inputs that validation rejects earlier; this helper turns that
// into a defined zero instead of propagating an infinity.
func Ratio(v, u float64) float64 {
	if u == 0 {
		return 0
	}
	return v / u
}
