package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
	"github.com/pressly/goose/v3"
)

var Pool *pgxpool.Pool

func Init() {
	// Retrieve DB password from secrets
	db_pw, err := os.ReadFile("/run/secrets/db-password")
	if err != nil {
		utils.Logger.Warn().Msg("Error reading secret, attempting to load from .env.dev")

		// Try reading from .env
		if err := godotenv.Load("config/.env.dev"); err != nil {
			utils.Logger.Panic().Err(err).Msg("Error reading from .env")
		}
		db_pw = []byte(os.Getenv("POSTGRES_PW"))
	}

	// Retrieve env variables
	postgres_db, exists := os.LookupEnv("POSTGRES_DB")
	if !exists {
		// Defaults to DEV_1
		utils.Logger.Warn().Msg("Error retrieving postgres database name, default name set.")
		postgres_db = "dev_1"
	}

	// Create connection pool to db
	pg_username := os.Getenv("PG_USERNAME")

	// IMPORTANT
	// Uncomment the correct connection string based on the environment
	// // Local
	// connectionString := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable", pg_username, string(db_pw), postgres_db)
	// Docker run
	connectionString := fmt.Sprintf("postgres://%s:%s@host.docker.internal:5432/%s?sslmode=disable", pg_username, string(db_pw), postgres_db)
	// // Docker compose
	// connectionString := fmt.Sprintf("postgres://%s:%s@db:5432/%s?sslmode=disable", string(db_pw), postgres_db)

	Pool, err = pgxpool.New(context.Background(), connectionString)
	if err != nil {
		utils.Logger.Panic().Err(err).Msg("Error creating connection pool")
	}

	// Initialize Goose and perform migrations
	if err := goose.SetDialect("postgres"); err != nil {
		utils.Logger.Panic().Err(err)
	}

	db := stdlib.OpenDBFromPool(Pool)

	// Check migration status
	if err := goose.Status(db, "migrations"); err != nil {
		utils.Logger.Panic().Err(err).Msg("Error checking migration status")
	}

	if err := goose.Up(db, "migrations"); err != nil {
		utils.Logger.Panic().Err(err)
	}
	if err := db.Close(); err != nil {
		utils.Logger.Panic().Err(err)
	}

	utils.Logger.Info().Msg("Database connection pool successfully created.")
}

func Close() {
	Pool.Close()
	utils.Logger.Info().Msg("Database connection pool successfully terminated.")
}
