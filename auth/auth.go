package auth

import (
	"context"
	"os"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
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

	// Query database
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
		// Retrieve role
		var role string
		if err = pool.QueryRow(context.Background(), "SELECT role FROM users where username=$1 and status='active' and password=$2", username, password).Scan(&role); err != nil {
			utils.Logger.Error().Err(err).Msg("Error retrieving role")
			response := LoginResponse{
				Success:  false,
				ErrorMsg: "Unexpected authentication error.",
			}
			c.JSON(200, &response)
			return
		}

		// Generate JWT
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": username,
			"role": role,
			"iat":  time.Now().Unix(),
		})

		secretKey := []byte(os.Getenv("JWT_KEY"))

		// Check if secret key was read correctly
		if len(secretKey) == 0 {
			utils.Logger.Warn().Msg("JWT secret key is empty!")
		}

		tokenString, err := token.SignedString(secretKey)
		if err != nil {
			utils.Logger.Error().Err(err).Msg("Error signing JWT token")
			c.JSON(500, "Internal server error")
		}

		response := LoginResponse{
			Success: true,
			Jwt:     tokenString,
		}
		c.JSON(200, &response)
		return
	}
}
