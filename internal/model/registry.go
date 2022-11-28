package model

type Collect func() any

var Registry = make(map[string]Collector)

func Register(key string, collector Collector) {
	Registry[key] = collector
}
