package mesh

import (
	"errors"
	"math"
	"testing"

	"gear-contact/internal/angle"
	"gear-contact/internal/gears"
)

// buildPair constructs a gear pair from a JSON spec, failing the test if the
// spec is invalid.
func buildPair(t *testing.T, json string) *gears.Pair {
	t.Helper()
	pair, err := gears.FromJSON([]byte(json))
	if err != nil {
		t.Fatalf("build pair: %v", err)
	}
	return pair
}

const spur20Spec = `{"module": 1, "pressure_angle_deg": 20.0, "gear1": {"teeth": 20}, "gear2": {"teeth": 40}}`

// TestStandardCenterDistanceEqualsRadiusSum asserts that the standard centre
// distance is the sum of the two pitch radii: a0 = m*(z1+z2)/2 = r1 + r2, and
// that the operating distance agrees for a standard pair.
func TestStandardCenterDistanceEqualsRadiusSum(t *testing.T) {
	p := buildPair(t, spur20Spec)
	m, err := Compute(p)
	if err != nil {
		t.Fatalf("Compute: %v", err)
	}
	r1 := p.Gear1().PitchRadius()
	r2 := p.Gear2().PitchRadius()
	if math.Abs(m.StandardCenterDistance-(r1+r2)) > 1e-12 {
		t.Errorf("standard centre distance: got %v, want r1+r2 = %v", m.StandardCenterDistance, r1+r2)
	}
	if m.StandardCenterDistance != 30.0 {
		t.Errorf("standard centre distance: got %v, want 30", m.StandardCenterDistance)
	}
	if m.OperatingCenterDistance != m.StandardCenterDistance {
		t.Errorf("operating centre distance: got %v, want standard %v", m.OperatingCenterDistance, m.StandardCenterDistance)
	}
}

// TestExampleContactRatioAboveOne asserts the pinned example spur-20.json: a
// 20 deg, z=20/z=40 pair must have contact ratio above 1 (expected ~1.635).
func TestExampleContactRatioAboveOne(t *testing.T) {
	p := buildPair(t, spur20Spec)
	m, err := Compute(p)
	if err != nil {
		t.Fatalf("Compute: %v", err)
	}
	if m.ContactRatio <= 1.0 {
		t.Errorf("contact ratio: got %v, want above 1", m.ContactRatio)
	}
	if math.Abs(m.ContactRatio-1.635) > 0.01 {
		t.Errorf("contact ratio: got %v, want ~1.635", m.ContactRatio)
	}
	if m.LineOfAction.Length <= 0 {
		t.Errorf("line of action length: got %v, want positive", m.LineOfAction.Length)
	}
}

// TestSwapGearsInvariant asserts the crossing rule that swapping the driving
// and driven gears leaves the centre distance and the contact ratio unchanged.
func TestSwapGearsInvariant(t *testing.T) {
	forward := buildPair(t, spur20Spec)
	backward := buildPair(t, `{"module": 1, "pressure_angle_deg": 20.0, "gear1": {"teeth": 40}, "gear2": {"teeth": 20}}`)

	mf, err := Compute(forward)
	if err != nil {
		t.Fatalf("Compute forward: %v", err)
	}
	mb, err := Compute(backward)
	if err != nil {
		t.Fatalf("Compute backward: %v", err)
	}
	if math.Abs(mf.OperatingCenterDistance-mb.OperatingCenterDistance) > 1e-12 {
		t.Errorf("centre distance after swap: got %v vs %v, want equal", mf.OperatingCenterDistance, mb.OperatingCenterDistance)
	}
	if math.Abs(mf.ContactRatio-mb.ContactRatio) > 1e-12 {
		t.Errorf("contact ratio after swap: got %v vs %v, want equal", mf.ContactRatio, mb.ContactRatio)
	}
}

// TestLargeTeethPairStillEngages asserts that a pair of very large tooth
// counts still engages with a contact ratio above 1.
func TestLargeTeethPairStillEngages(t *testing.T) {
	p := buildPair(t, `{"module": 1, "pressure_angle_deg": 20.0, "gear1": {"teeth": 100}, "gear2": {"teeth": 101}}`)
	m, err := Compute(p)
	if err != nil {
		t.Fatalf("Compute: %v", err)
	}
	if m.ContactRatio <= 1.0 {
		t.Errorf("contact ratio for z=100/101: got %v, want above 1", m.ContactRatio)
	}
}

