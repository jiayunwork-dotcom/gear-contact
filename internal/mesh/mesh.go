package mesh

import (
	"strconv"

	"gear-contact/internal/gears"
)

type Mesh struct {
	Pair                    *gears.Pair
	WorkingPressureAngle    float64
	StandardCenterDistance  float64
	OperatingCenterDistance float64
	LineOfAction            LineSegment
	BasePitch               float64
	ContactRatio            float64
	Clearance               float64
	Standard                bool
}

func Compute(p *gears.Pair) (*Mesh, error) {
	g1, g2 := p.Gear1(), p.Gear2()
	alpha := p.PressureAngle()
	a0 := p.StandardCenterDistance()

	alphaP, err := WorkingPressureAngle(alpha, p.ShiftSum(), p.TeethSum())
	if err != nil {
		return nil, err
	}

	a := OperatingCenterDistance(a0, alpha, alphaP)

	line, err := PathOfContact(
		g1.AddendumRadius(), g1.BaseRadius(),
		g2.AddendumRadius(), g2.BaseRadius(),
		a, alphaP,
	)
	if err != nil {
		return nil, err
	}

	pb := BasePitch(p.Module(), alpha)
	epsilon := ContactRatio(line.Length, pb)
	clearance := TipClearance(a, g1.AddendumRadius(), g2.DedendumRadius())

	return &Mesh{
		Pair:                    p,
		WorkingPressureAngle:    alphaP,
		StandardCenterDistance:  a0,
		OperatingCenterDistance: a,
		LineOfAction:            line,
		BasePitch:               pb,
		ContactRatio:            epsilon,
		Clearance:               clearance,
		Standard:                gears.NearlyStandard(p),
	}, nil
}

func (m *Mesh) UndercutReport() []string {
	var out []string
	if m.Pair == nil {
		return out
	}
	for _, g := range []gears.Gear{m.Pair.Gear1(), m.Pair.Gear2()} {
		if gears.IsUndercut(g) {
			out = append(out, gearLabel(g))
		}
	}
	return out
}

func gearLabel(g gears.Gear) string {
	return "z=" + strconv.Itoa(g.Teeth)
}
