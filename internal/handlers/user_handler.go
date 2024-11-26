package handlers

import (
	"net/http"
	"time"

	database "github.com/codepnw/go-auth-cookies/internal/db/migrations"
	"github.com/codepnw/go-auth-cookies/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *handlerConfig) SignupHandler(c *gin.Context) {
	var req models.UserRegisterReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	newUser, err := h.db.CreateUser(c, database.CreateUserParams{
		ID:        uuid.New(),
		Name:      req.Name,
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, newUser)
}
