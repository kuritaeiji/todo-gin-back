package repository

// mockgen -source=repository/card-repository.go -destination=./mock_repository/card-repository.go

import (
	"github.com/kuritaeiji/todo-gin-back/db"
	"github.com/kuritaeiji/todo-gin-back/model"
	"gorm.io/gorm"
)

type cardRepository struct {
	db *gorm.DB
}

type CardRepository interface {
	Create(*model.Card, *model.List) error
	// Update(card *model.Card, updatingCard *model.List) error
	// Destroy(card *model.Card) error
	Find(id int) (model.Card, error)
}

func NewCardRepository() CardRepository {
	return &cardRepository{db: db.GetDB()}
}

func (r *cardRepository) Create(card *model.Card, list *model.List) error {
	return r.db.Model(list).Association("Cards").Append(card)
}

func (r *cardRepository) Find(id int) (model.Card, error) {
	var card model.Card
	err := r.db.Model(model.Card{}).First(&card, id).Error

	return card, err
}
