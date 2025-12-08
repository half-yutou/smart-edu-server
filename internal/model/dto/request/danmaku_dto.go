package request

type SendDanmakuRequest struct {
	ResourceID int64   `json:"resource_id,string" binding:"required"`
	Content    string  `json:"content" binding:"required,max=255"`
	Time       float64 `json:"time" binding:"required,gte=0"`
	Color      string  `json:"color" binding:"omitempty,hexcolor"`
	Type       int     `json:"type" binding:"omitempty,oneof=0 1 2"`
}

type ListDanmakuRequest struct {
	ResourceID int64 `json:"resource_id,string" form:"resource_id" binding:"required"`
}
