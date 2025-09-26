package utils_test

import (
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/utils"
	"testing"
)

func TestExportToGeoJSON(t *testing.T) {
	wp0, _ := models.NewWaypointFromGeojson(`
		{
		"type": "Feature",
		"properties": {},
		"geometry": {
		"coordinates": [
			-3.680087389413245,
			40.41288213307115
		],
		"type": "Point"
		}
	}
	`)

	wp1, _ := models.NewWaypointFromGeojson(`
	{
      "type": "Feature",
      "properties": {
        "altitude": 120
      },
      "geometry": {
        "coordinates": [
          -3.680258352459191,
          40.41451466192274
        ],
        "type": "Point"
      }
    }
	`)

	wp2, _ := models.NewWaypointFromGeojson(`{
		"type": "Feature",
		"properties": {},
		"geometry": {
		"coordinates": [
			-3.684083650608301,
			40.415360741026745
		],
		"type": "Point"
		}
	}`)

	constraint, _ := models.NewConstraintFromGeojson(`
	{
    "type": "Feature",
    "geometry": {
    "type": "Polygon",
    "coordinates": [
        [
        [
            -3.682350868202377,
            40.41280363703126
        ],
        [
            -3.6803385702834817,
            40.41316452298395
        ],
        [
            -3.681329573921829,
            40.41581779687536
        ],
        [
            -3.6828215244555906,
            40.415560765544114
        ],
        [
            -3.682742913597508,
            40.41516672686038
        ],
        [
            -3.6841695232323843,
            40.413956171911565
        ],
        [
            -3.682350868202377,
            40.41280363703126
        ]
        ]
    ]
    },
    "properties": {
    "altitudeUnit": "mt",
    "maxAltitudeValue": 400,
    "minAltitudeValue": 500
    },
    "id": 0
}
	`)
	
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		waypoints []*models.Waypoint
		polygons  []*models.Feature3D
		filename  string
		lineBetweenWaypoints bool
		wantErr   bool
	}{		
		{
			name: "3 wp + 1 obstacle",
			waypoints: []*models.Waypoint{wp0, wp1, wp2},
			polygons: []*models.Feature3D{constraint},
			filename: "test0",
			lineBetweenWaypoints: true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := utils.ExportToGeoJSON(tt.waypoints, tt.polygons, tt.filename, tt.lineBetweenWaypoints)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ExportToGeoJSON() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ExportToGeoJSON() succeeded unexpectedly")
			}
		})
	}
}
