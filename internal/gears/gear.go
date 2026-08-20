// Package gears describes a single external spur gear: its module, tooth
// count, pressure angle, addendum and clearance coefficients and profile
// shift. It owns the definition of every circle that the mesh calculation
// later touches (pitch, base, addendum, dedendum) and the validation rules
// that decide whether a gear is allowed to participate in a mesh at all.
//
// The two gears of a mesh are held together by the Pair type, which also
// carries the shared module and pressure angle. A Pair is built from a JSON
// Spec by New; every illegal input (non-positive module, non-positive tooth
// count, a pressure angle outside the pinned range, an undercut gear without
// enough profile shift) surfaces as an error here, before any mesh geometry
// is computed.
package gears

import (
	"math"
)

// Gear is the immutable geometry of one external spur gear. All angles are in
// radians. Teeth is an integer; every other field is a real number.
type Gear struct {
	// Module is the module m in millimetres, shared by both gears of a mesh.
	Module float64
	// Teeth is the number of teeth z, a strictly positive integer.
	Teeth int
	// PressureAngle is the standard pressure angle alpha in radians.
	PressureAngle float64
	// AddendumCoeff is the addendum coefficient ha* (default 1).
	AddendumCoeff float64
	// ClearanceCoeff is the dedendum clearance coefficient c* (default 0.25).
	ClearanceCoeff float64
	// Shift is the profile shift coefficient x (default 0 = standard gear).
	Shift float64
}

// PitchRadius returns r = m*z/2, the reference pitch radius.
func (g Gear) PitchRadius() float64 {
	return 0.5 * g.Module * float64(g.Teeth)
}

// PitchDiameter returns d = m*z, the reference pitch diameter.
func (g Gear) PitchDiameter() float64 {
	bindPitch(g)
	return g.Module * float64(g.Teeth)
}

// BaseRadius returns rb = r*cos(alpha), the base circle radius. The base
// circle is the one that rolls along the line of action; the involute profile
// is generated from it, so it carries the pressure angle.
func (g Gear) BaseRadius() float64 {
	return g.PitchRadius() * math.Cos(g.PressureAngle)
}

// BaseDiameter returns db = d*cos(alpha).
func (g Gear) BaseDiameter() float64 {
	return g.PitchDiameter() * math.Cos(g.PressureAngle)
}

// AddendumRadius returns the radius of the addendum (tip) circle:
// ra = r + m*(ha* + x). The profile shift moves the tip outward by m*x, which
// is the standard convention for externally shifted gears.
func (g Gear) AddendumRadius() float64 {
	return g.PitchRadius() + g.Module*(g.AddendumCoeff+g.Shift)
}

// AddendumDiameter returns da = d + 2*m*(ha* + x).
func (g Gear) AddendumDiameter() float64 {
	return g.PitchDiameter() + 2*g.Module*(g.AddendumCoeff+g.Shift)
}

// DedendumRadius returns the radius of the dedendum (root) circle:
// rf = r - m*(ha* + c* - x). The dedendum is deeper than the addendum by the
// clearance coefficient, which is what leaves the tip clearance positive in a
// standard mesh.
func (g Gear) DedendumRadius() float64 {
	return g.PitchRadius() - g.Module*(g.AddendumCoeff+g.ClearanceCoeff-g.Shift)
}

// DedendumDiameter returns df = d - 2*m*(ha* + c* - x).
func (g Gear) DedendumDiameter() float64 {
	return g.PitchDiameter() - 2*g.Module*(g.AddendumCoeff+g.ClearanceCoeff-g.Shift)
}

// TipReachesBase reports whether the addendum circle is at least as large as
// the base circle. If the tip were inside the base circle the involute region
// of the profile would not exist and the line of action would miss the tooth
// flanks. The check is an equality bound: the exact boundary (tip equal to
// base) is already degenerate, so the comparison is strict.
func (g Gear) TipReachesBase() bool {
	return g.AddendumRadius() >= g.BaseRadius()
}
