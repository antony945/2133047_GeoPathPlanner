package storage_test

import (
	"fmt"
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/storage"
	"geopathplanner/routing/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStorage_NearestPoint(t *testing.T) {
	a, _ := models.NewAltitude(100, models.MT)
  high, _ := models.NewAltitude(1000, models.MT)

	w := models.MustNewWaypoint(4, 40.41470804648725, -3.7105863826680263, a)

	w2 := models.MustNewWaypoint(5, 40.4196041474128, -3.711548885387799, a)

  w_close2dbuthigh := models.MustNewWaypoint(6, 40.41521155294711, -3.709898376035113, high)
	
	w_list := []*models.Waypoint{
		models.MustNewWaypoint(0, 40.41916768849225, -3.71113552992486, a),
		models.MustNewWaypoint(1, 40.4196554439028, -3.7156507406992887, a),
		models.MustNewWaypoint(2, 40.41656626657965, -3.7131795780455263, a),
		models.MustNewWaypoint(3, 40.42051481387966, -3.7107084153917924, a),
    w_close2dbuthigh,
	}
	
	c_list := []*models.Feature3D{
		models.MustNewFeatureFromGeojson(`
			{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              -3.7125694144268095,
              40.420538039942926
            ],
            [
              -3.7121728080753087,
              40.41844766213757
            ],
            [
              -3.710372825402146,
              40.418354754949746
            ],
            [
              -3.710311809040263,
              40.42090965586385
            ],
            [
              -3.7125694144268095,
              40.420538039942926
            ]
          ]
        ],
        "type": "Polygon"
      }
    }
		`),
		models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              -3.716352428860091,
              40.41830830130769
            ],
            [
              -3.712447381703072,
              40.41798312491591
            ],
            [
              -3.7133016107692356,
              40.41977157562235
            ],
            [
              -3.716352428860091,
              40.41830830130769
            ]
          ]
        ],
        "type": "Polygon"
      }
    }`),
		models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              -3.711715185361953,
              40.4182618476336
            ],
            [
              -3.7134846598548847,
              40.414266711674
            ],
            [
              -3.7105863826680263,
              40.4182618476336
            ],
            [
              -3.711715185361953,
              40.4182618476336
            ]
          ]
        ],
        "type": "Polygon"
      }
    }`),
		models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              -3.7033559437936674,
              40.422628352796266
            ],
            [
              -3.6985661593914756,
              40.41960899115938
            ],
            [
              -3.6971322748894693,
              40.42332510931317
            ],
            [
              -3.700732240235908,
              40.42385928442218
            ],
            [
              -3.7033559437936674,
              40.422628352796266
            ]
          ]
        ],
        "type": "Polygon"
      }
    }`),
		models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              -3.7113490871907118,
              40.42176900981872
            ],
            [
              -3.7114101035526232,
              40.420212874327774
            ],
            [
              -3.7074440400350284,
              40.4197018966147
            ],
            [
              -3.7083287772814515,
              40.42290705626897
            ],
            [
              -3.7113490871907118,
              40.42176900981872
            ]
          ]
        ],
        "type": "Polygon"
      }
    }`),
	}
	
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		w_list []*models.Waypoint
		c_list []*models.Feature3D
		// Named input parameters for target function.
		p       *models.Waypoint
		want    *models.Waypoint
		wantErr bool
	}{
		{ name: "MEMORY - Find NearestPoint to external point", w_list: w_list, c_list: c_list, p: w, want: w_list[2], wantErr: false, },
		{ name: "MEMORY - Find NearestPoint to external point - closer", w_list: w_list, c_list: c_list, p: w2, want: w_list[0], wantErr: false, },
		{ name: "MEMORY - Find NearestPoint to internal point", w_list: w_list, c_list: c_list, p: w_list[0], want: w_list[0], wantErr: false, },
		{ name: "MEMORY - Find NearestPoint to external point in presence of a 2d but not 3d nearest", w_list: w_list, c_list: c_list, p: w, want: w_list[2], wantErr: false, },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := storage.NewMemoryStorage(tt.w_list, tt.c_list)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got, dist, gotErr := m.NearestPoint(tt.p)
      fmt.Printf("nearest dist: %f mt\n", dist)

			// TODO: For visually testing, export results in geojson
			utils.ExportToGeoJSON("storage", append(tt.w_list, tt.p), tt.c_list, tt.name, false)

			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NearestPoint() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NearestPoint() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("NearestPoint() = %v, want %v", got.ID, tt.want.ID)
			}
		})
	}
}

