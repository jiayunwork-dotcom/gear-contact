package angle

import (
	"math"
	"strconv"
)

// Formatting helpers for angles. The mesh report prints all angles in degrees
// with three decimals, which keeps the output readable while still exposing a
// change of a tenth of a degree in the working pressure angle.
const (
	// ReportDecimals is the fixed number of decimals used in reports.
	ReportDecimals = 3
	// ReportWidth is the field width used to right-align numbers in the text
	// report. Keeping it consistent makes column-aligned output.
	ReportWidth = 10
)

// FormatDeg renders an angle given in radians as a degrees string with
// ReportDecimals decimals. The conversion and the formatting happen in the
// same place so no other package can mix units by accident.
func FormatDeg(rad float64) string {
	return strconv.FormatFloat(RadToDeg(rad), 'f', ReportDecimals, 64)
}

// FormatDegSigned is FormatDeg with an explicit '+' for positive values, used
// for profile shift coefficients where the sign carries meaning.
func FormatDegSigned(rad float64) string {
	deg := RadToDeg(rad)
	if deg >= 0 {
		return "+" + strconv.FormatFloat(deg, 'f', ReportDecimals, 64)
	}
	return strconv.FormatFloat(deg, 'f', ReportDecimals, 64)
}

// FloatShort renders a float with up to 6 significant digits, trimming
// trailing zeros. It is used by error messages where a full-width float would
// be noise.
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

// FormatRad renders an angle in radians with ReportDecimals decimals. It is
// used when a radian value (such as an involute result) has to be surfaced in
// a message without a degree conversion.
func FormatRad(rad float64) string {
	return strconv.FormatFloat(rad, 'f', ReportDecimals, 64)
}
