package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jsawo/colin-server/data"

	"github.com/jsawo/colin-server/ws"

	_ "github.com/jsawo/colin-server/data/cmd"
	_ "github.com/jsawo/colin-server/data/cpu"
	_ "github.com/jsawo/colin-server/data/mem"
)

func main() {
	readInConfig()

	go StartServer()
	// go generateNoise()

	RunCollectors()

	<-done
}

func RunCollectors() {
	for _, collector := range CollectorConfigs {
		if collector.Enabled {
			go MonitorCollector(collector)
		}
	}
}

func MonitorCollector(collector CollectorConfig) {
	data.Registry[collector.Key].Setup(collector.Params)
	for {
		result := data.Registry[collector.Key].Collect()
		ws.WriteMessage(collector.Channel, result)
		time.Sleep(collector.Frequency)
	}
}

func generateNoise() {
	rand.Seed(time.Now().UnixNano())
	for {
		n := rand.Intn(10)
		time.Sleep(time.Duration(n) * time.Second)
		ws.WriteMessage("dummy", fmt.Sprintf("Message delayed %v seconds\n", n))
	}
}
