package middleware

import (
	"net/http"
	"strings"

	"backend/config"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

// AdminAuthRequired is a middleware to protect routes requiring staff/admin privileges
func AdminAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
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

		role, ok := claims["role"].(string)
		roleLower := strings.ToLower(role)
		if !ok || (roleLower != "admin" && roleLower != "staff") {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient privileges"})
			return
		}

		// Set claims in context for downstream handlers
		c.Set("username", claims["username"])
		c.Set("role", role)

		c.Next()
	}
}

// SuperAdminAuthRequired protects routes that require strict Admin privileges
func SuperAdminAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
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

		role, ok := claims["role"].(string)
		if !ok || strings.ToLower(role) != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Super Admin privileges required"})
			return
		}

		// Set claims in context for downstream handlers
		c.Set("username", claims["username"])
		c.Set("role", role)

		c.Next()
	}
}
