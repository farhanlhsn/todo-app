package models

import (
	"time"

	"gorm.io/gorm"
)

type Priority string

const (
	None   Priority = "none"
	Low    Priority = "low"
	Medium Priority = "medium"
	High   Priority = "high"
)

type Task struct {
	gorm.Model
	Title       string       `gorm:"not null"`
	Description string       `gorm:"null"`
	IsCompleted bool         `gorm:"default:false"`
	UserID      uint         `gorm:"not null"` // Foreign key
	User        User         `gorm:"foreignKey:UserID;references:id"`
	DueDate     *time.Time   `gorm:"type:datetime"` 
	CategoryID  *uint        // Foreign key
	Category    TaskCategory `gorm:"foreignKey:CategoryID;references:id"`
	Priority    Priority     `gorm:"type:enum('none', 'low', 'medium', 'high');default:'none'"`
}
