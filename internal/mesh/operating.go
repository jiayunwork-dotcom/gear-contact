package mesh

import "math"

func OperatingPitchRadius(standardRadius, alpha, alphaP float64) float64 {
	return standardRadius * math.Cos(alpha) / math.Cos(alphaP)
}

func OperatingPitchRadiusByDistance(standardRadius, a, r1, r2 float64) float64 {
	return a * standardRadius / (r1 + r2)
}

func (m *Mesh) OperatingPitchRadii() (float64, float64) {
	g1, g2 := m.Pair.Gear1(), m.Pair.Gear2()
	alpha := m.Pair.PressureAngle()
	r1 := OperatingPitchRadius(g1.PitchRadius(), alpha, m.WorkingPressureAngle)
	r2 := OperatingPitchRadius(g2.PitchRadius(), alpha, m.WorkingPressureAngle)
	return r1, r2
}

func OperatingPitchDiameter(standardDiameter, alpha, alphaP float64) float64 {
	return standardDiameter * math.Cos(alpha) / math.Cos(alphaP)
}
