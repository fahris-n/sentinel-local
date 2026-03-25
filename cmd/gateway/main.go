package main

import (
	"log"
	"net/http"

	"github.com/fahris-n/sentinel-local/internal/gateway"
	"github.com/fahris-n/sentinel-local/internal/proxy"
)

func main() {
	// create reverse proxy obj
	helloProxy, err := proxy.NewReverseProxy()
	if err != nil {
		log.Fatal(err)
	}

	// create handler. proxy obj is param to handler
	handler := gateway.NewHandler(helloProxy)
	http.HandleFunc("/api/hello", handler.HelloHandler)

	log.Println("gateway listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
