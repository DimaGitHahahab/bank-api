package handlers

import (
	"net/http"
	"strconv"

	"bank-api/internal/domain"

	"github.com/gin-gonic/gin"
)

type depositRequest struct {
	Amount int `json:"amount" binding:"required"`
}

func (h *Handler) Deposit() gin.HandlerFunc {
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
		var req depositRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}

		err = h.tr.ProcessTransaction(c, &domain.Transaction{
			UserId:      id,
			ToAccountId: accountId,
			Amount:      req.Amount,
			Type:        domain.Deposit,
		})
		if err != nil {
			code, message := handleError(err)
			c.JSON(code, gin.H{"message": message})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

type withdrawRequest struct {
	Amount int `json:"amount" binding:"required"`
}

func (h *Handler) Withdraw() gin.HandlerFunc {
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

		var req withdrawRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}

		err = h.tr.ProcessTransaction(c, &domain.Transaction{
			UserId:        id,
			FromAccountId: accountId,
			Amount:        req.Amount,
			Type:          domain.Withdraw,
		})
		if err != nil {
			code, message := handleError(err)
			c.JSON(code, gin.H{"message": message})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

type transferRequest struct {
	FromAccountId int `json:"from_account_id" binding:"required"`
	ToAccountId   int `json:"to_account_id" binding:"required"`
	Amount        int `json:"amount" binding:"required"`
}

func (h *Handler) Transfer() gin.HandlerFunc {
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

		err := h.tr.ProcessTransaction(c, &domain.Transaction{
			UserId:        id,
			FromAccountId: req.FromAccountId,
			ToAccountId:   req.ToAccountId,
			Amount:        req.Amount,
			Type:          domain.Transfer,
		})
		if err != nil {
			code, message := handleError(err)
			c.JSON(code, gin.H{"message": message})
			return
		}

		c.Status(http.StatusNoContent)
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

func (h *Handler) ListTransactions() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}
		id := int(userId.(float64))

		transactions, err := h.tr.ListTransactions(c, id)
		if err != nil {
			code, message := handleError(err)
			c.JSON(code, gin.H{"message": message})
			return
		}

		if len(transactions) == 0 {
			c.Status(http.StatusNoContent)
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
