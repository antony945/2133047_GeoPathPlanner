package models

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

type Feature3D struct {
	*geojson.Feature
	MinAltitude Altitude
	MaxAltitude Altitude
}

func NewFeatureFromGeojsonFile(filename string) (*Feature3D, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	return NewFeatureFromGeojson(string(data))
}

func NewFeatureFromGeojson(geojsonString string) (*Feature3D, error) {
	var c *Feature3D
	if err := json.Unmarshal([]byte(geojsonString), &c); err != nil {
		return nil, fmt.Errorf("unmarshaling geojson: %w", err)
	}

	// fmt.Printf("min altitude: %+v\n", c.MinAltitude)
	// fmt.Printf("max altitude: %+v\n", c.MaxAltitude)
	// fmt.Printf("constraint: %+v\n", c)
	return c, nil
}

func MustNewFeatureFromGeojson(geojsonString string) *Feature3D {
	f, err := NewFeatureFromGeojson(geojsonString)
    if err != nil {
        panic(err)
    }
    return f
}

func (c *Feature3D) Bound() orb.Bound {
	return c.Geometry.Bound()
}

func (c *Feature3D) ToPolygon() orb.Polygon {
	return c.Geometry.(orb.Polygon)
}

func (c *Feature3D) SetAltitude(min, max Altitude) error {
	// Make sure that they are both in the same unit
	if err := min.Unit.IsEqual(max.Unit); err != nil {
		return err
	}

	// Check that min is less than max (otherwise swap the two)
	// TODO: Check if swapping is good or if it's better to throw error
	if min.Compare(max) > 0 {
		min, max = max, min
	}

	// Assign min max to instance variable
	c.MinAltitude = min
	c.MaxAltitude = max

	// Write altitudes in the properties
	c.Feature.Properties["minAltitudeValue"] = float64(c.MinAltitude.Value)
	c.Feature.Properties["maxAltitudeValue"] = float64(c.MaxAltitude.Value)
	c.Feature.Properties["altitudeUnit"] = c.MinAltitude.Unit
	return nil
}

func (c *Feature3D) UnmarshalJSON(data []byte) error {
    if err := json.Unmarshal(data, &c.Feature); err != nil {
        return err
    }
	
	var err error

	unit := c.Feature.Properties.MustString("altitudeUnit", string(MT))
	min := c.Feature.Properties.MustFloat64("minAltitudeValue", DEFAULT_MIN_ALT)
	max := c.Feature.Properties.MustFloat64("maxAltitudeValue", DEFAULT_MAX_ALT)
	minAlt, err := NewAltitude(min, AltitudeUnit(unit))
	if err != nil {
		return err
	}
	maxAlt, err := NewAltitude(max, AltitudeUnit(unit))
	if err != nil {
		return err
	}
	if err := c.SetAltitude(minAlt, maxAlt); err != nil {
		return err
	}

	return nil
}

func (c *Feature3D) MarshalJSON() ([]byte, error) {
	return c.Feature.MarshalJSON()
}

