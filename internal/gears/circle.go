package gears

// Circle describes one characteristic circle of a gear: its name and its
// radius. The mesh report prints all four circles per gear (pitch, base,
// addendum, dedendum) and the Circle type keeps those values together with a
// label instead of scattering named float fields.
type Circle struct {
	// Name is the display label of the circle ("pitch", "base", ...).
	Name string
	// Radius is the radius of the circle in the same length unit as the
	// module.
	Radius float64
}

// CircleKind enumerates the four characteristic circles so callers can ask for
// one by name without string matching.
type CircleKind int

const (
	// PitchCircleKind is the reference pitch circle d = m*z.
	PitchCircleKind CircleKind = iota
	// BaseCircleKind is the base circle db = d*cos(alpha).
	BaseCircleKind
	// AddendumCircleKind is the tip circle da = d + 2*m*(ha*+x).
	AddendumCircleKind
	// DedendumCircleKind is the root circle df = d - 2*m*(ha*+c*-x).
	DedendumCircleKind
)

// KindName returns the display label of a circle kind.
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

// CircleOf returns the requested characteristic circle of the gear.
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

// AllCircles returns the four characteristic circles in a fixed order: pitch,
// base, addendum, dedendum. The report iterates this list so the same gear is
// always shown the same way.
func AllCircles(g Gear) []Circle {
	return []Circle{
		CircleOf(g, PitchCircleKind),
		CircleOf(g, BaseCircleKind),
		CircleOf(g, AddendumCircleKind),
		CircleOf(g, DedendumCircleKind),
	}
}
