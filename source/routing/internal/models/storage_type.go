package models

import (
	"encoding/json"
	"fmt"
)

type StorageType string

const (
	Memory StorageType = "memory"
	Redis  StorageType = "redis"
	RTree StorageType = "rtree"
	DEFAULT_STORAGE StorageType = RTree
)

// Validate algorithm type (enforce enum)
func (s StorageType) Validate() error {
	switch s {
	case Memory, Redis, RTree:
		return nil
	default:
		return fmt.Errorf("invalid storage type: %s, available options are %s, %s, %s", s, Memory, Redis, RTree)
	}
}

func (s *StorageType) UnmarshalJSON(data []byte) error {
    var value string
    if err := json.Unmarshal(data, &value); err != nil {
        return err
    }

	*s = StorageType(value)
	if err := s.Validate(); err != nil {
		return err
	}

	return nil
}