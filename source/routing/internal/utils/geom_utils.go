package utils

import (
	"fmt"
	"geopathplanner/routing/internal/models"
	"math"

	"github.com/engelsjk/polygol"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/geojson"
	"github.com/paulmach/orb/planar"
	"github.com/paulmach/orb/resample"
)

const (
	DEFAULT_LINE_DIVISION_MAX_STEP_SIZE_MT = 50
)

// Implement POINT-POLYGON intersection
func PointInPolygon(p *models.Waypoint, poly *models.Feature3D) bool {	
	// 1. First check intersection with BBox of poly -> if not, then immediately return false
	if !poly.Geometry.Bound().Contains(p.Point2D()) {
		p.Feature.Properties["inside"] = false
		return false 
	}

	// 2. If yes, check the altitude -> if point altitude is not within min and max poly altitude return false
	if !p.Alt.IsWithin(poly.MinAltitude, poly.MaxAltitude) {
		p.Feature.Properties["inside"] = false
		return false
	}

	// 3. Here you have to check exactly if it insersects: run PiP algorithm (PnPoly, uses RayTracing) algorithm to do that
	// Add "inside" property if it's inside the polygon
	isInside := PointInPolygon2D(p.Point2D(), poly.ToPolygon())
	// TODO: For now just for testing
	p.Feature.Properties["inside"] = isInside
	return isInside
}

func PointInPolygon2D(p orb.Point, poly orb.Polygon) bool {
	return planar.PolygonContains(poly, p);
}

// Implement LINE-POLYGON intersection
func LineInPolygon(p1, p2 *models.Waypoint, polygons ...*models.Feature3D) (bool, []*models.Waypoint) {
	// Use linebound to rapidly check if it's inside polygons or not
	bound_intersects := false
	for _, poly := range polygons {
		if bound_intersects = poly.Geometry.Bound().Intersects(p1.GetLineStringBound(p2)); bound_intersects {
			break
		}	
	}
	if !bound_intersects {
		p1.Feature.Properties["inside"] = false
		p2.Feature.Properties["inside"] = false
		return false, []*models.Waypoint{p1, p2}
	}
	
	// Divide line into point and then check if any individual point lies in polygon
	quantizedLine := DefaultResampleLineToInterval(p1, p2)
	// Do it in a smart way, by checking intermediate point first and then recursively the left and the right parts of the line
	// TODO: Think about keeping the recursive version or changing to iterative
	var inside bool
	for _, poly := range polygons {
		inside = _sublineInPolygonLinear(quantizedLine, poly)
		// inside = _sublineInPolygonRecursive(quantizedLine, poly, 0, len(quantizedLine)-1)
		if inside {
			break
		}
	}
	return inside, quantizedLine
}

func LineInPolygonRemoveLast(p1, p2 *models.Waypoint, polygons ...*models.Feature3D) (bool, []*models.Waypoint) {
	// Divide line into point and then check if any individual point lies in polygon
	quantizedLine := DefaultResampleLineToInterval(p1, p2)
	// Do it in a smart way, by checking intermediate point first and then recursively the left and the right parts of the line
	// TODO: Think about keeping the recursive version or changing to iterative
	var inside bool
	for _, poly := range polygons {
		inside = _sublineInPolygonLinear(quantizedLine[:len(quantizedLine)-1], poly)
		// inside = _sublineInPolygonRecursive(quantizedLine, poly, 0, len(quantizedLine)-1)
		if inside {
			break
		}
	}
	return inside, quantizedLine
}

func _sublineInPolygonLinear(quantizedLine []*models.Waypoint, poly *models.Feature3D) bool {
	for _, p := range quantizedLine {
		if PointInPolygon(p, poly) {
			return true
		}
	}

	return false
}

func _sublineInPolygonRecursive(quantizedLine []*models.Waypoint, poly *models.Feature3D, start, end int) bool {
	// Base case
	if start > end {
		return false
	}
	
	middle := (start + end) / 2
	if PointInPolygon(quantizedLine[middle], poly) {
		return true
	}

	// Here we need to left and right search
	if _sublineInPolygonRecursive(quantizedLine, poly, start, middle-1) {
		return true
	}

	if _sublineInPolygonRecursive(quantizedLine, poly, middle+1, end) {
		return true
	}

	return false
}

// Implement line division into evenly spaced points
func ResampleLineToInterval(p1, p2 *models.Waypoint, distMt float64) []*models.Waypoint {
	if (p1 == p2) {
		return []*models.Waypoint{p1, p2}
	}
	
	// Check if line's ends distance between each other is already less than step size (base case)
	dist := HaversineDistance3D(*p1, *p2)
	if (dist <= distMt) {
		return []*models.Waypoint{p1, p2}
	}

	// Define numStep, altStepSize and stepSize
	numStep := int(dist / distMt)
	altStepSizeMt := p1.Alt.Distance(p2.Alt).ConvertTo(models.MT).Value / float64(numStep)
	
	// Pythagorean theorem to find x stepSize
	stepSizeMt := math.Sqrt(distMt*distMt - altStepSizeMt*altStepSizeMt)

	quantizedPoints := make([]*models.Waypoint, 0, numStep)

	// 2D resample line	
	resampledLine := resample.ToInterval(p1.GetLineString(p2), geo.DistanceHaversine, stepSizeMt)
	
	// For each point in resampleLine, add the altitude interpolated
	startingAltVal := p1.Alt.ConvertTo(models.MT).Value
	for i, p := range resampledLine {
		// For altitude linearly interpolate
		altVal := startingAltVal + float64(i)*altStepSizeMt
        alt, _ := models.NewAltitude(altVal, models.MT)
		// TODO: Debug error in case
		wp, _ := models.NewWaypoint(p.Lat(), p.Lon(), alt)
		
		quantizedPoints = append(quantizedPoints, wp)
	}

	return quantizedPoints
}

