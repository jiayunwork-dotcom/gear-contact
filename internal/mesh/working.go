package mesh

import (
	"errors"
	"math"

	"gear-contact/internal/angle"
)

// Working pressure angle. For a standard pair (x1 + x2 = 0) the working
// pressure angle equals the standard one and the gears run at a0. For a
// profile-shifted pair the no-backlash condition is
//
//	inv(alpha') = inv(alpha) + 2*(x1+x2)*tan(alpha)/(z1+z2)
//
// which is solved numerically by the angle package. The equation is the same
// one used to derive the operating centre distance, so alpha' and a are
// mutually consistent by construction.

var (
	// ErrInvalidTeethSum is returned when the tooth count sum is not positive.
	ErrInvalidTeethSum = errors.New("tooth count sum must be positive for the involute equation")
	// ErrWorkingAngle wraps any failure of the inverse involute solve.
	ErrWorkingAngle = errors.New("could not solve the working pressure angle")
)

// WorkingPressureAngle returns the no-backlash working pressure angle in
// radians. With zero total profile shift the exact standard angle is returned
// without iterating, so a standard mesh never carries solver round-off.
func WorkingPressureAngle(alpha, shiftSum float64, teethSum int) (float64, error) {
	if teethSum <= 0 {
		return 0, ErrInvalidTeethSum
	}
	if shiftSum == 0 {
		return alpha, nil
	}
	target := angle.Inv(alpha) + 2*shiftSum*math.Tan(alpha)/float64(teethSum)
	alphaP, err := angle.InvSolve(target, alpha)
	if err != nil {
		return 0, wrapWorking(err)
	}
	return alphaP, nil
}

// WorkingPressureAngleDeg is WorkingPressureAngle for callers that work in
// degrees; it returns the angle and its degree equivalent together so no one
// re-converts and drifts.
func WorkingPressureAngleDeg(alphaDeg, shiftSum float64, teethSum int) (rad, deg float64, err error) {
	rad, err = WorkingPressureAngle(angle.DegToRad(alphaDeg), shiftSum, teethSum)
	if err != nil {
		return 0, 0, err
	}
	return rad, angle.RadToDeg(rad), nil
}

// wrapWorking attaches the working-angle sentinel to a solver error.
func wrapWorking(cause error) error {
	return errors.Join(cause, ErrWorkingAngle)
}
