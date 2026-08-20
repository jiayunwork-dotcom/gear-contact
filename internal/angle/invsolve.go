package angle

import (
	"errors"
	"math"
)

// Solver tuning. These constants bound the numerical inverse involute solve;
// they are deliberately conservative so the iteration either converges to
// machine precision or fails with a typed error instead of looping forever.
const (
	// InvMaxIter is the iteration cap for both the Newton and bisection paths.
	InvMaxIter = 64
	// InvTol is the absolute tolerance on the function residual.
	InvTol = 1e-13
	// InvStepTol is the tolerance on the step size between two iterates.
	InvStepTol = 1e-14
	// InvLowerBound is the smallest angle (in radians) the solver will return.
	InvLowerBound = 1e-9
	// InvUpperBound is the largest angle the solver will return: pi/2 minus a
	// margin so tan(x) stays finite.
	InvUpperBound = math.Pi/2 - 1e-6
)

var (
	// ErrInvNoConverge is returned when the iteration count is exhausted
	// without meeting either tolerance.
	ErrInvNoConverge = errors.New("inverse involute solve did not converge")
	// ErrInvTargetOutOfDomain is returned when the target value maps to an
	// angle outside the valid domain.
	ErrInvTargetOutOfDomain = errors.New("inverse involute target maps outside the valid domain")
)

// InvSolve finds x in (0, pi/2) with Inv(x) == target by Newton iteration,
// starting from guess. The guess should be near the answer; the mesh package
// passes the standard pressure angle, which is already a good first estimate
// for small profile shifts.
//
// The Newton step is x_{k+1} = x_k - (Inv(x_k) - target) / InvDeriv(x_k).
// Near x = 0 the derivative tan^2(x) vanishes, so a fallback bisection over
// the whole domain is used whenever the derivative drops below a threshold or
// the iterate leaves the domain. The returned error is either
// ErrInvTargetOutOfDomain (when the target is not reachable) or
// ErrInvNoConverge (iteration cap hit).
func InvSolve(target, guess float64) (float64, error) {
	if math.IsNaN(target) || math.IsInf(target, 0) {
		return 0, ErrInvTargetOutOfDomain
	}
	if target < 0 {
		// Inv(x) < 0 for every x in (0, pi/2); such a target is unreachable.
		return 0, ErrInvTargetOutOfDomain
	}

	x := guess
	if !InvDomain(x) {
		x = DegToRad(MaxPressureAngleDeg)
	}

	for i := 0; i < InvMaxIter; i++ {
		residual := Inv(x) - target
		if math.Abs(residual) < InvTol {
			return x, nil
		}
		deriv := InvDeriv(x)
		if math.Abs(deriv) < 1e-12 {
			// Newton would divide by a numerically zero derivative; switch to
			// the bracketed fallback over the whole domain.
			return InvSolveBisect(target)
		}
		step := residual / deriv
		next := x - step
		if !InvDomain(next) {
			// The step left the domain; bisection on the whole interval is the
			// safe fallback.
			return InvSolveBisect(target)
		}
		if math.Abs(step) < InvStepTol {
			return next, nil
		}
		x = next
	}
	return 0, ErrInvNoConverge
}

// InvSolveBisect finds x with Inv(x) == target inside the whole valid domain
// by plain bisection. Inv is strictly increasing on (0, pi/2), so for any
// reachable non-negative target the interval [InvLowerBound, InvUpperBound]
// brackets exactly one root and bisection is guaranteed to converge. The
// upper bound is capped below pi/2 so tan(x) remains finite.
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
