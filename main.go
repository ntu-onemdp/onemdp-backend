package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ntu-onemdp/onemdp-backend/users"
	utils "github.com/ntu-onemdp/onemdp-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/ntu-onemdp/onemdp-backend/auth"
	"github.com/pressly/goose/v3"
)

func main() {
	// Retrieve DB password from secrets
	db_pw, err := os.ReadFile("/run/secrets/db-password")
	if err != nil {
		utils.Logger.Error().Msg("Error reading secret")

		// Try reading from .env
		if err := godotenv.Load(".env.dev"); err != nil {
			utils.Logger.Panic().Err(err).Msg("Error reading from .env")
		}
		db_pw = []byte(os.Getenv("POSTGRES_PW"))
	}

	// Retrieve env variables
	postgres_db, exists := os.LookupEnv("POSTGRES_DB")
	if !exists {
		// Defaults to DEV_1
		utils.Logger.Error().Msg("Error retrieving postgres database name, default name set.")
		postgres_db = "dev_1"
	}

	// Create connection pool to db
	connection_string := fmt.Sprintf("postgres://postgres:%s@localhost:5432/%s?sslmode=disable", string(db_pw), postgres_db)
	// Use below if using container
	// connection_string := fmt.Sprintf("postgres://postgres:%s@db:5432/%s?sslmode=disable", string(db_pw), postgres_db)
	dbpool, err := pgxpool.New(context.Background(), connection_string)
	if err != nil {
		utils.Logger.Panic().Err(err).Msg("Error creating connection pool")
	}
	defer dbpool.Close()

	// Initialize Goose and perform migrations
	if err := goose.SetDialect("postgres"); err != nil {
		utils.Logger.Panic().Err(err)
	}

	db := stdlib.OpenDBFromPool(dbpool)
	if err := goose.Up(db, "migrations"); err != nil {
		utils.Logger.Panic().Err(err)
	}
	if err := db.Close(); err != nil {
		utils.Logger.Panic().Err(err)
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"content": "hello",
		})
	})

	r.POST("/api/v1/auth/login", func(c *gin.Context) {
		auth.HandleLogin(c, dbpool)
	})

	r.POST("/api/v1/users/create", func(c *gin.Context) {
		users.CreateUsers(c, dbpool)
	})
	r.Run("localhost:8080") // listen and serve on 0.0.0.0:8080
}
