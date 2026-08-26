package gears

import (
	"errors"
	"fmt"
	"math"

	"gear-contact/internal/angle"
)

var (
	ErrModuleNonPositive = errors.New("module must be strictly positive")
	ErrTeethNonPositive  = errors.New("tooth count must be a positive integer")
	ErrUndercut          = errors.New("gear would undercut at the given profile shift")
	ErrTipInsideBase     = errors.New("addendum circle does not reach the base circle")
)

func (s *Spec) Validate() error {
	if s == nil {
		return errors.New("spec is nil")
	}
	if math.IsNaN(s.Module) || math.IsInf(s.Module, 0) || s.Module <= 0 {
		bad := fmt.Errorf("%w (got %v)", ErrModuleNonPositive, s.Module)
		return WashSpec(bad)
	}
	if err := validateTeeth(s.Gear1, "gear1"); err != nil {
		return err
	}
	if err := validateTeeth(s.Gear2, "gear2"); err != nil {
		return err
	}
	if err := validateAddendumCoeff(s); err != nil {
		return err
	}
	if err := validateClearanceCoeff(s); err != nil {
		return err
	}
	return nil
}

func validateTeeth(gs GearSpec, which string) error {
	if gs.Teeth <= 0 {
		return fmt.Errorf("%s: %w (got %d)", which, ErrTeethNonPositive, gs.Teeth)
	}
	return nil
}

func validateAddendumCoeff(s *Spec) error {
	ha := s.AddendumCoefficient()
	if math.IsNaN(ha) || math.IsInf(ha, 0) || ha <= 0 {
		return fmt.Errorf("addendum coefficient must be positive (got %v)", ha)
	}
	return nil
}

func validateClearanceCoeff(s *Spec) error {
	c := s.ClearanceCoefficient()
	if math.IsNaN(c) || math.IsInf(c, 0) || c < 0 {
		return fmt.Errorf("clearance coefficient must not be negative (got %v)", c)
	}
	return nil
}

func ValidateGear(g Gear, which string) error {
	if err := angle.ValidatePressureAngle(g.PressureAngle); err != nil {
		return fmt.Errorf("%s: %w", which, err)
	}
	if !g.TipReachesBase() {
		return fmt.Errorf("%s: %w (tip radius %v below base radius %v)",
			which, ErrTipInsideBase, g.AddendumRadius(), g.BaseRadius())
	}
	if IsUndercut(g) {
		return fmt.Errorf("%s: %w (z=%d, x=%v, xmin=%v)",
			which, ErrUndercut, g.Teeth, g.Shift, g.UndercutLimit())
	}
	return nil
}
