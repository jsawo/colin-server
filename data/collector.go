package data

type Collector func() any

var Registry = make(map[string]func() any)

func register(topic string, callable Collector) {
	Registry[topic] = callable
}
