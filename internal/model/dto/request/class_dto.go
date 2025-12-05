package request

type CreateClassRequest struct {
	Name      string `json:"name" binding:"required,max=100"`
	SubjectID int    `json:"subject_id,string" binding:"required"`
	GradeID   int    `json:"grade_id,string" binding:"required"`
}

type DeleteClassRequest struct {
	ClassID int64 `json:"class_id,string" binding:"required"`
}

type UpdateClassRequest struct {
	ClassID   int64   `json:"class_id,string" binding:"required"`
	Name      *string `json:"name" binding:"omitempty,max=100"`
	GradeID   *int    `json:"grade_id,string" binding:"omitempty"`
	SubjectID *int    `json:"subject_id,string" binding:"omitempty"`
}

type GetClassByCodeRequest struct {
	Code string `json:"code" binding:"required,len=6"`
}

type GetClassByIDRequest struct {
	ID int64 `json:"id,string" binding:"required"`
}

type JoinClassByCodeRequest struct {
	Code string `json:"code" binding:"required,len=6"`
}

type QuitClassRequest struct {
	ClassID int64 `json:"class_id,string" binding:"required"`
}
