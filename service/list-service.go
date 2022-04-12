package service

// mockgen -source=service/list-service.go -destination=./mock_service/list-service.go

import (
	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/dto"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/repository"
)

type listService struct {
	rep repository.ListRepository
}

type ListService interface {
	Create(*gin.Context) (model.List, error)
}

func NewListService() ListService {
	return &listService{rep: repository.NewListRepository()}
}

func (s *listService) Create(ctx *gin.Context) (model.List, error) {
	var dtoList dto.List
	err := ctx.ShouldBindJSON(&dtoList)
	if err != nil {
		return model.List{}, err
	}

	var list model.List
	dtoList.Transfer(&list)
	currentUser := ctx.MustGet("currentUser").(model.User)
	err = s.rep.Create(&currentUser, &list)
	if err != nil {
		return model.List{}, err
	}

	return list, nil
}

// test
func TestNewListService(listRepository repository.ListRepository) ListService {
	return &listService{rep: listRepository}
}
