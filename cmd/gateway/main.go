package main

import (
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/fahris-n/sentinel-local/internal/gateway"
	"github.com/fahris-n/sentinel-local/internal/proxy"
)

func main() {
	// Create proxies for every path we want to support
	// on each backend we want to support
	helloProxy, err := proxy.NewReverseProxy("http://localhost:8081", "/hello")
	if err != nil {
		log.Fatal(err)
	}
	registerProxy, err := proxy.NewReverseProxy("http://localhost:8082", "/register")
	if err != nil {
		log.Fatal(err)
	}

	// generate proxy map
	routeMap := map[string]*httputil.ReverseProxy{
		"/api/hello":    helloProxy,
		"/api/register": registerProxy,
	}

	// create handler and pass proxy map
	handler := gateway.NewHandler(routeMap)
	http.HandleFunc("/api/", handler.HandleRequest)

	log.Println("gateway listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
