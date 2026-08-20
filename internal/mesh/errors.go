package mesh

import "errors"

// This file centralizes the mesh package sentinel errors. It also hosts a few
// predicates that turn an error into a decision without string matching.

var (
	// ErrNoMesh is returned when a nil Mesh is passed to a method that needs
	// one; it is a defensive error that should never be seen from the CLI.
	ErrNoMesh = errors.New("mesh is nil")
)

// IsContactError reports whether err is, or wraps, any of the line-of-action
// errors (contact not formed, empty path).
func IsContactError(err error) bool {
	return errors.Is(err, ErrContactNotFormed) || errors.Is(err, ErrPathNotPositive)
}

// IsWorkingError reports whether err is, or wraps, the working pressure angle
// error.
func IsWorkingError(err error) bool {
	return errors.Is(err, ErrWorkingAngle)
}
