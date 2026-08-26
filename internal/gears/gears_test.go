package gears

import (
	"errors"
	"math"
	"strconv"
	"testing"

	"gear-contact/internal/angle"
)

func TestPitchDiameter(t *testing.T) {
	g := Gear{Module: 1, Teeth: 20, PressureAngle: angle.DegToRad(20), AddendumCoeff: 1, ClearanceCoeff: 0.25}
	if got := g.PitchDiameter(); got != 20.0 {
		t.Errorf("pitch diameter: got %v, want 20", got)
	}
	if got := g.PitchRadius(); got != 10.0 {
		t.Errorf("pitch radius: got %v, want 10", got)
	}
	g2 := Gear{Module: 2, Teeth: 40, PressureAngle: angle.DegToRad(20), AddendumCoeff: 1, ClearanceCoeff: 0.25}
	if got := g2.PitchDiameter(); got != 80.0 {
		t.Errorf("pitch diameter for m=2 z=40: got %v, want 80", got)
	}
}

func TestBaseAndAddendumDiameter(t *testing.T) {
	alpha := angle.DegToRad(20.0)
	cos := math.Cos(alpha)
	g := Gear{Module: 1, Teeth: 20, PressureAngle: alpha, AddendumCoeff: 1, ClearanceCoeff: 0.25}
	wantDb := 20.0 * cos
	if got := g.BaseDiameter(); math.Abs(got-wantDb) > 1e-12 {
		t.Errorf("base diameter: got %v, want %v", got, wantDb)
	}
	if got := g.AddendumDiameter(); got != 22.0 {
		t.Errorf("addendum diameter at ha*=1: got %v, want 22", got)
	}

	gCustom := Gear{Module: 1, Teeth: 20, PressureAngle: alpha, AddendumCoeff: 1.2, ClearanceCoeff: 0.25}
	if got := gCustom.AddendumDiameter(); got != 22.4 {
		t.Errorf("addendum diameter at ha*=1.2: got %v, want 22.4", got)
	}

	gShift := Gear{Module: 1, Teeth: 20, PressureAngle: alpha, AddendumCoeff: 1, ClearanceCoeff: 0.25, Shift: 0.1}
	if got := gShift.AddendumDiameter(); got != 22.2 {
		t.Errorf("addendum diameter at x=0.1: got %v, want 22.2", got)
	}
}

func TestModuleZeroRejected(t *testing.T) {
	data := []byte(`{"module": 0, "pressure_angle_deg": 20.0, "gear1": {"teeth": 20}, "gear2": {"teeth": 40}}`)
	_, err := FromJSON(data)
	if err == nil {
		t.Fatal("FromJSON with module 0: got nil error, want ErrModuleNonPositive")
	}
	if !errors.Is(err, ErrModuleNonPositive) {
		t.Errorf("FromJSON with module 0: error %v, want ErrModuleNonPositive", err)
	}
}

func TestTeethNonPositiveRejected(t *testing.T) {
	for _, teeth := range []int{0, -5} {
		data := []byte(`{"module": 1, "pressure_angle_deg": 20.0, "gear1": {"teeth": ` + strconv.Itoa(teeth) + `}, "gear2": {"teeth": 40}}`)
		_, err := FromJSON(data)
		if err == nil {
			t.Errorf("FromJSON with z=%d: got nil error, want ErrTeethNonPositive", teeth)
			continue
		}
		if !errors.Is(err, ErrTeethNonPositive) {
			t.Errorf("FromJSON with z=%d: error %v, want ErrTeethNonPositive", teeth, err)
		}
	}
}

func TestUndercutRejected(t *testing.T) {
	data := []byte(`{"module": 1, "pressure_angle_deg": 20.0, "gear1": {"teeth": 8}, "gear2": {"teeth": 40}}`)
	_, err := FromJSON(data)
	if err == nil {
		t.Fatal("FromJSON with z1=8, x=0: got nil error, want ErrUndercut")
	}
	if !errors.Is(err, ErrUndercut) {
		t.Errorf("FromJSON with z1=8, x=0: error %v, want ErrUndercut", err)
	}

	shifted := []byte(`{"module": 1, "pressure_angle_deg": 20.0, "gear1": {"teeth": 8, "shift": 0.6}, "gear2": {"teeth": 40}}`)
	if _, err := FromJSON(shifted); err != nil {
		t.Errorf("FromJSON with z1=8, x=0.6: got error %v, want nil (enough shift avoids undercut)", err)
	}
}

func TestMinShiftMatchesClassicRule(t *testing.T) {
	alpha := angle.DegToRad(20.0)
	xmin := MinShift(alpha, 1.0, 8)
	classic := (17.0 - 8.0) / 17.0
	if math.Abs(xmin-classic) > 0.005 {
		t.Errorf("MinShift(20 deg, z=8): got %v, want within 0.005 of classic %v", xmin, classic)
	}
	if zmin := MinTeeth(alpha, 1.0); math.Abs(zmin-17.097) > 0.01 {
		t.Errorf("MinTeeth(20 deg): got %v, want ~17.097", zmin)
	}
}

func TestSwapRoundTrip(t *testing.T) {
	p := buildPair(t, `{"module": 1, "pressure_angle_deg": 20.0, "gear1": {"teeth": 20}, "gear2": {"teeth": 40}}`)
	swapped := p.Swap()
	if swapped.Gear1().Teeth != 40 || swapped.Gear2().Teeth != 20 {
		t.Errorf("swap: got z=%d/%d, want 40/20", swapped.Gear1().Teeth, swapped.Gear2().Teeth)
	}
	if !SamePair(p, swapped.Swap()) {
		t.Error("swap round trip: pair differs from the original")
	}
}

func buildPair(t *testing.T, spec string) *Pair {
	t.Helper()
	p, err := FromJSON([]byte(spec))
	if err != nil {
		t.Fatalf("build pair: %v", err)
	}
	return p
}
