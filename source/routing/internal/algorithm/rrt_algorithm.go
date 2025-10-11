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
	// TODO: To implement
	return &RRTAlgorithm{}, nil
}

// TODO: Implement RRT Algorithm
func (a *RRTAlgorithm) Compute(searchVolume *models.Feature3D, waypoints []*models.Waypoint, constraints []*models.Feature3D, parameters map[string]any, storage storage.Storage) ([]*models.Waypoint, float64, error) {
	// Create empty list of wps
	route := make([]*models.Waypoint, 0)
	cost := 0.0

	// 0. Load first wp
	route = append(route, waypoints[0])
	
	// TODO: Buffer constraints
	// 1. Load constraint into storage
	err := storage.AddConstraints(constraints)
	if err != nil {
		return nil, 0.0, err
	}

	// 2. For each pair of wp -> run rrt
	for i := 0; i < len(waypoints)-1; i++ {
		tmpRoute, tmpCost, err := a.Run(searchVolume, waypoints[i], waypoints[i+1], parameters, storage)
		if err != nil {
			// Return route until now
			return route, cost, fmt.Errorf("interrupted RRT for error between wp[%d] and wp[%d]: %w", i, i+1, err)
		}
		// Append new route but removing the first one
		route = append(route, tmpRoute[1:]...)
		cost += tmpCost
	}
	
	// 3. Return everything
	return route, cost, nil
}

func (a *RRTAlgorithm) Run(searchVolume *models.Feature3D, start, end *models.Waypoint, parameters map[string]any, storage storage.Storage) ([]*models.Waypoint, float64, error) {
	// HERE IMPLEMENT RRT
	// start from start
	// sample a new free wp
	// get nearest wp
	// check if connection from nearest to free it's possible
	// if yes connect them, if no sample new one
	// check if it's goal (approximately)
	// check if connection from free to goal is possible
	// if yes break from the loop, if not continue

	// TODO: Think if this is the correct place
	storage.ClearWaypoints()
	fmt.Printf("Storage has %d constraints and %d waypoints.\n\n", storage.ConstraintsLen(), storage.WaypointsLen())

	// TODO: Parameters
	MAX_ITERATIONS := 10000
	GOAL_BIAS := 0.10
	SAMPLER := utils.NewGoalBiasSampler(
		utils.NewUniformSamplerWithSeed(10),
		end,
		GOAL_BIAS,
	)
	STEP_SIZE_MT := 10.0

	// Add start to storage
	err := storage.AddWaypointWithPrevious(nil, start)
	if err != nil {
		return nil, 0.0, err
	}
	
	goal_found := false
	for current_iter := range MAX_ITERATIONS {
		// 1. Sample a new free wp
		sampled, err := storage.SampleFree(SAMPLER, searchVolume, start.Alt)
		if err != nil {
			// Impossible to sample 
			return nil, 0.0, err
		}

		// 2. Get nearest wp
		nearest, err := storage.NearestPoint(sampled)
		if err != nil {
			return nil, 0.0, err
		}

		// 3. Find wp starting from nearest in direction of sampled at distance STEP_SIZE_MT
		new := utils.GetPointInDirectionAtDistance(nearest, sampled, STEP_SIZE_MT)

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
		if a.isGoal(new, end, STEP_SIZE_MT) {
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
				fmt.Printf("Goal found at iteration %d/%d.\n\n", current_iter, MAX_ITERATIONS)
				break
			}
		}

	}

	// fmt.Printf("Storage at the end has %d constraints and %d waypoints.\n\n", storage.ConstraintsLen(), storage.WaypointsLen())

	// Now try to obtain the route from the goal node
	if goal_found {
		route, err := storage.GetPathToRoot(end)
		if err != nil {
			return nil, 0.0, err
		}
		cost := utils.TotalHaversineDistance(route)
		return route, cost, nil
	} else {
		return nil, 0.0, fmt.Errorf("goal not found with %d in iterations", MAX_ITERATIONS)
	}
}

func (a *RRTAlgorithm) isGoal(w, goal *models.Waypoint, tolerance_mt float64) bool {
	// is goal if distance is less than tolerance_mt
	return utils.HaversineDistance3D(*w, *goal) < tolerance_mt
}