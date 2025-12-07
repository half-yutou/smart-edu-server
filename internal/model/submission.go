package model

import "time"

// Submission 提交记录表
type Submission struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	HomeworkID  int64     `gorm:"not null;uniqueIndex:idx_submissions_unique" json:"homework_id"`
	StudentID   int64     `gorm:"not null;uniqueIndex:idx_submissions_unique" json:"student_id"`
	Status      string    `gorm:"type:varchar(20);default:'submitted'" json:"status"`
	TotalScore  int       `gorm:"default:0" json:"total_score"`
	Feedback    string    `gorm:"type:text;not null;default:''" json:"feedback"`
	AIFeedback  string    `gorm:"type:text;not null;default:''" json:"ai_feedback"`
	SubmittedAt time.Time `gorm:"autoCreateTime" json:"submitted_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联
	Details []SubmissionDetail `gorm:"foreignKey:SubmissionID" json:"details,omitempty"`
	Student *User              `gorm:"foreignKey:StudentID" json:"student,omitempty"`
}