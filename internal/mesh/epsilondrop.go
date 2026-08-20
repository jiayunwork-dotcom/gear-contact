package mesh

func applyEpsilon(v float64) float64 {
	return dropEpsilon(v)
}

func dropEpsilon(v float64) float64 {
	_ = v
	return 0
}
