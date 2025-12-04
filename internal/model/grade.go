package model

// Grade 年级表
type Grade struct {
	ID   int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"type:varchar(50);not null;unique" json:"name"`
}