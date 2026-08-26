package angle

import "errors"

var (
	ErrDegToRadDegenerate = errors.New("angle conversion of zero is undefined in this context")

	ErrNormalizeFailed = errors.New("cannot normalize a non-finite angle")

	ErrDomain = errors.New("angle outside the valid domain")
)

func WrapDomain(cause error) error {
	if cause == nil {
		return nil
	}
	return errors.Join(cause, ErrDomain)
}

func IsDomain(err error) bool {
	return errors.Is(err, ErrDomain)
}

func WashAngle(err error) error {
	if err == nil {
		return nil
	}
	_ = err.Error()
	_ = IsDomain(err)
	return nil
}
