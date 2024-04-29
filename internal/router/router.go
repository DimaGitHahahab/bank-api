package router

import (
	"bank-api/internal/handlers"
	"bank-api/internal/middleware"

	"github.com/gin-gonic/gin"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(h *handlers.Handler) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.RateLimiter(1000))

	r.POST("/user/signup", h.SignUp())
	r.POST("/user/login", h.Login())

	auth := r.Group("/")

	jwt := middleware.Jwt{Secret: h.JwtSecret}
	auth.Use(jwt.RequireAuth)
	{
		auth.GET("user", h.GetUser())
		auth.PUT("user", h.UpdateUser())
		auth.PATCH("user", h.UpdateUser())
		auth.DELETE("user", h.DeleteUser())

		auth.POST("account", h.NewAccount())
		auth.GET("account/:id", h.GetAccount())
		auth.DELETE("account/:id", h.DeleteAccount())

		auth.POST("account/:id/deposit", h.Deposit())
		auth.POST("account/:id/withdraw", h.Withdraw())

		auth.POST("account/transfer", h.Transfer())

		auth.GET("history", h.ListTransactions())
	}

	r.GET("/swagger/*any", gin.WrapH(httpSwagger.Handler(
		httpSwagger.URL("/swagger.yaml"), // The url pointing to API definition
	)))
	r.StaticFile("/swagger.yaml", "./docs/openapi.yaml")

	return r
}
