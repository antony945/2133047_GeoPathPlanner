package storage

import (
	"errors"
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/utils"
	"slices"
)

// MemoryStorage stores everything in-memory (cleared after each request)
type MemoryStorage struct {
	DefaultStorage
	waypoints  []*models.Waypoint
	constraints []*models.Constraint
}

func NewEmptyMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

func NewMemoryStorage(w_list []*models.Waypoint, c_list []*models.Constraint) *MemoryStorage {
	ms := NewEmptyMemoryStorage()
	ms.AddWaypoints(w_list)
	ms.AddConstraints(c_list)
	return ms
}

func (m *MemoryStorage) Clear() error {
	// m.mu.Lock()
	// defer m.mu.Unlock()

	// TODO: Think about putting them to nil
	m.waypoints = []*models.Waypoint{}
	m.constraints = []*models.Constraint{}
	return nil
}

func (m *MemoryStorage) AddWaypoint(w *models.Waypoint) error {
	// m.mu.Lock()
	// defer m.mu.Unlock()
	m.waypoints = append(m.waypoints, w)
	return nil
}

func (m *MemoryStorage) GetWaypoints() ([]*models.Waypoint, error) {
	// m.mu.Lock()
	// defer m.mu.Unlock()
	return append([]*models.Waypoint(nil), m.waypoints...), nil
}

func (m *MemoryStorage) AddConstraint(c *models.Constraint) error {
	// m.mu.Lock()
	// defer m.mu.Unlock()
	m.constraints = append(m.constraints, c)
	return nil
}

func (m *MemoryStorage) GetConstraints() ([]*models.Constraint, error) {
	// m.mu.Lock()
	// defer m.mu.Unlock()
	return append([]*models.Constraint(nil), m.constraints...), nil
}

// ================================================================= Geometric helpers

// Scan list of obstacle until you find someone that intersect
// O(N)
func (m *MemoryStorage) NearestPoint(p *models.Waypoint) (*models.Waypoint, error) {
	if len(m.waypoints) == 0 {
		return nil, errors.New("no waypoints")
	}

	var nearest *models.Waypoint
	minDist := float64(-1)

	for _, wp := range m.waypoints {
		dist := utils.HaversineDistance3D(p, wp)
		if minDist < 0 || dist < minDist {
			minDist = dist
			nearest = wp
		}
	}

	return nearest, nil
}

// Compute distance from p to every point in the list and then sort based on that distance, retaining only the first k.
// O(N*logN)
func (m *MemoryStorage) KNearestPoints(p *models.Waypoint, k int) ([]*models.Waypoint, error) {
	if len(m.waypoints) == 0 {
		return nil, errors.New("no waypoints")
	}
	if k <= 0 {
		return nil, errors.New("k must be positive")
	}
	if k > len(m.waypoints) {
		k = len(m.waypoints)
	}

	// Create a slice of waypoints with distances
	type pointDist struct {
		point    *models.Waypoint
		distance float64
	}
	
	points := make([]pointDist, len(m.waypoints))
	for i, wp := range m.waypoints {
		points[i] = pointDist{
			point:    wp,
			distance: utils.HaversineDistance3D(p, wp),
		}
	}

	// Sort by ascending distance
	slices.SortFunc(points, func(p1, p2 pointDist) int {
		return int(p1.distance - p2.distance)
	})

	// Return k nearest points
	result := make([]*models.Waypoint, k)
	for i := 0; i < k; i++ {
		result[i] = points[i].point
	}
	return result, nil
}

// Compute distance from p to every point in the list discarding everyone that is more distant than radius. Then sort based on that distance.
// O(N*logN)
func (m *MemoryStorage) NearestPointsInRadius(p *models.Waypoint, radius float64) ([]*models.Waypoint, error) {
	if len(m.waypoints) == 0 {
		return nil, errors.New("no waypoints")
	}
	if radius <= 0 {
		return nil, errors.New("radius must be positive")
	}

	// Create a slice of waypoints with distances
	type pointDist struct {
		point    *models.Waypoint
		distance float64
	}
	
	points := make([]pointDist, 0, len(m.waypoints))
	for _, wp := range m.waypoints {
		if distance := utils.HaversineDistance3D(p, wp); distance < radius {
			points = append(points, pointDist{
				point:    wp,
				distance: utils.HaversineDistance3D(p, wp),
			})
		}
	}

	// Sort by ascending distance
	slices.SortFunc(points, func(p1, p2 pointDist) int {
		return int(p1.distance - p1.distance)
	})

	// Return all nearest points
	result := make([]*models.Waypoint, len(points))
	for i := 0; i < len(result); i++ {
		result[i] = points[i].point
	}
	return result, nil
}

// Scan list of obstacle until you find someone that intersect
// O(#obstacles)
func (m *MemoryStorage) IsPointInObstacle(p *models.Waypoint) (bool, error) {
	for _, obstacle := range m.constraints {
		if utils.PointInPolygon(p, obstacle) {
			return true, nil
		}		
	}
	return false, nil
}

// Scan list of obstacle until you find someone that intersect line
// O(#obstacles)
func (m *MemoryStorage) IsLineInObstacle(p1, p2 *models.Waypoint) (bool, error) {
	for _, obstacle := range m.constraints {
		if utils.DefaultLineInPolygon(p1, p2, obstacle) {
			return true, nil
		}
	}
	return false, nil
}