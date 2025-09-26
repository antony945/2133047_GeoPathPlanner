package utils

import (
	"geopathplanner/routing/internal/models"
	"math"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
	"github.com/paulmach/orb/planar"
	"github.com/paulmach/orb/resample"
)

const (
	DEFAULT_LINE_DIVISION_MAX_STEP_SIZE_MT = 50
)

// Implement POINT-POLYGON intersection
func PointInPolygon(p models.Waypoint, poly *models.Feature3D) bool {
	// 1. First check intersection with BBox of poly -> if not, then immediately return false
	if !poly.Geometry.Bound().Contains(p.Point2D()) {
		return false 
	}

	// 2. If yes, check the altitude -> if point altitude is not within min and max poly altitude return false
	if !p.Alt.IsWithin(poly.MinAltitude, poly.MaxAltitude) {
		return false
	}

	// 3. Here you have to check exactly if it insersects: run PiP algorithm (PnPoly, uses RayTracing) algorithm to do that
	return PointInPolygon2D(p.Point2D(), poly.ToPolygon())
}

func PointInPolygon2D(p orb.Point, poly orb.Polygon) bool {
	return planar.PolygonContains(poly, p);
}

// Implement LINE-POLYGON intersection
func LineInPolygon(p1, p2 models.Waypoint, poly *models.Feature3D) bool {
	// Divide line into point and then check if any individual point lies in polygon
	quantizedLine := DefaultResampleLineToInterval(p1, p2)
	
	// Do it in a smart way, by checking intermediate point first and then recursively the left and the right parts of the line
	// TODO: Think about keeping the recursive version or changing to iterative
	return _sublineInPolygonRecursive(quantizedLine, poly, 0, len(quantizedLine)-1)
}

func _sublineInPolygonRecursive(quantizedLine []models.Waypoint, poly *models.Feature3D, start, end int) bool {
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
func ResampleLineToInterval(p1, p2 models.Waypoint, distMt float64) []models.Waypoint {
	// Check if line's ends distance between each other is already less than step size (base case)
	dist := HaversineDistance3D(p1, p2)
	if (dist <= distMt) {
		return []models.Waypoint{p1, p2}
	}

	// Define numStep, altStepSize and stepSize
	numStep := int(dist / distMt)
	altStepSizeMt := p1.Alt.Distance(p2.Alt).ConvertTo(models.MT).Value / float64(numStep)
	
	// Pythagorean theorem to find x stepSize
	stepSizeMt := math.Sqrt(distMt*distMt - altStepSizeMt*altStepSizeMt)

	quantizedPoints := make([]models.Waypoint, 0, numStep)

	// 2D resample line	
	resampledLine := resample.ToInterval(orb.LineString{p1.Point2D(), p2.Point2D()}, geo.DistanceHaversine, stepSizeMt)
	
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

func DefaultResampleLineToInterval(p1, p2 models.Waypoint) []models.Waypoint {
	return ResampleLineToInterval(p1, p2, DEFAULT_LINE_DIVISION_MAX_STEP_SIZE_MT)
}