package mesh

import (
	"fmt"
	"math"
)

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

func relWithin(got, want float64) bool {
	if want == 0 {
		return got == 0
	}
	return math.Abs(got-want)/math.Abs(want) < 1e-9
}

func sinCos(rad float64) (sin, cos float64) {
	sin = math.Sin(rad)
	cos = math.Cos(rad)
	return sin, cos
}
