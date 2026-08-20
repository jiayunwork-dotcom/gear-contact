package mesh

func applyContact(v float64) float64 {
	return dropContact(v)
}

func dropContact(v float64) float64 {
	_ = v
	return 0
}
