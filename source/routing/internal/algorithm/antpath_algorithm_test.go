package algorithm_test

import (
	"geopathplanner/routing/internal/algorithm"
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/storage"
	"geopathplanner/routing/internal/utils"
	"testing"
)

func TestAntPathAlgorithm_run(t *testing.T) {
	a, _ := models.NewAltitude(100, models.MT)
	w1 := models.MustNewWaypoint(0, 50.872778105839274, 4.433724687935722, a) // p1
	w2 := models.MustNewWaypoint(1, 50.884400404439646, 4.46992531620532, a)  // p2

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

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		start      *models.Waypoint
		end      *models.Waypoint
		constraints []*models.Feature3D
		want    []*models.Waypoint
		wantCost float64
		wantErr bool
	}{
		{name: "AntPath with non-overlapping obstacles", start: w1, end: w2, constraints: c_list, wantErr: false},
		{name: "AntPath with no obstacles", start: w1, end: w2, constraints: []*models.Feature3D{}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := algorithm.NewAntPathAlgorithm()
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			
			m, err := storage.NewEmptyMemoryStorage()
			if err != nil {
				t.Fatalf("could not construct memory storage: %v", err)
			}

			got, _, gotErr := a.Run(tt.start, tt.end, tt.constraints, nil, m)

			// TODO: For visually testing, export results in geojson
			utils.ExportToGeoJSON("algorithm", got, tt.constraints, tt.name, true)

			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetIntersectionPoints() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetIntersectionPoints() succeeded unexpectedly")
			}
		})
	}
}