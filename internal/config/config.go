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
	MaxTokens    int16    `yaml:"maxTokens"`
	RefillRate   int16    `yaml:"refillRate"`
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

func ValidateConfig(cfg Config) error {
	if len(cfg.Routes) == 0 {
		return fmt.Errorf("no routes defined")
	}

	for i, route := range cfg.Routes {
		if route.Path == "" {
			return fmt.Errorf("route at index %d is missing a path", i)
		}
		if route.Backend == "" {
			return fmt.Errorf("route at index %d is missing a backend", i)
		}
		if route.BackendPath == "" {
			return fmt.Errorf("route at index %d is missing a backend path", i)
		}
		if route.RequiresAuth && len(route.AllowedRoles) == 0 {
			return fmt.Errorf("route at index %d requires auth but does not have any defined allowed roles", i)
		}
		if route.MaxTokens <= 0 {
			return fmt.Errorf("route at index %d has no tokens", i)
		}
		if route.RefillRate <= 0 {
			return fmt.Errorf("route at index %d has no refill rate", i)
		}
	}

	return nil
}
