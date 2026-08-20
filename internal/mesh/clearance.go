package mesh

// Tip clearance. The clearance is the radial gap between the tip of one gear
// and the root of the other, measured along the line of centres:
//
//	c = a - ra1 - rf2   (equivalently c = a - ra2 - rf1)
//
// where a is the operating centre distance, ra the addendum (tip) radius of
// one gear and rf the dedendum (root) radius of the other. The two formulas
// are equal whenever both gears share the same module and coefficients, which
// the pair invariant guarantees.
//
// For a standard mesh (zero shift) this works out to c = c**m, positive as
// long as the dedendum is deeper than the addendum by the clearance
// coefficient. The default c* = 0.25 gives a standard clearance of 0.25*m.

// TipClearance returns c = a - ra1 - rf2.
func TipClearance(operatingCenter, ra1, rf2 float64) float64 {
	return operatingCenter - ra1 - rf2
}

// TipClearanceSymmetric is the clearance computed from the other pair of
// circles: c = a - ra2 - rf1. For a consistent pair it equals TipClearance;
// the tests use the two together to verify the pair invariant.
func TipClearanceSymmetric(operatingCenter, ra2, rf1 float64) float64 {
	return operatingCenter - ra2 - rf1
}

// ClearanceCoefficient recovers c*/m from a measured clearance: c* = c/m.
// Standard mounting with standard proportions yields 0.25.
func ClearanceCoefficient(clearance, module float64) float64 {
	return clearance / module
}

// ClearanceSign describes whether the clearance is positive, zero or negative
// as an int so the report can render a human word instead of a bare sign.
func ClearanceSign(clearance float64) int {
	switch {
	case clearance > 0:
		return 1
	case clearance < 0:
		return -1
	default:
		return 0
	}
}

// ClearanceWord maps a clearance sign to "positive", "zero" or "negative".
func ClearanceWord(clearance float64) string {
	switch ClearanceSign(clearance) {
	case 1:
		return "positive"
	case -1:
		return "negative"
	default:
		return "zero"
	}
}
