package repository

// mockgen -source=repository/list-repository.go -destination=./mock_repository/list-repository.go

import (
	"github.com/kuritaeiji/todo-gin-back/db"
	"github.com/kuritaeiji/todo-gin-back/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type listRepository struct {
	db *gorm.DB
}

type ListRepository interface {
	Create(*model.User, *model.List) error
	Update(list *model.List, updatingList model.List) error
	Destroy(*model.List) error
	DestroyLists(lists *[]model.List, tx *gorm.DB) error
	Move(list *model.List, toIndex int, currentUser *model.User) error
	Find(id int) (model.List, error)
	FindListsWithCards(*model.User) error
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
		err := tx.Select(clause.Associations).Delete(list).Error
		if err != nil {
			return err
		}

		err = tx.Model(model.List{}).Where("lists.index > ? AND lists.user_id = ?", list.Index, list.UserID).Updates(map[string]interface{}{"index": gorm.Expr("lists.index - ?", 1)}).Error
		if err != nil {
			return err
		}

		return nil
	})
}

// userを削除する際にリストとカードを一括削除するのに使う
func (r *listRepository) DestroyLists(lists *[]model.List, tx *gorm.DB) error {
	return tx.Select(clause.Associations).Delete(lists).Error
}

func (r *listRepository) Move(list *model.List, toIndex int, currentUser *model.User) error {
	if toIndex > list.Index {
		return r.db.Transaction(func(tx *gorm.DB) error {
			// 現在のリストの位置と移動後のリストの位置の間にあるリスト群を取り出し、indexを−1する
			err := tx.Model(&model.List{}).Where("lists.index <= ? AND lists.index > ? AND lists.user_id = ?", toIndex, list.Index, currentUser.ID).Updates(map[string]interface{}{"index": gorm.Expr("lists.index - ?", 1)}).Error
			if err != nil {
				return err
			}

			err = tx.Model(list).Update("index", toIndex).Error
			if err != nil {
				return err
			}

			return nil
		})
	}

	// toIndexがlistのindexよりも小さい時
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 現在のリストの位置と移動後のリストの位置の間にあるリスト群を取り出し、indexを+1する
		err := tx.Model(model.List{}).Where("lists.index >= ? AND lists.index < ? AND lists.user_id = ?", toIndex, list.ID, list.UserID).Updates(map[string]interface{}{"index": gorm.Expr("lists.index + ?", 1)}).Error
		if err != nil {
			return err
		}

		err = tx.Model(list).Update("index", toIndex).Error
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

func (r *listRepository) FindListsWithCards(user *model.User) error {
	// user.listsにlistsをsetする(cardもpreloadした状態で)
	return r.db.Where(model.List{UserID: user.ID}).Order("lists.index ASC").Preload("Cards", func(tx *gorm.DB) *gorm.DB {
		return tx.Order("cards.index ASC")
	}).Find(&user.Lists).Error
}
