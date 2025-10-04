package models

import (
	"encoding/json"
	"fmt"
	"math"
	"os"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/geojson"
	"github.com/tidwall/geodesic"
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

func MustNewWaypoint(id int, lat, lon float64, alt Altitude) *Waypoint {
    wp, err := NewWaypoint(lat, lon, alt)
    if err != nil {
        panic(err)
    }

	wp.ID = id
    return wp
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

// CirclePolygon generates a polygon feature representing a circle around a waypoint
// radius is in meters, numSides controls how smooth the circle approximation is
func (w *Waypoint) CircleAroundWaypoint(radiusMeters float64, numSides int) *geojson.Feature {
	if numSides < 8 {
		numSides = 64 // ensure good approximation
	}

	centerPoint := w.Point2D()
	radiusDeg := radiusMeters / 111320.0 // rough conversion meters->degrees (works for small areas)

	ring := make(orb.Ring, numSides+1)
	for i := 0; i <= numSides; i++ {
		theta := 2 * math.Pi * float64(i) / float64(numSides)
		dx := radiusDeg * math.Cos(theta)
		dy := radiusDeg * math.Sin(theta)

		ring[i] = orb.Point{
			centerPoint[0] + dx,
			centerPoint[1] + dy,
		}
	}

	poly := orb.Polygon{ring}
	feature := geojson.NewFeature(poly)
	feature.Properties["type"] = "circle"
	feature.Properties["center_lat"] = w.Lat
	feature.Properties["center_lon"] = w.Lon
	feature.Properties["radius_mt"] = radiusMeters
	return feature
}

// CirclePolygonGeodesic generates a geodesic-accurate circle polygon around a waypoint.
// radiusMeters = circle radius, numSides = polygon smoothness
func (w *Waypoint) CircleAroundWaypointGeodesic(radiusMeters float64, numSides int) *geojson.Feature {
	if numSides < 8 {
		numSides = 64
	}

	ring := make(orb.Ring, numSides+1)

	// step around the circle bearings 0..360°
	for i := 0; i <= numSides; i++ {
		azi := 360.0 * float64(i) / float64(numSides) // bearing in degrees

		var lat, lon float64
		geodesic.WGS84.Direct(w.Lat, w.Lon, azi, radiusMeters, &lat, &lon, nil)

		ring[i] = orb.Point{lon, lat}
	}

	poly := orb.Polygon{ring}
	feature := geojson.NewFeature(poly)
	feature.Properties["type"] = "circle"
	feature.Properties["center_lat"] = w.Lat
	feature.Properties["center_lon"] = w.Lon
	feature.Properties["radius_m"] = radiusMeters

	return feature
}

// BoundingBoxAroundWaypoint creates a bounding box polygon feature around a waypoint.
// Uses orb/geo.NewBoundAroundPoint (not a circle).
func (w *Waypoint) BBoxAroundWaypoint(radiusMeters float64) *geojson.Feature {
	// orb expects radius in meters, and Point as [lon, lat]
	pt := orb.Point{w.Lon, w.Lat}

	// get the bounding box
	bound := geo.NewBoundAroundPoint(pt, radiusMeters)

	// convert bound to a polygon (a closed Ring)
	ring := orb.Ring{
		{bound.Min[0], bound.Min[1]}, // bottom-left
		{bound.Max[0], bound.Min[1]}, // bottom-right
		{bound.Max[0], bound.Max[1]}, // top-right
		{bound.Min[0], bound.Max[1]}, // top-left
		{bound.Min[0], bound.Min[1]}, // close polygon
	}
	poly := orb.Polygon{ring}

	// wrap as GeoJSON
	feature := geojson.NewFeature(poly)
	feature.Properties["type"] = "bounding_box"
	feature.Properties["center_lat"] = w.Lat
	feature.Properties["center_lon"] = w.Lon
	feature.Properties["radius_m"] = radiusMeters

	return feature
}

func (w *Waypoint) GetLineString(w2 *Waypoint) orb.LineString {
	return orb.LineString{w.Point2D(), w2.Point2D()}
}

func (w *Waypoint) GetLineStringFeature(w2 *Waypoint) *geojson.Feature {
	feature := geojson.NewFeature(w.GetLineString(w2))
	return feature
}

func (w *Waypoint) GetLineStringBound(w2 *Waypoint) orb.Bound {
	// TODO: Implement in a smarter way, since now if the line is diagonal it will create a very big BB
	return w.GetLineString(w2).Bound()
}

func (w *Waypoint) GetLineStringBoundFeature(w2 *Waypoint) *geojson.Feature {
	return geojson.NewFeature(w.GetLineString(w2).Bound())
}