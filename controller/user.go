package controller

import (
	"net/http"
	"strconv"
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

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User Created", "result": gin.H{"user": user}})
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

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Users found", "result": gin.H{"users": users}})
	logger.Info("Find users api finished")
}

func (ctrl *userController) FindByID(c *gin.Context) {
	idStr := c.Param("id")
	logger.Info("API Request for finding user by id " + idStr + ".")

	id, err := strconv.ParseUint(idStr, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Invalid User ID", "result": gin.H{"error": err.Error()}})
		logger.Info("Find userby id api stopped")
		return
	}

	user, appErr := ctrl.svc.FindByID(uint(id))
	if appErr != nil {
		appErr.WriteHTTPResponse(c)
		logger.Info("Find user by id api stopped")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User Found", "result": gin.H{"user": user}})
	logger.Info("Find user by id api finished")
}
