package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/fahris-n/sentinel-local/internal/gateway"
	"github.com/fahris-n/sentinel-local/internal/middleware"
	"github.com/fahris-n/sentinel-local/internal/proxy"
	"github.com/joho/godotenv"
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

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	routeMap := map[string]*httputil.ReverseProxy{}
	cfg, err := LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	for _, route := range cfg.Routes {
		proxy, err := proxy.NewReverseProxy(route.Backend, route.BackendPath)
		if err != nil {
			log.Fatal(err)
		}
		routeMap[route.Path] = proxy
	}

	handler := gateway.NewHandler(routeMap)
	wrappedHandler := middleware.Chain(
		handler,
		middleware.Logging,
		middleware.AuthMiddleware,
	)

	http.Handle("/api/", wrappedHandler)

	log.Println("gateway listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
