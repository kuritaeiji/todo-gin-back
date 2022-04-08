package dto

import (
	"github.com/kuritaeiji/todo-gin-back/model"
	// "gorm.io/gorm"
)

type List struct {
	Title  string `json:"title" binding:"required,max=50"`
	UserID int    `json:"userID" binding:"required"`
}

func (dtoList List) Transfer(list *model.List) {
	list.Title = dtoList.Title
	list.UserID = dtoList.UserID
}
