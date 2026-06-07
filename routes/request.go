package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRequestRoutes sets up the document request routes
func RegisterRequestRoutes(router *gin.Engine) {
	reqGroup := router.Group("/api/v1/request")
	reqGroup.Use(middleware.AuthRequired())
	{
		reqGroup.POST("", controllers.SubmitRequest)
	}
}
