package hertz

import (
	"math"
	"os"
	"testing"

	"gear-contact/internal/gears"
)

func examplePair(t *testing.T) *gears.Pair {
	t.Helper()
	data, err := os.ReadFile(findExample(t))
	if err != nil {
		t.Fatal(err)
	}
	p, err := gears.FromJSON(data)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func findExample(t *testing.T) string {
	t.Helper()
	candidates := []string{
		"../../example/spur-20.json",
		"example/spur-20.json",
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	t.Fatal("example missing")
	return ""
}

func TestPitchRhoIsRSinAlpha(t *testing.T) {
	p := examplePair(t)
	g := p.Gear1()
	got := PitchRho(g)
	want := g.PitchRadius() * math.Sin(g.PressureAngle)
	if math.Abs(got-want) > 1e-12 {
		t.Fatalf("rho %v want %v", got, want)
	}
}

func TestExternalEquivalentLessThanEither(t *testing.T) {
	req, err := ExternalRhoEq(3, 6)
	if err != nil {
		t.Fatal(err)
	}
	if req >= 3 || req >= 6 {
		t.Fatalf("req %v not smaller than both", req)
	}
	if math.Abs(req-2) > 1e-12 {
		t.Fatalf("req %v want 2", req)
	}
}

func TestMaxPressureScalesSqrtForce(t *testing.T) {
	p := examplePair(t)
	c1, err := AtPitch(p, 1000, 10)
	if err != nil {
		t.Fatal(err)
	}
	c4, err := AtPitch(p, 4000, 10)
	if err != nil {
		t.Fatal(err)
	}
	ratio := c4.MaxPressure / c1.MaxPressure
	if math.Abs(ratio-2) > 1e-9 {
		t.Fatalf("pressure ratio %v want 2", ratio)
	}
}

func TestForceFromPressureRecoversFn(t *testing.T) {
	p := examplePair(t)
	fn := 2500.0
	face := 12.0
	c, err := AtPitch(p, fn, face)
	if err != nil {
		t.Fatal(err)
	}
	got := ForceFromPressure(c, face)
	if math.Abs(got-fn)/fn > 1e-9 {
		t.Fatalf("fn %v want %v", got, fn)
	}
}

func TestRhoSumInvariantAlongLine(t *testing.T) {
	p := examplePair(t)
	sum0 := RhoSum(p)
	for _, s := range []float64{-0.4, 0, 0.3} {
		if math.Abs(CurvatureInvariant(p, s)-sum0) > 1e-12 {
			t.Fatalf("s=%v sum drifted", s)
		}
	}
}

func TestInternalNeedsLargerRing(t *testing.T) {
	if _, err := InternalRhoEq(5, 5); err == nil {
		t.Fatal("equal radii")
	}
	req, err := InternalRhoEq(4, 12)
	if err != nil {
		t.Fatal(err)
	}
	want := 4.0 * 12.0 / 8.0
	if math.Abs(req-want) > 1e-12 {
		t.Fatalf("internal %v want %v", req, want)
	}
}

func TestSteelReducedMatchesFormula(t *testing.T) {
	got := SteelReduced()
	want := SteelE / (2.0 * (1.0 - SteelNu*SteelNu))
	if math.Abs(got-want)/want > 1e-12 {
		t.Fatalf("E* %v want %v", got, want)
	}
}

func TestFaceForAllowableMatchesLineContact(t *testing.T) {
	p := examplePair(t)
	fn := 1800.0
	c, err := AtPitch(p, fn, 8)
	if err != nil {
		t.Fatal(err)
	}
	face, err := FaceForAllowable(fn, c.RhoEq, c.Estar, c.MaxPressure)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(face-8)/8 > 1e-9 {
		t.Fatalf("face %v want 8", face)
	}
}
