package models

type PointDist struct {
	Point    *Waypoint
	Distance float64
}

func NewPointDist(point *Waypoint, distance float64) *PointDist {
	return &PointDist{
		Point:    point,
		Distance: distance,
	}
}
