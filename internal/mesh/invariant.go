package mesh

import (
	"fmt"
	"math"
)

// Base circle invariant. Every mesh that came out of Compute satisfies
//
//	a * cos(alpha') = a0 * cos(alpha)
//
// because the operating centre distance was built from exactly this relation.
// The verification below re-derives both sides from the mesh and compares
// them, so a regression in the centre distance logic cannot hide behind a
// consistent-looking report.

// BaseCircleInvariant returns the absolute difference between the two sides of
// the base circle invariant, and an error if the difference exceeds the
// tolerance the rest of the geometry works at.
func (m *Mesh) BaseCircleInvariant() (float64, error) {
	alpha := m.Pair.PressureAngle()
	lhs := m.OperatingCenterDistance * math.Cos(m.WorkingPressureAngle)
	rhs := m.StandardCenterDistance * math.Cos(alpha)
	diff := math.Abs(lhs - rhs)
	if diff > 1e-9 {
		return diff, fmt.Errorf("base circle invariant violated: a*cos(alpha') - a0*cos(alpha) = %v", diff)
	}
	return diff, nil
}

// InvariantOK is a predicate form of BaseCircleInvariant for callers that only
// need a boolean gate.
func (m *Mesh) InvariantOK() bool {
	_, err := m.BaseCircleInvariant()
	return err == nil
}

// ContactRatioConsistency returns the contact ratio recomputed from the two
// reported quantities (line length / base pitch) so a caller can confirm the
// report field is not an independent number.
func (m *Mesh) ContactRatioConsistency() float64 {
	return m.LineOfAction.Length / m.BasePitch
}
