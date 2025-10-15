package storage

import (
	"errors"
	"fmt"
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/utils"

	"github.com/dhconnelly/rtreego"
)

const (
	SPATIAL_DIMENSION = 2
	MINIMUM_BRANCHING_FACTOR = 25
	MAXIMUM_BRANCHING_FACTOR = 50
)

type RTreeStorage struct {
	// TODO: Think about using a list in case retrieving all waypoints at once is an important operation
	// waypoints []*models.Waypoint
	waypointsTree *rtreego.Rtree
	constraintsTree *rtreego.Rtree
	// TODO: Think if waypoints map is necessary or we could insert nodes in the RTree to avoid memory duplication
	waypointsMap map[*models.Waypoint]*models.PointDist
}

// ---------------------------------------------------------------- CONSTRUCTORS

func NewEmptyRTreeStorage() (*RTreeStorage, error) {
	return &RTreeStorage{
		waypointsTree: rtreego.NewTree(SPATIAL_DIMENSION, MINIMUM_BRANCHING_FACTOR, MAXIMUM_BRANCHING_FACTOR),
		constraintsTree: rtreego.NewTree(SPATIAL_DIMENSION, MINIMUM_BRANCHING_FACTOR, MAXIMUM_BRANCHING_FACTOR),
	}, nil
}

func NewRTreeStorage(w_list []*models.Waypoint, c_list []*models.Feature3D) (*RTreeStorage, error) {
	rs, err := NewEmptyRTreeStorage()
	if err != nil {
		return nil, err
	}

	if err := rs.AddWaypoints(w_list); err != nil {
		return nil, err
	}

	if err := rs.AddConstraints(c_list); err != nil {
		return nil, err
	}
	
	return rs, nil
}

// ---------------------------------------------------------------- GENERAL

func (r *RTreeStorage) Clear() error {
	r.constraintsTree = rtreego.NewTree(SPATIAL_DIMENSION, MINIMUM_BRANCHING_FACTOR, MAXIMUM_BRANCHING_FACTOR)
	return r.ClearWaypoints()
}

func (r *RTreeStorage) ClearWaypoints() error {
	r.waypointsMap = make(map[*models.Waypoint]*models.PointDist)
	r.waypointsTree = rtreego.NewTree(SPATIAL_DIMENSION, MINIMUM_BRANCHING_FACTOR, MAXIMUM_BRANCHING_FACTOR)
	return nil
}

func (r *RTreeStorage) WaypointsLen() int {
	return r.waypointsTree.Size()
}

func (r *RTreeStorage) ConstraintsLen() int {
	return r.constraintsTree.Size()
}

// ---------------------------------------------------------------- WAYPOINTS

func (r *RTreeStorage) AddWaypoint(w *models.Waypoint) error {
	r.waypointsTree.Insert(w)
	return nil
}

func (r *RTreeStorage) AddWaypoints(w_list []*models.Waypoint) error {
	if w_list == nil {
		return nil
	}

	for _, w := range w_list {
		if err := r.AddWaypoint(w); err != nil {
			return err
		}
	}
	return nil
}

// TODO: Not efficient as it have to scan the keys of the map, think about maintaining an array if it's needed
func (r *RTreeStorage) GetWaypoints() ([]*models.Waypoint, error) {
	keys := make([]*models.Waypoint, 0, len(r.waypointsMap))
	for w := range r.waypointsMap {
		keys = append(keys, w)
	}
	return keys, nil
}

// TODO: Not efficient as it have to scan the keys of the map, think about maintaining an array if it's needed
func (r *RTreeStorage) MustGetWaypoints() []*models.Waypoint {
	wps, err := r.GetWaypoints()
	if err != nil {
		panic(err)
	}
	return wps
}

// ---------------------------------------------------------------- CONSTRAINTS

func (r *RTreeStorage) AddConstraint(c *models.Feature3D) error {
	r.constraintsTree.Insert(c)
	return nil
}

func (r *RTreeStorage) AddConstraints(c_list []*models.Feature3D) error {
	if c_list == nil {
		return nil
	}

	for _, c := range c_list {
		if err := r.AddConstraint(c); err != nil {
			return err
		}
	}
	return nil
}

// ================================================================= RRT

func (r *RTreeStorage) AddWaypointWithPrevious(prev *models.Waypoint, w *models.Waypoint) error {
	// Add waypoint to the list of waypoints
	r.AddWaypoint(w)

	// But also add it in the map
	r.ChangePrevious(prev, w)
	return nil
}

func (r *RTreeStorage) ChangePrevious(new_prev *models.Waypoint, w *models.Waypoint) error {
	distance := 0.0
	if new_prev != nil {
		distance = utils.HaversineDistance3D(new_prev, w)
	}
	
	// Just add it to the map
	prev, ok := r.waypointsMap[w]
	if !ok {
		// w never added to the map, add it
		r.waypointsMap[w] = models.NewPointDist(new_prev, distance)
	} else {
		// w already added
		prev.Point = new_prev
		prev.Distance = distance
		r.waypointsMap[w] = prev
	}
	
	return nil
}

