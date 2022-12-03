package model

import (
	"fmt"
	"golang.org/x/exp/slices"
	"sync"
)

var (
	Registry           = make(map[string]Collector)
	CollectorInstances = make(map[string]CollectorBinding)
	CollectorsMutex    = &sync.RWMutex{}
)

func RegisterCollector(key string, collector Collector) {
	Registry[key] = collector
}

type CollectorBinding struct {
	Collector Collector
	Clients   []string
}

func SubscribeClient(client, topic string) {
	fmt.Printf("- client %v subscribes to %q \n", client, topic)
	if instance, ok := CollectorInstances[topic]; ok {
		if !slices.Contains(instance.Clients, client) {
			instance.Clients = append(instance.Clients, client)
			CollectorInstances[topic] = instance
		}
	}
}

func RemoveClient(client string) {
	for topic, instance := range CollectorInstances {
		for i := 0; i < len(instance.Clients); i++ {
			if instance.Clients[i] == client {
				instance.Clients = slices.Delete(instance.Clients, i, i+1)
				CollectorInstances[topic] = instance
			}
		}
	}
}
