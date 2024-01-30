package handlers

import (
	"bank-api/internal/bank"
	"bank-api/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type newAccountRequest struct {
	CurrencyName string `json:"currency_name" binding:"required"`
}

type accountInfoResponse struct {
	Id           int    `json:"id"`
	CurrencyName string `json:"currency_name"`
	Amount       int    `json:"amount"`
}

func NewAccount(bank *bank.AccountService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request (no id in context)"})
			return
		}
		id := int(userId.(float64))

		var req newAccountRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}

		cur := model.Currency(req.CurrencyName)
		account, err := (*bank).CreateAccount(c, id, cur)
		if err != nil {
			code, message := handleError(err)
			c.JSON(code, gin.H{"message": message})
			return
		}

		resp := accountInfoResponse{
			Id:           account.Id,
			CurrencyName: string(account.Cur),
			Amount:       account.Amount,
		}

		c.JSON(http.StatusOK, resp)
	}
}

type getAccountRequest struct {
	AccountId int `json:"id" binding:"required"`
}

func GetAccount(bank *bank.AccountService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request (no id in context)"})
			return
		}
		id := int(userId.(float64))

		var req getAccountRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}

		account, err := (*bank).GetAccount(c, id, req.AccountId)
		if err != nil {
			code, message := handleError(err)
			c.JSON(code, gin.H{"message": message})

		}

		resp := accountInfoResponse{
			Id:           account.Id,
			CurrencyName: string(account.Cur),
			Amount:       account.Amount,
		}

		c.JSON(http.StatusOK, resp)
	}
}

type deleteAccountRequest struct {
	AccountId int `json:"id" binding:"required"`
}

func DeleteAccount(bank *bank.AccountService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request (no id in context)"})
			return
		}
		id := int(userId.(float64))

		var req deleteAccountRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}
		err := (*bank).DeleteAccount(c, id, req.AccountId)
		if err != nil {
			code, message := handleError(err)
			c.JSON(code, gin.H{"message": message})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	}
}
