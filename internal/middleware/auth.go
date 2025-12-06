package middleware

import (
	"strings"

	"smarteduhub/internal/model/dto/response"
	"smarteduhub/internal/pkg/errno"

	"github.com/click33/sa-token-go/stputil"
	"github.com/gin-gonic/gin"
)

// Auth 认证中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取 Token
		token := c.GetHeader("Authorization")

		// 2. 处理 Bearer 前缀
		// 使用 strings.Fields 自动处理空格，更健壮
		fields := strings.Fields(token)
		if len(fields) == 2 && strings.ToUpper(fields[0]) == "BEARER" {
			token = fields[1]
		} else if len(fields) == 1 {
			// 也许前端只传了 token 没有 Bearer 前缀，直接使用
			token = fields[0]
		}

		// 3. 判空
		if token == "" {
			response.FailWithCode(c, errno.AuthFailed, errno.GetMsg(errno.AuthFailed))
			c.Abort()
			return
		}

		// 4. 校验 Token
		loginId, err := stputil.GetLoginID(token)
		if err != nil || loginId == "" {
			response.FailWithCode(c, errno.AuthFailed, errno.GetMsg(errno.AuthFailed))
			c.Abort()
			return
		}

		// 5. 将 loginId 存入上下文，方便后续使用
		c.Set("uid", loginId)

		c.Next()
	}
}

// AuthAdmin 管理员认证中间件
func AuthAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取 uid
		uid, exists := c.Get("uid")
		if !exists {
			response.FailWithCode(c, errno.AuthFailed, errno.GetMsg(errno.AuthFailed))
			c.Abort()
			return
		}

		// 2. 从会话中获取角色
		session, err := stputil.GetSession(uid.(string))
		if err != nil {
			response.FailWithCode(c, errno.AuthFailed, errno.GetMsg(errno.AuthFailed))
			c.Abort()
			return
		}

		// 3. 检查角色是否为管理员
		role, ok := session.Get("role")
		if !ok || role != "admin" {
			response.FailWithCode(c, errno.AuthFailed, errno.GetMsg(errno.AuthFailed))
			c.Abort()
			return
		}

		c.Next()
	}
}

// AuthTeacher 教师认证中间件
func AuthTeacher() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取 uid
		uid, exists := c.Get("uid")
		if !exists {
			response.FailWithCode(c, errno.AuthFailed, errno.GetMsg(errno.AuthFailed))
			c.Abort()
			return
		}

		// 2. 从会话中获取角色
		session, err := stputil.GetSession(uid.(string))
		if err != nil {
			response.FailWithCode(c, errno.AuthFailed, errno.GetMsg(errno.AuthFailed))
			c.Abort()
			return
		}

		// 3. 检查角色是否为教师
		role, ok := session.Get("role")
		if !ok || role != "teacher" {
			response.FailWithCode(c, errno.AuthFailed, errno.GetMsg(errno.AuthFailed))
			c.Abort()
			return
		}

		c.Next()
	}
}

// AuthStudent 学生认证中间件
func AuthStudent() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取 uid
		uid, exists := c.Get("uid")
		if !exists {
			response.FailWithCode(c, errno.AuthFailed, errno.GetMsg(errno.AuthFailed))
			c.Abort()
			return
		}

		// 2. 从会话中获取角色
		session, err := stputil.GetSession(uid.(string))
		if err != nil {
			response.FailWithCode(c, errno.AuthFailed, errno.GetMsg(errno.AuthFailed))
			c.Abort()
			return
		}

		// 3. 检查角色是否为学生
		role, ok := session.Get("role")
		if !ok || role != "student" {
			response.FailWithCode(c, errno.AuthFailed, errno.GetMsg(errno.AuthFailed))
			c.Abort()
			return
		}

		c.Next()
	}
}
