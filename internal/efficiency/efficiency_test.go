package efficiency

import (
	"testing"

	"gear-contact/internal/gears"
)

func pair(t *testing.T) *gears.Pair {
	t.Helper()
	p, err := gears.FromJSON([]byte(`{"module":1,"pressure_angle_deg":20,"gear1":{"teeth":20},"gear2":{"teeth":40}}`))
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func TestPitchIsUnityAndMuHurts(t *testing.T) {
	p := pair(t)
	if err := PitchIsUnity(p, 100, 0.06); err != nil {
		t.Fatal(err)
	}
	if err := HigherMuWorse(p, 100); err != nil {
		t.Fatal(err)
	}
	if err := ZeroMuNoLoss(p, 80); err != nil {
		t.Fatal(err)
	}
}
