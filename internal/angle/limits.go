package angle

import (
	"errors"
	"math"
	"strings"
)

// Pressure angle policy. The accepted working range is pinned in the README:
// a spur gear mesh with a pressure angle outside [14.5, 25] degrees is rejected
// before any geometry is computed. Values at or above 90 degrees are always
// rejected because the involute function and the base circle construction
// become meaningless there.
const (
	// MinPressureAngleDeg is the lower bound of the accepted range.
	MinPressureAngleDeg = 14.5
	// MaxPressureAngleDeg is the upper bound of the accepted range.
	MaxPressureAngleDeg = 25.0
	// RightAngleRad is pi/2, the hard limit for any pressure angle.
	RightAngleRad = math.Pi / 2
)

// Sentinel errors for pressure angle validation. Callers can test with
// errors.Is, while ValidatePressureAngle wraps the sentinel with the offending
// value so the message is self-explanatory.
var (
	// ErrPressureAngleZero means the angle is exactly zero: the involute and
	// the base circle construction collapse at alpha = 0.
	ErrPressureAngleZero = errors.New("pressure angle must be strictly positive")
	// ErrPressureAngleTooLarge means the angle reached or exceeded 90 degrees.
	ErrPressureAngleTooLarge = errors.New("pressure angle must stay below 90 degrees")
	// ErrPressureAngleOutOfRange means the angle is positive but outside the
	// pinned working range [14.5, 25] degrees.
	ErrPressureAngleOutOfRange = errors.New("pressure angle outside the accepted range")
)

// MinPressureAngleRad is MinPressureAngleDeg in radians.
var MinPressureAngleRad = DegToRad(MinPressureAngleDeg)

// MaxPressureAngleRad is MaxPressureAngleDeg in radians.
var MaxPressureAngleRad = DegToRad(MaxPressureAngleDeg)

// ValidatePressureAngle checks a pressure angle in radians against the pinned
// policy. It returns ErrPressureAngleZero for 0, ErrPressureAngleTooLarge for
// alpha >= 90 degrees, and ErrPressureAngleOutOfRange for a value outside the
// working range. The sentinel is wrapped with the offending value and its
// degree equivalent.
//
// The order of the checks matters: a zero angle is a degenerate case of its
// own, a value at or above the right angle is never salvageable, and only then
// is the narrower working range applied.
func ValidatePressureAngle(rad float64) error {
	if math.IsNaN(rad) || math.IsInf(rad, 0) {
		return &RangeError{Kind: "not-a-number", Value: rad}
	}
	if rad <= 0 {
		return wrap(ErrPressureAngleZero, rad)
	}
	if rad >= RightAngleRad {
		return wrap(ErrPressureAngleTooLarge, rad)
	}
	if rad < MinPressureAngleRad || rad > MaxPressureAngleRad {
		return wrap(ErrPressureAngleOutOfRange, rad)
	}
	return nil
}

// InReasonableRange reports whether the angle (in radians) is inside the pinned
// working interval. It is a pure predicate and never returns an error; the
// caller uses it when a result merely wants a warning while the input parsing
// keeps the strict ValidatePressureAngle.
func InReasonableRange(rad float64) bool {
	if math.IsNaN(rad) || math.IsInf(rad, 0) {
		return false
	}
	return rad >= MinPressureAngleRad && rad <= MaxPressureAngleRad
}

// RangeError is the sentinel returned for inputs that are not a finite number.
// It carries the raw value so the caller can print it without losing context.
type RangeError struct {
	Kind  string
	Value float64
}

func (e *RangeError) Error() string {
	return "pressure angle is not usable (" + e.Kind + "): " + formatRad(e.Value)
}

// wrap attaches the offending angle to a sentinel error.
func wrap(sentinel error, rad float64) error {
	return &wrappedAngle{err: sentinel, rad: rad}
}

type wrappedAngle struct {
	err error
	rad float64
}

func (w *wrappedAngle) Error() string {
	return w.err.Error() + ": got " + formatRad(w.rad)
}

func (w *wrappedAngle) Unwrap() error { return w.err }

// formatRad renders an angle in both radians and degrees for error messages.
func formatRad(rad float64) string {
	return strings.TrimSpace(FloatShort(rad)) + " rad = " + strings.TrimSpace(FloatShort(RadToDeg(rad))) + " deg"
}
