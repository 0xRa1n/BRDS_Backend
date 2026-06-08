package routes

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRequestRoutes sets up the document request routes
func RegisterRequestRoutes(router *gin.Engine) {
	// Public tracking route
	router.GET("/api/v1/track/:reference", controllers.TrackRequest)

	reqGroup := router.Group("/api/v1/request")
	reqGroup.Use(middleware.AuthRequired())
	{
		reqGroup.POST("", controllers.SubmitRequest)
		reqGroup.GET("", controllers.GetRequests)
	}
}
