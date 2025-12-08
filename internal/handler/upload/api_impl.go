package upload

import (
	"smarteduhub/internal/model/dto/response"
	"smarteduhub/internal/pkg/oss"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	UploadFile(c *gin.Context)
}

type handlerImpl struct{}

func NewHandler() Handler {
	return &handlerImpl{}
}

func (h *handlerImpl) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, "Get file error: "+err.Error())
		return
	}

	// 上传到 OSS
	url, err := oss.UploadFile(file)
	if err != nil {
		response.Fail(c, "Upload to OSS failed: "+err.Error())
		return
	}

	response.SuccessWithMsg(c, "upload successful", gin.H{
		"url": url,
	})
}
