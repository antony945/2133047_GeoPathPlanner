package service

import (
	"geopathplanner/routing/internal/algorithm"
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/validator"
)

type RoutingService struct {
}

func NewRoutingService() (*RoutingService, error) {
	return &RoutingService{}, nil
}

func (rs *RoutingService) HandleRoutingRequest(input *models.RoutingRequest) (*models.RoutingResponse) {
	// TODO: Think about this
	// 1. Validate waypoints and constraint
	val, err := validator.NewDefaultValidator()
	if err != nil {
		return models.NewRoutingResponseError(input.RequestID, input.ReceivedAt, err.Error())
	}
	err = val.ValidateInput(input.SearchVolume, input.Waypoints, input.Constraints)
	if err != nil {
		return models.NewRoutingResponseError(input.RequestID, input.ReceivedAt, err.Error())
	}

	// 2. Pick and create algorithm (from input)
	algo, err := algorithm.NewAlgorithm(input.Algorithm())
	if err != nil {
		return models.NewRoutingResponseError(input.RequestID, input.ReceivedAt, err.Error())
	}

	// 3. Compute route
	// TODO: Test with both compute and computeConcurrently
	route, cost, err := algo.ComputeConcurrently(input.SearchVolume, input.Waypoints, input.Constraints, input.Parameters, input.Storage(), 0)
	if err != nil {
		return models.NewRoutingResponseError(input.RequestID, input.ReceivedAt, err.Error())
	}

	// 4. Return route
	return models.NewRoutingResponseSuccess(input.RequestID, input.ReceivedAt, route, cost)
}