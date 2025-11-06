package algorithm_test

import (
	"geopathplanner/routing/internal/algorithm"
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/storage"
	"geopathplanner/routing/internal/utils"
	"testing"
)

func TestAntPathAlgorithm_run(t *testing.T) {
  _, w_list, c_list, c_overlapping := utils.SetupTestScenario()
  w1 := w_list[0]
  w2 := w_list[1]

	tests := []struct {
		name string // description of this test case
    storageType models.StorageType
		// Named input parameters for target function.
		start      *models.Waypoint
		end      *models.Waypoint
		constraints []*models.Feature3D
		want    []*models.Waypoint
		wantCost float64
		wantErr bool
	}{
		// {name: "AntPath with non-overlapping obstacles - LIST", storageType: models.List, start: w1, end: w2, constraints: c_list, wantErr: false},
		{name: "AntPath with non-overlapping obstacles - RTREE", storageType: models.RTree, start: w1, end: w2, constraints: c_list, wantErr: false},
		// {name: "AntPath with overlapping obstacles - LIST", storageType: models.List, start: w1, end: w2, constraints: append(c_list, c_overlapping...), wantErr: false},
		{name: "AntPath with overlapping obstacles - RTREE", storageType: models.RTree, start: w1, end: w2, constraints: append(c_list, c_overlapping...), wantErr: false},
		// {name: "AntPath with no obstacles - LIST", storageType: models.List, start: w1, end: w2, constraints: []*models.Feature3D{}, wantErr: false},
		{name: "AntPath with no obstacles - RTREE", storageType: models.RTree, start: w1, end: w2, constraints: []*models.Feature3D{}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := algorithm.NewAntPathAlgorithm()
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}

      s, err := storage.NewEmptyStorage(tt.storageType)
      if err != nil {
				t.Fatalf("could not construct storage: %v", err)
			}
      s.AddConstraints(tt.constraints)

			got, _, gotErr := a.Run(tt.start, tt.end, nil, s)

			// TODO: For visually testing, export results in geojson
			// utils.ExportToGeoJSON("algorithm", got, tt.constraints, tt.name, true)
      		utils.MarkWaypointsAsOriginal(tt.start, tt.end)
			utils.ExportToGeoJSONRoute("algorithm", got, tt.constraints, nil, tt.name, true)

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

func TestAntPathAlgorithm_Compute(t *testing.T) {
  _, w_list, c_list, c_overlapping := utils.SetupTestScenario()

	tests := []struct {
		name string // description of this test case
    storageType models.StorageType
		// Named input parameters for target function.
		waypoints   []*models.Waypoint
		constraints []*models.Feature3D
		want    []*models.Waypoint
		wantCost float64
		wantErr bool
	}{
		// {name: "AntPathFull with non-overlapping obstacles - LIST", storageType: models.List, waypoints: w_list, constraints: c_list, wantErr: false},
		{name: "AntPathFull with non-overlapping obstacles - RTREE", storageType: models.RTree, waypoints: w_list, constraints: c_list, wantErr: false},
		// {name: "AntPathFull with overlapping obstacles - LIST", storageType: models.List, waypoints: w_list, constraints: append(c_list, c_overlapping...), wantErr: false},
		{name: "AntPathFull with overlapping obstacles - RTREE", storageType: models.RTree, waypoints: w_list, constraints: append(c_list, c_overlapping...), wantErr: false},
		// {name: "AntPathFull with no obstacles - LIST", storageType: models.List, waypoints: w_list, constraints: []*models.Feature3D{}, wantErr: false},
		{name: "AntPathFull with no obstacles - RTREE", storageType: models.RTree, waypoints: w_list, constraints: []*models.Feature3D{}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := algorithm.NewAntPathAlgorithm()
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
      }

			got, _, gotErr := a.Compute(nil, tt.waypoints, tt.constraints, nil, tt.storageType)

			// TODO: For visually testing, export results in geojson
			// utils.ExportToGeoJSON("algorithm", got, tt.constraints, tt.name, true)
      utils.MarkWaypointsAsOriginal(tt.waypoints...)
			utils.ExportToGeoJSONRoute("algorithm", got, tt.constraints, nil, tt.name, true)

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

func TestAntPathAlgorithm_ComputeConcurrently(t *testing.T) {
  	_, w_list, c_list, c_overlapping := utils.SetupTestScenario()

	tests := []struct {
		name string // description of this test case
    storageType models.StorageType
		// Named input parameters for target function.
		waypoints   []*models.Waypoint
		constraints []*models.Feature3D
    maxWorkers int
		want    []*models.Waypoint
		wantCost float64
		wantErr bool
	}{
		{name: "ConcurrentAntPathFull -1 with no obstacles - RTREE", storageType: models.RTree, waypoints: w_list, constraints: []*models.Feature3D{}, maxWorkers: -1, wantErr: false},
		{name: "ConcurrentAntPathFull 1 with no obstacles - RTREE", storageType: models.RTree, waypoints: w_list, constraints: []*models.Feature3D{}, maxWorkers: 1, wantErr: false},
    {name: "ConcurrentAntPathFull 3 with no obstacles - RTREE", storageType: models.RTree, waypoints: w_list, constraints: []*models.Feature3D{}, maxWorkers: 3, wantErr: false},
		{name: "ConcurrentAntPathFull 3 with non-overlapping obstacles - RTREE", storageType: models.RTree, waypoints: w_list, constraints: c_list, maxWorkers: 3, wantErr: false},
		{name: "ConcurrentAntPathFull 3 with overlapping obstacles - RTREE", storageType: models.RTree, waypoints: w_list, constraints: append(c_list, c_overlapping...), maxWorkers: 3, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := algorithm.NewAntPathAlgorithm()
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}

			got, _, gotErr := a.ComputeConcurrently(nil, tt.waypoints, tt.constraints, nil, tt.storageType, tt.maxWorkers)

			// TODO: For visually testing, export results in geojson
      		utils.MarkWaypointsAsOriginal(tt.waypoints...)
			utils.ExportToGeoJSONRoute("algorithm", got, tt.constraints, nil, tt.name, true)

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