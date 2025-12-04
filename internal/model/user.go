package model

import "time"

// User 用户表
type User struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"type:varchar(50);not null;unique" json:"username"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"` // 密码不返回 JSON
	Role      string    `gorm:"type:varchar(20);not null" json:"role"` // student, teacher, admin
	Nickname  string    `gorm:"type:varchar(50);not null;default:''" json:"nickname"`
	AvatarURL string    `gorm:"type:varchar(255);not null;default:''" json:"avatar_url"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}