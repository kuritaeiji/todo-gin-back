package service_test

import (
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

type ListMiddlewareServiceTestSuite struct {
	suite.Suite
	service            service.ListMiddlewareServive
	listRepositoryMock *mock_repository.MockListRepository
	ctx                *gin.Context
}

func (suite *ListMiddlewareServiceTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func (suite *ListMiddlewareServiceTestSuite) SetupTest() {
	suite.listRepositoryMock = mock_repository.NewMockListRepository(gomock.NewController(suite.T()))
	suite.service = service.TestNewListMiddlewareService(suite.listRepositoryMock)
	suite.ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
}

func TestListMiddlewareService(t *testing.T) {
	suite.Run(t, new(ListMiddlewareServiceTestSuite))
}

func (suite *ListMiddlewareServiceTestSuite) TestSuccessAuthorizeWithListIDParam() {
	listID := 1
	suite.ctx.Params = gin.Params{gin.Param{Key: "listID", Value: strconv.Itoa(listID)}}
	list := factory.NewList(&factory.ListConfig{})
	suite.listRepositoryMock.EXPECT().Find(listID).Return(list, nil)
	user := factory.NewUser(&factory.UserConfig{})
	user.ID = list.ID
	suite.ctx.Set(config.CurrentUserKey, user)
	rlist, err := suite.service.Authorize(suite.ctx)

	suite.Equal(list, rlist)
	suite.Nil(err)
}

func (suite *ListMiddlewareServiceTestSuite) TestSuccessAuthorizeWithIDParam() {
	listID := 1
	suite.ctx.Params = gin.Params{gin.Param{Key: "id", Value: strconv.Itoa(listID)}}
	list := factory.NewList(&factory.ListConfig{})
	suite.listRepositoryMock.EXPECT().Find(listID).Return(list, nil)
	user := factory.NewUser(&factory.UserConfig{})
	user.ID = list.ID
	suite.ctx.Set(config.CurrentUserKey, user)
	rlist, err := suite.service.Authorize(suite.ctx)

	suite.Equal(list, rlist)
	suite.Nil(err)
}

func (suite *ListMiddlewareServiceTestSuite) TestBadAuthorizeWithIDToIntError() {
	suite.ctx.Params = gin.Params{gin.Param{Key: "id", Value: "a"}}
	_, err := suite.service.Authorize(suite.ctx)

	suite.IsType(&strconv.NumError{}, err)
}

func (suite *ListMiddlewareServiceTestSuite) TestBadAuthorizeWithNotFoundList() {
	listID := 1
	suite.ctx.Params = gin.Params{gin.Param{Key: "id", Value: strconv.Itoa(listID)}}
	suite.listRepositoryMock.EXPECT().Find(listID).Return(model.List{}, gorm.ErrRecordNotFound)
	_, err := suite.service.Authorize(suite.ctx)

	suite.Equal(gorm.ErrRecordNotFound, err)
}

func (suite *ListMiddlewareServiceTestSuite) TestBadAuthorizeWithForbiddenError() {
	listID := 1
	suite.ctx.Params = gin.Params{gin.Param{Key: "listID", Value: strconv.Itoa(listID)}}
	list := factory.NewList(&factory.ListConfig{})
	suite.listRepositoryMock.EXPECT().Find(listID).Return(list, nil)
	user := factory.NewUser(&factory.UserConfig{})
	user.ID = list.ID + 1
	suite.ctx.Set(config.CurrentUserKey, user)
	_, err := suite.service.Authorize(suite.ctx)

	suite.Equal(config.ForbiddenError, err)
}
