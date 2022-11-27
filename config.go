package main

import (
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	configPath = "./config.yaml"
)

var (
	collectorConfig  *AppConfig
	CollectorConfigs = map[string]CollectorConfig{}
)

type AppConfig struct {
	Entries []Entry `yaml:"collectors"`
}

type Entry map[string]string

func readInConfig() {
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

func parseCollector(entry Entry) CollectorConfig {
	collector := CollectorConfig{
		Params: map[string]any{},
	}

	validateEntry(entry)

	collectorType, ok := parseCollectorType(entry["type"])
	if !ok {
		log.Fatalf("unrecognized collector type: %s", entry["type"])
	}
	freq, err := time.ParseDuration(entry["frequency"])
	if err != nil {
		log.Fatalf("error when parsing collector frequency %q - %s", entry["frequency"], err.Error())
	}

	collector.Enabled = entry["enabled"] == "true"
	collector.Type = collectorType
	collector.Frequency = freq

	for key, value := range entry {
		switch key {
		case "key":
			collector.Key = entry[key]
		case "channel":
			collector.Channel = entry[key]
		case "title":
			collector.Title = entry[key]
		case "description":
			collector.Description = entry[key]
		default:
			collector.Params[key] = value
		}
	}

	return collector
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

	requireValue(entry, key, "channel")
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
