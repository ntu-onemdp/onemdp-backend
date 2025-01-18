package auth

import (
	"context"

	"github.com/alexedwards/argon2id"
	utils "github.com/ntu-onemdp/onemdp-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LoginForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"` // Plaintext password
}

type LoginResponse struct {
	// Status   string `json:"status"` // To be removed
	Success  bool   `json:"success"`
	ErrorMsg string `json:"error_msg"`
	Jwt      string `json:"jwt"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Name     string `json:"name"`
}

// Handle log in requests from frontend.
func HandleLogin(c *gin.Context, pool *pgxpool.Pool) {
	utils.Logger.Info().Msg("Login request received")
	var form LoginForm

	// Bind with form
	if err := c.ShouldBind(&form); err != nil {
		utils.Logger.Error().Err(err).Msg("Error processing login request")
		response := LoginResponse{
			Success:  false,
			ErrorMsg: "Malformed request",
		}

		c.JSON(200, &response)
		return
	}

	// Query database. The query fails if the status is not 'active' (deleted or inactive user).
	var username string
	var password string // Stored hashed password
	err := pool.QueryRow(context.Background(), "SELECT username, password from users where username=$1 and status='active'", form.Username).Scan(&username, &password)
	if err != nil {
		// Check logs. If error=="no rows in result set", the username does not exist in the database.
		utils.Logger.Error().Err(err).Msg("Error fetching username from database. ")
		response := LoginResponse{
			Success:  false,
			ErrorMsg: "Unauthorized: Incorrect username/password",
		}

		c.JSON(200, &response)
		return
	}

	// Authenticate user
	match, err := argon2id.ComparePasswordAndHash(form.Password, password)
	if !match && err != nil {
		utils.Logger.Trace().Msg("Invalid login attempt")
		response := LoginResponse{
			Success:  false,
			ErrorMsg: "Unauthorized: Incorrect username/password",
		}
		c.JSON(200, &response)
	} else {
		// Retrieve name and role
		var role string
		var name string
		if err = pool.QueryRow(context.Background(), "SELECT role, name FROM users where username=$1 and status='active' and password=$2", username, password).Scan(&role, &name); err != nil {
			utils.Logger.Error().Err(err).Msg("Error retrieving user details")
			response := LoginResponse{
				Success:  false,
				ErrorMsg: "Unexpected authentication error.",
			}
			c.JSON(200, &response)
			return
		}

		// Generate JWT
		tokenString, err := GenerateJwt(UserClaim{username, role})
		if err != nil {
			c.JSON(500, "Internal server error")
		}

		response := LoginResponse{
			Success:  true,
			Jwt:      tokenString,
			Username: username,
			Role:     role,
			Name:     name,
		}
		c.JSON(200, &response)
		return
	}
}
