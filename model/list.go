package model

import (
	"github.com/gin-gonic/gin"
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

func (list *List) ToJson() gin.H {
	return gin.H{
		"id":     list.ID,
		"title":  list.Title,
		"index":  list.Index,
		"userID": list.UserID,
	}
}

func ToJsonListSlice(listSlice []List) []gin.H {
	jsonListSlice := make([]gin.H, 0, len(listSlice))
	for _, list := range listSlice {
		jsonListSlice = append(jsonListSlice, list.ToJson())
	}
	return jsonListSlice
}
