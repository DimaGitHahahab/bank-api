package handlers

import (
	"bank-api/internal/bank"
	"errors"
	"net/http"
)

func handleError(err error) (int, string) {
	switch {
	case errors.Is(err, bank.ErrInvalidAccount):
		return http.StatusBadRequest, "Invalid account"
	case errors.Is(err, bank.ErrNoSuchCurrency):
		return http.StatusBadRequest, "No such currency"
	case errors.Is(err, bank.ErrNoSuchAccount):
		return http.StatusBadRequest, "No such account"
	case errors.Is(err, bank.ErrNotEnoughMoney):
		return http.StatusForbidden, "Not enough money"
	case errors.Is(err, bank.ErrUserAlreadyExists):
		return http.StatusConflict, "User already exists"
	case errors.Is(err, bank.ErrInvalidAmount):
		return http.StatusBadRequest, "Invalid amount"
	case errors.Is(err, bank.ErrNoTransactions):
		return http.StatusNotFound, "No transactions"
	default:
		return http.StatusInternalServerError, "Internal server error"
	}
}
