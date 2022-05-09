package rproxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// New returns modified gin.Engine which redirects requests with leading
// ServiceInfo.Name in path to ServiceInfo.URI() using its adress and protocol
func New(engine *gin.Engine, srvlist ServiceList) (*gin.Engine, error) {
	for _, service := range srvlist {
		serviceRoute := engine.Group("/" + service.Name)
		serviceProxyHandler, err := newServiceProxyHandler(service.URI())
		if err != nil {
			return nil, err
		}
		serviceRoute.Any(
			"*any",
			gin.WrapH(serviceProxyHandler),
		)
	}
	return engine, nil
}

// newSerciceProxyHandler makes redirect Handler with rawURL destination
func newServiceProxyHandler(rawURL string) (http.Handler, error) {
	url, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	serviceProxy := httputil.NewSingleHostReverseProxy(url)
	// Request Modification
	originalDirector := serviceProxy.Director
	serviceProxy.Director = func(req *http.Request) {
		originalDirector(req)
		trimURLPathRoot(req)
	}

	return serviceProxy, nil
}

// Modifies given request by trimming the leading part of request.URL.Path
func trimURLPathRoot(req *http.Request) {
	req.URL.Path = trimPathRoot(req.URL.Path)
}

func trimPathRoot(urlPath string) string {
	if len(urlPath) < 2 {
		return ""
	}
	for i := 1; i < len(urlPath); i++ {
		if urlPath[i] == '/' {
			return urlPath[i:]
		}
	}
	return ""
}
