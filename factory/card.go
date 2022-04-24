package factory

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/dto"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/repository"
)

type CardConfig struct {
	Title              string
	Index              int
	NotUseDefaultValue bool
}

const defaultCardTitle = "card title"

var index = 0

func (cardConfig *CardConfig) setDefaultValue() {
	if cardConfig.NotUseDefaultValue {
		return
	}

	if cardConfig.Title == "" {
		cardConfig.Title = "card title"
	}

	if cardConfig.Index == 0 {
		cardConfig.Index = index
		index++
	}
}

func NewDtoCard(cardConfig *CardConfig) dto.Card {
	cardConfig.setDefaultValue()
	return dto.Card{
		Title: cardConfig.Title,
		Index: cardConfig.Index,
	}
}

func NewCard(cardConfig *CardConfig) model.Card {
	dtoCard := NewDtoCard(cardConfig)
	var card model.Card
	dtoCard.Transfer(&card)
	return card
}

func CreateCard(cardConfig *CardConfig, list model.List) model.Card {
	card := NewCard(cardConfig)
	repository.NewCardRepository().Create(&card, &list)
	return card
}

func CreateCardRequestBody(cardConfig *CardConfig) io.Reader {
	cardConfig.setDefaultValue()
	body := gin.H{
		"title": cardConfig.Title,
		"index": cardConfig.Index,
	}
	json, _ := json.Marshal(body)
	return strings.NewReader(string(json))
}