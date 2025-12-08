package danmaku

import (
	"smarteduhub/internal/model/dto/request"
	"smarteduhub/internal/model/dto/response"
	"smarteduhub/internal/pkg/utils"
	danmakuService "smarteduhub/internal/service/danmaku"

	"github.com/gin-gonic/gin"
)

type handlerImpl struct {
	service danmakuService.Service
}

func NewHandler() Handler {
	return &handlerImpl{
		service: danmakuService.NewService(),
	}
}

func (h *handlerImpl) Send(c *gin.Context) {
	var req request.SendDanmakuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}

	uid, err := utils.GetLoginID(c)
	if err != nil {
		response.FailWithHTTPStatus(c, 401, -1, "unauthorized")
		return
	}

	if err := h.service.Send(uid, &req); err != nil {
		response.Fail(c, "发送失败")
		return
	}

	response.Success(c, nil)
}

func (h *handlerImpl) List(c *gin.Context) {
	var req request.ListDanmakuRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}

	list, err := h.service.List(req.ResourceID)
	if err != nil {
		response.Fail(c, "获取弹幕失败")
		return
	}

	response.Success(c, list)
}
