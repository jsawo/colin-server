package data

import (
	"fmt"
	"os"
	"time"

	"github.com/mackerelio/go-osstat/cpu"
)

const cpuKey = "cpu"

func init() {
	register(cpuKey, getCpu)
}

func getCpu() any {
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
