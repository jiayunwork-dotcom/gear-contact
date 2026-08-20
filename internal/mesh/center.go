package mesh

import (
	"math"

	"gear-contact/internal/angle"
)

// Centre distance geometry. The standard centre distance is the value at
// which the two reference pitch circles touch: a0 = m*(z1+z2)/2. For a pair
// with total profile shift the no-backlash requirement moves the gears apart
// or together, and the operating centre distance follows from the working
// pressure angle through the base circle relation a*cos(alpha') = a0*cos(alpha).
//
// The invariant a*cos(alpha') = a0*cos(alpha) is what makes the involute
// profiles stay tangent along the line of action; it is asserted by a test in
// this package.

// StandardCenterDistance returns a0 = m*(z1+z2)/2.
func StandardCenterDistance(module float64, z1, z2 int) float64 {
	return 0.5 * module * float64(z1+z2)
}

// OperatingCenterDistance returns the centre distance that satisfies the base
// circle invariant for the given working pressure angle:
// a = a0*cos(alpha)/cos(alpha'). For alpha' = alpha the result is a0 exactly,
// which is the standard mounting case.
func OperatingCenterDistance(a0, alpha, alphaP float64) float64 {
	cosA := math.Cos(alpha)
	cosAP := math.Cos(alphaP)
	if cosAP == 0 {
		return math.Inf(1)
	}
	return a0 * cosA / cosAP
}

// OperatingCenterInverse recovers the standard centre distance a0 from an
// operating value: a0 = a*cos(alpha')/cos(alpha). The mesh code does not need
// it, but the tests use it to verify the invariant both ways.
func OperatingCenterInverse(a, alpha, alphaP float64) float64 {
	cosA := math.Cos(alpha)
	if cosA == 0 {
		return math.Inf(1)
	}
	return a * math.Cos(alphaP) / cosA
}

// CentreDistanceRatio returns a/a0. It is > 1 when the gears are spread apart
// (positive total shift), < 1 when they are pulled together, and exactly 1 at
// standard mounting.
func CentreDistanceRatio(a, a0 float64) float64 {
	return angle.Ratio(a, a0)
}
