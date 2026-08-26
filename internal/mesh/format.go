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
