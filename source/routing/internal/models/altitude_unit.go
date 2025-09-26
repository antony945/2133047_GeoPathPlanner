package models

import "fmt"

type AltitudeUnit string

const (
	MT       AltitudeUnit = "mt"
	FT       AltitudeUnit = "ft"
	MT_TO_FT float64      = 3.28084
)

// Validate Altitude (unit must be just MT or FT)
func (a AltitudeUnit) Validate() error {
	switch a {
	case MT, FT:
		return nil
	default:
		return fmt.Errorf("invalid altitude unit: %s, available options are %s or %s", a, MT, FT)
	}
}