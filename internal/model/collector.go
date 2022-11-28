package model

import (
	"time"
)

type Collector interface {
	Setup(collectorConfig map[string]any)
	Collect() any
}

type CollectorConfig struct {
	Key         string
	Channel     string
	Title       string
	Description string
	Enabled     bool
	Type        CollectorType
	Frequency   time.Duration
	Params      map[string]any
}
