package storage

import "geopathplanner/routing/internal/models"

// RedisStorage stores everything in redis database (cleared after each request)
type RedisStorage struct {
	waypoints  []*models.Waypoint
	constraints []*models.Feature3D
}