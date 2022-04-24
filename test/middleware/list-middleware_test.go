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

type ListMiddlewareTestSuite struct {
	suite.Suite
	middleware                middleware.ListMiddleware
	listMiddlewareServiceMock *mock_service.MockListMiddlewareServive
	rec                       *httptest.ResponseRecorder
	ctx                       *gin.Context
}

func (suite *ListMiddlewareTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func (suite *ListMiddlewareTestSuite) SetupTest() {
	suite.listMiddlewareServiceMock = mock_service.NewMockListMiddlewareServive(gomock.NewController(suite.T()))
	suite.middleware = middleware.TestNewListMiddleware(suite.listMiddlewareServiceMock)
	suite.rec = httptest.NewRecorder()
	suite.ctx, _ = gin.CreateTestContext(suite.rec)
}

func TestListMiddleware(t *testing.T) {
	suite.Run(t, new(ListMiddlewareTestSuite))
}

func (suite *ListMiddlewareTestSuite) TestSuccessAuthorize() {
	list := factory.NewList(&factory.ListConfig{})
	suite.listMiddlewareServiceMock.EXPECT().Authorize(suite.ctx).Return(list, nil)
	suite.middleware.Authorize(suite.ctx)

	rList := suite.ctx.MustGet(config.ListKey).(model.List)
	suite.Equal(list, rList)
}

func (suite *ListMiddlewareTestSuite) TestBadAuthorize() {
	suite.listMiddlewareServiceMock.EXPECT().Authorize(suite.ctx).Return(model.List{}, gorm.ErrRecordNotFound)
	suite.middleware.Authorize(suite.ctx)

	suite.Equal(config.RecordNotFoundErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.RecordNotFoundErrorResponse.Json["content"])
}

func (suite *ListMiddlewareTestSuite) TestBadAuthorizeWithForbiddenError() {
	suite.listMiddlewareServiceMock.EXPECT().Authorize(suite.ctx).Return(model.List{}, config.ForbiddenError)
	suite.middleware.Authorize(suite.ctx)

	suite.Equal(config.ForbiddenErrorResponse.Code, suite.rec.Code)
	suite.Contains(suite.rec.Body.String(), config.ForbiddenErrorResponse.Json["content"])
}
