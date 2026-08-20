package mesh

import (
	"fmt"
	"math"
)

// Inspection helpers for the line of action. These functions do not compute
// new mesh quantities; they verify that the points and lengths already
// derived satisfy the geometry they claim to. The tests call them so an error
// in the coordinate frame shows up as a failed invariant instead of a wrong
// contact ratio.

// VerifyEndpoints checks that the two endpoints of the segment really lie on
// the addendum circles they are supposed to lie on, in the frame where O1 is
// the origin and O2 is at (a, 0):
//
//   - End must satisfy |End - O1| = ra1 (on gear 1's addendum circle),
//   - Start must satisfy |Start - O2| = ra2 (on gear 2's addendum circle).
//
// The check is exact up to a relative tolerance because the endpoints are
// derived from the same radii.
func (seg LineSegment) VerifyEndpoints(ra1, ra2, a float64) error {
	distStartToO2 := math.Hypot(seg.Start.X-a, seg.Start.Y)
	distEndToO1 := math.Hypot(seg.End.X, seg.End.Y)
	if !relWithin(distStartToO2, ra2) {
		return fmt.Errorf("segment start is not on gear 2 addendum circle: distance %v, want %v", distStartToO2, ra2)
	}
	if !relWithin(distEndToO1, ra1) {
		return fmt.Errorf("segment end is not on gear 1 addendum circle: distance %v, want %v", distEndToO1, ra1)
	}
	return nil
}

// VerifyTangentDistance checks that the line of action is tangent to both
// base circles: the perpendicular distance from O1 (respectively O2) to the
// line through the segment equals rb1 (respectively rb2). The line has unit
// normal n = (-cos(alpha'), sin(alpha')) and its constant term c is read from
// the segment's Start point, which must lie on the line; the two distances are
// then compared with the two base radii.
func (seg LineSegment) VerifyTangentDistance(rb1, rb2, a, alphaP float64) error {
	sin, cos := sinCos(alphaP)
	nx, ny := -cos, sin
	c := nx*seg.Start.X + ny*seg.Start.Y
	distO1 := math.Abs(nx*0 + ny*0 - c)
	distO2 := math.Abs(nx*a + ny*0 - c)
	if !relWithin(distO1, rb1) {
		return fmt.Errorf("line of action is not tangent to base circle 1: distance %v, want %v", distO1, rb1)
	}
	if !relWithin(distO2, rb2) {
		return fmt.Errorf("line of action is not tangent to base circle 2: distance %v, want %v", distO2, rb2)
	}
	return nil
}

// relWithin reports whether got equals want within a relative tolerance.
func relWithin(got, want float64) bool {
	if want == 0 {
		return got == 0
	}
	return math.Abs(got-want)/math.Abs(want) < 1e-9
}

// sinCos is a local two-value sincos to avoid importing the angle package in
// the inspection path.
func sinCos(rad float64) (sin, cos float64) {
	sin = math.Sin(rad)
	cos = math.Cos(rad)
	return sin, cos
}
