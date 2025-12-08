package homework

import "github.com/gin-gonic/gin"

type Handler interface {
	Create(c *gin.Context)
	Delete(c *gin.Context)
	Update(c *gin.Context)
	GetByID(c *gin.Context)
	ListByTeacher(c *gin.Context)
	ListByClass(c *gin.Context)
	Submit(c *gin.Context)
	GetSubmission(c *gin.Context)

	ListSubmissions(c *gin.Context)
	GradeSubmission(c *gin.Context)
}
