package gateway

import (
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/golang-jwt/jwt/v5/request"
	"github.com/fahris-n/sentinel-local/internal/auth"
)

type Handler struct {
	RouteMap map[string]*httputil.ReverseProxy
}

func NewHandler(routeMap map[string]*httputil.ReverseProxy) *Handler {
	return &Handler{
		RouteMap: routeMap,
	}
}

func (h *Handler) HandleRequest(w http.ResponseWriter, r *http.Request) {
	// check jwt
	extractor := &request.BearerExtractor{}
	tokenStr, err := extractor.ExtractToken(r)
	if err != nil {
		log.Println("TokenExtractor error:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	claims, err := auth.ValidateJWT(tokenStr)
	if err != nil {
		log.Println("JWT error:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// ^^later on we can use something like claims.Role to give user specific tooling/resources
	log.Println("UserID:", claims.UserID)

	proxy, ok := h.RouteMap[r.URL.Path]
	if !ok {
		http.NotFound(w, r)
		return
	}
	proxy.ServeHTTP(w, r)
}
