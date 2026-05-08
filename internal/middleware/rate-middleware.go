package middleware

import (
	"net"
	"net/http"

	"github.com/fahris-n/sentinel-local/internal/ratelimit"
	"github.com/fahris-n/sentinel-local/internal/routing"
)

func RateLimitMiddleware(limiter *ratelimit.Limiter, routeMap map[string]*routing.RouteEntry) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			route := routeMap[r.URL.Path]
			ip, _, _ := net.SplitHostPort(r.RemoteAddr)
			allowed, err := limiter.Allow(ip, route.MaxTokens, route.RefillRate)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if !allowed {
				http.Error(w, "Rate Limit Excedded", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
