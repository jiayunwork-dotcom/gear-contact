package gears

import (
	"fmt"
	"math"
)

type Pair struct {
	gear1 Gear
	gear2 Gear
}

func New(spec *Spec) (*Pair, error) {
	if err := spec.Validate(); err != nil {
		return nil, err
	}
	ha := spec.AddendumCoefficient()
	c := spec.ClearanceCoefficient()
	alpha := spec.PressureAngleRad()

	g1 := Gear{
		Module:         spec.Module,
		Teeth:          spec.Gear1.Teeth,
		PressureAngle:  alpha,
		AddendumCoeff:  ha,
		ClearanceCoeff: c,
		Shift:          spec.Gear1.Shift,
	}
	g2 := Gear{
		Module:         spec.Module,
		Teeth:          spec.Gear2.Teeth,
		PressureAngle:  alpha,
		AddendumCoeff:  ha,
		ClearanceCoeff: c,
		Shift:          spec.Gear2.Shift,
	}

	if err := ValidateGear(g1, "gear1"); err != nil {
		return nil, err
	}
	if err := ValidateGear(g2, "gear2"); err != nil {
		return nil, err
	}
	return &Pair{gear1: g1, gear2: g2}, nil
}

func FromJSON(data []byte) (*Pair, error) {
	spec, err := ParseSpec(data)
	if err != nil {
		return nil, err
	}
	return New(spec)
}

func (p *Pair) Gear1() Gear { return p.gear1 }

func (p *Pair) Gear2() Gear { return p.gear2 }

func (p *Pair) Module() float64 { return p.gear1.Module }

func (p *Pair) PressureAngle() float64 { return p.gear1.PressureAngle }

func (p *Pair) TeethSum() int { return p.gear1.Teeth + p.gear2.Teeth }

func (p *Pair) ShiftSum() float64 { return p.gear1.Shift + p.gear2.Shift }

func (p *Pair) StandardCenterDistance() float64 {
	return 0.5 * p.Module() * float64(p.TeethSum())
}

func (p *Pair) Describe() string {
	return fmt.Sprintf("m=%v, z1=%d, z2=%d, alpha=%.3f deg",
		p.Module(), p.Gear1().Teeth, p.Gear2().Teeth,
		radToDeg(p.PressureAngle()))
}

func radToDeg(rad float64) float64 {
	return rad * 180.0 / math.Pi
}
