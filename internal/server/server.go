package server

import (
	"context"
	"net"
	"net/http"
)

type Server struct {
	http *http.Server
}

func New(h http.Handler) *Server {
	return &Server{http: &http.Server{Handler: h}}
}

func (s *Server) Run(port string) error {
	s.http.Addr = net.JoinHostPort("", port)
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
