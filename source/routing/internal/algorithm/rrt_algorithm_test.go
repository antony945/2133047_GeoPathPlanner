package algorithm_test

import (
	"geopathplanner/routing/internal/algorithm"
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/storage"
	"geopathplanner/routing/internal/utils"
	"testing"
)

func TestRRTAlgorithm_run(t *testing.T) {
  sv, w_list, c_list, c_overlapping := utils.SetupTestScenario()
  w1 := w_list[0]
  w2 := w_list[1]

	tests := []struct {
		name string // description of this test case
    storageType models.StorageType
		// Named input parameters for target function.
    searchVolume *models.Feature3D
		start      *models.Waypoint
		end      *models.Waypoint
		constraints []*models.Feature3D
		want    []*models.Waypoint
		wantCost float64
		wantErr bool
	}{
		// {name: "RRT with non-overlapping obstacles - LIST", storageType: models.List, searchVolume: sv, start: w1, end: w2, constraints: c_list, wantErr: false},
		{name: "RRT with non-overlapping obstacles - RTREE", storageType: models.RTree, searchVolume: sv, start: w1, end: w2, constraints: c_list, wantErr: false},
		// {name: "RRT with overlapping obstacles - LIST",     storageType: models.List, searchVolume: sv, start: w1, end: w2, constraints: append(c_list, c_overlapping...), wantErr: false},
		{name: "RRT with overlapping obstacles - RTREE",     storageType: models.RTree, searchVolume: sv, start: w1, end: w2, constraints: append(c_list, c_overlapping...), wantErr: false},
		// {name: "RRT with no obstacles - LIST",              storageType: models.List, searchVolume: sv, start: w1, end: w2, constraints: []*models.Feature3D{}, wantErr: false},
		{name: "RRT with no obstacles - RTREE",              storageType: models.RTree, searchVolume: sv, start: w1, end: w2, constraints: []*models.Feature3D{}, wantErr: false},
  }
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := algorithm.NewRRTAlgorithm()
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

func TestRRTAlgorithm_Compute(t *testing.T) {
  sv, w_list, c_list, c_overlapping := utils.SetupTestScenario()

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
		// {name: "RRTFull with non-overlapping obstacles - LIST", storageType: models.List, searchVolume: sv, waypoints: w_list, constraints: c_list, wantErr: false},
		{name: "RRTFull with non-overlapping obstacles - RTREE", storageType: models.RTree, searchVolume: sv, waypoints: w_list, constraints: c_list, wantErr: false},

		// {name: "RRTFull with overlapping obstacles - LIST", storageType: models.List, searchVolume: sv, waypoints: w_list, constraints: append(c_list, c_overlapping...), wantErr: false},
		{name: "RRTFull with overlapping obstacles - RTREE", storageType: models.RTree, searchVolume: sv, waypoints: w_list, constraints: append(c_list, c_overlapping...), wantErr: false},

    // {name: "RRTFull with no obstacles - LIST", storageType: models.List, searchVolume: sv, waypoints: w_list, constraints: []*models.Feature3D{}, wantErr: false},
    {name: "RRTFull with no obstacles - RTREE", storageType: models.RTree, searchVolume: sv, waypoints: w_list, constraints: []*models.Feature3D{}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := algorithm.NewRRTAlgorithm()
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}

			got, _, gotErr := a.Compute(tt.searchVolume, tt.waypoints, tt.constraints, nil, tt.storageType)

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

func TestRRTAlgorithm_ComputeConcurrently(t *testing.T) {
  sv, w_list, c_list, c_overlapping := utils.SetupTestScenario()

	tests := []struct {
		name string // description of this test case
    storageType models.StorageType
		// Named input parameters for target function.
    searchVolume *models.Feature3D
		waypoints   []*models.Waypoint
		constraints []*models.Feature3D
    maxWorkers int
		want    []*models.Waypoint
		wantCost float64
		wantErr bool
	}{
		{name: "ConcurrentRRTFull -1 with no obstacles - RTREE", storageType: models.RTree, searchVolume: sv, waypoints: w_list, constraints: []*models.Feature3D{}, maxWorkers: -1, wantErr: false},
		{name: "ConcurrentRRTFull 1 with no obstacles - RTREE", storageType: models.RTree, searchVolume: sv, waypoints: w_list, constraints: []*models.Feature3D{}, maxWorkers: 1, wantErr: false},
    {name: "ConcurrentRRTFull 3 with no obstacles - RTREE", storageType: models.RTree, searchVolume: sv, waypoints: w_list, constraints: []*models.Feature3D{}, maxWorkers: 3, wantErr: false},
		{name: "ConcurrentRRTFull 3 with non-overlapping obstacles - RTREE", storageType: models.RTree, searchVolume: sv, waypoints: w_list, constraints: c_list, maxWorkers: 3, wantErr: false},
		{name: "ConcurrentRRTFull 3 with overlapping obstacles - RTREE", storageType: models.RTree, searchVolume: sv, waypoints: w_list, constraints: append(c_list, c_overlapping...), maxWorkers: 3, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := algorithm.NewRRTAlgorithm()
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}

			got, _, gotErr := a.ComputeConcurrently(tt.searchVolume, tt.waypoints, tt.constraints, nil, tt.storageType, tt.maxWorkers)

			// TODO: For visually testing, export results in geojson
      		utils.MarkWaypointsAsOriginal(tt.waypoints...)
			utils.ExportToGeoJSONRoute("algorithm", got, tt.constraints, tt.searchVolume, tt.name, true)

			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ComputeConcurrently() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ComputeConcurrently() succeeded unexpectedly")
			}
		})
	}
}