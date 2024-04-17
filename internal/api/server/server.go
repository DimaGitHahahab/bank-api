package server

import (
	"context"
	"net/http"

	"bank-api/internal/api/router"
	"bank-api/internal/service"
)

type Server struct {
	http *http.Server
}

func New(users service.UserService, accounts service.AccountService, transactions service.TransactionService) *Server {
	srv := &Server{}

	r := router.NewRouter(users, accounts, transactions)

	srv.http = &http.Server{
		Handler: r,
	}

	return srv
}

func (s *Server) Run(port string) error {
	s.http.Addr = ":" + port
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