func TestMemoryStorage_KNearestPoints(t *testing.T) {
	a, _ := models.NewAltitude(100, models.MT)
	w := models.MustNewWaypoint(4, 40.41470804648725, -3.7105863826680263, a)

	w2 := models.MustNewWaypoint(5, 40.4196041474128, -3.711548885387799, a)
	
	w_list := []*models.Waypoint{
		models.MustNewWaypoint(0, 40.41916768849225, -3.71113552992486, a),
		models.MustNewWaypoint(1, 40.4196554439028, -3.7156507406992887, a),
		models.MustNewWaypoint(2, 40.41656626657965, -3.7131795780455263, a),
		models.MustNewWaypoint(3, 40.42051481387966, -3.7107084153917924, a),
	}
	
	c_list := []*models.Feature3D{
		models.MustNewFeatureFromGeojson(`
			{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              -3.7125694144268095,
              40.420538039942926
            ],
            [
              -3.7121728080753087,
              40.41844766213757
            ],
            [
              -3.710372825402146,
              40.418354754949746
            ],
            [
              -3.710311809040263,
              40.42090965586385
            ],
            [
              -3.7125694144268095,
              40.420538039942926
            ]
          ]
        ],
        "type": "Polygon"
      }
    }
		`),
		models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              -3.716352428860091,
              40.41830830130769
            ],
            [
              -3.712447381703072,
              40.41798312491591
            ],
            [
              -3.7133016107692356,
              40.41977157562235
            ],
            [
              -3.716352428860091,
              40.41830830130769
            ]
          ]
        ],
        "type": "Polygon"
      }
    }`),
		models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              -3.711715185361953,
              40.4182618476336
            ],
            [
              -3.7134846598548847,
              40.414266711674
            ],
            [
              -3.7105863826680263,
              40.4182618476336
            ],
            [
              -3.711715185361953,
              40.4182618476336
            ]
          ]
        ],
        "type": "Polygon"
      }
    }`),
		models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              -3.7033559437936674,
              40.422628352796266
            ],
            [
              -3.6985661593914756,
              40.41960899115938
            ],
            [
              -3.6971322748894693,
              40.42332510931317
            ],
            [
              -3.700732240235908,
              40.42385928442218
            ],
            [
              -3.7033559437936674,
              40.422628352796266
            ]
          ]
        ],
        "type": "Polygon"
      }
    }`),
		models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              -3.7113490871907118,
              40.42176900981872
            ],
            [
              -3.7114101035526232,
              40.420212874327774
            ],
            [
              -3.7074440400350284,
              40.4197018966147
            ],
            [
              -3.7083287772814515,
              40.42290705626897
            ],
            [
              -3.7113490871907118,
              40.42176900981872
            ]
          ]
        ],
        "type": "Polygon"
      }
    }`),
	}
	
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		w_list []*models.Waypoint
		c_list []*models.Feature3D
		// Named input parameters for target function.
		p       *models.Waypoint
		k		int
		want    []*models.Waypoint
		wantErr bool
	}{
		{ name: "MEMORY - KNN - external point - k=0", w_list: w_list, c_list: c_list, p: w, k: 0, want: []*models.Waypoint{}, wantErr: false },
		{ name: "MEMORY - KNN - external point - k=1", w_list: w_list, c_list: c_list, p: w2, k: 1, want: []*models.Waypoint{w_list[0]}, wantErr: false },
		{ name: "MEMORY - KNN - external point - k=3", w_list: w_list, c_list: c_list, p: w2, k: 3, want: []*models.Waypoint{w_list[0], w_list[3], w_list[1]}, wantErr: false },
		{ name: "MEMORY - KNN - external point - k>#wps", w_list: w_list, c_list: c_list, p: w2, k: 7, want: []*models.Waypoint{w_list[0], w_list[3], w_list[1], w_list[2]}, wantErr: false },
		{ name: "MEMORY - KNN - internal point - k=1", w_list: w_list, c_list: c_list, p: w_list[0], k: 1, want: []*models.Waypoint{w_list[0]}, wantErr: false },
		{ name: "MEMORY - KNN - internal point - k=2", w_list: w_list, c_list: c_list, p: w_list[0], k: 2, want: []*models.Waypoint{w_list[0], w_list[3]}, wantErr: false },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := storage.NewMemoryStorage(tt.w_list, tt.c_list)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got, _, gotErr := m.KNearestPoints(tt.p, tt.k)

			// TODO: For visually testing, export results in geojson
			utils.ExportToGeoJSON("storage", append(tt.w_list, tt.p), tt.c_list, tt.name, false)

			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("KNearestPoints() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("KNearestPoints() succeeded unexpectedly")
			}

      // TODO: update the condition below to compare got with tt.want.
      if !assert.ElementsMatch(t, got, tt.want) {
				t.Errorf("KNearestPoints() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryStorage_NearestPointsInRadius(t *testing.T) {
	a, _ := models.NewAltitude(100, models.MT)
	w := models.MustNewWaypoint(4, 40.41470804648725, -3.7105863826680263, a) //far
	w2 := models.MustNewWaypoint(5, 40.4196041474128, -3.711548885387799, a) //closer
	
	w_list := []*models.Waypoint{
		models.MustNewWaypoint(0, 40.41916768849225, -3.71113552992486, a),
		models.MustNewWaypoint(1, 40.4196554439028, -3.7156507406992887, a),
		models.MustNewWaypoint(2, 40.41656626657965, -3.7131795780455263, a),
		models.MustNewWaypoint(3, 40.42051481387966, -3.7107084153917924, a),
	}
	
	c_list := []*models.Feature3D{
		models.MustNewFeatureFromGeojson(`
			{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              -3.7125694144268095,
              40.420538039942926
            ],
            [
              -3.7121728080753087,
              40.41844766213757
            ],
            [
              -3.710372825402146,
              40.418354754949746
            ],
            [
              -3.710311809040263,
              40.42090965586385
            ],
            [
              -3.7125694144268095,
              40.420538039942926
            ]
          ]
        ],
        "type": "Polygon"
      }
    }
		`),
		models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              -3.716352428860091,
              40.41830830130769
            ],
            [
              -3.712447381703072,
              40.41798312491591
            ],
            [
              -3.7133016107692356,
              40.41977157562235
            ],
            [
              -3.716352428860091,
              40.41830830130769
            ]
          ]
        ],
        "type": "Polygon"
      }
    }`),
		models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              -3.711715185361953,
              40.4182618476336
            ],
            [
              -3.7134846598548847,
              40.414266711674
            ],
            [
              -3.7105863826680263,
              40.4182618476336
            ],
            [
              -3.711715185361953,
              40.4182618476336
            ]
          ]
        ],
        "type": "Polygon"
      }
    }`),
		models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              -3.7033559437936674,
              40.422628352796266
            ],
            [
              -3.6985661593914756,
              40.41960899115938
            ],
            [
              -3.6971322748894693,
              40.42332510931317
            ],
            [
              -3.700732240235908,
              40.42385928442218
            ],
            [
              -3.7033559437936674,
              40.422628352796266
            ]
          ]
        ],
        "type": "Polygon"
      }
    }`),
		models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              -3.7113490871907118,
              40.42176900981872
            ],
            [
              -3.7114101035526232,
              40.420212874327774
            ],
            [
              -3.7074440400350284,
              40.4197018966147
            ],
            [
              -3.7083287772814515,
              40.42290705626897
            ],
            [
              -3.7113490871907118,
              40.42176900981872
            ]
          ]
        ],
        "type": "Polygon"
      }
    }`),
	}
	
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		w_list []*models.Waypoint
		c_list []*models.Feature3D
		// Named input parameters for target function.
		p       *models.Waypoint
		r float64
		want    *[]models.Waypoint
		wantErr bool
	}{
		{ name: "MEMORY - RadiusNN - external point - r=0", w_list: w_list, c_list: c_list, p: w2, r: 0, wantErr: false },
		{ name: "MEMORY - RadiusNN - external point - r<0", w_list: w_list, c_list: c_list, p: w2, r: -10.5, wantErr: true },
		{ name: "MEMORY - RadiusNN - external point - r=8.5", w_list: w_list, c_list: c_list, p: w2, r: 8.5, wantErr: false },
		{ name: "MEMORY - RadiusNN - external point - r=150", w_list: w_list, c_list: c_list, p: w2, r: 150, wantErr: false },
		{ name: "MEMORY - RadiusNN - external point - r=500", w_list: w_list, c_list: c_list, p: w2, r: 500, wantErr: false },
		{ name: "MEMORY - RadiusNN - external point2 - r=500", w_list: w_list, c_list: c_list, p: w, r: 500, wantErr: false },
		{ name: "MEMORY - RadiusNN - internal point - r=0.05", w_list: w_list, c_list: c_list, p: w_list[0], r: 0.05, wantErr: false },
		{ name: "MEMORY - RadiusNN - internal point - r=500", w_list: w_list, c_list: c_list, p: w_list[0], r: 500, wantErr: false },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := storage.NewMemoryStorage(tt.w_list, tt.c_list)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			_, _, gotErr := m.NearestPointsInRadius(tt.p, tt.r)

			// TODO: For visually testing, export results in geojson
			// Generate circle polygon

			utils.ExportToGeoJSON("storage", append(tt.w_list, tt.p), tt.c_list, tt.name, false, tt.p.CircleAroundWaypointGeodesic(tt.r).Feature)

			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("NearestPointsInRadius() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("NearestPointsInRadius() succeeded unexpectedly")
			}
		})
	}
}

