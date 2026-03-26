package proxy

import (
	"net/http/httputil"
	"net/url"
)

func NewReverseProxy(targetRaw string, backendPath string) (*httputil.ReverseProxy, error) {
	targetURL, err := url.Parse(targetRaw)
	if err != nil {
		return nil, err
	}

	proxy := &httputil.ReverseProxy{
		Rewrite: func(r *httputil.ProxyRequest) {
			r.SetURL(targetURL)
			r.SetXForwarded()
			r.Out.URL.Path = backendPath 
		},
	}

	return proxy, nil
}
