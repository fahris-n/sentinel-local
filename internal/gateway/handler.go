package gateway

import (
	"net/http"

	"github.com/fahris-n/sentinel-local/internal/routing"
)

type Handler struct {
	RouteMap map[string]*routing.RouteEntry
}

func NewHandler(routeMap map[string]*routing.RouteEntry) *Handler {
	return &Handler{
		RouteMap: routeMap,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route, ok := h.RouteMap[r.URL.Path]
	if !ok {
		http.NotFound(w, r)
		return
	}
	route.Proxy.ServeHTTP(w, r)
}
