package models

type LinePolygonIntersection struct {
	EnteringPoint *Waypoint
	ExitingPoint  *Waypoint
	Polygon       *Feature3D
}

func NewLinePolygonIntersection(entering, exiting *Waypoint, poly *Feature3D) *LinePolygonIntersection {
	return &LinePolygonIntersection{
		EnteringPoint: entering,
		ExitingPoint:  exiting,
		Polygon:       poly,
	}
}