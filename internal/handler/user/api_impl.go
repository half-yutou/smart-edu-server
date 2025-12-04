package user

import (
	"strconv"

	"smarteduhub/internal/model/dto/request"
	"smarteduhub/internal/model/dto/response"
	"smarteduhub/internal/pkg/errno"
	"smarteduhub/internal/pkg/validator"
	userService "smarteduhub/internal/service/user"

	"github.com/click33/sa-token-go/stputil"
	"github.com/gin-gonic/gin"
)

type handlerImpl struct {
	userService userService.Service
}

var _ Handler = (*handlerImpl)(nil)

func NewHandler() Handler {
	return &handlerImpl{
		userService: userService.NewService(),
	}
}

func (h *handlerImpl) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 使用翻译器将验证错误转换为中文
		response.Fail(c, validator.Translate(err))
		return
	}

	if err := h.userService.Register(&req); err != nil {
		// 处理特定业务错误
		if err.Error() == "username already exists" {
			response.FailWithCode(c, errno.UserAlreadyExists, errno.GetMsg(errno.UserAlreadyExists))
			return
		}
		response.Fail(c, err.Error())
		return
	}

	response.SuccessWithMsg(c, "registration successful", nil)
}

func (h *handlerImpl) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, validator.Translate(err))
		return
	}

	resp, err := h.userService.Login(&req)
	if err != nil {
		// 这里也可以根据 error 类型返回具体的 code
		response.Fail(c, err.Error())
		return
	}

	response.Success(c, resp)
}

func (h *handlerImpl) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		response.Fail(c, "token is required")
		return
	}

	if err := h.userService.Logout(token); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.SuccessWithMsg(c, "logout successful", nil)
}

func (h *handlerImpl) UpdateProfile(c *gin.Context) {
	var req request.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, validator.Translate(err))
		return
	}

	// 从 Sa-Token 获取当前用户ID
	uid, _ := stputil.GetLoginID(c.GetHeader("Authorization"))
	uidInt, _ := strconv.Atoi(uid)

	if err := h.userService.UpdateProfile(int64(uidInt), &req); err != nil {
		response.Fail(c, err.Error())
		return
	}

	response.SuccessWithMsg(c, "profile updated successfully", nil)
}
