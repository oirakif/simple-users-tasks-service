package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        int       `gorm:"primaryKey autoIncrement"`
	Name      string    `gorm:"type:varchar(255)"`
	Email     string    `gorm:"type:varchar(255) unique"  json:"email,omitempty"`
	Password  string    `gorm:"type:varchar(255)" json:"password,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:" autoUpdateTime"`
}

type Task struct {
	gorm.Model
	ID          int `gorm:"primaryKey autoIncrement"`
	UserID      int
	User        User      `gorm:"foreignKey:UserID"`
	Title       string    `gorm:"type:varchar(255)"`
	Description string    `gorm:"type:text"`
	Status      string    `gorm:"default:pending"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:" autoUpdateTime"`
}
