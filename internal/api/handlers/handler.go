package handlers

import (
	"bank-api/internal/domain"
	"bank-api/internal/service"
	"errors"
	"net/http"
)

type Handler struct {
	us service.UserService
	ac service.AccountService
	tr service.TransactionService
}

func NewHandler(us service.UserService, as service.AccountService, tr service.TransactionService) *Handler {
	return &Handler{
		us: us,
		ac: as,
		tr: tr,
	}
}

func handleError(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrInvalidAccount):
		return http.StatusBadRequest, "Invalid account"
	case errors.Is(err, domain.ErrNoSuchAccount):
		return http.StatusNotFound, "No such account"
	case errors.Is(err, domain.ErrNoSuchCurrency):
		return http.StatusNotFound, "No such currency"
	case errors.Is(err, domain.ErrNotEnoughMoney):
		return http.StatusForbidden, "Not enough money"
	case errors.Is(err, domain.ErrInvalidAmount):
		return http.StatusForbidden, "Invalid amount"
	case errors.Is(err, domain.ErrUserAlreadyExists):
		return http.StatusConflict, "User already exists"
	case errors.Is(err, domain.ErrNoSuchUser):
		return http.StatusNotFound, "No such user"
	case errors.Is(err, domain.ErrUserAlreadyExists):
		return http.StatusConflict, "User already exists"
	case errors.Is(err, domain.ErrInvalidEmail):
		return http.StatusBadRequest, "Invalid email"
	case errors.Is(err, domain.ErrEmptyPassword):
		return http.StatusBadRequest, "Empty password"
	case errors.Is(err, domain.ErrEmptyUserInfo):
		return http.StatusBadRequest, "Empty user info"
	default:
		return http.StatusInternalServerError, "Internal server error"
	}
}
