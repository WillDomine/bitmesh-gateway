package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// bitmesh-gateway port
type ServerConfig struct {
	Port string `yaml:"port"`
}

// config.yaml structure
type Config struct {
	Server   ServerConfig      `yaml:"server"`
	Services map[string]string `yaml:"services"`
}

func LoadConfig(path string) (*Config, error) {
	//Open to read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to read file: %w", err)
	}
	//Preallocate to parse YAML
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("Failed to parse yaml: %w", err)
	}

	return &cfg, nil
}
