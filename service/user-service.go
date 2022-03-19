package service

import (
	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/db"
	"github.com/kuritaeiji/todo-gin-back/dto"
	"github.com/kuritaeiji/todo-gin-back/model"
	"gorm.io/gorm"
)

type UserService interface {
	Create(ctx *gin.Context) (model.User, error)
	IsUniqueEmail(ctx *gin.Context) bool
}

type userService struct {
	db *gorm.DB
}

func NewUserService() UserService {
	return &userService{db.GetDB()}
}

func (s *userService) Create(ctx *gin.Context) (model.User, error) {
	var userProxy dto.User
	if err := ctx.ShouldBindJSON(&userProxy); err != nil {
		return model.User{}, err
	}
	var user model.User
	userProxy.Transfer(&user)
	s.db.Create(&user)
	return user, nil
}

func (s *userService) IsUniqueEmail(ctx *gin.Context) bool {
	var user dto.UniqueUser
	if err := ctx.ShouldBindQuery(&user); err != nil {
		return false
	}
	return true
}
