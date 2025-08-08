package router

import (
	"xinde/internal/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Ping test is now inside the /api/v1/health group

	// API v1 routes
	apiV1 := router.Group("/api/v1")
	{
		// Health check group
		health := apiV1.Group("/health")
		{
			health.GET("/ping", handler.Ping)
		}

		enroll := apiV1.Group("/enroll")
		{
			enroll.POST("/", handler.Enroll)
		}
		// You can add more groups here, for example:
		// userGroup := apiV1.Group("/users")
		// {
		// 	userGroup.GET("/", handler.GetUsers) // Example
		// }
	}

	return router
}
