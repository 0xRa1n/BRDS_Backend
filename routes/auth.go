package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes sets up the authentication routes
func RegisterAuthRoutes(router *gin.Engine) {
	authGroup := router.Group("/api/v1/auth")
	{
		authGroup.POST("/send-otp", controllers.SendOTP)
		authGroup.POST("/verify-otp", controllers.VerifyOTP)
		authGroup.POST("/logout", controllers.Logout)
	}
}
