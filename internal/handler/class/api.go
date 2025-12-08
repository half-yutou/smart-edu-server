package class

import "github.com/gin-gonic/gin"

type Handler interface {
	Create(c *gin.Context)
	ListForTeacher(c *gin.Context)
	DeleteForTeacherByID(c *gin.Context)
	UpdateForTeacherByID(c *gin.Context)
	ListMembers(c *gin.Context)

	GetByCode(c *gin.Context)
	GetByID(c *gin.Context)

	ListForStudent(c *gin.Context)
	Quit(c *gin.Context)
	JoinByCode(c *gin.Context)

	// 班级资源关联

	AddResource(c *gin.Context)
	RemoveResource(c *gin.Context)
	ListResources(c *gin.Context)
}