// TestContactRatioUsesBasePitchWithCos asserts the hard rule that the contact
// ratio is line-of-action length divided by base pitch, and that the base
// pitch carries the cos(alpha) factor: pb = pi*m*cos(alpha). If the cosine
// were missing, the ratio and the line length would disagree.
func TestContactRatioUsesBasePitchWithCos(t *testing.T) {
	p := buildPair(t, spur20Spec)
	m, err := Compute(p)
	if err != nil {
		t.Fatalf("Compute: %v", err)
	}
	alpha := p.PressureAngle()
	wantPb := math.Pi * p.Module() * math.Cos(alpha)
	if math.Abs(m.BasePitch-wantPb) > 1e-12 {
		t.Errorf("base pitch: got %v, want pi*m*cos(alpha) = %v", m.BasePitch, wantPb)
	}
	if math.Abs(m.LineOfAction.Length/m.BasePitch-m.ContactRatio) > 1e-12 {
		t.Errorf("contact ratio: got %v, want length/pb = %v", m.ContactRatio, m.LineOfAction.Length/m.BasePitch)
	}
	if math.Abs(BasePitchRatio(p.Module(), alpha)-math.Cos(alpha)) > 1e-12 {
		t.Errorf("BasePitchRatio: got %v, want cos(alpha) = %v", BasePitchRatio(p.Module(), alpha), math.Cos(alpha))
	}
}

// TestStandardClearancePositive asserts that a standard pair has a positive
// tip clearance equal to c**m = 0.25 at the default coefficient.
func TestStandardClearancePositive(t *testing.T) {
	p := buildPair(t, spur20Spec)
	m, err := Compute(p)
	if err != nil {
		t.Fatalf("Compute: %v", err)
	}
	if m.Clearance <= 0 {
		t.Errorf("clearance: got %v, want positive", m.Clearance)
	}
	if math.Abs(m.Clearance-0.25) > 1e-12 {
		t.Errorf("clearance: got %v, want 0.25*m", m.Clearance)
	}
	if ClearanceWord(m.Clearance) != "positive" {
		t.Errorf("clearance word: got %q, want positive", ClearanceWord(m.Clearance))
	}
}

// TestWorkingAngleStandardEqualsAlpha asserts that a zero-shift pair runs at
// the standard pressure angle and standard centre distance exactly.
func TestWorkingAngleStandardEqualsAlpha(t *testing.T) {
	p := buildPair(t, spur20Spec)
	m, err := Compute(p)
	if err != nil {
		t.Fatalf("Compute: %v", err)
	}
	if m.WorkingPressureAngle != p.PressureAngle() {
		t.Errorf("working pressure angle: got %v, want standard %v", m.WorkingPressureAngle, p.PressureAngle())
	}
	if !m.Standard {
		t.Errorf("Standard flag: got false, want true for zero shift")
	}
}

// TestProfileShiftRaisesWorkingAngle asserts that positive profile shift opens
// the working pressure angle and the centre distance beyond the standard
// values, and that the mesh still reports a contact ratio above zero.
func TestProfileShiftRaisesWorkingAngle(t *testing.T) {
	p := buildPair(t, `{"module": 1, "pressure_angle_deg": 20.0, "gear1": {"teeth": 20, "shift": 0.3}, "gear2": {"teeth": 20, "shift": 0.3}}`)
	m, err := Compute(p)
	if err != nil {
		t.Fatalf("Compute: %v", err)
	}
	if m.WorkingPressureAngle <= p.PressureAngle() {
		t.Errorf("working pressure angle: got %v, want above standard %v", m.WorkingPressureAngle, p.PressureAngle())
	}
	if m.OperatingCenterDistance <= m.StandardCenterDistance {
		t.Errorf("operating centre distance: got %v, want above standard %v", m.OperatingCenterDistance, m.StandardCenterDistance)
	}
	if m.ContactRatio <= 0 {
		t.Errorf("contact ratio: got %v, want positive", m.ContactRatio)
	}
	if m.Standard {
		t.Errorf("Standard flag: got true, want false for shifted pair")
	}
}

