package response

import "smarteduhub/internal/model"

type HomeworkInfo struct {
	*model.Homework
}

type HomeworkListItem struct {
	*model.Homework
	StudentStatus string `json:"student_status"` // "unsubmitted", "submitted", "graded"
}
