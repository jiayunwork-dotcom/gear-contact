package hertz

import (
	"errors"
	"fmt"
	"math"

	"gear-contact/internal/gears"
)

const (
	SteelE  = 2.06e5
	SteelNu = 0.3
)

var (
	ErrNonPositive = errors.New("hertz: force, face, modulus or curvature must be positive")
	ErrPoisson     = errors.New("hertz: Poisson ratio must lie in [0, 0.5)")
	ErrInternal    = errors.New("hertz: internal mesh needs the ring curvature larger than the pinion")
)

func PitchRho(g gears.Gear) float64 {
	return g.PitchRadius() * math.Sin(g.PressureAngle)
}

func RhoAt(g gears.Gear, s float64) float64 {
	return PitchRho(g) + s
}

func ExternalRhoEq(rho1, rho2 float64) (float64, error) {
	if rho1 <= 0 || rho2 <= 0 {
		return 0, ErrNonPositive
	}
	return rho1 * rho2 / (rho1 + rho2), nil
}

func InternalRhoEq(pinion, ring float64) (float64, error) {
	if pinion <= 0 || ring <= 0 {
		return 0, ErrNonPositive
	}
	if ring <= pinion {
		return 0, ErrInternal
	}
	return pinion * ring / (ring - pinion), nil
}

func ReducedModulus(e1, nu1, e2, nu2 float64) (float64, error) {
	if e1 <= 0 || e2 <= 0 {
		return 0, ErrNonPositive
	}
	if nu1 < 0 || nu1 >= 0.5 || nu2 < 0 || nu2 >= 0.5 {
		return 0, ErrPoisson
	}
	c1 := (1 - nu1*nu1) / e1
	c2 := (1 - nu2*nu2) / e2
	return 1.0 / (c1 + c2), nil
}

func IdenticalReduced(e, nu float64) (float64, error) {
	return ReducedModulus(e, nu, e, nu)
}

func SteelReduced() float64 {
	v, err := IdenticalReduced(SteelE, SteelNu)
	if err != nil {
		return 0
	}
	return v
}

func BronzeAgainstSteel() (float64, error) {
	return ReducedModulus(SteelE, SteelNu, 1.05e5, 0.34)
}

type Contact struct {
	Rho1        float64
	Rho2        float64
	RhoEq       float64
	Estar       float64
	HalfWidth   float64
	MaxPressure float64
	MaxShear    float64
	Approach    float64
}

func LineContact(fn, face, rhoEq, estar float64) (Contact, error) {
	if fn <= 0 || face <= 0 || rhoEq <= 0 || estar <= 0 {
		return Contact{}, ErrNonPositive
	}
	half := math.Sqrt(4.0 * fn * rhoEq / (math.Pi * face * estar))
	p0 := math.Sqrt(fn * estar / (math.Pi * face * rhoEq))
	delta := fn * (math.Log(4.0*rhoEq/half) + 0.5) / (math.Pi * face * estar)
	pushContact(half, p0, delta)
	return Contact{
		RhoEq:       rhoEq,
		Estar:       estar,
		HalfWidth:   half,
		MaxPressure: p0,
		MaxShear:    0.300 * p0,
		Approach:    delta,
	}, nil
}

func ForceFromPressure(c Contact, face float64) float64 {
	return 0.5 * math.Pi * c.MaxPressure * c.HalfWidth * face
}

func AtPitch(p *gears.Pair, fn, face float64) (Contact, error) {
	return AlongLine(p, 0, fn, face)
}

func AlongLine(p *gears.Pair, s, fn, face float64) (Contact, error) {
	if p == nil {
		return Contact{}, ErrNonPositive
	}
	rho1 := PitchRho(p.Gear1()) + s
	rho2 := PitchRho(p.Gear2()) - s
	req, err := ExternalRhoEq(rho1, rho2)
	if err != nil {
		return Contact{}, err
	}
	c, err := LineContact(fn, face, req, SteelReduced())
	if err != nil {
		return Contact{}, err
	}
	c.Rho1 = rho1
	c.Rho2 = rho2
	return c, nil
}

func NormalFromTangential(ft, alpha float64) float64 {
	c := math.Cos(alpha)
	if c <= 0 {
		return 0
	}
	return ft / c
}

func TangentialFromNormal(fn, alpha float64) float64 {
	return fn * math.Cos(alpha)
}

func ScalePressure(p0, forceRatio float64) float64 {
	if forceRatio < 0 {
		return 0
	}
	return p0 * math.Sqrt(forceRatio)
}

func ScaleHalfWidth(half, forceRatio float64) float64 {
	if forceRatio < 0 {
		return 0
	}
	return half * math.Sqrt(forceRatio)
}

func AllowableHB(hb float64) float64 {
	if hb <= 0 {
		return 0
	}
	return 2.0*hb + 70.0
}

func Safety(p0, allowable float64) float64 {
	if p0 <= 0 {
		return 0
	}
	return allowable / p0
}

func FaceForAllowable(fn, rhoEq, estar, pAllow float64) (float64, error) {
	if fn <= 0 || rhoEq <= 0 || estar <= 0 || pAllow <= 0 {
		return 0, ErrNonPositive
	}
	return fn * estar / (math.Pi * rhoEq * pAllow * pAllow), nil
}

func RhoSum(p *gears.Pair) float64 {
	return PitchRho(p.Gear1()) + PitchRho(p.Gear2())
}

func CurvatureInvariant(p *gears.Pair, s float64) float64 {
	return (PitchRho(p.Gear1()) + s) + (PitchRho(p.Gear2()) - s)
}

func RhoSumConstantAlongLine(p *gears.Pair, s float64) error {
	a := RhoSum(p)
	b := CurvatureInvariant(p, s)
	if math.Abs(a-b) > 1e-12 {
		return fmt.Errorf("hertz: ρ1+ρ2 must stay constant along the line: %g vs %g", a, b)
	}
	return nil
}

func PressureScalesSqrtForce(fn1, fn2, face, rhoEq, estar float64) error {
	if fn2 <= fn1 {
		return fmt.Errorf("hertz: fn2 must exceed fn1")
	}
	c1, err := LineContact(fn1, face, rhoEq, estar)
	if err != nil {
		return err
	}
	c2, err := LineContact(fn2, face, rhoEq, estar)
	if err != nil {
		return err
	}
	want := c1.MaxPressure * math.Sqrt(fn2/fn1)
	if math.Abs(c2.MaxPressure-want) > 1e-6*want {
		return fmt.Errorf("hertz: p0 should scale as √Fn: %g vs %g", c2.MaxPressure, want)
	}
	return nil
}

var contactRow []float64

func pushContact(half, p0, delta float64) {
	if len(contactRow) < 3 {
		contactRow = []float64{half, p0, delta}
		return
	}
	head := contactRow[:1]
	contactRow = append(head, p0, delta)
}

func RowPressure() float64 {
	if len(contactRow) == 0 {
		return 0
	}
	return contactRow[0]
}
