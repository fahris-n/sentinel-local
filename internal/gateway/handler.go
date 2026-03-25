package gateway

import (
	"log"
	"net/http"
	"net/http/httputil"
)

type Handler struct {
	HelloProxy *httputil.ReverseProxy
}

func NewHandler(helloProxy *httputil.ReverseProxy) *Handler {
	return &Handler{
		HelloProxy: helloProxy,
	}
}

func (h *Handler) HelloHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("HelloHandler hit")
	h.HelloProxy.ServeHTTP(w, r)
}
