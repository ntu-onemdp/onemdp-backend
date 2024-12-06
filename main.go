package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
)

func main() {
	// Configure logger
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	// Retrieve DB password from secrets
	db_pw, err := os.ReadFile("/run/secrets/db-password")
	if err != nil {
		logger.Error().Msg("Error reading secret")

		// Try reading from .env
		if err := godotenv.Load(".env.dev"); err != nil {
			logger.Panic().Err(err).Msg("Error reading from .env")
		}
		db_pw = []byte(os.Getenv("POSTGRES_PW"))
	}

	// Retrieve env variables
	postgres_db, exists := os.LookupEnv("POSTGRES_DB")
	if !exists {
		// Defaults to DEV_1
		logger.Error().Msg("Error retrieving postgres database name, default name set.")
		postgres_db = "DEV_1"
	}

	// Create connection pool to db
	connection_string := fmt.Sprintf("postgres://postgres:%s@db:5432/%s?sslmode=disable", string(db_pw), postgres_db)
	dbpool, err := pgxpool.New(context.Background(), connection_string)
	if err != nil {
		logger.Panic().Err(err).Msg("Error creating connection pool")
	}
	defer dbpool.Close()

	// Initialize Goose and perform migrations
	if err := goose.SetDialect("postgres"); err != nil {
		logger.Panic().Err(err)
	}

	db := stdlib.OpenDBFromPool(dbpool)
	if err := goose.Up(db, "migrations"); err != nil {
		logger.Panic().Err(err)
	}
	if err := db.Close(); err != nil {
		logger.Panic().Err(err)
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"content": "hello",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
