package models

import (
	"fmt"
	"math"
	"strings"
)

type AltitudeUnit string

const (
	MT AltitudeUnit = "MT"
	FT AltitudeUnit = "FT"
	MT_TO_FT float64 = 3.28084
)

type Altitude struct {
	Value float64      `json:"value"`
	Unit  AltitudeUnit `json:"unit"`
}

func NewAltitude(value float64, unit AltitudeUnit) (Altitude, error) {
	a := Altitude{
		Value: value,
		Unit: unit,
	}
	if err := a.Validate(); err != nil {
		return Altitude{}, err
	}

	return a, nil
}

// Validate Altitude (unit must be just MT or FT)
func (a Altitude) Validate() error {
	switch a.Unit {
	case MT, FT:
		return nil
	default:
		return fmt.Errorf("invalid altitude unit: %s, must be %s or %s", a.Unit, MT, FT)
	}
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

// Parse unit
func ParseUnit(unit string) (AltitudeUnit, error) {
    switch strings.ToUpper(unit) {
    case "MT":
        return MT, nil
    case "FT":
        return FT, nil
    default:
		return "", fmt.Errorf("invalid altitude unit: %s, must be %s or %s", unit, MT, FT)
    }
}

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