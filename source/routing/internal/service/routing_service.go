package service

import "geopathplanner/routing/internal/models"

type RoutingService struct {
}

func NewRoutingService() *RoutingService {
	return &RoutingService{}
}

func (rs *RoutingService) HandleRoutingRequest(input *models.RoutingRequest) (*models.RoutingResponse, error) {
	// 1. Pick algorithm (from input)
	// 2. Pick storage (from input)
	// 3. Validate waypoints and constraint
	// 4. Compute route
	// 5. Clear temporary storage if used
	return nil, nil
}