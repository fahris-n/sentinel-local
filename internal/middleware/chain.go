package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, m...Middleware) http.Handler {
	// wrap in reverse order so first middleware in list runs first
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}
	return h
}
