package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一返回结构体
type Response struct {
	Code int         `json:"code"`           // 业务码，0表示成功，非0表示失败
	Msg  string      `json:"msg"`            // 提示信息
	Data interface{} `json:"data,omitempty"` // 返回数据，omitempty表示如果为空则不返回该字段
}

const (
	CodeSuccess = 0
	CodeError   = -1
)

// Result 返回通用响应
func Result(c *gin.Context, httpStatus int, code int, msg string, data interface{}) {
	c.JSON(httpStatus, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

// Success 返回成功响应
func Success(c *gin.Context, data interface{}) {
	Result(c, http.StatusOK, CodeSuccess, "success", data)
}

// SuccessWithMsg 返回带消息的成功响应
func SuccessWithMsg(c *gin.Context, msg string, data interface{}) {
	Result(c, http.StatusOK, CodeSuccess, msg, data)
}

// Fail 返回失败响应
func Fail(c *gin.Context, msg string) {
	Result(c, http.StatusOK, CodeError, msg, nil)
}

// FailWithCode 返回带业务码的失败响应
func FailWithCode(c *gin.Context, code int, msg string) {
	Result(c, http.StatusOK, code, msg, nil)
}

// FailWithHTTPStatus 返回自定义HTTP状态码的失败响应
func FailWithHTTPStatus(c *gin.Context, httpStatus int, code int, msg string) {
	Result(c, httpStatus, code, msg, nil)
}
