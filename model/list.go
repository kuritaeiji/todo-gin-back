package model

import (
	"gorm.io/gorm"
)

type List struct {
	gorm.Model
	ID     int    `gorm:"primary_key;AUTO_INCREMANT;not null" json:"id"`
	Title  string `gorm:"type:varchar(50);not null" json:"title"`
	Index  int    `json:"index"`
	UserID int    `json:"userID"`
	User   User
}
