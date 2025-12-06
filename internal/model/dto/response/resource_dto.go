package response

import "time"

type ResourceInfo struct {
	ID           int64     `json:"id,string"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	ResType      string    `json:"res_type"`
	FileURL      string    `json:"file_url"`
	SubjectID    int       `json:"subject_id"`
	SubjectName  string    `json:"subject_name"`
	GradeID      int       `json:"grade_id"`
	GradeName    string    `json:"grade_name"`
	UploaderID   int64     `json:"uploader_id,string"`
	UploaderName string    `json:"uploader_name"`
	Duration     int       `json:"duration"`
	CreatedAt    time.Time `json:"created_at"`
}

type PageResult struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}