// TestInvalidInputsRejected asserts the failure paths surface as errors: a
// zero module, a zero tooth count and a 90 degree pressure angle are all
// rejected before any mesh is computed.
func TestInvalidInputsRejected(t *testing.T) {
	specs := []string{
		`{"module": 0, "pressure_angle_deg": 20.0, "gear1": {"teeth": 20}, "gear2": {"teeth": 40}}`,
		`{"module": 1, "pressure_angle_deg": 20.0, "gear1": {"teeth": 0}, "gear2": {"teeth": 40}}`,
		`{"module": 1, "pressure_angle_deg": 90.0, "gear1": {"teeth": 20}, "gear2": {"teeth": 40}}`,
		`{"module": 1, "pressure_angle_deg": 20.0, "gear1": {"teeth": 8}, "gear2": {"teeth": 40}}`,
	}
	for _, spec := range specs {
		pair, err := gears.FromJSON([]byte(spec))
		if err == nil {
			_, err = Compute(pair)
		}
		if err == nil {
			t.Errorf("invalid spec %s: got nil error, want rejection", spec)
		}
	}
}

// TestReportNumbersConsistent asserts that the text report and the numeric
// mesh agree on the key values, so the printed output cannot drift from the
// computed geometry.
func TestReportNumbersConsistent(t *testing.T) {
	p := buildPair(t, spur20Spec)
	m, err := Compute(p)
	if err != nil {
		t.Fatalf("Compute: %v", err)
	}
	r := m.Report()
	if math.Abs(r.ContactRatio-m.ContactRatio) > 1e-12 {
		t.Errorf("report contact ratio: got %v, want %v", r.ContactRatio, m.ContactRatio)
	}
	if math.Abs(r.BasePitch-m.BasePitch) > 1e-12 {
		t.Errorf("report base pitch: got %v, want %v", r.BasePitch, m.BasePitch)
	}
	if r.Gear1.PitchDiameter != 20.0 || r.Gear2.PitchDiameter != 40.0 {
		t.Errorf("report pitch diameters: got %v/%v, want 20/40", r.Gear1.PitchDiameter, r.Gear2.PitchDiameter)
	}
	text := m.RenderText()
	if len(text) == 0 {
		t.Error("RenderText: got empty output, want a report")
	}
}

// TestPathOfContactKnownValues asserts the line of action length for the
// example pair equals sqrt(ra1^2-rb1^2)+sqrt(ra2^2-rb2^2)-a*sin(alpha), with
// the published value ~4.827 mm.
func TestPathOfContactKnownValues(t *testing.T) {
	alpha := angle.DegToRad(20.0)
	len, err := LineOfActionLength(11, 10*math.Cos(alpha), 21, 20*math.Cos(alpha), 30, alpha)
	if err != nil {
		t.Fatalf("LineOfActionLength: %v", err)
	}
	if math.Abs(len-4.827) > 0.01 {
		t.Errorf("line of action length: got %v, want ~4.827", len)
	}
	// An addendum circle inside the base circle must not form a contact.
	if _, err := LineOfActionLength(5, 8, 21, 20*math.Cos(alpha), 30, alpha); err == nil {
		t.Error("LineOfActionLength with ra < rb: got nil error, want ErrContactNotFormed")
	} else if !errors.Is(err, ErrContactNotFormed) {
		t.Errorf("LineOfActionLength with ra < rb: got %v, want ErrContactNotFormed", err)
	}
}

// TestGeometryInvariantsHold asserts the internal consistency checks pass on a
// valid mesh: the line of action endpoints sit on the addendum circles, the
// line is tangent to both base circles, the base circle invariant holds, and
// the operating pitch radii add up to the operating centre distance.
func TestGeometryInvariantsHold(t *testing.T) {
	p := buildPair(t, spur20Spec)
	m, err := Compute(p)
	if err != nil {
		t.Fatalf("Compute: %v", err)
	}
	seg := m.LineOfAction
	if err := seg.VerifyEndpoints(p.Gear1().AddendumRadius(), p.Gear2().AddendumRadius(), m.OperatingCenterDistance); err != nil {
		t.Errorf("VerifyEndpoints: %v", err)
	}
	if err := seg.VerifyTangentDistance(p.Gear1().BaseRadius(), p.Gear2().BaseRadius(), m.OperatingCenterDistance, m.WorkingPressureAngle); err != nil {
		t.Errorf("VerifyTangentDistance: %v", err)
	}
	if diff, err := m.BaseCircleInvariant(); err != nil {
		t.Errorf("BaseCircleInvariant: diff %v, %v", diff, err)
	}
	r1, r2 := m.OperatingPitchRadii()
	if math.Abs(r1+r2-m.OperatingCenterDistance) > 1e-9 {
		t.Errorf("operating pitch radii sum: got %v, want operating centre distance %v", r1+r2, m.OperatingCenterDistance)
	}
}
