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
	repository            repository.CardRepository
	listMiddlewareService ListMiddlewareServive
}

type CardService interface {
	Create(*gin.Context) (model.Card, error)
	Update(*gin.Context) (model.Card, error)
	Destroy(*gin.Context) error
	Move(*gin.Context) error
}

func NewCardService() CardService {
	return &cardService{repository: repository.NewCardRepository(), listMiddlewareService: NewListMiddlewareService()}
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

func (s *cardService) Update(ctx *gin.Context) (model.Card, error) {
	var dtoCard dto.Card
	err := ctx.ShouldBindJSON(&dtoCard)
	if err != nil {
		return model.Card{}, err
	}

	var updatingCard model.Card
	dtoCard.Transfer(&updatingCard)
	card := ctx.MustGet(config.CardKey).(model.Card)
	err = s.repository.Update(&card, &updatingCard)

	return card, err
}

func (s *cardService) Destroy(ctx *gin.Context) error {
	card := ctx.MustGet(config.CardKey).(model.Card)
	return s.repository.Destroy(&card)
}

func (s *cardService) Move(ctx *gin.Context) error {
	var dtoMoveCard dto.MoveCard
	err := ctx.ShouldBindJSON(&dtoMoveCard)

	if err != nil {
		return err
	}

	// カレントユーザーが移動した先のリストを所有しているか確認する
	currentUser := ctx.MustGet(config.CurrentUserKey).(model.User)
	_, err = s.listMiddlewareService.FindAndAuthorizeList(dtoMoveCard.ToListID, currentUser)
	if err != nil {
		return err
	}

	card := ctx.MustGet(config.CardKey).(model.Card)
	return s.repository.Move(&card, dtoMoveCard.ToListID, dtoMoveCard.ToIndex)
}

// test
func TestNewCardService(cardRepository repository.CardRepository, listMiddlewareService ListMiddlewareServive) CardService {
	return &cardService{repository: cardRepository, listMiddlewareService: listMiddlewareService}
}
