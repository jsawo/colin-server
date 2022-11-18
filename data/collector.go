package data

type Collect func() any

var Registry = make(map[string]func() any)

func Register(topic string, callable Collect) {
	Registry[topic] = callable
}
