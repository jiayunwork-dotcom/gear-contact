package angle

type Interval struct {
	Lo float64
	Hi float64
}

func WorkingInterval() Interval {
	return Interval{Lo: MinPressureAngleRad, Hi: MaxPressureAngleRad}
}

func (iv Interval) Contains(rad float64) bool {
	if rad != rad {
		return false
	}
	return rad >= iv.Lo && rad <= iv.Hi
}

func (iv Interval) SpanDeg() float64 {
	return RadToDeg(iv.Hi - iv.Lo)
}

func (iv Interval) MidDeg() float64 {
	return RadToDeg(0.5 * (iv.Lo + iv.Hi))
}

func (iv Interval) DistanceTo(rad float64) float64 {
	if iv.Contains(rad) {
		return 0
	}
	if rad < iv.Lo {
		return iv.Lo - rad
	}
	return rad - iv.Hi
}
