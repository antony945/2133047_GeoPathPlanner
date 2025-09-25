package models

import "time"

type RoutingResponse struct {
	RequestID   string     `json:"request_id"`   // must match request
	RouteFound  bool       `json:"route_found"`  // true if route was computed
	Route       []Waypoint `json:"route"`        // final route if found
	CostKm      float64    `json:"cost_km"`      // optional, distance
	Message     string     `json:"message"`      // error or informational message
	ReceivedAt  time.Time  `json:"received_at"`  // when request arrived
	CompletedAt time.Time  `json:"completed_at"` // when response generated
}