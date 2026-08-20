package angle

import "errors"

// This file centralizes every error value the angle package can produce.
// Keeping the sentinels in one place makes the error contract of the package
// easy to review: a caller that imports angle and matches on these values can
// route each failure to the right message without string matching.
var (
	// ErrDegToRadDegenerate is returned by helpers that refuse a zero angle
	// where the geometry is undefined. It is kept separate from the pressure
	// angle policy errors because some callers legitimately pass zero (for
	// example when comparing a standard mesh) and should be able to tell the
	// two situations apart.
	ErrDegToRadDegenerate = errors.New("angle conversion of zero is undefined in this context")

	// ErrNormalizeFailed is a defensive error returned when an angle cannot be
	// brought into [0, 2*pi), which can only happen for NaN or infinite input.
	ErrNormalizeFailed = errors.New("cannot normalize a non-finite angle")

	// ErrDomain is the umbrella sentinel for values that leave the angle
	// domain. It wraps ErrPressureAngleZero / TooLarge / OutOfRange through
	// errors.Join so a single errors.Is(_, ErrDomain) check is enough to
	// detect any domain violation.
	ErrDomain = errors.New("angle outside the valid domain")
)

// WrapDomain attaches the umbrella domain sentinel to a pressure angle error
// so callers that only care about "some angle problem" can match once.
func WrapDomain(cause error) error {
	if cause == nil {
		return nil
	}
	return errors.Join(cause, ErrDomain)
}

// IsDomain reports whether the error chain contains the domain sentinel.
func IsDomain(err error) bool {
	return errors.Is(err, ErrDomain)
}
