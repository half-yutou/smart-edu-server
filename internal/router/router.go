package router

import (
	classHandler "smarteduhub/internal/handler/class"
	homeworkHandler "smarteduhub/internal/handler/homework"
	resourceHandler "smarteduhub/internal/handler/resource"
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
	r.Use(middleware.Cors()) // 跨域中间件

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
		classH := classHandler.NewHandler()
		resourceH := resourceHandler.NewHandler()
		homeworkH := homeworkHandler.NewHandler()

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
			// userGroup.POST("/upload/image", uploadH.UploadImage)
		}

		classGroup := apiGroup.Group("/class")
		classGroup.Use(middleware.Auth()) // 基础认证
		{
			classGroup.POST("/get", classH.GetByID)
			classGroup.POST("/code", classH.GetByCode)
			classGroup.GET("/resource/list", classH.ListResources) // 班级资源列表 (老师学生通用)
		}
		// 班级管理路由 (需要 登录+教师)
		classGroupForTeacher := apiGroup.Group("/class/teacher")
		classGroupForTeacher.Use(middleware.Auth(), middleware.AuthTeacher())
		{
			classGroupForTeacher.POST("/create", classH.Create)
			classGroupForTeacher.POST("/list", classH.ListForTeacher)
			classGroupForTeacher.POST("/delete", classH.DeleteForTeacherByID)
			classGroupForTeacher.POST("/update", classH.UpdateForTeacherByID)
			classGroupForTeacher.GET("/members", classH.ListMembers) // 班级成员列表

			// 班级资源管理
			classGroupForTeacher.POST("/resource/add", classH.AddResource)
			classGroupForTeacher.POST("/resource/remove", classH.RemoveResource)
		}
		// 班级管理路由 (需要 登录+学生)
		classGroupForStudent := apiGroup.Group("/class/student")
		classGroupForStudent.Use(middleware.Auth(), middleware.AuthStudent())
		{
			classGroupForStudent.POST("/list", classH.ListForStudent)
			classGroupForStudent.POST("/join", classH.JoinByCode)
			classGroupForStudent.POST("/quit", classH.Quit)
		}

		// 资源相关路由
		// 1. 公共部分 (无需登录，广场资源是公开的)
		resourcePublicGroup := apiGroup.Group("/resource")
		{
			resourcePublicGroup.GET("/list", resourceH.List)
			resourcePublicGroup.GET("/detail", resourceH.GetByID) // 使用 Query: ?id=123
		}

		// 2. 教师管理部分 (需要 登录+教师)
		resourceTeacherGroup := apiGroup.Group("/resource/teacher")
		resourceTeacherGroup.Use(middleware.Auth(), middleware.AuthTeacher())
		{
			resourceTeacherGroup.POST("/upload", uploadH.UploadFile) // 移动到这里
			resourceTeacherGroup.POST("/create", resourceH.Create)
			resourceTeacherGroup.POST("/update", resourceH.Update)
			resourceTeacherGroup.POST("/delete", resourceH.Delete)
			resourceTeacherGroup.GET("/my", resourceH.ListMyResources)
		}

		// 作业相关路由 (教师)
		homeworkTeacherGroup := apiGroup.Group("/homework/teacher")
		homeworkTeacherGroup.Use(middleware.Auth(), middleware.AuthTeacher())
		{
			homeworkTeacherGroup.POST("/create", homeworkH.Create)
			homeworkTeacherGroup.POST("/delete", homeworkH.Delete)
			homeworkTeacherGroup.POST("/update", homeworkH.Update)
			homeworkTeacherGroup.GET("/list", homeworkH.ListByTeacher)
			homeworkTeacherGroup.GET("/detail", homeworkH.GetByID)
			homeworkTeacherGroup.GET("/submissions", homeworkH.ListSubmissions)
			homeworkTeacherGroup.POST("/grade", homeworkH.GradeSubmission)
		}

		// 作业相关路由 (学生)
		homeworkStudentGroup := apiGroup.Group("/homework/student")
		homeworkStudentGroup.Use(middleware.Auth(), middleware.AuthStudent())
		{
			homeworkStudentGroup.GET("/list", homeworkH.ListByClass)
			homeworkStudentGroup.POST("/submit", homeworkH.Submit)
			homeworkStudentGroup.GET("/submission", homeworkH.GetSubmission)
		}
	}

	return r
}
