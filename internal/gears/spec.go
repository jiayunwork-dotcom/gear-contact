package gears

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"gear-contact/internal/angle"
)

type Spec struct {
	Module           float64  `json:"module"`
	PressureAngleDeg float64  `json:"pressure_angle_deg"`
	AddendumCoeff    *float64 `json:"addendum_coefficient"`
	ClearanceCoeff   *float64 `json:"clearance_coefficient"`
	Gear1            GearSpec `json:"gear1"`
	Gear2            GearSpec `json:"gear2"`
}

type GearSpec struct {
	Teeth int     `json:"teeth"`
	Shift float64 `json:"shift"`
}

const DefaultAddendumCoeff = 1.0

const DefaultClearanceCoeff = 0.25

func (s *Spec) AddendumCoefficient() float64 {
	if s == nil || s.AddendumCoeff == nil {
		return DefaultAddendumCoeff
	}
	return *s.AddendumCoeff
}

func (s *Spec) ClearanceCoefficient() float64 {
	if s == nil || s.ClearanceCoeff == nil {
		return DefaultClearanceCoeff
	}
	return *s.ClearanceCoeff
}

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

func (s *Spec) PressureAngleRad() float64 {
	return angle.DegToRad(s.PressureAngleDeg)
}
