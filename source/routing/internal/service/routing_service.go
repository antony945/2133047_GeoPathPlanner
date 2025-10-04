package service

import (
	"geopathplanner/routing/internal/algorithm"
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/storage"
	"geopathplanner/routing/internal/validator"
)

type RoutingService struct {
}

func NewRoutingService() *RoutingService {
	return &RoutingService{}
}

func (rs *RoutingService) HandleRoutingRequest(input *models.RoutingRequest) (*models.RoutingResponse) {
	// 1. Validate waypoints and constraint
	val, err := validator.NewDefaultValidator()
	if err != nil {
		return models.NewRoutingResponseError(input.RequestID, input.ReceivedAt, err.Error())
	}
	err = val.ValidateInput(input.SearchVolume, input.Waypoints, input.Constraints)
	if err != nil {
		return models.NewRoutingResponseError(input.RequestID, input.ReceivedAt, err.Error())
	}
	
	// 2. Pick and create storage (from input)	
	stor, err := storage.NewEmptyStorage(input.Storage)
	if err != nil {
		return models.NewRoutingResponseError(input.RequestID, input.ReceivedAt, err.Error())
	}

	// 3. Pick and create algorithm (from input)
	algo, err := algorithm.NewAlgorithm(input.Algorithm)
	if err != nil {
		return models.NewRoutingResponseError(input.RequestID, input.ReceivedAt, err.Error())
	}

	// 4. Compute route
	route, cost, err := algo.Compute(input.SearchVolume, input.Waypoints, input.Constraints, input.Parameters, stor)
	if err != nil {
		return models.NewRoutingResponseError(input.RequestID, input.ReceivedAt, err.Error())
	}


	// 5. If necessary, clear temporary storage if used
	// TODO: Think where to put the clearance, in algorithm maybe?

	return models.NewRoutingResponseSuccess(input.RequestID, input.ReceivedAt, route, cost)
}