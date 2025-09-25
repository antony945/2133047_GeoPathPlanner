package models

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/paulmach/orb/geojson"
)

type Constraint struct {
	*geojson.Feature
	MinAltitude Altitude
	MaxAltitude Altitude
}

func NewConstraintFromGeojsonFile(filename string) (*Constraint, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	return NewConstraintFromGeojson(string(data))
}

func NewConstraintFromGeojson(geojsonString string) (*Constraint, error) {
	var c *Constraint
	if err := json.Unmarshal([]byte(geojsonString), &c); err != nil {
		return nil, fmt.Errorf("unmarshaling geojson: %w", err)
	}

	// fmt.Printf("min altitude: %+v\n", c.MinAltitude)
	// fmt.Printf("max altitude: %+v\n", c.MaxAltitude)
	// fmt.Printf("constraint: %+v\n", c)
	return c, nil
}

func (c *Constraint) UnmarshalJSON(data []byte) error {
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

func (c *Constraint) MarshalJSON() ([]byte, error) {
	c.Feature.Properties["minAltitudeValue"] = float64(c.MinAltitude.ConvertTo(MT).Value)
	c.Feature.Properties["maxAltitudeValue"] = float64(c.MaxAltitude.ConvertTo(MT).Value)
	c.Feature.Properties["altitudeUnit"] = string(MT)
	return c.Feature.MarshalJSON()
}

