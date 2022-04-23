package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID             int    `gorm:"primary_key;AUTO_INCREMANT;not null" json:"id"`
	Email          string `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	PasswordDigest string `gorm:"type:varchar(256)" json:"passwordDigest"`
	Activated      bool   `gorm:"default:false" json:"activatedAt"`

	Lists []List
}

func (user *User) Authenticate(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password))
	if err == nil {
		return true
	}

	return false
}

func (user *User) HasList(list List) bool {
	return user.ID == list.UserID
}
