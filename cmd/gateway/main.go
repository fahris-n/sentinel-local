package main

import (
	"log"
	"net/http"

	"github.com/fahris-n/sentinel-local/internal/config"
	"github.com/fahris-n/sentinel-local/internal/gateway"
	"github.com/fahris-n/sentinel-local/internal/middleware"
	"github.com/fahris-n/sentinel-local/internal/proxy"
	"github.com/fahris-n/sentinel-local/internal/routing"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	routeMap := map[string]*routing.RouteEntry{}
	for _, route := range cfg.Routes {
		proxy, err := proxy.NewReverseProxy(route.Backend, route.BackendPath)
		if err != nil {
			log.Fatal(err)
		}
		routeMap[route.Path] = &routing.RouteEntry{
			Proxy:        proxy,
			RequiresAuth: route.RequiresAuth,
			AllowedRoles: route.AllowedRoles,
		}
	}

	handler := gateway.NewHandler(routeMap)
	wrappedHandler := middleware.Chain(
		handler,
		middleware.Logging,
		middleware.AuthMiddleware(routeMap),
	)

	http.Handle("/api/", wrappedHandler)

	log.Println("gateway listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
