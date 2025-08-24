package routes

import (
	"treeforms_billing/controller"

	"github.com/gin-gonic/gin"
)

func mountAuthenticationRoutes(r *gin.RouterGroup) {
	authenticationRoutes := r.Group("/authentication")

	ctrl := controller.NewAuthenticationController()
	authenticationRoutes.POST("/signup", ctrl.Signup)
	authenticationRoutes.POST("/login/email", ctrl.EmailLogin)
	authenticationRoutes.POST("/refresh-token", ctrl.RotateRefreshTokenWithNewAccessToken)
}
