package router

import (
	"bank-api/internal/api/handlers"
	"bank-api/internal/api/middleware"
	"bank-api/internal/service"
	"github.com/gin-gonic/gin"
)

func NewRouter(users service.UserService, accounts service.AccountService, transactions service.TransactionService) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.RateLimiter(2))

	r.POST("/user/signup", handlers.SignUp(&users))
	r.POST("/user/login", handlers.Login(&users))

	auth := r.Group("/")
	auth.Use(middleware.RequireAuth)
	{
		auth.GET("user", handlers.GetUser(&users))
		auth.PATCH("user", handlers.UpdateUser(&users))
		auth.DELETE("user", handlers.DeleteUser(&users))

		auth.POST("account", handlers.NewAccount(&accounts))
		auth.GET("account", handlers.GetAccount(&accounts))
		auth.DELETE("account", handlers.DeleteAccount(&accounts))

		auth.POST("account/deposit", handlers.Deposit(&transactions))
		auth.POST("account/withdraw", handlers.Withdraw(&transactions))

		auth.POST("account/transfer", handlers.Transfer(&transactions))

		auth.GET("account/history", handlers.ListTransactions(&transactions))
	}

	return r
}
