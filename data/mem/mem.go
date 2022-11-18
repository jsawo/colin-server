package data

import (
	"fmt"
	"os"

	"github.com/jsawo/colin-server/data"
	"github.com/mackerelio/go-osstat/memory"
)

const key = "mem"

func init() {
	data.Register(key, collect)
}

func collect() any {
	memory, err := memory.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return 0
	}

	return float64(memory.Used) / float64(memory.Total) * 100
}
