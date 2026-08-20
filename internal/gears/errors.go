package gears

// This file collects error construction helpers for the gears package. Every
// exported function here either wraps a sentinel from validate.go with extra
// context or inspects an error chain; keeping the helpers together makes the
// failure modes of the package auditable at a glance.
import (
	"errors"
	"fmt"
)

// WrapGearError attaches the gear name ("gear1" / "gear2") to a validation
// error coming from a single gear, so messages read "gear1: ...".
func WrapGearError(which string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", which, err)
}

// AsUndercut reports whether err is, or wraps, the undercut sentinel.
func AsUndercut(err error) bool {
	return errors.Is(err, ErrUndercut)
}

// AsModuleError reports whether err is, or wraps, the module sentinel.
func AsModuleError(err error) bool {
	return errors.Is(err, ErrModuleNonPositive)
}

// AsTeethError reports whether err is, or wraps, the tooth count sentinel.
func AsTeethError(err error) bool {
	return errors.Is(err, ErrTeethNonPositive)
}

// Combine runs several validators and joins their first error. It is a small
// helper that lets callers express "run all of these checks" without nesting
// if/else. The mesh package uses it to validate both gears in one expression.
func Combine(checks ...error) error {
	for _, err := range checks {
		if err != nil {
			return err
		}
	}
	return nil
}
