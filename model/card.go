package model

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Card struct {
	gorm.Model
	ID     int    `gorm:"primaryKey;autoIncrement;not null"`
	Title  string `gorm:"type:varchar(100)"`
	Index  int
	ListID int
	List   List `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (card *Card) ToJson() gin.H {
	return gin.H{
		"id":    card.ID,
		"title": card.Title,
	}
}

func ToJsonCardSlice(cards []Card) []gin.H {
	jsonCardSlice := make([]gin.H, 0, len(cards))
	for _, card := range cards {
		jsonCardSlice = append(jsonCardSlice, card.ToJson())
	}
	return jsonCardSlice
}
