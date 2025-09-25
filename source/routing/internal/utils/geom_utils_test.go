package utils

import (
	"geopathplanner/routing/internal/models"
	"testing"
)

func TestLineInPolygon(t *testing.T) {
	wp0json := `{
		"type": "Feature",
		"properties": {},
		"geometry": {
		"coordinates": [
			-3.680087389413245,
			40.41288213307115
		],
		"type": "Point"
		}
	}`
	wp1json := `{
	"type": "Feature",
	"properties": {},
	"geometry": {
		"coordinates": [
		-3.680258352459191,
		40.41451466192274
		],
		"type": "Point"
	}
	}`
	wp2json := `{
		"type": "Feature",
		"properties": {},
		"geometry": {
		"coordinates": [
			-3.684083650608301,
			40.415360741026745
		],
		"type": "Point"
		}
	}`

	wp3json := `{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          -3.683451478749049,
          40.4133217702433
        ],
        "type": "Point"
      }
    }`

	wp4json := `    {
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          -3.6831563623714487,
          40.412370308142016
        ],
        "type": "Point"
      }
    }`

	polyjson := `{
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
		"maxAltitudeValue": 999999,
		"minAltitudeValue": -999999
		},
		"id": 0
	}`

	tests := []struct {
		name     string
		p1       models.Waypoint
		p2       models.Waypoint
		poly     *models.Constraint
		expected bool
	}{
		{
			name: "Line completely outside polygon",
			p1: func() models.Waypoint {
				wp, _ := models.NewWaypointFromGeojson(wp0json)
				return wp
			}(),
			p2: func() models.Waypoint {
				wp, _ := models.NewWaypointFromGeojson(wp1json)
				return wp
			}(),
			poly: func() *models.Constraint {
				c, _ := models.NewConstraintFromGeojson(polyjson)
				return c
			}(),
			expected: false,
		},
		{
			name: "Line inside polygon bbox but outside polygon altitude",
			p1: func() models.Waypoint {
				wp, _ := models.NewWaypointFromGeojson(wp1json)
				return wp
			}(),
			p2: func() models.Waypoint {
				wp, _ := models.NewWaypointFromGeojson(wp2json)
				return wp
			}(),
			poly: func() *models.Constraint {
				c, _ := models.NewConstraintFromGeojson(polyjson)
				c.MinAltitude, _ = models.NewAltitude(400, models.MT)
				c.MaxAltitude, _ = models.NewAltitude(500, models.MT)
				return c
			}(),
			expected: false,
		},
		{
			name: "Line inside polygon bbox but outside polygon shape",
			p1: func() models.Waypoint {
				wp, _ := models.NewWaypointFromGeojson(wp3json)
				return wp
			}(),
			p2: func() models.Waypoint {
				wp, _ := models.NewWaypointFromGeojson(wp4json)
				return wp
			}(),
			poly: func() *models.Constraint {
				c, _ := models.NewConstraintFromGeojson(polyjson)
				return c
			}(),
			expected: false,
		},
		{
			name: "Line inside polygon shape",
			p1: func() models.Waypoint {
				wp, _ := models.NewWaypointFromGeojson(wp1json)
				return wp
			}(),
			p2: func() models.Waypoint {
				wp, _ := models.NewWaypointFromGeojson(wp2json)
				return wp
			}(),
			poly: func() *models.Constraint {
				c, _ := models.NewConstraintFromGeojson(polyjson)
				return c
			}(),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LineInPolygon(tt.p1, tt.p2, tt.poly)
			if result != tt.expected {
				t.Errorf("LineInPolygon() = %v, want %v", result, tt.expected)
			}
		})
	}
}