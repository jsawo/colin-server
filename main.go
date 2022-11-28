package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jsawo/colin-server/internal/config"
	"github.com/jsawo/colin-server/internal/model"
	"github.com/jsawo/colin-server/internal/server"
	"github.com/jsawo/colin-server/internal/ws"

	_ "github.com/jsawo/colin-server/internal/collector/cmd"
	_ "github.com/jsawo/colin-server/internal/collector/cpu"
	_ "github.com/jsawo/colin-server/internal/collector/mem"
)

var (
	done = make(chan interface{})
)

func main() {
	config.ReadInConfig()

	go server.StartServer()
	// go generateNoise()

	RunCollectors()

	<-done
}

func RunCollectors() {
	for _, col := range config.CollectorConfigs {
		if col.Enabled {
			go MonitorCollector(col)
		}
	}
}

func MonitorCollector(col model.CollectorConfig) {
	model.Registry[col.Key].Setup(col.Params)
	for {
		result := model.Registry[col.Key].Collect()
		ws.WriteMessage(col.Topic, result)
		time.Sleep(col.Frequency)
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
