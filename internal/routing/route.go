package routing

import "net/http/httputil"

type RouteEntry struct {
	Proxy        *httputil.ReverseProxy
	RequiresAuth bool
	AllowedRoles []string
}
