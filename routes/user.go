package routes

import (
	"treeforms_billing/controller"

	"github.com/gin-gonic/gin"
)

func mountUserRoutes(r *gin.RouterGroup) {
	userRoutes := r.Group("/user")
	userController := controller.NewUserController()

	userRoutes.POST("", userController.Create)
	userRoutes.GET("", userController.Find)
	userRoutes.GET("/:id", userController.FindByID)
	userRoutes.PATCH("/:id", userController.UpdateByID)
	userRoutes.DELETE("/:id", userController.DeleteByID)
}
