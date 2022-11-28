package model

import "strings"

const (
	TypeGauge = iota
	TypeText
	TypeJson
	TypeCounter
	TypeHistogram
)

type CollectorType int

var collectorTypeMap = map[string]CollectorType{
	"gauge":     TypeGauge,
	"text":      TypeText,
	"json":      TypeJson,
	"counter":   TypeCounter,
	"histogram": TypeHistogram,
}

func ParseCollectorType(str string) (CollectorType, bool) {
	c, ok := collectorTypeMap[strings.ToLower(str)]
	return c, ok
}

func (c CollectorType) ToString() string {
	for key, value := range collectorTypeMap {
		if value == c {
			return key
		}
	}

	return "unknown"
}
