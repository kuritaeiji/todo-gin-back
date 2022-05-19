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
	Activate(user *model.User) error
	Destroy(user *model.User) error
	FindOrCreateByOpenID(openID string) (model.User, error)
	IsUnique(email string) (bool, error)
	Find(id int) (model.User, error)
	FindByEmail(email string) (model.User, error)
	HasCard(card model.Card, user model.User) (bool, error)
}

type userRepository struct {
	db             *gorm.DB
	listRepository ListRepository
}

func NewUserRepository() UserRepository {
	return &userRepository{
		db:             db.GetDB(),
		listRepository: NewListRepository(),
	}
}

func (r *userRepository) Create(user *model.User) error {
	if result, _ := r.IsUnique(user.Email); !result {
		return config.UniqueUserError
	}
	return r.db.Create(&user).Error
}

func (r *userRepository) Activate(user *model.User) error {
	user.Activated = true
	return r.db.Save(user).Error
}

func (r *userRepository) Destroy(user *model.User) error {
	err := r.listRepository.FindListsWithCards(user)
	if err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		err := r.listRepository.DestroyLists(&user.Lists, tx)
		if err != nil {
			return err
		}

		return r.db.Delete(user).Error
	})
}

func (r *userRepository) FindOrCreateByOpenID(openID string) (model.User, error) {
	var user model.User
	err := r.db.Model(model.User{}).Where("open_id = ?", openID).First(&user).Error

	// ユーザーが見つかった場合とエラーが発生した場合
	if err != nil && err != gorm.ErrRecordNotFound {
		return user, err
	}

	// ユーザーが見つからなかった場合
	if err == gorm.ErrRecordNotFound {
		user.OpenID = openID
		user.Activated = true
		err = r.db.Create(&user).Error
	}

	return user, err
}

func (r *userRepository) IsUnique(email string) (bool, error) {
	var count int64
	err := r.db.Model(model.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 0, nil
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

func (r *userRepository) HasCard(card model.Card, user model.User) (bool, error) {
	err := r.db.Joins("List").First(&card).Error
	if err != nil {
		return false, err
	}

	return card.List.UserID == user.ID, nil
}
