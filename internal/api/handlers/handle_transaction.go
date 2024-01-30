package handlers

import (
	"bank-api/internal/bank"
	"bank-api/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type depositRequest struct {
	AccountId int `json:"id" binding:"required"`
	Amount    int `json:"amount" binding:"required"`
}

func Deposit(bank *bank.TransactionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request (no id in context)"})
			return
		}
		id := int(userId.(float64))
		var req depositRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}

		err := (*bank).ProcessTransaction(c, &model.Transaction{
			UserId:    id,
			AccountId: req.AccountId,
			Amount:    req.Amount,
			Type:      model.Deposit,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	}
}

type withdrawRequest struct {
	AccountId int `json:"id" binding:"required"`
	Amount    int `json:"amount" binding:"required"`
}

func Withdraw(bank *bank.TransactionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request (no id in context)"})
			return
		}
		id := int(userId.(float64))
		var req withdrawRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}

		err := (*bank).ProcessTransaction(c, &model.Transaction{
			UserId:    id,
			AccountId: req.AccountId,
			Amount:    req.Amount,
			Type:      model.Withdraw,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	}
}

type transferRequest struct {
	FromAccountId int `json:"from_account_id" binding:"required"`
	ToAccountId   int `json:"to_account_id" binding:"required"`
	Amount        int `json:"amount" binding:"required"`
}

func Transfer(bank *bank.TransactionService) gin.HandlerFunc {
	return func(c *gin.Context) {

		userId, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request (no id in context)"})
			return
		}
		id := int(userId.(float64))
		var req transferRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}

		err := (*bank).ProcessTransfer(c, &model.Transfer{
			UserId:      id,
			AccountId:   req.FromAccountId,
			ToAccountId: req.ToAccountId,
			Amount:      req.Amount,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	}
}
