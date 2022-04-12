package repository

// mockgen -source=repository/list-repository.go -destination=./mock_repository/list-repository.go

import (
	"github.com/kuritaeiji/todo-gin-back/db"
	"github.com/kuritaeiji/todo-gin-back/model"
	"gorm.io/gorm"
)

type listRepository struct {
	db *gorm.DB
}

type ListRepository interface {
	Create(*model.User, *model.List) error
}

func NewListRepository() ListRepository {
	return &listRepository{db: db.GetDB()}
}

func (r *listRepository) Create(user *model.User, list *model.List) error {
	return r.db.Model(user).Association("Lists").Append(list)
}
