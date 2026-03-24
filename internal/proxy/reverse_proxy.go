package proxy

import (
	"net/http/httputil"
	"net/url"
)

func NewReverseProxy() (*httputil.ReverseProxy, error) {
	targetURL, err := url.Parse("http://localhost:8081")
	if err != nil {
		return nil, err
	}

	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(targetURL)
			r.SetXForwarded()
		},
	}

	return proxy, nil
}
