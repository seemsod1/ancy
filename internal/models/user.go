package models

import "time"

type User struct {
	ID               int      `gorm:"primaryKey"`
	Username         string   `gorm:"size:255;not null,unique" json:"username"`
	Email            string   `gorm:"unique;not null"  json:"email"`
	Password         string   `gorm:"size:255;not null" json:"password"`
	ProfilePhotoPath string   `gorm:"size:255" `
	RoleID           int      `gorm:"default:1"`
	Role             UserRole `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type UserRole struct {
	ID   int    `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:255;not null;unique" json:"name"`
}
