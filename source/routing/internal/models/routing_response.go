package models

import "time"

// TODO: To delete, this is for Umberto's db
// type DbModel struct {
// 	// ------------------- RELATED TO REQUEST
// 	UserID string
// 	RequestID   string     `json:"request_id"`   // must match request
// 	Waypoints   []*Waypoint     `json:"waypoints"`   	// at least 2 waypoints
// 	Constraints []*Feature3D  	`json:"constraints"` 	// constraints
// 	SearchVolume *Feature3D 	`json:"search_volume"` 	// search area
// 	Parameters  map[string]any 	`json:"parameters"`  	// optional additional params (may be related to algorithm, may not)
// 	ReceivedAt  time.Time  `json:"received_at"`  // when request arrived
// 	// ------------------- RELATED TO RESPONSE
// 	RouteFound  bool       `json:"route_found"`  // true if route was computed
// 	Route       []*Waypoint `json:"route"`        // final route if found
// 	CostKm      float64    `json:"cost_km"`      // optional, distance
// 	Message     string     `json:"message"`      // error or informational message
// 	CompletedAt time.Time  `json:"completed_at"` // when response generated
// }

type RoutingResponse struct {
	RequestID   string     `json:"request_id"`   // must match request
	RouteFound  bool       `json:"route_found"`  // true if route was computed
	Route       []*Waypoint `json:"route"`        // final route if found
	CostKm      float64    `json:"cost_km"`      // optional, distance
	Message     string     `json:"message"`      // error or informational message
	ReceivedAt  time.Time  `json:"received_at"`  // when request arrived
	CompletedAt time.Time  `json:"completed_at"` // when response generated
}

// Success response
func NewRoutingResponseSuccess(requestID string, receivedAt time.Time, route []*Waypoint, costKm float64) *RoutingResponse {
	now := time.Now()
	return &RoutingResponse{
		RequestID:   requestID,
		RouteFound:  true,
		Route:       route,
		CostKm:      costKm,
		Message:     "Route computed successfully",
		ReceivedAt:  now,
		CompletedAt: now,
	}
}

// Error response
func NewRoutingResponseError(requestID string, receivedAt time.Time, message string) *RoutingResponse {
	now := time.Now()
	return &RoutingResponse{
		RequestID:   requestID,
		RouteFound:  false,
		Route:       nil,
		CostKm:      0,
		Message:     message,
		ReceivedAt:  now,
		CompletedAt: now,
	}
}