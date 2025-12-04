package errno

// 错误码定义
const (
	Success = 0
	Error   = -1

	// 100xx 通用错误
	InvalidParams = 10001 // 参数错误
	AuthFailed    = 10002 // 认证失败

	// 200xx 用户模块错误
	UserAlreadyExists = 20001 // 用户已存在
	UserNotFound      = 20002 // 用户不存在或密码错误
)

// 错误码对应消息
var msgFlags = map[int]string{
	Success:           "success",
	Error:             "fail",
	InvalidParams:     "参数错误",
	AuthFailed:        "权限认证失败-请登录或权限不足",
	UserAlreadyExists: "用户已存在",
	UserNotFound:      "用户不存在或密码错误",
}

// GetMsg 获取错误码对应的消息
func GetMsg(code int) string {
	msg, ok := msgFlags[code]
	if ok {
		return msg
	}
	return msgFlags[Error]
}
