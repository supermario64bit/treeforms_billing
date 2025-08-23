package routes

import "github.com/gin-gonic/gin"

func MountHTTPRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")

	mountUserRoutes(api)
	mountAuthenticationRoutes(api)
}
