package gears

import (
	"math"

	"gear-contact/internal/angle"
)

func BuildGear(spec *Spec, gs GearSpec) Gear {
	return Gear{
		Module:         spec.Module,
		Teeth:          gs.Teeth,
		PressureAngle:  spec.PressureAngleRad(),
		AddendumCoeff:  spec.AddendumCoefficient(),
		ClearanceCoeff: spec.ClearanceCoefficient(),
		Shift:          gs.Shift,
	}
}

func Height(g Gear) float64 {
	return (2*g.AddendumCoeff + g.ClearanceCoeff) * g.Module
}

func PitchCircle(g Gear) Circle    { return CircleOf(g, PitchCircleKind) }
func BaseCircle(g Gear) Circle     { return CircleOf(g, BaseCircleKind) }
func AddendumCircle(g Gear) Circle { return CircleOf(g, AddendumCircleKind) }
func DedendumCircle(g Gear) Circle { return CircleOf(g, DedendumCircleKind) }

func WorkingClearanceDenominator(g Gear) float64 {
	return g.DedendumRadius()
}

func NearlyStandard(p *Pair) bool {
	return math.Abs(p.ShiftSum()) < angle.DefaultTol
}
