package repository_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/db"
	"github.com/kuritaeiji/todo-gin-back/factory"
	"github.com/kuritaeiji/todo-gin-back/repository"
	"github.com/stretchr/testify/suite"
)

type CardRepositoryTestSuite struct {
	suite.Suite
	repository repository.CardRepository
}

func (suite *CardRepositoryTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	config.Init()
	db.Init()
	suite.repository = repository.NewCardRepository()
}

func (suite *CardRepositoryTestSuite) TearDownSuite() {
	db.CloseDB()
}

func (suite *CardRepositoryTestSuite) TearDownTest() {
	db.DeleteAll()
}

func TestCardRepository(t *testing.T) {
	suite.Run(t, new(CardRepositoryTestSuite))
}

func (suite *CardRepositoryTestSuite) TestSuccessCreate() {
	user := factory.CreateUser(&factory.UserConfig{})
	list := factory.CreateList(&factory.ListConfig{}, user)
	card := factory.NewCard(&factory.CardConfig{})
	err := suite.repository.Create(&card, &list)

	suite.Nil(err)
	rCard, _ := suite.repository.Find(card.ID)
	suite.Equal(card.ID, rCard.ID)
	suite.Equal(card.Title, rCard.Title)
	suite.Equal(list.ID, rCard.ListID)
}

func (suite *CardRepositoryTestSuite) TestSuccessFind() {
	user := factory.CreateUser(&factory.UserConfig{})
	list := factory.CreateList(&factory.ListConfig{}, user)
	card := factory.CreateCard(&factory.CardConfig{}, list)
	rCard, err := suite.repository.Find(card.ID)

	suite.Nil(err)
	suite.Equal(card, rCard)
}