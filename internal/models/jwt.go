package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type JwtClaim struct {
	Uid  string `json:"uid"`
	Role string `json:"role"`
	Name string `json:"name"`

	jwt.RegisteredClaims
}

func NewClaim(user *UserProfile, uid string) *JwtClaim {
	return &JwtClaim{
		Uid:  uid,
		Role: user.Role,
		Name: user.Name,
	}
}
