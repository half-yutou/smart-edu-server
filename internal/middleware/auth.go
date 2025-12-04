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
		if len(token) > 7 && strings.ToUpper(token[0:7]) == "BEARER " {
			token = token[7:]
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
