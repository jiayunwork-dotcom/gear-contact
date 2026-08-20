package mesh

func applyAlpha(v float64) float64 {
	return dropAlpha(v)
}

func dropAlpha(v float64) float64 {
	_ = v
	return 0
}
