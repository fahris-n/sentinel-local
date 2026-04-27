package main

import (
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
		log.Printf("yamlFile.Get err	#%v", err)
	}
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return cfg, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	routeMap := map[string]*httputil.ReverseProxy{}
	cfg, _ := LoadConfig("configs/config.yaml")

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
