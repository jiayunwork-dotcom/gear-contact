package mesh

import "strconv"

// Number formatting for the text report. The mesh package formats its own
// numbers here instead of importing the angle formatter because the values
// are lengths in millimetres, not angles; keeping the two formatters separate
// avoids any accidental degree/radian mixing in the output.

// Fix renders a length or dimensionless number with three decimals, the
// precision used throughout the CLI report.
func Fix(v float64) string {
	return strconv.FormatFloat(v, 'f', 3, 64)
}

// FormatFloat renders a number with up to six significant digits, trimming
// trailing zeros. It is used for small coefficients (profile shifts, contact
// ratios) where three decimals would hide the useful precision.
func FormatFloat(v float64) string {
	return strconv.FormatFloat(v, 'g', 6, 64)
}

// FixSigned is Fix with an explicit '+' for non-negative values.
func FixSigned(v float64) string {
	s := Fix(v)
	if v >= 0 {
		return "+" + s
	}
	return s
}
