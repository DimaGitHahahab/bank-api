package handlers

import (
	"net/http"
	"strconv"

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
		userId, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}
		id := int(userId.(float64))

		var req newAccountRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}

		cur := domain.Currency{Symbol: req.CurrencyName}
		account, err := h.ac.CreateAccount(c, id, cur)
		if err != nil {
			code, message := handleError(err)
			c.JSON(code, gin.H{"message": message})
			return
		}

		resp := accountInfoResponse{
			Id:           account.Id,
			CurrencyName: account.Cur.Symbol,
			Amount:       account.Amount,
		}

		c.JSON(http.StatusOK, resp)
	}
}

func (h *Handler) GetAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}
		id := int(userId.(float64))

		accountId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}

		account, err := h.ac.GetAccount(c, id, accountId)
		if err != nil {
			code, message := handleError(err)
			c.JSON(code, gin.H{"message": message})

		}

		resp := accountInfoResponse{
			Id:           account.Id,
			CurrencyName: account.Cur.Symbol,
			Amount:       account.Amount,
		}

		c.JSON(http.StatusOK, resp)
	}
}

func (h *Handler) DeleteAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}
		id := int(userId.(float64))

		accountId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}

		err = h.ac.DeleteAccount(c, id, accountId)
		if err != nil {
			code, message := handleError(err)
			c.JSON(code, gin.H{"message": message})
			return
		}

		c.Status(http.StatusNoContent)
	}
}
