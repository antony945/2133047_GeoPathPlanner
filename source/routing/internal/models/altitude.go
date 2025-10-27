package models

import (
	"math"
	"strings"
)

const (
	// TODO: To change
	DEFAULT_MIN_ALT = -999999
	DEFAULT_MAX_ALT = +999999
	DEFAULT_ALT = 100
)

type Altitude struct {
	Value float64      `json:"value"`
	Unit  AltitudeUnit `json:"unit"`
}

func NewAltitude(value float64, unit AltitudeUnit) (Altitude, error) {
	a := Altitude{
		Value: value,
		Unit: AltitudeUnit(strings.ToLower(string(unit))),
	}
	if err := a.Validate(); err != nil {
		return Altitude{}, err
	}

	return a, nil
}

func MustNewAltitude(value float64, unit AltitudeUnit) (Altitude) {
	a, err := NewAltitude(value, unit)
	if err != nil {
		panic(err)
	}

	return a
}

// Validate Altitude (unit must be just MT or FT)
func (a Altitude) Validate() error {
	return a.Unit.Validate()
}

// Convert to MT or FT
func (a Altitude) ConvertTo(target AltitudeUnit) Altitude {
	if a.Unit == target {
		return a
	}
	if a.Unit == MT && target == FT {
		alt, _ := NewAltitude(a.Value*MT_TO_FT, target)
		return alt
	}
	if a.Unit == FT && target == MT {
		alt, _ := NewAltitude(a.Value/MT_TO_FT, target)
		return alt
	}
	return a
}

// Normalize to default unit of measure (MT)
func (a Altitude) Normalize() Altitude {
	return a.ConvertTo(MT)
}

// Subtract calculates the difference between two altitudes
// It converts the second altitude (b) to the same unit as the first (a)
// and returns a new Altitude instance with the subtraction result.
// The resulting Altitude maintains the same unit as the first altitude (a).
func (a Altitude) Subtract(b Altitude) Altitude {
	bConverted := b.ConvertTo(a.Unit)
	result, _ := NewAltitude(a.Value - bConverted.Value, a.Unit)
	return result
}

func (a Altitude) Distance(b Altitude) Altitude {
	result := a.Subtract(b)
	result.Value = math.Abs(result.Value)
	return result
}

func (a Altitude) Compare(b Altitude) int {
	tmp := int(a.Subtract(b).Value)
	if tmp == 0 {
        return 0
    }
	if tmp > 0 {
		return 1
	} else {
		return -1
	}
}

func (a Altitude) IsWithin(min Altitude, max Altitude) bool {
	return a.Normalize().Compare(min.Normalize()) > 0 && a.Normalize().Compare(max.Normalize()) < 0
}