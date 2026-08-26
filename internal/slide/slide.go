package slide

import (
	"errors"
	"fmt"
	"math"

	"gear-contact/internal/gears"
	"gear-contact/internal/hertz"
)

var ErrPositive = errors.New("slide: speed, load or thermal property must be positive")

func OmegaFromRPM(rpm float64) float64 {
	return rpm * 2.0 * math.Pi / 60.0
}

func WheelOmega(p *gears.Pair, pinionOmega float64) float64 {
	if p == nil || p.Gear2().Teeth == 0 {
		return 0
	}
	return pinionOmega * float64(p.Gear1().Teeth) / float64(p.Gear2().Teeth)
}

func PitchLineVelocity(p *gears.Pair, pinionOmega float64) float64 {
	if p == nil {
		return 0
	}
	return math.Abs(pinionOmega) * p.Gear1().PitchRadius()
}

func Sliding(p *gears.Pair, s, pinionOmega float64) float64 {
	if p == nil {
		return 0
	}
	w1 := pinionOmega
	w2 := WheelOmega(p, pinionOmega)
	return (w1 + w2) * s
}

func Rolling(p *gears.Pair, s, pinionOmega float64) float64 {
	if p == nil {
		return 0
	}
	w1 := math.Abs(pinionOmega)
	w2 := math.Abs(WheelOmega(p, pinionOmega))
	rho1 := hertz.PitchRho(p.Gear1()) + s
	rho2 := hertz.PitchRho(p.Gear2()) - s
	return 0.5 * (w1*math.Abs(rho1) + w2*math.Abs(rho2))
}

func SurfaceSpeeds(p *gears.Pair, s, pinionOmega float64) (v1, v2 float64) {
	w1 := math.Abs(pinionOmega)
	w2 := math.Abs(WheelOmega(p, pinionOmega))
	rho1 := hertz.PitchRho(p.Gear1()) + s
	rho2 := hertz.PitchRho(p.Gear2()) - s
	return w1 * rho1, w2 * rho2
}

func SpecificSliding(p *gears.Pair, s, pinionOmega float64) (g1, g2 float64, err error) {
	v1, v2 := SurfaceSpeeds(p, s, pinionOmega)
	vs := v1 - v2
	if math.Abs(v1) < 1e-15 || math.Abs(v2) < 1e-15 {
		return 0, 0, ErrPositive
	}
	return vs / v1, -vs / v2, nil
}

func FlashBlok(mu, w, vs, k, dens, heat, vSum float64) (float64, error) {
	if mu < 0 || w <= 0 || k <= 0 || dens <= 0 || heat <= 0 {
		return 0, ErrPositive
	}
	if vSum <= 0 {
		return 0, nil
	}
	peclet := math.Sqrt(k * dens * heat * vSum)
	if peclet <= 0 {
		return 0, ErrPositive
	}
	return 1.11 * mu * w * math.Abs(vs) / peclet, nil
}

func FlashAt(p *gears.Pair, s, pinionOmega, mu, fn, face float64) (float64, error) {
	if face <= 0 || fn <= 0 {
		return 0, ErrPositive
	}
	vs := Sliding(p, s, pinionOmega)
	if vs == 0 {
		return 0, nil
	}
	c, err := hertz.AlongLine(p, s, fn, face)
	if err != nil {
		return 0, err
	}
	v1, v2 := SurfaceSpeeds(p, s, pinionOmega)
	vSum := math.Abs(v1) + math.Abs(v2)
	w := fn / face
	return FlashBlok(mu, w, vs, 45.0, 7.85e-6, 460.0, vSum+c.HalfWidth*0)
}

func SlideRatio(p *gears.Pair, s, pinionOmega float64) float64 {
	vr := Rolling(p, s, pinionOmega)
	if vr == 0 {
		return 0
	}
	return math.Abs(Sliding(p, s, pinionOmega)) / vr
}

func MaxSlideOnPath(p *gears.Pair, sStart, sEnd, pinionOmega float64) float64 {
	a := math.Abs(Sliding(p, sStart, pinionOmega))
	b := math.Abs(Sliding(p, sEnd, pinionOmega))
	if a > b {
		return a
	}
	return b
}

func SignsOpposeAcrossPitch(p *gears.Pair, s, pinionOmega float64) error {
	if s == 0 {
		return errors.New("slide: need a non-zero station")
	}
	a := Sliding(p, s, pinionOmega)
	b := Sliding(p, -s, pinionOmega)
	if a*b >= 0 {
		return fmt.Errorf("slide: sliding should reverse across the pitch point: %g vs %g", a, b)
	}
	return nil
}

func FasterSpinScalesSlide(p *gears.Pair, s float64) error {
	slow := math.Abs(Sliding(p, s, 50))
	fast := math.Abs(Sliding(p, s, 100))
	if math.Abs(fast-2*slow) > 1e-9*math.Max(1, fast) {
		return fmt.Errorf("slide: sliding should scale with ω: %g vs 2×%g", fast, slow)
	}
	return nil
}

func PitchLineMatchesOmegaR(p *gears.Pair, pinionOmega float64) error {
	got := PitchLineVelocity(p, pinionOmega)
	want := math.Abs(pinionOmega) * p.Gear1().PitchRadius()
	if math.Abs(got-want) > 1e-12 {
		return fmt.Errorf("slide: pitch-line velocity %g != ωr %g", got, want)
	}
	return nil
}
