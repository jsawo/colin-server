package cpu

import (
	"fmt"
	"os"
	"time"

	"github.com/jsawo/colin-server/internal/model"

	"github.com/mackerelio/go-osstat/cpu"
)

const key = "cpu"

type cpuCollector struct {
	config model.CollectorConfig
}

func init() {
	model.RegisterCollector(key, &cpuCollector{})
}

func (c *cpuCollector) NewCollector(config model.CollectorConfig) model.Collector {
	return c
}

func (c *cpuCollector) Collect() any {
	before, err := cpu.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return nil
	}
	time.Sleep(time.Duration(1) * time.Second)
	after, err := cpu.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return nil
	}
	total := float64(after.Total - before.Total)
	idle := float64(after.Idle-before.Idle) / total * 100

	return 100 - idle
}
