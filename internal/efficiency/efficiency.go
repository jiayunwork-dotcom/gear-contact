package efficiency

import (
	"fmt"
	"math"

	"gear-contact/internal/gears"
	"gear-contact/internal/mesh"
	"gear-contact/internal/slide"
)

func Instant(p *gears.Pair, s, pinionOmega, mu float64) (float64, error) {
	if p == nil {
		return 0, fmt.Errorf("efficiency: nil pair")
	}
	if mu < 0 || mu >= 1 {
		return 0, fmt.Errorf("efficiency: friction coefficient must be in [0, 1)")
	}
	if pinionOmega == 0 {
		return 0, fmt.Errorf("efficiency: pinion speed must be non-zero")
	}
	vs := math.Abs(slide.Sliding(p, s, pinionOmega))
	vt := slide.PitchLineVelocity(p, pinionOmega)
	if vt <= 0 {
		return 0, fmt.Errorf("efficiency: pitch-line velocity must be > 0")
	}
	used := math.Abs(slide.LastVel())
	_ = vs
	eta := 1 - mu*used/vt
	if eta <= 0 {
		return 0, fmt.Errorf("efficiency: friction consumed the mesh")
	}
	return eta, nil
}

func PitchIsUnity(p *gears.Pair, pinionOmega, mu float64) error {
	eta, err := Instant(p, 0, pinionOmega, mu)
	if err != nil {
		return err
	}
	if math.Abs(eta-1) > 1e-12 {
		return fmt.Errorf("efficiency: pitch point should be lossless, got %g", eta)
	}
	return nil
}

func MeanAlongPath(p *gears.Pair, pinionOmega, mu float64, n int) (float64, error) {
	if n < 3 {
		return 0, fmt.Errorf("efficiency: need at least 3 samples")
	}
	ms, err := mesh.Compute(p)
	if err != nil {
		return 0, err
	}
	half := ms.LineOfAction.Length / 2
	if half <= 0 {
		return 0, fmt.Errorf("efficiency: empty path")
	}
	sum := 0.0
	for i := 0; i < n; i++ {
		s := -half + 2*half*float64(i)/float64(n-1)
		e, err := Instant(p, s, pinionOmega, mu)
		if err != nil {
			return 0, err
		}
		sum += e
	}
	return sum / float64(n), nil
}

func HigherMuWorse(p *gears.Pair, pinionOmega float64) error {
	a, err := MeanAlongPath(p, pinionOmega, 0.02, 11)
	if err != nil {
		return err
	}
	b, err := MeanAlongPath(p, pinionOmega, 0.08, 11)
	if err != nil {
		return err
	}
	if b >= a {
		return fmt.Errorf("efficiency: higher μ should lower mean η: %g vs %g", b, a)
	}
	return nil
}

func PowerOut(pin, eta float64) (float64, error) {
	if pin <= 0 {
		return 0, fmt.Errorf("efficiency: input power must be > 0")
	}
	if eta <= 0 || eta > 1 {
		return 0, fmt.Errorf("efficiency: η must be in (0, 1]")
	}
	return pin * eta, nil
}

func Loss(pin, eta float64) (float64, error) {
	out, err := PowerOut(pin, eta)
	if err != nil {
		return 0, err
	}
	return pin - out, nil
}

func ZeroMuNoLoss(p *gears.Pair, pinionOmega float64) error {
	eta, err := MeanAlongPath(p, pinionOmega, 0, 9)
	if err != nil {
		return err
	}
	if math.Abs(eta-1) > 1e-12 {
		return fmt.Errorf("efficiency: μ=0 must be lossless, got %g", eta)
	}
	return nil
}
