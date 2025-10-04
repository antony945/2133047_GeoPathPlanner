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

func (rs *RoutingService) HandleRoutingRequest(input *models.RoutingRequest) (*models.RoutingResponse) {
	// 1. Pick and create algorithm (from input)
	algo, err := algorithm.NewAlgorithm(input.Algorithm)
	if err != nil {
		return models.NewRoutingResponseError(input.RequestID, input.ReceivedAt, err.Error())
	}

	// 1b. If necessary, validate waypoints and constraint

	// 2. Pick and create storage (from input)
	stor, err := storage.NewStorage(input.Waypoints, input.Constraints, input.Storage)
	if err != nil {
		return models.NewRoutingResponseError(input.RequestID, input.ReceivedAt, err.Error())
	}
	
	// 3. Compute route
	route, cost, err := algo.Compute(input.SearchVolume, input.Waypoints, input.Constraints, input.Parameters, stor)
	if err != nil {
		return models.NewRoutingResponseError(input.RequestID, input.ReceivedAt, err.Error())
	}


	// 4. If necessary, clear temporary storage if used
	// TODO: Think where to put the clearance, in algorithm maybe?



	return models.NewRoutingResponseSuccess(input.RequestID, input.ReceivedAt, route, cost)
}