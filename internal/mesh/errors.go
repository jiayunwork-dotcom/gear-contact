package mesh

import "errors"

var (
	ErrNoMesh = errors.New("mesh is nil")
)

func IsContactError(err error) bool {
	return errors.Is(err, ErrContactNotFormed) || errors.Is(err, ErrPathNotPositive)
}

func IsWorkingError(err error) bool {
	return errors.Is(err, ErrWorkingAngle)
}
