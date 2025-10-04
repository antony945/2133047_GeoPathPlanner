package validator

import (
	"fmt"
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/storage"
)

type Validator interface {
	ValidateInput(searchVolume *models.Feature3D, waypoints *models.Waypoint, constraints []*models.Feature3D) error
}

type DefaultValidator struct {
	s storage.Storage
}

func NewDefaultValidator() (*DefaultValidator, error) {
	// TODO: For now use MemoryStorage that's easy
	s, err := storage.NewEmptyMemoryStorage()
	if err != nil {
		return nil, fmt.Errorf("error while creating empty memory storage: %w", err)
	}
	return &DefaultValidator{
		s: s,
	}, nil
}

func (v *DefaultValidator) ValidateInput(searchVolume *models.Feature3D, waypoints []*models.Waypoint, constraints []*models.Feature3D) error {
	// TODO: Check search volume
	// TODO: Check constraints
	// TODO: Check if waypoints are blocked by constraints
	for i, wp := range waypoints {
		inside, poly, err := v.s.IsPointInObstacles(wp)
		if inside {
			return fmt.Errorf("wp[%d] blocked by poly %v", i, poly)
		}
		if err != nil {
			return err
		}
	}
	return nil
} 