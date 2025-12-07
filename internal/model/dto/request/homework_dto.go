package request

import "time"

type QuestionItem struct {
	ID            int64  `json:"id,string"` // 更新时使用，新增为0
	QuestionType  string `json:"question_type" binding:"required,oneof=choice text"`
	Content       string `json:"content" binding:"required"`
	Options       string `json:"options"` // JSON string, e.g. "{\"A\":\"x\", \"B\":\"y\"}"
	CorrectAnswer string `json:"correct_answer"`
	Score         int    `json:"score" binding:"required,min=1"`
}

type CreateHomeworkRequest struct {
	ClassID   int64          `json:"class_id,string" binding:"required"`
	Title     string         `json:"title" binding:"required"`
	Deadline  *time.Time     `json:"deadline"`
	Questions []QuestionItem `json:"questions" binding:"required,dive"`
}

type DeleteHomeworkRequest struct {
	HomeworkID int64 `json:"homework_id,string" binding:"required"`
}

type UpdateHomeworkRequest struct {
	HomeworkID int64          `json:"homework_id,string" binding:"required"`
	Title      string         `json:"title" binding:"required"`
	Deadline   *time.Time     `json:"deadline"`
	Questions  []QuestionItem `json:"questions" binding:"required,dive"`
}

type ListHomeworksRequest struct {
	ClassID int64 `json:"class_id,string" form:"class_id"`
}

type SubmissionDetailItem struct {
	QuestionID  int64  `json:"question_id,string" binding:"required"`
	ImageURL    string `json:"image_url"`
	TextContent string `json:"text_content"`
}

type SubmitHomeworkRequest struct {
	HomeworkID int64                  `json:"homework_id,string" binding:"required"`
	Details    []SubmissionDetailItem `json:"details" binding:"required,dive"`
}

type ManualGradeRequest struct {
	SubmissionID int64 `json:"submission_id,string" binding:"required"`
	Details      []struct {
		DetailID int64  `json:"detail_id,string" binding:"required"`
		Score    int    `json:"score"`
		Comment  string `json:"comment"`
	} `json:"details" binding:"required,dive"`
}