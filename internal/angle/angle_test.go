package angle

import (
	"errors"
	"math"
	"testing"
)

func TestDegToRadRoundTrip(t *testing.T) {
	deg := 20.0
	rad := DegToRad(deg)
	back := RadToDeg(rad)
	if math.Abs(back-deg) > 1e-12 {
		t.Errorf("DegToRad/RadToDeg round trip: got %v deg, want %v deg", back, deg)
	}
	if DegToRad(90.0) != math.Pi/2 {
		t.Errorf("DegToRad(90): got %v, want %v", DegToRad(90.0), math.Pi/2)
	}
	if DegToRad(180.0) != math.Pi {
		t.Errorf("DegToRad(180): got %v, want %v", DegToRad(180.0), math.Pi)
	}
}

func TestInvoluteFormula(t *testing.T) {
	cases := []struct {
		deg  float64
		want float64
	}{
		{20.0, 0.0149044},
		{25.0, 0.0299753},
	}
	for _, c := range cases {
		got := Inv(DegToRad(c.deg))
		if math.Abs(got-c.want) > 1e-5 {
			t.Errorf("Inv(%v deg): got %v, want %v", c.deg, got, c.want)
		}
	}
}

func TestInvSolveRecoversWorkingAngle(t *testing.T) {
	target := Inv(DegToRad(25.0))
	got, err := InvSolve(target, DegToRad(20.0))
	if err != nil {
		t.Fatalf("InvSolve returned an error for a reachable target: %v", err)
	}
	if math.Abs(RadToDeg(got)-25.0) > 1e-9 {
		t.Errorf("InvSolve target inv(25 deg): got %v deg, want 25 deg", RadToDeg(got))
	}
	residual := Inv(got) - target
	if math.Abs(residual) > 1e-12 {
		t.Errorf("InvSolve residual: got %v, want within 1e-12", residual)
	}
}

func TestPressureAnglePolicy(t *testing.T) {
	cases := []struct {
		deg   float64
		sent  error
		valid bool
	}{
		{0.0, ErrPressureAngleZero, false},
		{90.0, ErrPressureAngleTooLarge, false},
		{10.0, ErrPressureAngleOutOfRange, false},
		{14.5, nil, true},
		{20.0, nil, true},
		{25.0, nil, true},
		{25.5, ErrPressureAngleOutOfRange, false},
	}
	for _, c := range cases {
		err := ValidatePressureAngle(DegToRad(c.deg))
		if c.valid {
			if err != nil {
				t.Errorf("ValidatePressureAngle(%v deg): got error %v, want nil", c.deg, err)
			}
			continue
		}
		if err == nil {
			t.Errorf("ValidatePressureAngle(%v deg): got nil, want %v", c.deg, c.sent)
			continue
		}
		if !errors.Is(err, c.sent) {
			t.Errorf("ValidatePressureAngle(%v deg): error %v, want %v", c.deg, err, c.sent)
		}
	}
}
