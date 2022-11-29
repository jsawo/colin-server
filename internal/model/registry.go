package model

import "fmt"

var (
	Registry           = make(map[string]Collector)
	CollectorInstances = make(map[string]CollectorBinding)
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
		instance.Clients = append(instance.Clients, client)
		CollectorInstances[topic] = instance
	}
}
