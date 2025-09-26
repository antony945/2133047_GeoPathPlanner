package storage

import (
	"errors"
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/utils"
	"slices"
)

// MemoryStorage stores everything in-memory (cleared after each request)
type MemoryStorage struct {
	waypoints  []*models.Waypoint
	constraints []*models.Feature3D
}

func NewEmptyMemoryStorage() (*MemoryStorage, error) {
	return &MemoryStorage{}, nil
}

func NewMemoryStorage(w_list []*models.Waypoint, c_list []*models.Feature3D) (*MemoryStorage, error) {
	ms, err := NewEmptyMemoryStorage()
	if err != nil {
		return nil, err
	}

	if err := ms.AddWaypoints(w_list); err != nil {
		return nil, err
	}

	if err := ms.AddConstraints(c_list); err != nil {
		return nil, err
	}
	
	return ms, nil
}

func (m *MemoryStorage) Clear() error {
	// m.mu.Lock()
	// defer m.mu.Unlock()

	// TODO: Think about putting them to nil
	m.waypoints = []*models.Waypoint{}
	m.constraints = []*models.Feature3D{}
	return nil
}

func (m *MemoryStorage) WaypointsLen() int {
	return len(m.waypoints)
}

func (m *MemoryStorage) ConstraintsLen() int {
	return len(m.constraints)
}

func (m *MemoryStorage) AddWaypoint(w *models.Waypoint) error {
	// m.mu.Lock()
	// defer m.mu.Unlock()
	m.waypoints = append(m.waypoints, w)
	return nil
}

func (m *MemoryStorage) AddWaypoints(w_list []*models.Waypoint) error {
	for _, w := range w_list {
		if err := m.AddWaypoint(w); err != nil {
			return err
		}
	}
	return nil
}

func (m *MemoryStorage) GetWaypoints() ([]*models.Waypoint, error) {
	// m.mu.Lock()
	// defer m.mu.Unlock()
	return m.waypoints, nil
}

func (m *MemoryStorage) AddConstraint(c *models.Feature3D) error {
	// m.mu.Lock()
	// defer m.mu.Unlock()
	m.constraints = append(m.constraints, c)
	return nil
}

func (m *MemoryStorage) AddConstraints(c_list []*models.Feature3D) error {
	for _, c := range c_list {
		if err := m.AddConstraint(c); err != nil {
			return err
		}		
	}
	return nil
}

func (m *MemoryStorage) GetConstraints() ([]*models.Feature3D, error) {
	// m.mu.Lock()
	// defer m.mu.Unlock()
	return m.constraints, nil
}

// ================================================================= Geometric helpers

// Scan list of obstacle until you find someone that intersect
// O(N)
func (m *MemoryStorage) NearestPoint(p *models.Waypoint) (*models.Waypoint, error) {
	// TODO: Just for visual debug
	delete(p.Feature.Properties, "parameter")
	delete(p.Feature.Properties, "nearest")
	p.Feature.Properties["parameter"] = true
	
	if len(m.waypoints) == 0 {
		return nil, errors.New("no waypoints")
	}

	var nearest *models.Waypoint
	minDist := float64(-1)

	for _, wp := range m.waypoints {
		// TODO: Just for visual debug
		delete(wp.Feature.Properties, "parameter")
		delete(wp.Feature.Properties, "nearest")
		
		dist := utils.HaversineDistance3D(*p, *wp)
		if minDist < 0 || dist < minDist {
			minDist = dist
			nearest = wp
		}
	}

	// TODO: Just for visual debug
	nearest.Feature.Properties["nearest"] = true
	return nearest, nil
}

// Compute distance from p to every point in the list and then sort based on that distance, retaining only the first k.
// O(N*logN)
func (m *MemoryStorage) KNearestPoints(p *models.Waypoint, k int) ([]*models.Waypoint, error) {
	// TODO: Just for visual debug
	delete(p.Feature.Properties, "parameter")
	delete(p.Feature.Properties, "near")
	delete(p.Feature.Properties, "nearest")
	p.Feature.Properties["parameter"] = true
	
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
		// TODO: Just for visual debug
		delete(wp.Feature.Properties, "parameter")
		delete(wp.Feature.Properties, "near")
		delete(wp.Feature.Properties, "nearest")
		
		points[i] = pointDist{
			point:    wp,
			distance: utils.HaversineDistance3D(*p, *wp),
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
		// TODO: Just for visual debug
		result[i].Feature.Properties["near"] = true
		if i == 0 {
			result[i].Feature.Properties["nearest"] = true
		}
	}

	return result, nil
}

// Compute distance from p to every point in the list discarding everyone that is more distant than radius. Then sort based on that distance.
// O(N*logN)
func (m *MemoryStorage) NearestPointsInRadius(p *models.Waypoint, radius float64) ([]*models.Waypoint, error) {
	// TODO: Just for visual debug
	delete(p.Feature.Properties, "parameter")
	delete(p.Feature.Properties, "near")
	delete(p.Feature.Properties, "nearest")
	p.Feature.Properties["parameter"] = true
	
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
		// TODO: Just for visual debug
		delete(wp.Feature.Properties, "parameter")
		delete(wp.Feature.Properties, "near")
		delete(wp.Feature.Properties, "nearest")

		if distance := utils.HaversineDistance3D(*p, *wp); distance < radius {
			points = append(points, pointDist{
				point:    wp,
				distance: distance,
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
		// TODO: Just for visual debug
		result[i].Feature.Properties["near"] = true
		if i == 0 {
			result[i].Feature.Properties["nearest"] = true
		}
	}

	return result, nil
}

// Scan list of obstacle until you find someone that intersect
// O(#obstacles)
func (m *MemoryStorage) IsPointInObstacles(p *models.Waypoint) (bool, error) {
	for _, obstacle := range m.constraints {
		if utils.PointInPolygon(p, obstacle) {
			return true, nil
		}		
	}
	
	return false, nil
}

// Scan list of obstacle until you find someone that intersect line
// O(#obstacles)
func (m *MemoryStorage) IsLineInObstacles(p1, p2 *models.Waypoint) (bool, []*models.Waypoint, error) {
	in, line := utils.LineInPolygon(p1, p2, m.constraints...)
	return in, line, nil
}