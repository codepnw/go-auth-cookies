package main

import (
	"log"
	"os"

	"github.com/codepnw/go-auth-cookies/internal/db"
	"github.com/codepnw/go-auth-cookies/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	envPath = "cmd/dev.env"
	version = "v1"
)

func main() {
	postgres, err := db.NewPostgresConnect(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}

	redis := db.NewRedis()
	router := gin.Default()

	routes.NewRoutes(version, router, postgres, redis)

	port := os.Getenv("APP_PORT")
	log.Println("server starting on port: ", port)

	router.Run(":" + port)
}

func init() {
	if err := godotenv.Load(envPath); err != nil {
		log.Fatal("load env file failed: %w", err)
	}
}
