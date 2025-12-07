package homework

import (
	"strconv"

	"smarteduhub/internal/model/dto/request"
	"smarteduhub/internal/model/dto/response"
	"smarteduhub/internal/pkg/utils"
	"smarteduhub/internal/pkg/validator"
	homeworkService "smarteduhub/internal/service/homework"

	"github.com/gin-gonic/gin"
)

type handlerImpl struct {
	service homeworkService.Service
}

var _ Handler = (*handlerImpl)(nil)

func NewHandler() Handler {
	return &handlerImpl{
		service: homeworkService.NewService(),
	}
}

func (h *handlerImpl) Create(c *gin.Context) {
	var req request.CreateHomeworkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, validator.Translate(err))
		return
	}

	uid, err := utils.GetLoginID(c)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	if err := h.service.Create(uid, &req); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.SuccessWithMsg(c, "homework published successfully", nil)
}

func (h *handlerImpl) Delete(c *gin.Context) {
	var req request.DeleteHomeworkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, validator.Translate(err))
		return
	}

	uid, err := utils.GetLoginID(c)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	if err := h.service.Delete(uid, &req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.SuccessWithMsg(c, "homework deleted successfully", nil)
}

func (h *handlerImpl) Update(c *gin.Context) {
	var req request.UpdateHomeworkRequest
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
	response.SuccessWithMsg(c, "homework updated successfully", nil)
}

func (h *handlerImpl) GetByID(c *gin.Context) {
	// 支持 Query: ?id=123
	idStr := c.Query("id")
	if idStr == "" {
		response.Fail(c, "id is required")
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Fail(c, "invalid id")
		return
	}

	hw, err := h.service.GetByID(id)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, hw)
}

func (h *handlerImpl) ListByTeacher(c *gin.Context) {
	uid, err := utils.GetLoginID(c)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	list, err := h.service.ListByCreator(uid)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, list)
}

func (h *handlerImpl) ListByClass(c *gin.Context) {
	classIDStr := c.Query("class_id")
	if classIDStr == "" {
		response.Fail(c, "class_id is required")
		return
	}
	classID, err := strconv.ParseInt(classIDStr, 10, 64)
	if err != nil {
		response.Fail(c, "invalid class_id")
		return
	}

	uid, err := utils.GetLoginID(c)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	list, err := h.service.ListByClass(uid, classID)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, list)
}

func (h *handlerImpl) Submit(c *gin.Context) {
	var req request.SubmitHomeworkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, validator.Translate(err))
		return
	}

	uid, err := utils.GetLoginID(c)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	if err := h.service.Submit(uid, &req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.SuccessWithMsg(c, "homework submitted successfully", nil)
}

func (h *handlerImpl) GetSubmission(c *gin.Context) {
	hwIDStr := c.Query("homework_id")
	if hwIDStr == "" {
		response.Fail(c, "homework_id is required")
		return
	}
	hwID, err := strconv.ParseInt(hwIDStr, 10, 64)
	if err != nil {
		response.Fail(c, "invalid homework_id")
		return
	}

	uid, err := utils.GetLoginID(c)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	sub, err := h.service.GetSubmission(hwID, uid)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, sub)
}

func (h *handlerImpl) ListSubmissions(c *gin.Context) {
	hwIDStr := c.Query("homework_id")
	if hwIDStr == "" {
		response.Fail(c, "homework_id is required")
		return
	}
	hwID, err := strconv.ParseInt(hwIDStr, 10, 64)
	if err != nil {
		response.Fail(c, "invalid homework_id")
		return
	}

	uid, err := utils.GetLoginID(c)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	list, err := h.service.ListSubmissions(uid, hwID)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, list)
}

func (h *handlerImpl) GradeSubmission(c *gin.Context) {
	var req request.ManualGradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, validator.Translate(err))
		return
	}

	uid, err := utils.GetLoginID(c)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}

	if err := h.service.GradeSubmission(uid, &req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.SuccessWithMsg(c, "grade updated successfully", nil)
}
