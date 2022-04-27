package request_test

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/controller"
	"github.com/kuritaeiji/todo-gin-back/db"
	"github.com/kuritaeiji/todo-gin-back/dto"
	"github.com/kuritaeiji/todo-gin-back/factory"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/repository"
	"github.com/kuritaeiji/todo-gin-back/server"
	"github.com/kuritaeiji/todo-gin-back/validators"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type CardRequestTestSuite struct {
	suite.Suite
	rec            *httptest.ResponseRecorder
	router         *gin.Engine
	db             *gorm.DB
	repository     repository.CardRepository
	listRepository repository.ListRepository
}

func (suite *CardRequestTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	config.Init()
	validators.Init()
	db.Init()
	suite.router = server.RouterSetup(controller.NewUserController())
	suite.db = db.GetDB()
	suite.repository = repository.NewCardRepository()
	suite.listRepository = repository.NewListRepository()
}

func (suite *CardRequestTestSuite) SetupTest() {
	suite.rec = httptest.NewRecorder()
}

func (suite *CardRequestTestSuite) TearDownSuite() {
	db.CloseDB()
}

func (suite *CardRequestTestSuite) TearDownTest() {
	db.DeleteAll()
}

func TestCardRequest(t *testing.T) {
	suite.Run(t, new(CardRequestTestSuite))
}

func (suite *CardRequestTestSuite) TestSuccessCreate() {
	user := factory.CreateUser(&factory.UserConfig{})
	list := factory.CreateList(&factory.ListConfig{}, user)
	cardConfig := &factory.CardConfig{}
	req := httptest.NewRequest("POST", fmt.Sprintf("/api/lists/%v/cards", list.ID), factory.CreateCardRequestBody(cardConfig))
	req.Header.Add(config.TokenHeader, factory.CreateAccessToken(user))
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(200, suite.rec.Code)
	var card model.Card
	suite.db.Model(model.Card{}).First(&card)
	suite.Equal(cardConfig.Title, card.Title)
	suite.Equal(cardConfig.Index, card.Index)
	suite.Equal(list.ID, card.ListID)

	var rCard model.Card
	json.Unmarshal(suite.rec.Body.Bytes(), &rCard)
	suite.Equal(card.ID, rCard.ID)
	suite.Equal(card.Title, rCard.Title)
}

func (suite *ListRequestTestSuite) TestBadCreateWithNotFoundList() {
	user := factory.CreateUser(&factory.UserConfig{})
	req := httptest.NewRequest("POST", "/api/lists/1/cards", nil)
	req.Header.Add(config.TokenHeader, factory.CreateAccessToken(user))
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(config.RecordNotFoundErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.RecordNotFoundErrorResponse.Json["content"])
}

func (suite *ListRequestTestSuite) TestBadCreateWithListForbiddenError() {
	user := factory.CreateUser(&factory.UserConfig{})
	otherUser := factory.CreateUser(&factory.UserConfig{})
	list := factory.CreateList(&factory.ListConfig{}, otherUser)
	req := httptest.NewRequest("POST", fmt.Sprintf("/api/lists/%v/cards", list.ID), nil)
	req.Header.Add(config.TokenHeader, factory.CreateAccessToken(user))
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(config.ForbiddenErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.ForbiddenErrorResponse.Json["content"])
}

func (suite *ListRequestTestSuite) TestBadCardCreateWithValidationError() {
	user := factory.CreateUser(&factory.UserConfig{})
	list := factory.CreateList(&factory.ListConfig{}, user)
	req := httptest.NewRequest("POST", fmt.Sprintf("/api/lists/%v/cards", list.ID), factory.CreateCardRequestBody(&factory.CardConfig{NotUseDefaultValue: true}))
	req.Header.Add(config.TokenHeader, factory.CreateAccessToken(user))
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(config.ValidationErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.ValidationErrorResponse.Json["content"])
}

func (suite *ListRequestTestSuite) TestSuccessUpdateCard() {
	user := factory.CreateUser(&factory.UserConfig{})
	list := factory.CreateList(&factory.ListConfig{}, user)
	card := factory.CreateCard(&factory.CardConfig{}, list)
	updatingCardConfig := &factory.CardConfig{Title: "updated title"}
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/cards/%v", card.ID), factory.CreateCardRequestBody(updatingCardConfig))
	req.Header.Add(config.TokenHeader, factory.CreateAccessToken(user))
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(200, suite.rec.Code)
	var rCard model.Card
	json.Unmarshal(suite.rec.Body.Bytes(), &rCard)
	suite.Equal(card.ID, rCard.ID)
	suite.Equal(updatingCardConfig.Title, rCard.Title)
	dbCard, _ := suite.repository.Find(card.ID)
	suite.Equal(updatingCardConfig.Title, dbCard.Title)
}

func (suite *ListRequestTestSuite) TestBadUpdateCardWithNotFoundCard() {
	user := factory.CreateUser(&factory.UserConfig{})
	req := httptest.NewRequest("PUT", "/api/cards/1", nil)
	req.Header.Add(config.TokenHeader, factory.CreateAccessToken(user))
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(config.RecordNotFoundErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.RecordNotFoundErrorResponse.Json["content"])
}

