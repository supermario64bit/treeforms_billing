package middlewares

import (
	"net/http"
	"os"
	"strings"
	"treeforms_billing/application_types"
	"treeforms_billing/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type authenticationMiddleware struct {
	authenticationSvc services.AuthenticationService
	userSvc           services.UserService
}

type AuthenticationMiddleware interface {
	ValidateAccessToken(c *gin.Context)
}

func NewAuthenticationMiddleware() AuthenticationMiddleware {
	return &authenticationMiddleware{
		authenticationSvc: services.NewAuthenticationSevice(),
		userSvc:           services.NewUserService(),
	}
}

func (mw *authenticationMiddleware) ValidateAccessToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "failed", "message": "No access token found in the header."})
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "failed", "message": "Invalid authorization header format."})
		return
	}

	claims := application_types.AccessToken{}
	token, err := jwt.ParseWithClaims(parts[1], &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SIGNING_SECRET")), nil
	})

	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "failed", "message": "Invalid Token"})
		return
	}

	c.Set("userID", claims.ID)
	c.Set("userRole", claims.Role)
	c.Set("userName", claims.Name)

	c.Next()
}
