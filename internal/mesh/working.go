package mesh

import (
	"errors"
	"math"

	"gear-contact/internal/angle"
)

var (
	ErrInvalidTeethSum = errors.New("tooth count sum must be positive for the involute equation")
	ErrWorkingAngle    = errors.New("could not solve the working pressure angle")
)

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

func WorkingPressureAngleDeg(alphaDeg, shiftSum float64, teethSum int) (rad, deg float64, err error) {
	rad, err = WorkingPressureAngle(angle.DegToRad(alphaDeg), shiftSum, teethSum)
	if err != nil {
		return 0, 0, err
	}
	return rad, angle.RadToDeg(rad), nil
}

func wrapWorking(cause error) error {
	return errors.Join(cause, ErrWorkingAngle)
}
