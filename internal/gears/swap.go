package gears

// Swap builds a new pair with the driving and driven gears exchanged. The
// crossing rule of the engagement calculation says the centre distance and
// the contact ratio are unchanged by the swap; returning a real Pair lets the
// mesh package verify that rule on an actual object instead of on hand-built
// values.
func (p *Pair) Swap() *Pair {
	return &Pair{gear1: p.gear2, gear2: p.gear1}
}

// SamePair reports whether two pairs are geometrically identical: equal
// modules, pressure angles, tooth counts, coefficients and shifts on both
// gears. It is used by tests to confirm that a swap round trip (swapping
// twice) restores the original pair exactly.
func SamePair(a, b *Pair) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Module() == b.Module() &&
		a.PressureAngle() == b.PressureAngle() &&
		a.Gear1() == b.Gear1() &&
		a.Gear2() == b.Gear2()
}

// GearAt returns the gear by index 1 or 2; the zero value is returned for any
// other index. The report loop uses it to walk the two gears uniformly.
func (p *Pair) GearAt(index int) Gear {
	switch index {
	case 1:
		return p.gear1
	case 2:
		return p.gear2
	default:
		return Gear{}
	}
}
