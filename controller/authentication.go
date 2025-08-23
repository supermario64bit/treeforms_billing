package controller

import (
	"net/http"
	"treeforms_billing/dtos"
	"treeforms_billing/logger"
	"treeforms_billing/services"

	"github.com/gin-gonic/gin"
)

type authenticatioController struct {
	authSvc services.AuthenticationService
}

type AuthenticatioController interface {
	Signup(c *gin.Context)
	EmailLogin(c *gin.Context)
}

func NewAuthenticationController() AuthenticatioController {
	return &authenticatioController{
		authSvc: services.NewAuthenticationSevice(),
	}
}

func (ctrl *authenticatioController) Signup(c *gin.Context) {
	logger.Info("API Request for signup")
	var signDto dtos.SignupDTO
	if err := c.ShouldBindBodyWithJSON(&signDto); err != nil {
		logger.Danger("Invalid Payload. Message: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Unable to read request body!", "result": gin.H{"error": err.Error()}})
		return
	}

	if appErr := ctrl.authSvc.Signup(signDto); appErr != nil {
		logger.Danger("API Request for signup Stopped")
		appErr.WriteHTTPResponse(c)
		return
	}

	logger.Success("API Request for signup success.")
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Signup is successfull."})
}

func (ctrl *authenticatioController) EmailLogin(c *gin.Context) {
	logger.Info("API Request for Email Login")
	var loginDto dtos.LoginDTO
	if err := c.ShouldBindBodyWithJSON(&loginDto); err != nil {
		logger.Danger("Invalid Payload. Message: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Unable to read request body!", "result": gin.H{"error": err.Error()}})
		return
	}

	if appErr := ctrl.authSvc.EmailLogin(loginDto.Email, loginDto.Password); appErr != nil {
		logger.Danger("API Request for Email Login Stopped")
		appErr.WriteHTTPResponse(c)
		return
	}

	logger.Success("API Request for Email Login success.")
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Login is successfull."})
}
