package request

type CreateResourceRequest struct {
	Title       string `json:"title" binding:"required,max=100"`
	Description string `json:"description" binding:"max=500"`
	ResType     string `json:"res_type" binding:"required,oneof=video document"`
	FileURL     string `json:"file_url" binding:"required,url"` // 前端上传OSS后拿到的URL
	SubjectID   int    `json:"subject_id" binding:"required"`
	GradeID     int    `json:"grade_id" binding:"required"`
	Duration    int    `json:"duration" binding:"omitempty,gte=0"` // 视频时长，单位秒
}

type UpdateResourceRequest struct {
	ResourceID  int64   `json:"resource_id,string" binding:"required"`
	Title       *string `json:"title" binding:"omitempty,max=100"`
	Description *string `json:"description" binding:"omitempty,max=500"`
	SubjectID   *int    `json:"subject_id" binding:"omitempty"`
	GradeID     *int    `json:"grade_id" binding:"omitempty"`
	ResType     *string `json:"res_type" binding:"omitempty,oneof=video document"`
	FileURL     *string `json:"file_url" binding:"omitempty,url"`
}

type DeleteResourceRequest struct {
	ResourceID int64 `json:"resource_id,string" binding:"required"`
}

type ListResourcesRequest struct {
	Page      int    `json:"page" form:"page,default=1" binding:"gte=1"`
	PageSize  int    `json:"page_size" form:"page_size,default=10" binding:"gte=1,lte=100"`
	SubjectID int    `json:"subject_id,string" form:"subject_id" binding:"omitempty"`
	GradeID   int    `json:"grade_id,string" form:"grade_id" binding:"omitempty"`
	ResType   string `json:"res_type" form:"res_type" binding:"omitempty,oneof=video document"`
	Keyword   string `json:"keyword" form:"keyword" binding:"omitempty"` // 搜索标题
}
