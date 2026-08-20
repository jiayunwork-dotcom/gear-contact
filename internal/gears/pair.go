package gears

import (
	"fmt"
	"math"
)

// Pair is the two gears of one external mesh together with the values they
// share. The module and the pressure angle are stored on both Gear values; the
// Pair exists so the mesh package can treat the pair as a unit and so the
// invariant "both gears use the same module and pressure angle" is enforced in
// exactly one place (New).
type Pair struct {
	gear1 Gear
	gear2 Gear
}

// New validates a spec and builds the immutable gear pair. Every validation
// error from the spec and from the individual gears is surfaced here, before
// any mesh geometry is attempted.
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

// FromJSON is the single entry point the main command uses: parse the JSON
// spec and build the pair in one call.
func FromJSON(data []byte) (*Pair, error) {
	spec, err := ParseSpec(data)
	if err != nil {
		return nil, err
	}
	return New(spec)
}

// Gear1 returns the first gear.
func (p *Pair) Gear1() Gear { return p.gear1 }

// Gear2 returns the second gear.
func (p *Pair) Gear2() Gear { return p.gear2 }

// Module returns the shared module of the pair.
func (p *Pair) Module() float64 { return p.gear1.Module }

// PressureAngle returns the shared standard pressure angle in radians.
func (p *Pair) PressureAngle() float64 { return p.gear1.PressureAngle }

// TeethSum returns z1 + z2, used by the standard centre distance formula.
func (p *Pair) TeethSum() int { return p.gear1.Teeth + p.gear2.Teeth }

// ShiftSum returns x1 + x2, the argument of the no-backlash involute equation.
func (p *Pair) ShiftSum() float64 { return p.gear1.Shift + p.gear2.Shift }

// StandardCenterDistance returns a0 = m*(z1+z2)/2, the centre distance at
// which the two pitch circles touch. For a standard (zero shift) pair this is
// also the operating centre distance.
func (p *Pair) StandardCenterDistance() float64 {
	return 0.5 * p.Module() * float64(p.TeethSum())
}

// Describe is a short human readable summary used by error messages.
func (p *Pair) Describe() string {
	return fmt.Sprintf("m=%v, z1=%d, z2=%d, alpha=%.3f deg",
		p.Module(), p.Gear1().Teeth, p.Gear2().Teeth,
		radToDeg(p.PressureAngle()))
}

// radToDeg is a local helper so the Pair describe path does not need to import
// the angle package just for one conversion.
func radToDeg(rad float64) float64 {
	return rad * 180.0 / math.Pi
}