func (r *RTreeStorage) GetPrevious(w *models.Waypoint) (*models.Waypoint, error) {
	// Search in the map
	return r.waypointsMap[w].Point, nil
}

// TODO: Same as memoryStorage, maybe put it in the utils
func (r *RTreeStorage) GetPathToRoot(w *models.Waypoint) ([]*models.Waypoint, error) {
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
		prev, ok := r.waypointsMap[current]
		if !ok {
    		return nil, fmt.Errorf("waypoint %+v not found in map", w)
		}

		current = prev.Point
	}
}

// TODO: Same as memoryStorage, maybe put it in the utils
func (r *RTreeStorage) GetCostToRoot(w *models.Waypoint) (float64, error) {
	// Iteratively do GetPrevious and search for the cost
	current := w
	cost := 0.0
	for {
		if current == nil {
			// Reach root
			break
		}

		// Add distance
		prev, ok := r.waypointsMap[current]
		if !ok {
    		return 0.0, fmt.Errorf("waypoint %+v not found in map", w)
		}

		cost += prev.Distance
		current = prev.Point
	}

	return cost, nil
}

// ================================================================= Geometric helpers

// TODO: Make sure it's using the right distance function when doing so
// Use 1-nn with rtree to find nearest point
// O(logN)
func (r *RTreeStorage) NearestPoint(p *models.Waypoint) (*models.Waypoint, float64, error) {
	// TODO: Just for visual debug
	delete(p.Feature.Properties, "parameter")
	delete(p.Feature.Properties, "nearest")
	p.Feature.Properties["parameter"] = true
	
	if r.WaypointsLen() == 0 {
		return nil, 0.0, errors.New("no waypoints")
	}

	// Use rtree functionalities to improve 1-nn search
	nearest := r.waypointsTree.NearestNeighbor(p.RTreePoint()).(*models.Waypoint)
	minDist := utils.HaversineDistance3D(p, nearest)

	// TODO: Just for visual debug
	nearest.Feature.Properties["nearest"] = true
	return nearest, minDist, nil
}

// TODO: Make sure it's using the right distance function when doing so
// Use k-nn with rtree to find nearest points.
// O(logN)
func (r *RTreeStorage) KNearestPoints(p *models.Waypoint, k int) ([]*models.Waypoint, []float64, error) {
	// TODO: Just for visual debug
	delete(p.Feature.Properties, "parameter")
	delete(p.Feature.Properties, "near")
	delete(p.Feature.Properties, "nearest")
	p.Feature.Properties["parameter"] = true
	
	wpsLen := r.WaypointsLen()
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

	// Use r-tree k-nn
	points := r.waypointsTree.NearestNeighbors(k, p.RTreePoint())

	// Create list for points
	result := make([]*models.Waypoint, 0, len(points))
	// And list for distances
	distances := make([]float64, 0, len(points))

	// Convert from rtreego.Spatial to *Waypoint and compute dist
	for _, ps := range points {
		wp := ps.(*models.Waypoint)
		wp.Feature.Properties["near"] = true

		result = append(result, wp)
		distances = append(distances, utils.HaversineDistance3D(p, wp))
	}

	return result, distances, nil
}

// Find points that intersects with circle bbox using rtree. Then keep only the ones that actually intersects with it, not just the bbox.
// O(logN)
func (r *RTreeStorage) NearestPointsInRadius(p *models.Waypoint, radius float64) ([]*models.Waypoint, []float64, error) {
	// TODO: Just for visual debug
	delete(p.Feature.Properties, "parameter")
	delete(p.Feature.Properties, "near")
	delete(p.Feature.Properties, "nearest")
	p.Feature.Properties["parameter"] = true
	
	wpsLen := r.WaypointsLen()
	if wpsLen == 0 {
		return nil, nil, errors.New("no waypoints")
	}
	if radius < 0 {
		return nil, nil, errors.New("radius must be positive")
	}
	if radius == 0 {
		return []*models.Waypoint{}, []float64{0.0}, nil
	}

	// Create circle radius
	circle := p.CircleAroundWaypointGeodesic(radius)
	
	// Use searchIntersect to retrieve point in circle bbox and filter out points that are not exactly inside circle
	points := r.waypointsTree.SearchIntersect(circle.Bounds(), func(results []rtreego.Spatial, object rtreego.Spatial) (refuse bool, abort bool) {
		// Return only objects that are actually inside the circle
		return !utils.PointInPolygon(object.(*models.Waypoint), circle), false
	})

		// Create list for points
	result := make([]*models.Waypoint, 0, len(points))
	// And list for distances
	distances := make([]float64, 0, len(points))

	// Convert from rtreego.Spatial to *Waypoint and compute dist
	for _, ps := range points {
		wp := ps.(*models.Waypoint)
		wp.Feature.Properties["near"] = true

		result = append(result, wp)
		distances = append(distances, utils.HaversineDistance3D(p, wp))
	}

	return result, distances, nil
}

