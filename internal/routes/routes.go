package routes

import (
	database "github.com/codepnw/go-auth-cookies/internal/db/migrations"
	"github.com/codepnw/go-auth-cookies/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Version string
	Router  *gin.Engine
	DB      *database.Queries
	Redis   *redis.Client
}

func NewRoutes(cfg *Config) {
	handler := handlers.NewHandler(cfg.DB, cfg.Redis)

	middleware := cfg.Router.Group(cfg.Version, handler.AuthMiddleware())
	users := cfg.Router.Group(cfg.Version + "/users")

	middleware.GET("/healthcheck", handler.HealthCheck)

	users.POST("/signin", handler.SignInHandler)
	users.GET("/logout", handler.LogoutHandler)
	users.POST("/signup", handler.SignupHandler)
}
