package models

import (
	"encoding/json"
	"fmt"
)

type AlgorithmType string

const (
	RRT     AlgorithmType = "rrt"
	RRTStar AlgorithmType = "rrtstar"
	AntPath AlgorithmType = "antpath"
	// TODO: Decide which one
	DEFAULT_ALGORITHM AlgorithmType = RRTStar
)

// Validate algorithm type (enforce enum)
func (a AlgorithmType) Validate() error {
	switch a {
	case RRT, RRTStar, AntPath:
		return nil
	default:
		return fmt.Errorf("invalid algorithm type: %s, available options are %s, %s, %s", a, RRT, RRTStar, AntPath)
	}
}

func (a *AlgorithmType) UnmarshalJSON(data []byte) error {
    var value string
    if err := json.Unmarshal(data, &value); err != nil {
        return err
    }

	*a = AlgorithmType(value)
	if err := a.Validate(); err != nil {
		return err
	}

	return nil
}