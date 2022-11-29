package cmd

import (
	"log"
	"os/exec"

	"github.com/jsawo/colin-server/internal/model"
)

const key = "cmd"

type cmdCollector struct {
	config    model.CollectorConfig
	command   string
	directory string
}

func init() {
	model.RegisterCollector(key, &cmdCollector{})
}

func (c *cmdCollector) Setup(config model.CollectorConfig) model.Collector {
	if _, ok := config.Params["command"]; !ok {
		log.Fatal("'command' key is missing in 'cmd' collector configuration")
	}

	if dir, ok := config.Params["directory"]; ok {
		c.directory = dir.(string)
	}

	c.command = config.Params["command"].(string)
	return c
}

func (c *cmdCollector) Collect() any {
	cmd := exec.Command("sh", "-c", c.command)
	if c.directory != "" {
		cmd.Dir = c.directory
	}
	out, _ := cmd.CombinedOutput()

	return string(out)
}
