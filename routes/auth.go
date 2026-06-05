package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes sets up the authentication routes
func RegisterAuthRoutes(router *gin.Engine) {
	authGroup := router.Group("/api/v1/auth")
	{
		authGroup.POST("/otp/send", controllers.SendOTP)
		authGroup.POST("/otp/verify", controllers.VerifyOTP)
		authGroup.POST("/logout", controllers.Logout)
	}
}
