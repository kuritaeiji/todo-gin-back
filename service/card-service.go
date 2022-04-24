package service

// mockgen -source=service/card-service.go -destination=./mock_service/card-service.go

import (
	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/dto"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/repository"
)

type cardService struct {
	repository repository.CardRepository
}

type CardService interface {
	Create(*gin.Context) (model.Card, error)
}

func NewCardService() CardService {
	return &cardService{repository: repository.NewCardRepository()}
}

func (s *cardService) Create(ctx *gin.Context) (model.Card, error) {
	var cardDto dto.Card
	err := ctx.ShouldBindJSON(&cardDto)
	if err != nil {
		return model.Card{}, err
	}

	var card model.Card
	cardDto.Transfer(&card)

	list := ctx.MustGet(config.ListKey).(model.List)
	err = s.repository.Create(&card, &list)
	return card, err
}

// test
func TestNewCardService(cardRepository repository.CardRepository) CardService {
	return &cardService{repository: cardRepository}
}