package models

type AlgorithmType string
type StorageType string

const (
	RRT     AlgorithmType = "RRT"
	RRTStar AlgorithmType = "RRTStar"
	AntPath AlgorithmType = "AntPath"
	Memory  StorageType   = "Memory"
	Redis   StorageType   = "Redis"
)

type RoutingRequest struct {
	RoutingID   string         `json:"request_id"`  // unique ID for this request
	Waypoints   []*Waypoint    `json:"waypoints"`   // at least 2 waypoints
	Constraints []*Constraint  `json:"constraints"` // optional
	Algorithm   AlgorithmType  `json:"algorithm"`   // optional, dev/testing only
	Parameters  map[string]any `json:"parameters"`  // optional additional params related to algorithm
	Storage     StorageType    `json:"storage"`     // optional, dev/testing only
	ReceivedAt  int64          `json:"received_at"` // when request arrived (unix timestamp)
}
