package mesh

import (
	"errors"
	"math"
	"testing"

	"gear-contact/internal/angle"
	"gear-contact/internal/gears"
)

func buildPair(t *testing.T, json string) *gears.Pair {
	t.Helper()
	pair, err := gears.FromJSON([]byte(json))
	if err != nil {
		t.Fatalf("build pair: %v", err)
	}
	return pair
}

const spur20Spec = `{"module": 1, "pressure_angle_deg": 20.0, "gear1": {"teeth": 20}, "gear2": {"teeth": 40}}`

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

func TestPathOfContactKnownValues(t *testing.T) {
	alpha := angle.DegToRad(20.0)
	len, err := LineOfActionLength(11, 10*math.Cos(alpha), 21, 20*math.Cos(alpha), 30, alpha)
	if err != nil {
		t.Fatalf("LineOfActionLength: %v", err)
	}
	if math.Abs(len-4.827) > 0.01 {
		t.Errorf("line of action length: got %v, want ~4.827", len)
	}
	if _, err := LineOfActionLength(5, 8, 21, 20*math.Cos(alpha), 30, alpha); err == nil {
		t.Error("LineOfActionLength with ra < rb: got nil error, want ErrContactNotFormed")
	} else if !errors.Is(err, ErrContactNotFormed) {
		t.Errorf("LineOfActionLength with ra < rb: got %v, want ErrContactNotFormed", err)
	}
}

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
