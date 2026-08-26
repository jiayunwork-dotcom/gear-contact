package gears

import "math"

func MinTeeth(pressureAngle, addendumCoeff float64) float64 {
	sin := math.Sin(pressureAngle)
	return 2 * addendumCoeff / (sin * sin)
}

func MinShift(pressureAngle, addendumCoeff float64, teeth int) float64 {
	zmin := MinTeeth(pressureAngle, addendumCoeff)
	return addendumCoeff * (1 - float64(teeth)/zmin)
}

func (g Gear) UndercutLimit() float64 {
	return MinShift(g.PressureAngle, g.AddendumCoeff, g.Teeth)
}

func IsUndercut(g Gear) bool {
	return g.Shift < g.UndercutLimit()-1e-9
}

func UndercutMargin(g Gear) float64 {
	return g.Shift - g.UndercutLimit()
}
