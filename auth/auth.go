package auth

import (
	"context"
	"fmt"

	"github.com/alexedwards/argon2id"
	utils "github.com/ntu-onemdp/onemdp-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LoginForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"` // Plaintext password
}

type LoginResponse struct {
	Status string `json:"status"`
}

// Handle log in requests from frontend.
func HandleLogin(c *gin.Context, pool *pgxpool.Pool) {
	var form LoginForm
	var loginResponse LoginResponse

	// Bind with form
	if err := c.MustBindWith(&form, binding.FormPost); err != nil {
		utils.Logger.Error().Err(err).Msg("Error processing login request")
		loginResponse.Status = "Malformed request"

		c.JSON(400, &loginResponse)
		return
	}

	utils.Logger.Debug().Msg(fmt.Sprintf("Username: %s, Password: %s", form.Username, form.Password))

	// Query database
	var username string
	var password string // Stored hashed password h2(h1(pw) + s)
	err := pool.QueryRow(context.Background(), "SELECT username, password from users where username=$1 and status='active'", form.Username).Scan(&username, &password)
	if err != nil {
		// Check logs. If error=="no rows in result set", the username does not exist in the database.
		utils.Logger.Error().Err(err).Msg("Error fetching username from database. ")
		loginResponse.Status = "unauthorized"

		c.JSON(401, &loginResponse)
		return
	}

	// Authenticate user
	match, err := argon2id.ComparePasswordAndHash(form.Password, password)
	if !match && err != nil {
		utils.Logger.Trace().Msg("Invalid login attempt")
		loginResponse.Status = "unauthorized"
		c.JSON(401, &loginResponse)
	} else {
		// Retrieve role
		var role string
		if err = pool.QueryRow(context.Background(), "SELECT role FROM users where username=$1 and status='active' and password=$2", username, password).Scan(&role); err != nil {
			utils.Logger.Error().Err(err).Msg("Error retrieving role")
			loginResponse.Status = "unauthorized"

			c.JSON(401, &loginResponse)
			return
		}
		loginResponse.Status = "success"
		c.JSON(200, &loginResponse)
		return
	}
}
