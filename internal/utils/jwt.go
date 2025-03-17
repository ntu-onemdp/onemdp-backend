package utils

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JwtClaim struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Represent user claim before signing the JWT token.
type UserClaim struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

// Generate and sign jwt, returning a token as a string
func GenerateJwt(claim UserClaim) (string, error) {
	secretKey := getSecretKey()

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": claim.Username,
		"role":     claim.Role,
		"iat":      time.Now().Unix(),
	})

	// Sign key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		Logger.Error().Err(err).Msg("Error signing JWT token")
		return "", err
	}

	return tokenString, nil
}

// Get username from JWT token in request
func GetUsernameFromJwt(c *gin.Context) string {
	// Get username from JWT token
	jwt := c.Request.Header.Get("Authorization")
	claim, err := ParseJwt(jwt)
	if err != nil {
		Logger.Error().Err(err).Msg("Error parsing JWT token")
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Error parsing JWT token",
			"error":   err.Error(),
		})
		return ""
	}

	return claim.Username
}

// Validate if username matches jwt token. Automatically returns false if error.
func ValidateUsername(username string, tokenString string) bool {
	// Remove "Bearer " prefix if included
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	claim, err := ParseJwt(tokenString)
	if err != nil {
		Logger.Error().Err(err)
		return false
	}

	return claim.Username == username
}

// Parse signed jwt string
func ParseJwt(tokenString string) (*JwtClaim, error) {
	secretKey := getSecretKey()

	// Remove "Bearer " prefix if included
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	// Parse and verify the token
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaim{}, func(token *jwt.Token) (interface{}, error) {
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
	if claim, ok := token.Claims.(*JwtClaim); ok && token.Valid {
		return claim, nil
	}

	return nil, fmt.Errorf("error extracting claim from jwt")
}

// Helper function to retrieve secret key
func getSecretKey() []byte {
	// Load secret key
	secretKey := []byte(os.Getenv("JWT_KEY"))

	// Check if secret key was read correctly
	if len(secretKey) == 0 {
		Logger.Warn().Msg("JWT secret key is empty!")
	}

	return secretKey
}
