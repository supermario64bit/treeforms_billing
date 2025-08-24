package models

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	gorm.Model
	UserID    uint      `gorm:"not null;index"`
	TokenHash string    `gorm:"not null;unique"`
	ExpiresAt time.Time `gorm:"not null"`
}
