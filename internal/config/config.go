package config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/jsawo/colin-server/internal/model"

	"gopkg.in/yaml.v3"
)

const (
	configPath = "./config.yaml"
)

var (
	collectorConfig  *AppConfig
	CollectorConfigs = map[string]model.CollectorConfig{}
)

type AppConfig struct {
	Entries []Entry `yaml:"collectors"`
}

type Entry map[string]string

func ReadInConfig() {
	cfg, err := parseConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	collectorConfig = cfg

	for _, entry := range collectorConfig.Entries {
		parsedCollector := parseCollector(entry)

		if _, ok := CollectorConfigs[parsedCollector.Key]; ok {
			log.Fatalf("multiple collectors with the same key are not allowed: %q", parsedCollector.Key)
		}

		CollectorConfigs[parsedCollector.Key] = parsedCollector
	}
}

func parseCollector(entry Entry) model.CollectorConfig {
	col := model.CollectorConfig{
		Params: map[string]any{},
	}

	validateEntry(entry)

	collectorType, ok := model.ParseCollectorType(entry["type"])
	if !ok {
		log.Fatalf("unrecognized collector type: %s", entry["type"])
	}
	freq, err := time.ParseDuration(entry["frequency"])
	if err != nil {
		log.Fatalf("error when parsing collector frequency %q - %s", entry["frequency"], err.Error())
	}

	col.Enabled = entry["enabled"] == "true"
	col.Type = collectorType
	col.Frequency = freq

	for key, value := range entry {
		switch key {
		case "key":
			col.Key = entry[key]
		case "topic":
			col.Topic = entry[key]
		case "title":
			col.Title = entry[key]
		case "description":
			col.Description = entry[key]
		default:
			col.Params[key] = value
		}
	}

	return col
}

func parseConfig(configPath string) (*AppConfig, error) {
	config := &AppConfig{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func validateEntry(entry Entry) {
	key, ok := entry["key"]
	if !ok {
		log.Fatal("'key' value is missing in a collector configuration")
	}

	requireValue(entry, key, "topic")
	requireValue(entry, key, "title")
	requireValue(entry, key, "description")
}

func requireValue(entry Entry, collectorKey, entryKey string) {
	if _, ok := entry[entryKey]; !ok {
		log.Fatalf("%q is missing in a %q collector configuration", entryKey, collectorKey)
	}

	if strings.Trim(entry[entryKey], " ") == "" {
		log.Fatalf("%q cannot be empty in a %q collector configuration", entryKey, collectorKey)
	}
}
