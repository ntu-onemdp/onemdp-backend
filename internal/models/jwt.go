package models

import (
	"github.com/golang-jwt/jwt/v5"
)

// Jwt will store only Uid. Retrieve role from database instead of jwt.
type JwtClaim struct {
	Uid string `json:"uid"`

	jwt.RegisteredClaims
}

func NewClaim(uid string) *JwtClaim {
	return &JwtClaim{
		Uid: uid,
	}
}
