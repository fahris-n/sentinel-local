package middleware

import (
	"context"
	"net/http"

	"github.com/fahris-n/sentinel-local/internal/auth"
	"github.com/golang-jwt/jwt/v5/request"
)

type contextKey string

const ClaimsKey contextKey = "claims"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		// need to make a new context to pass our claims
		ctx := context.WithValue(r.Context(), ClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
