package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jsawo/colin-server/data"
	"github.com/jsawo/colin-server/ws"

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
	for _, collector := range AppConfig.Collectors {
		if collector.Enabled {
			go MonitorCollector(collector)
		}
	}
}

func MonitorCollector(collector Collector) {
	sleepDuration := collector.GetFrequency()
	for {
		result := data.Registry[collector.Key]()
		ws.WriteMessage(collector.Channel, result)
		time.Sleep(sleepDuration)
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
