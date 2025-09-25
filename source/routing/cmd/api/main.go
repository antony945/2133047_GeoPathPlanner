package main

import (
	"encoding/json"
	"fmt"
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/utils"
	"os"

	"github.com/paulmach/orb/geojson"
)

func main() {
	fmt.Println("ROUTING MICROSERVICE")
}

func debugConstraintParsing(input_file, output_file string) error {
	f, err := models.NewConstraintFromGeojsonFile(input_file)
	if err != nil {
		return fmt.Errorf("unmarshaling constraint: %w", err)
	}

	fc := geojson.NewFeatureCollection()
	fc.Append(f.Feature)
	for _, feature := range fc.Features {
		fmt.Printf("Feature properties: %+v\n", feature.Properties)
		fmt.Printf("Feature geometry type: %s\n", feature.Geometry.GeoJSONType())
		fmt.Printf("BBox: %+v\n", feature.Geometry.Bound())
	}

	// Marshal f again
	marshaled, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling constraint: %w", err)
	}
	
	if err := os.WriteFile(output_file, marshaled, 0644); err != nil {
		return fmt.Errorf("writing output file: %w", err)
	}
	fmt.Printf("Saved marshaled constraint to: %s\n", output_file)

	return nil
}

func debugWaypointParsing(input_file, output_file string) error {
	f, err := models.NewWaypointFromGeojsonFile(input_file)
	if err != nil {
		return fmt.Errorf("unmarshaling waypoint: %w", err)
	}

	// Marshal f again
	marshaled, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling waypoint: %w", err)
	}
	
	if err := os.WriteFile(output_file, marshaled, 0644); err != nil {
		return fmt.Errorf("writing output file: %w", err)
	}
	fmt.Printf("Saved marshaled waypoint to: %s\n", output_file)

	return nil
}

func testDistances(p1, p2 models.Waypoint) error {
	// Use geo distance first
	fmt.Printf("GeoDistance: %f\n", utils.FastDistance3D(p1, p2))
	fmt.Printf("GeoDistanceHaversine: %f\n\n", utils.HaversineDistance3D(p1, p2))

	fmt.Printf("p1: %+v\n", p1)
	fmt.Printf("p2: %+v\n\n", p2)

	// Create feature collection to display line in geojson
	fc := geojson.NewFeatureCollection()
	// Try to sample line
	sampledwps := utils.DefaultResampleLineToInterval(p1, p2)
	fmt.Printf("Divided line into %d points\n", len(sampledwps))
	for i, wp := range sampledwps {
		fmt.Printf("wp[%d] = %+v\n", i, wp)
		fc.Append(wp.Feature())
	}

	// Marshal feature collection
	marshaled, err := json.MarshalIndent(fc, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling feature collection: %w", err)
	}

	if err := os.WriteFile(utils.ResolvePath("dev/requests/sampled_line.geojson"), marshaled, 0644); err != nil {
		return fmt.Errorf("writing output file: %w", err)
	}
	fmt.Printf("Saved feature collection to: %s\n", utils.ResolvePath("dev/outputs/sampled_line.geojson"))

	return nil
}