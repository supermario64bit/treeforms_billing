package controller

import (
	"net/http"
	"treeforms_billing/dtos"
	"treeforms_billing/logger"
	"treeforms_billing/models"
	"treeforms_billing/services"

	"github.com/gin-gonic/gin"
)

type userController struct {
	svc services.UserService
}

type UserController interface {
	Create(c *gin.Context)
}

func NewUserController() UserController {
	return &userController{
		svc: services.NewUserService(),
	}
}

func (ctlr *userController) Create(c *gin.Context) {
	logger.Info("API Request for creating a user.")
	var userDTO *dtos.UserDTO
	err := c.ShouldBindBodyWithJSON(userDTO)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Invalid Request Body", "result": gin.H{"error": err.Error()}})
		logger.Info("Create user api stopped due to request body is invalid")
		return
	}

	user, appErr := ctlr.svc.Create(&models.User{Name: userDTO.Name, Email: userDTO.Email, Phone: userDTO.Phone,
		Role: userDTO.Role, Status: userDTO.Status})
	if appErr != nil {
		appErr.WriteHTTPResponse(c)
		logger.Info("Create user api stopped")
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"status": "success", "message": "User Created", "result": gin.H{"user": user}})
	logger.Info("Create user api finished")

}

func (ctrl *userController) Find(c *gin.Context) {
	logger.Info("API Request for finding users.")
	filter := &models.UserFilter{}
	err := c.ShouldBindBodyWithJSON(filter)
	if err != nil {
		logger.Info("Find user api stopped due to request body is invalid")
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Invalid Request Body", "result": gin.H{"error": err.Error()}})
		return
	}

	users, appErr := ctrl.svc.Find(*filter)
	if appErr != nil {
		appErr.WriteHTTPResponse(c)
		logger.Info("Find users api stopped")
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"status": "success", "message": "User Created", "result": gin.H{"users": users}})
	logger.Info("Find users api finished")
}
