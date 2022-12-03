package model

import (
	"time"
)

type Collector interface {
	NewCollector(collectorConfig CollectorConfig) Collector
	Collect() any
}

type CollectorConfig struct {
	Collector   string
	Topic       string
	Title       string
	Description string
	Enabled     bool
	Type        CollectorType
	Frequency   time.Duration
	Params      map[string]any
}
