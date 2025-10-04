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
func (a *RRTAlgorithm) Compute(searchVolume *models.Feature3D, waypoints []*models.Waypoint, constraints []*models.Feature3D, parameters map[string]any, storage storage.Storage) ([]*models.Waypoint, float64, error) {
	return []*models.Waypoint{}, 0.0, nil
}
