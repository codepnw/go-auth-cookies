package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	database "github.com/codepnw/go-auth-cookies/internal/db/migrations"
)

func NewPostgresConnect(dbURL string) (*database.Queries, error) {
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	fmt.Println("database connected...")
	return database.New(conn), nil
}
