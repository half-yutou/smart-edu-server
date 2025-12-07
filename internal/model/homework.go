package model

import "time"

// Homework 作业表
type Homework struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     string    `gorm:"type:varchar(100);not null" json:"title"`
	ClassID   int64     `gorm:"not null;index" json:"class_id"`
	CreatorID int64     `gorm:"not null" json:"creator_id"`
	Deadline  *time.Time `json:"deadline"` // 可空
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联
	Class     *Class     `gorm:"foreignKey:ClassID" json:"class,omitempty"`
	Questions []Question `gorm:"foreignKey:HomeworkID" json:"questions,omitempty"`
}