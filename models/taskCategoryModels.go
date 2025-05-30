package models

import "gorm.io/gorm"

type TaskCategory struct {
	gorm.Model
	Name      string `gorm:"not null;unique"`
	UserID    *uint  `gorm:"null"`
	User      User   `gorm:"foreignKey:UserID;references:id"`
	IsDefault bool   `gorm:"default:false"`
	Tasks     []Task `gorm:"foreignKey:CategoryID"`
}
