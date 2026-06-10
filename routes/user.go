package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes sets up the user profile routes
func RegisterUserRoutes(router *gin.Engine) {
	userGroup := router.Group("/api/v1/users")
	userGroup.Use(middleware.AuthRequired())
	{
		userGroup.GET("/profile", controllers.GetProfile)
		userGroup.PUT("/profile", controllers.UpdateProfile)
	}
}
