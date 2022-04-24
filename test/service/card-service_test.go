package service_test

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/factory"
	"github.com/kuritaeiji/todo-gin-back/mock_repository"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/service"
	"github.com/kuritaeiji/todo-gin-back/validators"
	"github.com/stretchr/testify/suite"
)

type CardServiceTestSuite struct {
	suite.Suite
	service            service.CardService
	cardRepositoryMock *mock_repository.MockCardRepository
	ctx                *gin.Context
}

func (suite *CardServiceTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	validators.Init()
}

func (suite *CardServiceTestSuite) SetupTest() {
	suite.cardRepositoryMock = mock_repository.NewMockCardRepository(gomock.NewController(suite.T()))
	suite.service = service.TestNewCardService(suite.cardRepositoryMock)
	suite.ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
}

func TestCardService(t *testing.T) {
	suite.Run(t, new(CardServiceTestSuite))
}

func (suite *CardServiceTestSuite) TestSuccessCreate() {
	cardFactory := &factory.CardConfig{}
	req := httptest.NewRequest("POST", "/api/lists/:listID/cards", factory.CreateCardRequestBody(cardFactory))
	suite.ctx.Request = req
	list := factory.NewList(&factory.ListConfig{})
	suite.ctx.Set(config.ListKey, list)
	suite.cardRepositoryMock.EXPECT().Create(gomock.Any(), &list).Return(nil).Do(func(argCard *model.Card, argList *model.List) {
		suite.Equal(cardFactory.Title, argCard.Title)
		suite.Equal(cardFactory.Index, argCard.Index)
	})
	rCard, err := suite.service.Create(suite.ctx)

	suite.Nil(err)
	suite.Equal(cardFactory.Title, rCard.Title)
	suite.Equal(cardFactory.Index, rCard.Index)
}

func (suite *CardServiceTestSuite) TestBadCreateWithValidation() {
	suite.ctx.Request = httptest.NewRequest("POST", "/api/lists/listID/cards", factory.CreateCardRequestBody(&factory.CardConfig{NotUseDefaultValue: true}))
	_, err := suite.service.Create(suite.ctx)

	suite.IsType(validator.ValidationErrors{}, err)
}

func (suite *CardServiceTestSuite) TestBadCreateWithDBError() {
	suite.ctx.Request = httptest.NewRequest("POST", "/api/lists/listID/cards", factory.CreateCardRequestBody(&factory.CardConfig{}))
	list := factory.NewList(&factory.ListConfig{})
	suite.ctx.Set(config.ListKey, list)
	err := errors.New("db error")
	suite.cardRepositoryMock.EXPECT().Create(gomock.Any(), &list).Return(err)
	_, rerr := suite.service.Create(suite.ctx)

	suite.Equal(err, rerr)
}