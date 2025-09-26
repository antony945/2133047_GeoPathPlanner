package storage

import (
	"fmt"
	"geopathplanner/routing/internal/models"
	"sync"
)

// Define common interface for storing and querying geospatial data
type Storage interface {
	AddWaypoint(w *models.Waypoint) error
	AddConstraint(c *models.Feature3D) error
	Clear() error // Clear temporary data for a request

	NearestPoint(p *models.Waypoint) (*models.Waypoint, error)
    KNearestPoints(p *models.Waypoint, k int) ([]*models.Waypoint, error)
	NearestPointsInRadius(p *models.Waypoint, radius_mt float64) ([]*models.Waypoint, error)
    IsPointInObstacles(p *models.Waypoint) (bool, error)
    IsLineInObstacles(p1, p2 *models.Waypoint) (bool, error)
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

// DefaultStorage must implement Storage, so it can define a default for some methods
type DefaultStorage struct {
	mu         sync.Mutex
}

func (ds *DefaultStorage) AddWaypoint(w *models.Waypoint) error {
	return nil
}

func (ds *DefaultStorage) AddConstraint(c *models.Feature3D) error {
	return nil
}

func (ds *DefaultStorage) Clear(w *models.Waypoint) error {
	return nil
}

// ====================== Default methods

func (d *DefaultStorage) AddWaypoints(w_list []*models.Waypoint) error {
	for _, w := range w_list {
		if err := d.AddWaypoint(w); err != nil {
			return err
		}
	}
	return nil
}

func (d *DefaultStorage) AddConstraints(c_list []*models.Feature3D) error {
	for _, c := range c_list {
		if err := d.AddConstraint(c); err != nil {
			return err
		}		
	}
	return nil
}