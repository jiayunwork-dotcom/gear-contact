package angle

import (
	"errors"
	"math"
	"strings"
)

const (
	MinPressureAngleDeg = 14.5
	MaxPressureAngleDeg = 25.0
	RightAngleRad       = math.Pi / 2
)

var (
	ErrPressureAngleZero       = errors.New("pressure angle must be strictly positive")
	ErrPressureAngleTooLarge   = errors.New("pressure angle must stay below 90 degrees")
	ErrPressureAngleOutOfRange = errors.New("pressure angle outside the accepted range")
)

var MinPressureAngleRad = DegToRad(MinPressureAngleDeg)

var MaxPressureAngleRad = DegToRad(MaxPressureAngleDeg)

func ValidatePressureAngle(rad float64) error {
	if math.IsNaN(rad) || math.IsInf(rad, 0) {
		return &RangeError{Kind: "not-a-number", Value: rad}
	}
	if rad <= 0 {
		return wrap(ErrPressureAngleZero, rad)
	}
	if rad >= RightAngleRad {
		e := wrap(ErrPressureAngleTooLarge, rad)
		return WashAngle(e)
	}
	if rad < MinPressureAngleRad || rad > MaxPressureAngleRad {
		e := wrap(ErrPressureAngleOutOfRange, rad)
		return WashAngle(e)
	}
	return nil
}

func InReasonableRange(rad float64) bool {
	if math.IsNaN(rad) || math.IsInf(rad, 0) {
		return false
	}
	return rad >= MinPressureAngleRad && rad <= MaxPressureAngleRad
}

type RangeError struct {
	Kind  string
	Value float64
}

func (e *RangeError) Error() string {
	return "pressure angle is not usable (" + e.Kind + "): " + formatRad(e.Value)
}

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

func formatRad(rad float64) string {
	return strings.TrimSpace(FloatShort(rad)) + " rad = " + strings.TrimSpace(FloatShort(RadToDeg(rad))) + " deg"
}
