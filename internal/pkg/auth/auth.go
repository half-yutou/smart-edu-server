package auth

import (
	"github.com/click33/sa-token-go/core"
	"github.com/click33/sa-token-go/storage/memory"
	"github.com/click33/sa-token-go/stputil"
)

func InitSaToken() {
	stputil.SetManager(
		core.NewBuilder().
			Storage(memory.NewStorage()).
			TokenStyle(core.TokenStyleJWT). // 使用 JWT
			JwtSecretKey("jwt-secret-key"). // JWT 密钥
			Timeout(24 * 60 * 60).          // 24小时过期
			ActiveTimeout(12 * 60 * 60).    // 12小时无操作过期
			MaxRefresh(60 * 60).            // 1小时自动续期
			Build(),
	)
}
