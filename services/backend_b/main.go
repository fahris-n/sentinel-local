package main

import (
	"net/http"
	"encoding/json"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := map[string]string {
		"service": "backend_b",
		"message": "hello from backend_b",
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.ListenAndServe(":8082", nil)
}
