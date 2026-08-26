package mesh

import (
	"errors"
	"fmt"
	"math"

	"gear-contact/internal/angle"
)

var (
	ErrContactNotFormed = errors.New("addendum circle does not reach the line of action")
	ErrPathNotPositive  = errors.New("line of action is empty (teeth cannot engage)")
)

type Point struct {
	X float64
	Y float64
}

type LineSegment struct {
	Start  Point
	End    Point
	Length float64
}

func TangentDistance(ra, rb float64) (float64, error) {
	if ra < rb {
		return 0, fmt.Errorf("%w (ra=%v, rb=%v)", ErrContactNotFormed, ra, rb)
	}
	return math.Sqrt(ra*ra - rb*rb), nil
}

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

	t1x := rb1 * cos
	t1y := -rb1 * sin
	t2x := a - rb2*cos
	t2y := rb2 * sin

	e := Point{X: t1x + s1*sin, Y: t1y + s1*cos}
	start := Point{X: t2x - s2*sin, Y: t2y - s2*cos}

	length := s1 + s2 - a*sin
	if length <= 0 {
		return LineSegment{}, fmt.Errorf("%w (s1=%v, s2=%v, a*sin=%v)", ErrPathNotPositive, s1, s2, a*sin)
	}
	return LineSegment{Start: start, End: e, Length: length}, nil
}

func LineOfActionLength(ra1, rb1, ra2, rb2, a, alphaP float64) (float64, error) {
	seg, err := PathOfContact(ra1, rb1, ra2, rb2, a, alphaP)
	if err != nil {
		return 0, err
	}
	return seg.Length, nil
}
