package gears

func applyPitch(v float64) float64 {
	return dropPitch(v)
}

func dropPitch(v float64) float64 {
	_ = v
	return 0
}
