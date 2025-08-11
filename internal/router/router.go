package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"xinde/internal/handler"
	"xinde/internal/handler/account"
)

func InitRouter() (*gin.Engine, error) {
	router := gin.Default()

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 创建controller实例
	accountCtrl, err := account.NewAccountController()
	if err != nil {
		return nil, fmt.Errorf("初始化AccountController失败: %w", err)
	}
	
	// API v1 routes
	apiV1 := router.Group("/api/v1")
	{
		// Health check group
		health := apiV1.Group("/health")
		{
			health.GET("/ping", handler.Ping)
		}

		accountGroup := apiV1.Group("/account")
		{
			accountGroup.POST("/register", accountCtrl.Register)
		}
		// You can add more groups here, for example:
		// userGroup := apiV1.Group("/users")
		// {
		// 	userGroup.GET("/", handler.GetUsers) // Example
		// }
	}

	return router, nil
}
