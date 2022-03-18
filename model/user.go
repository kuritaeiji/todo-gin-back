package model

import "time"

type User struct {
	ID             int
	Email          string `gorm:"unique_index"`
	PasswordDigest string
	Activated      bool `gorm:"default:false"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
