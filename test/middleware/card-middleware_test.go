package middleware_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kuritaeiji/todo-gin-back/config"
	"github.com/kuritaeiji/todo-gin-back/factory"
	"github.com/kuritaeiji/todo-gin-back/middleware"
	"github.com/kuritaeiji/todo-gin-back/mock_service"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type CardMiddlewareTestSuite struct {
	suite.Suite
	middleware                middleware.CardMiddleware
	CardMiddlewareServiceMock *mock_service.MockCardMiddlewareService
	rec                       *httptest.ResponseRecorder
	ctx                       *gin.Context
}

func (suite *CardMiddlewareTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func (suite *CardMiddlewareTestSuite) SetupTest() {
	suite.CardMiddlewareServiceMock = mock_service.NewMockCardMiddlewareService(gomock.NewController(suite.T()))
	suite.middleware = middleware.TestNewCardMiddleware(suite.CardMiddlewareServiceMock)
	suite.rec = httptest.NewRecorder()
	suite.ctx, _ = gin.CreateTestContext(suite.rec)
}

func TestCardMiddleware(t *testing.T) {
	suite.Run(t, new(CardMiddlewareTestSuite))
}

func (suite *CardMiddlewareTestSuite) TestSuccessAuthorize() {
	card := factory.NewCard(&factory.CardConfig{})
	suite.CardMiddlewareServiceMock.EXPECT().Authorize(suite.ctx).Return(card, nil)
	suite.middleware.Authorize(suite.ctx)

	rCard := suite.ctx.MustGet(config.CardKey).(model.Card)
	suite.Equal(card, rCard)
}

func (suite *CardMiddlewareTestSuite) TestBadAuthorize() {
	suite.CardMiddlewareServiceMock.EXPECT().Authorize(suite.ctx).Return(model.Card{}, gorm.ErrRecordNotFound)
	suite.middleware.Authorize(suite.ctx)

	suite.Equal(config.RecordNotFoundErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.RecordNotFoundErrorResponse.Json["content"])
}

func (suite *CardMiddlewareTestSuite) TestBadAuthorizeWithForbiddenError() {
	suite.CardMiddlewareServiceMock.EXPECT().Authorize(suite.ctx).Return(model.Card{}, config.ForbiddenError)
	suite.middleware.Authorize(suite.ctx)

	suite.Equal(config.ForbiddenErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.ForbiddenErrorResponse.Json["content"])
}
