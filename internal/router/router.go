package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"xinde/internal/handler"
	"xinde/internal/handler/account"
	"xinde/internal/handler/attachment"
	"xinde/internal/handler/company"
	"xinde/internal/handler/device"
	"xinde/internal/handler/group"
	"xinde/internal/handler/price"
	"xinde/internal/middleware/auth"
)

func InitRouter() (*gin.Engine, error) {
	router := gin.Default()

	uploadUrlPrefix := viper.GetString("attachment.upload_url_prefix") // -> "/static/uploads"
	savePath := viper.GetString("attachment.save_path")                // -> "uploads"
	router.Static(uploadUrlPrefix, savePath)

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 创建controller实例
	accountCtrl, err := account.NewAccountController()
	if err != nil {
		return nil, fmt.Errorf("初始化AccountController失败: %w", err)
	}
	companyCtrl, err := company.NewCompanyController()
	if err != nil {
		return nil, fmt.Errorf("初始化CompanyController失败: %w", err)
	}
	priceCtrl, err := price.NewController()
	if err != nil {
		return nil, fmt.Errorf("初始化PriceController失败: %w", err)
	}
	attachmentCtrl, err := attachment.NewAttachmentController()
	if err != nil {
		return nil, fmt.Errorf("初始化AttachmentController失败: %w", err)
	}
	groupCtrl, err := group.NewGroupController()
	if err != nil {
		return nil, fmt.Errorf("初始化GroupController失败: %w", err)
	}
	deviceCtrl, err := device.NewDeviceController()
	if err != nil {
		return nil, fmt.Errorf("初始化DeviceController失败: %w", err)
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
			//TODO 用户访问记录
			adminAccountGroup := adminGroup.Group("/account")
			{
				adminAccountGroup.GET("/list", accountCtrl.List) //TODO 接入用户访问记录
				adminAccountGroup.GET("/approval/list", accountCtrl.ApprovalList)
				adminAccountGroup.POST("/approval/:id", accountCtrl.Approve)
				adminAccountGroup.DELETE("/:id", accountCtrl.DeleteUser)
				adminAccountGroup.POST("/reset/password/:id", accountCtrl.ResetPassword)
				adminAccountGroup.PATCH("/remark/:id", accountCtrl.ResetRemark)
				adminAccountGroup.PATCH("/password/:id", accountCtrl.UpdatePassword)
			}

			adminCompanyGroup := adminGroup.Group("/company")
			{
				adminCompanyGroup.GET("/list", companyCtrl.List)
				adminCompanyGroup.PATCH("/price/level/:id", companyCtrl.UpdatePriceLevel)
			}

			adminPriceGroup := adminGroup.Group("/price")
			{
				adminPriceGroup.GET("/list", priceCtrl.List)
				adminPriceGroup.POST("/import", priceCtrl.Import)
			}

			attachmentGroup := adminGroup.Group("/attachment")
			{
				attachmentGroup.GET("/list", attachmentCtrl.List)
				attachmentGroup.GET("/download/:id", attachmentCtrl.Download)
				attachmentGroup.DELETE("/:id", attachmentCtrl.Delete)
				attachmentGroup.GET("/scan/invalid", attachmentCtrl.ScanInvalid)
				attachmentGroup.POST("/fix/orphan", attachmentCtrl.FixOrphan)
			}

			groupGroup := adminGroup.Group("/group")
			{
				groupGroup.POST("/create", groupCtrl.Create)
				groupGroup.GET("/tree", groupCtrl.GetTree)
				groupGroup.GET("/list", groupCtrl.List)
				groupGroup.PUT("/update/:id", groupCtrl.Update)
				groupGroup.DELETE("/delete/:id", groupCtrl.Delete)
				groupGroup.GET("/device/list/:id", deviceCtrl.GroupDeviceList)
			}

			deviceGroup := adminGroup.Group("/device")
			{
				deviceGroup.POST("/import", deviceCtrl.Import)
				deviceGroup.GET("/list", deviceCtrl.List)
				deviceGroup.PUT("/import/:id", deviceCtrl.UpdateImport)
				deviceGroup.PATCH("/update/group/:id", deviceCtrl.UpdateGroup)
				deviceGroup.PATCH("/update/name/:id", deviceCtrl.UpdateName)
				deviceGroup.POST("/update/image/:id", deviceCtrl.UpdateImage)
				deviceGroup.DELETE("/delete/:id", deviceCtrl.Delete)
			}

			filterImageGroup := adminGroup.Group("/filter_image")
			{
				filterImageGroup.POST("/create", deviceCtrl.CreateFilterImage)
				filterImageGroup.GET("/list", deviceCtrl.FilterImageList)
				filterImageGroup.DELETE("/delete/:id", deviceCtrl.DeleteFilterImage)
				filterImageGroup.PATCH("/change/device_type/:id", deviceCtrl.ChangeFilterImageDevice)
			}
		}

		// ========== 需要认证的接口 ==========

		// ========== 可选认证接口（有token更好，没有也行）==========

	}

	return router, nil
}
