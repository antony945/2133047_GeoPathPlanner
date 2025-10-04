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

// TODO: Implement AntPath Algorithm
func (a *AntPathAlgorithm) Compute(searchVolume *models.Feature3D, waypoints []*models.Waypoint, constraints []*models.Feature3D, parameters map[string]any, storage storage.Storage) ([]*models.Waypoint, float64, error) {	
	// Create empty list of wps
	route := make([]*models.Waypoint, 0)
	cost := 0.0

	// 0. Load first wp
	route = append(route, waypoints[0])
	
	// 1. Load constraint into storage

	// 2. For each pair of wp -> run antpath (with )
	for i := 0; i < len(waypoints)-1; i++ {
		tmpRoute, tmpCost, err := a.Run(waypoints[i], waypoints[i+1], constraints, parameters, storage)
		if err != nil {
			// Return route until now
			return route, cost, fmt.Errorf("interrupted antpath for error between wp[%d] and wp[%d]: %w", i, i+1, err)
		}
		// Append new route but removing the first one
		route = append(route, tmpRoute[1:]...)
		cost += tmpCost
	}
	
	// 3. Return everything
	return route, cost, nil
}

func (a *AntPathAlgorithm) Run(start, end *models.Waypoint, constraints []*models.Feature3D, parameters map[string]any, storage storage.Storage) ([]*models.Waypoint, float64, error) {
	route := make([]*models.Waypoint, 0)
	// cost := 0.0

	// TODO: Buffer constraints
	// TODO: Union constraints

	// Add constraints to storage
	err := storage.AddConstraints(constraints)
	if err != nil {
		return nil, 0.0, err
	}

	// Get intersection points
	intersectionPoints, err := storage.GetIntersectionPoints(start, end)
	if err != nil {
		return nil, 0.0, err
	}

	route = append(route, start)
	
	// For every intersectionPoint struct get best way to go around obstacle
	for _, ip := range intersectionPoints {
		bestForPolygonWay := utils.GetBestWayToGoAroundPolygon(ip.Polygon, ip.EnteringPoint, ip.ExitingPoint)
		route = append(route, bestForPolygonWay...)
	}

	route = append(route, end)

	cost := utils.TotalHaversineDistance(route)
	return route, cost, nil
}