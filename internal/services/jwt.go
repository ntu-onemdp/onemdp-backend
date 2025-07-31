package services

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ntu-onemdp/onemdp-backend/internal/models"
	"github.com/ntu-onemdp/onemdp-backend/internal/utils"
)

// JWT object storing secret key
type Jwt struct {
	secretKey []byte
}

// Store global JWT handler instance
var JwtHandler *Jwt

// Load secret key from file. Panics if secret key cannot be found.
func InitJwt() {
	// Get app env
	env, found := os.LookupEnv("ENV")
	if !found {
		// Default environment: PROD
		env = "PROD"
	}

	var path string
	switch env {
	case "PROD", "QA":
		path = "mnt/secrets/jwt-key"
	case "DEV":
		path = "config/jwt-key.txt"
	default:
		path = "run/secrets/jwt-key" // only used when running from docker compose, can be removed.
	}
	key, err := os.ReadFile(path)
	if err != nil {
		utils.Logger.Warn().Err(err).Msgf("Error reading secret from %s", path)
	}

	// Check if secret key was read correctly
	if len(key) == 0 {
		utils.Logger.Warn().Msg("JWT secret key is empty!")
	}

	JwtHandler = &Jwt{
		secretKey: key,
	}
}

// Generate and sign jwt, returning a token as a string
func (j *Jwt) GenerateJwt(claim *models.JwtClaim) (string, error) {
	secretKey := j.secretKey

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": claim.Uid,
		"iat": time.Now().Unix(),
	})

	// Sign key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error signing JWT token")
		return "", err
	}

	return tokenString, nil
}

// Get uid from JWT token in request.
// This acts as a middleware, so it automatically returns a 401 Unauthorized response if the JWT is invalid or missing.
func (j *Jwt) GetUidFromJwt(c *gin.Context) string {
	// Get Uid from JWT token
	jwt := c.Request.Header.Get("Authorization")
	claim, err := j.ParseJwt(jwt)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("Error parsing JWT token")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Error parsing JWT token",
			"error":   err.Error(),
		})
		return ""
	}

	return claim.Uid
}

// Parse signed jwt string
func (j *Jwt) ParseJwt(tokenString string) (*models.JwtClaim, error) {
	secretKey := j.secretKey

	// Remove "Bearer " prefix if included
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	// Parse and verify the token
	token, err := jwt.ParseWithClaims(tokenString, &models.JwtClaim{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Extract the claims
	if claim, ok := token.Claims.(*models.JwtClaim); ok && token.Valid {
		return claim, nil
	}

	return nil, fmt.Errorf("error extracting claim from jwt")
}
