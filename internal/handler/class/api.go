package class

import "github.com/gin-gonic/gin"

type Handler interface {
	Create(c *gin.Context)
	ListForTeacher(c *gin.Context)
	DeleteForTeacherByID(c *gin.Context)
	UpdateForTeacherByID(c *gin.Context)

	GetByCode(c *gin.Context)
	GetByID(c *gin.Context)

	ListForStudent(c *gin.Context)
	Quit(c *gin.Context)
	JoinByCode(c *gin.Context)
}
