package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterAdminRoutes sets up endpoints for admin/staff functionality
func RegisterAdminRoutes(router *gin.Engine) {
	adminGroup := router.Group("/api/admin")
	
	// Public admin route for login
	adminGroup.POST("/login", controllers.AdminLogin)

	// Protected admin routes
	protectedAdminGroup := adminGroup.Group("/")
	protectedAdminGroup.Use(middleware.AdminAuthRequired())
	{
		protectedAdminGroup.GET("/requests", controllers.AdminGetAllRequests)
		protectedAdminGroup.POST("/requests", controllers.AdminCreateRequest)
	}
}
