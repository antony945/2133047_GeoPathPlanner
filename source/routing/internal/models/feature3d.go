package models

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
	"github.com/twpayne/go-geom"
)

type Feature3D struct {
	*geojson.Feature
	MinAltitude Altitude
	MaxAltitude Altitude
}

func NewFeatureFromGeojsonFeature(feature *geojson.Feature) (*Feature3D, error) {
	c := &Feature3D{}
	c.Feature = feature

	var err error

	unit := c.Feature.Properties.MustString("altitudeUnit", string(MT))
	min := c.Feature.Properties.MustFloat64("minAltitudeValue", DEFAULT_MIN_ALT)
	max := c.Feature.Properties.MustFloat64("maxAltitudeValue", DEFAULT_MAX_ALT)
	minAlt, err := NewAltitude(min, AltitudeUnit(unit))
	if err != nil {
		return nil, err
	}
	maxAlt, err := NewAltitude(max, AltitudeUnit(unit))
	if err != nil {
		return nil, err
	}
	if err := c.SetAltitude(minAlt, maxAlt); err != nil {
		return nil, err
	}

	return c, nil
}

func MustNewFeatureFromGeojsonFeature(feature *geojson.Feature) *Feature3D {
	f, err := NewFeatureFromGeojsonFeature(feature)
    if err != nil {
        panic(err)
    }
    return f
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

// ---------------------------------------------------------------
func (c *Feature3D) ToGeomPolygon() *geom.Polygon {
	// Create geom polygon
	// g := geom.NewPolygon(geom.XY)
	// data, err := c.MarshalJSON()
	// if err != nil {
	// 	panic("impossibile marshal json")
	// }
	
	// var g geom.T
	// gj.Unmarshal(data, &g)
	// return g

	// OrbToGeomPolygon converts an orb.Polygon to a gogeos geom.Polygon
	p := c.ToPolygon()
	
	if len(p) == 0 {
		return nil
	}

	// geom expects slices of slices of geom.Coord (rings)
	rings := make([][]geom.Coord, len(p))
	for i, ring := range p {
		coords := make([]geom.Coord, len(ring))
		for j, pt := range ring {
			coords[j] = geom.Coord{pt.X(), pt.Y()}
		}
		rings[i] = coords
	}

	return geom.NewPolygon(geom.XY).MustSetCoords(rings)
	// geom.NewPolygon(geom.XY).MustSetCoords(c.ToPolygon()[0][0])
}

func (c *Feature3D) GetVertices(alt Altitude, reversed bool) []*Waypoint {
	polygon := c.ToPolygon()
	if len(polygon) == 0 {
		return nil
	}

	// Get outer ring vertices (skip last one as it's same as first)
	vertices := make([]*Waypoint, 0, len(polygon[0])-1)
	for _, pt := range polygon[0][:len(polygon[0])-1] {
		wp, _ := NewWaypoint(pt.Lat(), pt.Lon(), alt)
		vertices = append(vertices, wp)
	}

	if reversed {
		left := 0
		right := len(vertices) - 1
		for left < right {
			vertices[left], vertices[right] = vertices[right], vertices[left]
			left++
			right--
	}
	}

	return vertices
}