package middleware

import (
	"log"
	"net/http"
	"time"
)

type StatusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *StatusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		log.Printf("started %s %s", r.Method, r.URL.Path)
		rec := StatusRecorder{w, 200}

		next.ServeHTTP(&rec, r)
		log.Printf("completed %s %s: status: %d, time: %v", r.Method, r.URL.Path, rec.status, time.Since(startTime))
	})
}
