package users

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	utils "github.com/ntu-onemdp/onemdp-backend/utils"
)

// Users from JSON reqeust
type Users struct {
	Usernames []string `json:"usernames" binding:"required"`
}

// User struct to be inserted into DB
type UserRow struct {
	username string
	password string
	role     string
	date     time.Time
	status   string
}

func CreateUsers(c *gin.Context, pool *pgxpool.Pool) {
	utils.Logger.Trace().Msg("Create users function called")

	var users Users
	if err := c.BindJSON(&users); err != nil {
		utils.Logger.Error().Err(err).Msg("Error binding")
	}

	// Create file to store passwords
	file, err := os.OpenFile("new_users.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error creating file for password storage")
	}
	defer file.Close()

	// Initialize rows for batch insert (using copy function)
	rows := []UserRow{}

	for i := range users.Usernames {
		username := users.Usernames[i]
		password := utils.GeneratePassword()
		date := time.Now()

		// Write username and password to file
		file.WriteString(fmt.Sprintf("%s,%s\n", username, password))

		// h(p)
		// Note: Set custom params for prod
		stored_hashed_pw, err := argon2id.CreateHash(password, argon2id.DefaultParams)
		if err != nil {
			utils.Logger.Error().Err(err).Msg("Error hashing password")
		}

		// Add to row
		row := UserRow{
			username: username,
			password: string(stored_hashed_pw),
			role:     "student",
			date:     date,
			status:   "active",
		}
		rows = append(rows, row)
	}

	// Perform bulk insert to DB
	num_success, err := pool.CopyFrom(context.Background(), pgx.Identifier{"users"}, []string{"username", "password", "role", "date_created", "status"}, pgx.CopyFromSlice(len(rows), func(i int) ([]any, error) {
		return []any{rows[i].username, rows[i].password, rows[i].role, rows[i].date, rows[i].status}, nil
	}))
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error performing bulk insert to DB")
	}

	utils.Logger.Info().Msg(fmt.Sprintf("%d rows successfully inserted", num_success))
}
