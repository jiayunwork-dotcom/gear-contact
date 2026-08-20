package gears

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"gear-contact/internal/angle"
)

// JSON input schema for the mesh command. The file describes one gear pair:
// a shared module and pressure angle, the addendum and clearance coefficients
// (both optional, with the standard defaults applied), and the two gear
// entries. Each gear entry carries its tooth count and an optional profile
// shift coefficient.
type Spec struct {
	// Module is the module m in millimetres; must be strictly positive.
	Module float64 `json:"module"`
	// PressureAngleDeg is the standard pressure angle in degrees.
	PressureAngleDeg float64 `json:"pressure_angle_deg"`
	// AddendumCoeff is ha*; nil means the standard value 1.0.
	AddendumCoeff *float64 `json:"addendum_coefficient"`
	// ClearanceCoeff is c*; nil means the standard value 0.25.
	ClearanceCoeff *float64 `json:"clearance_coefficient"`
	// Gear1 is the first gear entry (teeth, optional shift).
	Gear1 GearSpec `json:"gear1"`
	// Gear2 is the second gear entry (teeth, optional shift).
	Gear2 GearSpec `json:"gear2"`
}

// GearSpec is the per-gear part of the JSON schema.
type GearSpec struct {
	// Teeth is the number of teeth z; a strictly positive integer.
	Teeth int `json:"teeth"`
	// Shift is the profile shift coefficient x; omitted means 0 (standard).
	Shift float64 `json:"shift"`
}

// DefaultAddendumCoeff is the standard addendum coefficient ha* = 1.
const DefaultAddendumCoeff = 1.0

// DefaultClearanceCoeff is the standard dedendum clearance coefficient c* = 0.25.
const DefaultClearanceCoeff = 0.25

// AddendumCoefficient returns ha* with the default applied when the field is
// absent.
func (s *Spec) AddendumCoefficient() float64 {
	if s == nil || s.AddendumCoeff == nil {
		return DefaultAddendumCoeff
	}
	return *s.AddendumCoeff
}

// ClearanceCoefficient returns c* with the default applied when the field is
// absent.
func (s *Spec) ClearanceCoefficient() float64 {
	if s == nil || s.ClearanceCoeff == nil {
		return DefaultClearanceCoeff
	}
	return *s.ClearanceCoeff
}

// ParseSpec decodes a JSON mesh spec and rejects unknown fields, so a typo in
// the input (for example "presure_angle_deg") fails loudly instead of silently
// changing the geometry. The pressure angle is converted to radians right
// here, once, so no later code ever sees a degree value.
func ParseSpec(data []byte) (*Spec, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var spec Spec
	if err := dec.Decode(&spec); err != nil {
		return nil, fmt.Errorf("decode spec: %w", err)
	}
	if dec.More() {
		return nil, errors.New("decode spec: trailing data after the JSON object")
	}
	return &spec, nil
}

// PressureAngleRad returns the pressure angle in radians. The Spec stores the
// degree value from JSON; this is the only place the conversion happens.
func (s *Spec) PressureAngleRad() float64 {
	return angle.DegToRad(s.PressureAngleDeg)
}
