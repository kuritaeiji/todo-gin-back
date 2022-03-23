package repository

// mockgen -source=repository/user-repository.go -destination=mock_repository/user-repository.go

import (
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/db"
	"github.com/kuritaeiji/todo-gin-back/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	IsUnique(email string) (bool, error)
	Activate(user *model.User) error
	Find(id int) (model.User, error)
	FindByEmail(email string) (model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{
		db: db.GetDB(),
	}
}

func (r *userRepository) Create(user *model.User) error {
	if result, _ := r.IsUnique(user.Email); !result {
		return config.UniqueUserError
	}
	return r.db.Create(&user).Error
}

func (r *userRepository) IsUnique(email string) (bool, error) {
	var count int64
	err := r.db.Model(model.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func (r *userRepository) Activate(user *model.User) error {
	user.Activated = true
	return r.db.Save(user).Error
}

func (r *userRepository) Find(id int) (model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *userRepository) FindByEmail(email string) (model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
