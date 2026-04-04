package middleware

import (
	"log"
	"net/http"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("started %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("completed %s %s", r.Method, r.URL.Path)
	})
}
