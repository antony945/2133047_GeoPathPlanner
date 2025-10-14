package algorithm

import (
	"fmt"
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/storage"
	"geopathplanner/routing/internal/utils"
	"math"
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

func (a *RRTStarAlgorithm) Compute(searchVolume *models.Feature3D, waypoints []*models.Waypoint, constraints []*models.Feature3D, parameters map[string]any, storage storage.Storage) ([]*models.Waypoint, float64, error) {
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
			return route, cost, fmt.Errorf("interrupted RRTStar for error between wp[%d] and wp[%d]: %w", i, i+1, err)
		}
		// Append new route but removing the first one
		route = append(route, tmpRoute[1:]...)
		cost += tmpCost
	}
	
	// 3. Return everything
	return route, cost, nil
}

func (a *RRTStarAlgorithm) Run(searchVolume *models.Feature3D, start, end *models.Waypoint, parameters map[string]any, storage storage.Storage) ([]*models.Waypoint, float64, error) {
	// TODO: Think if this is the correct place
	storage.ClearWaypoints()
	fmt.Printf("Storage has %d constraints and %d waypoints.\n\n", storage.ConstraintsLen(), storage.WaypointsLen())

	// TODO: Parameters
	MAX_ITERATIONS := 10000
	GOAL_BIAS := 0.10
	SAMPLER := utils.NewGoalBiasSampler(
		utils.NewUniformSampler(),
		// utils.NewUniformSamplerWithSeed(10),
		// utils.NewHaltonSampler(),
		end,
		GOAL_BIAS,
	)
	STEP_SIZE_MT := 10.0
	K_INIT := 2*math.E
	R_INIT_MT := 2.0

	fmt.Printf("PARAMETERS\n")
	fmt.Printf("max_iterations: %d\n", MAX_ITERATIONS)
	fmt.Printf("step_size_mt: %f\n", STEP_SIZE_MT)
	fmt.Printf("goal_bias: %f\n", GOAL_BIAS)
	fmt.Printf("sampler: %+v\n", SAMPLER)
	fmt.Printf("k_init: %d\n", int(K_INIT))
	fmt.Printf("r_init_mt: %f\n", R_INIT_MT)
	fmt.Printf("--------------------------------------------------------\n")

	// Add start to storage
	err := storage.AddWaypointWithPrevious(nil, start)
	if err != nil {
		return nil, 0.0, err
	}
	goal_found := false

	// ------------------------------------------------------------------------------------------------------

	for current_iter := range MAX_ITERATIONS {
		// Change K according to cardinality of V (no. of nodes)
		K := int(K_INIT * math.Log(float64(storage.WaypointsLen())))+1
		// R := math.Max(R_INIT_MT * math.Sqrt(math.Log(float64(storage.WaypointsLen()))/float64(storage.WaypointsLen())), STEP_SIZE_MT)

		if current_iter % 100 == 0 {
			if goal_found {
				route, _ := storage.GetPathToRoot(end)
				cost_km := utils.TotalHaversineDistance(route)/1000
				fmt.Printf("[%d/%d] k: %d, #wps: %d, #routeWps: %d, routeCost: %.3f km\n", current_iter, MAX_ITERATIONS, K, storage.WaypointsLen(), len(route), cost_km)
				// fmt.Printf("[%d/%d] radius: %.2fmt, #wps: %d, cost: %.3f km\n", current_iter, MAX_ITERATIONS, R, len(route), cost_km)
			
			} else {
				fmt.Printf("[%d/%d] k: %d, #wps: %d, goal not found yet\n", current_iter, MAX_ITERATIONS, K, storage.WaypointsLen())
				// fmt.Printf("[%d/%d] radius: %.2fmt, goal not found yet\n", current_iter, MAX_ITERATIONS, R)
			}
		}

		// 1. Sample a new free wp
		sampled, err := storage.SampleFree(SAMPLER, searchVolume, start.Alt)
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
		// 	cost_km := utils.TotalHaversineDistance(route)/1000
		// 	fmt.Printf("Rewired Tree\n#wps: %d, cost: %.3f km\n", len(route), cost_km)
		// }

		// 6. Check if it's goal
		if !goal_found && a.isGoal(new, end, STEP_SIZE_MT) {
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
				cost_km := utils.TotalHaversineDistance(route)/1000
				fmt.Printf("New goal found at iteration %d/%d.\n", current_iter, MAX_ITERATIONS)
				fmt.Printf("#wps: %d, cost: %.3f km\n", len(route), cost_km)
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
		return nil, 0.0, fmt.Errorf("goal not found with %d iterations", MAX_ITERATIONS)
	}
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