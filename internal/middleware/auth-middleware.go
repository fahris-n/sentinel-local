package middleware

import (
	"context"
	"net/http"

	"github.com/fahris-n/sentinel-local/internal/auth"
	"github.com/fahris-n/sentinel-local/internal/routing"
	"github.com/golang-jwt/jwt/v5/request"
)

type contextKey string

const ClaimsKey contextKey = "claims"

func AuthMiddleware(routeMap map[string]*routing.RouteEntry) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			route, ok := routeMap[r.URL.Path]
			if !ok || !route.RequiresAuth {
				next.ServeHTTP(w, r)
				return
			}

			// extract Bearer token
			extractor := &request.BearerExtractor{}
			tokenStr, err := extractor.ExtractToken(r)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// validate JWT
			claims, err := auth.ValidateJWT(tokenStr)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// check role against roles
			if len(route.AllowedRoles) > 0 {
				if !contains(route.AllowedRoles, claims.Role) {
					http.Error(w, "Forbidden", http.StatusForbidden)
					return
				}
			}

			ctx := context.WithValue(r.Context(), ClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func contains(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}
