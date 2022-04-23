package dto

import "github.com/kuritaeiji/todo-gin-back/model"

type Card struct {
	Title string `json:"title" binding:"required,max=100"`
	Index int    `json:"index" binding:"gte=0"`
}

func (dtoCard Card) Transfer(card *model.Card) {
	card.Title = dtoCard.Title
	card.Index = dtoCard.Index
}
