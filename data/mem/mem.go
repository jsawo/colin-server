package data

import (
	"fmt"
	"github.com/jsawo/colin-server/data"
	"os"

	"github.com/mackerelio/go-osstat/memory"
)

const key = "mem"

type memCollector struct{}

func init() {
	data.Register(key, &memCollector{})
}

func (c *memCollector) Setup(params map[string]any) {}

func (*memCollector) Collect() any {
	mem, err := memory.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return 0
	}

	return float64(mem.Used) / float64(mem.Total) * 100
}
