package main

import (
	"encoding/json"
	"net/http"
)

func helloHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := map[string]string {
		"service": "backend-a",
		"message": "hello from backend-a",
	}
	// Should change this to handle errors
	_ = json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.ListenAndServe(":8081", nil)
}
