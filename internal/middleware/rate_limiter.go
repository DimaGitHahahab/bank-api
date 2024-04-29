package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func RateLimiter(maxRequests int) gin.HandlerFunc {
	rate := limiter.Rate{
		Period: 1 * time.Second,
		Limit:  int64(maxRequests),
	}
	store := memory.NewStore()
	instance := limiter.New(store, rate)

	return func(c *gin.Context) {
		context, err := instance.Get(c, c.ClientIP())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			return
		}

		if context.Reached {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "too many requests"})
			return
		}

		c.Next()
	}
}
