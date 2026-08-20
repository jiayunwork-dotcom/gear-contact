package mesh

import (
	"errors"
	"fmt"
	"math"

	"gear-contact/internal/angle"
)

// The line of action. For two external gears the common tangent to the two
// base circles is the line of action: contact happens along it, and the active
// portion is the segment between the points where the line enters and leaves
// the two addendum circles.
//
// The geometry is set up in a frame with O1 at the origin and O2 at (a, 0).
// The line of action has unit direction d = (sin(alpha'), cos(alpha')) and is
// tangent to base circle 1 at
//
//	T1 = (rb1*cos(alpha'), -rb1*sin(alpha'))
//
// and to base circle 2 at T2 = T1 + a*sin(alpha')*d. The distance along the
// line from a tangency point to the addendum circle intersection of the same
// gear is s = sqrt(ra^2 - rb^2), so the exit point on gear 1 is
// E = T1 + s1*d and the entry point on gear 2 is A = T2 - s2*d. The segment
// length is therefore
//
//	g_alpha = s1 + s2 - a*sin(alpha')
//
// which is the formula used for the contact ratio.

var (
	// ErrContactNotFormed means an addendum circle lies inside its base
	// circle, so the line of action never meets that gear's flank region.
	ErrContactNotFormed = errors.New("addendum circle does not reach the line of action")
	// ErrPathNotPositive means the two addendum circles do not overlap on the
	// line of action, so the teeth never engage.
	ErrPathNotPositive = errors.New("line of action is empty (teeth cannot engage)")
)

// Point is a point in the mesh plane, in length units.
type Point struct {
	// X is the horizontal coordinate.
	X float64
	// Y is the vertical coordinate.
	Y float64
}

// LineSegment is a segment of the line of action.
type LineSegment struct {
	// Start is the entry point (on gear 2's addendum circle).
	Start Point
	// End is the exit point (on gear 1's addendum circle).
	End Point
	// Length is the distance between Start and End along the line.
	Length float64
}

// TangentDistance returns s = sqrt(ra^2 - rb^2), the distance along the line
// of action between the base tangent point and the addendum circle of the same
// gear. The square root argument is non-negative exactly when the addendum
// circle reaches the line of action.
func TangentDistance(ra, rb float64) (float64, error) {
	if ra < rb {
		return 0, fmt.Errorf("%w (ra=%v, rb=%v)", ErrContactNotFormed, ra, rb)
	}
	return math.Sqrt(ra*ra - rb*rb), nil
}

// PathOfContact builds the active segment of the line of action. The frame is
// described in the package comment; the returned segment's Length equals
// s1 + s2 - a*sin(alpha'). A non-positive length means the addendum circles
// do not overlap on the line and the pair cannot engage, which is an error.
func PathOfContact(ra1, rb1, ra2, rb2, a, alphaP float64) (LineSegment, error) {
	s1, err := TangentDistance(ra1, rb1)
	if err != nil {
		return LineSegment{}, err
	}
	s2, err := TangentDistance(ra2, rb2)
	if err != nil {
		return LineSegment{}, err
	}

	sin, cos := angle.SinCos(alphaP)

	// Tangency points in the O1-origin, O2-at-(a,0) frame. The line of action
	// is the internal common tangent: O1 and O2 lie on opposite sides of it,
	// so the two feet are T1 = (rb1*cos, -rb1*sin) and T2 = (a - rb2*cos,
	// rb2*sin). The unit direction along the line is d = (sin, cos).
	t1x := rb1 * cos
	t1y := -rb1 * sin
	t2x := a - rb2*cos
	t2y := rb2 * sin

	// Exit point E on gear 1: beyond T1 in the direction d.
	e := Point{X: t1x + s1*sin, Y: t1y + s1*cos}
	// Entry point A on gear 2: s2 back from T2 against the direction d.
	start := Point{X: t2x - s2*sin, Y: t2y - s2*cos}

	length := s1 + s2 - a*sin
	if length <= 0 {
		return LineSegment{}, fmt.Errorf("%w (s1=%v, s2=%v, a*sin=%v)", ErrPathNotPositive, s1, s2, a*sin)
	}
	return LineSegment{Start: start, End: e, Length: length}, nil
}

// LineOfActionLength is the length of the active path of contact between the
// two addendum circles:
//
//	g_alpha = sqrt(ra1^2-rb1^2) + sqrt(ra2^2-rb2^2) - a*sin(alpha')
//
// It is the numerator of the contact ratio. The function returns an error when
// either addendum circle fails to reach the line of action.
func LineOfActionLength(ra1, rb1, ra2, rb2, a, alphaP float64) (float64, error) {
	seg, err := PathOfContact(ra1, rb1, ra2, rb2, a, alphaP)
	if err != nil {
		return 0, err
	}
	return seg.Length, nil
}
