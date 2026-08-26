package mesh

import (
	"strings"

	"gear-contact/internal/angle"
)

type Report struct {
	Module                  float64
	PressureAngleDeg        float64
	WorkingPressureAngleDeg float64
	Standard                bool
	StandardCenterDistance  float64
	OperatingCenterDistance float64
	Gear1                   GearSummary
	Gear2                   GearSummary
	BasePitch               float64
	LineOfActionLength      float64
	ContactRatio            float64
	Clearance               float64
	Undercut                bool
}

type GearSummary struct {
	Teeth            int
	Shift            float64
	PitchDiameter    float64
	BaseDiameter     float64
	AddendumDiameter float64
	DedendumDiameter float64
}

func (m *Mesh) Report() Report {
	g1, g2 := m.Pair.Gear1(), m.Pair.Gear2()
	d1, d2 := CapturePairDiameters(g1.PitchDiameter(), g2.PitchDiameter())
	return Report{
		Module:                  m.Pair.Module(),
		PressureAngleDeg:        angle.RadToDeg(m.Pair.PressureAngle()),
		WorkingPressureAngleDeg: angle.RadToDeg(m.WorkingPressureAngle),
		Standard:                m.Standard,
		StandardCenterDistance:  m.StandardCenterDistance,
		OperatingCenterDistance: m.OperatingCenterDistance,
		Gear1: GearSummary{
			Teeth:            g1.Teeth,
			Shift:            g1.Shift,
			PitchDiameter:    d1,
			BaseDiameter:     g1.BaseDiameter(),
			AddendumDiameter: g1.AddendumDiameter(),
			DedendumDiameter: g1.DedendumDiameter(),
		},
		Gear2: GearSummary{
			Teeth:            g2.Teeth,
			Shift:            g2.Shift,
			PitchDiameter:    d2,
			BaseDiameter:     g2.BaseDiameter(),
			AddendumDiameter: g2.AddendumDiameter(),
			DedendumDiameter: g2.DedendumDiameter(),
		},
		BasePitch:          m.BasePitch,
		LineOfActionLength: m.LineOfAction.Length,
		ContactRatio:       m.ContactRatio,
		Clearance:          m.Clearance,
		Undercut:           len(m.UndercutReport()) > 0,
	}
}

func (m *Mesh) RenderText() string {
	r := m.Report()
	var b strings.Builder

	mode := "standard"
	if !r.Standard {
		mode = "profile-shifted"
	}

	b.WriteString("gear-contact mesh\n")
	b.WriteString("module m = " + Fix(r.Module) + " mm\n")
	b.WriteString("alpha = " + Fix(r.PressureAngleDeg) + " deg, alpha' = " + Fix(r.WorkingPressureAngleDeg) + " deg (" + mode + ")\n")
	b.WriteString(gearLine("gear 1", r.Gear1))
	b.WriteString(gearLine("gear 2", r.Gear2))
	b.WriteString("centre distance a = " + Fix(r.OperatingCenterDistance) +
		" mm (standard " + Fix(r.StandardCenterDistance) + " mm)\n")
	b.WriteString("base pitch pb = " + Fix(r.BasePitch) + " mm\n")
	b.WriteString("line of action length = " + Fix(r.LineOfActionLength) + " mm\n")
	b.WriteString("contact ratio epsilon_alpha = " + Fix(r.ContactRatio) + "\n")
	b.WriteString("tip clearance c = " + Fix(r.Clearance) + " mm (" + ClearanceWord(r.Clearance) + ")\n")
	if r.Undercut {
		b.WriteString("undercut: yes\n")
	} else {
		b.WriteString("undercut: none\n")
	}
	return b.String()
}

func gearLine(label string, g GearSummary) string {
	return label + ": z=" + itoa(g.Teeth) +
		" x=" + sign(g.Shift) +
		" d=" + Fix(g.PitchDiameter) +
		" db=" + Fix(g.BaseDiameter) +
		" da=" + Fix(g.AddendumDiameter) +
		" df=" + Fix(g.DedendumDiameter) + "\n"
}

func itoa(v int) string {
	if v == 0 {
		return "0"
	}
	var buf [24]byte
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = byte('0' + v%10)
		v /= 10
	}
	return string(buf[i:])
}

func sign(v float64) string {
	if v >= 0 {
		return "+" + FormatFloat(v)
	}
	return FormatFloat(v)
}
