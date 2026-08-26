package share

import (
	"math"
	"os"
	"testing"

	"gear-contact/internal/gears"
	"gear-contact/internal/hertz"
	"gear-contact/internal/mesh"
)

func exampleMesh(t *testing.T) *mesh.Mesh {
	t.Helper()
	var data []byte
	var err error
	for _, c := range []string{"../../example/spur-20.json", "example/spur-20.json"} {
		data, err = os.ReadFile(c)
		if err == nil {
			break
		}
	}
	if err != nil {
		t.Fatal(err)
	}
	p, err := gears.FromJSON(data)
	if err != nil {
		t.Fatal(err)
	}
	m, err := mesh.Compute(p)
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func TestDoublePairHalfLoadAtStart(t *testing.T) {
	g, pb := 4.827, 2.952
	if Instant(0, g, pb) != 0.5 {
		t.Fatalf("start share %v", Instant(0, g, pb))
	}
	if Instant(g, g, pb) != 0.5 {
		t.Fatalf("end share %v", Instant(g, g, pb))
	}
}

func TestSinglePairFullLoadInMiddle(t *testing.T) {
	g, pb := 4.827, 2.952
	mid := 0.5 * g
	if Instant(mid, g, pb) != 1 {
		t.Fatalf("mid share %v", Instant(mid, g, pb))
	}
}

func TestShareZeroOutsidePath(t *testing.T) {
	if Instant(-0.1, 4, 2.5) != 0 || Instant(4.1, 4, 2.5) != 0 {
		t.Fatal("outside")
	}
}

func TestMeanShareIsOneOverEpsilon(t *testing.T) {
	m := exampleMesh(t)
	got := MeanShare(m.LineOfAction.Length, m.BasePitch)
	want := 1.0 / m.ContactRatio
	if math.Abs(got-want) > 1e-12 {
		t.Fatalf("mean %v want %v", got, want)
	}
}

func TestPeakHertzOffPitchDueToLoadShare(t *testing.T) {
	m := exampleMesh(t)
	samples, err := Walk(m, 1200, 10, 41)
	if err != nil {
		t.Fatal(err)
	}
	peak, ok := PeakPressure(samples)
	if !ok {
		t.Fatal("empty")
	}
	pitch, err := hertz.AtPitch(m.Pair, hertz.NormalFromTangential(1200, m.WorkingPressureAngle), 10)
	if err != nil {
		t.Fatal(err)
	}
	if peak.P0 <= pitch.MaxPressure {
		t.Fatalf("peak %v at s=%v not above undivided pitch %v", peak.P0, peak.S, pitch.MaxPressure)
	}
	if math.Abs(peak.S) < 1e-9 && peak.Share != 1 {
		t.Fatal("expected the single-pair zone, not s=0 with half load")
	}
}

func TestLimitsSpanMatchesLineLength(t *testing.T) {
	m := exampleMesh(t)
	s0, s1, err := Limits(m.Pair, m.WorkingPressureAngle)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs((s1-s0)-m.LineOfAction.Length) > 1e-9 {
		t.Fatalf("span %v length %v", s1-s0, m.LineOfAction.Length)
	}
}
