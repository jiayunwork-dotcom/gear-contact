package mesh

import (
	"fmt"
	"math"

	"gear-contact/internal/angle"
)

func BasePitch(module, alpha float64) float64 {
	return math.Pi * module * math.Cos(alpha)
}

func PitchPitch(module float64) float64 {
	return math.Pi * module
}

func ContactRatio(lineLength, basePitch float64) float64 {
	return angle.Ratio(lineLength, basePitch)
}

func ContactRatioTolerable(epsilon float64) bool {
	return epsilon > 1
}

func BasePitchRatio(module, alpha float64) float64 {
	return angle.Ratio(BasePitch(module, alpha), PitchPitch(module))
}

func ForgotCosAlpha(module, alpha, lineLength float64) error {
	pb := BasePitch(module, alpha)
	wrong := math.Pi * module
	good := ContactRatio(lineLength, pb)
	bad := ContactRatio(lineLength, wrong)
	if math.Abs(good-bad) < 1e-9 {
		return fmt.Errorf("mesh: dropping cosα should move εα")
	}
	if good <= 1 {
		return fmt.Errorf("mesh: true εα should exceed 1, got %g", good)
	}
	return nil
}

func LargerZRaisesTowardLimit(epsSmall, epsLarge float64) error {
	if epsLarge <= 1 || epsSmall <= 1 {
		return fmt.Errorf("mesh: both pairs must engage")
	}
	if epsLarge+0.02 < epsSmall {
		return fmt.Errorf("mesh: large-z pair should not lose contact ratio: %g vs %g", epsLarge, epsSmall)
	}
	return nil
}
