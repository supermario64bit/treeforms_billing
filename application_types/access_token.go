package application_types

import "github.com/golang-jwt/jwt/v4"

type AccessToken struct {
	UserID uint   `json:"sub"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
