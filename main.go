package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jsawo/colin-server/data"
	"github.com/jsawo/colin-server/ws"
)

func main() {
	readInConfig()

	go StartServer()
	go generateNoise()

	RunCollectors()

	<-done
}

func RunCollectors() {
	for _, collector := range AppConfig.Collectors {
		fmt.Printf("starting monitoring for a collector: %v \n", collector.Key)
		go MonitorCollector(collector)
	}
}

func MonitorCollector(collector Collector) {
	for {
		result := data.Registry[collector.Key]()
		ws.WriteMessage(collector.Channel, result)
		time.Sleep(time.Duration(5) * time.Second)
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
