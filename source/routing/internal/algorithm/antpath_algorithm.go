package algorithm

import (
	"fmt"
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/storage"
	"geopathplanner/routing/internal/utils"
)

type AntPathAlgorithm struct {}

func NewAntPathAlgorithm() (*AntPathAlgorithm, error) {
	// TODO: To implement
	return &AntPathAlgorithm{}, nil
}

// Implement AntPath Algorithm
func (a *AntPathAlgorithm) Compute(searchVolume *models.Feature3D, waypoints []*models.Waypoint, constraints []*models.Feature3D, parameters map[string]any, storageType models.StorageType) ([]*models.Waypoint, float64, error) {	
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
