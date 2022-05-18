package repository

// mockgen -source=repository/card-repository.go -destination=./mock_repository/card-repository.go

import (
	"github.com/kuritaeiji/todo-gin-back/db"
	"github.com/kuritaeiji/todo-gin-back/model"
	"gorm.io/gorm"
)

type cardRepository struct {
	db             *gorm.DB
	listRepository ListRepository
}

type CardRepository interface {
	Create(*model.Card, *model.List) error
	Update(card *model.Card, updatingCard *model.Card) error
	Destroy(card *model.Card) error
	Move(card *model.Card, toListID int, toIndex int) error
	Find(id int) (model.Card, error)
}

func NewCardRepository() CardRepository {
	return &cardRepository{db: db.GetDB(), listRepository: NewListRepository()}
}

var (
	plusIndexExpr  = map[string]interface{}{"index": gorm.Expr("cards.index + ?", 1)}
	minusIndexExpr = map[string]interface{}{"index": gorm.Expr("cards.index - ?", 1)}
)

func (r *cardRepository) Create(card *model.Card, list *model.List) error {
	return r.db.Model(list).Association("Cards").Append(card)
}

func (r *cardRepository) Update(card *model.Card, updatingCard *model.Card) error {
	return r.db.Model(&card).Select("title").Updates(updatingCard).Error
}

func (r *cardRepository) Destroy(card *model.Card) error {
	return r.db.Delete(&card).Error
}

func (r *cardRepository) Move(card *model.Card, toListID int, toIndex int) error {
	if card.ListID == toListID {
		if toIndex > card.Index {
			return r.moveWhenIncreaseIndex(card, toIndex)
		}

		return r.moveWhenDecreaseIndex(card, toIndex)
	}

	return r.moveWhenChangeList(card, toListID, toIndex)
}

func (r *cardRepository) moveWhenIncreaseIndex(card *model.Card, toIndex int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(model.Card{}).Where("cards.index > ? AND cards.index <= ? AND cards.list_id = ?", card.Index, toIndex, card.ListID).Updates(minusIndexExpr).Error
		if err != nil {
			return err
		}

		return tx.Model(card).Update("index", toIndex).Error
	})
}

func (r *cardRepository) moveWhenDecreaseIndex(card *model.Card, toIndex int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(model.Card{}).Where("cards.index < ? AND cards.index >= ? AND cards.list_id = ?", card.Index, toIndex, card.ListID).Updates(plusIndexExpr).Error
		if err != nil {
			return err
		}

		return tx.Model(card).Update("index", toIndex).Error
	})
}

func (r *cardRepository) moveWhenChangeList(card *model.Card, toListID int, toIndex int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(model.Card{}).Where("cards.index > ? AND cards.list_id = ?", card.Index, card.ListID).Updates(minusIndexExpr).Error
		if err != nil {
			return err
		}

		err = tx.Model(model.Card{}).Where("cards.index >= ? AND cards.list_id = ?", toIndex, toListID).Updates(plusIndexExpr).Error
		if err != nil {
			return err
		}

		err = tx.Model(card).Select("Index", "ListID").Updates(model.Card{Index: toIndex, ListID: toListID}).Error
		return err
	})
}

func (r *cardRepository) Find(id int) (model.Card, error) {
	var card model.Card
	err := r.db.Model(model.Card{}).First(&card, id).Error

	return card, err
}
