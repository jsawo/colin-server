package data

import (
	"fmt"
	"os"

	"github.com/jsawo/colin-server/internal/model"

	"github.com/mackerelio/go-osstat/memory"
)

const key = "mem"

type memCollector struct {
	config model.CollectorConfig
}

func init() {
	model.RegisterCollector(key, &memCollector{})
}

func (c *memCollector) NewCollector(config model.CollectorConfig) model.Collector {
	return c
}

func (c *memCollector) Collect() any {
	mem, err := memory.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return 0
	}

	return float64(mem.Used) / float64(mem.Total) * 100
}
