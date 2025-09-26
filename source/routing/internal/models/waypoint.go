package models

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

type Waypoint struct {
	*geojson.Feature
	Lat float64
	Lon float64
	Alt Altitude
}

// NewWaypoint is a constructor with validation
func NewWaypoint(lat, lon float64, alt Altitude) (*Waypoint, error) {
	wp := &Waypoint{
		Lat: lat,
		Lon: lon,
		Alt: alt,
	}
	
	if err := wp.Validate(); err != nil {
		return nil, err
	}

	wp.Feature = geojson.NewFeature(wp.Point2D())
	
	// No need to check again alt, as it was alread validated before
	wp.SetAltitude(alt)
	return wp, nil
}

func NewWaypointFromGeojsonFile(filename string) (*Waypoint, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	return NewWaypointFromGeojson(string(data))
}

func NewWaypointFromGeojson(geojsonString string) (*Waypoint, error) {
	var w *Waypoint
	if err := json.Unmarshal([]byte(geojsonString), &w); err != nil {
		return nil, fmt.Errorf("unmarshaling geojson: %w", err)
	}

	// fmt.Printf("altitude: %+v\n", w.Alt)
	// fmt.Printf("waypoint: %+v\n", w)
	return w, nil
}

// Validate checks if lat/lon values are in valid range
func (w *Waypoint) Validate() error {
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

func (w *Waypoint) SetAltitude(alt Altitude) error {
	// Make sure to validate altitude
	if err := alt.Validate(); err != nil {
		return err
	}

	// Assign alt to instance variable
	w.Alt = alt

	// Write altitudes in the properties
	w.Feature.Properties["altitudeValue"] = float64(w.Alt.Value)
	w.Feature.Properties["altitudeUnit"] = w.Alt.Unit
	return nil
}

func (w *Waypoint) Point2D() orb.Point {
	return orb.Point{
		w.Lon,
		w.Lat,
	}
}

func (w *Waypoint) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &w.Feature); err != nil {
        return err
    }

	var err error
	unit := w.Feature.Properties.MustString("altitudeUnit", string(MT))
	alt_value := w.Feature.Properties.MustFloat64("altitudeValue", DEFAULT_ALT)
	alt, err := NewAltitude(alt_value, AltitudeUnit(unit))
	if err != nil {
		return err
	}
	
	w.Lat = w.Feature.Point().Lat()
	w.Lon = w.Feature.Point().Lon()
	w.SetAltitude(alt)

	return nil
}

func (w *Waypoint) MarshalJSON() ([]byte, error) {
	return w.Feature.MarshalJSON()
}