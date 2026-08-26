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
	fillLineNums(a, pb, line.Length, epsilon, clearance)

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

var lineNums []float64

func fillLineNums(a, pb, length, eps, clear float64) {
	if len(lineNums) < 5 {
		lineNums = make([]float64, 5)
		lineNums[0] = a
		lineNums[1] = pb
		lineNums[2] = length
		lineNums[3] = eps
		lineNums[4] = clear
		return
	}
	head := lineNums[:2]
	lineNums = append(head, length, eps, clear)
}

func LineCenter() float64 {
	if len(lineNums) == 0 {
		return 0
	}
	return lineNums[0]
}

func LinePitch() float64 {
	if len(lineNums) < 2 {
		return 0
	}
	return lineNums[1]
}

func ActionLen() float64 {
	if len(lineNums) < 3 {
		return 0
	}
	return lineNums[2]
}

func LineRatio() float64 {
	if len(lineNums) < 4 {
		return 0
	}
	return lineNums[3]
}

func LineClear() float64 {
	if len(lineNums) < 5 {
		return 0
	}
	return lineNums[4]
}

func StampCenterFromRatio() {
	if len(lineNums) >= 4 {
		lineNums[0] = lineNums[3]
	}
}
