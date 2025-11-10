package algorithm

import (
	"fmt"
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/storage"
	"geopathplanner/routing/internal/utils"
	"math"
	"runtime"
	"sync"
)

const (
	K_INIT float64 = 2*math.E
	R_INIT_MT float64 = 2.0
)

type RRTStarAlgorithm struct {
	*RRTAlgorithm
}


func NewRRTStarAlgorithm() (*RRTStarAlgorithm, error) {
	a, err := NewRRTAlgorithm()
	if err != nil {
		return nil, err
	}

	return &RRTStarAlgorithm{
		a,
	}, nil
}

// Concurrency version of Compute function, where every pair of wps is processed in a separate goroutine.
// TODO: Still in testing
func (a *RRTStarAlgorithm) ComputeConcurrently(searchVolume *models.Feature3D, waypoints []*models.Waypoint, constraints []*models.Feature3D, parameters map[string]any, storageType models.StorageType, maxWorkers int) ([]*models.Waypoint, float64, error) {
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
				tmpRoute, tmpCost, err := a.Run(searchVolume, j.startWP, j.endWP, parameters, storage.Clone())
				if err != nil {
					results <- result{i: j.i, err: fmt.Errorf("worker %d: run RRT*: %w", workerID, err)}
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

func (a *RRTStarAlgorithm) Compute(searchVolume *models.Feature3D, waypoints []*models.Waypoint, constraints []*models.Feature3D, parameters map[string]any, storageType models.StorageType) ([]*models.Waypoint, float64, error) {
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

	// 2. For each pair of wp -> run rrt
	for i := 0; i < len(waypoints)-1; i++ {
		tmpRoute, tmpCost, err := a.Run(searchVolume, waypoints[i], waypoints[i+1], parameters, storage)
		if err != nil {
			// Return route until now
			return route, cost, fmt.Errorf("interrupted RRTStar for error between wp[%d] and wp[%d]: %w", i, i+1, err)
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

func (a *RRTStarAlgorithm) Run(searchVolume *models.Feature3D, start, end *models.Waypoint, parameters map[string]any, storage storage.Storage) ([]*models.Waypoint, float64, error) {
	// TODO: Think if this is the correct place
	storage.ClearWaypoints()
	fmt.Printf("wpA: %v, wpB: %v, storage starts with %d constraints and %d sampled waypoints.\n", start, end, storage.ConstraintsLen(), storage.WaypointsLen())

	// First thing to do if to check if a straight line connection is possible
	if obstacleBetweenStartEnd, _, _ := storage.IsLineInObstacles(start, end); !obstacleBetweenStartEnd {
		fmt.Printf("Goal immediately found: straight line collision-free\n\n")
		return []*models.Waypoint{start, end}, utils.HaversineDistance3D(start, end), nil
	}

	// TODO: Parameters
	// TODO: Think about not to use max_iterations directly, rather continue until a certain condition happen (e.g. cost of route stopped decreasing for a while) 
	// Get Parameters
	sampler, max_iterations, step_size_mt, _ := a.GetParameters(parameters, end)

	// Add start to storage
	err := storage.AddWaypointWithPrevious(nil, start)
	if err != nil {
		return nil, 0.0, err
	}
	goal_found := false

	// ------------------------------------------------------------------------------------------------------

	for current_iter := range max_iterations {
		// Change K according to cardinality of V (no. of nodes)
		K := int(K_INIT * math.Log(float64(storage.WaypointsLen())))+1
		// R := math.Max(R_INIT_MT * math.Sqrt(math.Log(float64(storage.WaypointsLen()))/float64(storage.WaypointsLen())), step_size_mt)

		if current_iter % 1000 == 0 {
			if goal_found {
				route, _ := storage.GetPathToRoot(end)
				cost_km := utils.TotalHaversineDistance(route)
				fmt.Printf("[%d/%d] k: %d, #wps: %d, #routeWps: %d, routeCost: %.3f mt\n", current_iter, max_iterations, K, storage.WaypointsLen(), len(route), cost_km)
				// fmt.Printf("[%d/%d] radius: %.2fmt, #wps: %d, cost: %.3f mt\n", current_iter, MAX_ITERATIONS, R, len(route), cost_km)
			
			} else {
				fmt.Printf("[%d/%d] k: %d, #wps: %d, goal not found yet\n", current_iter, max_iterations, K, storage.WaypointsLen())
				// fmt.Printf("[%d/%d] radius: %.2fmt, goal not found yet\n", current_iter, MAX_ITERATIONS, R)
			}
		}

		// 1. Sample a new free wp
		sampled, err := storage.SampleFree(sampler, searchVolume, end.Alt)
		if err != nil {
			// Impossible to sample 
			return nil, 0.0, err
		}

		// 2. Get nearest wp
		nearest, _, err := storage.NearestPoint(sampled)
		if err != nil {
			return nil, 0.0, err
		}

		// 3. Find wp starting from nearest in direction of sampled at distance (steering) step_size_mt
		new := utils.GetPointInDirectionAtDistance(nearest, sampled, step_size_mt)

		// 4. Check if connection from nearest to new is possible
		isInObstacles, _, err := storage.IsLineInObstacles(nearest, new)
		if err != nil {
			return nil, 0.0, err
		}
		if isInObstacles {
			continue
		}

		// Here we know that nearest can connect to new
		// But we need to check if in the neighbors of new there could be a better one
		// 5. Add new to nearest
		// err = storage.AddWaypointWithPrevious(nearest, new)
		// if err != nil {
		// 	return nil, 0.0, err
		// }

		// 5. Here you check all the neighbors to connect to min cost path and rewire the tree
		_, err = a.ConnectAndRewire(new, nearest, K, storage)
		if err != nil {
			return nil, 0.0, err
		}
		// if rewired && goal_found {
		// 	route, _ := storage.GetPathToRoot(end)
		// 	cost_km := utils.TotalHaversineDistance(route)
		// 	fmt.Printf("Rewired Tree\n#wps: %d, cost: %.3f mt\n", len(route), cost_km)
		// }

		// 6. Check if it's goal
		if !goal_found && a.isGoal(new, end, step_size_mt) {
			// 7. Check if can be connected to goal
			isInObstacles, _, err := storage.IsLineInObstacles(new, end)
			if err != nil {
				return nil, 0.0, err
			}
			// If yes, connect to goal, found!!! but not break
			if !isInObstacles {
				err := storage.AddWaypointWithPrevious(new, end)
				if err != nil {
					return nil, 0.0, err
				}
				
				route, _ := storage.GetPathToRoot(end)
				cost_km := utils.TotalHaversineDistance(route)
				fmt.Printf("New goal found at iteration %d/%d.\n", current_iter, max_iterations)
				fmt.Printf("#wps: %d, cost: %.3f mt\n", len(route), cost_km)
				goal_found = true
			}
		}
	}

	// Now try to obtain the route from the goal node
	if goal_found {
		route, err := storage.GetPathToRoot(end)
		if err != nil {
			return nil, 0.0, err
		}
		cost := utils.TotalHaversineDistance(route)
		return route, cost, nil
	} else {
		return nil, 0.0, fmt.Errorf("goal not found with %d iterations", max_iterations)
	}
}

func (a *RRTStarAlgorithm) GetParameters(parameters map[string]any, goal *models.Waypoint) (utils.Sampler, int, float64, float64) {
	SAMPLER, MAX_ITERATIONS, STEP_SIZE_MT, GOAL_BIAS := a.RRTAlgorithm.GetParameters(parameters, goal)
	
	// TODO: Delete this, just for debug
	fmt.Printf("k_init: %d\n", int(math.Floor(K_INIT)))
	fmt.Printf("r_init_mt: %f\n", R_INIT_MT)
	fmt.Printf("--------------------------------------------------------\n")

	return SAMPLER, MAX_ITERATIONS, STEP_SIZE_MT, GOAL_BIAS
}

func (a *RRTStarAlgorithm) ConnectAndRewire(new, nearest *models.Waypoint, k int, storage storage.Storage) (bool, error) {
	// TODO: Test also with radius
	neighbors, distances, err := storage.KNearestPoints(new, k)
	if err != nil {
		return false, fmt.Errorf("error while getting the %d-nn of %v: %+w", k, new, err)
	}

	return a.connectAndRewireWithNeighbors(new, nearest, neighbors, distances, storage)
}

func (a *RRTStarAlgorithm) ConnectAndRewireInRadius(new, nearest *models.Waypoint, radius_mt float64, storage storage.Storage) (bool, error) {
	neighbors, distances, err := storage.NearestPointsInRadius(new, radius_mt)
	if err != nil {
		return false, fmt.Errorf("error while getting point within %.2f mt of %v: %+w", radius_mt, new, err)
	}

	return a.connectAndRewireWithNeighbors(new, nearest, neighbors, distances, storage)
}

func (a *RRTStarAlgorithm) connectAndRewireWithNeighbors(new, nearest *models.Waypoint, neighbors []*models.Waypoint, distances []float64, storage storage.Storage) (bool, error) {
	// CONNECT
	// For every neighbor check if connecting to new via that would be better compared to connect to nearest
	minCostWp := nearest
	minCost, err := a.getCost(minCostWp, storage)
	if err != nil {
		return false, err
	}
	minCost += distances[0]

	// Scan neighbors
	// TODO: you can skip first one as it will be the nearest
	isInObstacleList := make(map[int]bool)

	for idx, near := range neighbors {
		// Check if near can be connected to new
		isInObstacles, _, err := storage.IsLineInObstacles(near, new)
		if err != nil {
			return false, err
		}
		isInObstacleList[idx] = isInObstacles
		if isInObstacles {
			continue
		}

		// Here near can be connected to new, check the cost
		currentCost, err := a.getCost(near, storage)
		if err != nil {
			return false, err
		}
		// TODO: Think if it make sense
		// currentCost += a.getLineCost(near, new)
		currentCost += distances[idx]

		if currentCost < minCost {
			minCostWp = near
			minCost = currentCost
		}
	}

	// Here connect to minCostWp
	err = storage.AddWaypointWithPrevious(minCostWp, new)
	if err != nil {
		return false, err
	}

	// -----------------------------------------------------------

	// Now rewire the tree
	rewired := false
	newCost, err := a.getCost(new, storage)
	if err != nil {
		return rewired, err
	}

	for idx, near := range neighbors {
		// Check if new can be connected to near
		// Rehuse list of before
		if (isInObstacleList[idx]) {
			continue
		}

		// Here near can be connected to new, check the cost
		newCost += a.getLineCost(new, near)
		currentCost, err := a.getCost(near, storage)
		if err != nil {
			return rewired, err
		}

		if newCost < currentCost {
			// New becomes the parent of near
			storage.ChangePrevious(new, near)
			rewired = true
		}
	}
	
	return rewired, nil
}

func (a *RRTStarAlgorithm) getCost(wp *models.Waypoint, storage storage.Storage) (float64, error) {
	// route, err := storage.GetPathToRoot(wp)
	// if err != nil {
	// 	return 0.0, err
	// }
	// cost := utils.TotalHaversineDistance(route)
	// return cost, nil

	return storage.GetCostToRoot(wp)
}

func (a *RRTStarAlgorithm) getLineCost(start, end *models.Waypoint) (float64) {
	// get cost of connecting start to end (haversine distance cost)
	return utils.HaversineDistance3D(start, end)
}