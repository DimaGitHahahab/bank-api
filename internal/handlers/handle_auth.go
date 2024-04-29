package handlers

import (
	"net/http"
	"time"

	"bank-api/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const expireDeadline = time.Hour * 24

type signUpRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type userResponse struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *Handler) SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req signUpRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			returnBadRequest(c)
			return
		}

		user, err := h.us.CreateUser(c, &domain.UserInfo{
			Name:     req.Name,
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			returnError(c, err)
			return
		}

		c.JSON(http.StatusCreated, userResponse{
			Id:        user.Id,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		})
	}
}

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req loginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			returnBadRequest(c)
			return
		}

		u, err := h.us.AuthenticateUser(c, &domain.UserInfo{
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			returnError(c, err)
			return
		}

		t, err := h.generateJWT(u)
		if err != nil {
			returnError(c, err)
			return
		}

		setCookieToken(c, t)
		c.Status(http.StatusOK)
	}
}

func (h *Handler) generateJWT(u *domain.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": u.Id,
		"exp": time.Now().Add(expireDeadline).Unix(),
	})

	return token.SignedString([]byte(h.JwtSecret))
}

func setCookieToken(c *gin.Context, token string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", token, 3600*24, "", "", true, true)
}
