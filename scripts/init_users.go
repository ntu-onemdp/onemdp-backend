package main

import (
	"context"
	"fmt"
	"os"

	"github.com/alexedwards/argon2id"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
	"github.com/pressly/goose/v3"
)

// Script to force initialize 1 user into the base. Use this to initialize the admin.
func main() {
	// Command line arguments (3): username, name, password
	username := os.Args[1]
	name := os.Args[2]
	password := os.Args[3]

	// Load dotenv file
	if err := godotenv.Load("config/.env.dev"); err != nil {
		utils.Logger.Panic().Err(err).Msg("Error reading from .env")
	}
	db_pw := os.Getenv("POSTGRES_PW")
	pg_db := os.Getenv("POSTGRES_DB")
	pg_username := os.Getenv("PG_USERNAME")

	connectionString := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable", pg_username, db_pw, pg_db)
	dbpool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		utils.Logger.Panic().Err(err).Msg("Error creating connection pool")
	}
	defer dbpool.Close()

	// Initialize Goose and perform migrations
	if err := goose.SetDialect("postgres"); err != nil {
		utils.Logger.Panic().Err(err)
	}

	stored_hashed_pw, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error hashing password")
	}

	// Add to users
	res, err := dbpool.Exec(context.Background(), "insert into users values ($1, $2, now())", username, name)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error inserting into users table")
	} else {
		utils.Logger.Trace().Msg(res.String())
	}

	// Add to auth
	res, err = dbpool.Exec(context.Background(), "insert into auth values ($1, $2, 'student')", username, stored_hashed_pw)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error inserting into auth table")
	} else {
		utils.Logger.Trace().Msg(res.String())
	}

	utils.Logger.Info().Msg("Script completed")
}
