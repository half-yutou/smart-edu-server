package class

import (
	"smarteduhub/internal/model/dto/request"
	"smarteduhub/internal/model/dto/response"
	"smarteduhub/internal/pkg/utils"
	"smarteduhub/internal/pkg/validator"
	classService "smarteduhub/internal/service/class"

	"github.com/gin-gonic/gin"
)

type handlerImpl struct {
	service classService.Service
}

var _ Handler = (*handlerImpl)(nil)

func NewHandler() Handler {
	return &handlerImpl{
		service: classService.NewService(),
	}
}

// region 教师-班级操作

func (h *handlerImpl) Create(c *gin.Context) {
	var req request.CreateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, validator.Translate(err))
		return
	}

	uid, err := utils.GetLoginID(c)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	class, err := h.service.Create(uid, &req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.SuccessWithMsg(c, "class created successfully", class)
}

func (h *handlerImpl) ListForTeacher(c *gin.Context) {
	uid, err := utils.GetLoginID(c)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	classes, err := h.service.ListForTeacher(uid)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.Success(c, classes)
}

func (h *handlerImpl) DeleteForTeacherByID(c *gin.Context) {
	var req request.DeleteClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, validator.Translate(err))
		return
	}

	uid, uidErr := utils.GetLoginID(c)
	if uidErr != nil {
		response.Fail(c, uidErr.Error())
		return
	}

	if err := h.service.DeleteByID(uid, req.ClassID); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.SuccessWithMsg(c, "class deleted successfully", nil)
}

func (h *handlerImpl) UpdateForTeacherByID(c *gin.Context) {
	var req request.UpdateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, validator.Translate(err))
		return
	}

	uid, uidErr := utils.GetLoginID(c)
	if uidErr != nil {
		response.Fail(c, uidErr.Error())
		return
	}

	if err := h.service.UpdateByID(uid, &req); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.SuccessWithMsg(c, "class updated successfully", nil)
}

// endregion

// region 单个班级查询

func (h *handlerImpl) GetByCode(c *gin.Context) {
	var req request.GetClassByCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, validator.Translate(err))
		return
	}

	class, err := h.service.GetByCode(req.Code)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.Success(c, class)
}

func (h *handlerImpl) GetByID(c *gin.Context) {
	var req request.GetClassByIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, validator.Translate(err))
		return
	}

	uid, uidErr := utils.GetLoginID(c)
	if uidErr != nil {
		response.Fail(c, uidErr.Error())
		return
	}

	class, err := h.service.GetByID(uid, req.ID)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.Success(c, class)
}

// endregion

// region 学生-班级操作

// ListForStudent 获取学生加入的班级
func (h *handlerImpl) ListForStudent(c *gin.Context) {
	uid, err := utils.GetLoginID(c)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	classes, err := h.service.ListForStudent(uid)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.Success(c, classes)
}

// Quit 学生退出班级
func (h *handlerImpl) Quit(c *gin.Context) {
	var req request.QuitClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, validator.Translate(err))
		return
	}

	uid, err := utils.GetLoginID(c)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	if err := h.service.Quit(uid, req.ClassID); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.SuccessWithMsg(c, "class quit successfully", nil)
}

// JoinByCode 学生通过邀请码加入班级
func (h *handlerImpl) JoinByCode(c *gin.Context) {
	var req request.JoinClassByCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, validator.Translate(err))
		return
	}

	uid, err := utils.GetLoginID(c)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	if err := h.service.JoinByCode(uid, req.Code); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.SuccessWithMsg(c, "class joined successfully", nil)
}

// region 班级资源操作

func (h *handlerImpl) AddResource(c *gin.Context) {
	var req request.AddResourceToClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, validator.Translate(err))
		return
	}

	uid, err := utils.GetLoginID(c)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	if err := h.service.AddResource(uid, &req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.SuccessWithMsg(c, "resource added to class successfully", nil)
}

func (h *handlerImpl) RemoveResource(c *gin.Context) {
	var req request.RemoveResourceFromClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, validator.Translate(err))
		return
	}

	uid, err := utils.GetLoginID(c)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	if err := h.service.RemoveResource(uid, &req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.SuccessWithMsg(c, "resource removed from class successfully", nil)
}

func (h *handlerImpl) ListResources(c *gin.Context) {
	var req request.ListClassResourcesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, validator.Translate(err))
		return
	}

	uid, err := utils.GetLoginID(c)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	result, err := h.service.ListResources(uid, &req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, result)
}

// endregion

// endregion
