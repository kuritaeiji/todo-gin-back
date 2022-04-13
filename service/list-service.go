package service

// mockgen -source=service/list-service.go -destination=./mock_service/list-service.go

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/dto"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/repository"
)

type listService struct {
	rep repository.ListRepository
}

type ListService interface {
	Index(*gin.Context) ([]model.List, error)
	Create(*gin.Context) (model.List, error)
	Update(*gin.Context) (model.List, error)
}

func NewListService() ListService {
	return &listService{rep: repository.NewListRepository()}
}

func (s *listService) Index(ctx *gin.Context) ([]model.List, error) {
	currentUser := ctx.MustGet(config.CurrentUserKey).(model.User)
	err := s.rep.FindLists(&currentUser)
	return currentUser.Lists, err
}

func (s *listService) Create(ctx *gin.Context) (model.List, error) {
	var dtoList dto.List
	err := ctx.ShouldBindJSON(&dtoList)
	if err != nil {
		return model.List{}, err
	}

	var list model.List
	dtoList.Transfer(&list)
	currentUser := ctx.MustGet(config.CurrentUserKey).(model.User)
	err = s.rep.Create(&currentUser, &list)
	if err != nil {
		return model.List{}, err
	}

	return list, nil
}

func (s *listService) Update(ctx *gin.Context) (model.List, error) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return model.List{}, err
	}

	list, err := s.rep.Find(id)
	if err != nil {
		return list, err
	}

	currentUser := ctx.MustGet(config.CurrentUserKey).(model.User)
	if !currentUser.HasList(list) {
		return list, config.ForbiddenError
	}

	var dtoList dto.List
	err = ctx.ShouldBindJSON(&dtoList)
	if err != nil {
		return list, err
	}

	var updatingList model.List
	dtoList.Transfer(&updatingList)
	err = s.rep.Update(&list, updatingList)
	if err != nil {
		return list, err
	}

	return list, nil
}

// test
func TestNewListService(listRepository repository.ListRepository) ListService {
	return &listService{rep: listRepository}
}
