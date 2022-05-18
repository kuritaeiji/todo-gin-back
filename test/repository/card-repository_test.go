package repository_test

import (
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/db"
	"github.com/kuritaeiji/todo-gin-back/factory"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/repository"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type CardRepositoryTestSuite struct {
	suite.Suite
	repository     repository.CardRepository
	listRepository repository.ListRepository
}

func (suite *CardRepositoryTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	config.Init()
	db.Init()
	suite.repository = repository.NewCardRepository()
	suite.listRepository = repository.NewListRepository()
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

func (suite *CardRepositoryTestSuite) TestSuccessMoveWhenIncreaseIndex() {
	user := factory.CreateUser(&factory.UserConfig{})
	list := factory.CreateList(&factory.ListConfig{}, user)
	cards := make([]model.Card, 0, 5)
	for i := 0; i <= 4; i++ {
		cards = append(cards, factory.CreateCard(&factory.CardConfig{Index: i, Title: strconv.Itoa(i)}, list))
	}
	err := suite.repository.Move(&cards[1], cards[1].ListID, 2)

	suite.Nil(err)
	suite.listRepository.FindListsWithCards(&user)
	dbCards := user.Lists[0].Cards
	suite.Equal("0", dbCards[0].Title)
	suite.Equal("2", dbCards[1].Title)
	suite.Equal("1", dbCards[2].Title)
	suite.Equal("3", dbCards[3].Title)
	suite.Equal("4", dbCards[4].Title)
}

func (suite *CardRepositoryTestSuite) TestSuccessMoveWhenDecreaseIndex() {
	user := factory.CreateUser(&factory.UserConfig{})
	list := factory.CreateList(&factory.ListConfig{}, user)
	cards := make([]model.Card, 0, 5)
	for i := 0; i <= 4; i++ {
		cards = append(cards, factory.CreateCard(&factory.CardConfig{Index: i, Title: strconv.Itoa(i)}, list))
	}
	err := suite.repository.Move(&cards[2], list.ID, 1)

	suite.Nil(err)
	suite.listRepository.FindListsWithCards(&user)
	dbCards := user.Lists[0].Cards
	suite.Equal("0", dbCards[0].Title)
	suite.Equal("2", dbCards[1].Title)
	suite.Equal("1", dbCards[2].Title)
	suite.Equal("3", dbCards[3].Title)
	suite.Equal("4", dbCards[4].Title)
}

func (suite *CardRepositoryTestSuite) TestSuccessMoveWhenChangeList() {
	user := factory.CreateUser(&factory.UserConfig{})
	list := factory.CreateList(&factory.ListConfig{Index: 0}, user)
	toList := factory.CreateList(&factory.ListConfig{Index: 1}, user)
	cards := make([]model.Card, 0, 3)
	toListCards := make([]model.Card, 0, 4)
	for i := 0; i <= 2; i++ {
		iString := strconv.Itoa(i)
		cards = append(cards, factory.CreateCard(&factory.CardConfig{Index: i, Title: "card" + iString}, list))
		toListCards = append(toListCards, factory.CreateCard(&factory.CardConfig{Index: i, Title: "toListCard" + iString}, toList))
	}
	err := suite.repository.Move(&cards[1], toList.ID, 2)

	suite.Nil(err)
	suite.listRepository.FindListsWithCards(&user)
	cards = user.Lists[0].Cards
	toListCards = user.Lists[1].Cards
	suite.Equal("card0", cards[0].Title)
	suite.Equal("card2", cards[1].Title)
	suite.Equal("toListCard0", toListCards[0].Title)
	suite.Equal("toListCard1", toListCards[1].Title)
	suite.Equal("card1", toListCards[2].Title)
	suite.Equal("toListCard2", toListCards[3].Title)
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
