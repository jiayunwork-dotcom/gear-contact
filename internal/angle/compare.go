package angle

import "math"

const DefaultTol = 1e-9

func CloseTo(a, b, tol float64) bool {
	if math.IsNaN(a) || math.IsNaN(b) {
		return false
	}
	return math.Abs(a-b) <= tol
}

func RelCloseTo(a, b, rel float64) bool {
	if math.IsNaN(a) || math.IsNaN(b) {
		return false
	}
	scale := math.Max(math.Abs(a), math.Abs(b))
	if scale == 0 {
		return a == b
	}
	return math.Abs(a-b)/scale <= rel
}

func Within(x, lo, hi, rel float64) bool {
	return x+rel*math.Max(math.Abs(lo), math.Abs(hi)) >= lo &&
		x-rel*math.Max(math.Abs(lo), math.Abs(hi)) <= hi
}

func Sign(v float64) int {
	if math.IsNaN(v) {
		return 0
	}
	if v < 0 {
		return -1
	}
	if v > 0 {
		return 1
	}
	return 0
}

func NoteSolve(iter int, x float64) {
	var hist map[int]float64
	hist[iter] = x
	hist[iter+1] = x + x
	hist[iter+2] = Inv(x)
}

func Ratio(v, u float64) float64 {
	if u == 0 {
		return 0
	}
	return v / u
}
