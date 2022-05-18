package service

// mockgen -source=service/list-service.go -destination=./mock_service/list-service.go

import (
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
	Destroy(*gin.Context) error
	Move(*gin.Context) error
}

func NewListService() ListService {
	return &listService{rep: repository.NewListRepository()}
}

func (s *listService) Index(ctx *gin.Context) ([]model.List, error) {
	currentUser := ctx.MustGet(config.CurrentUserKey).(model.User)
	err := s.rep.FindListsWithCards(&currentUser)
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
	var dtoList dto.List
	err := ctx.ShouldBindJSON(&dtoList)
	if err != nil {
		return model.List{}, err
	}

	var updatingList model.List
	dtoList.Transfer(&updatingList)

	list := ctx.MustGet(config.ListKey).(model.List)
	err = s.rep.Update(&list, updatingList)
	if err != nil {
		return list, err
	}

	return list, nil
}

func (s *listService) Destroy(ctx *gin.Context) error {
	list := ctx.MustGet(config.ListKey).(model.List)
	return s.rep.Destroy(&list)
}

func (s *listService) Move(ctx *gin.Context) error {
	var moveList dto.MoveList
	err := ctx.ShouldBindJSON(&moveList)
	if err != nil {
		return err
	}

	list := ctx.MustGet(config.ListKey).(model.List)
	currentUser := ctx.MustGet(config.CurrentUserKey).(model.User)
	return s.rep.Move(&list, moveList.Index, &currentUser)
}

// test
func TestNewListService(listRepository repository.ListRepository) ListService {
	return &listService{rep: listRepository}
}
