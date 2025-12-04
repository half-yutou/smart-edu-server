package model

import "time"

// Resource 资源表
type Resource struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string    `gorm:"type:varchar(100);not null" json:"title"`
	Description string    `gorm:"type:text;not null;default:''" json:"description"`
	ResType     string    `gorm:"type:varchar(20);not null" json:"res_type"` // video, document
	FileURL     string    `gorm:"type:varchar(255);not null" json:"file_url"`
	SubjectID   int       `gorm:"not null;index:idx_resources_subject_grade" json:"subject_id"`
	GradeID     int       `gorm:"not null;index:idx_resources_subject_grade" json:"grade_id"`
	UploaderID  int64     `gorm:"not null;index" json:"uploader_id"`
	Duration    int       `gorm:"default:0" json:"duration"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}