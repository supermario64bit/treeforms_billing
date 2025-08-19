package db

import (
	"os"
	"treeforms_billing/logger"
	"treeforms_billing/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Get() *gorm.DB {
	if db != nil {
		return db
	}
	var err error

	dsn := os.Getenv("DB_DSN")
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Danger("Database connection failed. Message: " + err.Error())
		return nil
	}
	return db
}

func Automigrate() {
	if db == nil {
		Get()
	}
	db.AutoMigrate(
		models.User{},
	)

	passwordsTableCreateQuery := `
	CREATE TABLE IF NOT EXISTS passwords (
	    id BIGSERIAL PRIMARY KEY,
	    hash TEXT NOT NULL,
	    user_id BIGINT UNIQUE NOT NULL
	);`

	if err := db.Exec(passwordsTableCreateQuery).Error; err != nil {
		logger.HighlightedDanger("failed to run migration:" + err.Error())
	}
}