// Use Rtree to efficiently get constraints whose bbox intersects with point. Then check if the point actually intersect the constraint, and not just its bbox. Stop as soon as you find one.
// O(logM)
func (r *RTreeStorage) IsPointInObstacles(p *models.Waypoint) (bool, *models.Feature3D, error) {	
	// Get obstacles that the point intersects with their bbox
	intersectedConstraintsBBox := r.constraintsTree.SearchIntersect(p.Bounds(), func(results []rtreego.Spatial, object rtreego.Spatial) (refuse bool, abort bool) {
		// Check if object actually intersect with point
		if utils.PointInPolygon(p, object.(*models.Feature3D)) {
			// If yes, abort operation
			return false, true
		} else {
			return true, false
		}
	})

	if len(intersectedConstraintsBBox) >= 1 {
		return true, intersectedConstraintsBBox[0].(*models.Feature3D), nil
	}
	
	return false, nil, nil
}

// Use Rtree to efficiently get constraints whose bbox intersects with point. Then check if the point actually intersect the constraint, and not just its bbox.
// O(logM)
func (r *RTreeStorage) GetAllObstaclesContainingPoint(p *models.Waypoint) ([]*models.Feature3D, error) {
	// Get obstacles that the point intersects with their bbox
	intersectedConstraintsBBox := r.constraintsTree.SearchIntersect(p.Bounds())
	
	obstacles := make([]*models.Feature3D, 0, len(intersectedConstraintsBBox))
	
	for _, obs := range intersectedConstraintsBBox {
		obstacle := obs.(*models.Feature3D)
		if utils.PointInPolygon(p, obstacle) {
			obstacles = append(obstacles, obstacle)
		}
	}

	return obstacles, nil
}

// Scan list of obstacle until you find someone that intersect line
// O(logM)
func (r *RTreeStorage) IsLineInObstacles(p1, p2 *models.Waypoint) (bool, []*models.Waypoint, error) {
	// TODO: Two methods: generate line first and then use rtree intersect for every point you get
	// or: use line bounding box to have constraints in advance and then check if the points collide
	// For now let's use the second one that's easier
	
	// Get obstacles that the line intersects with their bbox
	intersectedConstraintsBBox := r.constraintsTree.SearchIntersect(p1.GetLineStringFeature3D(p2).Bounds())
	constraints := make([]*models.Feature3D, 0, len(intersectedConstraintsBBox))
	for _, c := range intersectedConstraintsBBox {
		constraints = append(constraints, c.(*models.Feature3D))
	}

	in, line := utils.LineInPolygon(p1, p2, constraints...)
	return in, line, nil
}

// Get intersection points (useful for AntPath)
func (r *RTreeStorage)	GetIntersectionPoints(p1, p2 *models.Waypoint) ([]*models.LinePolygonIntersection, error) {
	// Divide line into point and then check if any individual point lies in polygon
	quantizedLine := utils.DefaultResampleLineToInterval(p1, p2)
	// For every point, check if it intersect with a polygon and mark it as:
	// Start Point: if previous point was not intersecting but current is
	// End Point: if previous point was intersecting but current is not
	// Return startpoint, endpoint and the polygon that the line is intersecting
	
	// If first is already inside -> error, impossible
	if firstInside, _, _ := r.IsPointInObstacles(quantizedLine[0]); firstInside {
		return nil, errors.New("first point is already in obstacles")
	}
	if lastInside, _, _ := r.IsPointInObstacles(quantizedLine[0]); lastInside {
		return nil, errors.New("last point is already in obstacles")
	}
	
	lpi_list := make([]*models.LinePolygonIntersection, 0)
	obstacles := make(models.PolygonSet)
	var startPoint, endPoint *models.Waypoint
	previousInside := false
	for i := 1; i < len(quantizedLine); i++ {
		// Check if current point is intersecting
		currentObstacles, _ := r.GetAllObstaclesContainingPoint(quantizedLine[i])
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

func (r *RTreeStorage) Sample(sampler utils.Sampler, sampleVolume *models.Feature3D, alt models.Altitude) (*models.Waypoint, error) {	
	// TODO: Decide which to use
	sampled, err := utils.SampleWithAltitude2D(sampler, sampleVolume.Geometry, alt)
	if err != nil {
		return nil, fmt.Errorf("unexpected error during memoryStorage Sample: %w", err)
	}

	// TODO: Check if sampled was already present
	return sampled, nil
}

func (r *RTreeStorage) SampleFree(sampler utils.Sampler, sampleVolume *models.Feature3D, alt models.Altitude) (*models.Waypoint, error) {
	isInObstacle := true
	var sampled *models.Waypoint
	var err error
	// sample until you found a point that is not in an obstacle
	for isInObstacle {
		sampled, err = r.Sample(sampler, sampleVolume, alt)
		if err != nil {
			return nil, fmt.Errorf("unexpected error during memoryStorage SampleFree: %w", err)
		}
		isInObstacle, _, _ = r.IsPointInObstacles(sampled)
	}
	
	return sampled, nil
}