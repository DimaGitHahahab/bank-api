package handlers

import (
	"net/http"

	"bank-api/internal/domain"

	"github.com/gin-gonic/gin"
)

type newAccountRequest struct {
	CurrencyName string `json:"currency_name" binding:"required"`
}

type accountInfoResponse struct {
	Id           int    `json:"id"`
	CurrencyName string `json:"currency_name"`
	Amount       int    `json:"amount"`
}

func (h *Handler) NewAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var id int
		if ok := getUserId(c, &id); !ok {
			returnBadRequest(c)
			return
		}

		var req newAccountRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			returnBadRequest(c)
			return
		}

		cur := domain.Currency{Symbol: req.CurrencyName}
		account, err := h.ac.CreateAccount(c, id, cur)
		if err != nil {
			returnError(c, err)
			return
		}

		c.JSON(http.StatusOK, accountInfoResponse{
			Id:           account.Id,
			CurrencyName: account.Cur.Symbol,
			Amount:       account.Amount,
		})
	}
}

func (h *Handler) GetAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var id int
		if ok := getUserId(c, &id); !ok {
			returnBadRequest(c)
			return
		}

		var accountId int
		if ok := getAccountId(c, &accountId); !ok {
			returnBadRequest(c)
			return
		}

		account, err := h.ac.GetAccount(c, id, accountId)
		if err != nil {
			returnError(c, err)
			return
		}

		c.JSON(http.StatusOK, accountInfoResponse{
			Id:           account.Id,
			CurrencyName: account.Cur.Symbol,
			Amount:       account.Amount,
		})
	}
}

func (h *Handler) DeleteAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		var id int
		if ok := getUserId(c, &id); !ok {
			returnBadRequest(c)
			return
		}

		var accountId int
		if ok := getAccountId(c, &accountId); !ok {
			returnBadRequest(c)
			return
		}

		if err := h.ac.DeleteAccount(c, id, accountId); err != nil {
			returnError(c, err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}
