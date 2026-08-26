package mesh

import (
	"math"

	"gear-contact/internal/angle"
)

func StandardCenterDistance(module float64, z1, z2 int) float64 {
	return 0.5 * module * float64(z1+z2)
}

func OperatingCenterDistance(a0, alpha, alphaP float64) float64 {
	cosA := math.Cos(alpha)
	cosAP := math.Cos(alphaP)
	if cosAP == 0 {
		return math.Inf(1)
	}
	return a0 * cosA / cosAP
}

func OperatingCenterInverse(a, alpha, alphaP float64) float64 {
	cosA := math.Cos(alpha)
	if cosA == 0 {
		return math.Inf(1)
	}
	return a * math.Cos(alphaP) / cosA
}

func CentreDistanceRatio(a, a0 float64) float64 {
	return angle.Ratio(a, a0)
}
