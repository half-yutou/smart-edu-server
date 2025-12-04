package user

import "github.com/gin-gonic/gin"

type Handler interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	UpdateProfile(c *gin.Context)
}
