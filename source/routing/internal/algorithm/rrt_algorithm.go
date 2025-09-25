package algorithm

import (
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/storage"
)

type RRTAlgorithm struct {
}

func NewRRTAlgorithm() (*RRTAlgorithm, error) {
	// TODO: To implement
	return &RRTAlgorithm{}, nil
}

// TODO: Implement RRT Algorithm
func (a *RRTAlgorithm) Compute(routingRequest *models.RoutingRequest, storage storage.Storage) (*models.RoutingResponse, error) {
	return &models.RoutingResponse{}, nil
}