func TestMemoryStorage_IsPointInObstacles(t *testing.T) {
	a, _ := models.NewAltitude(100, models.MT)
	w := models.MustNewWaypoint(4, 40.41470804648725, -3.7105863826680263, a) //far
	w2 := models.MustNewWaypoint(5, 40.4196041474128, -3.711548885387799, a) //closer
	
	w_list := []*models.Waypoint{
		models.MustNewWaypoint(0, 40.41916768849225, -3.71113552992486, a),
		models.MustNewWaypoint(1, 40.4196554439028, -3.7156507406992887, a),
		models.MustNewWaypoint(2, 40.41656626657965, -3.7131795780455263, a),
		models.MustNewWaypoint(3, 40.42051481387966, -3.7107084153917924, a),
	}
	
	c_list := []*models.Feature3D{
		models.MustNewFeatureFromGeojson(`
			{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              -3.7125694144268095,
              40.420538039942926
            ],
            [
              -3.7121728080753087,
              40.41844766213757
            ],
            [
              -3.710372825402146,
              40.418354754949746
            ],
            [
              -3.710311809040263,
              40.42090965586385
            ],
            [
              -3.7125694144268095,
              40.420538039942926
            ]
          ]
        ],
        "type": "Polygon"
      }
    }
		`),
		models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              -3.716352428860091,
              40.41830830130769
            ],
            [
              -3.712447381703072,
              40.41798312491591
            ],
            [
              -3.7133016107692356,
              40.41977157562235
            ],
            [
              -3.716352428860091,
              40.41830830130769
            ]
          ]
        ],
        "type": "Polygon"
      }
    }`),
		models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              -3.711715185361953,
              40.4182618476336
            ],
            [
              -3.7134846598548847,
              40.414266711674
            ],
            [
              -3.7105863826680263,
              40.4182618476336
            ],
            [
              -3.711715185361953,
              40.4182618476336
            ]
          ]
        ],
        "type": "Polygon"
      }
    }`),
		models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              -3.7033559437936674,
              40.422628352796266
            ],
            [
              -3.6985661593914756,
              40.41960899115938
            ],
            [
              -3.6971322748894693,
              40.42332510931317
            ],
            [
              -3.700732240235908,
              40.42385928442218
            ],
            [
              -3.7033559437936674,
              40.422628352796266
            ]
          ]
        ],
        "type": "Polygon"
      }
    }`),
		models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              -3.7113490871907118,
              40.42176900981872
            ],
            [
              -3.7114101035526232,
              40.420212874327774
            ],
            [
              -3.7074440400350284,
              40.4197018966147
            ],
            [
              -3.7083287772814515,
              40.42290705626897
            ],
            [
              -3.7113490871907118,
              40.42176900981872
            ]
          ]
        ],
        "type": "Polygon"
      }
    }`),
	}
	
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		w_list []*models.Waypoint
		c_list []*models.Feature3D
		// Named input parameters for target function.
		p       *models.Waypoint
		want    bool
		wantErr bool
	}{
		{ name: "MEMORY - PiP - external point0", w_list: w_list, c_list: c_list, p: w, wantErr: false },
		{ name: "MEMORY - PiP - external point1", w_list: w_list, c_list: c_list, p: w2, wantErr: false },
		{ name: "MEMORY - PiP - external point2", w_list: w_list, c_list: c_list, p: w_list[0], wantErr: false },
		{ name: "MEMORY - PiP - external point3", w_list: w_list, c_list: c_list, p: w_list[1], wantErr: false },
		{ name: "MEMORY - PiP - external point4", w_list: w_list, c_list: c_list, p: w_list[2], wantErr: false },
		{ name: "MEMORY - PiP - external point5", w_list: w_list, c_list: c_list, p: w_list[3], wantErr: false },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := storage.NewMemoryStorage(tt.w_list, tt.c_list)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			_, _, gotErr := m.IsPointInObstacles(tt.p)

			// TODO: For visually testing, export results in geojson
			utils.ExportToGeoJSON("storage", append(tt.w_list, tt.p), tt.c_list, tt.name, false)

			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("IsPointInObstacles() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("IsPointInObstacles() succeeded unexpectedly")
			}
		})
	}
}

