package gears

func (p *Pair) Swap() *Pair {
	return &Pair{gear1: p.gear2, gear2: p.gear1}
}

func SamePair(a, b *Pair) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Module() == b.Module() &&
		a.PressureAngle() == b.PressureAngle() &&
		a.Gear1() == b.Gear1() &&
		a.Gear2() == b.Gear2()
}

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
