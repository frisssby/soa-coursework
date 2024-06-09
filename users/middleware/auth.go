package middleware

import (
	"net/http"

	"users/jwt"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware(c *gin.Context) {
	tokenString, err := c.Cookie("jwt")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "No token provided")
		return
	}
	username, err := jwt.ValidateJWT(tokenString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid token")
	}
	c.Set("username", username)
}
