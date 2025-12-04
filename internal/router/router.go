package router

import (
	uploadHandler "smarteduhub/internal/handler/upload"
	userHandler "smarteduhub/internal/handler/user"
	"smarteduhub/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// 全局中间件
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// 静态资源服务 (用于访问上传的图片)
	r.Static("/uploads", "./uploads")

	// 健康检查
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// API 路由组
	apiGroup := r.Group("/api/v1")
	{
		userH := userHandler.NewHandler()
		uploadH := uploadHandler.NewHandler()

		// 公开路由 (无需登录)
		authGroup := apiGroup.Group("/user")
		{
			authGroup.POST("/register", userH.Register)
			authGroup.POST("/login", userH.Login)
		}

		// 用户相关路由 (需要登录)
		userGroup := apiGroup.Group("/user")
		userGroup.Use(middleware.Auth())
		{
			userGroup.GET("/hello", func(c *gin.Context) {
				uid, _ := c.Get("uid")
				c.JSON(200, gin.H{
					"message": "Hello, SmartEduHub!",
					"uid":     uid,
				})
			})

			// 用户相关
			userGroup.POST("/user/logout", userH.Logout)
			userGroup.POST("/user/profile", userH.UpdateProfile)

			// 上传头像相关
			userGroup.POST("/upload/image", uploadH.UploadImage)
		}

		// 班级相关路由 (需要登录)
		classGroup := apiGroup.Group("/class")
		classGroup.Use(middleware.Auth())
		{
			// 班级相关操作
		}
	}

	return r
}
