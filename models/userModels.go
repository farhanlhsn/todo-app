package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Nama     string
	Email    string `gorm:"unique"`
	Password string
	IsLoggedIn bool `gorm:"default:false"`
}
