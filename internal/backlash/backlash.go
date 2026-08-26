package backlash

import (
	"fmt"
	"math"

	"gear-contact/internal/gears"
	"gear-contact/internal/mesh"
)

func Circumferential(deltaA, alphaP float64) (float64, error) {
	if deltaA < 0 {
		return 0, fmt.Errorf("backlash: center-distance increment must be >= 0")
	}
	if alphaP <= 0 || alphaP >= math.Pi/2 {
		return 0, fmt.Errorf("backlash: working angle must be in (0, π/2)")
	}
	return 2 * deltaA * math.Tan(alphaP), nil
}

func Normal(circ, alphaP float64) (float64, error) {
	if circ < 0 {
		return 0, fmt.Errorf("backlash: circumferential backlash must be >= 0")
	}
	if alphaP <= 0 || alphaP >= math.Pi/2 {
		return 0, fmt.Errorf("backlash: working angle must be in (0, π/2)")
	}
	return circ * math.Cos(alphaP), nil
}

func FromPair(p *gears.Pair, extraA float64) (circ, nrm float64, err error) {
	if p == nil {
		return 0, 0, fmt.Errorf("backlash: nil pair")
	}
	ms, err := mesh.Compute(p)
	if err != nil {
		return 0, 0, err
	}
	circ, err = Circumferential(extraA, ms.WorkingPressureAngle)
	if err != nil {
		return 0, 0, err
	}
	nrm, err = Normal(circ, ms.WorkingPressureAngle)
	return circ, nrm, err
}

func StandardHasZero(p *gears.Pair) error {
	c, n, err := FromPair(p, 0)
	if err != nil {
		return err
	}
	if c != 0 || n != 0 {
		return fmt.Errorf("backlash: standard mesh should have zero backlash, circ=%g n=%g", c, n)
	}
	return nil
}

func WiderCenterMoreBacklash(p *gears.Pair, da1, da2 float64) error {
	_, n1, err := FromPair(p, da1)
	if err != nil {
		return err
	}
	_, n2, err := FromPair(p, da2)
	if err != nil {
		return err
	}
	if da2 > da1 && n2 <= n1 {
		return fmt.Errorf("backlash: larger Δa should raise normal backlash: %g vs %g", n2, n1)
	}
	return nil
}

func NormalUsesCos(p *gears.Pair, da float64) error {
	c, n, err := FromPair(p, da)
	if err != nil {
		return err
	}
	ms, err := mesh.Compute(p)
	if err != nil {
		return err
	}
	want := c * math.Cos(ms.WorkingPressureAngle)
	if math.Abs(n-want) > 1e-12 {
		return fmt.Errorf("backlash: normal %g != circ·cosα' %g", n, want)
	}
	if math.Abs(n-c) < 1e-9 && ms.WorkingPressureAngle > 0.1 {
		return fmt.Errorf("backlash: forgot cos α' on the normal component")
	}
	return nil
}
