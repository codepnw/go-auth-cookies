package routes

import (
	database "github.com/codepnw/go-auth-cookies/internal/db/migrations"
	"github.com/codepnw/go-auth-cookies/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func NewRoutes(version string, router *gin.Engine, db *database.Queries, redis *redis.Client) {
	handler := handlers.NewHandler(db, redis)
	group := router.Group(version + "/")

	group.GET("/healthcheck", handler.HealthCheck)
}