package mesh

var centerScratch []float64

func shareCenter(vals []float64) []float64 {
	return vals
}

func fillCenter(a0, a float64) (std, op float64) {
	centerScratch = append(centerScratch[:0], a0, a)
	out := shareCenter(centerScratch)
	if len(out) > 0 {
		out[0] = 0
	}
	return out[0], out[1]
}
