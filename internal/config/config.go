package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type RouteConfig struct {
	Path         string   `yaml:"path"`
	Backend      string   `yaml:"backend"`
	BackendPath  string   `yaml:"backendPath"`
	RequiresAuth bool     `yaml:"requiresAuth"`
	AllowedRoles []string `yaml:"allowedRoles"`
}

type Config struct {
	Routes []RouteConfig `yaml:"routes"`
}

func LoadConfig(path string) (Config, error) {
	var cfg Config

	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("failed to read config file: %w", err)
	}

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("failed to parse config file: %w", err)
	}

	return cfg, nil
}
