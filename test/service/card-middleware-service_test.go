package service_test

import (
	"errors"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/factory"
	"github.com/kuritaeiji/todo-gin-back/mock_repository"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/service"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type CardMiddlewareServiceTestSuite struct {
	suite.Suite
	service            service.CardMiddlewareService
	cardRepositoryMock *mock_repository.MockCardRepository
	userRepositoryMock *mock_repository.MockUserRepository
	ctx                *gin.Context
}

func (suite *CardMiddlewareServiceTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func (suite *CardMiddlewareServiceTestSuite) SetupTest() {
	suite.cardRepositoryMock = mock_repository.NewMockCardRepository(gomock.NewController(suite.T()))
	suite.userRepositoryMock = mock_repository.NewMockUserRepository(gomock.NewController(suite.T()))
	suite.service = service.TestNewCardMiddlewareService(suite.cardRepositoryMock, suite.userRepositoryMock)
	suite.ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
}

func TestCardMiddlewareService(t *testing.T) {
	suite.Run(t, new(CardMiddlewareServiceTestSuite))
}

func (suite *CardMiddlewareServiceTestSuite) TestSuccessAuthorize() {
	cardID := 1
	suite.ctx.Params = gin.Params{gin.Param{Key: "id", Value: strconv.Itoa(cardID)}}
	card := factory.NewCard(&factory.CardConfig{})
	suite.cardRepositoryMock.EXPECT().Find(cardID).Return(card, nil)
	currentUser := factory.NewUser(&factory.UserConfig{})
	suite.ctx.Set(config.CurrentUserKey, currentUser)
	suite.userRepositoryMock.EXPECT().HasCard(card, currentUser).Return(true, nil)
	rCard, err := suite.service.Authorize(suite.ctx)

	suite.Equal(card, rCard)
	suite.Nil(err)
}

func (suite *CardMiddlewareServiceTestSuite) TestBadAuthorizeWithIDToIntError() {
	suite.ctx.Params = gin.Params{gin.Param{Key: "id", Value: "a"}}
	_, err := suite.service.Authorize(suite.ctx)

	suite.IsType(&strconv.NumError{}, err)
}

func (suite *CardMiddlewareServiceTestSuite) TestBadAuthorizeWithNotFoundCard() {
	cardID := 1
	suite.ctx.Params = gin.Params{gin.Param{Key: "id", Value: strconv.Itoa(cardID)}}
	suite.cardRepositoryMock.EXPECT().Find(cardID).Return(model.Card{}, gorm.ErrRecordNotFound)
	_, err := suite.service.Authorize(suite.ctx)

	suite.Equal(gorm.ErrRecordNotFound, err)
}

func (suite *CardMiddlewareServiceTestSuite) TestBadAuthorizeWithForbiddenError() {
	cardID := 1
	suite.ctx.Params = gin.Params{gin.Param{Key: "id", Value: strconv.Itoa(cardID)}}
	card := factory.NewCard(&factory.CardConfig{})
	suite.cardRepositoryMock.EXPECT().Find(cardID).Return(card, nil)
	currentUser := factory.NewUser(&factory.UserConfig{})
	suite.ctx.Set(config.CurrentUserKey, currentUser)
	suite.userRepositoryMock.EXPECT().HasCard(card, currentUser).Return(false, nil)
	_, err := suite.service.Authorize(suite.ctx)

	suite.Equal(config.ForbiddenError, err)
}

func (suite *CardMiddlewareServiceTestSuite) TestBadAuthorizeWithDBError() {
	cardID := 1
	suite.ctx.Params = gin.Params{gin.Param{Key: "id", Value: strconv.Itoa(cardID)}}
	card := factory.NewCard(&factory.CardConfig{})
	suite.cardRepositoryMock.EXPECT().Find(cardID).Return(card, nil)
	currentUser := factory.NewUser(&factory.UserConfig{})
	suite.ctx.Set(config.CurrentUserKey, currentUser)
	err := errors.New("db error")
	suite.userRepositoryMock.EXPECT().HasCard(card, currentUser).Return(false, err)
	_, rerr := suite.service.Authorize(suite.ctx)

	suite.Equal(err, rerr)
}
