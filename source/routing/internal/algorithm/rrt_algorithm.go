package algorithm

import (
	"fmt"
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/storage"
	"geopathplanner/routing/internal/utils"
)

type RRTAlgorithm struct {
}

func NewRRTAlgorithm() (*RRTAlgorithm, error) {
	return &RRTAlgorithm{}, nil
}

func (a *RRTAlgorithm) Compute(searchVolume *models.Feature3D, waypoints []*models.Waypoint, constraints []*models.Feature3D, parameters map[string]any, storageType models.StorageType) ([]*models.Waypoint, float64, error) {
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
	
	// HERE IMPLEMENT RRT
	// start from start
	// sample a new free wp
	// get nearest wp
	// check if connection from nearest to free it's possible
	// if yes connect them, if no sample new one
	// check if it's goal (approximately)
	// check if connection from free to goal is possible
	// if yes break from the loop, if not continue

	fmt.Printf("Storage has %d constraints and %d waypoints.\n\n", storage.ConstraintsLen(), storage.WaypointsLen())

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
	MAX_ITERATIONS := utils.GetOrDefault(parameters, "max_iterations", 20000)
	GOAL_BIAS := utils.GetOrDefault(parameters, "goal_bias", 0.10)
	STEP_SIZE_MT := utils.GetOrDefault(parameters, "goal_bias", 10.0)
	SAMPLER_TYPE := utils.GetOrDefault(parameters, "sampler_type", models.Uniform)
	SEED := utils.GetOrDefault(parameters, "seed", 945)

	// use sampler_type and sampler_seed
	base_sampler, err := utils.NewSampler(SAMPLER_TYPE, int64(SEED))
	if err != nil {
		// TODO: HANDLE THIS
	}

	SAMPLER := utils.NewGoalBiasSampler(
		base_sampler,
		goal,
		GOAL_BIAS,
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