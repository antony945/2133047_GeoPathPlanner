package algorithm

import (
	"fmt"
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/storage"
)

type Algorithm interface {
	Compute(*models.RoutingRequest, storage.Storage) (*models.RoutingResponse, error)
}

func NewAlgorithm(algorithmType models.AlgorithmType) (Algorithm, error) {
	switch algorithmType {
	case models.RRT:
		return NewRRTAlgorithm()
	case models.RRTStar, models.AntPath:
		return nil, fmt.Errorf("algorithm currently not implemented: %s", algorithmType)
	default:
		return nil, fmt.Errorf("algorithm not recognized: %s", algorithmType)

	}
}