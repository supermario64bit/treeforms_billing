package controller

import (
	"net/http"
	"treeforms_billing/dtos"
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
	var userDTO *dtos.UserDTO
	err := c.ShouldBindBodyWithJSON(userDTO)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Invalid Request Body", "result": gin.H{"error": err.Error()}})
		return
	}

	user, appErr := ctlr.svc.Create(&models.User{Name: userDTO.Name, Email: userDTO.Email, Phone: userDTO.Phone,
		Role: userDTO.Role, Status: userDTO.Status})
	if appErr != nil {
		appErr.WriteHTTPResponse(c)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"status": "success", "message": "User Created", "result": gin.H{"user": user}})
}
