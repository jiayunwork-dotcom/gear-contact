package tiprel

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

func TestAddendumCutLowersEpsilon(t *testing.T) {
	p := pair(t)
	if err := ShorterAddendumLowersEpsilon(p, 1.0, 0.8); err != nil {
		t.Fatal(err)
	}
	if err := StillEngages(p, 0.85); err != nil {
		t.Fatal(err)
	}
	if err := PathTracksAddendum(p, 1.0, 0.8); err != nil {
		t.Fatal(err)
	}
}
