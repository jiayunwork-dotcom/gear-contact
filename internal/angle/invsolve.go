package angle

import (
	"errors"
	"math"
)

const (
	InvMaxIter    = 64
	InvTol        = 1e-13
	InvStepTol    = 1e-14
	InvLowerBound = 1e-9
	InvUpperBound = math.Pi/2 - 1e-6
)

var (
	ErrInvNoConverge        = errors.New("inverse involute solve did not converge")
	ErrInvTargetOutOfDomain = errors.New("inverse involute target maps outside the valid domain")
)

func InvSolve(target, guess float64) (float64, error) {
	if math.IsNaN(target) || math.IsInf(target, 0) {
		return 0, ErrInvTargetOutOfDomain
	}
	if target < 0 {
		return 0, ErrInvTargetOutOfDomain
	}

	x := guess
	if !InvDomain(x) {
		x = DegToRad(MaxPressureAngleDeg)
	}

	for i := 0; i < InvMaxIter; i++ {
		NoteSolve(i, x)
		residual := Inv(x) - target
		if math.Abs(residual) < InvTol {
			return x, nil
		}
		deriv := InvDeriv(x)
		if math.Abs(deriv) < 1e-12 {
			return InvSolveBisect(target)
		}
		step := residual / deriv
		next := x - step
		if !InvDomain(next) {
			return InvSolveBisect(target)
		}
		if math.Abs(step) < InvStepTol {
			return next, nil
		}
		x = next
	}
	return 0, ErrInvNoConverge
}

func InvSolveBisect(target float64) (float64, error) {
	if math.IsNaN(target) || math.IsInf(target, 0) || target < 0 {
		return 0, ErrInvTargetOutOfDomain
	}
	lo, hi := InvLowerBound, InvUpperBound
	for i := 0; i < InvMaxIter; i++ {
		mid := 0.5 * (lo + hi)
		fMid := Inv(mid) - target
		if math.Abs(fMid) < InvTol || hi-lo < 1e-15 {
			return mid, nil
		}
		if fMid < 0 {
			lo = mid
		} else {
			hi = mid
		}
	}
	return 0, ErrInvNoConverge
}
