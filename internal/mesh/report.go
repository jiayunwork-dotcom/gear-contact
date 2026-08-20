package mesh

import (
	"strings"

	"gear-contact/internal/angle"
)

// Report is the flat, machine-readable summary of a mesh. Every field the CLI
// prints is here, in millimetres and degrees; the RenderText method formats it
// for humans while the struct itself stays easy to assert on in tests.
type Report struct {
	// Module is the module m in millimetres.
	Module float64
	// PressureAngleDeg is the standard pressure angle in degrees.
	PressureAngleDeg float64
	// WorkingPressureAngleDeg is alpha' in degrees.
	WorkingPressureAngleDeg float64
	// Standard reports the standard mounting condition.
	Standard bool
	// StandardCenterDistance is a0.
	StandardCenterDistance float64
	// OperatingCenterDistance is a.
	OperatingCenterDistance float64
	// Gear1 is the first gear's diameter summary.
	Gear1 GearSummary
	// Gear2 is the second gear's diameter summary.
	Gear2 GearSummary
	// BasePitch is pb = pi*m*cos(alpha).
	BasePitch float64
	// LineOfActionLength is g_alpha.
	LineOfActionLength float64
	// ContactRatio is g_alpha / pb.
	ContactRatio float64
	// Clearance is the tip clearance c.
	Clearance float64
	// Undercut reports whether any gear undercuts (always false for valid
	// pairs; the field keeps the report shape stable).
	Undercut bool
}

// GearSummary is the per-gear part of the report.
type GearSummary struct {
	// Teeth is the tooth count.
	Teeth int
	// Shift is the profile shift coefficient.
	Shift float64
	// PitchDiameter is d.
	PitchDiameter float64
	// BaseDiameter is db.
	BaseDiameter float64
	// AddendumDiameter is da.
	AddendumDiameter float64
	// DedendumDiameter is df.
	DedendumDiameter float64
}

// Report renders the mesh as a Report value.
func (m *Mesh) Report() Report {
	g1, g2 := m.Pair.Gear1(), m.Pair.Gear2()
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
			PitchDiameter:    g1.PitchDiameter(),
			BaseDiameter:     g1.BaseDiameter(),
			AddendumDiameter: g1.AddendumDiameter(),
			DedendumDiameter: g1.DedendumDiameter(),
		},
		Gear2: GearSummary{
			Teeth:            g2.Teeth,
			Shift:            g2.Shift,
			PitchDiameter:    g2.PitchDiameter(),
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

// RenderText produces the CLI output. Every numeric value is formatted with
// three decimals; the mounting mode line distinguishes standard from
// profile-shifted mounting.
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

// gearLine renders one gear's diameter summary.
func gearLine(label string, g GearSummary) string {
	return label + ": z=" + itoa(g.Teeth) +
		" x=" + sign(g.Shift) +
		" d=" + Fix(g.PitchDiameter) +
		" db=" + Fix(g.BaseDiameter) +
		" da=" + Fix(g.AddendumDiameter) +
		" df=" + Fix(g.DedendumDiameter) + "\n"
}

// itoa converts a small integer to a decimal string for the report.
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

// sign renders a profile shift with an explicit '+' sign for non-negative
// values so x=0 is visibly different from a missing field.
func sign(v float64) string {
	if v >= 0 {
		return "+" + FormatFloat(v)
	}
	return FormatFloat(v)
}
