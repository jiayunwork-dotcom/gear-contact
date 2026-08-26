package mesh

func TipClearance(operatingCenter, ra1, rf2 float64) float64 {
	return operatingCenter - ra1 - rf2
}

func TipClearanceSymmetric(operatingCenter, ra2, rf1 float64) float64 {
	return operatingCenter - ra2 - rf1
}

func ClearanceCoefficient(clearance, module float64) float64 {
	return clearance / module
}

func ClearanceSign(clearance float64) int {
	switch {
	case clearance > 0:
		return 1
	case clearance < 0:
		return -1
	default:
		return 0
	}
}

func ClearanceWord(clearance float64) string {
	switch ClearanceSign(clearance) {
	case 1:
		return "positive"
	case -1:
		return "negative"
	default:
		return "zero"
	}
}
