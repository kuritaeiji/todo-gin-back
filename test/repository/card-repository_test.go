package repository_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/db"
	"github.com/kuritaeiji/todo-gin-back/factory"
	"github.com/kuritaeiji/todo-gin-back/repository"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
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

func (suite *CardRepositoryTestSuite) TestSuccessUpdate() {
	user := factory.CreateUser(&factory.UserConfig{})
	list := factory.CreateList(&factory.ListConfig{}, user)
	card := factory.CreateCard(&factory.CardConfig{}, list)
	updatingCard := factory.NewCard(&factory.CardConfig{Title: "updated title"})
	err := suite.repository.Update(&card, &updatingCard)

	suite.Nil(err)
	rCard, _ := suite.repository.Find(card.ID)
	suite.Equal(updatingCard.Title, rCard.Title)
}

func (suite *CardRepositoryTestSuite) TestSuccessFind() {
	user := factory.CreateUser(&factory.UserConfig{})
	list := factory.CreateList(&factory.ListConfig{}, user)
	card := factory.CreateCard(&factory.CardConfig{}, list)
	rCard, err := suite.repository.Find(card.ID)

	suite.Nil(err)
	suite.Equal(card, rCard)
}

func (suite *CardRepositoryTestSuite) TestSuccessDestroy() {
	user := factory.CreateUser(&factory.UserConfig{})
	list := factory.CreateList(&factory.ListConfig{}, user)
	card := factory.CreateCard(&factory.CardConfig{}, list)
	err := suite.repository.Destroy(&card)

	suite.Nil(err)
	_, err = suite.repository.Find(card.ID)
	suite.Equal(gorm.ErrRecordNotFound, err)
}
