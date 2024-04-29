package handlers

import (
	"net/http"

	"bank-api/internal/domain"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var id int
		if ok := getUserId(c, &id); !ok {
			returnBadRequest(c)
			return
		}

		user, err := h.us.GetUserById(c, id)
		if err != nil {
			returnError(c, err)
			return
		}

		c.JSON(http.StatusOK, userResponse{
			Id:        user.Id,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		})
	}
}

type updateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *Handler) UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var id int
		if ok := getUserId(c, &id); !ok {
			returnBadRequest(c)
			return
		}

		var req updateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			returnBadRequest(c)
			return
		}

		account, err := h.us.UpdateUserInfo(c, id, &domain.UserInfo{
			Name:  req.Name,
			Email: req.Email,
		})
		if err != nil {
			returnError(c, err)
			return
		}

		c.JSON(http.StatusOK, userResponse{
			Id:        account.Id,
			Name:      account.Name,
			Email:     account.Email,
			CreatedAt: account.CreatedAt,
		})
	}
}

func (h *Handler) DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var id int
		if ok := getUserId(c, &id); !ok {
			returnBadRequest(c)
			return
		}

		if err := h.us.DeleteUserById(c, id); err != nil {
			returnError(c, err)
			return
		}

		c.Status(http.StatusNoContent)
	}
}
