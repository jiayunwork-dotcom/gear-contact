package backlash

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

func TestStandardHasZero(t *testing.T) {
	p := pair(t)
	if err := StandardHasZero(p); err != nil {
		t.Fatal(err)
	}
	if err := WiderCenterMoreBacklash(p, 0.05, 0.12); err != nil {
		t.Fatal(err)
	}
	if err := NormalUsesCos(p, 0.08); err != nil {
		t.Fatal(err)
	}
}
