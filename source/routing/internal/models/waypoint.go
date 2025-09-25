package models

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

type Waypoint struct {
	Lat float64
	Lon float64
	Alt Altitude
}

// NewWaypoint is a constructor with validation
func NewWaypoint(lat, lon float64, alt Altitude) (Waypoint, error) {
	wp := Waypoint{
		Lat: lat,
		Lon: lon,
		Alt: alt,
	}		
	
	if err := wp.Validate(); err != nil {
		return Waypoint{}, err
	}
	return wp, nil
}

func NewWaypointFromGeojsonFile(filename string) (Waypoint, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return Waypoint{}, fmt.Errorf("reading file: %w", err)
	}

	return NewWaypointFromGeojson(string(data))
}

func NewWaypointFromGeojson(geojsonString string) (Waypoint, error) {
	var w *Waypoint
	if err := json.Unmarshal([]byte(geojsonString), &w); err != nil {
		return Waypoint{}, fmt.Errorf("unmarshaling geojson: %w", err)
	}

	// fmt.Printf("altitude: %+v\n", w.Alt)
	// fmt.Printf("waypoint: %+v\n", w)
	return *w, nil
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

func (w Waypoint) Point2D() orb.Point {
	return orb.Point{
		w.Lon,
		w.Lat,
	}
}

func (w Waypoint) Feature() *geojson.Feature {
	f := geojson.NewFeature(w.Point2D())
	f.Properties["altitudeValue"] = float64(w.Alt.Value)
	f.Properties["altitudeUnit"] = string(w.Alt.Unit)
	return f
}

func (w *Waypoint) UnmarshalJSON(data []byte) error {
    // Create a new feature
	var f *geojson.Feature
	
	if err := json.Unmarshal(data, &f); err != nil {
        return err
    }

	var err error
	unit := f.Properties.MustString("altitudeUnit", string(MT))
	alt_value := f.Properties.MustFloat64("altitudeValue", DEFAULT_ALT)
	alt, err := NewAltitude(alt_value, AltitudeUnit(unit))
	if err != nil {
		return err
	}
	
	w.Lat = f.Point().Lat()
	w.Lon = f.Point().Lon()
	w.Alt = alt

	return nil
}

func (w *Waypoint) MarshalJSON() ([]byte, error) {
	return w.Feature().MarshalJSON()
}