package model

import "time"

// ClassResource 班级资源关联表
type ClassResource struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ClassID    int64     `gorm:"not null;uniqueIndex:idx_class_resources_unique" json:"class_id"`
	ResourceID int64     `gorm:"not null;uniqueIndex:idx_class_resources_unique" json:"resource_id"`
	AddedAt    time.Time `gorm:"autoCreateTime" json:"added_at"`
}