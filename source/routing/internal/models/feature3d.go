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

func NewConstraintFromGeojsonFile(filename string) (*Feature3D, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	return NewConstraintFromGeojson(string(data))
}

func NewConstraintFromGeojson(geojsonString string) (*Feature3D, error) {
	var c *Feature3D
	if err := json.Unmarshal([]byte(geojsonString), &c); err != nil {
		return nil, fmt.Errorf("unmarshaling geojson: %w", err)
	}

	// fmt.Printf("min altitude: %+v\n", c.MinAltitude)
	// fmt.Printf("max altitude: %+v\n", c.MaxAltitude)
	// fmt.Printf("constraint: %+v\n", c)
	return c, nil
}

func (c *Feature3D) Bound() orb.Bound {
	return c.Geometry.Bound()
}

func (c *Feature3D) ToPolygon() orb.Polygon {
	return c.Geometry.(orb.Polygon)
}

func (c *Feature3D) UnmarshalJSON(data []byte) error {
    if err := json.Unmarshal(data, &c.Feature); err != nil {
        return err
    }
	
	var err error

	unit := c.Feature.Properties.MustString("altitudeUnit", string(MT))
	
	min := c.Feature.Properties.MustFloat64("minAltitudeValue", DEFAULT_MIN_ALT)
	c.MinAltitude, err = NewAltitude(min, AltitudeUnit(unit))
	if err != nil {
		return err
	}
	
	max := c.Feature.Properties.MustFloat64("maxAltitudeValue", DEFAULT_MAX_ALT)
	c.MaxAltitude, err = NewAltitude(max, AltitudeUnit(unit))
	if err != nil {
		return err
	}

	return nil
}

func (c *Feature3D) MarshalJSON() ([]byte, error) {
	c.Feature.Properties["minAltitudeValue"] = float64(c.MinAltitude.ConvertTo(MT).Value)
	c.Feature.Properties["maxAltitudeValue"] = float64(c.MaxAltitude.ConvertTo(MT).Value)
	c.Feature.Properties["altitudeUnit"] = string(MT)
	return c.Feature.MarshalJSON()
}

