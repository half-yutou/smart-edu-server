package utils

import (
	"strconv"
	"strings"

	"github.com/click33/sa-token-go/stputil"
	"github.com/gin-gonic/gin"
)

func GetLoginID(c *gin.Context) (int64, error) {
	token := c.GetHeader("Authorization")
	// 处理 Bearer 前缀
	// 使用 strings.Fields 自动处理空格，更健壮
	fields := strings.Fields(token)
	if len(fields) == 2 && strings.ToUpper(fields[0]) == "BEARER" {
		token = fields[1]
	} else if len(fields) == 1 {
		// 前端只传了 token 没有 Bearer 前缀，直接使用
		token = fields[0]
	}
	uidStr, err := stputil.GetLoginID(token)
	if err != nil {
		return 0, err
	}
	uidInt, err := strconv.Atoi(uidStr)
	if err != nil {
		return 0, err
	}
	return int64(uidInt), nil
}
