package upload

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"smarteduhub/internal/model/dto/response"
	"smarteduhub/internal/pkg/oss"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	UploadImage(c *gin.Context)
	UploadFile(c *gin.Context)
}

type handlerImpl struct{}

func NewHandler() Handler {
	return &handlerImpl{}
}

func (h *handlerImpl) UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, "Get file error: "+err.Error())
		return
	}

	// 检查文件大小 (例如限制 5MB)
	if file.Size > 5*1024*1024 {
		response.Fail(c, "File size exceeds 5MB limit")
		return
	}

	// 检查文件类型 (简单的扩展名检查)
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		response.Fail(c, "Only jpg/png images are allowed")
		return
	}

	// 创建保存目录
	saveDir := "uploads/images"
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		response.Fail(c, "Create directory error: "+err.Error())
		return
	}

	// 生成唯一文件名
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	savePath := filepath.Join(saveDir, filename)

	// 保存文件
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		response.Fail(c, "Save file error: "+err.Error())
		return
	}

	// 返回可访问的 URL (假设服务器地址是配置好的，这里暂时返回相对路径或绝对路径)
	// 在生产环境中，这里应该是 OSS 的 URL
	// 本地开发时，我们需要在路由中配置静态资源服务
	fileURL := fmt.Sprintf("/uploads/images/%s", filename)

	response.SuccessWithMsg(c, "upload successful", gin.H{
		"url": fileURL,
	})
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
