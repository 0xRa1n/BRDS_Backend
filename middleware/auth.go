package middleware

import (
	"net/http"

	"backend/config"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

// AuthRequired is a middleware to protect routes requiring JWT authentication
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("jwt_token")
		if err != nil || tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization cookie required"})
			return
		}
		claims, err := utils.ExtractJWTClaims(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		jti, ok := claims["jti"].(string)
		if !ok || config.AppCache.IsJWTBlocked(jti) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token has been invalidated"})
			return
		}

		// Set claims in context for downstream handlers
		c.Set("phone", claims["phone"])
		c.Set("role", claims["role"])

		c.Next()
	}
}
