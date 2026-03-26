package main

import (
	"log"
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
	log.Println("backend path:", r.URL.Path)
}

func main() {
	http.HandleFunc("/register", registerUserHandler)
	log.Println("backend_b listening on :8082")
	http.ListenAndServe(":8082", nil)
}
