package model_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/factory"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/stretchr/testify/suite"
)

type CardModelTestSuite struct {
	suite.Suite
}

func (suite *CardModelTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func TestCardModel(t *testing.T) {
	suite.Run(t, new(CardModelTestSuite))
}

func (suite *CardModelTestSuite) TestToJson() {
	card := factory.NewCard(&factory.CardConfig{})
	cardJson := card.ToJson()

	suite.Equal(gin.H{"id": card.ID, "title": card.Title}, cardJson)
}

func (suite *CardModelTestSuite) TestToJsonCardSlice() {
	cards := make([]model.Card, 0, 2)
	for i := 1; i <= 2; i++ {
		cards = append(cards, factory.NewCard(&factory.CardConfig{}))
	}
	cardsJson := model.ToJsonCardSlice(cards)

	suite.Equal([]gin.H{{"id": cards[0].ID, "title": cards[0].Title}, {"id": cards[1].ID, "title": cards[1].Title}}, cardsJson)
}
