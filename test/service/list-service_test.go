package service_test

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/kuritaeiji/todo-gin-back/factory"
	"github.com/kuritaeiji/todo-gin-back/middleware"
	"github.com/kuritaeiji/todo-gin-back/mock_repository"
	"github.com/kuritaeiji/todo-gin-back/model"
	"github.com/kuritaeiji/todo-gin-back/service"
	"github.com/kuritaeiji/todo-gin-back/validators"
	"github.com/stretchr/testify/suite"
)

type ListServiceTestSuite struct {
	suite.Suite
	service            service.ListService
	listRepositoryMock *mock_repository.MockListRepository
	ctx                *gin.Context
}

func (suite *ListServiceTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	validators.Init()
}

func (suite *ListServiceTestSuite) SetupTest() {
	suite.listRepositoryMock = mock_repository.NewMockListRepository(gomock.NewController(suite.T()))
	suite.service = service.TestNewListService(suite.listRepositoryMock)
	suite.ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
}

func TestListServiceSuite(t *testing.T) {
	suite.Run(t, new(ListServiceTestSuite))
}

func (suite *ListServiceTestSuite) TestSuccessCreate() {
	currentUser := model.User{}
	suite.ctx.Set(middleware.CurrentUserKey, currentUser)
	listConfig := &factory.ListConfig{}
	req := httptest.NewRequest("POST", "/api/lists", factory.CreateListRequestBody(listConfig))
	suite.ctx.Request = req
	list := factory.NewList(listConfig)
	suite.listRepositoryMock.EXPECT().Create(&currentUser, &list).Return(nil)
	rList, err := suite.service.Create(suite.ctx)

	suite.Nil(err)
	suite.Equal(list, rList)
}

func (suite *ListServiceTestSuite) TestBadCreateWithValidationError() {
	req := httptest.NewRequest("POST", "/api/lists", factory.CreateListRequestBody(&factory.ListConfig{NotUseDefaultValue: true}))
	suite.ctx.Request = req
	list, err := suite.service.Create(suite.ctx)

	suite.Equal(model.List{}, list)
	suite.IsType(validator.ValidationErrors{}, err)
}

func (suite *ListServiceTestSuite) TestBadCreateWithDBError() {
	listConfig := &factory.ListConfig{}
	req := httptest.NewRequest("POST", "/api/lists", factory.CreateListRequestBody(listConfig))
	suite.ctx.Request = req
	var currentUser model.User
	suite.ctx.Set(middleware.CurrentUserKey, currentUser)
	err := errors.New("DB error")
	list := factory.NewList(listConfig)
	suite.listRepositoryMock.EXPECT().Create(&currentUser, &list).Return(err)
	_, rerr := suite.service.Create(suite.ctx)

	suite.Equal(err, rerr)
}
