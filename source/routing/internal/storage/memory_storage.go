package storage

import (
	"errors"
	"fmt"
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/utils"
	"slices"
)

// MemoryStorage stores everything in-memory (cleared after each request)
type MemoryStorage struct {
	waypoints  []*models.Waypoint
	constraints []*models.Feature3D
	waypointsMap map[*models.Waypoint]*models.PointDist
}

func NewEmptyMemoryStorage() (*MemoryStorage, error) {
	return &MemoryStorage{
		waypoints: make([]*models.Waypoint, 0),
		constraints: make([]*models.Feature3D, 0),
		waypointsMap: make(map[*models.Waypoint]*models.PointDist),
	}, nil
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
	m.constraints = make([]*models.Feature3D, 0)
	m.ClearWaypoints()
	return nil
}

func (m *MemoryStorage) ClearWaypoints() error {
	m.waypoints = make([]*models.Waypoint, 0)
	m.waypointsMap = make(map[*models.Waypoint]*models.PointDist)
	return nil
}

func (m *MemoryStorage) WaypointsLen() int {
	return len(m.waypointsMap)
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
	if w_list == nil {
		return nil
	}
	
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

func (m *MemoryStorage) MustGetWaypoints() []*models.Waypoint {
	wps, err := m.GetWaypoints()
	if err != nil {
		panic(err)
	}
	return wps
}

func (m *MemoryStorage) AddConstraint(c *models.Feature3D) error {
	// m.mu.Lock()
	// defer m.mu.Unlock()
	m.constraints = append(m.constraints, c)
	return nil
}

func (m *MemoryStorage) AddConstraints(c_list []*models.Feature3D) error {
	if c_list == nil {
		return nil
	}
	
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

func (m *MemoryStorage) MustGetConstraints() []*models.Feature3D {
	c, err := m.GetConstraints()
	if err != nil {
		panic(err)
	}
	return c
}

// ================================================================= RRT

func (m *MemoryStorage) AddWaypointWithPrevious(prev *models.Waypoint, w *models.Waypoint) error {
	// Add waypoint to the list of waypoints
	m.AddWaypoint(w)

	// But also add it in the map
	m.ChangePrevious(prev, w)
	return nil
}

func (m *MemoryStorage)	ChangePrevious(new_prev *models.Waypoint, w *models.Waypoint) error {
	distance := 0.0
	if new_prev != nil {
		distance = utils.HaversineDistance3D(new_prev, w)
	}
	
	// Just add it to the map
	prev, ok := m.waypointsMap[w]
	if !ok {
		// w never added to the map, add it
		m.waypointsMap[w] = models.NewPointDist(new_prev, distance)
	} else {
		// w already added
		prev.Point = new_prev
		prev.Distance = distance
		m.waypointsMap[w] = prev
	}
	
	return nil
}

func (m *MemoryStorage) GetPrevious(w *models.Waypoint) (*models.Waypoint, error) {
	// Search in the map
	return m.waypointsMap[w].Point, nil
}

func (m *MemoryStorage) GetPathToRoot(w *models.Waypoint) ([]*models.Waypoint, error) {
	// Start from w and find its previous until there are no more previous
	route := make([]*models.Waypoint, 0)
	current := w
	for {
		// Check if current is nil, in case return
		if current == nil {
			// End of route
			// Reverse route and return it
			left := 0
			right := len(route) - 1
			for left < right {
				route[left], route[right] = route[right], route[left]
				left++
				right--
			}
			return route, nil
		}

		route = append(route, current)

		// Check if prev was already in the map, if no error
		prev, ok := m.waypointsMap[current]
		if !ok {
    		return nil, fmt.Errorf("waypoint %+v not found in map", w)
		}

		current = prev.Point
	}
}

func (m *MemoryStorage) GetCostToRoot(w *models.Waypoint) (float64, error) {
	// Iteratively do GetPrevious and search for the cost
	current := w
	cost := 0.0
	for {
		if current == nil {
			// Reach root
			break
		}

		// Add distance
		prev, ok := m.waypointsMap[current]
		if !ok {
    		return 0.0, fmt.Errorf("waypoint %+v not found in map", w)
		}

		cost += prev.Distance
		current = prev.Point
	}

	return cost, nil
}

// ================================================================= Geometric helpers

// Scan list of obstacle until you find someone that intersect
// O(N)
func (m *MemoryStorage) NearestPoint(p *models.Waypoint) (*models.Waypoint, float64, error) {
	// TODO: Just for visual debug
	delete(p.Feature.Properties, "parameter")
	delete(p.Feature.Properties, "nearest")
	p.Feature.Properties["parameter"] = true
	
	if m.WaypointsLen() == 0 {
		return nil, 0.0, errors.New("no waypoints")
	}

	var nearest *models.Waypoint
	minDist := float64(-1)

	for _, wp := range m.MustGetWaypoints() {
		// TODO: Just for visual debug
		delete(wp.Feature.Properties, "parameter")
		delete(wp.Feature.Properties, "nearest")
		
		dist := utils.HaversineDistance3D(p, wp)
		if minDist < 0 || dist < minDist {
			minDist = dist
			nearest = wp
		}
	}

	// TODO: Just for visual debug
	nearest.Feature.Properties["nearest"] = true
	return nearest, minDist, nil
}

// Compute distance from p to every point in the list and then sort based on that distance, retaining only the first k.
// O(N*logN)
func (m *MemoryStorage) KNearestPoints(p *models.Waypoint, k int) ([]*models.Waypoint, []float64, error) {
	// TODO: Just for visual debug
	delete(p.Feature.Properties, "parameter")
	delete(p.Feature.Properties, "near")
	delete(p.Feature.Properties, "nearest")
	p.Feature.Properties["parameter"] = true
	
	wpsLen := m.WaypointsLen()
	if wpsLen == 0 {
		return nil, nil, errors.New("no waypoints")
	}
	if k < 0 {
		return nil, nil, errors.New("k must be positive")
	}
	if k == 0 {
		return []*models.Waypoint{}, []float64{0.0}, nil
	}
	if k > wpsLen {
		k = wpsLen
	}

	// Create a slice of waypoints with distances
	points_dist := make(map[*models.Waypoint]float64, wpsLen)
	
	for _, wp := range m.MustGetWaypoints() {
		// TODO: Just for visual debug
		delete(wp.Feature.Properties, "parameter")
		delete(wp.Feature.Properties, "near")
		delete(wp.Feature.Properties, "nearest")
		
		points_dist[wp] = utils.HaversineDistance3D(p, wp)
	}

	// Sort by ascending distance
	points := make([]*models.Waypoint, 0, len(points_dist))
	for wp := range points_dist {
		points = append(points, wp)
	}

	slices.SortFunc(points, func(w1, w2 *models.Waypoint) int {
		return int(points_dist[w1] - points_dist[w2])
	})

	// Return k nearest points
	result := make([]*models.Waypoint, k)
	for i := 0; i < k; i++ {
		result[i] = points[i]
		// TODO: Just for visual debug
		result[i].Feature.Properties["near"] = true
		if i == 0 {
			result[i].Feature.Properties["nearest"] = true
		}
	}

	// Return distances
	distances := make([]float64, 0, k)
	for _, near := range result {
		distances = append(distances, points_dist[near])
	}

	return result, distances, nil
}

// Compute distance from p to every point in the list discarding everyone that is more distant than radius. Then sort based on that distance.
// O(N*logN)
func (m *MemoryStorage) NearestPointsInRadius(p *models.Waypoint, radius float64) ([]*models.Waypoint, []float64, error) {
	// TODO: Just for visual debug
	delete(p.Feature.Properties, "parameter")
	delete(p.Feature.Properties, "near")
	delete(p.Feature.Properties, "nearest")
	p.Feature.Properties["parameter"] = true
	
	wpsLen := m.WaypointsLen()
	if wpsLen == 0 {
		return nil, nil, errors.New("no waypoints")
	}
	if radius < 0 {
		return nil, nil, errors.New("radius must be positive")
	}
	if radius == 0 {
		return []*models.Waypoint{}, []float64{0.0}, nil
	}

	// Create a slice of waypoints with distances
	points_dist := make(map[*models.Waypoint]float64, wpsLen)

	for _, wp := range m.MustGetWaypoints() {
		// TODO: Just for visual debug
		delete(wp.Feature.Properties, "parameter")
		delete(wp.Feature.Properties, "near")
		delete(wp.Feature.Properties, "nearest")

		if distance := utils.HaversineDistance3D(p, wp); distance < radius {
			points_dist[wp] = distance
		}
	}

	// Sort by ascending distance
	points := make([]*models.Waypoint, 0, len(points_dist))
	for wp := range points_dist {
		points = append(points, wp)
	}

	slices.SortFunc(points, func(w1, w2 *models.Waypoint) int {
		return int(points_dist[w1] - points_dist[w2])
	})

	// Return all nearest points
	result := make([]*models.Waypoint, 0, len(points))
	for i := 0; i < len(result); i++ {
		result = append(result, points[i]) 

		// TODO: Just for visual debug
		result[i].Feature.Properties["near"] = true
		if i == 0 {
			result[i].Feature.Properties["nearest"] = true
		}
	}

	// Return distances
	distances := make([]float64, 0, len(points))
	for _, near := range result {
		distances = append(distances, points_dist[near])
	}

	return result, distances, nil
}

// Scan list of obstacle until you find someone that intersect
// O(#obstacles)
func (m *MemoryStorage) IsPointInObstacles(p *models.Waypoint) (bool, *models.Feature3D, error) {
	for _, obstacle := range m.constraints {
		if utils.PointInPolygon(p, obstacle) {
			return true, obstacle, nil
		}
	}
	
	return false, nil, nil
}

// Scan list of obstacle and return every obstacle that collide with point
// O(#obstacles)
func (m *MemoryStorage) GetAllObstaclesContainingPoint(p *models.Waypoint) ([]*models.Feature3D, error) {
	obstacles := make([]*models.Feature3D, 0)
	
	for _, obstacle := range m.constraints {
		if utils.PointInPolygon(p, obstacle) {
			obstacles = append(obstacles, obstacle)
		}
	}

	return obstacles, nil
}

// Scan list of obstacle until you find someone that intersect line
// O(#obstacles)
func (m *MemoryStorage) IsLineInObstacles(p1, p2 *models.Waypoint) (bool, []*models.Waypoint, error) {
	// TODO: First check line bounds with polygon bounds
	in, line := utils.LineInPolygon(p1, p2, m.constraints...)
	return in, line, nil
}

// Get intersection ponits (useful for AntPath)
func (m *MemoryStorage)	GetIntersectionPoints(p1, p2 *models.Waypoint) ([]*models.LinePolygonIntersection, error) {
	// Divide line into point and then check if any individual point lies in polygon
	quantizedLine := utils.DefaultResampleLineToInterval(p1, p2)
	// For every point, check if it intersect with a polygon and mark it as:
	// Start Point: if previous point was not intersecting but current is
	// End Point: if previous point was intersecting but current is not
	// Return startpoint, endpoint and the polygon that the line is intersecting
	
	// If first is already inside -> error, impossible
	if firstInside, _, _ := m.IsPointInObstacles(quantizedLine[0]); firstInside {
		return nil, errors.New("first point is already in obstacles")
	}
	if lastInside, _, _ := m.IsPointInObstacles(quantizedLine[0]); lastInside {
		return nil, errors.New("last point is already in obstacles")
	}
	
	lpi_list := make([]*models.LinePolygonIntersection, 0)
	obstacles := make(models.PolygonSet)
	var startPoint, endPoint *models.Waypoint
	previousInside := false
	for i := 1; i < len(quantizedLine); i++ {
		// Check if current point is intersecting
		currentObstacles, _ := m.GetAllObstaclesContainingPoint(quantizedLine[i])
		// If there are obstacles add them to the ongoing list, otherwise currentObstacles will be empty
		obstacles.AddAll(currentObstacles...)
		currentInside := len(currentObstacles) > 0

		if (currentInside && !previousInside) {
			// Previous is startPoint
			startPoint = quantizedLine[i-1]
		} else if (!currentInside && previousInside) { // Current is endPoint
			endPoint = quantizedLine[i]
			
			// Once you found the endPoint create the LinePolygonIntersection
			lpi_list = append(lpi_list, models.NewLinePolygonIntersection(startPoint, endPoint, obstacles.Values()))

			// Reset everything to search for other obstacles
			obstacles.Clear()
		}

		previousInside = currentInside
	}

	return lpi_list, nil
}

func (m *MemoryStorage) Sample(sampler utils.Sampler, sampleVolume *models.Feature3D, alt models.Altitude) (*models.Waypoint, error) {	
	// TODO: Decide which to use
	sampled, err := utils.SampleWithAltitude2D(sampler, sampleVolume.Geometry, alt)
	if err != nil {
		return nil, fmt.Errorf("unexpected error during memoryStorage Sample: %w", err)
	}

	// TODO: Check if sampled was already present
	return sampled, nil
}

func (m *MemoryStorage) SampleFree(sampler utils.Sampler, sampleVolume *models.Feature3D, alt models.Altitude) (*models.Waypoint, error) {
	isInObstacle := true
	var sampled *models.Waypoint
	var err error
	// sample until you found a point that is not in an obstacle
	for isInObstacle {
		sampled, err = m.Sample(sampler, sampleVolume, alt)
		if err != nil {
			return nil, fmt.Errorf("unexpected error during memoryStorage SampleFree: %w", err)
		}
		isInObstacle, _, _ = m.IsPointInObstacles(sampled)
	}
	
	return sampled, nil
}