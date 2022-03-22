package service

// mockgen -source=service/user-service.go -destination=./mock_service/user-service.go

import (
	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/dto"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/repository"
)

type UserService interface {
	Create(*gin.Context) (model.User, error)
	IsUnique(*gin.Context) bool
}

type userService struct {
	repository repository.UserRepository
	dto        dto.User
}

func NewUserService() UserService {
	return &userService{
		repository: repository.NewUserRepository(),
		dto:        dto.User{},
	}
}

func (s *userService) Create(ctx *gin.Context) (model.User, error) {
	if err := ctx.ShouldBindJSON(&s.dto); err != nil {
		return model.User{}, err
	}
	var user model.User
	s.dto.Transfer(&user)
	if err := s.repository.Create(&user); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (s *userService) IsUnique(ctx *gin.Context) bool {
	result, _ := s.repository.IsUnique(ctx.Query("email"))
	return result
}

// test用
func TestNewUserService(r repository.UserRepository) UserService {
	return &userService{
		repository: r,
		dto:        dto.User{},
	}
}
