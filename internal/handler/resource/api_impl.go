package resource

import (
	"smarteduhub/internal/model/dto/request"
	"smarteduhub/internal/model/dto/response"
	"smarteduhub/internal/pkg/utils"
	"smarteduhub/internal/pkg/validator"
	resourceService "smarteduhub/internal/service/resource"

	"github.com/gin-gonic/gin"
)

type handlerImpl struct {
	service resourceService.Service
}

var _ Handler = (*handlerImpl)(nil)

func NewHandler() Handler {
	return &handlerImpl{
		service: resourceService.NewService(),
	}
}

func (h *handlerImpl) Create(c *gin.Context) {
	var req request.CreateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, validator.Translate(err))
		return
	}

	uid, err := utils.GetLoginID(c)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	res, err := h.service.Create(uid, &req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.SuccessWithMsg(c, "resource created successfully", res)
}

func (h *handlerImpl) Update(c *gin.Context) {
	var req request.UpdateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, validator.Translate(err))
		return
	}

	uid, err := utils.GetLoginID(c)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	if err := h.service.Update(uid, &req); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.SuccessWithMsg(c, "resource updated successfully", nil)
}

func (h *handlerImpl) Delete(c *gin.Context) {
	var req request.DeleteResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, validator.Translate(err))
		return
	}

	uid, err := utils.GetLoginID(c)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	if err := h.service.Delete(uid, req.ResourceID); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.SuccessWithMsg(c, "resource deleted successfully", nil)
}

func (h *handlerImpl) GetByID(c *gin.Context) {
	type GetIDReq struct {
		ID int64 `json:"id,string" form:"id" binding:"required"`
	}
	var idReq GetIDReq

	// 绑定Query
	if err := c.ShouldBindQuery(&idReq); err != nil {
		response.Fail(c, validator.Translate(err))
		return
	}

	res, err := h.service.GetByID(idReq.ID)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.Success(c, res)
}

func (h *handlerImpl) List(c *gin.Context) {
	var req request.ListResourcesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, validator.Translate(err))
		return
	}

	pageResult, err := h.service.List(&req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.Success(c, pageResult)
}

func (h *handlerImpl) ListMyResources(c *gin.Context) {
	uid, err := utils.GetLoginID(c)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	list, err := h.service.ListMyResources(uid)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.Success(c, list)
}
