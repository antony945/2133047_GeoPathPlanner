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

func (rs *RoutingService) HandleRoutingRequest(input *models.RoutingRequest, val validator.Validator) (*models.RoutingResponse, bool) {
	// TODO: Think about this
	// 1. Validate waypoints and constraint
	wps, constraints, err := val.ValidateInput(input.SearchVolume, input.Waypoints, input.Constraints)
	if err != nil {
		return models.NewRoutingResponseError(input, err.Error()), false
	}

	// 2. Pick and create algorithm (from input)
	algo, err := algorithm.NewAlgorithm(input.Algorithm())
	if err != nil {
		return models.NewRoutingResponseError(input, err.Error()), false
	}

	// 3. Compute route
	// TODO: Test with both compute and computeConcurrently
	route, cost, err := algo.ComputeConcurrently(input.SearchVolume, wps, constraints, input.Parameters, input.Storage(), 0)
	if err != nil {
		return models.NewRoutingResponseError(input, err.Error()), false
	}

	// 4. Return route
	return models.NewRoutingResponseSuccess(input, route, cost), true
}