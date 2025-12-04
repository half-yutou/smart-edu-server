package model

import "time"

// Class 班级表
type Class struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Code        string    `gorm:"type:varchar(10);not null;unique" json:"code"`
	TeacherID   int64     `gorm:"not null;index" json:"teacher_id"`
	SubjectID   int       `gorm:"default:0" json:"subject_id"`
	GradeID     int       `gorm:"default:0" json:"grade_id"`
	Description string    `gorm:"type:text;not null;default:''" json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}