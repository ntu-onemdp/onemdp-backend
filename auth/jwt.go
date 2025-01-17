package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	utils "github.com/ntu-onemdp/onemdp-backend/utils"
)

type JwtClaim struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

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
		utils.Logger.Error().Err(err).Msg("Error signing JWT token")
		return "", err
	}

	return tokenString, nil
}

// Parse signed jwt string
func ParseJwt(tokenString string) (*JwtClaim, error) {
	secretKey := getSecretKey()

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
		utils.Logger.Warn().Msg("JWT secret key is empty!")
	}

	return secretKey
}
