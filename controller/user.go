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
	Find(c *gin.Context)
	FindByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	DeleteByID(c *gin.Context)
}

func NewUserController() UserController {
	return &userController{
		svc: services.NewUserService(),
	}
}

func (ctlr *userController) Create(c *gin.Context) {
	logger.Info("API Request for creating a user.")
	userDTO := &dtos.UserDTO{}
	err := c.ShouldBindBodyWithJSON(userDTO)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Invalid Request Body", "result": gin.H{"error": err.Error()}})
		logger.Info("Create user api stopped due to request body is invalid")
		return
	}

	user, appErr := ctlr.svc.Create(userDTO)
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

func (ctrl *userController) UpdateByID(c *gin.Context) {
	idStr := c.Param("id")
	logger.Info("API Request for updating a user by ID " + idStr + ".")

	userDTO := &dtos.UserDTO{}
	err := c.ShouldBindBodyWithJSON(userDTO)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Invalid Request Body", "result": gin.H{"error": err.Error()}})
		logger.Info("Update user by id api stopped due to request body is invalid")
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Invalid User ID", "result": gin.H{"error": err.Error()}})
		logger.Info("update user by id api stopped")
		return
	}

	updateUser, appErr := ctrl.svc.UpdateByID(uint(id), userDTO)
	if appErr != nil {
		appErr.WriteHTTPResponse(c)
		logger.Info("Update user by id api stopped")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User Updated", "result": gin.H{"user": updateUser}})
	logger.Info("Update user by id api finished")
}

func (ctrl *userController) DeleteByID(c *gin.Context) {
	idStr := c.Param("id")
	logger.Info("API Request for deleting a user by ID " + idStr + ".")

	id, err := strconv.ParseUint(idStr, 10, 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Invalid User ID", "result": gin.H{"error": err.Error()}})
		logger.Info("Delete user by id api stopped")
		return
	}

	appErr := ctrl.svc.DeleteByID(uint(id))
	if appErr != nil {
		appErr.WriteHTTPResponse(c)
		logger.Info("Delete user by id api stopped")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User Deleted"})
	logger.Info("Delete user by id api finished")
}
