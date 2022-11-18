package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	configPath = "./config.yaml"
)

var (
	AppConfig *Config
)

type Config struct {
	Collectors []Collector `yaml:"collectors"`
}

type Collector struct {
	Key         string `yaml:"key"`
	Channel     string `yaml:"channel"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Enabled     bool   `yaml:"enabled"`
	Type        string `yaml:"type"`
	Frequency   string `yaml:"frequency"`
}

func (c *Collector) GetFrequency() time.Duration {
	freq, err := time.ParseDuration(c.Frequency)
	if err != nil {
		log.Fatalf("Error when parsing collector frequency %q - %s", c.Frequency, err.Error())
	}

	return freq
}

func readInConfig() {
	cfg, err := NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	AppConfig = cfg

	fmt.Printf("Parsed config: %+v \n", AppConfig)
}

func NewConfig(configPath string) (*Config, error) {
	config := &Config{}

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

func ValidateConfigPath(path string) error {
	fileStat, err := os.Stat(path)
	if err != nil {
		return err
	}
	if fileStat.IsDir() {
		return fmt.Errorf("%q is a directory, file expected", path)
	}
	return nil
}
