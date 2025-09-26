package utils

import (
	"geopathplanner/routing/internal/models"
	"math"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
)

type DistanceFunc3D func(models.Waypoint, models.Waypoint) float64
type DistanceFunc2D orb.DistanceFunc

// Distance3D calculates the between two waypoints using any 2d distance function.
// Then add altitude difference to find the 3D distance. Returns result in mt.
// TODO: Check if it's handled the case when p1=p2

func Distance3D(p1, p2 models.Waypoint, distanceFunc2D DistanceFunc2D) float64 {
	// Use a 2D distance
	distance2D := distanceFunc2D(p1.Point2D(), p2.Point2D())

	// Include elevation difference
	elevDiff := p2.Alt.Subtract(p1.Alt).Value
	
	// Calculate 3D distance using Pythagorean theorem
	distance3D_mt := math.Sqrt(distance2D*distance2D + elevDiff*elevDiff)
	return distance3D_mt
}

// HaversineDistance3D calculates the 3D distance between two waypoints using the haversine formula. Returns result in mt.
func HaversineDistance3D(p1, p2 models.Waypoint) float64 {
	return Distance3D(p1, p2, geo.DistanceHaversine)
}

// FastDistance3D calculates the 3D distance between two waypoints using a fast distance formula. Returns result in mt.
func FastDistance3D(p1, p2 models.Waypoint) float64 {
	return Distance3D(p1, p2, geo.Distance)
}