func GetPointInDirectionAtDistance(p1, p2 *models.Waypoint, distMt float64) *models.Waypoint {
	return ResampleLineToInterval(p1, p2, distMt)[1]
}

func DefaultResampleLineToInterval(p1, p2 *models.Waypoint) []*models.Waypoint {
	return ResampleLineToInterval(p1, p2, DEFAULT_LINE_DIVISION_MAX_STEP_SIZE_MT)
}

func GetNearestFreeVertexIndex(c *models.Feature3D, p *models.Waypoint, reversed bool) int {
	vertices := c.GetVertices(p.Alt, reversed)

	minIndex := -1
	minDist := 1e9
	for i := range vertices {
		// make sure that a line between p and vertices[i] can be draw
		if blocked, _ := LineInPolygonRemoveLast(p, vertices[i], c); blocked {
			continue
		}

		dist := HaversineDistance3D(*p, *vertices[i])
		// Check distance between p and v
		if dist < minDist {
			// min = v
			minDist = dist
			minIndex = i
		}
	}
	return minIndex
}

func GetBestWayToGoAroundPolygon(c *models.Feature3D, enteringPoint, exitingPoint *models.Waypoint) []*models.Waypoint {
	// Consider vertices in both ways
	normalDirection := getWayToGoAroundPolygon(c, enteringPoint, exitingPoint, false)
	oppositeDirection := getWayToGoAroundPolygon(c, enteringPoint, exitingPoint, true) 

	if TotalHaversineDistance(normalDirection) < TotalHaversineDistance(oppositeDirection) {
		return normalDirection
	} else {
		return oppositeDirection
	}
}

func getWayToGoAroundPolygon(c *models.Feature3D, enteringPoint, exitingPoint *models.Waypoint, reversed bool) []*models.Waypoint {
	bestWay := []*models.Waypoint{enteringPoint}

	// Get start and end vertex index
	vertices := c.GetVertices(enteringPoint.Alt, reversed)
	startVertexIndex := GetNearestFreeVertexIndex(c, enteringPoint, reversed)
	endVertexIndex := GetNearestFreeVertexIndex(c, exitingPoint, reversed)
	
	// 3 cases
	// startIndex == endIndex
	if (startVertexIndex == endVertexIndex) {
		// 1 single vertex
		bestWay = append(bestWay, vertices[startVertexIndex])
	} else if (startVertexIndex < endVertexIndex) {
		// take the vertices in between them
		bestWay = append(bestWay, vertices[startVertexIndex:endVertexIndex+1]...)
	} else {
		// startVertexIndex > endVertexIndex
		// First append from start index to end of array
		bestWay = append(bestWay, vertices[startVertexIndex:]...)
		// Then append from beginning to end index
		bestWay = append(bestWay, vertices[:endVertexIndex+1]...)
	}

	bestWay = append(bestWay, exitingPoint)
	return bestWay
}

func FindMinMaxAltitude(features []*models.Feature3D) (models.Altitude, models.Altitude) {
	// Find min and max altitude first
	var minAlt, maxAlt models.Altitude
	minAlt, _ = models.NewAltitude(models.DEFAULT_MAX_ALT, models.MT)
	maxAlt, _ = models.NewAltitude(models.DEFAULT_MIN_ALT, models.MT)

	for _, f := range features {
		minAltCurrent := f.MinAltitude.Normalize()
		maxAltCurrent := f.MaxAltitude.Normalize()

		if minAltCurrent.Compare(minAlt) < 0 {
			// new min alt
			minAlt = minAltCurrent
		}
		if maxAltCurrent.Compare(maxAlt) > 0 {
			// new max alt
			maxAlt = maxAltCurrent
		}
	}

	return minAlt, maxAlt
}

func UnionFeatures(features []*models.Feature3D) ([]*models.Feature3D, error) {
	// Get min max altitude that will be set for all the features
	minAlt, maxAlt := FindMinMaxAltitude(features)

	// Convert features to polygol geom
	polygons := make([]polygol.Geom, 0, len(features))
	for _, f := range features {
		polygons = append(polygons, f.ToPolygol())
	}
	
	// Perform union
	result, err := polygol.Union(polygons[0], polygons[1:]...)
	if err != nil {
		return nil, fmt.Errorf("error while performing features union: %w", err)
	}

	// Convert back to list of features
	unionedFeatures, err := PolygolToListOfFeature(result, minAlt, maxAlt)
	if err != nil {
		return nil, fmt.Errorf("error while converting from polygol.Geom to Feature3D: %w", err)
	}

	return unionedFeatures, nil
}

func PolygolToListOfFeature(p [][][][]float64, minAltitude, maxAltitude models.Altitude) ([]*models.Feature3D, error) {
	feature_list := make([]*models.Feature3D, 0, len(p))

	// FOR EVERY POLYGON IN THE MULTIPOLYGON
	for _, polyData := range p {
		polygon := make(orb.Polygon, len(polyData))
		// FOR EVERY RING IN THE POLYGON
		for j, ringData := range polyData {
			polygon[j] = make(orb.Ring, len(ringData))
			// FOR EVERY POINT IN THE RING
			for k, pointData := range ringData {
				point := orb.Point{pointData[0], pointData[1]}
				polygon[j][k] = point
			}
		}

		// Here we have our polygon
		f, err := models.NewFeatureFromGeojsonFeature(geojson.NewFeature(polygon))
		if err != nil {
			return nil, err
		}
		f.SetAltitude(minAltitude, maxAltitude)
		
		// Set min, max altitude
		feature_list = append(feature_list, f)
	}

	return feature_list, nil
}