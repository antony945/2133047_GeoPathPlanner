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
	// RequestID   string     `json:"request_id"`   // must match request
	// ReceivedAt  time.Time  `json:"received_at"`  // when request arrived
	*RoutingRequest
	RouteFound  bool       `json:"route_found"`  // true if route was computed
	Route       []*Waypoint `json:"route"`        // final route if found
	CostKm      float64    `json:"cost_km"`      // optional, distance
	Message     string     `json:"message"`      // error or informational message
	CompletedAt time.Time  `json:"completed_at"` // when response generated
}

// Success response
func NewRoutingResponseSuccess(routingRequest *RoutingRequest, route []*Waypoint, costKm float64) *RoutingResponse {
	now := time.Now()
	return &RoutingResponse{
		RoutingRequest: routingRequest,
		RouteFound:  true,
		Route:       route,
		CostKm:      costKm,
		Message:     "Route computed successfully",
		CompletedAt: now,
	}
}

// Error response
func NewRoutingResponseError(routingRequest *RoutingRequest, message string) *RoutingResponse {
	now := time.Now()
	return &RoutingResponse{
		RoutingRequest: routingRequest,
		RouteFound:  false,
		Route:       nil,
		CostKm:      0,
		Message:     message,
		CompletedAt: now,
	}
}