package tiprel

import (
	"fmt"
	"math"

	"gear-contact/internal/gears"
	"gear-contact/internal/mesh"
)

func WithAddendum(p *gears.Pair, ha float64) (*gears.Pair, error) {
	if p == nil {
		return nil, fmt.Errorf("tiprel: nil pair")
	}
	if !(ha > 0) {
		return nil, fmt.Errorf("tiprel: addendum coefficient must be > 0")
	}
	g1 := p.Gear1()
	g2 := p.Gear2()
	spec := &gears.Spec{
		Module:           g1.Module,
		PressureAngleDeg: g1.PressureAngle * 180 / math.Pi,
		AddendumCoeff:    &ha,
		Gear1:            gears.GearSpec{Teeth: g1.Teeth, Shift: g1.Shift},
		Gear2:            gears.GearSpec{Teeth: g2.Teeth, Shift: g2.Shift},
	}
	return gears.New(spec)
}

func ContactRatioAt(p *gears.Pair, ha float64) (float64, error) {
	q, err := WithAddendum(p, ha)
	if err != nil {
		return 0, err
	}
	ms, err := mesh.Compute(q)
	if err != nil {
		return 0, err
	}
	return ms.ContactRatio, nil
}

func ShorterAddendumLowersEpsilon(p *gears.Pair, haFull, haCut float64) error {
	if haCut >= haFull {
		return fmt.Errorf("tiprel: cut addendum must be shorter")
	}
	full, err := ContactRatioAt(p, haFull)
	if err != nil {
		return err
	}
	cut, err := ContactRatioAt(p, haCut)
	if err != nil {
		return err
	}
	if cut >= full {
		return fmt.Errorf("tiprel: shorter addendum should cut εα: %g vs %g", cut, full)
	}
	return nil
}

func StillEngages(p *gears.Pair, ha float64) error {
	eps, err := ContactRatioAt(p, ha)
	if err != nil {
		return err
	}
	if eps <= 1 {
		return fmt.Errorf("tiprel: pair should still engage, εα=%g", eps)
	}
	return nil
}

func PathTracksAddendum(p *gears.Pair, haFull, haCut float64) error {
	a, err := WithAddendum(p, haFull)
	if err != nil {
		return err
	}
	b, err := WithAddendum(p, haCut)
	if err != nil {
		return err
	}
	ma, err := mesh.Compute(a)
	if err != nil {
		return err
	}
	mb, err := mesh.Compute(b)
	if err != nil {
		return err
	}
	if mb.LineOfAction.Length >= ma.LineOfAction.Length {
		return fmt.Errorf("tiprel: shorter addendum should shorten the path: %g vs %g", mb.LineOfAction.Length, ma.LineOfAction.Length)
	}
	if math.Abs(ma.BasePitch-mb.BasePitch) > 1e-12 {
		return fmt.Errorf("tiprel: base pitch must stay a module/α quantity")
	}
	return nil
}
