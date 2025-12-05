package utils

import (
	"strconv"

	"github.com/click33/sa-token-go/stputil"
	"github.com/gin-gonic/gin"
)

func GetLoginID(c *gin.Context) (int64, error) {
	uidStr, err := stputil.GetLoginID(c.GetHeader("Authorization"))
	if err != nil {
		return 0, err
	}
	uidInt, err := strconv.Atoi(uidStr)
	if err != nil {
		return 0, err
	}
	return int64(uidInt), nil
}
