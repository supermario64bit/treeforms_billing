package routes

import (
	"treeforms_billing/middlewares"

	"github.com/gin-gonic/gin"
)

func MountHTTPRoutes(r *gin.Engine) {
	authenticationMiddleware := middlewares.NewAuthenticationMiddleware()
	api := r.Group("/api/v1")
	apiProtected := r.Group("/api/v1", authenticationMiddleware.ValidateAccessToken)

	mountUserRoutes(apiProtected)
	mountAuthenticationRoutes(api)
}
