package gears

import "fmt"

type Circle struct {
	Name   string
	Radius float64
}

type CircleKind int

const (
	PitchCircleKind CircleKind = iota
	BaseCircleKind
	AddendumCircleKind
	DedendumCircleKind
)

func (k CircleKind) KindName() string {
	switch k {
	case PitchCircleKind:
		return "pitch"
	case BaseCircleKind:
		return "base"
	case AddendumCircleKind:
		return "addendum"
	case DedendumCircleKind:
		return "dedendum"
	default:
		return "unknown"
	}
}

func CircleOf(g Gear, kind CircleKind) Circle {
	switch kind {
	case PitchCircleKind:
		return Circle{Name: kind.KindName(), Radius: g.PitchRadius()}
	case BaseCircleKind:
		return Circle{Name: kind.KindName(), Radius: g.BaseRadius()}
	case AddendumCircleKind:
		return Circle{Name: kind.KindName(), Radius: g.AddendumRadius()}
	case DedendumCircleKind:
		return Circle{Name: kind.KindName(), Radius: g.DedendumRadius()}
	default:
		return Circle{Name: kind.KindName(), Radius: 0}
	}
}

func AllCircles(g Gear) []Circle {
	return []Circle{
		CircleOf(g, PitchCircleKind),
		CircleOf(g, BaseCircleKind),
		CircleOf(g, AddendumCircleKind),
		CircleOf(g, DedendumCircleKind),
	}
}

func TipOutsidePitch(g Gear) error {
	if g.AddendumRadius() <= g.PitchRadius() {
		return fmt.Errorf("gears: addendum must sit outside the pitch circle")
	}
	if g.BaseRadius() >= g.PitchRadius() {
		return fmt.Errorf("gears: base circle must sit inside the pitch circle")
	}
	if g.DedendumRadius() >= g.PitchRadius() {
		return fmt.Errorf("gears: dedendum must sit inside the pitch circle")
	}
	return nil
}

func CirclesOrdered(g Gear) error {
	if !(g.DedendumRadius() < g.BaseRadius() || g.Teeth < 20) {
		if g.DedendumRadius() >= g.PitchRadius() {
			return fmt.Errorf("gears: root above pitch")
		}
	}
	if g.AddendumRadius() <= g.BaseRadius() {
		return fmt.Errorf("gears: tip does not reach outside the base circle")
	}
	return nil
}
