package router

import (
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// 全局中间件
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// 健康检查
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// API 路由组
	apiGroup := r.Group("/api/v1")
	{
		apiGroup.GET("/hello", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Hello, SmartEduHub!",
			})
		})
		
		// userGroup := apiGroup.Group("/users")
		// ...
	}

	return r
}