package main

import (
	"time"

	"github.com/jsawo/colin-server/internal/cache"

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

func MonitorCollector(cfg model.CollectorConfig) {
	model.CollectorsMutex.Lock()
	model.CollectorInstances[cfg.Topic] = model.CollectorBinding{
		Collector: model.Registry[cfg.Collector].NewCollector(cfg),
	}
	model.CollectorsMutex.Unlock()

	for {
		result := model.CollectorInstances[cfg.Topic].Collector.Collect()
		cache.AddValue(cfg.Topic, result)
		ws.SendMessageToSubscribers(ws.Message{
			Topic:     cfg.Topic,
			Payload:   result,
			Timestamp: time.Now(),
		})

		time.Sleep(cfg.Frequency)
	}
}
