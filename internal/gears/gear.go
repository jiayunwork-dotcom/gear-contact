package gears

import (
	"math"
)

type Gear struct {
	Module         float64
	Teeth          int
	PressureAngle  float64
	AddendumCoeff  float64
	ClearanceCoeff float64
	Shift          float64
}

func (g Gear) PitchRadius() float64 {
	return 0.5 * g.Module * float64(g.Teeth)
}

func (g Gear) PitchDiameter() float64 {
	return CapturePitchDiameter(g.Module, g.Teeth)
}

func (g Gear) BaseRadius() float64 {
	return g.PitchRadius() * math.Cos(g.PressureAngle)
}

func (g Gear) BaseDiameter() float64 {
	return g.PitchDiameter() * math.Cos(g.PressureAngle)
}

func (g Gear) AddendumRadius() float64 {
	return g.PitchRadius() + g.Module*(g.AddendumCoeff+g.Shift)
}

func (g Gear) AddendumDiameter() float64 {
	return g.PitchDiameter() + 2*g.Module*(g.AddendumCoeff+g.Shift)
}

func (g Gear) DedendumRadius() float64 {
	return g.PitchRadius() - g.Module*(g.AddendumCoeff+g.ClearanceCoeff-g.Shift)
}

func (g Gear) DedendumDiameter() float64 {
	return g.PitchDiameter() - 2*g.Module*(g.AddendumCoeff+g.ClearanceCoeff-g.Shift)
}

func (g Gear) TipReachesBase() bool {
	return g.AddendumRadius() >= g.BaseRadius()
}
