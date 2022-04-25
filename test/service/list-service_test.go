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
	suite.listRepositoryMock.EXPECT().FindListsWithCards(&user).Return(nil)

	lists, err := suite.service.Index(suite.ctx)
	suite.Nil(err)
	suite.Equal(user.Lists, lists)
}

func (suite *ListServiceTestSuite) TestBadIndexWithDBError() {
	user := factory.NewUser(&factory.UserConfig{})
	suite.ctx.Set(config.CurrentUserKey, user)
	err := errors.New("db error")
	suite.listRepositoryMock.EXPECT().FindListsWithCards(&user).Return(err)
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
	updatingListConfig := &factory.ListConfig{Title: "test title"}
	body := factory.CreateListRequestBody(updatingListConfig)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/api/lists/%v", list.ID), body)
	suite.ctx.Request = req
	suite.ctx.Set(config.CurrentUserKey, user)
	suite.ctx.Set(config.ListKey, list)
	suite.listRepositoryMock.EXPECT().Update(&list, gomock.Any()).Do(func(l *model.List, updatingList model.List) {
		suite.Equal(updatingListConfig.Title, updatingList.Title)
	})
	rList, err := suite.service.Update(suite.ctx)

	suite.Nil(err)
	suite.Equal(list.ID, rList.ID)
}

func (suite *ListServiceTestSuite) TestBadUpdateWithValidationError() {
	var user model.User
	var list model.List
	suite.ctx.Set(config.CurrentUserKey, user)
	suite.ctx.Set(config.ListKey, list)
	req := httptest.NewRequest("PUT", "/api/lists/:id", factory.CreateListRequestBody(&factory.ListConfig{NotUseDefaultValue: true}))
	suite.ctx.Request = req
	_, err := suite.service.Update(suite.ctx)

	suite.IsType(validator.ValidationErrors{}, err)
}

func (suite *ListServiceTestSuite) TestBadUpdateWithDBError() {
	var user model.User
	var list model.List
	suite.ctx.Set(config.ListKey, list)
	suite.ctx.Set(config.CurrentUserKey, user)
	req := httptest.NewRequest("PUT", "/api/lists/:id", factory.CreateListRequestBody(&factory.ListConfig{}))
	suite.ctx.Request = req
	err := errors.New("db error")
	suite.listRepositoryMock.EXPECT().Update(&list, gomock.Any()).Return(err)
	_, rerr := suite.service.Update(suite.ctx)

	suite.Equal(err, rerr)
}

func (suite *ListServiceTestSuite) TestSuccessDestroy() {
	id := 1
	suite.ctx.Params = gin.Params{gin.Param{Key: "id", Value: strconv.Itoa(id)}}
	list := factory.NewList(&factory.ListConfig{})
	suite.ctx.Set(config.ListKey, list)
	currentUser := factory.NewUser(&factory.UserConfig{})
	suite.ctx.Set(config.CurrentUserKey, currentUser)
	list.UserID = currentUser.ID
	suite.listRepositoryMock.EXPECT().Destroy(&list).Return(nil)
	err := suite.service.Destroy(suite.ctx)

	suite.Nil(err)
}

func (suite *ListServiceTestSuite) TestBadDestroyWithDBError() {
	currentUser := factory.NewUser(&factory.UserConfig{})
	list := factory.NewList(&factory.ListConfig{})
	suite.ctx.Set(config.ListKey, list)
	suite.ctx.Set(config.CurrentUserKey, currentUser)
	err := errors.New("db error")
	suite.listRepositoryMock.EXPECT().Destroy(&list).Return(err)
	rerr := suite.service.Destroy(suite.ctx)

	suite.Equal(err, rerr)
}

func (suite *ListServiceTestSuite) TestSuccessMove() {
	list := factory.NewList(&factory.ListConfig{})
	user := factory.NewUser(&factory.UserConfig{})
	suite.ctx.Set(config.CurrentUserKey, user)
	suite.ctx.Set(config.ListKey, list)
	toIndex := 1
	req := httptest.NewRequest("PUT", "/api/lists/1/move", factory.CreateListRequestBody(&factory.ListConfig{Index: toIndex}))
	suite.ctx.Request = req
	suite.listRepositoryMock.EXPECT().Move(&list, toIndex, &user).Return(nil)
	err := suite.service.Move(suite.ctx)

	suite.Nil(err)
}

func (suite *ListServiceTestSuite) TestBadMoveWithValidationError() {
	list := factory.NewList(&factory.ListConfig{})
	currentUser := factory.NewUser(&factory.UserConfig{})
	suite.ctx.Set(config.CurrentUserKey, currentUser)
	suite.ctx.Set(config.ListKey, list)
	req := httptest.NewRequest("PUT", "/api/lists/1/move", factory.CreateListRequestBody(&factory.ListConfig{Index: -1}))
	suite.ctx.Request = req
	err := suite.service.Move(suite.ctx)

	suite.IsType(validator.ValidationErrors{}, err)
}

func (suite *ListServiceTestSuite) TestBadMoveWithDBError() {
	list := factory.NewList(&factory.ListConfig{})
	currentUser := factory.NewUser(&factory.UserConfig{})
	suite.ctx.Set(config.ListKey, list)
	suite.ctx.Set(config.CurrentUserKey, currentUser)
	toIndex := 1
	req := httptest.NewRequest("PUT", "/api/lists/1/move", factory.CreateListRequestBody(&factory.ListConfig{Index: toIndex}))
	suite.ctx.Request = req
	err := errors.New("db error")
	suite.listRepositoryMock.EXPECT().Move(&list, toIndex, &currentUser).Return(err)
	rerr := suite.service.Move(suite.ctx)

	suite.Equal(err, rerr)
}
