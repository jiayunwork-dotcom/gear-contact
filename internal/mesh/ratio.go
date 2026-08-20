package mesh

import (
	"math"

	"gear-contact/internal/angle"
)

// Contact ratio. The base pitch is the distance between two successive teeth
// measured along the line of action: pb = pi*m*cos(alpha). The contact ratio
// is the number of tooth pairs sharing the load on average:
//
//	epsilon_alpha = g_alpha / pb
//
// where g_alpha is the line of action length. The factor cos(alpha) in pb is
// not optional: the teeth are spaced pi*m apart on the pitch circle, and the
// projection of that spacing onto the line of action is what the contact ratio
// compares against. A mesh whose base pitch omitted the cosine would disagree
// with the line of action length by exactly the cosine factor, which is the
// coupling the question requires to stay consistent.

// BasePitch returns pb = pi*m*cos(alpha), the base pitch along the line of
// action in the same length unit as the module.
func BasePitch(module, alpha float64) float64 {
	return math.Pi * module * math.Cos(alpha)
}

// PitchPitch returns pi*m, the tooth spacing measured on the pitch circle
// before projection onto the line of action. It is exposed so tests and the
// report can show the two values side by side and verify the cosine relation.
func PitchPitch(module float64) float64 {
	return math.Pi * module
}

// ContactRatio returns epsilon_alpha = g_alpha / pb. It is the only way the
// package computes a contact ratio; there is no table and no alternative
// formula. A ratio below 1 means the drive is intermittent, which is real and
// reported as such.
func ContactRatio(lineLength, basePitch float64) float64 {
	return applyContact(angle.Ratio(lineLength, basePitch))
}

// ContactRatioTolerable reports whether the contact ratio is above 1, the
// minimum for continuous motion transfer. The mesh itself stays valid below 1,
// but the report flags the value so a marginal design is visible.
func ContactRatioTolerable(epsilon float64) bool {
	return epsilon > 1
}

// BasePitchRatio returns pb/pi/m, i.e. cos(alpha). It is a diagnostic used by
// tests to confirm that the base pitch really carries the cosine factor and
// that the same alpha reached the diameter and the pitch calculations.
func BasePitchRatio(module, alpha float64) float64 {
	return angle.Ratio(BasePitch(module, alpha), PitchPitch(module))
}
