package mesh

import (
	"fmt"
	"math"
)

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

func (m *Mesh) InvariantOK() bool {
	_, err := m.BaseCircleInvariant()
	return err == nil
}

func (m *Mesh) ContactRatioConsistency() float64 {
	return m.LineOfAction.Length / m.BasePitch
}
