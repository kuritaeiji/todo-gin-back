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
	Update(list *model.List, updatingList model.List) error
	Destroy(*model.List) error
	Find(id int) (model.List, error)
	FindLists(*model.User) error
}

func NewListRepository() ListRepository {
	return &listRepository{db: db.GetDB()}
}

func (r *listRepository) Create(user *model.User, list *model.List) error {
	return r.db.Model(user).Association("Lists").Append(list)
}

func (r *listRepository) Update(list *model.List, updatingList model.List) error {
	return r.db.Model(&list).Select("title").Updates(updatingList).Error
}

func (r *listRepository) Destroy(list *model.List) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := r.db.Delete(&list).Error
		if err != nil {
			return err
		}

		err = r.db.Model(model.List{}).Where("lists.index > ? AND lists.user_id = ?", list.Index, list.UserID).Updates(map[string]interface{}{"index": gorm.Expr("lists.index - ?", 1)}).Error
		if err != nil {
			return err
		}

		return nil
	})
}

func (r *listRepository) Find(id int) (model.List, error) {
	var list model.List
	err := r.db.First(&list, id).Error
	return list, err
}

func (r *listRepository) FindLists(user *model.User) error {
	return r.db.Where(model.List{UserID: user.ID}).Order("lists.index ASC").Find(&user.Lists).Error
}
