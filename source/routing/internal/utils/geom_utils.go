package utils

import (
	"geopathplanner/routing/internal/models"
	"math"

	"github.com/tidwall/geodesic"
)

const (
	DEFAULT_LINE_DIVISION_MAX_STEP_SIZE_MT = 50
)

// HaversineDistance3D calculates the 3D distance between two waypoints using the haversine formula
// and incorporates elevation differences. Returns result in mt.
func HaversineDistance3D(p1, p2 *models.Waypoint) float64 {
	const earthRadius = 6371000 // Earth's radius in meters

	lat1 := p1.Lat * math.Pi / 180.0
	lon1 := p1.Lon * math.Pi / 180.0
	lat2 := p2.Lat * math.Pi / 180.0
	lon2 := p2.Lon * math.Pi / 180.0

	// Haversine formula
	dlat := lat2 - lat1
	dlon := lon2 - lon1
	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	
	// 2D distance
	distance2D := earthRadius * c

	// Include elevation difference
	elevDiff := p2.Alt.Subtract(p1.Alt).Value
	
	// Calculate 3D distance using Pythagorean theorem
	distance3D_mt := math.Sqrt(distance2D*distance2D + elevDiff*elevDiff)
	return distance3D_mt
}

// TODO: Implement POINT-POLYGON intersection
func PointInPolygon(p *models.Waypoint, poly *models.Constraint) bool {
	return false
}

// TODO: Implement LINE-POLYGON intersection
func LineInPolygon(p1, p2 *models.Waypoint, poly *models.Constraint, stepSizeMt float64) bool {
	// Divide line into point and then check if any individual point lies in polygon
	quantizedLine := DivideLineInPoints(p1, p2, stepSizeMt)
	
	// Do it in a smart way, by checking intermediate point first and then recursively the left and the right parts of the line
	// TODO: Think about keeping the recursive version or changing to iterative
	return _sublineInPolygon(quantizedLine, poly, 0, len(quantizedLine)-1)
}

func DefaultLineInPolygon(p1, p2 *models.Waypoint, poly *models.Constraint) bool {
	return LineInPolygon(p1, p2, poly, DEFAULT_LINE_DIVISION_MAX_STEP_SIZE_MT);
}

// Implement line division into evenly spaced points
// TODO: To test
func DivideLineInPoints(p1, p2 *models.Waypoint, stepSizeMt float64) []*models.Waypoint {
	// Check if line's ends distance between each other is already less than step size (base case)
	dist := HaversineDistance3D(p1, p2)
	if (dist <= stepSizeMt) {
		return []*models.Waypoint{p1, p2}
	}

	numStep := int(dist / stepSizeMt)
	stepSizeMt = dist / float64(numStep)
	altStepSizeMt := p1.Alt.Distance(p2.Alt).ConvertTo(models.MT).Value / float64(numStep)
	
	quantizedPoints := make([]*models.Waypoint, 0, numStep)
	quantizedPoints = append(quantizedPoints, p1)
	
	// Find solutions to direct geodesic problem to get initial bearing
	var azi1, azi2 float64
	geodesic.WGS84.Inverse(p1.Lat, p1.Lon, p2.Lat, p2.Lon, nil, &azi1, &azi2)

	// Idea is to iteratively find point X moving ourselves from p1 in p2 direction by just the max_step_size_mt amount.
	currentPoint := p1
	for i := 0; i < numStep; i++ {
		var lat, lon float64
		geodesic.WGS84.Direct(currentPoint.Lat, currentPoint.Lon, azi1, stepSizeMt, &lat, &lon, nil)

		// For altitude linearly interpolate
		altVal := p1.Alt.ConvertTo(models.MT).Value + altStepSizeMt
        alt, _ := models.NewAltitude(altVal, models.MT)

		// Create new waypoint
		currentPoint, _ = models.NewWaypoint(lat, lon, alt)
		// TODO: Debug error in case

		quantizedPoints = append(quantizedPoints, currentPoint)
	}

	quantizedPoints = append(quantizedPoints, p2)
	return quantizedPoints
}

func DefaultDivideLineInPoints(p1, p2 *models.Waypoint) []*models.Waypoint {
	return DivideLineInPoints(p1, p2, DEFAULT_LINE_DIVISION_MAX_STEP_SIZE_MT)
}

func _sublineInPolygon(quantizedLine []*models.Waypoint, poly *models.Constraint, start, end int) bool {
	// Base case
	if start > end {
		return false
	}
	
	middle := (start + end) / 2
	if PointInPolygon(quantizedLine[middle], poly) {
		return true
	}

	// Here we need to left and right search
	if _sublineInPolygon(quantizedLine, poly, start, middle-1) {
		return true
	}

	if _sublineInPolygon(quantizedLine, poly, middle+1, end) {
		return true
	}

	return false
}