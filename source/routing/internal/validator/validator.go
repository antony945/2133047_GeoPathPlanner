package validator

import (
	"fmt"
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/storage"
)

type Validator interface {
	ValidateInput(searchVolume *models.Feature3D, waypoints *models.Waypoint, constraints []*models.Feature3D) error
}

type DefaultValidator struct {
	s storage.Storage
}

func NewDefaultValidator() (*DefaultValidator, error) {
	// TODO: For now use ListStorage that's easy
	s, err := storage.NewEmptyListStorage()
	if err != nil {
		return nil, fmt.Errorf("error while creating empty list storage: %w", err)
	}
	return &DefaultValidator{
		s: s,
	}, nil
}

func (v *DefaultValidator) ValidateInput(searchVolume *models.Feature3D, waypoints []*models.Waypoint, constraints []*models.Feature3D) error {
	// TODO: Check search volume
	
	// TODO: Check constraints
	// newConstraints, err := v.ConstraintUnion(constraints)
	// if err != nil {
	// 	return err
	// }

	// TODO: Check if waypoints are blocked by constraints
	for i, wp := range waypoints {
		inside, poly, err := v.s.IsPointInObstacles(wp)
		if inside {
			return fmt.Errorf("wp[%d] blocked by poly %v", i, poly)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// func (v *DefaultValidator) ConstraintUnion(constraints []*models.Feature3D) ([]*models.Feature3D, error) {
// 	// // Find min and max altitude first
// 	// var minAlt, maxAlt models.Altitude
// 	// minAlt, _ = models.NewAltitude(models.DEFAULT_MAX_ALT, models.MT)
// 	// maxAlt, _ = models.NewAltitude(models.DEFAULT_MIN_ALT, models.MT)

// 	// for _, f := range constraints {
// 	// 	minAltCurrent := f.MinAltitude.Normalize()
// 	// 	maxAltCurrent := f.MaxAltitude.Normalize()

// 	// 	if minAltCurrent.Compare(minAlt) < 0 {
// 	// 		// new min alt
// 	// 		minAlt = minAltCurrent
// 	// 	}
// 	// 	if maxAltCurrent.Compare(maxAlt) > 0 {
// 	// 		// new max alt
// 	// 		maxAlt = maxAltCurrent
// 	// 	}
// 	// }

// 	// Convert every constraint to polygol to perform union 
// 	unioned, err := polygol.Union()
// 	if err != nil {
// 		return err
// 	}

// }

// func p2g(p [][][][]float64) orb.Geometry {

// 	g := make(orb.MultiPolygon, len(p))

// 	for i := range p {
// 		g[i] = make([]orb.Ring, len(p[i]))
// 		for j := range p[i] {
// 			g[i][j] = make([]orb.Point, len(p[i][j]))
// 			for k := range p[i][j] {
// 				pt := p[i][j][k]
// 				point := orb.Point{pt[0], pt[1]}
// 				g[i][j][k] = point
// 			}
// 		}
// 	}
// 	return g
// }