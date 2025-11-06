package middleware

import (
	"net/http"
	"strings"

	"github.com/febry3/gamingin/internal/helpers"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwt *helpers.JwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "missing or invalid authorization header"})
			return
		}

		tokenString := strings.Split(authHeader, "Bearer ")[1]

		user, err := jwt.VerifyToken(tokenString)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid or expired token"})
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
