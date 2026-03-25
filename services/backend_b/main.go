package main

import (
	"net/http"
	"encoding/json"
)

func registerUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := map[string]string {
		"service": "backend_b",
		"message": "register user",
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/register", registerUserHandler)
	http.ListenAndServe(":8082", nil)
}
