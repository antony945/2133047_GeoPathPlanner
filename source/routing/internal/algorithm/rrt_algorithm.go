package algorithm

import (
	"fmt"
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/storage"
	"geopathplanner/routing/internal/utils"
	"runtime"
	"sync"
)

type RRTAlgorithm struct {
}

func NewRRTAlgorithm() (*RRTAlgorithm, error) {
	return &RRTAlgorithm{}, nil
}

// Concurrency version of Compute function, where every pair of wps is processed in a separate goroutine.
// TODO: Still in testing
func (a *RRTAlgorithm) ComputeConcurrently(searchVolume *models.Feature3D, waypoints []*models.Waypoint, constraints []*models.Feature3D, parameters map[string]any, storageType models.StorageType, maxWorkers int) ([]*models.Waypoint, float64, error) {
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
					results <- result{i: j.i, err: fmt.Errorf("worker %d: run RRT: %w", workerID, err)}
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

	// TODO: Since here we know what pair had an error, we could think about rerunning just that portion.
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

func (a *RRTAlgorithm) Compute(searchVolume *models.Feature3D, waypoints []*models.Waypoint, constraints []*models.Feature3D, parameters map[string]any, storageType models.StorageType) ([]*models.Waypoint, float64, error) {
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
		tmpRoute, tmpCost, err := a.Run(searchVolume, waypoints[i], waypoints[i+1], parameters, storage.Clone())
		if err != nil {
			// Return route until now
			return route, cost, fmt.Errorf("interrupted RRT for error between wp[%d] and wp[%d]: %w", i, i+1, err)
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

func (a *RRTAlgorithm) Run(searchVolume *models.Feature3D, start, end *models.Waypoint, parameters map[string]any, storage storage.Storage) ([]*models.Waypoint, float64, error) {
	// TODO: Think if this is the correct place
	storage.ClearWaypoints()
	fmt.Printf("Storage has %d constraints and %d waypoints.\n\n", storage.ConstraintsLen(), storage.WaypointsLen())

	// First thing to do if to check if a straight line connection is possible
	if obstacleBetweenStartEnd, _, _ := storage.IsLineInObstacles(start, end); !obstacleBetweenStartEnd {
		return []*models.Waypoint{start, end}, utils.HaversineDistance3D(start, end), nil
	}
	
	// HERE IMPLEMENT RRT
	// start from start
	// sample a new free wp
	// get nearest wp
	// check if connection from nearest to free it's possible
	// if yes connect them, if no sample new one
	// check if it's goal (approximately)
	// check if connection from free to goal is possible
	// if yes break from the loop, if not continue


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
		if current_iter % 1000 == 0 {
			fmt.Printf("[%d/%d] #wps: %d, goal not found yet\n", current_iter, max_iterations, storage.WaypointsLen())
			// fmt.Printf("[%d/%d] radius: %.2fmt, goal not found yet\n", current_iter, MAX_ITERATIONS, R)
		}

		// 1. Sample a new free wp
		sampled, err := storage.SampleFree(sampler, searchVolume, start.Alt)
		if err != nil {
			// Impossible to sample 
			return nil, 0.0, err
		}

		// 2. Get nearest wp
		nearest, _, err := storage.NearestPoint(sampled)
		if err != nil {
			return nil, 0.0, err
		}

		// 3. Find wp starting from nearest in direction of sampled at distance (steering) STEP_SIZE_MT
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
		// 5. Add new to nearest
		err = storage.AddWaypointWithPrevious(nearest, new)
		if err != nil {
			return nil, 0.0, err
		}

		// 6. Check if it's goal
		if a.isGoal(new, end, step_size_mt) {
			// 7. Check if can be connected to goal
			isInObstacles, _, err := storage.IsLineInObstacles(new, end)
			if err != nil {
				return nil, 0.0, err
			}
			// If yes, connect to goal, found!!!
			if !isInObstacles {
				err := storage.AddWaypointWithPrevious(new, end)
				if err != nil {
					return nil, 0.0, err
				}
				goal_found = true
				fmt.Printf("Goal found at iteration %d/%d.\n\n", current_iter, max_iterations)
				break
			}
		}

	}

	// fmt.Printf("Storage at the end has %d constraints and %d waypoints.\n\n", storage.ConstraintsLen(), storage.WaypointsLen())

	// ------------------------------------------------------------------------------------------------------

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

func (a *RRTAlgorithm) GetParameters(parameters map[string]any, goal *models.Waypoint) (utils.Sampler, int, float64, float64) {
	// TODO: Take parameters from the actual map, and use defaults if not found
	fmt.Printf("parameters: %+v\n", parameters)
	MAX_ITERATIONS := int(utils.GetOrDefault(parameters, "max_iterations", 100000.0))
	GOAL_BIAS := utils.GetOrDefault(parameters, "goal_bias", 0.10)
	STEP_SIZE_MT := utils.GetOrDefault(parameters, "step_size_mt", 20.0)
	SAMPLER_TYPE := utils.GetOrDefault(parameters, "sampler_type", models.Uniform)
	SEED := utils.GetOrDefault(parameters, "seed", 945)

	// use sampler_type and seed
	base_sampler, err := utils.NewSampler(SAMPLER_TYPE, int64(SEED))
	if err != nil {
		// TODO: HANDLE THIS
	}

	SAMPLER := utils.NewGoalBiasSampler(
		base_sampler,
		goal,
		GOAL_BIAS,
		int64(SEED),
	)

	fmt.Printf("PARAMETERS\n")
	fmt.Printf("max_iterations: %d\n", MAX_ITERATIONS)
	fmt.Printf("step_size_mt: %f\n", STEP_SIZE_MT)
	fmt.Printf("goal_bias: %f\n", GOAL_BIAS)
	fmt.Printf("sampler: %+v\n", SAMPLER)
	fmt.Printf("--------------------------------------------------------\n")

	return SAMPLER, MAX_ITERATIONS, STEP_SIZE_MT, GOAL_BIAS
}

func (a *RRTAlgorithm) isGoal(w, goal *models.Waypoint, tolerance_mt float64) bool {
	// is goal if distance is less than tolerance_mt
	return utils.HaversineDistance3D(w, goal) < tolerance_mt
}