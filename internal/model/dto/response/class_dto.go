package response

import (
	"time"
)

type ClassInfo struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	SubjectID   int       `json:"subject_id"`
	SubjectName string    `json:"subject_name"`
	GradeID     int       `json:"grade_id"`
	GradeName   string    `json:"grade_name"`
	TeacherID   int64     `json:"teacher_id"`
	TeacherName string    `json:"teacher_name"`
	MemberCount int       `json:"member_count"`
	CreatedAt   time.Time `json:"created_at"`
}
