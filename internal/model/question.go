package model

import "time"

// Question 题目表
type Question struct {
	ID            int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	HomeworkID    int64     `gorm:"not null;index" json:"homework_id"`
	QuestionType  string    `gorm:"type:varchar(20);not null" json:"question_type"` // choice, text
	Content       string    `gorm:"type:text;not null" json:"content"`
	Options       string    `gorm:"type:jsonb" json:"options"` // JSON string: {"A":".."}
	CorrectAnswer string    `gorm:"type:text" json:"correct_answer"`
	Score         int       `gorm:"not null" json:"score"`
	OrderNum      int       `gorm:"default:0" json:"order_num"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
}