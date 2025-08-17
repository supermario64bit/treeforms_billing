package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name   string `json:"name" validate:"required" gorm:"not null"`
	Email  string `json:"email" validate:"required,email" gorm:"unique;not null"`
	Phone  string `json:"phone" validate:"required" gorm:"unique;not null" `
	Role   string `json:"role" validate:"required,oneof=superadmin admin" gorm:"not null"`
	Status string `json:"role" validate:"required,oneof=active inactive" gorm:"not null"`
}

func (u *User) ValidateFields() error {
	err := validate.Struct(u)
	return err
}
