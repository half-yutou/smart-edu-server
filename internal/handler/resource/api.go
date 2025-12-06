package resource

import "github.com/gin-gonic/gin"

type Handler interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
	ListMyResources(c *gin.Context)
}