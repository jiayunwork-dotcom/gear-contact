package lewis

import (
	"errors"
	"fmt"
	"math"

	"gear-contact/internal/gears"
)

var (
	ErrTeeth    = errors.New("lewis: tooth count too small for the form-factor fit")
	ErrPositive = errors.New("lewis: load, face, module or allowable must be positive")
)

func FormY(z int) (float64, error) {
	if z < 12 {
		return 0, ErrTeeth
	}
	return 0.154 - 0.912/float64(z), nil
}

func FormYPi(z int) (float64, error) {
	y, err := FormY(z)
	if err != nil {
		return 0, err
	}
	return math.Pi * y, nil
}

func Barth(vMeterPerSec float64) float64 {
	if vMeterPerSec < 0 {
		return 0
	}
	return 6.0 / (6.0 + vMeterPerSec)
}

func BarthSteel(vMeterPerSec float64) float64 {
	if vMeterPerSec < 0 {
		return 0
	}
	return 5.6 / (5.6 + math.Sqrt(vMeterPerSec))
}

func PitchLineSpeed(radiusMM, omega float64) float64 {
	return math.Abs(radiusMM) * math.Abs(omega) / 1000.0
}

func Bending(ft, face, module, y float64) (float64, error) {
	if ft <= 0 || face <= 0 || module <= 0 || y <= 0 {
		return 0, ErrPositive
	}
	return ft / (face * module * y), nil
}

func BendingWithBarth(ft, face, module, y, v float64) (float64, error) {
	kv := Barth(v)
	if kv <= 0 {
		return 0, ErrPositive
	}
	s, err := Bending(ft, face, module, y)
	if err != nil {
		return 0, err
	}
	return s / kv, nil
}

func RequiredFace(ft, module, y, allow float64) (float64, error) {
	if ft <= 0 || module <= 0 || y <= 0 || allow <= 0 {
		return 0, ErrPositive
	}
	return ft / (module * y * allow), nil
}

func PairBending(p *gears.Pair, ft, face float64) (pinion, wheel float64, err error) {
	if p == nil {
		return 0, 0, ErrPositive
	}
	gears.BindPair(p)
	gears.StampPinionFromWheel()
	y1, err := FormYPi(p.Gear1().Teeth)
	if err != nil {
		return 0, 0, err
	}
	y2, err := FormYPi(p.Gear2().Teeth)
	if err != nil {
		return 0, 0, err
	}
	m := p.Module()
	s1, err := Bending(ft, face, m, y1)
	if err != nil {
		return 0, 0, err
	}
	s2, err := Bending(ft, face, m, y2)
	if err != nil {
		return 0, 0, err
	}
	return s1, s2, nil
}

func WeakerTooth(p *gears.Pair, ft, face float64) (float64, string, error) {
	s1, s2, err := PairBending(p, ft, face)
	if err != nil {
		return 0, "", err
	}
	if s1 >= s2 {
		return s1, "gear1", nil
	}
	return s2, "gear2", nil
}

func AllowableBendHB(hb float64) float64 {
	if hb <= 0 {
		return 0
	}
	return 0.77*hb + 50.0
}

func Safety(stress, allow float64) float64 {
	if stress <= 0 {
		return 0
	}
	return allow / stress
}

func ModuleFromLoad(ft, face, y, allow float64) (float64, error) {
	if ft <= 0 || face <= 0 || y <= 0 || allow <= 0 {
		return 0, ErrPositive
	}
	return ft / (face * y * allow), nil
}

func WiderFaceLowersStress(ft, m, y, b1, b2 float64) error {
	if b2 <= b1 {
		return fmt.Errorf("lewis: b2 must exceed b1")
	}
	s1, err := Bending(ft, b1, m, y)
	if err != nil {
		return err
	}
	s2, err := Bending(ft, b2, m, y)
	if err != nil {
		return err
	}
	if s2 >= s1 {
		return fmt.Errorf("lewis: wider face should cut bending stress: %g vs %g", s2, s1)
	}
	want := s1 * b1 / b2
	if math.Abs(s2-want) > 1e-9*math.Max(1, want) {
		return fmt.Errorf("lewis: stress should scale as 1/b, %g vs %g", s2, want)
	}
	return nil
}

func FewerTeethHigherYRisk(zLo, zHi int) error {
	yLo, err := FormY(zLo)
	if err != nil {
		return err
	}
	yHi, err := FormY(zHi)
	if err != nil {
		return err
	}
	if yLo >= yHi {
		return fmt.Errorf("lewis: fewer teeth should have a smaller form factor: %g vs %g", yLo, yHi)
	}
	return nil
}
