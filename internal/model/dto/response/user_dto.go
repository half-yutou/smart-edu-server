package response

type LoginResponse struct {
	ID        int64  `json:"id"`
	Token     string `json:"token"`
	Role      string `json:"role"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
}
