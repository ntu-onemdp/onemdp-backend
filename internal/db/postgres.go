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
	// Get app env
	env, found := os.LookupEnv("ENV")
	if !found {
		// Default environment: PROD
		env = "PROD"
	}
	utils.Logger.Info().Str("environment", env).Msgf("App environment loaded: %s", env)

	// Retrieve DB password from secrets
	db_pw, err := os.ReadFile("/run/secrets/db-password")
	if err != nil {
		utils.Logger.Warn().Msg("Error reading secret, attempting to load from .env")

		// Try reading from .env
		switch env {
		case "PROD":
			if err := godotenv.Load("config/.env"); err != nil {
				utils.Logger.Panic().Err(err).Msg("Error reading from .env")
			}
		case "QA":
			if err := godotenv.Load("config/.env.qa"); err != nil {
				utils.Logger.Panic().Err(err).Msg("Error reading from .env.qa")
			}
		case "DEV":
			if err := godotenv.Load("config/.env.dev"); err != nil {
				utils.Logger.Panic().Err(err).Msg("Error reading from .env.dev")
			}
		default:
			if err := godotenv.Load("config/.env"); err != nil {
				utils.Logger.Panic().Err(err).Msg("Error reading from .env")
			}
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

	netloc, exists := os.LookupEnv("POSTGRES_NETLOC")
	if !exists {
		// Defaults to host.docker.internal
		utils.Logger.Warn().Msg("Error retrieving postgres netloc, default netloc set.")
		netloc = "host.docker.internal"
	}

	pg_username := os.Getenv("PG_USERNAME")
	pg_port := os.Getenv("PG_PORT")

	// IMPORTANT
	// Set the correct host for the database depending on where the application is running. Set it in .env
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require", pg_username, string(db_pw), netloc, pg_port, postgres_db)
	utils.Logger.Debug().Str("connection string", connectionString).Msg("")

	// Create connection pool to db
	Pool, err = pgxpool.New(context.Background(), connectionString)
	if err != nil {
		utils.Logger.Panic().Err(err).Msg("Error creating connection pool")
	}
	utils.Logger.Debug().Msg("Postgres connection pool created")

	// Initialize Goose and perform migrations
	if err := goose.SetDialect("postgres"); err != nil {
		utils.Logger.Panic().Err(err)
	}

	db := stdlib.OpenDBFromPool(Pool)

	// Perform migrations
	if err := goose.Up(db, "migrations"); err != nil {
		utils.Logger.Panic().Err(err)
	}
	utils.Logger.Debug().Msg("Goose migrations applied")

	// Check migration status
	if err := goose.Status(db, "migrations"); err != nil {
		utils.Logger.Panic().Err(err).Msg("Error checking migration status")
	}
	if err := db.Close(); err != nil {
		utils.Logger.Panic().Err(err)
	}

	utils.Logger.Info().Msg("Postgres database fully initialized.")
}

func Close() {
	Pool.Close()
	utils.Logger.Info().Msg("Database connection pool successfully terminated.")
}
