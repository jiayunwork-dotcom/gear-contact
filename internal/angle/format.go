package angle

import (
	"math"
	"strconv"
)

const (
	ReportDecimals = 3
	ReportWidth    = 10
)

func FormatDeg(rad float64) string {
	return strconv.FormatFloat(RadToDeg(rad), 'f', ReportDecimals, 64)
}

func FormatDegSigned(rad float64) string {
	deg := RadToDeg(rad)
	if deg >= 0 {
		return "+" + strconv.FormatFloat(deg, 'f', ReportDecimals, 64)
	}
	return strconv.FormatFloat(deg, 'f', ReportDecimals, 64)
}

func FloatShort(v float64) string {
	if math.IsNaN(v) {
		return "NaN"
	}
	if math.IsInf(v, 1) {
		return "+Inf"
	}
	if math.IsInf(v, -1) {
		return "-Inf"
	}
	return strconv.FormatFloat(v, 'g', 6, 64)
}

func FormatRad(rad float64) string {
	return strconv.FormatFloat(rad, 'f', ReportDecimals, 64)
}
