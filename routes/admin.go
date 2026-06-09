package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterAdminRoutes sets up endpoints for admin/staff functionality
func RegisterAdminRoutes(router *gin.Engine) {
	adminGroup := router.Group("/api/v1/admin")
	
	// Public admin route for login
	adminGroup.POST("/login", controllers.AdminLogin)

	// Protected admin routes
	protectedAdminGroup := adminGroup.Group("/")
	protectedAdminGroup.Use(middleware.AdminAuthRequired())
	{
		protectedAdminGroup.GET("/requests", controllers.AdminGetAllRequests)
		protectedAdminGroup.POST("/requests", controllers.AdminCreateRequest)
		protectedAdminGroup.GET("/requests/:id", controllers.AdminGetRequest)
		protectedAdminGroup.PUT("/requests/:id/status", controllers.AdminUpdateStatus)
		protectedAdminGroup.PUT("/requests/:id/appointment", controllers.AdminSetAppointment)
	}

	// Super Admin routes
	superAdminRoot := adminGroup.Group("/")
	superAdminRoot.Use(middleware.SuperAdminAuthRequired())
	{
		usersGroup := superAdminRoot.Group("/users")
		{
			usersGroup.GET("", controllers.AdminGetUsers)
			usersGroup.POST("", controllers.AdminCreateUser)
			usersGroup.PUT("/:id", controllers.AdminUpdateUser)
			usersGroup.DELETE("/:id", controllers.AdminDeleteUser)
		}

		// Portal Users
		superAdminRoot.GET("/portal-users", controllers.AdminGetPortalUsers)
		superAdminRoot.PUT("/portal-users/:id", controllers.AdminUpdatePortalUser)
		superAdminRoot.DELETE("/portal-users/:id", controllers.AdminDeletePortalUser)

		// Archives
		superAdminRoot.GET("/archives/users", controllers.AdminGetArchivedAdmins)
		superAdminRoot.POST("/archives/users/:id/recover", controllers.AdminRecoverAdmin)
		superAdminRoot.GET("/archives/portal-users", controllers.AdminGetArchivedPortalUsers)
		superAdminRoot.POST("/archives/portal-users/:id/recover", controllers.AdminRecoverPortalUser)
	}
}
