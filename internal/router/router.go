package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"xinde/internal/handler"
	"xinde/internal/handler/account"
	"xinde/internal/middleware/auth"
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

		// ========== 公开接口（不需要认证）==========
		accountGroup := apiV1.Group("/account")
		{
			accountGroup.POST("/register", accountCtrl.Register)
			accountGroup.POST("/login", accountCtrl.Login)
		}

		// ========== 管理员接口（需要管理员权限）==========
		adminGroup := apiV1.Group("/admin")
		adminGroup.Use(auth.JWTAuth(), auth.AdminAuth())
		{
			adminAccountGroup := adminGroup.Group("/account")
			{
				adminAccountGroup.GET("/list", accountCtrl.List)
			}
		}

		// ========== 需要认证的接口 ==========

		// ========== 可选认证接口（有token更好，没有也行）==========

	}

	return router, nil
}
