package storage

import (
	"fmt"
	"geopathplanner/routing/internal/models"
)

// Define common interface for storing and querying geospatial data
type Storage interface {
	AddWaypoint(w *models.Waypoint) error
	AddWaypoints(w_list []*models.Waypoint) error
	AddConstraint(c *models.Feature3D) error
	AddConstraints(c_list []*models.Feature3D) error
	WaypointsLen() int
	ConstraintsLen() int
	Clear() error // Clear temporary data for a request

	NearestPoint(p *models.Waypoint) (*models.Waypoint, error)
    KNearestPoints(p *models.Waypoint, k int) ([]*models.Waypoint, error)
	NearestPointsInRadius(p *models.Waypoint, radius_mt float64) ([]*models.Waypoint, error)
    IsPointInObstacles(p *models.Waypoint) (bool, *models.Feature3D, error)
    IsLineInObstacles(p1, p2 *models.Waypoint) (bool, []*models.Waypoint, error)
	GetIntersectionPoints(p1, p2 *models.Waypoint) ([]*models.LinePolygonIntersection, error)
	GetAllObstaclesContainingPoint(p *models.Waypoint) ([]*models.Feature3D, error)
}

func NewEmptyStorage(storageType models.StorageType) (Storage, error) {
	switch storageType {
		case models.Memory:
			return NewEmptyMemoryStorage()
		case models.Redis:
			return nil, fmt.Errorf("storage currently not implemented: %s", storageType)
		default:
			return nil, fmt.Errorf("storage not recognized: %s", storageType)
	}
}

func NewStorage(w_list []*models.Waypoint, c_list []*models.Feature3D, storageType models.StorageType) (Storage, error) {
	switch storageType {
		case models.Memory:
			return NewMemoryStorage(w_list, c_list)
		case models.Redis:
			return nil, fmt.Errorf("storage currently not implemented: %s", storageType)
		default:
			return nil, fmt.Errorf("storage not recognized: %s", storageType)
	}
}