package gateway

import (
	"log"
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
	log.Println("HelloHandler hit")
	proxy := h.RouteMap[r.URL.Path]
	proxy.ServeHTTP(w, r)
}
