package slide

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

func TestPitchPointSlidingIsZero(t *testing.T) {
	p := examplePair(t)
	vs := Sliding(p, 0, 100)
	if math.Abs(vs) > 1e-12 {
		t.Fatalf("pitch slide %v", vs)
	}
}

func TestSlidingReversesAcrossPitch(t *testing.T) {
	p := examplePair(t)
	a := Sliding(p, 0.4, 80)
	b := Sliding(p, -0.4, 80)
	if a == 0 || a != -b {
		t.Fatalf("a=%v b=%v", a, b)
	}
}

func TestWheelOmegaInverseRatio(t *testing.T) {
	p := examplePair(t)
	w2 := WheelOmega(p, 40)
	want := 40.0 * 20.0 / 40.0
	if math.Abs(w2-want) > 1e-12 {
		t.Fatalf("w2 %v want %v", w2, want)
	}
}

func TestFlashZeroWhenMuZero(t *testing.T) {
	p := examplePair(t)
	dt, err := FlashAt(p, 0.3, 50, 0, 1000, 10)
	if err != nil {
		t.Fatal(err)
	}
	if dt != 0 {
		t.Fatalf("flash %v", dt)
	}
}

func TestFlashZeroAtPitch(t *testing.T) {
	p := examplePair(t)
	dt, err := FlashAt(p, 0, 50, 0.08, 1000, 10)
	if err != nil {
		t.Fatal(err)
	}
	if dt != 0 {
		t.Fatalf("pitch flash %v", dt)
	}
}

func TestSpecificSlidingSignsOppose(t *testing.T) {
	p := examplePair(t)
	g1, g2, err := SpecificSliding(p, 0.25, 30)
	if err != nil {
		t.Fatal(err)
	}
	if g1*g2 >= 0 {
		t.Fatalf("g1=%v g2=%v", g1, g2)
	}
}

func TestOmegaFromRPM(t *testing.T) {
	w := OmegaFromRPM(60)
	if math.Abs(w-2*math.Pi) > 1e-12 {
		t.Fatalf("omega %v", w)
	}
}
