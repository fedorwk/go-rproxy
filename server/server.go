package server

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/fedorwk/ngrok-single-endpoint/srvlist"
	"github.com/gin-gonic/gin"
)

type Server struct {
	srv    *http.Server
	client *http.Client

	services srvlist.ServiceList

	RequestContext context.Context
}

func New(srv *http.Server, client *http.Client, srvlist srvlist.ServiceList) *Server {
	server := &Server{
		srv:    srv,
		client: client,

		services: srvlist,

		RequestContext: context.TODO(),
	}
	return server
}

func (s *Server) Run() error {
	s.setupServiceRouter()
	err := s.srv.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) setupServiceRouter() {
	router := gin.Default()
	for name, addr := range s.services {
		service := router.Group("/" + name)
		service.Any("/", gin.WrapH(s.newServiceHandler(name, addr)))
	}
	s.srv.Handler = router
}

func (s *Server) newServiceHandler(name string, addr string) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		// TODO: This way doesn't work. Try to inspect and use http/httputil/ReverseProxy

		path := strings.TrimPrefix(r.URL.Path, "/"+name)
		uri := addr + path
		serviceRequest := r.Clone(s.RequestContext)
		serviceRequest.RequestURI = uri
		serviceRequest.URL.Path = path

		response, err := s.client.Do(serviceRequest)
		if err != nil {
			log.Println(err)
			return
		}
		err = response.Write(w)
		if err != nil {
			log.Println(err)
			return
		}
	}
	return http.HandlerFunc(handlerFunc)
}
