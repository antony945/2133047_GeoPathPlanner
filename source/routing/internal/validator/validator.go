package validator

import (
	"fmt"
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/storage"
)

type Validator interface {
	ValidateMessage(data []byte) (*models.RoutingRequest, error)
	ValidateInput(searchVolume *models.Feature3D, waypoints []*models.Waypoint, constraints []*models.Feature3D) ([]*models.Waypoint, []*models.Feature3D, error)
}

type DefaultValidator struct {
}

func NewDefaultValidator() *DefaultValidator {
	return &DefaultValidator{}
}

func (v *DefaultValidator) ValidateMessage(data []byte) (*models.RoutingRequest, error) {
	req, err := models.NewRoutingRequestFromJson(string(data))
	if err != nil {
		return nil, err
	}

	// TODO: Ensure that data is here (at least routingID, at least 2 wps, at least search_volume, ..)
	return req, nil

	// if !req.IsValid() {
	// 	fmt.Printf("⚠️ Ignoring invalid RoutingRequest: %+v", req)
	// }
}

func (v *DefaultValidator) ValidateInput(searchVolume *models.Feature3D, waypoints []*models.Waypoint, constraints []*models.Feature3D) ([]*models.Waypoint, []*models.Feature3D, error) {
	// Create temp RTree storage
	s, err := storage.NewEmptyRTreeStorage()
	if err != nil {
		return nil, nil, fmt.Errorf("error while creating empty rtree storage in validator: %w", err)
	}
	s.AddConstraints(constraints)
	s.AddWaypoints(waypoints)

	// TODO: 1. Check search volume
	
	// 2. Check constraints, discard ones that are not in search volume
	// for _, obs := range constraints {
	// 	fmt.Printf("old_obstacle: %v\n", obs)
	// }
	validatedConstraints, err := s.GetAllObstaclesInSearchVolume(searchVolume)
	if err != nil {
		return nil, nil, err
	}
	fmt.Printf("%d/%d constraints are in search volume\n", len(validatedConstraints), len(constraints))

	// 3. Check waypoints, discard ones that are not in search volume
	validatedWaypoints, err := s.GetAllWaypointsInSearchVolume(searchVolume)
	if err != nil {
		return nil, nil, err
	}
	fmt.Printf("%d/%d waypoints are in search volume\n", len(validatedWaypoints), len(waypoints))
	
	// TODO: Check if waypoints are blocked by constraints
	// for i, wp := range waypoints {
	// 	inside, poly, err := s.IsPointInObstacles(wp)
	// 	if err != nil {
	// 		return nil, nil, err
	// 	}

	// 	// If wp is "good" then add it to validated ones
	// 	if inside {
	// 		fmt.Printf("wp[%d] blocked by poly %v\n", i, poly)
	// 	} else {
	// 		validatedWaypoints = append(validatedWaypoints, wp)
	// 	}
	// }

	// // If wps < 2:
	// if len(validatedWaypoints) < 2 {
	// 	return nil, nil, fmt.Errorf("error while validating waypoints, only %d/%d are valid (not blocked by constraints): %w", len(validatedWaypoints), len(waypoints), err)
	// }

	return validatedWaypoints, validatedConstraints, nil
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