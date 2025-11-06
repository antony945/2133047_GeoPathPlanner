package models

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type RoutingRequest struct {
	RequestID   string         	`json:"request_id"`  	// unique ID for this request
	Waypoints   []*Waypoint     `json:"waypoints"`   	// at least 2 waypoints
	Constraints []*Feature3D  	`json:"constraints"` 	// constraints
	SearchVolume *Feature3D 	`json:"search_volume"` 	// search area
	Parameters  map[string]any 	`json:"parameters"`  	// optional additional params (may be related to algorithm, may not)
	ReceivedAt  time.Time      	`json:"received_at"` 	// when request arrived (unix timestamp)
}

func NewRoutingRequestFromJsonFile(filename string) (*RoutingRequest, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	return NewRoutingRequestFromJson(string(data))
}

func NewRoutingRequestFromJson(jsonString string) (*RoutingRequest, error) {
	var r *RoutingRequest
	if err := json.Unmarshal([]byte(jsonString), &r); err != nil {
		return nil, fmt.Errorf("unmarshaling json: %w", err)
	}

	// fmt.Printf("routingRequest: %+v\n", r)
	return r, nil
}

func MustNewRoutingRequestFromJson(jsonString string) (*RoutingRequest) {
	r, err := NewRoutingRequestFromJson(jsonString)
	if err != nil {
		panic(err)
	}

	return r
}

func (r *RoutingRequest) Algorithm() AlgorithmType {
	if r.Parameters == nil {
		// no parameters -> default algorithm
		return DEFAULT_ALGORITHM
	}

	// Try to get the "algorithm" key
	if val, ok := r.Parameters["algorithm"]; ok {
		// Convert to string if possible
		if s, ok := val.(string); ok {
			a := AlgorithmType(s)
			if err := a.Validate(); err == nil {
				return a
			}
			// invalid value -> fall back to default
		}
		// if it's not a string -> ignore
	}

	// default algorithm
	return DEFAULT_ALGORITHM
}

func (r *RoutingRequest) Storage() StorageType {
	if r.Parameters == nil {
		// no parameters -> default storage
		return DEFAULT_STORAGE
	}

	// Try to get the "storage" key
	if val, ok := r.Parameters["storage"]; ok {
		// Convert to string if possible
		if s, ok := val.(string); ok {
			a := StorageType(s)
			if err := a.Validate(); err == nil {
				return a
			}
			// invalid value -> fall back to default
		}
		// if it's not a string -> ignore
	}

	// default storage
	return DEFAULT_STORAGE
}