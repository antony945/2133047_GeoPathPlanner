package algorithm_test

import (
	"geopathplanner/routing/internal/algorithm"
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/storage"
	"geopathplanner/routing/internal/utils"
	"testing"
)

func TestRRTStarAlgorithm_run(t *testing.T) {
	a, _ := models.NewAltitude(100, models.MT)
	w1 := models.MustNewWaypoint(0, 50.872778105839274, 4.433724687935722, a) // p1
	w2 := models.MustNewWaypoint(1, 50.884400404439646, 4.46992531620532, a)  // p2
	// w3 := models.MustNewWaypoint(2, 50.890383059561145, 4.45503208121508, a)  // p3
	// w4 := models.MustNewWaypoint(3, 50.888976908715335, 4.434797815001701, a)  // p4

	c_list := []*models.Feature3D{
		models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              4.452166432497052,
              50.87565312347229
            ],
            [
              4.45389986799907,
              50.874291381256256
            ],
            [
              4.4498672299481825,
              50.88096516631941
            ],
            [
              4.451848154983281,
              50.8813222592305
            ],
            [
              4.45478419822345,
              50.87636729027514
            ],
            [
              4.4560577382486315,
              50.87672438629892
            ],
            [
              4.451812846246014,
              50.883420200613955
            ],
            [
              4.447815509645977,
              50.88214806230394
            ],
            [
              4.452166432497052,
              50.87565312347229
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
		models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              4.444525609715754,
              50.88774978175786
            ],
            [
              4.445126892739296,
              50.89074037823755
            ],
            [
              4.452025073805913,
              50.88824057502404
            ],
            [
              4.447107954561972,
              50.89129841394083
            ],
            [
              4.4405988631421,
              50.89109722235199
            ],
            [
              4.444525609715754,
              50.88774978175786
            ]
          ]
        ],
        "type": "Polygon"
      }
    }`),
	}

	c_overlapping := []*models.Feature3D{
		models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              4.45503183842132,
              50.883978094618584
            ],
            [
              4.459630799842046,
              50.878934067344744
            ],
            [
              4.464052999538779,
              50.88092054086215
            ],
            [
              4.462001401679402,
              50.883665840890075
            ],
            [
              4.45503183842132,
              50.883978094618584
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
                4.446718832978945,
                50.889579933133234
              ],
              [
                4.448204632144268,
                50.88801756100412
              ],
              [
                4.449089057639668,
                50.89082977389265
              ],
              [
                4.451388478480737,
                50.89094112448484
              ],
              [
                4.4508578609232075,
                50.89129825498483
              ],
              [
                4.446400418720572,
                50.89078514625584
              ],
              [
                4.446718832978945,
                50.889579933133234
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
              4.445763705007465,
              50.88904423628691
            ],
            [
              4.445480667924812,
              50.890383286240564
            ],
            [
              4.438865573479632,
              50.89085162671603
            ],
            [
              4.438618309202809,
              50.88687920409194
            ],
            [
              4.445763705007465,
              50.88904423628691
            ]
          ]
        ],
        "type": "Polygon"
      }
    }`),
	}

	sv := models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              4.424480223895728,
              50.89367115387381
            ],
            [
              4.424480223895728,
              50.867778999101745
            ],
            [
              4.480653371960557,
              50.867778999101745
            ],
            [
              4.480653371960557,
              50.89367115387381
            ],
            [
              4.424480223895728,
              50.89367115387381
            ]
          ]
        ],
        "type": "Polygon"
      }
    }`)

	tests := []struct {
		name string // description of this test case
    storageType models.StorageType
		// Named input parameters for target function.
		searchVolume *models.Feature3D
		start        *models.Waypoint
		end          *models.Waypoint
		constraints  []*models.Feature3D
		want         []*models.Waypoint
		wantCost     float64
		wantErr      bool
	}{
		{name: "RRTStar with non-overlapping obstacles - MEMORY", storageType: models.Memory, searchVolume: sv, start: w1, end: w2, constraints: c_list, wantErr: false},
		{name: "RRTStar with non-overlapping obstacles - RTREE", storageType: models.RTree, searchVolume: sv, start: w1, end: w2, constraints: c_list, wantErr: false},
		{name: "RRTStar with overlapping obstacles - MEMORY", storageType: models.Memory, searchVolume: sv, start: w1, end: w2, constraints: append(c_list, c_overlapping...), wantErr: false},
		{name: "RRTStar with overlapping obstacles - RTREE", storageType: models.RTree, searchVolume: sv, start: w1, end: w2, constraints: append(c_list, c_overlapping...), wantErr: false},
		{name: "RRTStar with no obstacles - MEMORY", storageType: models.Memory, searchVolume: sv, start: w1, end: w2, constraints: []*models.Feature3D{}, wantErr: false},
		{name: "RRTStar with no obstacles - RTREE", storageType: models.RTree, searchVolume: sv, start: w1, end: w2, constraints: []*models.Feature3D{}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := algorithm.NewRRTStarAlgorithm()
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}

			s, err := storage.NewEmptyStorage(tt.storageType)
			if err != nil {
				t.Fatalf("could not construct storage: %v", err)
			}
			s.AddConstraints(tt.constraints)

			got, _, gotErr := a.Run(tt.searchVolume, tt.start, tt.end, nil, s)

			// TODO: For visually testing, export results in geojson
			utils.MarkWaypointsAsOriginal(tt.start, tt.end)
			utils.ExportToGeoJSONRoute("algorithm", got, tt.constraints, tt.searchVolume, tt.name, true)

			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Run() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Run() succeeded unexpectedly")
			}
		})
	}
}

func TestRRTStarAlgorithm_Compute(t *testing.T) {
	a, _ := models.NewAltitude(100, models.MT)
	w_list := []*models.Waypoint{
    models.MustNewWaypoint(0, 50.872778105839274, 4.433724687935722, a),
    models.MustNewWaypoint(1, 50.884400404439646, 4.46992531620532, a),
    models.MustNewWaypoint(2, 50.890383059561145, 4.45503208121508, a),
    models.MustNewWaypoint(3, 50.888976908715335, 4.434797815001701, a),
  	}

	c_list := []*models.Feature3D{
    models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              4.452166432497052,
              50.87565312347229
            ],
            [
              4.45389986799907,
              50.874291381256256
            ],
            [
              4.4498672299481825,
              50.88096516631941
            ],
            [
              4.451848154983281,
              50.8813222592305
            ],
            [
              4.45478419822345,
              50.87636729027514
            ],
            [
              4.4560577382486315,
              50.87672438629892
            ],
            [
              4.451812846246014,
              50.883420200613955
            ],
            [
              4.447815509645977,
              50.88214806230394
            ],
            [
              4.452166432497052,
              50.87565312347229
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
    models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              4.444525609715754,
              50.88774978175786
            ],
            [
              4.445126892739296,
              50.89074037823755
            ],
            [
              4.452025073805913,
              50.88824057502404
            ],
            [
              4.447107954561972,
              50.89129841394083
            ],
            [
              4.4405988631421,
              50.89109722235199
            ],
            [
              4.444525609715754,
              50.88774978175786
            ]
          ]
        ],
        "type": "Polygon"
      }
    }`),
	}

	c_overlapping := []*models.Feature3D{
    models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              4.45503183842132,
              50.883978094618584
            ],
            [
              4.459630799842046,
              50.878934067344744
            ],
            [
              4.464052999538779,
              50.88092054086215
            ],
            [
              4.462001401679402,
              50.883665840890075
            ],
            [
              4.45503183842132,
              50.883978094618584
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
                4.446718832978945,
                50.889579933133234
              ],
              [
                4.448204632144268,
                50.88801756100412
              ],
              [
                4.449089057639668,
                50.89082977389265
              ],
              [
                4.451388478480737,
                50.89094112448484
              ],
              [
                4.4508578609232075,
                50.89129825498483
              ],
              [
                4.446400418720572,
                50.89078514625584
              ],
              [
                4.446718832978945,
                50.889579933133234
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
              4.445763705007465,
              50.88904423628691
            ],
            [
              4.445480667924812,
              50.890383286240564
            ],
            [
              4.438865573479632,
              50.89085162671603
            ],
            [
              4.438618309202809,
              50.88687920409194
            ],
            [
              4.445763705007465,
              50.88904423628691
            ]
          ]
        ],
        "type": "Polygon"
      }
    }`),
  	}

  	sv := models.MustNewFeatureFromGeojson(`{
      "type": "Feature",
      "properties": {},
      "geometry": {
        "coordinates": [
          [
            [
              4.424480223895728,
              50.89367115387381
            ],
            [
              4.424480223895728,
              50.867778999101745
            ],
            [
              4.480653371960557,
              50.867778999101745
            ],
            [
              4.480653371960557,
              50.89367115387381
            ],
            [
              4.424480223895728,
              50.89367115387381
            ]
          ]
        ],
        "type": "Polygon"
      }
    }`)

	tests := []struct {
		name string // description of this test case
    storageType models.StorageType
		// Named input parameters for target function.
    searchVolume *models.Feature3D
		waypoints   []*models.Waypoint
		constraints []*models.Feature3D
		want    []*models.Waypoint
		wantCost float64
		wantErr bool
	}{
		{name: "RRTStarFull with non-overlapping obstacles - MEMORY", storageType: models.Memory, searchVolume: sv, waypoints: w_list, constraints: c_list, wantErr: false},
		{name: "RRTStarFull with non-overlapping obstacles - RTREE", storageType: models.RTree, searchVolume: sv, waypoints: w_list, constraints: c_list, wantErr: false},
		{name: "RRTStarFull with overlapping obstacles - MEMORY", storageType: models.Memory, searchVolume: sv, waypoints: w_list, constraints: append(c_list, c_overlapping...), wantErr: false},
		{name: "RRTStarFull with overlapping obstacles - RTREE", storageType: models.RTree, searchVolume: sv, waypoints: w_list, constraints: append(c_list, c_overlapping...), wantErr: false},
		{name: "RRTStarFull with no obstacles - MEMORY", storageType: models.Memory, searchVolume: sv, waypoints: w_list, constraints: []*models.Feature3D{}, wantErr: false},
		{name: "RRTStarFull with no obstacles - RTREE", storageType: models.RTree, searchVolume: sv, waypoints: w_list, constraints: []*models.Feature3D{}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := algorithm.NewRRTStarAlgorithm()
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			
			s, err := storage.NewEmptyStorage(tt.storageType)
			if err != nil {
				t.Fatalf("could not construct storage: %v", err)
			}

			got, _, gotErr := a.Compute(tt.searchVolume, tt.waypoints, tt.constraints, nil, s)

			// TODO: For visually testing, export results in geojson
      		utils.MarkWaypointsAsOriginal(tt.waypoints...)
			utils.ExportToGeoJSONRoute("algorithm", got, tt.constraints, tt.searchVolume, tt.name, true)

			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Compute() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Compute() succeeded unexpectedly")
			}
		})
	}
}