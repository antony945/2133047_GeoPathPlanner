package utils

import (
	"encoding/json"
	"fmt"
	"geopathplanner/routing/internal/models"
	"os"
	"path/filepath"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

type Styler interface {
	Style(f *geojson.Feature) *geojson.Feature
}

type WaypointStyler struct {}
func (s WaypointStyler) Style(f *geojson.Feature) *geojson.Feature {
	color := "#656eec"
	stroke_width := 5
	
	if isInsidePolygon := f.Properties["inside"]; isInsidePolygon != nil {
		flag, _ := isInsidePolygon.(bool)
		if flag {
			color = "#d24141"
		} else {
			color = "#3ca32e"
		}
	}

	if f.Properties["parameter"] != nil {
		color = "#ffffff"
	}
	if f.Properties["near"] != nil {
		color = "#ffce0a"
	}
	if f.Properties["nearest"] != nil {
		color = "#ff8d0a"
	}

	return s.AdvancedStyle(f, color, stroke_width)
}

func (s WaypointStyler) AdvancedStyle(f *geojson.Feature, color string, stroke_width int) *geojson.Feature {
	f.Properties["fill"] = color
	f.Properties["fill-opacity"] = 0.8
	f.Properties["stroke"] = "#555555"
	f.Properties["stroke-width"] = stroke_width
	f.Properties["stroke-opacity"] = 0.8
	// f.Properties["marker-color"] = color
    // f.Properties["marker-size"] = "medium"
    // f.Properties["marker-symbol"] = "circle"
	return f
}

type PolygonStyler struct {}
func (s PolygonStyler) Style(f *geojson.Feature) *geojson.Feature {
	color := "#555555"
	return s.AdvancedStyle(f, color)
}

func (s PolygonStyler) AdvancedStyle(f *geojson.Feature, color string) *geojson.Feature {
	f.Properties["fill"] = color
	f.Properties["fill-opacity"] = 0.3
	f.Properties["stroke"] = "#555555"
	f.Properties["stroke-width"] = 2
	f.Properties["stroke-opacity"] = 0.8
	return f
}

type LineStyler struct {}
func (s LineStyler) Style(f *geojson.Feature) *geojson.Feature {
	f.Properties["stroke"] = "#3388ff"
	f.Properties["stroke-width"] = 5
	f.Properties["stroke-opacity"] = 0.8
	return f
}

// ExportToGeoJSON takes waypoints and constraints and writes them into a FeatureCollection
// saved under dev/results/<folder>/<filename>.geojson.
func ExportToGeoJSON(folder string, waypoints []*models.Waypoint, polygons []*models.Feature3D, filename string, lineBetweenWaypoints bool, otherFeatures ...*geojson.Feature) error {
	fc := CreateFeatureCollection(waypoints, polygons, lineBetweenWaypoints, otherFeatures...)

	// Marshal collection
	data, err := json.MarshalIndent(fc, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal feature collection: %w", err)
	}

	// Ensure dev/results exists
	outDir := ResolvePath(filepath.Join("dev", "results", folder))
	if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
		return fmt.Errorf("create results dir: %w", err)
	}
	
	// Write to file
	outPath := filepath.Join(outDir, filename+".geojson")
	if err := os.WriteFile(outPath, data, 0644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}

func CreateFeatureCollection(waypoints []*models.Waypoint, polygons []*models.Feature3D, lineBetweenWaypoints bool, otherFeatures ...*geojson.Feature) (*geojson.FeatureCollection) {
	// Create a FeatureCollection
	fc := geojson.NewFeatureCollection()
	var styler Styler
	
	// Wrap waypoints (MarshalJSON already defines how they look)
	styler = WaypointStyler{}
	for _, wp := range waypoints {
		fc.Append(styler.Style(wp.Feature))
	}

	// Draw also the line that passes through the point if needed
	styler = LineStyler{}
	if lineBetweenWaypoints {
		line := make(orb.LineString, len(waypoints))
		for i, wp := range waypoints {
			line[i] = wp.Point2D()
		}
		lineFeature := geojson.NewFeature(line)
		fc.Append(styler.Style(lineFeature))
	}

	// Wrap constraints
	styler = PolygonStyler{}
	for _, p := range polygons {
		if p == nil {
			continue
		}
		fc.Append(styler.Style(p.Feature))
	}

	// Append other features if we have some
	styler = PolygonStyler{}
	for _, f := range otherFeatures {
		fc.Append(styler.Style(f))
	}

	return fc
}
