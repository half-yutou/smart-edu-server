package model

import "time"

// ClassMember 班级成员表
type ClassMember struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ClassID   int64     `gorm:"not null;uniqueIndex:idx_class_members_unique" json:"class_id"`
	StudentID int64     `gorm:"not null;uniqueIndex:idx_class_members_unique;index" json:"student_id"`
	JoinedAt  time.Time `gorm:"autoCreateTime" json:"joined_at"`
}