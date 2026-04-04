package gateway

import (
	"net/http"
	"net/http/httputil"
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
	proxy, ok := h.RouteMap[r.URL.Path]
	if !ok {
		http.NotFound(w, r)
		return
	}
	proxy.ServeHTTP(w, r)
}
