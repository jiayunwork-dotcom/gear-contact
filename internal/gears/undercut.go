package gears

import "math"

// Undercut geometry. An undercut appears when the generating rack tooth cuts
// into the involute flank near the root; the practical consequence is a loss
// of usable flank length and a weakening of the tooth. For a standard rack the
// borderline is the minimum number of teeth z_min = 2*ha*/sin^2(alpha) at
// which the involute starts exactly at the base circle. At alpha = 20 degrees
// and ha* = 1 that is 17.097 teeth, which is where the classical "17 tooth"
// rule of thumb comes from.
//
// A profile shift x moves the boundary: the flank is safe when
// x >= x_min = ha**(1 - z/z_min). Rearranged this reads
// x_min = ha* - z*sin^2(alpha)/2, which is the formula the README pins.

// MinTeeth returns the smallest tooth count that meshes without undercut at
// zero profile shift: z_min = 2*ha*/sin^2(alpha).
func MinTeeth(pressureAngle, addendumCoeff float64) float64 {
	sin := math.Sin(pressureAngle)
	return 2 * addendumCoeff / (sin * sin)
}

// MinShift returns the smallest profile shift that avoids undercut for a gear
// with the given tooth count: x_min = ha**(1 - z/z_min). A negative result
// means the gear is already safe at standard (zero) shift; a positive result
// means that much positive shift is required.
func MinShift(pressureAngle, addendumCoeff float64, teeth int) float64 {
	zmin := MinTeeth(pressureAngle, addendumCoeff)
	return addendumCoeff * (1 - float64(teeth)/zmin)
}

// UndercutLimit returns x_min for this gear, rounding the tooth count to the
// continuous value the formula expects.
func (g Gear) UndercutLimit() float64 {
	return MinShift(g.PressureAngle, g.AddendumCoeff, g.Teeth)
}

// IsUndercut reports whether the gear's profile shift is below the limit at
// which the flank stays free of undercut. A gear that sits exactly on the
// boundary (x == x_min) is accepted: the involute then starts exactly at the
// base circle and there is nothing cut away. The comparison therefore uses a
// tiny margin so floating point does not reject a boundary gear.
func IsUndercut(g Gear) bool {
	return g.Shift < g.UndercutLimit()-1e-9
}

// UndercutMargin returns how far above (positive) or below (negative) the
// minimum shift the gear sits: x - x_min. The report prints it as a
// diagnostic so a marginal gear is visible before the validation even runs.
func UndercutMargin(g Gear) float64 {
	return g.Shift - g.UndercutLimit()
}
