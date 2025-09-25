package service

import (
	"geopathplanner/routing/internal/algorithm"
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/storage"
)

type RoutingService struct {
}

func NewRoutingService() *RoutingService {
	return &RoutingService{}
}

func (rs *RoutingService) HandleRoutingRequest(input *models.RoutingRequest) (*models.RoutingResponse, error) {
	// 1. Pick and create algorithm (from input)
	algo, err := algorithm.NewAlgorithm(input.Algorithm)
	if err != nil {
		return nil, err
	}

	// 1b. If necessary, validate waypoints and constraint

	// 2. Pick and create storage (from input)
	stor, err := storage.NewStorage(input.Waypoints, input.Constraints, input.Storage)
	if err != nil {
		return nil, err
	}
	
	// 3. Compute route
	output, err := algo.Compute(input, stor)

	// 4. If necessary, clear temporary storage if used
	// TODO: Think where to put the clearance, in algorithm maybe?
	return output, err
}