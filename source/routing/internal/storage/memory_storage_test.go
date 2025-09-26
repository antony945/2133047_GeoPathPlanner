package storage_test

import (
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/storage"
	"geopathplanner/routing/internal/utils"
	"testing"
)

func TestMemoryStorage_NearestPoint(t *testing.T) {
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
		want    *models.Waypoint
		wantErr bool
	}{
		{
			name: "Find NearestPoint to external point",
			w_list: w_list,
			c_list: c_list,
			p: w,
			want: w_list[2],
			wantErr: false,
		},
		{
			name: "Find NearestPoint to external point - closer",
			w_list: w_list,
			c_list: c_list,
			p: w2,
			want: w_list[0],
			wantErr: false,
		},
		{
			name: "Find NearestPoint to internal point",
			w_list: w_list,
			c_list: c_list,
			p: w_list[0],
			want: w_list[0],
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := storage.NewMemoryStorage(tt.w_list, tt.c_list)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got, gotErr := m.NearestPoint(tt.p)

			// TODO: For visually testing, export results in geojson
			utils.ExportToGeoJSON(append(tt.w_list, tt.p), tt.c_list, tt.name, false)

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
		{ name: "KNN - external point - k=0", w_list: w_list, c_list: c_list, p: w, k: 0, wantErr: true },
		{ name: "KNN - external point - k=1", w_list: w_list, c_list: c_list, p: w2, k: 1, wantErr: false },
		{ name: "KNN - external point - k=3", w_list: w_list, c_list: c_list, p: w2, k: 3, wantErr: false },
		{ name: "KNN - external point - k>#wps", w_list: w_list, c_list: c_list, p: w2, k: 7, wantErr: false },
		{ name: "KNN - internal point - k=1", w_list: w_list, c_list: c_list, p: w_list[0], k: 1, wantErr: false },
		{ name: "KNN - internal point - k=2", w_list: w_list, c_list: c_list, p: w_list[0], k: 2, wantErr: false },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := storage.NewMemoryStorage(tt.w_list, tt.c_list)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			_, gotErr := m.KNearestPoints(tt.p, tt.k)

			// TODO: For visually testing, export results in geojson
			utils.ExportToGeoJSON(append(tt.w_list, tt.p), tt.c_list, tt.name, false)

			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("KNearestPoints() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("KNearestPoints() succeeded unexpectedly")
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
		{ name: "RadiusNN - external point - r=0", w_list: w_list, c_list: c_list, p: w2, r: 0, wantErr: true },
		{ name: "RadiusNN - external point - r<0", w_list: w_list, c_list: c_list, p: w2, r: -10.5, wantErr: true },
		{ name: "RadiusNN - external point - r=8.5", w_list: w_list, c_list: c_list, p: w2, r: 8.5, wantErr: false },
		{ name: "RadiusNN - external point - r=150", w_list: w_list, c_list: c_list, p: w2, r: 150, wantErr: false },
		{ name: "RadiusNN - external point - r=500", w_list: w_list, c_list: c_list, p: w2, r: 500, wantErr: false },
		{ name: "RadiusNN - external point2 - r=500", w_list: w_list, c_list: c_list, p: w, r: 500, wantErr: false },
		{ name: "RadiusNN - internal point - r=0.05", w_list: w_list, c_list: c_list, p: w_list[0], r: 0.05, wantErr: false },
		{ name: "RadiusNN - internal point - r=500", w_list: w_list, c_list: c_list, p: w_list[0], r: 500, wantErr: false },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := storage.NewMemoryStorage(tt.w_list, tt.c_list)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			_, gotErr := m.NearestPointsInRadius(tt.p, tt.r)

			// TODO: For visually testing, export results in geojson
			// Generate circle polygon

			utils.ExportToGeoJSON(append(tt.w_list, tt.p), tt.c_list, tt.name, false, tt.p.CircleAroundWaypointGeodesic(tt.r, 64))

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
		{ name: "PiP - external point0", w_list: w_list, c_list: c_list, p: w, wantErr: false },
		{ name: "PiP - external point1", w_list: w_list, c_list: c_list, p: w2, wantErr: false },
		{ name: "PiP - external point2", w_list: w_list, c_list: c_list, p: w_list[0], wantErr: false },
		{ name: "PiP - external point3", w_list: w_list, c_list: c_list, p: w_list[1], wantErr: false },
		{ name: "PiP - external point4", w_list: w_list, c_list: c_list, p: w_list[2], wantErr: false },
		{ name: "PiP - external point5", w_list: w_list, c_list: c_list, p: w_list[3], wantErr: false },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := storage.NewMemoryStorage(tt.w_list, tt.c_list)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			_, gotErr := m.IsPointInObstacles(tt.p)

			// TODO: For visually testing, export results in geojson
			utils.ExportToGeoJSON(append(tt.w_list, tt.p), tt.c_list, tt.name, false)

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
		{ name: "LiP - samePoint", w_list: w_list, c_list: c_list, p: w, p2: w, want: false, wantErr: false },
		{ name: "LiP - w-w2", w_list: w_list, c_list: c_list, p: w, p2: w2, want: true, wantErr: false },
		{ name: "LiP - 0-1", w_list: w_list, c_list: c_list, p: w_list[0], p2: w_list[1], want: true, wantErr: false },
		{ name: "LiP - 1-2", w_list: w_list, c_list: c_list, p: w_list[1], p2: w_list[2], want: true, wantErr: false },
		{ name: "LiP - 2-3", w_list: w_list, c_list: c_list, p: w_list[2], p2: w_list[3], want: true, wantErr: false },
		{ name: "LiP - 3-w", w_list: w_list, c_list: c_list, p: w_list[3], p2: w, want: true, wantErr: false },
		{ name: "LiP - w-w3", w_list: w_list, c_list: c_list, p: w, p2: w3, want: false, wantErr: false },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := storage.NewMemoryStorage(tt.w_list, tt.c_list)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got, line, gotErr := m.IsLineInObstacles(tt.p, tt.p2)

			// TODO: For visually testing, export results in geojson
			utils.ExportToGeoJSON(line, tt.c_list, tt.name, true)

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
