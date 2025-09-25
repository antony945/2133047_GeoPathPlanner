package algorithm

import (
	"geopathplanner/routing/internal/models"
	"geopathplanner/routing/internal/storage"
)

type Algorithm interface {
	Compute(*models.RoutingRequest, storage.Storage) (*models.RoutingResponse, error)
}