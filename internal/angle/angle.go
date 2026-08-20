// Package angle holds the angular mathematics shared by every stage of the
// gear mesh calculation. All geometry inside the tool works in radians; the
// JSON input and the text report use degrees, so conversion lives here too.
//
// The package owns three responsibilities:
//   - degree <-> radian conversion and angle formatting,
//   - the pressure angle policy (what range is acceptable, what is rejected),
//   - the involute function inv(x)=tan(x)-x and the numerical inverse that
//     the profile-shift centre-distance calculation needs.
//
// Everything in this package is pure: no state, no I/O, only float64 math
// with explicit errors for inputs that cannot be handled.
package angle

import "math"

// TwoPi is a full turn in radians, used by angle normalization helpers.
const TwoPi = 2 * math.Pi

// DegToRad converts an angle in degrees to radians.
//
// The conversion is identity-safe: DegToRad(0) is 0 and DegToRad(90) is
// exactly pi/2 in IEEE-754 arithmetic, so callers can compare against the
// right-angle limit without worrying about round-off.
func DegToRad(deg float64) float64 {
	return deg * math.Pi / 180.0
}

// RadToDeg converts an angle in radians to degrees.
func RadToDeg(rad float64) float64 {
	return rad * 180.0 / math.Pi
}

// Normalize wraps an angle in radians into the interval [0, 2*pi). It is used
// when a caller receives an angle from a formula (for example the inverse
// involute solve) and wants to display it without a sign ambiguity.
func Normalize(rad float64) float64 {
	rad = math.Mod(rad, TwoPi)
	if rad < 0 {
		rad += TwoPi
	}
	return rad
}

// SinCos returns the sine and cosine of an angle in radians in one call.
// Keeping the pair together avoids the classic mistake of computing the sine
// from a different angle than the cosine; the mesh code uses the two values
// together to build the line of action.
func SinCos(rad float64) (sin, cos float64) {
	return math.Sincos(rad)
}
