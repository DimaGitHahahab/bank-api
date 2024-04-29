package handlers

import (
	"bank-api/internal/service"
)

type Handler struct {
	us service.UserService
	ac service.AccountService
	tr service.TransactionService

	JwtSecret string
}

func NewHandler(jwtSecrete string, us service.UserService, as service.AccountService, tr service.TransactionService) *Handler {
	return &Handler{
		us:        us,
		ac:        as,
		tr:        tr,
		JwtSecret: jwtSecrete,
	}
}
