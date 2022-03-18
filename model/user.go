package model

import "time"

type User struct {
	ID             int
	Email          string `gorm:"type:varchar(100);uniqueIndex"`
	PasswordDigest string `gorm:"type:varchar(50)"`
	Activated      bool   `gorm:"default:false"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
