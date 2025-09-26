package models

import (
	"time"
)

type RoutingRequest struct {
	RequestID   string         	`json:"request_id"`  	// unique ID for this request
	Waypoints   []*Waypoint     `json:"waypoints"`   	// at least 2 waypoints
	Constraints []*Feature3D  	`json:"constraints"` 	// constraints
	SearchVolume *Feature3D 	`json:"search_volume"` 	// search area
	Parameters  map[string]any 	`json:"parameters"`  	// optional additional params (may be related to algorithm, may not)
	Algorithm   AlgorithmType  	`json:"algorithm"`   	// optional, dev/testing only
	Storage     StorageType    	`json:"storage"`     	// optional, dev/testing only
	ReceivedAt  time.Time      	`json:"received_at"` 	// when request arrived (unix timestamp)
}
