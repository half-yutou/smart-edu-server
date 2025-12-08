package danmaku

import "github.com/gin-gonic/gin"

type Handler interface {
	Send(c *gin.Context)
	List(c *gin.Context)
}
