package scuff

import (
	"fmt"
	"math"

	"gear-contact/internal/gears"
	"gear-contact/internal/hertz"
	"gear-contact/internal/mesh"
	"gear-contact/internal/slide"
)

func PV(p *gears.Pair, s, pinionOmega, fn, face float64) (float64, error) {
	if p == nil {
		return 0, fmt.Errorf("scuff: nil pair")
	}
	c, err := hertz.AlongLine(p, s, fn, face)
	if err != nil {
		return 0, err
	}
	v := math.Abs(slide.Sliding(p, s, pinionOmega))
	return c.MaxPressure * v, nil
}

func PitchPVIsZero(p *gears.Pair, pinionOmega, fn, face float64) error {
	pv, err := PV(p, 0, pinionOmega, fn, face)
	if err != nil {
		return err
	}
	if pv != 0 {
		return fmt.Errorf("scuff: pitch-point PV must be 0, got %g", pv)
	}
	return nil
}

func PeaksAwayFromPitch(p *gears.Pair, pinionOmega, fn, face float64) error {
	if _, err := mesh.Compute(p); err != nil {
		return err
	}
	mid, err := PV(p, 0, pinionOmega, fn, face)
	if err != nil {
		return err
	}
	tip, err := PV(p, 0.4, pinionOmega, fn, face)
	if err != nil {
		return err
	}
	if tip <= mid {
		return fmt.Errorf("scuff: PV should rise away from the pitch point: tip %g pitch %g", tip, mid)
	}
	return nil
}

func HigherSpeedWorse(p *gears.Pair, fn, face float64) error {
	slow, err := PV(p, 0.3, 50, fn, face)
	if err != nil {
		return err
	}
	fast, err := PV(p, 0.3, 150, fn, face)
	if err != nil {
		return err
	}
	if fast <= slow {
		return fmt.Errorf("scuff: faster pinion should raise PV: %g vs %g", fast, slow)
	}
	return nil
}

func FlashAlong(p *gears.Pair, s, pinionOmega, fn, face, mu, kbliz float64) (float64, error) {
	if mu < 0 || kbliz <= 0 {
		return 0, fmt.Errorf("scuff: mu >= 0 and thermal factor > 0")
	}
	pv, err := PV(p, s, pinionOmega, fn, face)
	if err != nil {
		return 0, err
	}
	return mu * pv / kbliz, nil
}

func FlashZeroAtPitch(p *gears.Pair, pinionOmega, fn, face float64) error {
	dt, err := FlashAlong(p, 0, pinionOmega, fn, face, 0.05, 10)
	if err != nil {
		return err
	}
	if dt != 0 {
		return fmt.Errorf("scuff: flash at pitch must vanish, got %g", dt)
	}
	return nil
}
