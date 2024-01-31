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
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}
		id := int(userId.(float64))
		var req depositRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}

		err := (*bank).ProcessTransaction(c, &model.Transaction{
			UserId:      id,
			ToAccountId: req.AccountId,
			Amount:      req.Amount,
			Type:        model.Deposit,
		})
		if err != nil {
			code, message := handleError(err)
			c.JSON(code, gin.H{"message": message})
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
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}
		id := int(userId.(float64))
		var req withdrawRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}

		err := (*bank).ProcessTransaction(c, &model.Transaction{
			UserId:        id,
			FromAccountId: req.AccountId,
			Amount:        req.Amount,
			Type:          model.Withdraw,
		})
		if err != nil {
			code, message := handleError(err)
			c.JSON(code, gin.H{"message": message})
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
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}
		id := int(userId.(float64))
		var req transferRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}

		err := (*bank).ProcessTransaction(c, &model.Transaction{
			UserId:        id,
			FromAccountId: req.FromAccountId,
			ToAccountId:   req.ToAccountId,
			Amount:        req.Amount,
			Type:          model.Transfer,
		})
		if err != nil {
			code, message := handleError(err)
			c.JSON(code, gin.H{"message": message})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	}
}

type listTransactionsResponse struct {
	Transactions []transaction `json:"transactions"`
}

type transaction struct {
	FromAccountId  int    `json:"from_account_id"`
	ToAccountId    int    `json:"to_account_id"`
	CurrencySymbol string `json:"currency_name"`
	Amount         int    `json:"amount"`
	Time           string `json:"processed_at"`
}

func ListTransactions(bank *bank.TransactionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}
		id := int(userId.(float64))

		transactions, err := (*bank).ListTransactions(c, id)
		if err != nil {
			code, message := handleError(err)
			c.JSON(code, gin.H{"message": message})
			return
		}

		var resp listTransactionsResponse
		resp.Transactions = make([]transaction, len(transactions))
		for i := range resp.Transactions {
			resp.Transactions[i] = transaction{
				FromAccountId:  transactions[i].FromAccountId,
				ToAccountId:    transactions[i].ToAccountId,
				CurrencySymbol: transactions[i].Cur.Symbol,
				Amount:         transactions[i].Amount,
				Time:           transactions[i].Time.Format("2006-01-02 15:04:05"),
			}
		}

		c.JSON(http.StatusOK, resp)
	}
}
