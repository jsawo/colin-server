package data

import (
	"fmt"
	"os"

	"github.com/mackerelio/go-osstat/memory"
)

const memKey = "mem"

func init() {
	register(memKey, getMem)
}

func getMem() any {
	memory, err := memory.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return 0
	}

	return float64(memory.Used) / float64(memory.Total) * 100
}
