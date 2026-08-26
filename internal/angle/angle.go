package angle

import "math"

const TwoPi = 2 * math.Pi

func DegToRad(deg float64) float64 {
	return deg * math.Pi / 180.0
}

func RadToDeg(rad float64) float64 {
	return rad * 180.0 / math.Pi
}

func Normalize(rad float64) float64 {
	rad = math.Mod(rad, TwoPi)
	if rad < 0 {
		rad += TwoPi
	}
	return rad
}

func SinCos(rad float64) (sin, cos float64) {
	return math.Sincos(rad)
}
