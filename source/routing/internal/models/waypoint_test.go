package models_test

import (
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWaypoint_GetLineStringBoundFeature(t *testing.T) {
	tests := []struct {
		name     string
		wp1      *models.Waypoint
		wp2      *models.Waypoint
		wantType string
	}{
		{
			name:     "Simple line string bound",
			wp1:      models.MustNewWaypoint(1, 41.902782, 12.496366, models.Altitude{Value: 100, Unit: models.MT}), // Rome
			wp2:      models.MustNewWaypoint(2, 48.858370, 2.294481,  models.Altitude{Value: 100, Unit: models.MT}),  // Paris
			wantType: "Polygon",
		},
		{
			name:     "Same point bound",
			wp1:      models.MustNewWaypoint(1, 41.902782, 12.496366, models.Altitude{Value: 100, Unit: models.MT}),
			wp2:      models.MustNewWaypoint(2, 41.902782, 12.496366, models.Altitude{Value: 100, Unit: models.MT}),
			wantType: "Polygon",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.wp1.GetLineStringBoundFeature(tt.wp2)

			// TODO: For visually testing, export results in geojson
			utils.ExportToGeoJSON("models", []*models.Waypoint{tt.wp1, tt.wp2}, []*models.Feature3D{models.MustNewFeatureFromGeojsonFeature(got)}, tt.name, true)

			// Check that we got a feature
			assert.NotNil(t, got)

			// Verify the geometry type is Polygon (bounds are always polygons)
			assert.Equal(t, tt.wantType, got.Geometry.GeoJSONType())

			// // For same points, bound should be very small
			// if tt.wp1.Lat == tt.wp2.Lat && tt.wp1.Lon == tt.wp2.Lon {
			// 	bound := got.Geometry.(orb.Polygon)[0]
			// 	assert.Equal(t, bound[0], bound[4]) // First and last points should be same (closed polygon)
			// }
		})
	}
}