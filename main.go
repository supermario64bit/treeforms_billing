package main

import (
	"net/http"
	"os"
	"treeforms_billing/db"
	"treeforms_billing/logger"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		logger.HighlightedDanger("Error while loading .env file. Message: " + err.Error())
		return
	}

	// Automigrate DB
	db.Automigrate()

	r := gin.Default()

	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
		logger.Info("Running in release mode")
	} else {
		logger.Info("Running in development mode")
	}

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello World!"})
	})

	r.Run()
}
