package algorithm

import "geopathplanner/routing/internal/models"

type job struct {
	i       int
	startWP *models.Waypoint
	endWP   *models.Waypoint
}

type result struct {
	i     int
	route []*models.Waypoint
	cost  float64
	err   error
}