func TestMemoryStorage_IsLineInObstacles(t *testing.T) {
	a, _ := models.NewAltitude(100, models.MT)
	w := models.MustNewWaypoint(4, 40.41470804648725, -3.7105863826680263, a) //far
	w2 := models.MustNewWaypoint(5, 40.4196041474128, -3.711548885387799, a) //closer
  w3 := models.MustNewWaypoint(6, 40.41961575535774, -3.709500083865919, a)
	
	w_list := []*models.Waypoint{
		models.MustNewWaypoint(0, 40.41916768849225, -3.71113552992486, a),
		models.MustNewWaypoint(1, 40.4196554439028, -3.7156507406992887, a),
		models.MustNewWaypoint(2, 40.41656626657965, -3.7131795780455263, a),
		models.MustNewWaypoint(3, 40.42051481387966, -3.7107084153917924, a),
	}
	
	c_list := []*models.Feature3D{
		models.MustNewFeatureFromGeojson(`
			{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              -3.7125694144268095,
              40.420538039942926
            ],
            [
              -3.7121728080753087,
              40.41844766213757
            ],
            [
              -3.710372825402146,
              40.418354754949746
            ],
            [
              -3.710311809040263,
              40.42090965586385
            ],
            [
              -3.7125694144268095,
              40.420538039942926
            ]
          ]
        ],
        "type": "Polygon"
      }
    }
		`),
		models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              -3.716352428860091,
              40.41830830130769
            ],
            [
              -3.712447381703072,
              40.41798312491591
            ],
            [
              -3.7133016107692356,
              40.41977157562235
            ],
            [
              -3.716352428860091,
              40.41830830130769
            ]
          ]
        ],
        "type": "Polygon"
      }
    }`),
		models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              -3.711715185361953,
              40.4182618476336
            ],
            [
              -3.7134846598548847,
              40.414266711674
            ],
            [
              -3.7105863826680263,
              40.4182618476336
            ],
            [
              -3.711715185361953,
              40.4182618476336
            ]
          ]
        ],
        "type": "Polygon"
      }
    }`),
		models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              -3.7033559437936674,
              40.422628352796266
            ],
            [
              -3.6985661593914756,
              40.41960899115938
            ],
            [
              -3.6971322748894693,
              40.42332510931317
            ],
            [
              -3.700732240235908,
              40.42385928442218
            ],
            [
              -3.7033559437936674,
              40.422628352796266
            ]
          ]
        ],
        "type": "Polygon"
      }
    }`),
		models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              -3.7113490871907118,
              40.42176900981872
            ],
            [
              -3.7114101035526232,
              40.420212874327774
            ],
            [
              -3.7074440400350284,
              40.4197018966147
            ],
            [
              -3.7083287772814515,
              40.42290705626897
            ],
            [
              -3.7113490871907118,
              40.42176900981872
            ]
          ]
        ],
        "type": "Polygon"
      }
    }`),
	}
	
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		w_list []*models.Waypoint
		c_list []*models.Feature3D
		// Named input parameters for target function.
		p       *models.Waypoint
		p2       *models.Waypoint
		want    bool
		wantErr bool
	}{
		{ name: "MEMORY - LiP - samePoint", w_list: w_list, c_list: c_list, p: w, p2: w, want: false, wantErr: false },
		{ name: "MEMORY - LiP - w-w2", w_list: w_list, c_list: c_list, p: w, p2: w2, want: true, wantErr: false },
		{ name: "MEMORY - LiP - 0-1", w_list: w_list, c_list: c_list, p: w_list[0], p2: w_list[1], want: true, wantErr: false },
		{ name: "MEMORY - LiP - 1-2", w_list: w_list, c_list: c_list, p: w_list[1], p2: w_list[2], want: true, wantErr: false },
		{ name: "MEMORY - LiP - 2-3", w_list: w_list, c_list: c_list, p: w_list[2], p2: w_list[3], want: true, wantErr: false },
		{ name: "MEMORY - LiP - 3-w", w_list: w_list, c_list: c_list, p: w_list[3], p2: w, want: true, wantErr: false },
		{ name: "MEMORY - LiP - w-w3", w_list: w_list, c_list: c_list, p: w, p2: w3, want: false, wantErr: false },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := storage.NewMemoryStorage(tt.w_list, tt.c_list)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got, line, gotErr := m.IsLineInObstacles(tt.p, tt.p2)

			// TODO: For visually testing, export results in geojson
			utils.ExportToGeoJSON("storage", line, tt.c_list, tt.name, true)

			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("IsLineInObstacles() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("IsLineInObstacles() succeeded unexpectedly")
			}
      // TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("IsLineInObstacles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryStorage_SampleFree(t *testing.T) {
  a, _ := models.NewAltitude(100, models.MT)
  w1 := models.MustNewWaypoint(0, 50.872778105839274, 4.433724687935722, a) // p1
	w2 := models.MustNewWaypoint(1, 50.884400404439646, 4.46992531620532, a) // p2

	c_list := []*models.Feature3D{
		models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              4.435823054794525,
              50.87917754349243
            ],
            [
              4.435999901698551,
              50.876186530052024
            ],
            [
              4.443605337154025,
              50.878195458959425
            ],
            [
              4.439678842862065,
              50.88446720530686
            ],
            [
              4.435823054794525,
              50.87917754349243
            ]
          ]
        ],
        "type": "Polygon"
      }
    }`),
		models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              4.443180802952469,
              50.87710169486769
            ],
            [
              4.445940351819587,
              50.874043594523414
            ],
            [
              4.449867003560769,
              50.875472226326934
            ],
            [
              4.447425915923816,
              50.88058383172759
            ],
            [
              4.443180802952469,
              50.87710169486769
            ]
          ]
        ],
        "type": "Polygon"
      }
    }`),
		models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              4.454253542517819,
              50.88096314072149
            ],
            [
              4.456623787928493,
              50.87712402157996
            ],
            [
              4.462737042916956,
              50.87863279943002
            ],
            [
              4.460756068008692,
              50.884413843656176
            ],
            [
              4.454253542517819,
              50.88096314072149
            ]
          ]
        ],
        "type": "Polygon"
      }
    }`),
		models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              4.465548569604977,
              50.88086437277474
            ],
            [
              4.466185163129978,
              50.8824937224808
            ],
            [
              4.469156850257434,
              50.88035101122094
            ],
            [
              4.465619026164603,
              50.88738133071263
            ],
            [
              4.460914241378333,
              50.88548451715948
            ],
            [
              4.465548569604977,
              50.88086437277474
            ]
          ]
        ],
        "type": "Polygon"
      }
    }`),
		models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              4.467564748066309,
              50.88892119085034
            ],
            [
              4.463628522001983,
              50.88991457616726
            ],
            [
              4.46840349799416,
              50.887481856932595
            ],
            [
              4.467564748066309,
              50.88892119085034
            ]
          ]
        ],
        "type": "Polygon"
      }
    }`),
	}

  search_volume := models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              4.4281679963691545,
              50.89062392756719
            ],
            [
              4.4281679963691545,
              50.86983921138935
            ],
            [
              4.473729291371427,
              50.86983921138935
            ],
            [
              4.473729291371427,
              50.89062392756719
            ],
            [
              4.4281679963691545,
              50.89062392756719
            ]
          ]
        ],
        "type": "Polygon"
      }
    }`)
	
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		w_list []*models.Waypoint
		c_list []*models.Feature3D
		// Named input parameters for target function.
		sampler utils.Sampler
    sampleVolume *models.Feature3D
    alt models.Altitude
	}{
		{ name: "MEMORY - SampleFree with obstacles - uniform", w_list: []*models.Waypoint{w1, w2}, c_list: c_list, sampler: utils.NewUniformSamplerWithSeed(10), sampleVolume: search_volume, alt: a},
    { name: "MEMORY - SampleFree with obstacles - halton", w_list: []*models.Waypoint{w1, w2}, c_list: c_list, sampler: utils.NewHaltonSampler(), sampleVolume: search_volume, alt: a},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := storage.NewMemoryStorage(tt.w_list, tt.c_list)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}

      // Run 50 sample free
      gotList := make([]*models.Waypoint, 0, 200)
      for i := 0; i < 200; i++ {
        got, _ := m.SampleFree(tt.sampler, tt.sampleVolume, tt.alt)
        gotList = append(gotList, got)
      }
      
			// TODO: For visually testing, export results in geojson
      utils.ExportToGeoJSON("storage", gotList, append(tt.c_list, tt.sampleVolume), tt.name, false)
		})
	}
}