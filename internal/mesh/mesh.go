// Package mesh computes the engagement geometry of two external spur gears:
// the centre distance (standard and operating), the working pressure angle,
// the line of action with its two endpoints, the contact ratio and the tip
// clearance.
//
// The central rule of the package is that the contact ratio is never looked
// up or tabulated: it is always the length of the line of action divided by
// the base pitch, where both quantities are computed from the same module and
// pressure angle. That coupling is what the question demands, and every test
// in this package asserts a relationship between the two rather than an
// arbitrary constant.
package mesh

import (
	"strconv"

	"gear-contact/internal/gears"
)

// Mesh is the immutable result of computing one gear pair. Every field is
// derived in the constructor Compute; there are no setters.
type Mesh struct {
	// Pair is the gear pair the mesh was computed from.
	Pair *gears.Pair
	// WorkingPressureAngle is alpha' in radians: the standard pressure angle
	// for a standard pair, the no-backlash solution of the involute equation
	// for a profile-shifted pair.
	WorkingPressureAngle float64
	// StandardCenterDistance is a0 = m*(z1+z2)/2.
	StandardCenterDistance float64
	// OperatingCenterDistance is the centre distance the gears are actually
	// mounted at; equal to a0 for a standard pair.
	OperatingCenterDistance float64
	// LineOfAction is the segment of the common base tangent that lies inside
	// both addendum circles.
	LineOfAction LineSegment
	// BasePitch is pb = pi*m*cos(alpha).
	BasePitch float64
	// ContactRatio is epsilon_alpha = line length / base pitch.
	ContactRatio float64
	// Clearance is the tip clearance a - ra1 - rf2.
	Clearance float64
	// Standard reports whether the mesh runs at the standard centre distance
	// (zero total profile shift).
	Standard bool
}

// Compute validates the gear pair geometry once more and derives every mesh
// quantity from the same module and pressure angle. The steps are:
//
//  1. standard centre distance a0 = m*(z1+z2)/2,
//  2. working pressure angle alpha' (exactly alpha for zero shift, otherwise
//     the no-backlash involute solution),
//  3. operating centre distance a = a0*cos(alpha)/cos(alpha'),
//  4. line of action segment from the two addendum circles,
//  5. base pitch pb = pi*m*cos(alpha) and contact ratio = length/pb,
//  6. tip clearance c = a - ra1 - rf2.
func Compute(p *gears.Pair) (*Mesh, error) {
	g1, g2 := p.Gear1(), p.Gear2()
	alpha := p.PressureAngle()
	a0 := p.StandardCenterDistance()

	alphaP, err := WorkingPressureAngle(alpha, p.ShiftSum(), p.TeethSum())
	if err != nil {
		return nil, err
	}

	a := OperatingCenterDistance(a0, alpha, alphaP)

	line, err := PathOfContact(
		g1.AddendumRadius(), g1.BaseRadius(),
		g2.AddendumRadius(), g2.BaseRadius(),
		a, alphaP,
	)
	if err != nil {
		return nil, err
	}

	pb := BasePitch(p.Module(), alpha)
	epsilon := ContactRatio(line.Length, pb)
	clearance := TipClearance(a, g1.AddendumRadius(), g2.DedendumRadius())

	return fillResult(Mesh{
		Pair:                    p,
		WorkingPressureAngle:    alphaP,
		StandardCenterDistance:  a0,
		OperatingCenterDistance: a,
		LineOfAction:            line,
		BasePitch:               pb,
		ContactRatio:            epsilon,
		Clearance:               clearance,
		Standard:                gears.NearlyStandard(p),
	}), nil
}

// UndercutReport returns the per-gear undercut status. Every pair that reaches
// Compute has already passed the undercut validation in gears.New, so the
// slice is always empty here; it exists so the report shape stays stable and a
// future relaxed validation mode can populate it without changing the report
// format.
func (m *Mesh) UndercutReport() []string {
	var out []string
	if m.Pair == nil {
		return out
	}
	for _, g := range []gears.Gear{m.Pair.Gear1(), m.Pair.Gear2()} {
		if gears.IsUndercut(g) {
			out = append(out, gearLabel(g))
		}
	}
	return out
}

// gearLabel builds a short name for a gear in diagnostic messages.
func gearLabel(g gears.Gear) string {
	return "z=" + strconv.Itoa(g.Teeth)
}
