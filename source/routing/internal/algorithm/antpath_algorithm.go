package algorithm

import (
	"fmt"
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/storage"
	"geopathplanner/routing/internal/utils"
	"runtime"
	"sync"
)

type AntPathAlgorithm struct {}

func NewAntPathAlgorithm() (*AntPathAlgorithm, error) {
	return &AntPathAlgorithm{}, nil
}

// Concurrency version of Compute function, where every pair of wps is processed in a separate goroutine.
// TODO: Still in testing
func (a *AntPathAlgorithm) ComputeConcurrently(searchVolume *models.Feature3D, waypoints []*models.Waypoint, constraints []*models.Feature3D, parameters map[string]any, storageType models.StorageType, maxWorkers int) ([]*models.Waypoint, float64, error) {
	// Check if waypoints are at least 2
	numPairs := len(waypoints) - 1
	if numPairs <= 0 {
		return nil, 0.0, fmt.Errorf("less than 2 waypoints submitted (%d): abort", len(waypoints))
	}
	
	maxCPU := runtime.NumCPU()
	if maxWorkers <= 0 {
		maxWorkers = min(maxCPU, numPairs)
	} else {
		// fmt.Printf("[WARN] Requested %d workers exceeds %d cores, limiting to %d",
		maxWorkers = min(maxWorkers, maxCPU, numPairs)
	}

	// If just 1 worker, use the normal version
	if maxWorkers == 1 {
		return a.Compute(searchVolume, waypoints, constraints, parameters, storageType)
	}

	// TODO: Buffer constraints
	// TODO: Think about creating and adding constraints in Run function so to parallelize that function
	// Create storage and load constraint into it
	storage, err := storage.NewEmptyStorage(storageType)
	if err != nil {
		return nil, 0.0, err
	}
	err = storage.AddConstraints(constraints)
	if err != nil {
		return nil, 0.0, err
	}

	// === Channels and synchronization structures ===
	jobs := make(chan job, numPairs)       // channel for distributing work
	results := make(chan result, numPairs) // channel to collect computed results
	var wg sync.WaitGroup // ensures all workers complete before closing results

	// 1. Create and start the workers
	for w := 0; w < maxWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for j := range jobs {
				// Run RRT for this pair of waypoints
				tmpRoute, tmpCost, err := a.Run(j.startWP, j.endWP, parameters, storage.Clone())
				if err != nil {
					results <- result{i: j.i, err: fmt.Errorf("worker %d: run AntPath: %w", workerID, err)}
					continue
				}

				results <- result{
					i:     j.i,
					route: tmpRoute,
					cost:  tmpCost,
					err:   nil,
				}
			}
		}(w)
	}

	// 2. Send jobs to workers
	for i := 0; i < numPairs; i++ {
		jobs <- job{
			i:       i,
			startWP: waypoints[i],
			endWP:   waypoints[i+1],
		}
	}
	close(jobs) // no more jobs to send

	// 3. Collect results
	go func() {
		wg.Wait()      // wait for all workers to finish
		close(results) // then close result channel
	}()

	// Store results in correct order
	routeSegments := make([][]*models.Waypoint, numPairs)
	costs := make([]float64, numPairs)
	var firstErr error

	for res := range results {
		if res.err != nil && firstErr == nil {
			firstErr = res.err
		}
		routeSegments[res.i] = res.route
		costs[res.i] = res.cost
	}

	if firstErr != nil {
		return nil, 0, firstErr
	}

	// 4. Merge results
	finalRoute := make([]*models.Waypoint, 0)
	totalCost := 0.0

	// Start from first waypoint
	if len(routeSegments) > 0 && len(routeSegments[0]) > 0 {
		finalRoute = append(finalRoute, routeSegments[0][0])
	}

	for i, seg := range routeSegments {
		if len(seg) == 0 {
			continue
		}
		totalCost += costs[i]
		// Skip first element to avoid duplicates
		finalRoute = append(finalRoute, seg[1:]...)
	}

	return finalRoute, totalCost, nil
}

func (a *AntPathAlgorithm) Compute(searchVolume *models.Feature3D, waypoints []*models.Waypoint, constraints []*models.Feature3D, parameters map[string]any, storageType models.StorageType) ([]*models.Waypoint, float64, error) {	
	// Check if waypoints are at least 2
	numPairs := len(waypoints) - 1
	if numPairs <= 0 {
		return nil, 0.0, fmt.Errorf("less than 2 waypoints submitted (%d): abort", len(waypoints))
	}
	
	// Create empty list of wps
	route := make([]*models.Waypoint, 0)
	cost := 0.0

	// 0. Load first wp
	route = append(route, waypoints[0])
	
	// TODO: Buffer constraints
	// TODO: Think about creating and adding constraints in Run function so to parallelize that function
	// 1. Create storage and load constraint into it
	storage, err := storage.NewEmptyStorage(storageType)
	if err != nil {
		return nil, 0.0, err
	}
	err = storage.AddConstraints(constraints)
	if err != nil {
		return nil, 0.0, err
	}

	// 2. For each pair of wp -> run antpath
	for i := 0; i < len(waypoints)-1; i++ {
		tmpRoute, tmpCost, err := a.Run(waypoints[i], waypoints[i+1], parameters, storage)
		if err != nil {
			// Return route until now
			return route, cost, fmt.Errorf("interrupted antpath for error between wp[%d] and wp[%d]: %w", i, i+1, err)
		}
		// Append new route but removing the first one
		route = append(route, tmpRoute[1:]...)
		cost += tmpCost
	}

	// TODO: Think if this is the correct place
	storage.Clear()
	
	// 3. Return everything
	return route, cost, nil
}

func (a *AntPathAlgorithm) Run(start, end *models.Waypoint, parameters map[string]any, storage storage.Storage) ([]*models.Waypoint, float64, error) {
	route := make([]*models.Waypoint, 0)
	// cost := 0.0

	// ------------------------------------------------------------------------------------------------------

	// Get intersection points
	intersectionPoints, err := storage.GetIntersectionPoints(start, end)
	if err != nil {
		return nil, 0.0, err
	}
	route = append(route, start)
	
	// For every intersectionPoint struct get best way to go around obstacle
	for i, ip := range intersectionPoints {
		fmt.Printf("ip[%d]: intersects with %d polygons\n", i, len(ip.Polygons))
		var polygonToCheck *models.Feature3D
		if len(ip.Polygons) > 1 {
			// Union the polygons
			// fmt.Printf("trying to union them...")
			unionedFeatures, err := utils.UnionFeatures(ip.Polygons)
			if err != nil {
				return nil, 0.0, err
			}
			// TODO: Check that we have just 1 unionedFeatures
			// if len(unionedFeatures) > 0 {
				// return nil, 0.0, fmt.Errorf("features after union are more than 1, results may not be accurate")
			// }
			fmt.Printf("ip[%d]: union finished. got %d polygons\n", i, len(unionedFeatures))
			// for i, v := range unionedFeatures {
				// fmt.Printf("[%d]: %+v\n", i, v)
			// }
			polygonToCheck = unionedFeatures[0]
		} else {
			polygonToCheck = ip.Polygons[0]
		}

		bestForPolygonWay := utils.GetBestWayToGoAroundPolygon(polygonToCheck, ip.EnteringPoint, ip.ExitingPoint)
		route = append(route, bestForPolygonWay...)
	}

	// ------------------------------------------------------------------------------------------------------

	route = append(route, end)
	cost := utils.TotalHaversineDistance(route)
	return route, cost, nil
}
