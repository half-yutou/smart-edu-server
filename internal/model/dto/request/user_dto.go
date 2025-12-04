package request

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=student teacher admin"`
	Nickname string `json:"nickname"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	ID    int64  `json:"id"`
	Token string `json:"token"`
	Role  string `json:"role"`
	Name  string `json:"name"`
}

type UpdateProfileRequest struct {
	Nickname  string `json:"nickname" binding:"max=50"`
	AvatarURL string `json:"avatar_url" binding:"omitempty,url"` // 校验是否为有效URL，允许为空
}
