package service

// mockgen -source=service/card-middleware-service.go -destination=mock_service/card-middleware-service.go

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/repository"
)

type cardMiddlewareService struct {
	repository     repository.CardRepository
	userRepository repository.UserRepository
}

type CardMiddlewareService interface {
	Authorize(*gin.Context) (model.Card, error)
}

func NewCardMiddlewareService() CardMiddlewareService {
	return &cardMiddlewareService{repository: repository.NewCardRepository(), userRepository: repository.NewUserRepository()}
}

func (s *cardMiddlewareService) Authorize(ctx *gin.Context) (model.Card, error) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return model.Card{}, err
	}

	card, err := s.repository.Find(id)
	if err != nil {
		return card, err
	}

	currentUser := ctx.MustGet(config.CurrentUserKey).(model.User)
	hasCard, err := s.userRepository.HasCard(card, currentUser)
	if err != nil {
		return card, err
	}
	if !hasCard {
		return card, config.ForbiddenError
	}

	return card, nil
}

// test
func TestNewCardMiddlewareService(cardRepository repository.CardRepository, userRepository repository.UserRepository) CardMiddlewareService {
	return &cardMiddlewareService{repository: cardRepository, userRepository: userRepository}
}
