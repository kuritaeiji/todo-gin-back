package service

// mockgen -source=service/list-middleware-service.go -destination=./mock_service/list-middleware-service.go

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/repository"
)

type listMiddlewareServive struct {
	repository repository.ListRepository
}

type ListMiddlewareServive interface {
	Authorize(*gin.Context) (model.List, error)
	FindAndAuthorizeList(id int, currentUser model.User) (model.List, error)
}

func NewListMiddlewareService() ListMiddlewareServive {
	return &listMiddlewareServive{repository: repository.NewListRepository()}
}

func (s *listMiddlewareServive) Authorize(ctx *gin.Context) (model.List, error) {
	var (
		id  int
		err error
	)

	if idString := ctx.Param("listID"); idString != "" {
		id, err = strconv.Atoi(idString)
	} else {
		id, err = strconv.Atoi(ctx.Param("id"))
	}

	if err != nil {
		return model.List{}, err
	}

	currentUser := ctx.MustGet(config.CurrentUserKey).(model.User)
	return s.FindAndAuthorizeList(id, currentUser)
}

func (s *listMiddlewareServive) FindAndAuthorizeList(id int, currentUser model.User) (model.List, error) {
	list, err := s.repository.Find(id)
	if err != nil {
		return model.List{}, err
	}

	if !currentUser.HasList(list) {
		return model.List{}, config.ForbiddenError
	}

	return list, nil
}

// test
func TestNewListMiddlewareService(listRepository repository.ListRepository) ListMiddlewareServive {
	return &listMiddlewareServive{repository: listRepository}
}
