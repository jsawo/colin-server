package data

type Collect func() any

type Collector interface {
	Setup(collectorConfig map[string]any)
	Collect() any
}

var Registry = make(map[string]Collector)

func Register(key string, collector Collector) {
	Registry[key] = collector
}
