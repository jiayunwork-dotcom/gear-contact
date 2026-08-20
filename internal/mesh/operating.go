package mesh

import "math"

// Operating pitch circles. When a pair runs at the operating centre distance,
// the actual rolling circles are not the reference pitch circles any more
// (unless the pair is standard). The operating pitch radius of gear i is
//
//	r_i' = r_i * cos(alpha) / cos(alpha')
//
// equivalently r_i' = a * r_i / (r1 + r2), which is the same value. The
// operating pitch circles are the circles that roll on each other without
// slipping at the given mounting distance, so they are part of a complete
// engagement calculation.

// OperatingPitchRadius scales a reference pitch radius to the operating
// mounting: r' = r*cos(alpha)/cos(alpha'). The two formulas in the package
// comment are equal; this one makes the dependence on the working pressure
// angle explicit.
func OperatingPitchRadius(standardRadius, alpha, alphaP float64) float64 {
	return standardRadius * math.Cos(alpha) / math.Cos(alphaP)
}

// OperatingPitchRadiusByDistance is the same value expressed through the
// centre distance: r' = a*r/(r1+r2). Dividing by zero is avoided by the
// caller's guarantee that both pitch radii are positive.
func OperatingPitchRadiusByDistance(standardRadius, a, r1, r2 float64) float64 {
	return a * standardRadius / (r1 + r2)
}

// OperatingPitchRadii returns the operating pitch radii of both gears. For a
// standard mesh they equal the reference pitch radii.
func (m *Mesh) OperatingPitchRadii() (float64, float64) {
	g1, g2 := m.Pair.Gear1(), m.Pair.Gear2()
	alpha := m.Pair.PressureAngle()
	r1 := OperatingPitchRadius(g1.PitchRadius(), alpha, m.WorkingPressureAngle)
	r2 := OperatingPitchRadius(g2.PitchRadius(), alpha, m.WorkingPressureAngle)
	return r1, r2
}

// OperatingPitchDiameter is the diameter form of the operating pitch radius.
func OperatingPitchDiameter(standardDiameter, alpha, alphaP float64) float64 {
	return standardDiameter * math.Cos(alpha) / math.Cos(alphaP)
}
