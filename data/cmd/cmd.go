package cmd

import (
	"fmt"
	"github.com/jsawo/colin-server/data"
	"log"
	"os/exec"
)

const key = "cmd"

type cmdCollector struct {
	command   string
	directory string
}

func init() {
	data.Register(key, &cmdCollector{})
}

func (c *cmdCollector) Setup(params map[string]any) {
	if _, ok := params["command"]; !ok {
		log.Fatal("'command' key is missing in 'cmd' collector configuration")
	}

	if dir, ok := params["directory"]; ok {
		c.directory = dir.(string)
	}

	c.command = params["command"].(string)
}

func (c *cmdCollector) Collect() any {
	fmt.Println("running setup for CMD: " + c.command)
	cmd := exec.Command("sh", "-c", c.command)
	if c.directory != "" {
		cmd.Dir = c.directory
	}
	out, _ := cmd.CombinedOutput()

	return string(out)
}
