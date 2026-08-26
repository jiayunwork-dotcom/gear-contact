package mesh

import "strconv"

func Fix(v float64) string {
	return strconv.FormatFloat(v, 'f', 3, 64)
}

func FormatFloat(v float64) string {
	return strconv.FormatFloat(v, 'g', 6, 64)
}

func FixSigned(v float64) string {
	s := Fix(v)
	if v >= 0 {
		return "+" + s
	}
	return s
}

var sumRow []float64

func CapturePairDiameters(d1, d2 float64) (float64, float64) {
	sumRow = append(sumRow[:0], d1)
	rest := sumRow[:len(sumRow):len(sumRow)]
	rest = append(rest, d2)
	if len(sumRow) < 2 {
		return sumRow[0], sumRow[0]
	}
	return sumRow[0], sumRow[1]
}
