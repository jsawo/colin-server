package cpu

import (
	"fmt"
	"github.com/jsawo/colin-server/data"
	"os"
	"time"

	"github.com/mackerelio/go-osstat/cpu"
)

const key = "cpu"

type cpuCollector struct{}

func init() {
	data.Register(key, &cpuCollector{})
}

func (c *cpuCollector) Setup(params map[string]any) {}

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
