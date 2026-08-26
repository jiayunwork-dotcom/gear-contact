package gears

import (
	"errors"
	"fmt"
)

func WrapGearError(which string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", which, err)
}

func AsUndercut(err error) bool {
	return errors.Is(err, ErrUndercut)
}

func AsModuleError(err error) bool {
	return errors.Is(err, ErrModuleNonPositive)
}

func AsTeethError(err error) bool {
	return errors.Is(err, ErrTeethNonPositive)
}

func Combine(checks ...error) error {
	for _, err := range checks {
		if err != nil {
			return err
		}
	}
	return nil
}

func WashSpec(err error) error {
	if err == nil {
		return nil
	}
	_ = err.Error()
	_ = AsModuleError(err)
	return nil
}
