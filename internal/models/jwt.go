package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type JwtClaim struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Name     string `json:"name"`
	Semester int    `json:"semester"`

	jwt.RegisteredClaims
}

func NewClaim(user *UserProfile, role string) *JwtClaim {
	return &JwtClaim{
		Username: user.Username,
		Role:     role,
		Name:     user.Name,
		Semester: user.Semester,
	}
}

// // Represent user claim before signing the JWT token.
// type UserClaim struct {
// 	Username string `json:"username"`
// 	Role     string `json:"role"`
// }
