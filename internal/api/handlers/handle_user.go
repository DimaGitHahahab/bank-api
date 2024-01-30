package handlers

import (
	"bank-api/internal/bank"
	"bank-api/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUser(bank *bank.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {

		userId, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}
		id := int(userId.(float64))
		user, err := (*bank).GetUserById(c, id)
		if err != nil {
			code, message := handleError(err)
			c.JSON(code, gin.H{"message": message})
			return
		}

		c.JSON(http.StatusOK, userInfoResponse{
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

func UpdateUser(bank *bank.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}
		id := int(userId.(float64))
		var req updateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}

		account, err := (*bank).UpdateUserInfo(c, id, &model.UserInfo{
			Name:  req.Name,
			Email: req.Email,
		})
		if err != nil {
			code, message := handleError(err)
			c.JSON(code, gin.H{"message": message})
			return
		}

		c.JSON(http.StatusOK, userInfoResponse{
			Id:        account.Id,
			Name:      account.Name,
			Email:     account.Email,
			CreatedAt: account.CreatedAt,
		})
	}
}

func DeleteUser(bank *bank.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
			return
		}
		id := int(userId.(float64))

		if err := (*bank).DeleteUserById(c, id); err != nil {
			code, message := handleError(err)
			c.JSON(code, gin.H{"message": message})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	}
}
