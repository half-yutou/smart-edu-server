package model

// SubmissionDetail 答题详情表
type SubmissionDetail struct {
	ID             int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	SubmissionID   int64  `gorm:"not null;index" json:"submission_id"`
	QuestionID     int64  `gorm:"not null" json:"question_id"`
	AnswerContent  string `gorm:"type:text;not null;default:''" json:"answer_content"`
	AnswerImageURL string `gorm:"type:varchar(255);not null;default:''" json:"answer_image_url"`
	IsCorrect      bool   `gorm:"default:false" json:"is_correct"`
	Score          int    `gorm:"default:0" json:"score"`
	Comment        string `gorm:"type:varchar(255);not null;default:''" json:"comment"`
}