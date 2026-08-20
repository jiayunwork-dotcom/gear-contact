package angle

// Pressure angle interval. The accepted working range is a named value so
// callers and the README agree on exactly one definition instead of quoting
// the bounds independently.

// Interval is a closed range of angles in radians.
type Interval struct {
	// Lo is the lower bound of the interval, inclusive.
	Lo float64
	// Hi is the upper bound of the interval, inclusive.
	Hi float64
}

// WorkingInterval returns the pinned working pressure angle interval
// [14.5, 25] degrees, in radians.
func WorkingInterval() Interval {
	return Interval{Lo: MinPressureAngleRad, Hi: MaxPressureAngleRad}
}

// Contains reports whether the interval contains the angle, inclusive of both
// bounds. NaN is never contained.
func (iv Interval) Contains(rad float64) bool {
	if rad != rad { // NaN check
		return false
	}
	return rad >= iv.Lo && rad <= iv.Hi
}

// SpanDeg returns the width of the interval in degrees.
func (iv Interval) SpanDeg() float64 {
	return RadToDeg(iv.Hi - iv.Lo)
}

// MidDeg returns the middle of the interval in degrees.
func (iv Interval) MidDeg() float64 {
	return RadToDeg(0.5 * (iv.Lo + iv.Hi))
}

// DistanceTo returns how far the angle is outside the interval in radians;
// zero means the angle is inside. The mesh uses it to phrase warnings about
// angles that just miss the working range.
func (iv Interval) DistanceTo(rad float64) float64 {
	if iv.Contains(rad) {
		return 0
	}
	if rad < iv.Lo {
		return iv.Lo - rad
	}
	return rad - iv.Hi
}