func (suite *ListRequestTestSuite) TestBadUpdateCardWithForbiddenCardError() {
	user := factory.CreateUser(&factory.UserConfig{})
	otherUser := factory.CreateUser(&factory.UserConfig{})
	list := factory.CreateList(&factory.ListConfig{}, otherUser)
	card := factory.CreateCard(&factory.CardConfig{}, list)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/cards/%v", card.ID), nil)
	req.Header.Add(config.TokenHeader, factory.CreateAccessToken(user))
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(config.ForbiddenErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.ForbiddenErrorResponse.Json["content"])
}

func (suite *ListRequestTestSuite) TestBadUpdateCardWithValidationError() {
	user := factory.CreateUser(&factory.UserConfig{})
	list := factory.CreateList(&factory.ListConfig{}, user)
	card := factory.CreateCard(&factory.CardConfig{}, list)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/cards/%v", card.ID), factory.CreateCardRequestBody(&factory.CardConfig{NotUseDefaultValue: true}))
	req.Header.Add(config.TokenHeader, factory.CreateAccessToken(user))
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(config.ValidationErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.ValidationErrorResponse.Json["content"])
}

func (suite *CardRequestTestSuite) TestSuccessDestroyCard() {
	user := factory.CreateUser(&factory.UserConfig{})
	list := factory.CreateList(&factory.ListConfig{}, user)
	card := factory.CreateCard(&factory.CardConfig{}, list)
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/cards/%v", card.ID), nil)
	req.Header.Add(config.TokenHeader, factory.CreateAccessToken(user))
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(200, suite.rec.Code)
	_, err := suite.repository.Find(card.ID)
	suite.Equal(gorm.ErrRecordNotFound, err)
}

func (suite *CardRequestTestSuite) TestSuccessMoveCardWhenIncreaseIndex() {
	user := factory.CreateUser(&factory.UserConfig{})
	list := factory.CreateList(&factory.ListConfig{}, user)
	cards := make([]model.Card, 0, 5)
	for i := 0; i <= 4; i++ {
		cards = append(cards, factory.CreateCard(&factory.CardConfig{Title: strconv.Itoa(i), Index: i}, list))
	}
	dtoMoveCard := &dto.MoveCard{ToIndex: 3, ToListID: list.ID}
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/cards/%v/move", cards[1].ID), factory.CreateMoveCardRequestBody(dtoMoveCard))
	req.Header.Add(config.TokenHeader, factory.CreateAccessToken(user))
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(200, suite.rec.Code)
	suite.listRepository.FindListsWithCards(&user)
	cards = user.Lists[0].Cards
	suite.Equal("0", cards[0].Title)
	suite.Equal("2", cards[1].Title)
	suite.Equal("3", cards[2].Title)
	suite.Equal("1", cards[3].Title)
	suite.Equal("4", cards[4].Title)
}

func (suite *CardRequestTestSuite) TestSuccessMoveCardWhenDecreaseIndex() {
	user := factory.CreateUser(&factory.UserConfig{})
	list := factory.CreateList(&factory.ListConfig{}, user)
	cards := make([]model.Card, 0, 5)
	for i := 0; i <= 4; i++ {
		cards = append(cards, factory.CreateCard(&factory.CardConfig{Title: strconv.Itoa(i), Index: i}, list))
	}
	dtoMoveCard := &dto.MoveCard{ToIndex: 1, ToListID: list.ID}
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/cards/%v/move", cards[3].ID), factory.CreateMoveCardRequestBody(dtoMoveCard))
	req.Header.Add(config.TokenHeader, factory.CreateAccessToken(user))
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(200, suite.rec.Code)
	fmt.Println(suite.rec.Body.String())
	suite.listRepository.FindListsWithCards(&user)
	cards = user.Lists[0].Cards
	suite.Equal("0", cards[0].Title)
	suite.Equal("3", cards[1].Title)
	suite.Equal("1", cards[2].Title)
	suite.Equal("2", cards[3].Title)
	suite.Equal("4", cards[4].Title)
}

func (suite *CardRequestTestSuite) TestSuccessMoveCardWhenChangeList() {
	user := factory.CreateUser(&factory.UserConfig{})
	list1 := factory.CreateList(&factory.ListConfig{}, user)
	list2 := factory.CreateList(&factory.ListConfig{Index: 1}, user)
	cards1 := make([]model.Card, 0, 3)
	cards2 := make([]model.Card, 0, 3)
	for i := 0; i <= 2; i++ {
		iString := strconv.Itoa(i)
		cards1 = append(cards1, factory.CreateCard(&factory.CardConfig{Index: i, Title: iString + "card1"}, list1))
		cards2 = append(cards2, factory.CreateCard(&factory.CardConfig{Index: i, Title: iString + "card2"}, list2))
	}
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/cards/%v/move", cards1[1].ID), factory.CreateMoveCardRequestBody(&dto.MoveCard{ToIndex: 2, ToListID: list2.ID}))
	req.Header.Add(config.TokenHeader, factory.CreateAccessToken(user))
	suite.router.ServeHTTP(suite.rec, req)

	suite.Equal(200, suite.rec.Code)
	suite.listRepository.FindListsWithCards(&user)
	cards1 = user.Lists[0].Cards
	cards2 = user.Lists[1].Cards
	suite.Equal("0card1", cards1[0].Title)
	suite.Equal("2card1", cards1[1].Title)
	suite.Equal("0card2", cards2[0].Title)
	suite.Equal("1card2", cards2[1].Title)
	suite.Equal("1card1", cards2[2].Title)
	suite.Equal("2card2", cards2[3].Title)
}
