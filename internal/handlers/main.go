package handlers

import (
	"net/http"

	database "github.com/codepnw/go-auth-cookies/internal/db/migrations"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type handlerConfig struct {
	db          *database.Queries
	redisClient *redis.Client
}

func NewHandler(db *database.Queries, redisClient *redis.Client) *handlerConfig {
	return &handlerConfig{
		db:          db,
		redisClient: redisClient,
	}
}

func (h *handlerConfig) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
