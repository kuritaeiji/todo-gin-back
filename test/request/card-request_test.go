package request_test

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/controller"
	"github.com/kuritaeiji/todo-gin-back/db"
	"github.com/kuritaeiji/todo-gin-back/factory"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/server"
	"github.com/kuritaeiji/todo-gin-back/validators"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type CardRequestTestSuite struct {
	suite.Suite
	rec    *httptest.ResponseRecorder
	router *gin.Engine
	db     *gorm.DB
}

func (suite *CardRequestTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	config.Init()
	validators.Init()
	db.Init()
	suite.router = server.RouterSetup(controller.NewUserController())
	suite.db = db.GetDB()
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
	suite.Equal(card.ListID, rCard.ListID)
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
