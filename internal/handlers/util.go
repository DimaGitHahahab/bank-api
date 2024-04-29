package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"bank-api/internal/domain"

	"github.com/gin-gonic/gin"
)

func returnBadRequest(c *gin.Context) {
	c.String(http.StatusBadRequest, "invalid request")
}

func returnError(c *gin.Context, err error) {
	code, msg := getCodeAndMessage(err)
	c.String(code, msg)
}

func getUserId(c *gin.Context, id *int) bool {
	userIdClaim, ok := c.Get("user_id")
	if !ok {
		return false
	}
	*id = int(userIdClaim.(float64))
	return true
}

func getAccountId(c *gin.Context, id *int) bool {
	accountId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return false
	}

	*id = accountId
	return true
}

func getCodeAndMessage(err error) (int, string) {
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
	case errors.Is(err, domain.ErrWrongPassword):
		return http.StatusUnauthorized, "Wrong password"
	default:
		return http.StatusInternalServerError, "Internal server error"
	}
}
