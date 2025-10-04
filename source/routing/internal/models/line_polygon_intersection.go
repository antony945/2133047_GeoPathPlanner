package models

type PolygonSet map[*Feature3D]struct{}

func (s PolygonSet) Add(o *Feature3D) {
	s[o] = struct{}{}
}

func (s PolygonSet) AddAll(features ...*Feature3D) {
	for _, f := range features {
		s.Add(f)
	}
}

func (s PolygonSet) Has(o *Feature3D) bool {
	_, ok := s[o]
	return ok
}

func (s PolygonSet) Values() []*Feature3D {
	values := make([]*Feature3D, 0, len(s))
	for o := range s {
		values = append(values, o)
	}
	return values
}

func (s PolygonSet) Clear() {
	for k := range s {
		delete(s, k)
	}
}

type LinePolygonIntersection struct {
	EnteringPoint *Waypoint
	ExitingPoint  *Waypoint
	Polygons      []*Feature3D
}

func NewLinePolygonIntersection(entering, exiting *Waypoint, polygons []*Feature3D) *LinePolygonIntersection {
	return &LinePolygonIntersection{
		EnteringPoint: entering,
		ExitingPoint:  exiting,
		Polygons:      polygons,
	}
}