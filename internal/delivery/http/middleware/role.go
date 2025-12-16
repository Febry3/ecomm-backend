package middleware

import (
	"net/http"

	"github.com/febry3/gamingin/internal/dto"
	"github.com/gin-gonic/gin"
)

func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  false,
				"message": "user not found in context",
			})
			return
		}

		jwt, ok := user.(*dto.JwtPayload)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "invalid user context",
			})
			return
		}

		if jwt.Role != requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":  false,
				"message": "access denied: " + requiredRole + " role required",
			})
			return
		}

		c.Next()
	}
}
