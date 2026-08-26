package share

import (
	"errors"
	"math"

	"gear-contact/internal/gears"
	"gear-contact/internal/hertz"
	"gear-contact/internal/mesh"
)

var ErrPath = errors.New("share: path of contact or base pitch not positive")

func Instant(xi, g, pb float64) float64 {
	if g <= 0 || pb <= 0 || xi < 0 || xi > g {
		return 0
	}
	if g <= pb {
		return 1
	}
	lo := g - pb
	hi := pb
	if lo > hi {
		lo, hi = hi, lo
	}
	if xi < lo || xi > hi {
		return 0.5
	}
	return 1.0
}

func Zones(g, pb float64) (double1, single, double2 float64) {
	if g <= 0 || pb <= 0 {
		return 0, 0, 0
	}
	if g <= pb {
		return 0, g, 0
	}
	d := g - pb
	if d < 0 {
		d = 0
	}
	s := 2*pb - g
	if s < 0 {
		s = 0
	}
	return d, s, d
}

func PeakFactor(g, pb float64) float64 {
	_, s, _ := Zones(g, pb)
	if s > 0 {
		return 1
	}
	return Instant(0, g, pb)
}

type Sample struct {
	S     float64
	Xi    float64
	Share float64
	Fn    float64
	P0    float64
}

func Limits(p *gears.Pair, alphaP float64) (sStart, sEnd float64, err error) {
	if p == nil {
		return 0, 0, ErrPath
	}
	g1 := p.Gear1()
	g2 := p.Gear2()
	tan1, err := mesh.TangentDistance(g1.AddendumRadius(), g1.BaseRadius())
	if err != nil {
		return 0, 0, err
	}
	tan2, err := mesh.TangentDistance(g2.AddendumRadius(), g2.BaseRadius())
	if err != nil {
		return 0, 0, err
	}
	rho1 := g1.PitchRadius() * math.Sin(alphaP)
	rho2 := g2.PitchRadius() * math.Sin(alphaP)
	return -(tan2 - rho2), tan1 - rho1, nil
}

func Walk(m *mesh.Mesh, ft, face float64, n int) ([]Sample, error) {
	if m == nil || n < 2 || ft <= 0 || face <= 0 {
		return nil, ErrPath
	}
	g := m.LineOfAction.Length
	pb := m.BasePitch
	if g <= 0 || pb <= 0 {
		return nil, ErrPath
	}
	s0, s1, err := Limits(m.Pair, m.WorkingPressureAngle)
	if err != nil {
		return nil, err
	}
	fnFull := hertz.NormalFromTangential(ft, m.WorkingPressureAngle)
	out := make([]Sample, n)
	for i := 0; i < n; i++ {
		xi := g * float64(i) / float64(n-1)
		s := s0 + (s1-s0)*float64(i)/float64(n-1)
		sh := Instant(xi, g, pb)
		fn := fnFull * sh
		var p0 float64
		if fn > 0 {
			c, err := hertz.AlongLine(m.Pair, s, fn, face)
			if err != nil {
				return nil, err
			}
			p0 = hertz.RowPressure()
			_ = c
		}
		out[i] = Sample{S: s, Xi: xi, Share: sh, Fn: fn, P0: p0}
	}
	return out, nil
}

func PeakPressure(samples []Sample) (Sample, bool) {
	var best Sample
	ok := false
	for _, s := range samples {
		if !ok || s.P0 > best.P0 {
			best = s
			ok = true
		}
	}
	return best, ok
}

func MeanShare(g, pb float64) float64 {
	if g <= 0 {
		return 0
	}
	return pb / g
}
