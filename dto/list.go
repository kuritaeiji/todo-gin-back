package dto

import (
	"github.com/kuritaeiji/todo-gin-back/model"
)

type List struct {
	Title string `json:"title" binding:"required,max=50"`
	Index int    `json:"index" binding:"gte=0"`
}

func (dtoList List) Transfer(list *model.List) {
	list.Title = dtoList.Title
	list.Index = dtoList.Index
}

type MoveList struct {
	Index int `json:"index" binding:"gte=0"`
}
