package models

import (
	"os"
	"time"
	"treeforms_billing/application_types"
	"treeforms_billing/logger"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name   string `json:"name" validate:"required" gorm:"not null"`
	Email  string `json:"email" validate:"required,email" gorm:"not null"`
	Phone  string `json:"phone" validate:"required" gorm:"not null" `
	Role   string `json:"role" validate:"required,oneof=superadmin admin" gorm:"not null"`
	Status string `json:"role" validate:"required,oneof=active inactive" gorm:"not null"`
}

func (u *User) ValidateFields() error {
	err := validate.Struct(u)
	return err
}

func (u *User) NewAccessToken() (string, error) {
	token := application_types.AccessToken{
		UserID: u.ID,
		Name:   u.Name,
		Role:   u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)), // 5 minutes expiry
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "Treeforms Billing Software",
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, token)
	secret := []byte(os.Getenv("JWT_SIGNING_SECRET"))
	tokenString, err := jwtToken.SignedString(secret)
	if err != nil {
		logger.HighlightedDanger("Unable to sign the token. Message: " + err.Error())
		return "", err
	}

	return tokenString, nil
}
