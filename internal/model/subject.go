package model

// Subject 学科表
type Subject struct {
	ID   int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"type:varchar(50);not null;unique" json:"name"`
}