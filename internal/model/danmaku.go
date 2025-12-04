package model

import "time"

// Danmaku 弹幕表
type Danmaku struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ResourceID int64     `gorm:"not null;index:idx_danmakus_resource_time" json:"resource_id"`
	UserID     int64     `gorm:"not null" json:"user_id"`
	Content    string    `gorm:"type:varchar(255);not null" json:"content"`
	TimePoint  float64   `gorm:"not null;index:idx_danmakus_resource_time" json:"time_point"`
	Color      string    `gorm:"type:varchar(10);default:'#FFFFFF'" json:"color"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}