package gears

import (
	"math"

	"gear-contact/internal/angle"
)

// build.go turns a validated Spec into the two Gear values. The construction
// is kept separate from Pair.New so the mesh package can build individual gear
// values without going through a full spec when it only needs a radius.

// BuildGear constructs a Gear from a spec and one gear entry. The pressure
// angle is shared by construction: the pair always passes the same alpha to
// both gears, which is the invariant the question demands (the two gears must
// not end up with inconsistent pressure angles or modules).
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

// Height returns the tooth height h = (2*ha* + c*) - 0 + x terms cancel when
// both gears share the same coefficients, so the tooth height of a gear is
// h = (2*ha* + c*) * m. It is printed as a diagnostic column in the report.
func Height(g Gear) float64 {
	return (2*g.AddendumCoeff + g.ClearanceCoeff) * g.Module
}

// PitchCircle / BaseCircle / AddendumCircle / DedendumCircle are convenience
// aliases that keep the circle API at call sites short.
func PitchCircle(g Gear) Circle    { return CircleOf(g, PitchCircleKind) }
func BaseCircle(g Gear) Circle     { return CircleOf(g, BaseCircleKind) }
func AddendumCircle(g Gear) Circle { return CircleOf(g, AddendumCircleKind) }
func DedendumCircle(g Gear) Circle { return CircleOf(g, DedendumCircleKind) }

// WorkingClearanceDenominator is a helper for the mesh package: the tip
// clearance c = a - (ra1 + rf2) needs the dedendum radius of the *other* gear,
// and this alias documents that the two gears deliberately use different
// circles in that formula.
func WorkingClearanceDenominator(g Gear) float64 {
	return g.DedendumRadius()
}

// NearlyStandard reports whether both gears have zero profile shift, i.e. the
// pair is mounted at the standard centre distance with alpha' = alpha. The
// mesh package uses this flag to skip the involute solve and report the exact
// standard values.
func NearlyStandard(p *Pair) bool {
	return math.Abs(p.ShiftSum()) < angle.DefaultTol
}
