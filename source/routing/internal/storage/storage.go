package storage

import (
	"fmt"
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/utils"
)

// Define common interface for storing and querying geospatial data
type Storage interface {
	AddWaypoint(w *models.Waypoint) error
	AddWaypoints(w_list []*models.Waypoint) error
	GetWaypoints() ([]*models.Waypoint, error)
	MustGetWaypoints() ([]*models.Waypoint)

	AddConstraint(c *models.Feature3D) error
	AddConstraints(c_list []*models.Feature3D) error
	GetConstraints() ([]*models.Feature3D, error)
	MustGetConstraints() ([]*models.Feature3D)

	WaypointsLen() int
	
	ConstraintsLen() int
	
	Clear() error // Clear temporary data for a request
	ClearWaypoints() error
	ClearConstraints() error

	Clone() Storage // Clone storage

	AddWaypointWithPrevious(prev *models.Waypoint, w *models.Waypoint) error	
	ChangePrevious(new_prev *models.Waypoint, w *models.Waypoint) error
	GetPrevious(p *models.Waypoint) (*models.Waypoint, error)
	GetPathToRoot(w *models.Waypoint) ([]*models.Waypoint, error)
	GetCostToRoot(w *models.Waypoint) (float64, error)

	NearestPoint(p *models.Waypoint) (*models.Waypoint, float64, error)
    KNearestPoints(p *models.Waypoint, k int) ([]*models.Waypoint, []float64, error)
	NearestPointsInRadius(p *models.Waypoint, radius_mt float64) ([]*models.Waypoint, []float64, error)
    
	IsPointInObstacles(p *models.Waypoint) (bool, *models.Feature3D, error)
    IsLineInObstacles(p1, p2 *models.Waypoint) (bool, []*models.Waypoint, error)
	
	GetIntersectionPoints(p1, p2 *models.Waypoint) ([]*models.LinePolygonIntersection, error)
	GetAllObstaclesContainingPoint(p *models.Waypoint) ([]*models.Feature3D, error)
	
	Sample(sampler utils.Sampler, sampleVolume *models.Feature3D, alt models.Altitude) (*models.Waypoint, error)
	SampleFree(sampler utils.Sampler, sampleVolume *models.Feature3D, alt models.Altitude) (*models.Waypoint, error)
}

func NewEmptyStorage(storageType models.StorageType) (Storage, error) {
	switch storageType {
		case models.Memory:
			return NewEmptyMemoryStorage()
		case models.RTree:
			return NewEmptyRTreeStorage()
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
		case models.RTree:
			return NewRTreeStorage(w_list, c_list)
		case models.Redis:
			return nil, fmt.Errorf("storage currently not implemented: %s", storageType)
		default:
			return nil, fmt.Errorf("storage not recognized: %s", storageType)
	}
}