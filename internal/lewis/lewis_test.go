package lewis

import (
	"math"
	"os"
	"testing"

	"gear-contact/internal/gears"
)

func examplePair(t *testing.T) *gears.Pair {
	t.Helper()
	for _, c := range []string{"../../example/spur-20.json", "example/spur-20.json"} {
		data, err := os.ReadFile(c)
		if err != nil {
			continue
		}
		p, err := gears.FromJSON(data)
		if err != nil {
			t.Fatal(err)
		}
		return p
	}
	t.Fatal("example missing")
	return nil
}

func TestFormFactorDecreasesWithFewerTeeth(t *testing.T) {
	y20, err := FormY(20)
	if err != nil {
		t.Fatal(err)
	}
	y40, err := FormY(40)
	if err != nil {
		t.Fatal(err)
	}
	if y20 >= y40 {
		t.Fatalf("y(20)=%v y(40)=%v", y20, y40)
	}
}

func TestFormYKnownValue(t *testing.T) {
	y, err := FormY(20)
	if err != nil {
		t.Fatal(err)
	}
	want := 0.154 - 0.912/20.0
	if math.Abs(y-want) > 1e-12 {
		t.Fatalf("y %v want %v", y, want)
	}
}

func TestBendingInverselyFaceWidth(t *testing.T) {
	s1, err := Bending(1000, 10, 2, 0.3)
	if err != nil {
		t.Fatal(err)
	}
	s2, err := Bending(1000, 20, 2, 0.3)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(s1/s2-2) > 1e-12 {
		t.Fatalf("ratio %v", s1/s2)
	}
}

func TestBarthAtZeroIsOne(t *testing.T) {
	if Barth(0) != 1 {
		t.Fatalf("Barth(0)=%v", Barth(0))
	}
	if Barth(6) != 0.5 {
		t.Fatalf("Barth(6)=%v", Barth(6))
	}
}

func TestRequiredFaceRecoversInput(t *testing.T) {
	ft, module, y, face := 1200.0, 2.0, 0.35, 15.0
	s, err := Bending(ft, face, module, y)
	if err != nil {
		t.Fatal(err)
	}
	got, err := RequiredFace(ft, module, y, s)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(got-face) > 1e-12 {
		t.Fatalf("face %v want %v", got, face)
	}
}

func TestPinionBendsMoreThanWheel(t *testing.T) {
	p := examplePair(t)
	s1, s2, err := PairBending(p, 800, 10)
	if err != nil {
		t.Fatal(err)
	}
	if s1 <= s2 {
		t.Fatalf("pinion %v wheel %v", s1, s2)
	}
	weak, who, err := WeakerTooth(p, 800, 10)
	if err != nil {
		t.Fatal(err)
	}
	if who != "gear1" || weak != s1 {
		t.Fatalf("weaker %s %v", who, weak)
	}
}

func TestSmallTeethRejected(t *testing.T) {
	if _, err := FormY(11); err == nil {
		t.Fatal("z=11")
	}
}
