package model

import "time"

// Danmaku 弹幕表
type Danmaku struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ResourceID int64     `gorm:"not null;index" json:"resource_id"`        // 关联资源
	UserID     int64     `gorm:"column:user_id;not null" json:"sender_id"` // 发送者 (数据库字段: user_id)
	Content    string    `gorm:"type:varchar(255);not null" json:"content"`
	TimePoint  float64   `gorm:"column:time_point;not null" json:"time"` // 视频内时间点(秒) (数据库字段: time_point)
	Color      string    `gorm:"type:varchar(10);default:'#ffffff'" json:"color"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`

	// 关联
	Sender *User `gorm:"foreignKey:UserID" json:"sender,omitempty"`
}
