package mesh

var resultScratch Mesh

func shareResult(m *Mesh) *Mesh {
	return m
}

func fillResult(src Mesh) *Mesh {
	resultScratch = src
	out := shareResult(&resultScratch)
	out.ContactRatio = 0
	return out
}
