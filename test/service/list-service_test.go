package service_test

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"strconv"
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
	"gorm.io/gorm"
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

func (suite *ListServiceTestSuite) TestSuccessIndex() {
	user := factory.NewUser(&factory.UserConfig{})
	suite.ctx.Set(config.CurrentUserKey, user)
	suite.listRepositoryMock.EXPECT().FindLists(&user).Return(nil)

	lists, err := suite.service.Index(suite.ctx)
	suite.Nil(err)
	suite.Equal(user.Lists, lists)
}

func (suite *ListServiceTestSuite) TestBadIndexWithDBError() {
	user := factory.NewUser(&factory.UserConfig{})
	suite.ctx.Set(config.CurrentUserKey, user)
	err := errors.New("db error")
	suite.listRepositoryMock.EXPECT().FindLists(&user).Return(err)
	lists, rerr := suite.service.Index(suite.ctx)

	suite.Equal(err, rerr)
	suite.Equal(user.Lists, lists)
}

func (suite *ListServiceTestSuite) TestSuccessCreate() {
	currentUser := model.User{}
	suite.ctx.Set(config.CurrentUserKey, currentUser)
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
	suite.ctx.Set(config.CurrentUserKey, currentUser)
	err := errors.New("DB error")
	list := factory.NewList(listConfig)
	suite.listRepositoryMock.EXPECT().Create(&currentUser, &list).Return(err)
	_, rerr := suite.service.Create(suite.ctx)

	suite.Equal(err, rerr)
}

func (suite *ListServiceTestSuite) TestSuccessUpdate() {
	user := factory.NewUser(&factory.UserConfig{})
	list := factory.NewList(&factory.ListConfig{})
	list.UserID = user.ID
	updatingListConfig := &factory.ListConfig{Title: "test title"}
	body := factory.CreateListRequestBody(updatingListConfig)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/lists/%v", list.ID), body)
	suite.ctx.Request = req
	suite.ctx.Set(config.CurrentUserKey, user)
	suite.ctx.Params = gin.Params{gin.Param{Key: "id", Value: strconv.Itoa(list.ID)}}
	suite.listRepositoryMock.EXPECT().Find(list.ID).Return(list, nil)
	suite.listRepositoryMock.EXPECT().Update(&list, gomock.Any()).Do(func(l *model.List, updatingList model.List) {
		suite.Equal(updatingListConfig.Title, updatingList.Title)
	})
	rList, err := suite.service.Update(suite.ctx)

	suite.Nil(err)
	suite.Equal(list.ID, rList.ID)
}

func (suite *ListServiceTestSuite) TestBadUpdateWithIDToIntError() {
	_, err := suite.service.Update(suite.ctx)

	suite.Error(err)
	suite.IsType(&strconv.NumError{}, err)
}

func (suite *ListServiceTestSuite) TestBadUpdateWithCannotFindList() {
	suite.ctx.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	err := gorm.ErrRecordNotFound
	suite.listRepositoryMock.EXPECT().Find(1).Return(model.List{}, err)
	_, rerr := suite.service.Update(suite.ctx)

	suite.Equal(err, rerr)
}

func (suite *ListServiceTestSuite) TestBadUpdateWithCurrentUserIsForbidden() {
	user := factory.NewUser(&factory.UserConfig{})
	list := factory.NewList(&factory.ListConfig{})
	list.UserID = user.ID + 1
	suite.ctx.Params = gin.Params{gin.Param{Key: "id", Value: strconv.Itoa(list.ID)}}
	suite.listRepositoryMock.EXPECT().Find(list.ID).Return(list, nil)
	suite.ctx.Set(config.CurrentUserKey, user)
	_, err := suite.service.Update(suite.ctx)

	suite.Equal(config.ForbiddenError, err)
}

func (suite *ListServiceTestSuite) TestBadUpdateWithValidationError() {
	var user model.User
	var list model.List
	suite.ctx.Params = gin.Params{gin.Param{Key: "id", Value: strconv.Itoa(list.ID)}}
	suite.listRepositoryMock.EXPECT().Find(list.ID).Return(list, nil)
	suite.ctx.Set(config.CurrentUserKey, user)
	req := httptest.NewRequest("PUT", "/api/lists/:id", factory.CreateListRequestBody(&factory.ListConfig{NotUseDefaultValue: true}))
	suite.ctx.Request = req
	_, err := suite.service.Update(suite.ctx)

	suite.IsType(validator.ValidationErrors{}, err)
}

func (suite *ListServiceTestSuite) TestBadUpdateWithDBError() {
	var user model.User
	var list model.List
	suite.ctx.Params = gin.Params{gin.Param{Key: "id", Value: strconv.Itoa(list.ID)}}
	suite.listRepositoryMock.EXPECT().Find(list.ID).Return(list, nil)
	suite.ctx.Set(config.CurrentUserKey, user)
	req := httptest.NewRequest("PUT", "/api/lists/:id", factory.CreateListRequestBody(&factory.ListConfig{}))
	suite.ctx.Request = req
	err := errors.New("db error")
	suite.listRepositoryMock.EXPECT().Update(&list, gomock.Any()).Return(err)
	_, rerr := suite.service.Update(suite.ctx)

	suite.Equal(err, rerr)
}
