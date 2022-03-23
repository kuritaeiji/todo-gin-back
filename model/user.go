package model

import "time"

type User struct {
	ID             int       `json:"id"`
	Email          string    `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	PasswordDigest string    `gorm:"type:varchar(256)" json:"passwordDigest"`
	Activated      bool      `gorm:"default:false" json:"activatedAt"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
