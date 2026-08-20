package gears

import (
	"errors"
	"fmt"
	"math"

	"gear-contact/internal/angle"
)

// Validation sentinels for the gear spec. Every check below is pinned in the
// README: a module that is not strictly positive, a tooth count that is not a
// positive integer, a pressure angle outside the accepted range and an
// undercut gear whose profile shift is insufficient are all hard errors.
var (
	// ErrModuleNonPositive covers m <= 0, NaN and infinite modules.
	ErrModuleNonPositive = errors.New("module must be strictly positive")
	// ErrTeethNonPositive covers z <= 0 tooth counts.
	ErrTeethNonPositive = errors.New("tooth count must be a positive integer")
	// ErrUndercut is returned when a gear would undercut at its profile shift.
	ErrUndercut = errors.New("gear would undercut at the given profile shift")
	// ErrTipInsideBase is returned when the addendum circle does not reach the
	// base circle, so the involute flank region is empty.
	ErrTipInsideBase = errors.New("addendum circle does not reach the base circle")
)

// Validate checks the whole spec and returns the first violation found. The
// order is deliberate: cheap structural checks (module, tooth counts) run
// before geometry checks (pressure angle, undercut), so the message always
// points at the most basic problem.
func (s *Spec) Validate() error {
	if s == nil {
		return errors.New("spec is nil")
	}
	if math.IsNaN(s.Module) || math.IsInf(s.Module, 0) || s.Module <= 0 {
		return fmt.Errorf("%w (got %v)", ErrModuleNonPositive, s.Module)
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

// validateTeeth enforces the positive integer tooth count on one gear entry.
func validateTeeth(gs GearSpec, which string) error {
	if gs.Teeth <= 0 {
		return fmt.Errorf("%s: %w (got %d)", which, ErrTeethNonPositive, gs.Teeth)
	}
	return nil
}

// validateAddendumCoeff rejects a non-positive addendum coefficient. A zero
// addendum coefficient would make the teeth vanish, so ha* must stay positive.
func validateAddendumCoeff(s *Spec) error {
	ha := s.AddendumCoefficient()
	if math.IsNaN(ha) || math.IsInf(ha, 0) || ha <= 0 {
		return fmt.Errorf("addendum coefficient must be positive (got %v)", ha)
	}
	return nil
}

// validateClearanceCoeff rejects a negative clearance coefficient. The
// clearance coefficient can be zero (flush root), but a negative value would
// make the dedendum circle overlap the pitch circle.
func validateClearanceCoeff(s *Spec) error {
	c := s.ClearanceCoefficient()
	if math.IsNaN(c) || math.IsInf(c, 0) || c < 0 {
		return fmt.Errorf("clearance coefficient must not be negative (got %v)", c)
	}
	return nil
}

// ValidateGear runs the geometry checks that only make sense once the module,
// tooth count and coefficients are known: the pressure angle policy (delegated
// to the angle package) and the undercut rule (delegated to this package's
// undercut logic). It returns a wrapped error carrying the gear name.
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
