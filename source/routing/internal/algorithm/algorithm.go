package algorithm

import (
	"fmt"
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/storage"
)

type Algorithm interface {
	Compute(searchVolume *models.Feature3D, waypoints []*models.Waypoint, constraints []*models.Feature3D, parameters map[string]any, storage storage.Storage) ([]*models.Waypoint, float64, error)
}

func NewAlgorithm(algorithmType models.AlgorithmType) (Algorithm, error) {
	switch algorithmType {
	case models.RRT:
		return NewRRTAlgorithm()
	case models.AntPath:
		return NewAntPathAlgorithm()
	case models.RRTStar:
		return nil, fmt.Errorf("algorithm currently not implemented: %s", algorithmType)
	default:
		return nil, fmt.Errorf("algorithm not recognized: %s", algorithmType)

	}
}