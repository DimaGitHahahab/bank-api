package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Jwt struct {
	Secret string
}

func (j *Jwt) RequireAuth(c *gin.Context) {
	tokenStr, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(j.Secret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if time.Now().Unix() > int64(claims["exp"].(float64)) {
			c.AbortWithStatus(http.StatusUnauthorized)
		} else {
			c.Set("user_id", claims["sub"])
		}
	}

	c.Next()
}
