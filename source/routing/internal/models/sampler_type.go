package models

import (
	"encoding/json"
	"fmt"
)

type SamplerType string

const (
	Uniform SamplerType = "uniform"
	Halton  SamplerType = "halton"
)

// Validate sampler type (enforce enum)
func (s SamplerType) Validate() error {
	switch s {
	case Uniform, Halton:
		return nil
	default:
		return fmt.Errorf("invalid sampler type: %s, available options are %s, %s", s, Uniform, Halton)
	}
}

func (s *SamplerType) UnmarshalJSON(data []byte) error {
    var value string
    if err := json.Unmarshal(data, &value); err != nil {
        return err
    }

	*s = SamplerType(value)
	if err := s.Validate(); err != nil {
		return err
	}

	return nil
}