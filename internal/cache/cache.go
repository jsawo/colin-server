package cache

var storage = make(map[string]any)

func AddValue(topic string, value any) {
	storage[topic] = value
}

func GetLatestValue(topic string) any {
	if value, ok := storage[topic]; ok {
		return value
	}

	return ""
}
