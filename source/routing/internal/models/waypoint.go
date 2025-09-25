package models

import "fmt"

type Waypoint struct {
	Lat float64  `json:"lat"`
	Lon float64  `json:"lon"`
	Alt Altitude `json:"alt"`
}

// NewWaypoint is a constructor with validation
func NewWaypoint(lat, lon float64, alt Altitude) (*Waypoint, error) {
	wp := Waypoint{
		Lat: lat,
		Lon: lon,
		Alt: alt,
	}
	if err := wp.Validate(); err != nil {
		return nil, err
	}
	return &wp, nil
}

// Validate checks if lat/lon values are in valid range
func (w Waypoint) Validate() error {
	if w.Lat < -90 || w.Lat > 90 {
		return fmt.Errorf("invalid latitude: %.6f (must be between -90 and 90)", w.Lat)
	}
	if w.Lon < -180 || w.Lon > 180 {
		return fmt.Errorf("invalid longitude: %.6f (must be between -180 and 180)", w.Lon)
	}
	if err := w.Alt.Validate(); err != nil {
		return fmt.Errorf("invalid altitude: %w", err)
	}
	return nil
}