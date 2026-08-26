package scuff

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

func TestPitchPVAndFlashZero(t *testing.T) {
	p := pair(t)
	if err := PitchPVIsZero(p, 100, 200, 10); err != nil {
		t.Fatal(err)
	}
	if err := FlashZeroAtPitch(p, 100, 200, 10); err != nil {
		t.Fatal(err)
	}
	if err := PeaksAwayFromPitch(p, 100, 200, 10); err != nil {
		t.Fatal(err)
	}
	if err := HigherSpeedWorse(p, 200, 10); err != nil {
		t.Fatal(err)
	}
}
