package service_test

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/dto"
	"github.com/kuritaeiji/todo-gin-back/factory"
	"github.com/kuritaeiji/todo-gin-back/mock_repository"
	"github.com/kuritaeiji/todo-gin-back/mock_service"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/service"
	"github.com/kuritaeiji/todo-gin-back/validators"
	"github.com/stretchr/testify/suite"
)

type CardServiceTestSuite struct {
	suite.Suite
	service                   service.CardService
	cardRepositoryMock        *mock_repository.MockCardRepository
	listMiddlewareServiceMock *mock_service.MockListMiddlewareServive
	ctx                       *gin.Context
}

func (suite *CardServiceTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	validators.Init()
}

func (suite *CardServiceTestSuite) SetupTest() {
	suite.cardRepositoryMock = mock_repository.NewMockCardRepository(gomock.NewController(suite.T()))
	suite.listMiddlewareServiceMock = mock_service.NewMockListMiddlewareServive(gomock.NewController(suite.T()))
	suite.service = service.TestNewCardService(suite.cardRepositoryMock, suite.listMiddlewareServiceMock)
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

func (suite *CardServiceTestSuite) TestSuccessUpdate() {
	updatingCardConfig := &factory.CardConfig{Title: "updated title"}
	req := httptest.NewRequest("PUT", "/api/cards/1", factory.CreateCardRequestBody(updatingCardConfig))
	suite.ctx.Request = req
	card := factory.NewCard(&factory.CardConfig{})
	suite.ctx.Set(config.CardKey, card)
	suite.cardRepositoryMock.EXPECT().Update(&card, gomock.Any()).Return(nil).Do(func(card *model.Card, updatingCard *model.Card) {
		suite.Equal(updatingCardConfig.Title, updatingCard.Title)
	})
	rCard, err := suite.service.Update(suite.ctx)

	suite.Equal(card, rCard)
	suite.Nil(err)
}

func (suite *CardServiceTestSuite) TestBadUpdateWithValidationError() {
	req := httptest.NewRequest("PUT", "/api/cards/1", factory.CreateCardRequestBody(&factory.CardConfig{NotUseDefaultValue: true}))
	suite.ctx.Request = req
	_, err := suite.service.Update(suite.ctx)

	suite.IsType(validator.ValidationErrors{}, err)
}

func (suite *CardServiceTestSuite) TestBadUpdateWithDBError() {
	req := httptest.NewRequest("PUT", "/api/cards/1", factory.CreateCardRequestBody(&factory.CardConfig{}))
	suite.ctx.Request = req
	card := factory.NewCard(&factory.CardConfig{})
	suite.ctx.Set(config.CardKey, card)
	err := errors.New("db error")
	suite.cardRepositoryMock.EXPECT().Update(&card, gomock.Any()).Return(err)
	_, rerr := suite.service.Update(suite.ctx)

	suite.Equal(err, rerr)
}

func (suite *CardServiceTestSuite) TestSuccessDestroyCard() {
	var card model.Card
	suite.ctx.Set(config.CardKey, card)
	suite.cardRepositoryMock.EXPECT().Destroy(&card).Return(nil)
	err := suite.service.Destroy(suite.ctx)

	suite.Nil(err)
}

func (suite *CardServiceTestSuite) TestBadDestroyCardWithDBError() {
	var card model.Card
	suite.ctx.Set(config.CardKey, card)
	err := errors.New("db error")
	suite.cardRepositoryMock.EXPECT().Destroy(&card).Return(err)
	rerr := suite.service.Destroy(suite.ctx)

	suite.Equal(err, rerr)
}

func (suite *CardServiceTestSuite) TestSuccessMoveCard() {
	var dtoMoveCard dto.MoveCard
	req := httptest.NewRequest("PUT", "/api/cards/1/move", factory.CreateMoveCardRequestBody(&dtoMoveCard))
	suite.ctx.Request = req
	var currentUser model.User
	suite.ctx.Set(config.CurrentUserKey, currentUser)
	suite.listMiddlewareServiceMock.EXPECT().FindAndAuthorizeList(dtoMoveCard.ToListID, currentUser).Return(model.List{}, nil)
	var card model.Card
	suite.ctx.Set(config.CardKey, card)
	suite.cardRepositoryMock.EXPECT().Move(&card, dtoMoveCard.ToListID, dtoMoveCard.ToIndex).Return(nil)
	err := suite.service.Move(suite.ctx)

	suite.Nil(err)
}

func (suite *CardServiceTestSuite) TestBadMoveCardWithValidationError() {
	dtoMoveCard := dto.MoveCard{ToIndex: -1}
	req := httptest.NewRequest("PUT", "/api/cards/1/move", factory.CreateMoveCardRequestBody(&dtoMoveCard))
	suite.ctx.Request = req
	err := suite.service.Move(suite.ctx)

	suite.IsType(validator.ValidationErrors{}, err)
}

func (suite *CardServiceTestSuite) TestBadMoveCardWithToListNotAuthorized() {
	dtoMoveCard := &dto.MoveCard{}
	req := httptest.NewRequest("PUT", "/api/cards/1/move", factory.CreateMoveCardRequestBody(dtoMoveCard))
	suite.ctx.Request = req
	var user model.User
	suite.ctx.Set(config.CurrentUserKey, user)
	suite.listMiddlewareServiceMock.EXPECT().FindAndAuthorizeList(dtoMoveCard.ToListID, user).Return(model.List{}, config.ForbiddenError)
	err := suite.service.Move(suite.ctx)

	suite.Equal(config.ForbiddenError, err)
}

func (suite *CardServiceTestSuite) TestBadMoveCardWithDBError() {
	var dtoMoveCard dto.MoveCard
	req := httptest.NewRequest("PUT", "/api/cards/1/move", factory.CreateMoveCardRequestBody(&dtoMoveCard))
	suite.ctx.Request = req
	var currentUser model.User
	suite.ctx.Set(config.CurrentUserKey, currentUser)
	suite.listMiddlewareServiceMock.EXPECT().FindAndAuthorizeList(dtoMoveCard.ToListID, currentUser).Return(model.List{}, nil)
	var card model.Card
	suite.ctx.Set(config.CardKey, card)
	err := errors.New("db error")
	suite.cardRepositoryMock.EXPECT().Move(&card, dtoMoveCard.ToListID, dtoMoveCard.ToIndex).Return(err)
	rerr := suite.service.Move(suite.ctx)

	suite.Equal(err, rerr)
}
