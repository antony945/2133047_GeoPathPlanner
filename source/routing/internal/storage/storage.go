package storage

import (
	"geopathplanner/routing/internal/models"
	"sync"
)

// Define common interface for storing and querying geospatial data
type Storage interface {
	AddWaypoint(w *models.Waypoint) error
	AddConstraint(c *models.Constraint) error
	Clear() error // Clear temporary data for a request

	NearestPoint(p *models.Waypoint) (*models.Waypoint, error)
    KNearestPoints(p *models.Waypoint, k int) ([]*models.Waypoint, error)
	NearestPointsInRadius(p *models.Waypoint, radius_mt float64) ([]*models.Waypoint, error)
    IsPointInObstacle(p *models.Waypoint) (bool, error)
    IsLineInObstacle(p1, p2 *models.Waypoint) (bool, error)
}

// DefaultStorage must implement Storage, so it can define a default for some methods
type DefaultStorage struct {
	mu         sync.Mutex
}

func (ds *DefaultStorage) AddWaypoint(w *models.Waypoint) error {
	return nil
}

func (ds *DefaultStorage) AddConstraint(c *models.Constraint) error {
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

func (d *DefaultStorage) AddConstraints(c_list []*models.Constraint) error {
	for _, c := range c_list {
		if err := d.AddConstraint(c); err != nil {
			return err
		}		
	}
	return nil
}